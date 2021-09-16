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
            ondblclick="chat('${item['gid']}','groupChat')">
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