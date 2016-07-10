var config = {
    host: "http://127.0.0.1:8081"
}

var mininet = {};

function ajax(method, path, data, success, fail){
    $.ajax({
        method: method,
        url: config.host + path,
        contentType: "application/x-www-form-urlencoded",
        data: data,
        success: success,
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

function renderHtmlNavbar(route){
    var $siderbar = $("#sidebar-nav");
    var navbar = '<ul id="dashboard-menu">' +
                    '<li class="index"><a href="/"><i class="icon-home"></i><span>首页</span></a></li>' +
                    '<li class="user"><a href="/user-list.html"><i class="icon-group"></i><span>用户</span></a></li>' +      
                    '<li class="order"><a href="/order-list.html"><i class="icon-signal"></i><span>收入</span></a></li>' +
                    '<li class="activity"><a href="/activity"><i class="icon-th-large"></i><span>活动</span></a></li>' +
                    '<li class="logout"><a href="/logout"><i class="icon-share-alt"></i><span>退出</span></a></li>' +
                    '</ul>';
    var pointer = '<div class="pointer"><div class="arrow"></div><div class="arrow_border"></div></div>';
    $siderbar.html(navbar);

    if (route){
        var current = $("#dashboard-menu ." + route);
        current.addClass('active');
        current.prepend(pointer);
    }
}

mininet.ajax = ajax;
mininet.formatGender = formatGender;
mininet.formatPlat = formatPlat;
mininet.formatChannel = formatChannel;
mininet.renderHtmlNavbar = renderHtmlNavbar;