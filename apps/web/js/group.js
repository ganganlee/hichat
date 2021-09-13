/**
 * 展示创建群聊模态框
 */
function createGroup(res) {
    if(typeof res !== "undefined"){
        jqtoast('群创建成功');
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

function updateGroup() {
    ws.send('{"type":"EditGroup","service":"UserGroupsService","content":"{\\"name\\":\\"修改群\\",\\"description\\":\\"修改群\\",\\"avatar\\":\\"https://p5.toutiaoimg.com/origin/pgc-image/06be98af5dd4491993ac131c9a3410cf\\",\\"gid\\":\\"38de63fd-d1a5-4ccb-9885-334f91ae0bda\\"}"}')

    changeModalStatus('#user-group-hook', 'hide');
}

//获取群列表
function groupList() {
    AjaxGet("/v1/group/list", "", (json) => {

        if (json.code != 200) {
            jqtoast(json.msg);
            return false;
        }

        const data = json.result;
        let html = '';
        for (let i in data) {
            let info = data[i];
            info['id'] = 0;

            FRIENDS[info['token']] = info;
            html += `
                <div 
                oncontextmenu="customMenu(event,'${info['token']}','friend')" 
                id="friend-${info['token']}}" 
                class="friends_box" 
                ondblclick="chat('${info['token']}',${ChatType.group})">
                    <div class="user_head"><img src="${info['head_img']}" alt=""></div>
                    <div class="friends_text">
                        <p class="user_name">${info['username']}</p>
                    </div>
                </div>
            `;
        }

        //存在数据，往好友列表中添加
        if (html !== '') {
            html = `
                <li>
                    <p>群聊</p>
                    ${html}
                </li>
            `;

            $('.friends_list').prepend(html);
        }
    })
}