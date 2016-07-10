$(document).ready(function(){
    mininet.renderHtmlNavbar('index');

    mininet.ajax("post", "/", {}, function(rsp){
        if (rsp.ret == 0){
            debugger
            var key;
            for (key in rsp.data){
                $("#" + key).text(rsp.data[key]);
            }
        } else {
            // TODO 非正常处理
        }
    })
})