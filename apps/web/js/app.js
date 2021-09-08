/**
 * ajax post异步请求
 * @param url
 * @param data
 * @param callback
 * @constructor
 */
function AjaxPost(url, data, callback) {

    //数据容错
    if (data && typeof data !== "string") {
        data = JSON.stringify(data);
    }

    const token = getCookie('token');
    $.ajax(
        {
            url: API + url,
            'data': data,
            'type': 'POST',
            'processData': false,
            'dataType': 'json',
            'contentType': 'application/json',
            'async':true,
            headers:{
                Authorization:token
            },
            success: function (json) {
                callback(json);
            },
            error: function (err) {
                jqtoast(err);
            }
        }
    );
}

function AjaxGet(url, data, callback) {

    //数据容错
    if (data && typeof data !== "string") {
        data = JSON.stringify(data);
    }

    const token = getCookie('token');
    $.ajax(
        {
            url: API + url,
            'data': data,
            'type': 'GET',
            'processData': false,
            'dataType': 'json',
            'contentType': 'application/json',
            'async':true,
            headers:{
                Authorization:token
            },
            success: function (json) {
                callback(json);
            },
            error: function (err) {
                jqtoast(err);
            }
        }
    );
}

/**
 * 发送消息
 * @param url
 * @param data
 * @param callback
 * @constructor
 */
function AjaxMsg(url, data, callback,method) {

    //数据容错
    if (data && typeof data !== "string") {
        data = JSON.stringify(data);
    }

    if(typeof method === "undefined"){
        method = 'POST';
    }

    const token = getCookie('token');
    $.ajax(
        {
            url: MSGAPI + url,
            'data': data,
            'type': method,
            'processData': false,
            'dataType': 'json',
            'contentType': 'application/json',
            'async':true,
            headers:{
                Authorization:token
            },
            success: function (json) {
                callback(json);
            },
            error: function (err) {
                jqtoast(err);
            }
        }
    );
}

/**
 * 设置cookie
 * @param cname
 * @param value
 * @param expire
 */
function setCookie(cname, value, expire) {
    if (typeof expire === "undefined") {
        expire = 24 * 60 * 60 * 30;
    }

    let d = new Date();
    expire = d.getTime() + (expire * 1000);
    let expires = "expires=" + expire;
    document.cookie = cname + "=" + value + "; " + expires;
}

/**
 * 获取cookie
 * @param cname
 * @returns {string}
 */
function getCookie(cname) {
    let name = cname + "=";
    let ca = document.cookie.split(';');
    for (let i = 0; i < ca.length; i++) {
        let c = ca[i].trim();
        if (c.indexOf(name) == 0) {
            return c.substring(name.length, c.length);
        }
    }

    return "";
}

/**
 * 删除cookie
 * @param name
 */
function delCookie(name)//删除cookie
{
    document.cookie = name + "=;expires=";
}

/**
 * 切换模态框状态
 * @param el
 * @param status
 */
function changeModalStatus(el,status){
    if(status === 'show'){
        $(el).css('display','flex');
    }else {
        $(el).css('display','none');
    }
}