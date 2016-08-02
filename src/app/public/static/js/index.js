$(document).ready(function() {
    mininet.renderHtmlNavbar('index');

    
    mininet.ajax("post", "/", {}, function(rsp) {
        debugger
        if (rsp.ret == 0) {
            $("#newUserCount").text(rsp.data["newUserCount"]);
            $("#newOrderCount").text(rsp.data["newOrderCount"]);
            $("#newAmount").text((rsp.data["newAmount"] / 100.0).toFixed(2));
        } else {
            // TODO 非正常处理
        }
    })

    var $orderChat = $("#orderChat");
    var $incomeChat = $("#incomeChat");
    var orderData = [];
    var orderTicks = [];
    var incomeData = [];
    var incomeTicks = [];

    mininet.ajax("get", "/order/chart/graph", {}, function(rsp){
        debugger
        if (rsp.ret == 0){
            for (var i = 0; i < rsp.data.length; i++){
                orderData.push([i+1, rsp.data[i].count]);
                orderTicks.push([i+1, rsp.data[i].date])   
            }
            var options = {
                title: "订单趋势图"
            }
            showChat($orderChat, orderData, orderTicks, options)
        } 
    })

    mininet.ajax("get", "/income/chart/graph", {}, function(rsp){
        debugger
        if (rsp.ret == 0){
            for (var i = 0; i < rsp.data.length; i++){
                rsp.data[i].count = (rsp.data[i].count / 100.0).toFixed(2);
                incomeData.push([i+1, rsp.data[i].count])   
                incomeTicks.push([i+1, rsp.data[i].date])   
            }
            var options = {
                title: "收入趋势图"
            }
           showChat($incomeChat, incomeData, incomeTicks, options)
        } 
    })

})

function showChat($element, orderData, ticks, options) {
    var visits = orderData
    var plot = $.plot($element, [
        { data: orderData, label: options.title },
        // { data: visitors, label: "收入趋势图" }
    ], {
        series: {
            lines: {
                show: true,
                lineWidth: 1,
                fill: true,
                fillColor: { colors: [{ opacity: 0.1 }, { opacity: 0.13 }] }
            },
            points: {
                show: true,
                lineWidth: 2,
                radius: 3
            },
            shadowSize: 0,
            stack: true
        },
        grid: {
            hoverable: true,
            clickable: true,
            tickColor: "#f9f9f9",
            borderWidth: 0
        },
        legend: {
            // show: false
            labelBoxBorderColor: "#fff"
        },
        colors: ["#30a0eb", "#a7b5c5"],
        xaxis: {
            ticks: ticks,
            font: {
                size: 12,
                family: "Open Sans, Arial",
                variant: "small-caps",
                color: "#697695"
            }
        },
        yaxis: {
            ticks: 3,
            tickDecimals: 0,
            font: { size: 12, color: "#9da3a9" }
        }
    });

    function showTooltip(x, y, contents) {
        $('<div id="tooltip">' + contents + '</div>').css({
            position: 'absolute',
            display: 'none',
            top: y - 30,
            left: x - 50,
            color: "#fff",
            padding: '2px 5px',
            'border-radius': '6px',
            'background-color': '#000',
            opacity: 0.80
        }).appendTo("body").fadeIn(200);
    }

    var previousPoint = null;
    $element.bind("plothover", function(event, pos, item) {
        if (item) {
            if (previousPoint != item.dataIndex) {
                previousPoint = item.dataIndex;

                $("#tooltip").remove();
                var x = item.datapoint[0].toFixed(0),
                    y = item.datapoint[1].toFixed(0);

                var month = item.series.xaxis.ticks[item.dataIndex].label;

                showTooltip(item.pageX, item.pageY,
                    month + " - " + y);
            }
        } else {
            $("#tooltip").remove();
            previousPoint = null;
        }
    });
}
