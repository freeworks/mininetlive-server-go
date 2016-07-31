$(document).ready(function(){
    var params = _.parseUrlParams();

    mininet.ajax("get", "/share/wx/activity/" + params.aid, {}, function(rsp){
        if (rsp.ret == 0){
            debugger
            var activity =rsp.data;
            renderHtml(activity);
            $(".container").show();
        } else {
            // TODO 非正常处理
        }
    })
});


function renderHtml(activity){
    $("#frontCover").css("height", window.innerWidth * 100 / 187.5).css("background-image", "url(" + activity.frontCover + ")");
    $("#title").text(activity.title); document.title = activity.title;
    $("#date").text(formateDate(activity.date));
    $("#owner_avatar").attr("src", activity.owner.avatar);
    $("#desc").text(activity.desc);
    $("#qrcode").attr("src", activity.owner.qrcode);

    if (activity.activityType == 0){
        $("#playCount").text(activity.playCount);
        $(".state1").show();
    } else if (activity.appoinState == 0){
        
        $("#appointmentCount").text(activity.appointmentCount);
        $(".state2").show();
    } else {
        $("#onlineCount").text(activity.onlineCount);
        $(".state3").show();
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
    return formateTwo(date.getMonth() + 1) + "-" + formateTwo(date.getDate()) + " " + formateTwo(date.getHours()) + ":" + formateTwo(date.getMinutes());
}