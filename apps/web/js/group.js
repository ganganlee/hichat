/**
 * 展示创建群聊模态框
 */
function createGroup(res) {
    if(typeof res !== "undefined"){
        jqtoast('群创建成功');
        ws.send('{"type":"Groups","service":"UserGroupsService","content":""}');
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
    if(name === ""){
        jqtoast("群名称不能为空")
        return false;
    }
    if(description === ""){
        jqtoast("群描述不能为空")
        return false;
    }
    if(avatar === ""){
        jqtoast("群头像不能为空")
        return false;
    }

    ws.send('{"type":"CreateGroup","service":"UserGroupsService","content":"{\\"name\\":\\"'+name+'\\",\\"description\\":\\"'+description+'\\",\\"avatar\\":\\"'+avatar+'\\"}"}')

    changeModalStatus('#user-group-hook', 'hide');
}

//修改群信息
function updateGroup() {
    ws.send('{"type":"EditGroup","service":"UserGroupsService","content":"{\\"name\\":\\"修改群\\",\\"description\\":\\"修改群\\",\\"avatar\\":\\"https://p5.toutiaoimg.com/origin/pgc-image/06be98af5dd4491993ac131c9a3410cf\\",\\"gid\\":\\"38de63fd-d1a5-4ccb-9885-334f91ae0bda\\"}"}')

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
    for(let i in groups){
        let item = groups[i];
        GROUPS[item['gid']] = {
            avatar: item.avatar,
            username: item.name,
            uuid: item.gid,
            status :1,
            group : true
        };

        html += `
            <div 
            oncontextmenu="customMenu(event,'${item['gid']}','friend')" 
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
 * 获取群成员
 * @param data
 * @returns {boolean}
 * @constructor
 */
function GroupMembers(data){
    if (typeof data === "string"){
        //发送获取群成员请求
        ws.send('{"type":"GroupMembers","service":"UserGroupMemberService","content":"'+data+'"}')
        $('#sidebar-tool').fadeIn();
        return false;
    }

    if(data.length === 0){
        return  false;
    }

    console.log(data);
    let users = [];
    let html = '';
    for(let i in data){
        let item  = data[i];
        users.push(item.uuid);
        html += `
            <li class="sidebar-user-item">
                <img src="${item.avatar}" alt="">
                <p>${item.username}</p>
            </li>
        `;
    }
    $('.sidebar-users .add-member').before(html)

    const groupInfo = GROUPS[CHATInfo.uuid];
    //设置群名
    $('.sidebar-group-name').html(CHATInfo.username+'<span> > </span>')
    //设置群描述
    // $('.sidebar-group-description').text(CHATInfo+'')

    //添加好友元素
    let friendHtml = '';
    for(let i in FRIENDS){
        let item = FRIENDS[i];
        if(users.indexOf(item.uuid) !== -1){
            continue;
        }
        friendHtml += `
            <li>
                <p class="select">
                    <input type="checkbox" value="${item.uuid}">
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

