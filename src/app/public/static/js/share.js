$(document).ready(function(){
    var params = _.parseUrlParams();

    mininet.ajax("get", "/share/wx/activity/" + params.aid, {}, function(rsp){
        debugger
        if (rsp.ret == 0){
            var activity =rsp.data;
            renderHtml(activity);
            $(".container").show();
            share(rsp.data.title, location.href, rsp.data.title, rsp.data.frontCover+"?iopcmd=thumbnail&type=8&width=80&height=80");
        } else {
            // TODO 非正常处理
        }
    })

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


function renderHtml(activity){
    var params = _.parseUrlParams();
    
    if (activity.videoPath){
        $("#video").css("width", "100%")
        .css("height", window.innerWidth * 100 / 187.5)
        .css("background-image", activity.frontCover)
        .attr("poster", activity.frontCover)
        .attr("src", activity.videoPath).show();
    } else {
        $("#frontCover").css("height", window.innerWidth * 100 / 187.5).css("background-image", "url(" + activity.frontCover + ")").show();
    }
    
    $("#title").text(activity.title); document.title = activity.title;
    $("#date").text(formateDate(activity.date));

    // var img = document.createElement('img');
    // img.src = activity.owner.avatar;
    // img.onload = function(){
        // debugg
        $(".avatar").attr("src", activity.owner.avatar + "?iopcmd=thumbnail&type=8&width=64&height=64");
    // }
    
    $("#desc").text(activity.desc);
    // $("#qrcode").attr("src", activity.owner.qrcode);
    
    if (activity.streamType == 0){
        // 直播
        if (activity.activityState == 0){
        // 未开播
            $("#appointmentCount").text(activity.appointmentCount);
            $(".state2").show();
        }
        if (activity.activityState == 1){
            // 直播中
            $(".live").show();
            $("#onlineCount").text(activity.onlineCount);
            $(".state3").show();
        }

        if (activity.activityState == 2) {
        // 直播结束
        
        }
    } else {
        $("#playCount").text(activity.playCount);
        $(".state1").show();
        // 点播
    }

    if (activity.activityType == 0){
        // 免费
    } else {
        $(".price").text("￥" + (activity.price / 100.0).toFixed(2)).show();
    }

    if (params.icode){
        $("#invite_code").show();
        $("#icode").text(params.icode);
    }
}

function formateTwo(string){
    string += "";
    if (string.length == 1){
        return 0 + string;
    }
    return string;
}

function formateDate(date){
    date = new Date(date);
    if (date == "Invalid Date"){
        date = new Date();
    }
    return formateTwo(date.getMonth() + 1) + "-" + formateTwo(date.getDate()) + " " + formateTwo(date.getHours()) + ":" + formateTwo(date.getMinutes());
}

function share(title, link, desc, imgUrl){
    wx.ready(function(){
        wx.onMenuShareTimeline({
            title: title, // 分享标题
            link: link, // 分享链接
            imgUrl: imgUrl // 分享图标
            // success: function () { 
            //     // 用户确认分享后执行的回调函数
            // },
            // cancel: function () { 
            //     // 用户取消分享后执行的回调函数
            // }
        });

        wx.onMenuShareAppMessage({
            title: title, // 分享标题
            desc: desc, // 分享描述
            link: link, // 分享链接
            imgUrl: imgUrl // 分享图标
            // type: '', // 分享类型,music、video或link，不填默认为link
            // dataUrl: '', // 如果type是music或video，则要提供数据链接，默认为空
            // success: function () { 
            //     // 用户确认分享后执行的回调函数
            // },
            // cancel: function () { 
            //     // 用户取消分享后执行的回调函数
            // }
        });
    });
}
