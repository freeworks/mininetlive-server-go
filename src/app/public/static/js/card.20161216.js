function Event(t) {
    if (t.originalEvent.changedTouches) {
        var e = t.originalEvent.changedTouches[0];
        this.y = e.pageY,
        this.x = e.pageX
    } else this.y = t.originalEvent.y,
    this.x = t.originalEvent.x;
    this.time = t.timeStamp
}
function Gesture(t) {
    function e(t) {
        return Math.abs(t)
    }
    this.element = t,
    this.mLastMotionX,
    this.mLastMotionY,
    this.mCurrentEvent,
    this.mFirstEvent,
    this.mIsScrolling,
    this.started = !1,
    this.mTouchSlopSquare = 900 * window.devicePixelRatio,
    this.getAngle = function(t, e) {
        var i = e.x - t.x,
        n = e.y - t.y;
        return 180 * Math.atan2(n, i) / Math.PI
    },
    this.getDirection = function(t, i) {
        return t === i ? "none": e(t) >= e(i) ? t < 0 ? "left": "right": i < 0 ? "up": "down"
    },
    this.getVelocity = function(t, e) {
        var i = e.x - t.x,
        n = e.y - t.y,
        o = e.time - t.time;
        if (0 == o) var a = 0,
        s = 0;
        else var a = 1e3 * i / o || 0,
        s = 1e3 * n / o || 0;
        return {
            x: a,
            y: s,
            abs: Math.sqrt(a * a + s * s),
            direction: this.getDirection(a, s)
        }
    };
    var i = this;
    this._handleTouchStart = function(t) {
        i.onStart(new Event(t))
    },
    this._handleTouchMove = function(t) {
        i.started && (t.preventDefault(), i.onMove(new Event(t)))
    },
    this._handleTouchEnd = function(t) {
        i.started && i.onEnd(new Event(t))
    },
    $(t).on("mousedown touchstart", this._handleTouchStart).on("mousemove touchmove", this._handleTouchMove).on("mouseup touchend touchcancel", this._handleTouchEnd)
}
function switchAlbum(t) {
    var e = $album.$imgs.length;
    if (! (e <= 1)) {
        var i = $album.$imgs.pop();
        i.css3({
            transform: t,
            transition: "all 0.6s",
            opacity: 0
        });
        var e = $album.$imgs.length;
        i.one("webkitTransitionEnd transitionEnd",
        function() {
            $album.prepend(i),
            $album.$imgs.unshift(i),
            i.css3({
                transform: "rotate(" + 5 * e + "deg)",
                opacity: 1
            })
        });
        var n = 5 * e,
        o = .2 * e;
        $.each($album.$imgs,
        function(t, e) {
            n -= 5,
            e.css3({
                transform: "rotate(" + n + "deg)",
                transition: "all 0.2s " + o + "s",
                opacity: 1
            }),
            o -= .2
        })
    }
}
function switchAlbumNext() {
    switchAlbum("translateX(-140%)")
}
function switchAlbumBack() {
    switchAlbum("translateX(140%)")
}
function cpage(t) {
    t = t || {},
    t.type = t.type || 2,
    t.id = t.id || "",
    t.bgcolor = t.bgcolor || "#f3f3f3";
    var e;
    switch (t.type) {
    case 4:
        e = pageTpl4,
        e.css({
            opacity: 1
        }),
        e.show();
        break;
    case 5:
        return;
    default:
        e = $('<div class="page"></div>')
    }
    switch (e.attr("page_id", t.id), e.type = t.type, t.bgimg ? e.css({
        "background-image": 'url("' + t.bgimg + '")',
        "background-color": t.bgcolor,
        "background-size": "contain",
        "background-repeat": "no-repeat"
    }) : e.css("background", t.bgcolor), e.beforeHide = $.noop, e.beforeShow = $.noop, e.afterHide = $.noop, e.afterShow = $.noop, e.swipeUp = $.noop, e.swipeLeft = $.noop, e.swipeRight = $.noop, e.swipeDown = $.noop, t.type) {
    case 3:
        $album = $('<div class="stack-gallery" id="album"></div>');
        var i = $("html body").height();
        if ($album.css({
            "line-height": i + "px"
        }), $album.$imgs = [], e.append($album), $.each(t.album,
        function(t, e) {
            var i = $('<div class="wrap"></div>'),
            n = $('<img src="' + e + '" />');
            n.css3({
                transform: "scale(0.7)"
            }),
            i.append(n),
            $album.append(i),
            $album.$imgs.push(i),
            i.css3({
                transform: "translateX(-130%) translateY(60%) rotate(-30deg)"
            })
        }), $album.total = t.album.length, $album.total > 1) {
            var n = $album.$leftArrow = $('<img src="http://ecard.source.izhaowo.com/s5_arrow_h.png" style="width: 20px; height: 32px; visibility: visible; z-index: 10000; opacity: 1; transition: 0.8s; position: absolute; left: 20px; top: 50%;">');
            e.append(n);
            var o = $album.$rightArrow = $('<img src="http://ecard.source.izhaowo.com/s5_arrow_h.png" style="width: 20px; height: 32px; visibility: visible; z-index: 10000; opacity: 1; transition: 0.8s; position: absolute; right: 20px; top: 50%;">');
            o.css3({
                transform: "rotate(180deg)"
            }),
            e.append(o)
        }
        e.afterShow = function() {
            $album.addClass("active");
            var t = $album.$imgs.length,
            e = 5 * t;
            $.each($album.$imgs,
            function(t, i) {
                e -= 5,
                i.css3({
                    transform: "rotate(" + e + "deg)",
                    transition: "all 0.6s " + .3 * t + "s"
                })
            });
            var i = $album.$imgs[t - 1];
            $album.total > 1 && i.one("webkitTransitionEnd transitionEnd",
            function() {
                n.addClass("keyframe40"),
                o.addClass("keyframe41")
            })
        },
        e.swipeUp = switchNext,
        e.swipeDown = switchBack,
        e.swipeLeft = switchAlbumNext,
        e.swipeRight = switchAlbumBack;
        break;
    case 4:
        e.find(".u-weddate").text(t.weddate),
        e.find(".u-message").text(t.message),
        e.find("#avatar").attr("src", avatar),
        e.find("#address").text("" + t.address);
        var a = window.map,
        s = window.longitude,
        r = window.latitude;
        if (a) {
            e.find("img.navi").attr("src", a + "&size=542*268");
            var l = navigator.userAgent.toLowerCase();
            e.find(".navi").click(function() {
                "micromessenger" == l.match(/MicroMessenger/i) ? (r = parseFloat(r), s = parseFloat(s), wx.openLocation({
                    latitude: r,
                    longitude: s,
                    name: t.address,
                    address: "酒店"
                })) : location.href = "http://m.amap.com/navi/?dest=" + s + "," + r + "&destName=" + t.address + "&key=4828d3efad639020a1d175b68009bfa1"
            })
        } else e.find(".navi").hide();
        e.afterShow = function() {
            $(".u-guideWrap").hide()
        },
        e.swipeUp = switchNext,
        e.swipeDown = switchBack;
        break;
    default:
        e.swipeUp = switchNext,
        e.swipeDown = switchBack
    }
    switch (e.beforeShow = function() {
        var t = this;
        t.show()
    },
    e.afterHide = function() {
        var t = this;
        t.hide()
    },
    t.type) {
    case 3:
        e.beforeHide = function() {
            $album.total > 1 && ($album.$leftArrow.removeClass("keyframe40"), $album.$rightArrow.removeClass("keyframe41")),
            $.each($album.$imgs,
            function(t, e) {
                e.css3({
                    transform: "translateX(-130%) translateY(60%) rotate(-30deg)",
                    transition: "none"
                })
            })
        };
        break;
    case 4:
        e.beforeHide = function() {
            $(".u-guideWrap").show()
        }
    }
    return e
}
function preparePages() {
    var t = [];
    if ($.each(window.data,
    function(e, i) {
        var n = cpage(i);
        n && t.push(n)
    }), $.each(t,
    function(e, i) {
        var n = void 0,
        o = void 0;
        0 != e && (n = t[e - 1]),
        e != t.length - 1 && (o = t[e + 1]),
        i.prev = n,
        i.next = o
    }), t.length >= 2) {
        var e = t[0],
        i = t[t.length - 1];
        e.prev = i,
        i.next = e
    }
    return t
}
function getPage() {
    return $page
}
function getNextPage() {
    return $page.next
}
function getPrevPage() {
    return $page.prev
}
function slideZoomNext() {
    ready = !1;
    var t = getPage(),
    e = getNextPage();
    $viewport.append(e),
    t.addClass("flipOutUp"),
    e.addClass("flipInDown"),
    t.beforeHide(),
    e.beforeShow(),
    setTimeout(function() {
        t.removeClass("flipOutUp"),
        e.removeClass("flipInDown"),
        t.afterHide(),
        e.afterShow(),
        $page = e,
        ready = !0
    },
    800)
}
function slideZoomBack() {
    ready = !1;
    var t = getPage(),
    e = getPrevPage();
    $viewport.append(e),
    t.addClass("flipOutDown"),
    e.addClass("flipInUp"),
    t.beforeHide(),
    e.beforeShow(),
    setTimeout(function() {
        t.removeClass("flipOutDown"),
        e.removeClass("flipInUp"),
        t.afterHide(),
        e.afterShow(),
        $page = e,
        ready = !0
    },
    800)
}
function cacheCompeleted() {
    var t = preparePages();
    $page = t[0],
    $loading.fadeOut("slow",
    function() {
        if ($page = getPage(), $page.beforeShow(), $viewport.append($page), $page.afterShow(), music) {
            var t = $("#audioplayer")[0];
            t.src = music,
            $window.trigger("musicon")
        }
    })
}
var music = window.music,
avatar = window.avatar,
flexible = function() {},
windowW = window.screen.width,
windowH = window.screen.height;
flexible.init = function(t) {
    t = t || {};
    var e = (t.pageWidth || 640, 320),
    i = 520,
    n = document.getElementsByTagName("body")[0],
    o = function() {
        var t = n.clientWidth;
        t > e && windowW > windowH && (n.style.maxWidth = e + "px", $("body").css({
            maxHeight: i + "px",
            top: ($(window).height() - i) / 2 + "px",
            width: windowW + "px"
        })),
        $(n).css({
            visibility: "visible"
        }),
        $("html").css({
            background: "black"
        })
    };
    window.onresize = o()
},
String.prototype.hashCode = function() {
    var t = 0;
    if (0 == this.length) return t;
    for (var e = 0; e < this.length; e++) {
        var i = this.charCodeAt(e);
        t = (t << 5) - t + i,
        t &= t
    }
    return t
};
var lvl;
if (window.console) lvl = {
    n: function() {},
    l: function(t) {
        window.console.log(t)
    },
    i: function(t) {
        window.console.info(t)
    },
    d: function(t) {
        window.console.debug(t)
    },
    w: function(t) {
        window.console.warn(t)
    },
    e: function(t) {
        window.console.error(t)
    }
};
else {
    var print = function(t) {
        alert(t)
    };
    lvl = {
        n: function() {},
        l: print,
        i: print,
        d: print,
        w: print,
        e: print
    }
}
var log = function(t, e) {
    return "undefined" == typeof e ? void lvl.l(t) : void e(t)
};
log.level = lvl.d,
log.i = function(t) {
    log(t, lvl.i)
},
log.d = function(t) {
    log(t, lvl.d)
},
log.w = function(t) {
    log(t, lvl.w)
},
log.e = function(t) {
    log(t, lvl.e)
};
var native = function() {
    this.status = 0,
    this.options = {}
};
native.prototype.ready = function(t, e) {
    var i = this,
    n = navigator.userAgent,
    o = n.indexOf("zwios") > -1,
    a = n.indexOf("zwandroid") > -1;
    if (!o && !a) return void(e && setTimeout(e.faild, 0));
    i.options = {
        plantform: o ? "ios": a ? "android": "",
        token: "",
        userId: "",
        phone: "",
        build: ""
    },
    "undefined" == typeof window.ZWCallbacks && (window.ZWCallbacks = []);
    var s = function() {
        "ready" !== i.status && (i.status = i.status + 1, i.status < 3 ? (i.callNative("ready", {},
        function(t, n) {
            return n ? (i.status = "ready", i.options.token = n.token, i.options.userId = n.userId, i.options.phone = n.phone, i.options.build = n.build, void(e && setTimeout(e.success, 0))) : void(e && setTimeout(e.faild, 0))
        }), setTimeout(s, 1e3)) : 3 == i.status && e && setTimeout(e.faild, 0))
    };
    s()
};
var CallAndroidNative = function(t, e, i) {
    i = i ||
    function() {};
    var n = encodeURIComponent(JSON.stringify(e)),
    o = (t + n).hashCode().toString();
    window.ZWCallbacks[o] = i,
    log("call native function:" + t),
    window.Native.call(t, JSON.stringify(e), o)
},
CallIOSNative = function(t, e, i) {
    i = i ||
    function() {};
    var n = encodeURIComponent(JSON.stringify(e)),
    o = (t + n).hashCode().toString();
    window.ZWCallbacks[o] = i;
    var a = "zwscheme://nativeapi?cb=" + o + "&method=" + t + "&data=" + n,
    s = document.createElement("iframe");
    s.style.display = "none",
    s.src = a,
    setTimeout(function() {
        document.documentElement.appendChild(s),
        setTimeout(function() {
            document.documentElement.removeChild(s)
        },
        0)
    },
    0)
};
native.prototype.callNative = function(t, e, i) {
    var n = this;
    return "android" === n.options.plantform ? CallAndroidNative(t, e, i) : "ios" === n.options.plantform ? CallIOSNative(t, e, i) : log.w("native not bind yet!")
},
Gesture.prototype = {
    onStart: function(t) {
        this.started = !0;
        var e = t.y,
        i = t.x;
        this.mLastMotionX = i,
        this.mLastMotionY = e,
        this.mCurrentEvent = t,
        this.mFirstEvent = t,
        this.mIsScrolling = !1,
        $(this.element).trigger("pressdown", [t])
    },
    onMove: function(t) {
        var e = t.y,
        i = t.x;
        if (!this.mIsScrolling) {
            var n = i - this.mFirstEvent.x,
            o = e - this.mFirstEvent.y,
            a = n * n + o * o;
            a > this.mTouchSlopSquare && (this.mIsScrolling = !0, $(this.element).trigger("scroll", [t, this.getDirection(n, o)]))
        }
        var n = i - this.mCurrentEvent.x,
        o = e - this.mCurrentEvent.y;
        this.mIsScrolling || $(this.element).trigger("pressmove", [t, n, o]),
        this.mLastMotionX = i,
        this.mLastMotionY = e,
        this.mCurrentEvent = t
    },
    onEnd: function(t) {
        this.started = !1,
        $(this.element).trigger("pressup", [t])
    },
    destroy: function() {
        this.started = !1,
        $(this.element).off("mousedown touchstart", this._handleTouchStart).off("mousemove touchmove", this._handleTouchMove).off("mouseup touchend touchcancel", this._handleTouchEnd)
    },
    on: function(t, e) {
        return $(this.element).on(t, e),
        this
    }
},
$.fn.extend({
    css3: function(t) {
        var e = t || {};
        return $.each(e,
        function(t, i) {
            e["-webkit-" + t] = i
        }),
        this.each(function() {
            var t = $(this);
            t.css(e)
        })
    }
}),
$.cacheImage = function(t, e) {
    var i = new Image;
    return i.src = t,
    i.complete ? void e.call(i) : void(i.onabort = i.onerror = i.onload = function() {
        e.call(i)
    })
},
$.cacheImages = function(t, e) {
    var i = e || {};
    i.one = i.one || $.noop,
    i.all = i.all || $.noop;
    var n = t || [];
    if (0 == n.length) return void i.all.call(n);
    var o = [];
    $.each(n,
    function(e, a) {
        $.cacheImage(a,
        function() {
            var e = this;
            o.push(e),
            n.length == o.length ? i.all.call(e, o, t) : i.one.call(e, o, t)
        })
    })
},
$.alert = function(t, e) {
    var i = '<div class="m-layer z-show"><table><tbody><tr><td><article class="lywrap"><section class="lyct"></section><footer class="lybt"><div class="lybtns"><button type="button" class="u-btn">确定</button></div></footer></article></td></tr></tbody></table></div>',
    n = $(i);
    n.find(".lyct").html("<p>" + t + "</p>"),
    n.find(".u-btn").click(function() {
        n.remove(),
        e || e()
    }),
    $(document.body).append(n)
};
var pageTpl4, $page, $loading = $("#loading"),
$viewport = $(".viewport"),
$window = $(window),
$album,
ready = !0,
H = $("html body").height(),
scale = 1e3 / (1e3 + .5 * H),
dest = 1e3 * H * .5 / (1e3 + .5 * H),
switchNext = slideZoomNext,
switchBack = slideZoomBack;
music && $window.on("togglemusic",
function() {
    var t = $("#audioplayer")[0];
    t.paused ? (t.play(), $("#audio").addClass("active")) : ($("#audio").removeClass("active"), t.pause())
}).on("musicoff",
function() {
    var t = $("#audioplayer")[0];
    t.pause(),
    $("#audio").removeClass("active")
}).on("musicon",
function() {
    $("#audioplayer")[0].play(),
    $("#audio").addClass("active")
});
var fal = 0,
gesture = new Gesture($viewport).on("scroll",
function(t, e, i) {
    if (ready) try {
        "up" == i ? getPage().swipeUp() : "down" == i ? getPage().swipeDown() : "left" == i ? getPage().swipeLeft() : "right" == i && getPage().swipeRight()
    } catch(t) {
        throw ready = !0,
        t
    }
}).on("touchstart",
function() {
    0 == fal && ($("#audioplayer")[0].play(), fal = 1)
});
$(function() {
    window.onload = function() {
        var t = new native;
        pageTpl4 = $("#lastpage");
        var e = function() {
            flexible.init(),
            "ios" !== t.options.plantform && "android" !== t.options.plantform ? ($("#lastpage .createbtn").html("我要制作"), pageTpl4.remove()) : ($(".createbtn").remove(), pageTpl4.remove());
            var e = window.imgs || [];
            if ($.cacheImages(e, {
                one: function(t, e) {
                    $("#progress").text(parseInt(100 * t.length / e.length) + "%")
                },
                all: function(t, e) {
                    $("#progress").text("100%"),
                    cacheCompeleted()
                }
            }), windowW < windowH) {
                var i, n = parseInt(parseInt($(".page").css("width")) / 750 * 1217);
                n > $(window).height() ? (n = $(window).height(), i = parseInt(750 * n / 1217), $("body").css({
                    width: i + "px",
                    height: n + "px",
                    margin: "0 auto"
                }), $(".page").css({
                    width: i + "px",
                    height: n + "px"
                })) : (i = parseInt(750 * n / 1217), $("body").css({
                    width: i + "px",
                    height: n + "px"
                }), $(".page").css({
                    width: i + "px",
                    height: n + "px"
                }))
            }
        };
        t.ready(window.appinfo, {
            success: e,
            faild: e
        })
    }
});