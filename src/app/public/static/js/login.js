$(document).ready(function(){
    var $loginbtn = $("#loginbtn");
    var $phone = $("#phone");
    var $password = $("#password");

    var params = _.parseUrlParams();debugger

    $loginbtn.click(function(){
        var phone = $phone.val().trim();
        var password = $password.val().trim();
        // TODO 手机号码，密码校验
        var data = {
            phone: phone,
            password: password
        }

        mininet.ajax("post", "/login", data, function(rsp){
            if (rsp.ret == 0){
                window.location.href = "/"
                // window.location.href = params.next || "/index.html";
            } else {
                // TODO 错误提示
            }
        })
    });
})