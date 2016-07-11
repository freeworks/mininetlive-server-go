$(document).ready(function() {
    mininet.renderHtmlNavbar('activity');
    
    var $activityDetail = $("#activityDetail");

    var params = _.parseUrlParams();
    var aid = params.aid;

    mininet.ajax("get", "/activity/detail/" + aid, {}, function(rsp) {
        debugger
        if (rsp.ret == 0) {
            var activity = rsp.data;
            $activityDetail.append(renderHtmlActivity(activity));
        }
        // TODO 非正常处理
    })
})


function renderHtmlActivity(activity){
    var html = "";
    html += '<tr><td>' + '标题' + '</td><td>' + activity.title + '</td></tr>';
    html += '<tr><td>' + '时间' + '</td><td>' + activity.date + '</td></tr>';
    html += '<tr><td>' + '描述' + '</td><td>' + activity.desc + '</td></tr>';
    html += '<tr><td>' + '预览图' + '</td><td>' + '<img src=" ' + activity.frontCover +'" style="max-width: 300px">' + '</td></tr>';
    html += '<tr><td>' + '类型' + '</td><td>' + mininet.formatStreamType(activity.streamType) + '</td></tr>';
    if (activity.streamType == 1){
        html += '<tr><td>' + '路径' + '</td><td>' + activity.videoPath + '</td></tr>';
    }
    if (activity.streamType == 0){
        html += '<tr><td>' + '路径' + '</td><td>' + activity.livePushPath + '</td></tr>';
        html += '<tr><td>' + '状态' + '</td><td>' + activity.activityState + '</td></tr>';
        html += '<tr><td>' + '当前观看时间' + '</td><td>' + activity.onlineCount + '</td></tr>';
    }
    
    html += '<tr><td>' + '播放次数' + '</td><td>' + activity.playCount + '</td></tr>';
    html += '<tr><td>' + '预约人数' + '</td><td>' + activity.appointmentCount + '</td></tr>';
    
    return html;
}