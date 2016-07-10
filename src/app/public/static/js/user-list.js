$(document).ready(function(){
    mininet.renderHtmlNavbar('user');

    var $userList = $("#userList");

    mininet.ajax("get", "/user/list", {
        pageSize: 10,
        pageIndex: 1
    }, function(rsp){
        if (rsp.ret == 0){
            var userList = rsp.data;
            userList.forEach(function(user){
                $userList.append(renderHtmlUserRow(user));
            })
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