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

    var img = document.createElement('img');
    img.src = activity.owner.avatar;
    img.onload = function(){
        debugger
        $("#owner_avatar").attr("src", activity.owner.avatar + "?iopcmd=thumbnail&type=8&width=64&height=64");
    }
    
    $("#desc").text(activity.desc);
    // $("#qrcode").attr("src", activity.owner.qrcode);

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
    if (date == "Invalid Date"){
        date = new Date();
    }
    return formateTwo(date.getMonth() + 1) + "-" + formateTwo(date.getDate()) + " " + formateTwo(date.getHours()) + ":" + formateTwo(date.getMinutes());
}