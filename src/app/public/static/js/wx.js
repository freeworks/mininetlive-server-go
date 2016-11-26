$(document).ready(function(){
    $.ajax({
        url: "/wxpub/jsconfig",
        method: "POST",
        contentType: "application/x-www-form-urlencoded",
        data: {
            url: location.href
        },
        success: function(rsp){
            wx.config({
                debug: true, // 开启调试模式,调用的所有api的返回值会在客户端alert出来，若要查看传入的参数，可以在pc端打开，参数信息会通过log打出，仅在pc端时才会打印。
                appId: rsp.data.appId, // 必填，公众号的唯一标识
                timestamp: rsp.data.timestamp, // 必填，生成签名的时间戳
                nonceStr: rsp.data.nonceStr, // 必填，生成签名的随机串
                signature: rsp.data.signature,// 必填，签名，见附录1
                jsApiList: ["onMenuShareTimeline", "onMenuShareAppMessage"] // 必填，需要使用的JS接口列表，所有JS接口列表见附录2
            });
        }
    })
});