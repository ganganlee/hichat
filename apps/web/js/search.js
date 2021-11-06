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

    ws.send(JSON.stringify(data))
}

/**
 * 处理搜索结果
 * @param result
 */
function searchResponse(result) {

    //设置搜索条数
    let list = JSON.parse(result.result);
    $('#search-content').show();
    if (list.length == 0) {
        let keywords = $('input[name=message-search]').val();
        $('#search-total').text(`没有找到"${keywords}"相关结果`);
        $("#search-content ul").html('');

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

    $('#search-total').text(result.total + '条相关记录');
    let html = '';
    for (let i in list) {
        console.log(list[i]);
        let t = (new Date(list[i].create_time)).getTime() * 1000000;
        let d = formatDateByTimeStamp(t);
        let content;
        switch (list[i].content_type) {
            case 'mp3':
                content = JSON.parse(list[i].content);
                content = `
                    <div>
                        <p>${content.name}</p>
                        <audio style="width: 100%" controls src="${UPLOADAPI}${content.path}"></audio>
                    </div>
                `;
                break;
            case 'mp4':
                content = JSON.parse(list[i].content);
                content = `
                    <div>
                        <p>${content.name}</p>
                        <video style="width: 100%" controls src="${UPLOADAPI}${content.path}"></video>
                    </div>
                `;
                break;
            case 'img':
                content = JSON.parse(list[i].content);
                content = `
                    <div>
                        <p>${content.name}</p>
                        <img src="${UPLOADAPI}${content.path}" alt="" style="max-width: 100%">
                    </div>
                `;
                break;
            case 'file':
                content = JSON.parse(list[i].content);
                content = `
                    <div>
                        <p><a href="${UPLOADAPI}${content.path}">${content.name}</a></p>
                    </div>
                `;
                break;
            default:
                content = list[i].content;
                break;
        }

        html += `
            <li>
                <p>${members[list[i]['from_id']].username} <span>${d}</span></p>
                ${content}
            </li>
        `;
    }
    $("#search-content ul").html(html);
}