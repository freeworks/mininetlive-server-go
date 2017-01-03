$(document).ready(function() {
    mininet.renderHtmlNavbar('activity');
    
    var $activityDetail = $("#activityDetail");

    var params = _.parseUrlParams();
    var aid = params.aid;

    $("#editActivity").attr("href", "/activity-edit.html?aid=" + aid);

    mininet.ajax("get", "/activity/detail/" + aid, {}, function(rsp) {
        debugger
        if (rsp.ret == 0) {
            var activity = rsp.data;
            $activityDetail.append(renderHtmlActivity(activity));
        }else{
			alert(rsp.msg);
		}
    })
});


function renderHtmlActivity(activity){
    var html = "";
    html += '<tr><td>' + '标题' + '</td><td>' + activity.title + '</td></tr>';
    html += '<tr><td>' + '时间' + '</td><td>' + activity.date + '</td></tr>';
    html += '<tr><td>' + '描述' + '</td><td>' + activity.desc + '</td></tr>';
    html += '<tr><td>' + '预览图' + '</td><td>' + '<img src=" ' + activity.frontCover +'" style="max-width: 300px">' + '</td></tr>';
    html += '<tr><td>' + '类型' + '</td><td>' + mininet.formatStreamType(activity.streamType) + '</td></tr>';
    if (activity.streamType == 1){
        html += '<tr><td>' + '视频路径' + '</td><td>' + activity.videoPath + '</td></tr>';
    }
    if (activity.streamType == 0){
        html += '<tr><td>' + '推流地址' + '</td><td>' + activity.livePullPath + '</td></tr>';
		html += '<tr><td>' + '拉流地址' + '</td><td>' + activity.livePushPath + '</td></tr>';
        html += '<tr><td>' + '状态' + '</td><td>' + formatActivityState(activity.activityState) + '</td></tr>';
        html += '<tr><td>' + '在线观看时间' + '</td><td>' + activity.onlineCount + '</td></tr>';
    }
	if (activity.activityType == 0){
        html += '<tr><td>价格</td><td>免费</td></tr>';
    }else{
		html += '<tr><td>价格</td><td>'+(activity.price / 100).toFixed(2)+'元</td></tr>';
	}
	if (activity.isrecommend == 0){
        html += '<tr><td>推荐</td><td>否</td></tr>';
    }else{
		html += '<tr><td>推荐</td><td>是</td></tr>';
	}
    html += '<tr><td>' + '播放次数' + '</td><td>' + activity.playCount + '</td></tr>';
    html += '<tr><td>' + '预约人数' + '</td><td>' + activity.appointmentCount + '</td></tr>';
    
    return html;
}