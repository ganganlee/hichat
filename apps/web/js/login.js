$(function () {
    let token = getCookie('token');
    if (token) {
        window.location.href = '/';
    }
});

/**
 * 用户注册
 */
function register() {
    //获取用户名
    let username = $('input[name=username]').val();
    //获取密码
    let password = $('input[name=password]').val();
    let confirmPassword = $('input[name=confirmPassword]').val();

    if (!username || !password) {
        jqtoast("用户名或密码不能为空！");
        return;
    }

    if (password !== confirmPassword) {
        jqtoast("用户名或密码不能为空！");
        return;
    }

    AjaxPost('/v1/user/register', {
        'username': username,
        'password': password,
        'head_img': 'https://p5.toutiaoimg.com/origin/pgc-image/04d64e2b4d094debbd5aa3025b3dbce7'
    }, function (json) {
        if (json.code !== 200) {
            jqtoast(json.msg);
            return;
        }

        jqtoast('注册成功！！');

        //跳转至首页
        setTimeout(function () {
            window.location.href = '/login.html';
        }, 1500)
    });
}

/**
 * 用户登录
 */
function login() {
    //获取用户名
    let username = $('input[name=username]').val();
    //获取密码
    let password = $('input[name=password]').val();

    if (!username || !password) {
        jqtoast("用户名或密码不能为空！");
        return;
    }

    AjaxPost('/v1/user/login', {'username': username, 'password': password}, function (json) {
        if (json.code !== 200) {
            jqtoast(json.msg);
            return;
        }
        jqtoast('登陆成功！');

        //保存token
        setCookie('token', json.result.token, 24 * 60 * 60);
        setCookie('messageHost', json.result.host, 24 * 60 * 60);

        //跳转至首页
        setTimeout(function () {
            window.location.href = '/';
        }, 1500)
    });
}