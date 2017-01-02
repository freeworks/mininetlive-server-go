$(document).ready(function(){
    mininet.renderHtmlNavbar('activity');
    initPriceTypeChange();
	initStreamTypeChange();
    initDateTimePicker();
    
    var $activityAdd = $(".activityAdd");
    var $frontCover = $("#frontCover");
    var $uploadContainer = $(".uploadContainer");
    var params = {};

    $("body").on('click', '.uploadContainer', function(){
        $frontCover.trigger('click');
    });

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
            if (rsp.ret == 0){
                $uploadContainer.css('background-image', "url(" + rsp.data.url + ")");
                $("#frontCoverString").val(rsp.data.url);
            }
        });
    })

    $activityAdd.on('click', function(){
        params = _.parseParams($("#activityForm").serialize());
		if(params.streamType == 1 && params.videoPath == null){
			alert("填写正确的视频地址！")
			return;
		}
		params.isrecommend = $("#isrecommend").is(':checked')?1:0
        params.frontCover = $("#frontCoverString").val();
        // params.date = (new Date(params.date.replace("+", " "))).getTime() / 1000;
        params.date = params.date.replace("+", " ");
        params.price = params.price * 100;
        $(".field-box").removeClass('error');
        mininet.ajax("post", "/activity/new", params, function(rsp){
            if (rsp.ret == 0){
                window.location.href = "/activity-detail.html?aid=" + rsp.data.id;
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
    });
})

function initPriceTypeChange(){
    var $radio = $("input[name=activityType]");
    var $priceContainer = $("#priceContainer");
    $radio.on('change', function(){
       $priceContainer.toggle();
    })
}

function initStreamTypeChange(){
    var $radio = $("input[name=streamType]");
    var $streamContainer = $("#streamContainer");
    $radio.on('change', function(){
       $streamContainer.toggle();
    })
}

function initDateTimePicker(){
    $('.datepicker').datetimepicker({
        language: 'zh-CN',
        initialDate: new Date(),
        format: "yyyy-mm-dd hh:ii",
        startDate: new Date(),
    }).on('changeDate', function (ev) {
        $(this).datetimepicker('hide');
    });
}