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
        }else{
			alert(rsp.msg)
		}
    })
    
    var $activityAdd = $(".activityAdd");
    var $frontCover = $("#frontCover");
    var $frontCoverImg = $("#frontCoverImg");
    var $uploadContainer = $(".uploadContainer");
    var $videoFilePath = $("#videoFilePath")
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
		if(params.streamType == 1 && params.videoPath == null){
			alert("填写正确的视频地址！")
			return;
		}
		params.price = parseFloat(params.price)*100
//        params.activityType = $('input[name="activityType"]:checked').val();
        params.frontCover = $("#frontCoverString").val();
        // params.date = (new Date(params.date.replace("+", " "))).getTime() / 1000;
        params.date = params.date.replace("+", " ");
		params.isrecommend = $("#isrecommend").is(':checked')?1:0
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

    $videoFilePath.on('change', function(){
        var formData = new FormData();
        var file = this.files[0];
        name = file.name;
        size = file.size;
        type = file.type;
        if(type != "video/mp4"){
            alert("视频文件类型错误!");
            return;
        }
        formData.append('file', file)
        var contentType = false;
        $("#videoFilePathLabel")[0].innerText="上传中..."
        mininet.ajaxFile("post", "/uploadVideo", formData, function(rsp){
            debugger
            if (rsp.ret == 0){
                $("#videoPath").val(rsp.data.url);
            }else{
                 alert("视频文件上传失败!");
            }
            $("#videoFilePathLabel")[0].innerText="选择文件"
        });
    })
});

function renderHtmlActivityForm(activity){
    debugger
    $(".uploadContainer").css('background-image', "url(" + activity.frontCover+'?iopcmd=thumbnail&type=8&width=300&height=150' 
	+ ")");
    $("#frontCoverString").val(activity.frontCover);
    $("#title").val(activity.title);
    $("#desc").val(activity.desc);
    $("#price").val((activity.price / 100).toFixed(2));
    if (activity.activityType == 1){
        $("#priceContainer").show();
    }
	$("input[name=activityType][value=" + activity.activityType + "]").prop("checked", true);
    initDateTimePicker(activity.date);
    $("#date").val(activity.date);
	$("#isrecommend").prop("checked", activity.isrecommend);
	$("input[name=streamType][value=" + activity.streamType + "]").prop("checked", true);
	if (activity.streamType == 1){
        $("#streamContainer").show();
    }
	$("#videoPath").val(activity.videoPath)
}

function initPriceTypeChange(){
    var $radio = $("input[name=activityType]");
    var $priceContainer = $("#priceContainer");
    $radio.on('change', function(){
       $priceContainer.toggle();
    })
	$radio = $("input[name=streamType]");
    $streamContainer = $("#streamContainer");
	$radio.on('change', function(){
       $streamContainer.toggle();
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