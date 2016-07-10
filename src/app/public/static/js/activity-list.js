$(document).ready(function(){
    mininet.renderHtmlNavbar('activity');

    var $activityList = $("#activityList");
    var $pagination = $("#pagination");

    var params = mininet.parseUrlParams();
    params.pageSize = params.pageSize || 2;
    if (params.pageIndex > 0){
        params.pageIndex = params.pageIndex - 1;
    } else {
        params.pageIndex = 0;
    }

    mininet.ajax("get", "/activity/list", {
        pageSize: 10,
        pageIndex: 0
    }, function(rsp){
        debugger
        if (rsp.ret == 0){
            var activityList = rsp.data.activityList;
            activityList.forEach(function(activity){
                $activityList.append(renderHtmlActivityRow(activity));
            })
            $pagination.append(mininet.renderHtmlPagination(rsp.data.totalPageCount, params.pageIndex + 1, params.pageSize));
        } else {
            // TODO 非正常处理
        }
    })
})

function renderHtmlActivityRow(activity){
    return '<tr class="first">' + 
        '<td><img src="' + activity.frontCover + '" class="img-circle avatar hidden-phone" />' +
            '<a href="user-profile.html" class="name">' + activity.title + '</a>' + 
        '</td>' +
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