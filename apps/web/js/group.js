var ChatMembers = [];

/**
 * 展示创建群聊模态框
 */
function createGroup(res) {
    if (typeof res !== "undefined") {
        jqtoast('群创建成功');
        ws.send('{"type":"Groups","services":"UserGroupsService","content":""}');
        return
    }

    $("#search-friends-hook").hide();
    $("#user-group-hook").css('display', 'flex');
}

/**
 * 创建群
 */
function sendCreateGroup() {
    const name = $('input[name=group-name]').val();
    const description = $('textarea[name=group-description]').val();
    const avatar = $('input[name=group-avatar]').val();
    if (name === "") {
        jqtoast("群名称不能为空")
        return false;
    }
    if (description === "") {
        jqtoast("群描述不能为空")
        return false;
    }
    if (avatar === "") {
        jqtoast("群头像不能为空")
        return false;
    }

    ws.send('{"type":"CreateGroup","services":"UserGroupsService","content":"{\\"name\\":\\"' + name + '\\",\\"description\\":\\"' + description + '\\",\\"avatar\\":\\"' + avatar + '\\"}"}')

    changeModalStatus('#user-group-hook', 'hide');
}

//修改群信息
function updateGroup() {
    ws.send('{"type":"EditGroup","services":"UserGroupsService","content":"{\\"name\\":\\"修改群\\",\\"description\\":\\"修改群\\",\\"avatar\\":\\"https://p5.toutiaoimg.com/origin/pgc-image/06be98af5dd4491993ac131c9a3410cf\\",\\"gid\\":\\"38de63fd-d1a5-4ccb-9885-334f91ae0bda\\"}"}')

    changeModalStatus('#user-group-hook', 'hide');
}

/**
 * 渲染群列表
 * @param data
 * @constructor
 */
function Groups(data) {
    const groups = data.groups;
    let html = '';
    for (let i in groups) {
        let item = groups[i];
        GROUPS[item['gid']] = {
            avatar: item.avatar,
            username: item.name,
            uuid: item.gid,
            adminId: item.uuid,
            status: 1,
            group: true
        };

        html += `
            <div 
            oncontextmenu="customMenu(event,'${item['gid']}','group')" 
            id="friend-${item['gid']}}" 
            class="friends_box" 
            ondblclick="chat('${item['gid']}','groupMessage')">
                <div class="user_head"><img src="${item['avatar']}" alt=""></div>
                <div class="friends_text">
                    <p class="user_name">${item['name']}</p>
                </div>
            </div>
        `;
    }

    let dom = $('.group_list li');
    dom.children('.friends_box').remove();
    dom.append(html);
}

/**
 * 展开群设置
 * @param data
 * @param e
 */
function groupSetting(gid, e,type) {
    e.stopPropagation();

    if(typeof type === "undefined"){
        //获取群成员
        GroupMembers(gid);
        $('#sidebar-information').show();
        $('#sidebar-tool').fadeIn();
        return false;
    }

    //搜索消息
    $('#sidebar-information').hide();
    $('#sidebar-search').show();
    $('#search-content').hide();
    $('#sidebar-tool').fadeIn();
}

/**
 * 获取群成员
 * @param data
 * @returns {boolean}
 * @constructor
 */
function GroupMembers(data) {
    if (typeof data === "string") {
        //发送获取群成员请求
        ws.send('{"type":"GroupMembers","services":"UserGroupMemberService","content":"' + data + '"}')
        return false;
    }

    ChatMembers = [];
    if (data.length === 0) {
        return false;
    }

    //判断当前用户是当前群的管理员，则添加删除群样式
    let close = '';
    if (GROUPS[CHATInfo.uuid].adminId === USERInfo.Uuid) {
        //是当前群的管理员，增加删除群成员操作
        close = "<i onclick='removeMembers(this)'>删除</i>";
    }

    let users = [];
    let html = '';
    for (let i in data) {
        let item = data[i];
        users.push(item.uuid);
        ChatMembers[item.uuid] = item;
        html += `
            <li class="sidebar-user-item">
                <img src="${item.avatar}" alt="">
                <p data-uuid="${item.uuid}">${item.username}</p>
                ${close}
            </li>
        `;
    }
    $('.sidebar-users .sidebar-user-item').remove();
    $('.sidebar-users .add-member').before(html)

    const groupInfo = GROUPS[CHATInfo.uuid];
    //设置群名
    $('.sidebar-group-name').html(CHATInfo.username + '<span> > </span>')
    //设置群描述
    // $('.sidebar-group-description').text(CHATInfo+'')

    //添加好友元素
    let friendHtml = '';
    for (let i in FRIENDS) {
        let item = FRIENDS[i];
        if (users.indexOf(item.uuid) !== -1) {
            continue;
        }
        friendHtml += `
            <li>
                <p class="select">
                    <input type="checkbox" name="addMember" value="${item.uuid}" onchange="listenAddMember()">
                </p>
                <div class="user-wrapper">
                    <img src="${item.avatar}" alt="">
                    <p>${item.username}</p>
                </div>
            </li>
        `;
    }
    $('.sidebar-users-wrapper ul').html(friendHtml);
}

/**
 * 监听用户选择的添加群的成员
 */
function listenAddMember() {
    let data = getAddMembers();
    $('.push-group-member').text('完成（' + data.length + '）');
}

/**
 * 获取已选的用户
 */
function getAddMembers() {
    let data = [];
    $('input[name=addMember]:checked').each(function () {
        let val = $(this).val();
        data.push(val);
    });

    return data;
}

/**
 * 发送添加成员请求
 */
function addMembers() {
    let data = getAddMembers();
    if (data.length === 0) {
        jqtoast('请先选择需要添加的成员');
        return false;
    }

    for (let i in data) {
        //发送添加成员请求
        let json = {
            "type": "AddMember",
            "service": "UserGroupMemberService",
            "content": JSON.stringify({"uuid": data[i], "gid": CHATInfo.uuid})
        }
        ws.send(JSON.stringify(json))
    }

    //修改页面样式
    $('.sidebar-users-wrapper').fadeOut();
    $('input[name=addMember]').each(function () {
        let val = $(this).attr('checked', false);
    });
    $('.push-group-member').text('完成（0）');
}

/**
 * 添加成员通知
 * @param msg
 * @constructor
 */
function AddMember(msg) {
    jqtoast("邀请成功！");

    setTimeout(function () {
        //重新获取群成员
        const obj = JSON.parse(msg);
        GroupMembers(obj.gid);

        chat(obj.gid, 'groupMessage');
    },1000)
}


/**
 * 监听鼠标点击事件，控制展示对应的元素
 */
$(".sidebar-users-wrapper").click(function (e) {
    e.stopPropagation();
});

$("#sidebar-tool").click(function (e) {
    e.stopPropagation();
});

/**
 * 展示添加群成员dom
 * @param e
 */
function showAddMembersDom(e) {
    e.stopPropagation();
    $('.sidebar-users-wrapper').fadeIn();
}

/**
 * 删除群成员
 * @param el
 */
function removeMembers(el) {
    jqalert({
        title: "删除提示",
        content: "确认要删除此成员吗？此操作不可逆",
        yestext: "确认",
        notext: "算了",
        yesfn: function (event) {
            event.stopPropagation();
            let uuid = $(el).prev().data('uuid');
            //发送添加成员请求
            let json = {
                "type": "RemoveMember",
                "service": "UserGroupMemberService",
                "content": JSON.stringify({"uuid": uuid, "gid": CHATInfo.uuid})
            }
            ws.send(JSON.stringify(json))
            $(el).parent().remove();
        },
        nofn: function (event) {
            event.stopPropagation();
        }
    });
}

//退出群聊
function outGroup(gid) {
    jqalert({
        title: "退出提示",
        content: "确认要退出该群吗？退出后将不在接收此群消息",
        yestext: "确认",
        notext: "算了",
        yesfn: function (event) {
            event.stopPropagation();

            //发送添加成员请求
            let json = {
                "type": "OutGroup",
                "service": "UserGroupMemberService",
                "content": gid
            }
            ws.send(JSON.stringify(json))
        },
        nofn: function (event) {
            event.stopPropagation();
        }
    });
}

/**
 * 退出群通知
 * @param msg
 * @constructor
 */
function OutGroup(msg) {
    // jqtoast("退出群聊成功");
    //获取群列表
    ws.send('{"type":"Groups","services":"UserGroupsService","content":""}');
    //重新渲染聊天记录列表
    ws.send('{"type":"List","services":"HistoryRecordService","content":""}');
}

//删除群
function delGroup(gid) {
    jqalert({
        title: "解散群提示",
        content: "确认要解散该群吗？",
        yestext: "确认",
        notext: "算了",
        yesfn: function (event) {
            event.stopPropagation();

            //发送添加成员请求
            let json = {
                "type": "DelGroup",
                "service": "UserGroupsService",
                "content": gid
            }
            ws.send(JSON.stringify(json))
        },
        nofn: function (event) {
            event.stopPropagation();
        }
    });
}

function DelGroup(msg) {
    jqtoast("解散群聊成功");
}