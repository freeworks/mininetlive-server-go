function wxready() {
    $.ajax({
        url: "/wxpub/jsconfig",
        type: "POST",
        contentType: "application/x-www-form-urlencoded",
        data: {
            url: window.location.href
        },
        error: function(e, n, i) {
            debugger
            wx_debug && alert("网络错误：" + n)
        },
        success: function(resp) {
            if (resp && 0 == resp.ret) {
               load_wx({
                debug: false,
                appId: resp.data.appId,
                timestamp: resp.data.timestamp, 
                nonceStr: resp.data.nonceStr,
                signature: resp.data.signature,
                jsApiList: ["onMenuShareTimeline", "onMenuShareAppMessage"]
            })
            } else wx_debug && alert("初始化参数错误：" + e.code + ":" + e.desc)
        }
    })
}
function load_wx(e) {
    wx.config(e),
    wx.error(function(e) {
        wx_debug && alert(JSON.stringify(e))
    }),
    wx.ready(function() {
        wx_debug && alert("微信分享初始化成功"),
        wx.onMenuShareAppMessage({
            title: shareTitle,
            desc: descContent,
            link: message_link,
            imgUrl: imgUrl,
            trigger: function(e) {},
            success: function(e) {
                shareback()
            },
            cancel: function(e) {},
            fail: function(e) {}
        }),
        wx.onMenuShareTimeline({
            title: descContent,
            link: message_link,
            imgUrl: imgUrl,
            trigger: function(e) {},
            success: function(e) {
                shareback()
            },
            cancel: function(e) {},
            fail: function(e) {}
        }),
        wx.onMenuShareQQ({
            title: shareTitle,
            desc: descContent,
            link: message_link,
            imgUrl: imgUrl,
            trigger: function(e) {},
            complete: function(e) {},
            success: function(e) {
                shareback()
            },
            cancel: function(e) {},
            fail: function(e) {}
        }),
        wx.onMenuShareWeibo({
            title: shareTitle,
            desc: descContent,
            link: message_link,
            imgUrl: imgUrl,
            trigger: function(e) {},
            complete: function(e) {},
            success: function(e) {
                shareback()
            },
            cancel: function(e) {},
            fail: function(e) {}
        }),
        wx.checkJsApi({
            jsApiList: e.jsApiList,
            success: function(e) {}
        })
    })
} !
function(e, n) {
    "function" == typeof define && (define.amd || define.cmd) ? define(function() {
        return n(e)
    }) : n(e, !0)
} (this,
function(e, n) {
    function i(n, i, t) {
        e.WeixinJSBridge ? WeixinJSBridge.invoke(n, o(i),
        function(e) {
            c(n, e, t)
        }) : d(n, t)
    }
    function t(n, i, t) {
        e.WeixinJSBridge ? WeixinJSBridge.on(n,
        function(e) {
            t && t.trigger && t.trigger(e),
            c(n, e, i)
        }) : t ? d(n, t) : d(n, i)
    }
    function o(e) {
        return e = e || {},
        e.appId = x.appId,
        e.verifyAppId = x.appId,
        e.verifySignType = "sha1",
        e.verifyTimestamp = x.timestamp + "",
        e.verifyNonceStr = x.nonceStr,
        e.verifySignature = x.signature,
        e
    }
    function r(e) {
        return {
            timeStamp: e.timestamp + "",
            nonceStr: e.nonceStr,
            "package": e["package"],
            paySign: e.paySign,
            signType: e.signType || "SHA1"
        }
    }
    function c(e, n, i) {
        var t, o, r;
        switch (delete n.err_code, delete n.err_desc, delete n.err_detail, t = n.errMsg, t || (t = n.err_msg, delete n.err_msg, t = s(e, t), n.errMsg = t), i = i || {},
        i._complete && (i._complete(n), delete i._complete), t = n.errMsg || "", x.debug && !i.isInnerInvoke && alert(JSON.stringify(n)), o = t.indexOf(":"), r = t.substring(o + 1)) {
        case "ok":
            i.success && i.success(n);
            break;
        case "cancel":
            i.cancel && i.cancel(n);
            break;
        default:
            i.fail && i.fail(n)
        }
        i.complete && i.complete(n)
    }
    function s(e, n) {
        var i, t, o = e,
        r = g[o];
        return r && (o = r),
        i = "ok",
        n && (t = n.indexOf(":"), i = n.substring(t + 1), "confirm" == i && (i = "ok"), "failed" == i && (i = "fail"), -1 != i.indexOf("failed_") && (i = i.substring(7)), -1 != i.indexOf("fail_") && (i = i.substring(5)), i = i.replace(/_/g, " "), i = i.toLowerCase(), ("access denied" == i || "no permission to execute" == i) && (i = "permission denied"), "config" == o && "function not exist" == i && (i = "ok"), "" == i && (i = "fail")),
        n = o + ":" + i
    }
    function a(e) {
        var n, i, t, o;
        if (e) {
            for (n = 0, i = e.length; i > n; ++n) t = e[n],
            o = m[t],
            o && (e[n] = o);
            return e
        }
    }
    function d(e, n) {
        if (! (!x.debug || n && n.isInnerInvoke)) {
            var i = g[e];
            i && (e = i),
            n && n._complete && delete n._complete,
            console.log('"' + e + '",', n || "")
        }
    }
    function u() {
        if (! ("6.0.2" > T || I.systemType < 0)) {
            var e = new Image;
            I.appId = x.appId,
            I.initTime = v.initEndTime - v.initStartTime,
            I.preVerifyTime = v.preVerifyEndTime - v.preVerifyStartTime,
            V.getNetworkType({
                isInnerInvoke: !0,
                success: function(n) {
                    I.networkType = n.networkType;
                    var i = "https://open.weixin.qq.com/sdk/report?v=" + I.version + "&o=" + I.isPreVerifyOk + "&s=" + I.systemType + "&c=" + I.clientVersion + "&a=" + I.appId + "&n=" + I.networkType + "&i=" + I.initTime + "&p=" + I.preVerifyTime + "&u=" + I.url;
                    e.src = i
                }
            })
        }
    }
    function l() {
        return (new Date).getTime()
    }
    function f(n) {
        _ && (e.WeixinJSBridge ? n() : h.addEventListener && h.addEventListener("WeixinJSBridgeReady", n, !1))
    }
    function p() {
        V.invoke || (V.invoke = function(n, i, t) {
            e.WeixinJSBridge && WeixinJSBridge.invoke(n, o(i), t)
        },
        V.on = function(n, i) {
            e.WeixinJSBridge && WeixinJSBridge.on(n, i)
        })
    }
    var m, g, h, w, y, _, S, k, T, v, I, x, b, M, V;
    return e.jWeixin ? void 0 : (m = {
        config: "preVerifyJSAPI",
        onMenuShareTimeline: "menu:share:timeline",
        onMenuShareAppMessage: "menu:share:appmessage",
        onMenuShareQQ: "menu:share:qq",
        onMenuShareWeibo: "menu:share:weiboApp",
        onMenuShareQZone: "menu:share:QZone",
        previewImage: "imagePreview",
        getLocation: "geoLocation",
        openProductSpecificView: "openProductViewWithPid",
        addCard: "batchAddCard",
        openCard: "batchViewCard",
        chooseWXPay: "getBrandWCPayRequest"
    },
    g = function() {
        var e, n = {};
        for (e in m) n[m[e]] = e;
        return n
    } (), h = e.document, w = h.title, y = navigator.userAgent.toLowerCase(), _ = -1 != y.indexOf("micromessenger"), S = -1 != y.indexOf("android"), k = -1 != y.indexOf("iphone") || -1 != y.indexOf("ipad"), T = function() {
        var e = y.match(/micromessenger\/(\d+\.\d+\.\d+)/) || y.match(/micromessenger\/(\d+\.\d+)/);
        return e ? e[1] : ""
    } (), v = {
        initStartTime: l(),
        initEndTime: 0,
        preVerifyStartTime: 0,
        preVerifyEndTime: 0
    },
    I = {
        version: 1,
        appId: "",
        initTime: 0,
        preVerifyTime: 0,
        networkType: "",
        isPreVerifyOk: 1,
        systemType: k ? 1 : S ? 2 : -1,
        clientVersion: T,
        url: encodeURIComponent(location.href)
    },
    x = {},
    b = {
        _completes: []
    },
    M = {
        state: 0,
        res: {}
    },
    f(function() {
        v.initEndTime = l()
    }), V = {
        config: function(e) {
            x = e,
            d("config", e);
            var n = x.check !== !1;
            f(function() {
                var e, t, o;
                if (n) i(m.config, {
                    verifyJsApiList: a(x.jsApiList)
                },
                function() {
                    b._complete = function(e) {
                        v.preVerifyEndTime = l(),
                        M.state = 1,
                        M.res = e
                    },
                    b.success = function() {
                        I.isPreVerifyOk = 0
                    },
                    b.fail = function(e) {
                        b._fail ? b._fail(e) : M.state = -1
                    };
                    var e = b._completes;
                    return e.push(function() {
                        x.debug || u()
                    }),
                    b.complete = function() {
                        for (var n = 0,
                        i = e.length; i > n; ++n) e[n]();
                        b._completes = []
                    },
                    b
                } ()),
                v.preVerifyStartTime = l();
                else {
                    for (M.state = 1, e = b._completes, t = 0, o = e.length; o > t; ++t) e[t]();
                    b._completes = []
                }
            }),
            x.beta && p()
        },
        ready: function(e) {
            0 != M.state ? e() : (b._completes.push(e), !_ && x.debug && e())
        },
        error: function(e) {
            "6.0.2" > T || ( - 1 == M.state ? e(M.res) : b._fail = e)
        },
        checkJsApi: function(e) {
            var n = function(e) {
                var n, i, t = e.checkResult;
                for (n in t) i = g[n],
                i && (t[i] = t[n], delete t[n]);
                return e
            };
            i("checkJsApi", {
                jsApiList: a(e.jsApiList)
            },
            function() {
                return e._complete = function(e) {
                    if (S) {
                        var i = e.checkResult;
                        i && (e.checkResult = JSON.parse(i))
                    }
                    e = n(e)
                },
                e
            } ())
        },
        onMenuShareTimeline: function(e) {
            t(m.onMenuShareTimeline, {
                complete: function() {
                    i("shareTimeline", {
                        title: e.title || w,
                        desc: e.title || w,
                        img_url: e.imgUrl || "",
                        link: e.link || location.href,
                        type: e.type || "link",
                        data_url: e.dataUrl || ""
                    },
                    e)
                }
            },
            e)
        },
        onMenuShareAppMessage: function(e) {
            t(m.onMenuShareAppMessage, {
                complete: function() {
                    i("sendAppMessage", {
                        title: e.title || w,
                        desc: e.desc || "",
                        link: e.link || location.href,
                        img_url: e.imgUrl || "",
                        type: e.type || "link",
                        data_url: e.dataUrl || ""
                    },
                    e)
                }
            },
            e)
        },
        onMenuShareQQ: function(e) {
            t(m.onMenuShareQQ, {
                complete: function() {
                    i("shareQQ", {
                        title: e.title || w,
                        desc: e.desc || "",
                        img_url: e.imgUrl || "",
                        link: e.link || location.href
                    },
                    e)
                }
            },
            e)
        },
        onMenuShareWeibo: function(e) {
            t(m.onMenuShareWeibo, {
                complete: function() {
                    i("shareWeiboApp", {
                        title: e.title || w,
                        desc: e.desc || "",
                        img_url: e.imgUrl || "",
                        link: e.link || location.href
                    },
                    e)
                }
            },
            e)
        },
        onMenuShareQZone: function(e) {
            t(m.onMenuShareQZone, {
                complete: function() {
                    i("shareQZone", {
                        title: e.title || w,
                        desc: e.desc || "",
                        img_url: e.imgUrl || "",
                        link: e.link || location.href
                    },
                    e)
                }
            },
            e)
        },
        startRecord: function(e) {
            i("startRecord", {},
            e)
        },
        stopRecord: function(e) {
            i("stopRecord", {},
            e)
        },
        onVoiceRecordEnd: function(e) {
            t("onVoiceRecordEnd", e)
        },
        playVoice: function(e) {
            i("playVoice", {
                localId: e.localId
            },
            e)
        },
        pauseVoice: function(e) {
            i("pauseVoice", {
                localId: e.localId
            },
            e)
        },
        stopVoice: function(e) {
            i("stopVoice", {
                localId: e.localId
            },
            e)
        },
        onVoicePlayEnd: function(e) {
            t("onVoicePlayEnd", e)
        },
        uploadVoice: function(e) {
            i("uploadVoice", {
                localId: e.localId,
                isShowProgressTips: 0 == e.isShowProgressTips ? 0 : 1
            },
            e)
        },
        downloadVoice: function(e) {
            i("downloadVoice", {
                serverId: e.serverId,
                isShowProgressTips: 0 == e.isShowProgressTips ? 0 : 1
            },
            e)
        },
        translateVoice: function(e) {
            i("translateVoice", {
                localId: e.localId,
                isShowProgressTips: 0 == e.isShowProgressTips ? 0 : 1
            },
            e)
        },
        chooseImage: function(e) {
            i("chooseImage", {
                scene: "1|2",
                count: e.count || 9,
                sizeType: e.sizeType || ["original", "compressed"],
                sourceType: e.sourceType || ["album", "camera"]
            },
            function() {
                return e._complete = function(e) {
                    if (S) {
                        var n = e.localIds;
                        n && (e.localIds = JSON.parse(n))
                    }
                },
                e
            } ())
        },
        previewImage: function(e) {
            i(m.previewImage, {
                current: e.current,
                urls: e.urls
            },
            e)
        },
        uploadImage: function(e) {
            i("uploadImage", {
                localId: e.localId,
                isShowProgressTips: 0 == e.isShowProgressTips ? 0 : 1
            },
            e)
        },
        downloadImage: function(e) {
            i("downloadImage", {
                serverId: e.serverId,
                isShowProgressTips: 0 == e.isShowProgressTips ? 0 : 1
            },
            e)
        },
        getNetworkType: function(e) {
            var n = function(e) {
                var n, i, t, o = e.errMsg;
                if (e.errMsg = "getNetworkType:ok", n = e.subtype, delete e.subtype, n) e.networkType = n;
                else switch (i = o.indexOf(":"), t = o.substring(i + 1)) {
                case "wifi":
                case "edge":
                case "wwan":
                    e.networkType = t;
                    break;
                default:
                    e.errMsg = "getNetworkType:fail"
                }
                return e
            };
            i("getNetworkType", {},
            function() {
                return e._complete = function(e) {
                    e = n(e)
                },
                e
            } ())
        },
        openLocation: function(e) {
            i("openLocation", {
                latitude: e.latitude,
                longitude: e.longitude,
                name: e.name || "",
                address: e.address || "",
                scale: e.scale || 28,
                infoUrl: e.infoUrl || ""
            },
            e)
        },
        getLocation: function(e) {
            e = e || {},
            i(m.getLocation, {
                type: e.type || "wgs84"
            },
            function() {
                return e._complete = function(e) {
                    delete e.type
                },
                e
            } ())
        },
        hideOptionMenu: function(e) {
            i("hideOptionMenu", {},
            e)
        },
        showOptionMenu: function(e) {
            i("showOptionMenu", {},
            e)
        },
        closeWindow: function(e) {
            e = e || {},
            i("closeWindow", {
                immediate_close: e.immediateClose || 0
            },
            e)
        },
        hideMenuItems: function(e) {
            i("hideMenuItems", {
                menuList: e.menuList
            },
            e)
        },
        showMenuItems: function(e) {
            i("showMenuItems", {
                menuList: e.menuList
            },
            e)
        },
        hideAllNonBaseMenuItem: function(e) {
            i("hideAllNonBaseMenuItem", {},
            e)
        },
        showAllNonBaseMenuItem: function(e) {
            i("showAllNonBaseMenuItem", {},
            e)
        },
        scanQRCode: function(e) {
            e = e || {},
            i("scanQRCode", {
                needResult: e.needResult || 0,
                scanType: e.scanType || ["qrCode", "barCode"]
            },
            function() {
                return e._complete = function(e) {
                    var n, i;
                    k && (n = e.resultStr, n && (i = JSON.parse(n), e.resultStr = i && i.scan_code && i.scan_code.scan_result))
                },
                e
            } ())
        },
        openProductSpecificView: function(e) {
            i(m.openProductSpecificView, {
                pid: e.productId,
                view_type: e.viewType || 0,
                ext_info: e.extInfo
            },
            e)
        },
        addCard: function(e) {
            var n, t, o, r, c = e.cardList,
            s = [];
            for (n = 0, t = c.length; t > n; ++n) o = c[n],
            r = {
                card_id: o.cardId,
                card_ext: o.cardExt
            },
            s.push(r);
            i(m.addCard, {
                card_list: s
            },
            function() {
                return e._complete = function(e) {
                    var n, i, t, o = e.card_list;
                    if (o) {
                        for (o = JSON.parse(o), n = 0, i = o.length; i > n; ++n) t = o[n],
                        t.cardId = t.card_id,
                        t.cardExt = t.card_ext,
                        t.isSuccess = !!t.is_succ,
                        delete t.card_id,
                        delete t.card_ext,
                        delete t.is_succ;
                        e.cardList = o,
                        delete e.card_list
                    }
                },
                e
            } ())
        },
        chooseCard: function(e) {
            i("chooseCard", {
                app_id: x.appId,
                location_id: e.shopId || "",
                sign_type: e.signType || "SHA1",
                card_id: e.cardId || "",
                card_type: e.cardType || "",
                card_sign: e.cardSign,
                time_stamp: e.timestamp + "",
                nonce_str: e.nonceStr
            },
            function() {
                return e._complete = function(e) {
                    e.cardList = e.choose_card_info,
                    delete e.choose_card_info
                },
                e
            } ())
        },
        openCard: function(e) {
            var n, t, o, r, c = e.cardList,
            s = [];
            for (n = 0, t = c.length; t > n; ++n) o = c[n],
            r = {
                card_id: o.cardId,
                code: o.code
            },
            s.push(r);
            i(m.openCard, {
                card_list: s
            },
            e)
        },
        chooseWXPay: function(e) {
            i(m.chooseWXPay, r(e), e)
        }
    },
    n && (e.wx = e.jWeixin = V), V)
});
var query = function(e) {
    var n = new RegExp("(^|&)" + e + "=([^&]*)(&|$)"),
    i = window.location.search.substr(1).match(n);
    return i ? unescape(i[2]) : void 0
},
// wx_debug = "undefined" != typeof query("debug"),
wx_debug = false;
shareback = function() {
//    var e = window.ecard + "/ecard/tj?t=1&tid=" + id;
//    $("html body").append($('<script src="' + e + '"></script>'))
},
viewback = function() {
//    var e = window.ecard + "/ecard/tj?t=2&tid=" + id;
//    $("html body").append($('<script src="' + e + '"></script>'))
};
wxready(),
viewback();