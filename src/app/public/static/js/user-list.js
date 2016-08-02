$(document).ready(function(){
    mininet.renderHtmlNavbar('user');

    var $userList = $("#userList");
    var params = _.parseUrlParams();
    params.pageSize = params.pageSize || 5;

    mininet.ajax("get", "/user/list", params, function(rsp){
        debugger
        if (rsp.ret == 0){
            var userList = rsp.data.userList;
            userList.forEach(function(user){
                $userList.append(renderHtmlUserRow(user));
            })

            var $pagination = $("#pagination");
            $pagination.append(mininet.renderHtmlPagination(rsp.data.totalPageCount, params.pageIndex, params.pageSize));
        } else {
            // TODO 非正常处理
        }
    })
})

function renderHtmlUserRow(user){
    return '<tr class="first">' + 
        '<td><img src="' + user.avatar + '" class="img-circle avatar hidden-phone" />' +
            '<a href="user-profile.html" class="name" style="margin-left: 15px">' + user.nickname + '</a>' + 
        '</td>' +
        '<td>' + mininet.formatGender(user.gender) + '</td>' + 
        '<td>' + user.phone +'</td>' +
        '<td>￥ ' + (user.balance / 100).toFixed(2) + '元</td>' +
        '<td>' + mininet.formatPlat(user.plat) + '</td>' +
        '<td class="align-right">' + user.createTime + '</td>' +
    '</tr>'
}