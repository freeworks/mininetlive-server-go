$(document).ready(function(){
    mininet.renderHtmlNavbar('order');

    var $userList = $("#userList");
    var params = _.parseUrlParams();
    params.pageSize = params.pageSize || 10;

    mininet.ajax("get", "/order/list", params, function(rsp){
        debugger
        if (rsp.ret == 0){
            var orderList = rsp.data.orderList;
			if(orderList == null){
				alert("订单数据异常！")
				return;
			}
            orderList.forEach(function(order){
                $userList.append(renderHtmlOrderRow(order));
            })

            var $pagination = $("#pagination");
            $pagination.append(mininet.renderHtmlPagination(rsp.data.totalPageCount, params.pageIndex, params.pageSize));
        } else {
			alert(rsp.msg);
        }
    })
})

function renderHtmlOrderRow(order){
    return '<tr class="first">' + 
        '<td>' + order.no + '</td>' + 
        '<td>' + order.subject + '</td>' + 
        '<td>' + mininet.formatChannel(order.channel) + '</td>' + 
        '<td>￥ ' + (order.amount / 100.0).toFixed(2)+ '元</td>' + 
        '<td>' + mininet.formateOrderState(order.type) + '</td>' + 
        '<td class="align-right">' + order.createTime + '</td>' + 
    '</tr>'
}