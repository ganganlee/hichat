//当前登录用户
let USERInfo = {};

//当前聊天用户
let CHATInfo = {};

//好友列表
let FRIENDS = [];
//聊天列表
let HistoryList = [];
//音频对象
let messageAudio;

let ws;

$(function () {
    //判断是否登录，未登录的情况下跳转登录
    TOKEN = getCookie('token');

    if (!TOKEN) {
        window.location.href = '/login.html';
    }

    render();

    //获取音乐对象
    messageAudio = document.getElementById('messageAudio');
});

//监听搜索用户输入框，当输入框按下回车时进行搜索
$('input[name=searchUser]').bind('keydown', function (event) {
    if (event.keyCode == '13') {
        searchUser();
    }
});

/**
 * 页面初始化方法
 */
function render() {
    //获取用户信息
    renderUserInfo();

    //渲染emoji表情
    initEmoji();
}

/**
 * 渲染好友列表
 */
function friends(friends) {

        let friendList = {
            '*': [],
            'A': [],
            'B': [],
            'C': [],
            'D': [],
            'E': [],
            'F': [],
            'G': [],
            'H': [],
            'I': [],
            'J': [],
            'K': [],
            'L': [],
            'M': [],
            'N': [],
            'O': [],
            'P': [],
            'Q': [],
            'R': [],
            'S': [],
            'T': [],
            'U': [],
            'V': [],
            'W': [],
            'X': [],
            'Y': [],
            'Z': [],
            '#': [],
        };

        //循环列表，渲染好友数据
        for (let i in friends) {
            let friend = friends[i];
            FRIENDS[friend['uuid']] = friend;
            let html = `
                <div 
                oncontextmenu="customMenu(event,'${friend['uuid']}','friend')" 
                id="friend-${friend['uuid']}}" 
                class="friends_box" 
                ondblclick="chat('${friend['uuid']}')">
                    <div class="user_head"><img src="${friend['avatar']}" alt=""></div>
                    <div class="friends_text">
                        <p class="user_name">${friend['username']}</p>
                    </div>
                </div>
            `;
            let py = getPy(friend['username']);
            if (!friend.hasOwnProperty('status') || friend.status === 0) {
                py = ['*'];
                html = `
                    <div 
                    oncontextmenu="friendApprove('${friend['uuid']}','${friend.username}',this)" 
                    id="friend-${friend['uuid']}}" 
                    class="friends_box" 
                    onclick="friendApprove('${friend['uuid']}','${friend.username}',this)">
                        <div class="user_head"><img src="${friend['avatar']}"></div>
                        <div class="friends_text">
                            <p class="user_name">${friend['username']}</p>
                        </div>
                    </div>
                `;
            }

            friendList[py[0]].push(html);
        }

        let html = '';
        for (let i in friendList) {

            //判断数据是否存在
            if (friendList[i].length === 0) {
                continue;
            }

            let title = i;
            if (i === '*') {
                title = '好友申请';
            }

            html += `
                <li>
                    <p>${title}</p>
                    ${friendList[i].join('')}
                </li>
            `;
        }

        //渲染好友列表
        $('.friends_list').html(html);

        //判断是否存在好友申请，存在时添加角标
        if (friendList['*'].length > 0) {
            let l = friendList['*'].length;
            $('#si_2 span').text(l);
            $('#si_2 span').show();
        }else {
            $('#si_2 span').hide();
        }

        //获取群列表
        // groupList();

        //获取历史聊天列表
        // renderHistory();
}

/**
 * 开始聊天
 * @param userInfo
 */
function chat(token,msgType) {

    //设置聊天类型，1 私聊，2 群聊
    if(typeof msgType === "undefined"){
        msgType = ChatType.chat;
    }

    $('#send').data('chatType',msgType);

    //修改title
    $('title').text('微聊');

    //保存聊天对象
    CHATInfo = FRIENDS[token];

    //判断当前用户是否在聊天列表，不在聊天列表则添加
    if (!HistoryList[token]) {
        let res = {
            date: (new Date().getTime()) * 1000000,
            msg: "",
            token: token,
            unread: 0
        };

        HistoryList[token] = res;
        AppendHistoryHtml(res);
    }

    //设置聊天dom的用户信息
    //设置昵称
    $('#message-user').text(CHATInfo.username);

    //情况聊天记录
    $('#chat-wrapper').html('');

    //去除未读消息
    const dom = $(`.history-${token} .unread-message`);
    dom.text(0);
    dom.hide();

    //增加当前列表聊天样式
    $('.user_list li').removeClass('user_active');
    $(`.history-${token}`).addClass('user_active');

    //获取聊天记录
    AjaxMsg('/v1/msg/history', {
        "from_token": USERInfo.token,
        "from_user_id": USERInfo.id,
        "friend_token": CHATInfo.token,
        "friend_user_id": CHATInfo.id,
        "page": 1,
        "page_size": 20
    }, function (json) {
        if (json.code !== 200) {
            jqtoast(json.msg);
            return;
        }
        if (!json.result) {
            return;
        }

        const data = json.result;
        for (let i = data.length - 1; i >= 0; i--) {
            let message = JSON.parse(data[i].message);

            //判断接收者
            let position = 'other';
            let headImg = CHATInfo.head_img;
            let nickname = CHATInfo.username;
            if (message.receive_token === CHATInfo.token) {
                position = 'me';
                headImg = USERInfo.head_img;
                nickname = USERInfo.username;
            }

            //渲染消息
            fillingMsg(message.type, message.content, headImg, nickname, position);
        }

        //dom滚动至底部
        scrollToFooter('#chat-wrapper');
    });
}

/**
 * 连接消息服务器
 */
function websocket() {
    ws = new WebSocket(WEBSOCKETAPI + "/v1/listen?token=" + TOKEN);

    //连接打开时触发
    ws.onopen = function (evt) {
        console.log("Connection open ...");
        messageAudio.muted = false;

        //获取好友列表
        ws.send('{"type":"Friends","service":"UserService"}');
    };

    //接收到消息时触发
    ws.onmessage = function (evt) {
        let data = JSON.parse(evt.data);
        console.log(data);
        //数据操作是失败！
        if(data.type === 'err'){
            jqtoast(data.result);
            return;
        }

        //数据操作成功提示，没有其他操作
        if(data.type === 'success'){
            jqtoast(data.result)
            return;
        }

        eval(data.type+'(data.result)');
        return;

        $('iframe').remove();
        //音频提示
        //messageAudio.play();
        let iframe = document.createElement('iframe');
        iframe.src = "/media/message.mp3";
        document.body.appendChild(iframe);

        //判断当前用户是否在聊天列表，不在聊天列表则添加
        if (!HistoryList[data.token]) {
            let res = {
                date: (new Date().getTime()) * 1000000,
                msg: data.body.content,
                token: data.token,
                unread: 1
            };
            HistoryList[data.token] = res;
            AppendHistoryHtml(res);
            //增加未读消息
            pushUnread(data.token);
            return;
        }

        //判断当前的聊天对象是不是接收到消息的用户，如果不是，则增加角标
        if (CHATInfo.token !== data.token) {
            setUnreadMessage(data.token, data.body.type, data.body.content);

            //增加未读消息
            pushUnread(data.token);
            return;
        }

        fillingMsg(data.body.type, data.body.content, data.user.head_img, data.user.username, 'other')

        //同步左边聊天列表
        setUnreadMessage(data.token, data.body.type, data.body.content, true)

        //dom滚动至底部
        scrollToFooter('#chat-wrapper');
    };

    //连接关闭时触发
    ws.onclose = function (evt) {
        jqtoast('服务连接失败，正在重试！');
        setTimeout(function () {
            websocket();
        }, 1000);

    };
}

/**
 * 渲染聊天消息
 * @param token
 * @param type
 * @param content
 * @param clear
 */
function setUnreadMessage(token, type, content, clear) {

    let dom = $(`.history-${token} .unread-message`);

    //获取房前未读消息，
    let num = dom.text();
    if (!num) {
        num = 0;
    }

    //将数据转为int,方便运算
    num = parseInt(num) + 1;

    //当clear存在时，说明当前就是此用户再聊天，需要去除未读消息
    if (typeof clear !== 'undefined') {
        num = 0;
    }
    dom.text(num);

    if (num > 0) {
        dom.show();
    } else {
        dom.hide();
    }

    //将消息渲染到列表
    let message = '';
    switch (type) {
        case 2:
            message = "[图片]"
            break
        case 3:
            message = "[音频]"
            break
        case 4:
            message = "[视频]"
            break
        case 5:
            message = "[文件]"
            break
        default:
            message = content
    }

    $(`.history-${token} .user_message`).html(message);

    //修改时间
    let date = (new Date().getTime()) * 1000000;
    let d = formatDateByTimeStamp(date)
    $(`.history-${token} .user_time`).text(d);
}

/**
 * 填充消息内容
 * @param type
 * @param content
 * @param headImg
 * @param nickname
 * @param position
 */
function fillingMsg(type, content, headImg, nickname, position) {

    //判断type不等于1时说明消息为非文本消息，需要将消息解析对对象
    if (type !== 1) {
        content = JSON.parse(content);
        switch (type) {
            case 2://图片
                content = `<img src="${UPLOADAPI + content.path}" alt="${content.name}" style="max-width: 200px"/>`;

                break;
            case 3://音频
                content = `
                    <audio src="${UPLOADAPI + content.path}" controls="controls" style="width: 200px">
                        您的浏览器不支持 audio 标签。
                    </audio>
                `;
                break;
            case 4://视频
                content = `
                    <video src="${UPLOADAPI + content.path}" controls="controls" style="max-width: 200px">
                        您的浏览器不支持
                    </video>
                `;
                break;
            default://文件
                content = `
                    <a href="${UPLOADAPI + content.path}" download>${content.name}</a>
                `;
                break;
        }
    }

    if (typeof position === "undefined") {
        position = 'me';
    }

    let msg = `
        <li class="${position}">
            <img src="${headImg}" title="${nickname}">
            <span>${content}</span>
        </li>
    `;

    $('#chat-wrapper').append(msg);
}

/**
 * 发送消息
 * @param type
 * @returns {boolean}
 */
function sendMsg(type) {

    //关闭emoji模态框
    $('.emoji-wrapper').css('display', 'none');

    //给type设置默认值
    if (typeof type === "undefined") {
        type = 1;
    }

    let msg = $('#input_box').html();
    if (msg === '') {
        return false;
    }

    let content = {
        content: msg,
        to_token: CHATInfo.token,
        type: type,
        receive_token: CHATInfo.token
    };

    //获取聊天类型
    let chatType = $('#send').data('chatType');

    //定义聊天类型
    let data ={
        type : $('#send').data('chatType'),
        content : JSON.stringify(content)
    };

    //发送数据
    ws.send(JSON.stringify(data));

    //清空输入框
    setTimeout(function (){
        $('#input_box').html('');
        $('#input_box div').remove();
    },100)

    //添加样式
    fillingMsg(type, msg, USERInfo['head_img'], USERInfo['username']);

    //将消息渲染到聊天列表
    setUnreadMessage(CHATInfo.token, type, msg, true);

    //dom滚动至底部
    scrollToFooter('#chat-wrapper');
}

/**
 * 监听回车键，回车键直接发送数据
 * @param event
 */
function listenKeyDown(event) {
    if (event.keyCode === 13) {
        sendMsg();
    }
}

/**
 * 获取聊天记录
 */
function renderHistory() {
    AjaxMsg('/v1/msg/list', '', function (json) {
        if (json.code !== 200) {
            jqtoast(json.msg);
            return;
        }

        let data = json.result;
        if (!data) {
            return false;
        }

        //按照时间戳排序
        data.sort(sortByField("date"))

        $('.user_list').html('');
        for (let i in data) {
            let user = data[i];
            //将数据保存进入全局变量中
            HistoryList[user['token']] = user;

            AppendHistoryHtml(user);
        }

        //将消息的第一个渲染为聊天对象
        setTimeout(function () {
            chat(`${data[data.length - 1]['token']}`);
        }, 500)
    }, 'GET');
}

/**
 * 向聊天列表中添加html
 * @param user
 * @constructor
 */
function AppendHistoryHtml(user) {
    //获取时间
    const d = formatDateByTimeStamp(user.date)
    console.log(user);

    //判断未读消息数大于0时，显示未读消息
    let unreadStatus = 'hidden';
    if (user.unread > 0) {
        unreadStatus = 'block';
    }

    let html = `
        <li 
        class="history-${user['token']}" 
        ondblclick="chat('${user['token']}')"
        oncontextmenu="customMenu(event,'${user['token']}','history')">
            <div class="user_head">
                <img src="${FRIENDS[user['token']]['head_img']}" alt="">
                <span class="unread-message" style="display: ${unreadStatus}">${user.unread}</span>
            </div>
            <div class="user_text">
                <p class="user_name">${FRIENDS[user['token']]['username']}</p>
                <p class="user_message">${user.msg}</p>
            </div>
            <div class="user_time">${d}</div>
        </li>
    `;

    //填充数据
    $('.user_list').prepend(html);
}

/**
 * 数组排序
 * @param field
 * @returns {(function(*, *): (number))|*}
 */
function sortByField(field) {
    return function (obj1, obj2) {
        let val1 = obj1[field];
        let val2 = obj2[field];
        if (val1 < val2) {
            return -1;
        } else if (val1 > val2) {
            return 1;
        } else {
            return 0;
        }
    }
}

/**
 * 获取用户信息
 */
function renderUserInfo() {
    AjaxGet('/v1/user/info', '', function (json) {
        if (json.code !== 200) {
            jqtoast(json.msg);
            return;
        }

        USERInfo = json.result;

        console.log(USERInfo);

        //设置用户的连接属性
        WEBSOCKETAPI = 'ws://' + USERInfo['Host'];

        //设置用户信息
        $('.own_head').attr('src', USERInfo['Avatar']);
        $('.own_name').text(USERInfo['Username']);
        $('.own_numb').text('ID：' + USERInfo['uuid']);
        $('.user-head-img img').attr('src', USERInfo['Avatar']);

        //连接websocket
        websocket();

    });
}

/**
 * 搜索用户
 */
function searchUser() {
    let val = $('input[name=searchUser]').val();
    if (!val) {
        $('#').modal('hide');
        changeModalStatus('#search-friends-hook', 'hide');
        return;
    }

    ws.send('{"type":"FindByName","service":"UserService","content":"'+val+'"}');
}

/**
 * 渲染用户搜索
 * @param users
 */
function findUser(users) {
    //判断是否存在用户
    if (users === null) {
        jqtoast('用户不存在');
        return;
    }

    let html = '';
    for (let i in users) {
        let user = users[i];
        let status = '未添加';
        let click = 'selectUser(this)';
        //判断当前用户是否为自己
        if (user.uuid === USERInfo.uuid) {
            status = '自己';
            click = '';
        }

        //判断当前用户是否已经是好友
        if (FRIENDS[user.uuid]) {
            status = '已添加';
            click = '';
        }

        html += `
            <div data-token="${user.uuid}" onclick="${click}" class="search-user-hook" style="padding: 15px;border: 1px solid #ccc;border-radius: 5px;text-align: center;margin: 15px;">
                <img style="width: 100px;height: 100px;object-fit: cover;border-radius: 100px;" src="${user.avatar}" alt="">
                <p>${user.username}</p>
                <p>${status}</p>
            </div>
        `;
    }

    //渲染列表
    $('#friends-hook .model-content').html(html);

    //展示模态框
    changeModalStatus('#friends-hook', 'show')
    changeModalStatus('#search-friends-hook', 'hide')
}

/**
 * 选择需要添加的好友
 */
function selectUser(el) {
    $(el).css('border-color', 'red');
    $(el).addClass('active');
}

/**
 * 添加好友
 */
function addUser() {
    if (!$('.search-user-hook').hasClass('active')) {
        jqtoast('请先选择需要添加的用户', 3000);
        return;
    }

    const token = $('#friends-hook .active').data('token');
    ws.send('{"type":"ApplyFriend","service":"UserService","content":"'+token +'"}');
    changeModalStatus('#friends-hook', 'hide');
}

/**
 * 展示好友申请
 * @param token
 * @param username
 * @param el
 */
function friendApprove(token, username, el) {

    jqalert({
        title: "好友申请",
        content: username + '请求添加你为好友，是否同意',
        yestext: '同意',
        notext: '拒绝',
        yesfn: function () {
            //同意
            ws.send('{"type":"ApproveFriend","service":"UserService","content":"'+token+'"}');
            //增加提示
            setTimeout(function (){
                ws.send('{"type":"Friends","service":"UserService"}')
            },1500);
            return false;
        },
        nofn: function () {
            //拒绝
            ws.send('{"type":"RefuseFriend","service":"UserService","content":"'+token+'"}');
            //增加提示
            setTimeout(function (){
                ws.send('{"type":"Friends","service":"UserService"}')
            },1500);
        }
    });

    document.oncontextmenu = function(e){
        return false;
    };
}

/**
 * 设置聊天框中的滚动条始终在底部
 */
function scrollToFooter(el) {
    let parentHeight = $('#chat-container').height();
    let height = $(el).height() + 20;
    if (height > parentHeight) {
        $(el).animate({
            top: parentHeight - height
        }, 500);
    }
}

/**
 * 文件上传
 */
function upload(id, func) {
    let formData = new FormData();
    formData.append("file", $(id)[0].files[0]);
    const token = getCookie('token');
    add_load();
    $.ajax({
        url: UPLOADAPI + '/v1/upload', /*接口域名地址*/
        type: 'post',
        data: formData,
        contentType: false,
        processData: false,
        dataType: 'json',
        headers: {
            Authorization: token
        },
        success: function (json) {
            remove_load();
            if (json.code !== 200) {
                jqtoast(json.msg);
                return false;
            }

            //文件上传成功，渲染图片
            func(json.result)
        },
        error: function (res) {
            remove_load();
            console.log(res);
        }
    })
}

/**
 * 上传头像
 * @param data
 */
function handleHeadImg(data) {
    const img = UPLOADAPI + data.path;
    $('input[name=update-head-img]').val(img);
    $('.show-user-img').attr('src', img)
}

/**
 * 上传群头像
 * @param data
 */
function handleGroupAvatar(data) {
    const img = UPLOADAPI + data.path;
    $('input[name=group-avatar]').val(img);
    $('.group-avatar').attr('src', img);
    $('.group-avatar').show();
}

/**
 * 修改用户信息
 */
function updateUserInfo() {
    let username = $('input[name=update-username]').val();
    let password = $('input[name=update-password]').val();
    let head_img = $('input[name=update-head-img]').val();
    if (!username && !password && !head_img) {
        jqtoast('未设置任何参数');
        return false;
    }
    AjaxPost('/v1/user/updateInfo', {
        "username": username,
        "password": password,
        "head_img": head_img
    }, function (json) {
        console.log(json);
        if (json.code !== 200) {
            jqtoast(json.msg);
            return false;
        }

        jqtoast('操作成功！')
        //关闭模态框
        changeModalStatus('#update-info-hook', 'hide')
        //重新获取用户数据
        renderUserInfo();
    });
}

/**
 * 展示用户信息
 */
function showUserInfo(event) {
    const model = $('.user-head-img-model');
    if (model.css('display') === 'none') {
        model.show();
        $('input[name=update-username]').val(USERInfo.username);
        $('input[name=update-head-img]').val(USERInfo.head_img);
        $('.show-user-img').attr('src', USERInfo['head_img']);
    } else {
        model.hide();
    }

    event.stopPropagation();
}

/**
 * 聊天窗口上传文件
 * @param json
 */
function handleChatUpload(data) {
    //文件上传成功，组织发送的数据
    let message = {
        "path": data.path,
        "name": data.name
    }
    message = JSON.stringify(message);
    $('#input_box').text(message);
    sendMsg(TYPES[data.type]);
}

/**
 * 初始化emoji表情
 */
function initEmoji() {
    let html = '';
    for (let i in EMOJIS) {
        html += `
            <p class="emoji-item" onclick="pushEmoji(${EMOJIS[i]})">
                <img src="images/emoji/${EMOJIS[i]}.gif" alt="">
            </p>
        `;
    }

    $('.emoji-wrapper').html(html);
}

/**
 * 切换emoji表情模态框状态
 */
function selectEmojiStatus() {
    let status = $('.emoji-wrapper').css('display');
    if (status === 'flex') {
        $('.emoji-wrapper').css('display', 'none');
    } else {
        $('.emoji-wrapper').css('display', 'flex');
    }
}

/**
 * 往输入框插入表情
 * @param index
 */
function pushEmoji(index) {
    $('#input_box').append(`<img src="images/emoji/${index}.gif" />`);
}

/**
 * 删除好友
 * @param friendToken
 */
function delFriend(friendToken) {
    AjaxGet("/v1/friends/delFriend/" + friendToken, "", (json) => {
        if (json.code === 200) {
            jqtoast('删除成功');
            $(`#friend-${friendToken}`).remove();
            $(`.history-${friendToken}`).remove();
            return false;
        }

        jqtoast(json.msg);
    });
}

/**
 * 自定义鼠标右键菜单
 * @param event
 * @param friendToken
 */
function customMenu(event, friendToken, origin) {

    //关闭鼠标右键
    event.preventDefault();

    // 显示自定义的菜单调整位置
    let scrollTop = document.documentElement.scrollTop || document.body.scrollTop;// 获取垂直滚动条位置
    let scrollLeft = document.documentElement.scrollLeft || document.body.scrollLeft;// 获取水平滚动条位置
    const dom = $('#cursor');

    if (origin === 'friend') {
        //点击好友列表
        $('#cursor-del-friend').attr('onclick', `delFriend('${friendToken}')`);
        $('#cursor-del-friend').show();
        $('#cursor-del-history').hide();
    } else {
        //点击消息列表
        $('#cursor-del-friend').hide();
        $('#cursor-del-history').show();
        $('#cursor-del-history').attr('onclick', `delHistory('${friendToken}')`);
    }

    //展示自定义菜单
    dom.show();
    dom.css({left: event.clientX + scrollLeft + 10 + 'px', top: event.clientY + scrollTop + 'px'});

    $('#cursor-send').attr('onclick', `chat('${friendToken}')`);
}

/**
 * 删除聊天记录
 * @param friendToken
 */
function delHistory(friendToken) {
    AjaxMsg('/v1/msg/list/' + friendToken, '', (json) => {
        if (json.code !== 200) {
            jqtoast(json.msg);
            return false;
        }

        jqtoast('删除成功');
        $(`.history-${friendToken}`).remove();
        delete HistoryList[friendToken];
    }, "GET")
}

// 鼠标点击其他位置时隐藏菜单
let userHeadImgModel = $(".user-head-img-model");
document.onclick = function (event) {
    $('#cursor').hide();
    userHeadImgModel.hide();
}
userHeadImgModel.bind("click", function (event) {
    event = event || window.event;
    event.stopPropagation();
});

/**
 * 增加未读消息
 * @param token
 */
function pushUnread(token) {
    //修改网页title
    $('title').text('您有新消息 - 微聊');

    AjaxMsg('/v1/msg/unread/' + token, '', (json) => {
        console.log(json);
    }, 'GET');
}

/**
 * 退出登录
 */
function logout() {
    //关闭socket连接
    console.log(ws);
    ws.close();
    //删除token
    delCookie('token');
    //将页面跳转至登录界面
    window.location.href = '/login.html';
}