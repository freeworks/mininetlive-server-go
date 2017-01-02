var config = {
    host: "http://127.0.0.1:8081"
}

var mininet = {};

function ajax(method, path, data, success, error){
    debugger
    $.ajax({
        url: path,
        method: method,
        contentType: "application/x-www-form-urlencoded",
        data: data,
        success: function(rsp){
            if (rsp.ret == -1){
                var next = rsp.data.redirectPath;
                window.location.href = "/login.html";
                // window.location.href = "/login.html?next=" + next;
            } else {
                success(rsp);
            }
        },
        error: error
    })
}


function ajaxFile(method, path, data, success, fail, contentType){
    // contentType = contentType || "application/x-www-form-urlencoded";
    debugger
    $.ajax({
        url: path,
        method: method,
        contentType: false,
        processData: false,
        data: data,
        success: function(rsp){
            if (rsp.ret == -1){
                var next = rsp.data.redirectPath;
                window.location.href = "/login.html?next=" + next;
            } else {
                success(rsp);
            }
        },
        fail: fail
    })
}


function formatGender(gender){
    switch(gender){
        case 0:
            return "女"
        case 1:
            return "男"
        default:
            return "未知"
    }
}

function formatPlat(plat){
    switch(plat){
        case "QQ":
            return "QQ"
        case "Wechat":
            return "微信"
        case "SinaWeibo":
            return "新浪微博"
        default:
            return plat;
    }
}

function formatChannel(plat){
    switch(plat){
        case "wx":
            return "微信"
        default:
            return plat;
    }
}

function formatStreamType(type){
    switch(type){
        case 1:
            return "点播"
        case 0:
            return "直播"
        default:
            return type;
    }
}

function formatActivityState(state){
    switch(state){
        case 0:
            return "未开播";
        case 1:
            return "直播中";
        case 2:
            return "已结束"
    }
}

function changeTow(number){
    number += "";
    if (number.length == 1){
        return "0" + number;
    }
    return number;
}

function formatDateTime(date){
    date = new Date(date);
    return  date.getFullYear() + "-" + (date.getMonth() + 1) + "-" + date.getDate() + " " + date.getHours() + ":" + changeTow(date.getMinutes());
}

function renderHtmlNavbar(route){
    var $siderbar = $("#sidebar-nav");
    var navbar = '<ul id="dashboard-menu">' +
                    '<li class="index"><a href="/"><i class="icon-home"></i><span>首页</span></a></li>' +
                    '<li class="user"><a href="/user-list.html"><i class="icon-group"></i><span>用户</span></a></li>' +      
                    '<li class="order"><a href="/order-list.html"><i class="icon-signal"></i><span>收入</span></a></li>' +
                    '<li class="activity"><a href="/activity-list.html"><i class="icon-th-large"></i><span>活动</span></a></li>' +
                    '<li class="logout"><a href="#"><i class="icon-share-alt"></i><span>退出</span></a></li>' +
                    '</ul>';
    var pointer = '<div class="pointer"><div class="arrow"></div><div class="arrow_border"></div></div>';
    $siderbar.html(navbar);

    if (route){
        var current = $("#dashboard-menu ." + route);
        current.addClass('active');
        current.prepend(pointer);
    }

    // 退出
    $siderbar.on('click', '.logout', function(e){
        debugger
        // e.preventDefault();
        ajax("post", "/logout", {}, function(rsp){
            debugger
            if (rsp.ret == 0){
                window.location.href = "/login.html";
            }
        }, function(){
            debugger
        })
    })
}

function renderHtmlPagination(total, current, pageSize){
    total = parseInt(total);
    current = parseInt(current || 1);
    pageSize = parseInt(pageSize);
    
    var params = {
        pageSize: pageSize || 10
    };

    var html = "<ul>";

    // 前一页
    params.pageIndex = current - 1 || 1;
    html += '<li><a href="' + newLocationPath(params) + '">&#8249;</a></li>';

    
    for (var i = 1; i <= total; i++){
        params.pageIndex = i;
        if (current != i){
            html += '<li><a href="' + newLocationPath(params) + '">' + i + '</a></li>';
        } else {
            html += '<li><a href="' + newLocationPath(params) + '" class="active" >' + i + '</a></li>';
        }   
    }

    // 后一页
    params.pageIndex = current + 1;
    if (params.pageIndex > total){
        params.pageIndex = total;
    }
    html += '<li><a href="' + newLocationPath(params) + '">&#8250;</a></li>';

    html += '</ul>';
    return html
}


function newLocationPath(params){
    return location.pathname + "?" + _.stringifyUrlParams(params);
}

function formatActivityState(state){
    switch(state){
        case 0:
            return "未开播"
        case 1:
            return "直播中"
        case 2:
            return "直播结束"
        default:
            return state;
    }
}

function formateActivityType(state){
    switch(state){
        case 0:
            return "免费"
        case 1:
            return "收费"
        default:
            return state;
    }
}

function formateStreamType(state){
    switch(state){
        case 0:
            return "直播"
        case 1:
            return "视频"
        default:
            return state;
    }
}

function formateAppoinState(state){
    switch(state){
        case 0:
            return "未预约"
        case 1:
            return "已预约"
        default:
            return state;
    }
}

function formatePayState(state){
    switch(state){
        case 0:
            return "未支付"
        case 1:
            return "已支付"
        default:
            return state;
    }
}

function formateOrderState(state){
    switch(state){
        case 0:
            return "未支付"
        case 1:
            return "已支付"
        default:
            return state;
    }
}

/**
* 实时动态强制更改用户录入
* arg1 inputObject
**/
function amount(th){
    var regStrs = [
        ['^0(\\d+)$', '$1'], //禁止录入整数部分两位以上，但首位为0
        ['[^\\d\\.]+$', ''], //禁止录入任何非数字和点
        ['\\.(\\d?)\\.+', '.$1'], //禁止录入两个以上的点
        ['^(\\d+\\.\\d{2}).+', '$1'] //禁止录入小数点后两位以上
    ];
    for(i=0; i<regStrs.length; i++){
        var reg = new RegExp(regStrs[i][0]);
        th.value = th.value.replace(reg, regStrs[i][1]);
    }
}
 
/**
* 录入完成后，输入模式失去焦点后对录入进行判断并强制更改，并对小数点进行0补全
* arg1 inputObject
* 这个函数写得很傻，是我很早以前写的了，没有进行优化，但功能十分齐全，你尝试着使用
* 其实有一种可以更快速的JavaScript内置函数可以提取杂乱数据中的数字：
* parseFloat('10');
**/
function overFormat(th){
    var v = th.value;
    if(v === ''){
        v = '0.00';
    }else if(v === '0'){
        v = '0.00';
    }else if(v === '0.'){
        v = '0.00';
    }else if(/^0+\d+\.?\d*.*$/.test(v)){
        v = v.replace(/^0+(\d+\.?\d*).*$/, '$1');
        v = inp.getRightPriceFormat(v).val;
    }else if(/^0\.\d$/.test(v)){
        v = v + '0';
    }else if(!/^\d+\.\d{2}$/.test(v)){
        if(/^\d+\.\d{2}.+/.test(v)){
            v = v.replace(/^(\d+\.\d{2}).*$/, '$1');
        }else if(/^\d+$/.test(v)){
            v = v + '.00';
        }else if(/^\d+\.$/.test(v)){
            v = v + '00';
        }else if(/^\d+\.\d$/.test(v)){
            v = v + '0';
        }else if(/^[^\d]+\d+\.?\d*$/.test(v)){
            v = v.replace(/^[^\d]+(\d+\.?\d*)$/, '$1');
        }else if(/\d+/.test(v)){
            v = v.replace(/^[^\d]*(\d+\.?\d*).*$/, '$1');
            ty = false;
        }else if(/^0+\d+\.?\d*$/.test(v)){
            v = v.replace(/^0+(\d+\.?\d*)$/, '$1');
            ty = false;
        }else{
            v = '0.00';
        }
    }
    th.value = v; 
}


mininet.ajax = ajax;
mininet.ajaxFile = ajaxFile;
mininet.formatGender = formatGender;
mininet.formatPlat = formatPlat;
mininet.formatChannel = formatChannel;
mininet.renderHtmlNavbar = renderHtmlNavbar;
mininet.renderHtmlPagination = renderHtmlPagination;
mininet.formatStreamType = formatStreamType;
mininet.formatDateTime = formatDateTime;
mininet.formatActivityState = formatActivityState;
mininet.formateActivityType = formateActivityType;
mininet.formateStreamType = formateStreamType;
mininet.formateAppoinState = formateAppoinState;
mininet.formatePayState = formatePayState;
mininet.formateOrderState = formateOrderState;
mininet.amount = amount;
mininet.overFormat = overFormat;
