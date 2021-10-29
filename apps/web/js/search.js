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
    let data = {
        "type": "Search",
        "services": "messageSearchService",
        "content": JSON.stringify({
            "keywords": keywords,
            "to_id": toId
        })
    };

    ws.send(JSON.stringify(data))
}