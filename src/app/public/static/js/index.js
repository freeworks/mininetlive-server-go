$(document).ready(function(){
    mininet.renderHtmlNavbar('index');

    mininet.ajax("post", "/", {}, function(rsp){
        if (rsp.ret == 0){
            console.log()
        } else {
            // TODO 非正常处理
        }
    })
})