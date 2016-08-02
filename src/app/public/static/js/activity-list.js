$(document).ready(function(){
    mininet.renderHtmlNavbar('activity');

    var $activityList = $("#activityList");
    var params = _.parseUrlParams();
    params.pageSize = params.pageSize || 5;

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
    });

    $("body").on('click', '.delete', function(){
        var $this = $(this);
        mininet.ajax("delete", "/activity/delete/" + $(this).data("aid"), {}, function(rsp){
            if (rsp.ret == 0){
               $this.parent().parent().hide();
            } else {
                alert("删除失败");
            }
        });
    })
})

function renderHtmlActivityRow(activity){
    return '<tr class="first">' + 
        '<td>' + '<a href="/activity-detail.html?aid=' + activity.aid + '" class="name">' + activity.title + '</a>' + '</td>' +
        '<td><img src="' + activity.frontCover + '" style="width: 150px;height:75px;" class="avatar hidden-phone" />' + '</td>' +
        '<td>' + mininet.formatActivityState(activity.activityState) +'</td>' +
        '<td>' + mininet.formateActivityType(activity.activityType) +'</td>' +
        // '<td>' + mininet.formateAppoinState(activity.appoinState) +'</td>' +
        '<td>' + activity.appointmentCount +'</td>' +
        '<td>￥' + (activity.price / 100).toFixed(2) +'元</td>' +
        // '<td>' + mininet.formatePayState(activity.payState) +'</td>' +
        '<td>' + activity.owner.nickname +'</td>' +
        '<td>' + activity.createTime + '</td>' +
        '<td class="align-right"><a class="btn-flat warning delete" data-aid="' + activity.aid + '">删除</a></tr>' + 
    '</tr>'
}