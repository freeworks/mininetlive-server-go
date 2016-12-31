$(document).ready(function(){
    mininet.renderHtmlNavbar('activity');
    initPriceTypeChange();
    // initDateTimePicker();

    var params = _.parseUrlParams();
    var aid = params.aid;
    var id;

    mininet.ajax("get", "/activity/detail/" + aid, {}, function(rsp) {
        debugger
        if (rsp.ret == 0) {
            id = rsp.data.id;
            renderHtmlActivityForm(rsp.data);
        }
        // TODO 非正常处理
    })
    
    var $activityAdd = $(".activityAdd");
    var $frontCover = $("#frontCover");
    var $frontCoverImg = $("#frontCoverImg");
    var $uploadContainer = $(".uploadContainer");
    var params = {};

    $frontCover.on('change', function(){
        var formData = new FormData();
        var file = this.files[0];
        name = file.name;
        size = file.size;
        type = file.type;

        formData.append('file', file)

        var contentType = false;

        mininet.ajaxFile("post", "/upload", formData, function(rsp){
            debugger
             $uploadContainer.css('background-image', "url(" + rsp.data.url + ")");
            $("#frontCoverString").val(rsp.data.url);
        });
    })

    $activityAdd.on('click', function(){
        params = _.parseParams($("#activityForm").serialize());
        params.activityType = $(this).data("activitytype");
        params.frontCover = $("#frontCoverString").val();
        // params.date = (new Date(params.date.replace("+", " "))).getTime() / 1000;
        params.date = params.date.replace("+", " ");

        debugger
        mininet.ajax("put", "/activity/update/" + id, params, function(rsp){
            debugger
            if (rsp.ret == 0){
                window.location.href = "/activity-detail.html?aid=" + aid;
            } else {
				alert(rsp.msg)
            }
        }, function(rsp){
            if (rsp.status == 422){
                var errors = rsp.responseJSON;
                errors.forEach(function(error){
                    $("." + error.fieldNames[0] + "Container").addClass('error');
                })
            }
        })
    })
});

function renderHtmlActivityForm(activity){
    debugger
    $(".uploadContainer").css('background-image', "url(" + activity.frontCover + ")");
    $("#frontCoverString").val(activity.frontCover);
    $("#title").val(activity.title);
    $("#desc").val(activity.desc);
    $("#price").val(activity.price);
    if (activity.activityType == 1){
        $("#priceContainer").show();
    }
    initDateTimePicker(activity.date);
    $("#date").val(activity.date);
    $("input[value=" + activity.activityType + "]").prop("checked", true);
}

function initPriceTypeChange(){
    var $radio = $("input[name=activityType]");
    var $priceContainer = $("#priceContainer");
    $radio.on('change', function(){
       $priceContainer.toggle();
    })
}

function initDateTimePicker(){
    $('.datepicker').datetimepicker({
        language: 'zh-CN',
        format: "yyyy-mm-dd hh:ii",
        startDate: new Date(),
    }).on('changeDate', function (ev) {
        $(this).datetimepicker('hide');
    });
}