$(document).ready(function(){

    $("#submit").on("click", function(){

        showErrorCode("");

        var phone = $("#phone").val().trim();
        var vcode = $("#vcode").val().trim();

        if(!validatePhone(phone)){
            showErrorCode("请输入正确的手机号码");
            return;
        }

        if(!validateCode(vcode)){
            showErrorCode("请输入六位验证码");
            return;
        }

		var openId
  		var reg = new RegExp("(^|&)=([^&]*)(&|$)", "i");
    	var r = location.search.substr(1).match(reg);
    	if (r != null) 
			openId = unescape(r[2]);
    	}else{
			alert("链接错误!");
			return;
		}
        $.ajax({
            url: "/wxpub/bindphone",
            method: "POST",
            contentType: "application/x-www-form-urlencoded",
            data: {
                phone: phone,
                vcode: vcode,
                openId: openId
            },
            success: function(rsp){
                if (rsp.ret === 0){
                    alert("绑定成功");
                    wx.closeWindow();
                } else {
                    showErrorCode(rsp.msg);
                }
            }
        })
    })

    $("#send").on("click", sendCallback);
});

function sendCallback(){
    showErrorCode("");

    var phone = $("#phone").val().trim();
    if(!validatePhone(phone)){
        showErrorCode("请输入正确的手机号码");
        return;
    }
    $.ajax({
        url: "/wxpub/vcode",
        method: "POST",
        contentType: "application/x-www-form-urlencoded",
        data: {
            phone: phone
        },
        success: function(rsp){
            if (rsp.ret === 0){
                var time = 10;
                $("#send").off("click", sendCallback);
                var timeout = setInterval(function(){
                    showSendCode(time + "秒后可重新发送");
                    time--;
                    if (time == 0){
                        $("#send").on("click", sendCallback);
                        clearInterval(timeout);
                        showSendCode("发送验证码");
                    }
                }, 1000)
            } else {
                showErrorCode(rsp.msg);
            }
        }
    })
}

function showErrorCode(text){
    $("#error").text(text);
}

function showSendCode(text){
    $("#send").text(text);
}

function validatePhone(phone){
    if (phone && /^1\d{10}$/.test(phone)){
        return true;
    } else {
        return false;
    }
}

function validateCode(code){
    if (code && /^\d{6}$/.test(code)){
        return true;
    } else {
        return false;
    }
}

function getVCode(){
    debugger
    var phone = $("#phone").val().trim();
    $.ajax({
        url: "/wxpub/vcode",
        method: "POST",
        contentType: "application/x-www-form-urlencoded",
        data: {
            phone: phone
        },
        success: function(rsp){
            alert("success");
            if (rsp.ret === 0){
                var next = rsp.data.redirectPath;
                window.location.href = "/login.html";
                // window.location.href = "/login.html?next=" + next;
            } else {
                success(rsp);
            }
        }
    })
}