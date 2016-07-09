$(document).ready(function(){
    var $loginbtn = $("#loginbtn");
    var $phone = $("#phone");
    var $password = $("#password");

    $loginbtn.click(function(){
        var phone = $phone.val().trim();
        var password = $password.val().trim();
        // TODO 手机号码，密码校验
        var data = {
            phone: phone,
            password: password
        }

        mininet.ajax("post", "/login", data, function(rsp){
            console.log(111111111)
            console.log(rsp);
            if (rsp.ret == 0){
                var redirectPath = rsp.data.redirectPath || "/index.html";
                window.location.href = redirectPath;
            } else {
                // TODO 错误提示
            }
        })
    });
})