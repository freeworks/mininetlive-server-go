$(document).ready(function(){
    mininet.renderHtmlNavbar('order');

    var $userList = $("#userList");

    mininet.ajax("get", "/order/list", {
        pageSize: 10,
        pageIndex: 0
    }, function(rsp){
        debugger
        if (rsp.ret == 0){
            var orderList = rsp.data;
            orderList.forEach(function(order){
                $userList.append(renderHtmlOrderRow(order));
            })
        } else {
            // TODO 非正常处理
        }
    })
})

function renderHtmlOrderRow(order){
    return '<tr class="first">' + 
        '<td>' + order.no + '</td>' + 
        '<td>' + order.subject + '</td>' + 
        '<td>' + mininet.formatChannel(order.channel) + '</td>' + 
        '<td>￥ ' + order.amount + '</td>' + 
        '<td>' + order.type + '</td>' + 
        '<td class="align-right">' + order.createTime + '</td>' + 
    '</tr>'
}