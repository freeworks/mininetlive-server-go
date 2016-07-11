$(document).ready(function(){
    mininet.renderHtmlNavbar('activity');
    
    var $activityAdd = $(".activityAdd");
    var $frontCover = $("#frontCover");
    var params = {};

    $frontCover.on('change', function(){
        var file = this.files[0];
        name = file.name;
        size = file.size;
        type = file.type;
        debugger
        mininet.ajax("post", "/upload");
    })

    $activityAdd.on('click', function(){
        params = _.parseParams($("#activityForm").serialize());
        params.activityType = $(this).data("activitytype");

        debugger
        mininet.ajax("post", "/activity/new", params, function(rsp){
            if (rsp.ret == 0){
                window.location.href = "/activity-detail.html?aid=" + rsp.data.id;
            } else {
                // TODO 非正常处理
            }
        })
    })

    
})