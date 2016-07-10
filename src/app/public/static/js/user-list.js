$(document).ready(function(){
    mininet.renderHtmlNavbar('user');

    var $userList = $("#userList");
    var $pagination = $("#pagination");

    var params = mininet.parseUrlParams();
    params.pageSize = params.pageSize || 2;
    if (params.pageIndex > 0){
        params.pageIndex = params.pageIndex - 1;
    } else {
        params.pageIndex = 0;
    }

    mininet.ajax("get", "/user/list", params, function(rsp){
        debugger
        if (rsp.ret == 0){
            var userList = rsp.data.userList;
            userList.forEach(function(user){
                $userList.append(renderHtmlUserRow(user));
            })
            $pagination.append(mininet.renderHtmlPagination(rsp.data.totalPageCount, params.pageIndex + 1, params.pageSize));
        } else {
            // TODO 非正常处理
        }
    })
})

function renderHtmlUserRow(user){
    return '<tr class="first">' + 
        '<td><img src="' + user.avatar + '" class="img-circle avatar hidden-phone" />' +
            '<a href="user-profile.html" class="name">' + user.nickname + '</a>' + 
        '</td>' +
        '<td>' + mininet.formatGender(user.gender) + '</td>' + 
        '<td>' + user.phone +'</td>' +
        '<td>￥ ' + user.balance + '</td>' +
        '<td>' + mininet.formatPlat(user.plat) + '</td>' +
        '<td class="align-right">' + user.createTime + '</td>' +
    '</tr>'
}