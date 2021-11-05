/**
 * 搜索消息
 */
function searchMessage() {
    let keywords = $('input[name=message-search]').val();
    if (!keywords) {
        jqtoast('请输入需要搜索的内容！')
        return false;
    }

    const toId = CHATInfo.uuid;

    //判断是否为群聊
    let isGroup = false;
    if (CHATInfo.hasOwnProperty('group') && CHATInfo['group']) {
        isGroup = true
        console.log(122);
    }

    let data = {
        "type": "Search",
        "services": "messageSearchService",
        "content": JSON.stringify({
            "keywords": keywords,
            "to_id": toId,
            "is_group": isGroup
        })
    };

    console.log(data);

    ws.send(JSON.stringify(data))
}

/**
 * 处理搜索结果
 * @param result
 */
function searchResponse(result) {

    //设置搜索条数
    let list = JSON.parse(result.result);
    if (list.length == 0) {
        jqtoast('未搜索到数据');
        return false;
    }

    let members = FRIENDS;
    //将自己的信息加入到好友列表中
    members[USERInfo.Uuid] = {
        username: USERInfo.Username,
        avatar: USERInfo.Avatar,
        uuid: USERInfo.Uuid
    };

    //判断是否为群聊
    if (CHATInfo.hasOwnProperty('group') && CHATInfo['group']) {
        members = ChatMembers;
    }
    $('#search-content').show();
    $('.search-total').text(result.total + '条相关记录');
    let html = '';
    for (let i in list) {
        let t = (new Date(list[i].create_time)).getTime() * 1000000;
        let d = formatDateByTimeStamp(t);
        html += `
            <li>
                <p>${members[list[i]['from_id']].username} <span>${d}</span></p>
                <div>${list[i].content}</div>
            </li>
        `;
    }
    $("#search-content ul").html(html);
}