$(document).ready(function(){
    mininet.renderHtmlNavbar('activity');

    var $activityList = $("#activityList");
    var params = _.parseUrlParams();
    params.pageSize = params.pageSize || 10;

    mininet.ajax("get", "/activity/list", params, function(rsp){
        debugger
        if (rsp.ret == 0){
            var activityList = rsp.data.activityList;
            activityList.forEach(function(activity){
                $activityList.append(renderHtmlActivityRow(activity));
            })

            var $pagination = $("#pagination");
            $pagination.append(mininet.renderHtmlPagination(rsp.data.totalPageCount, params.pageIndex, params.pageSize));
        } else {
            // TODO 非正常处理
        }
    })
})

function renderHtmlActivityRow(activity){
    return '<tr class="first">' + 
        '<td>' + '<a href="/activity-detail.html?aid=' + activity.aid + '" class="name">' + activity.title + '</a>' + '</td>' +
        '<td><img src="' + activity.frontCover + '" style="width: 150px;height:75px;" class="avatar hidden-phone" />' + '</td>' +
        '<td>' + activity.activityState +'</td>' +
        '<td>' + activity.activityType +'</td>' +
        '<td>' + activity.appoinState +'</td>' +
        '<td>' + activity.appointmentCount +'</td>' +
        '<td>' + activity.price +'</td>' +
        '<td>' + activity.payState +'</td>' +
        '<td>' + activity.owner.nickname +'</td>' +
        '<td class="align-right">' + activity.createTime + '</td>' +
    '</tr>'
}