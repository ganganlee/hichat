function HistoryRecord(data){
    //过滤数据为空的列表
    let keys = Object.keys(data);
    if(keys.length === 0){
        return false;
    }

    let list = [];
    for(let i in data){
        list.push(JSON.parse(data[i]));
    }

    //按照时间戳排序
    list.sort(sortByField("date"))

    $('.user_list').html('');
    for (let i in list) {
        let item = list[i];

        let avatar,username,info;
        switch (item.message_type) {
            case 'groupMessage':
                info = GROUPS[item.id]
                console.log(item.id);
                console.log(info);
                console.log(GROUPS);

                username = info.username;
                avatar  = info.avatar;
                break;
            case 'privateMessage':
                info = FRIENDS[item.id]
                username = info.username;
                avatar  = info.avatar;
                break;
        }

        let res = {
            date: (new Date(item.date).getTime()) * 1000000,
            msg: item.content,
            token: item.id,
            unread: 0,
            avatar:avatar,
            username:username,
            message_type:item.message_type
        };

        //将数据保存进入全局变量中
        HistoryList[item['id']] = res;

        AppendHistoryHtml(res);
    }

    //将消息的第一个渲染为聊天对象
    // setTimeout(function () {
    //     chat(`${data[data.length - 1]['token']}`);
    // }, 500)
}