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

mininet.ajax = ajax;