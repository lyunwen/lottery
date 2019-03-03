$(function () {

    var leaderArray = [
        {min: 'bobli', 'name': '李波'},
        {min: 'rickyyan', 'name': '闫敏'},
        {min: 'shawnli', 'name': '李晓佑'},
        {min: 'mandywang', 'name': '王影'},
        {min: 'hymanxie', 'name': '谢怀遥'},
        {min: 'silonli', 'name': '李善林'},
        {min: 'lynnning', 'name': '宁强'},
        {min: 'swordwu', 'name': '吴剑'},
        {min: 'greyqiao', 'name': '乔成'},
        {min: 'louisjiang', 'name': '江能沐'},
        {min: 'capableguo', 'name': '郭亿清'},
        {min: 'ianwu', 'name': '吴建国'},
        {min: 'zeikzheng', 'name': '郑旭'},
        {min: 'vickychi', 'name': '迟冠群'},
        {min: 'slarkzhang', 'name': '张梁'},
        {min: 'vickyhhli', 'name': '李慧慧'},
        {min: 'xiaomeizhang', 'name': '张小妹'},
        {min: 'rockenmeng', 'name': '孟同亮'}

    ];

    if (!Array.prototype.filter) {
        Array.prototype.filter = function (fun /*, thisArg */) {
            "use strict";

            if (this === void 0 || this === null)
                throw new TypeError();

            var t = Object(this);
            var len = t.length >>> 0;
            if (typeof fun !== "function") {
                throw new TypeError();
            }

            var res = [];
            var thisArg = arguments.length >= 2 ? arguments[1] : void 0;
            for (var i = 0; i < len; i++) {
                if (i in t) {
                    var val = t[i];
                    if (fun.call(thisArg, val, i, t))
                        res.push(val);
                }
            }

            return res;
        };
    }

    /**
     * 自动选择组件构造方法
     * @param options   默认参数
     * @constructor
     */
    function MinAutoComplete(options) {

        this.options = {
            targetElement: '.js-autocomplete-input',
            originData: [], //默认数据
            renderData: [] //临时渲染数据
        };

        this.options = $.extend({}, this.options, options);
    };

    MinAutoComplete.prototype = {
        init: function () {
            this.beforeMount();
            this.assignInputEvent();//input修改事件
            this.assignKeyEvent();//键盘事件
            this.assignChangeEvnet();//焦点事件
            this.assignClickEvent();//min选中事件
        },
        assignChangeEvnet: function () {
            var that = this,
                targetElement = that.options.targetElement;

            $(document).delegate(targetElement, 'blur', function () {
                var node = $(this).next();
                if (node.hasClass('drop')) {
                    //延迟消失，先触发点击
                    setTimeout(function () {
                        node.toggle();
                    }, 200)
                }
            });

            $(document).delegate(targetElement, 'focus', function () {
                var node = $(this).next();
                if (node.hasClass('drop')) {
                    node.toggle();
                }
            });
        },
        assignKeyEvent: function () {
            var targetElement = this.options.targetElement,
                url = this.options.url,
                that = this,
                target = this.options.targetElement.substring(1, this.options.targetElement.length);

            $(document).delegate(targetElement, 'keydown', function (event) {

                if (event.keyCode == 38) {
                    var ul = $(".js-min-item-" + target);
                    var node = ul.find('.active');
                    var index = node && node.index();
                    var next = index - 1;
                    if (next >= 0) {
                        node.removeClass('active');
                        ul.find('li').eq(next).addClass('active');

                        var scrollNodeNum = ul.scrollTop() / 30;//被卷去的节点个数

                        //如果被卷曲的元素个数等于当前元素index,
                        //那么说明当前被选中的元素已经是容器内的最上面一个元素，需要开始需要滚动滚动条
                        if (scrollNodeNum == index) {
                            ul.scrollTop((index - 1) * 30);
                        }
                    }

                    event.preventDefault();
                    event.stopPropagation();
                } else if (event.keyCode == 40) {
                    var ul = $(".js-min-item-" + target);
                    var node = ul.find('.active');
                    var index = node && node.index();
                    var next = index + 1;
                    var height = 30;
                    if (next < ul.find('li').length) {
                        node.removeClass('active');
                        ul.find('li').eq(next).addClass('active');

                        //当前元素索引 减去 容器内最大能包含的元素个数 - 1 (因为包含了本身，所以需要多减去1，如：最大元素个数为7，则是减去6),
                        //然后再加上 1
                        //然后乘以单个元素的高度，即为应该卷去的高度
                        var extra = (index - 6 + 1) * 30;
                        ul.scrollTop(extra);
                    }
                } else if (event.keyCode == 13) {

                    var node = $(".js-min-item-" + target).find('.active');
                    that.assignEnter(node, this);
                    event.preventDefault();
                    event.stopPropagation();
                }
            });
        },
        getInputMinVal: function (elem) {
            var originVal = $(elem).val(),
                originValArr = typeof originVal != "undefined" ? originVal.split(';') : [];

            return originValArr.pop();

        },
        assignInputEvent: function () {
            var targetElement = this.options.targetElement,
                that = this;

            $(document).delegate(targetElement, 'input propertychange', function () {

                var elem = this,
                    min = that.getInputMinVal(elem),
                    minStr = new String(min);//强制转换为字符串

                //长度检测,长度为1,请求后台数据，否则直接使用本地数据
                if (typeof min != "undefined") {
                    var currMin = that.getInputMinVal(elem);
                    that.match(currMin);
                    that.createHtml(elem);
                    that.render(elem);
                } else if (typeof min == "undefined" || min.length == 0) {
                    that.remove();
                }
            });
        },
        assignClickEvent: function () {
            var target = this.options.targetElement.substring(1, this.options.targetElement.length),
                that = this;
            $("body").delegate('.js-min-item-' + target + ' li', 'click', function () {
                that.assignSelect(this);
            });
        },
        assignSelect: function (elem) {
            var that = this,
                min = $(elem).data('min'),
                targetNode = $(that.options.targetElement + "[data-autoid='" + $(elem).data('target') + "']");//目标元素节点
            $(that.options.targetElement + "[data-autoid='" + $(elem).data('target') + "']").val(min);
            that.remove();
        },
        /**
         * 处理Enter事件
         * @param elem
         * @param target
         */
        assignEnter: function (elem, target) {
            var that = this,
                min = $(elem).data('min'),
                targetNode = $(target),//目标元素节点
                originVal = typeof targetNode.val() == "undefined" ? '' : targetNode.val(),
                originValArr = originVal.split(';');
            str = '';

            //弹出组后一个元素,放入选中的元素
            if (null != min) {
                originValArr.pop();
                originValArr.push(min);
            }
            if (that.options.multiple == true) {
                str = originValArr.join(';') + ";";
            } else {
                str = originValArr.join(';');
            }

            $(target).val(str);
            that.remove();
        },
        match: function (min) {
            var data = this.options.originData,
                targetElement = this.options.targetElement,
                min = new String(min),
                regex = new RegExp('^' + min, 'i');

            //filter兼容处理
            var filterData = data.filter(function (item, index, data) {
                var name = item['min'];
                return regex.test(name);
            });

            this.options.renderData = filterData;
        },
        render: function (elm) {
            //先移除
            this.remove();
            var targetElement = this.options.targetElement,
                that = this;

            $(elm).after(that.template);
        },
        beforeMount: function () {
            //预加载，处理节点，给第一个节点加一个data-autoid
            var options = this.options;
            $(options['targetElement']).each(function (index, item) {
                var target = options['targetElement'].substring(1, options['targetElement'].length);
                var virtualId = target + '_min_auto_complete' + index;
                $(this).attr('data-autoid', virtualId);
            });
        },
        remove: function () {
            var node = $(this.options.targetElement).next();
            if (node.hasClass('drop')) {
                node.remove();
            }
        },
        createHtml: function (elem) {
            //创建节点
            var data = this.options.renderData,
                html = '',
                target = this.options.targetElement.substring(1, this.options.targetElement.length);


            if (Array.isArray(data) && data.length > 0) {
                var temp = '';
                for (var i = 0; i < data.length; i++) {
                    if (i == 0) {
                        temp += '<li class="active" data-target="' + $(elem).attr('data-autoid') + '" data-min="' + data[i]['min'] + '">' + data[i]['min'] + '[' + data[i]['name'] + '] </li>'
                    } else {
                        temp += '<li data-target="' + $(elem).attr('data-autoid') + '" data-min="' + data[i]['min'] + '">' + data[i]['min'] + '[' + data[i]['name'] + '] </li>'
                    }
                }

                html = '<div class="drop">' +
                    '<ul class="js-min-item-' + target + '">' +
                    temp +
                    '</ul>' +
                    '</div>';
            } else if (Array.isArray(data) && !(data.length > 0)) {
                html = '';
            } else {
                throw new Error('data type must be specified as Array');
            }

            this.template = html;
        }
    }

    var Page = {
        Init: function () {
            Page.Connection = new WebSocket("ws://localhost:12345/ws");

            this.QueryInit();
            this.EXDrawInit();
            this.NDrawInit();
            this.PoolDrawInit();
            this.AddPoolMoneyInit();
            this.AddEggMoneyInit();
            Page.QueryAndShowStatics();
            Page.NGetAndSetNextAward();
            Page.ExDrawSelectDataInit();
            $("#btn-data-init").bind("click", function () {
                Page.DataReset();
            });
            $("#btn-view-init").bind("click", function () {
                Page.ViewReset();
            });
            $("#btn-view-on").bind("click", function () {
                Page.ViewOn();
            });
            $("#btn-view-off").bind("click", function () {
                Page.ViewOff();
            });
            $("#btn-view-empty").bind("click", function () {
                Page.ViewEmpty();
            });
            $("#btn-draw-stop").bind("click", function () {
                Page.DrawStop();
            });
            $("#btn-data-count").unbind().bind("click", function () {
                Page.QueryAndShowStatics();
            });

            var auto = new MinAutoComplete({
                originData: leaderArray
            });
            auto.init();
        },

        InitObj: {
            Type: "init",
            PoolMoney: 0
        },
        ViewOjb: {
            Type: "",
            PoolMoney: 0
        },
        ReadyObj: {
            Type: "ready",
            AllPeopleCount: 0,
            AwardID: 0,
            Drawer: "",
            BackMoney: 0
        },
        StartObj: {
            Type: "start",
            Data: null,
        },
        StopObj: {
            Type: "stop",
            PoolMoney: 0,
        },
        EditPoolObj: {
            Type: "pool",
            Value: 0,
            PoolMoney: 0,
        },

        DataReset: function () {
            $.ajax({
                url: "/api/init",
                data: {pswd: $("#input-data-init").val(),},
                dataType: "json",
                type: "get",
                success: function (info) {
                    if (info.code == 0) {
                        Page.ShowStaticsMsg("数据初始化成功..", null);
                    } else {
                        Page.ShowStaticsMsg("数据初始化失败", info.msg);
                    }
                },
                error: function (e) {
                    Page.ShowStaticsMsg(null, "数据初始化异常");
                }
            });
        },

        ViewReset: function () {
            Page.InitObj.PoolMoney = Page.GetPoolMoney();
            Page.Connection.send(JSON.stringify(Page.InitObj));
        },
        ViewOn: function () {
            Page.ViewOjb.PoolMoney = Page.GetPoolMoney();
            Page.ViewOjb.Type = "view_on";
            Page.Connection.send(JSON.stringify(Page.ViewOjb));
        },
        ViewOff: function () {
            Page.ViewOjb.PoolMoney = Page.GetPoolMoney();
            Page.ViewOjb.Type = "view_off";
            Page.Connection.send(JSON.stringify(Page.ViewOjb));
        },
        ViewEmpty: function () {
            Page.ViewOjb.PoolMoney = Page.GetPoolMoney();
            Page.ViewOjb.Type = "view_empty";
            Page.Connection.send(JSON.stringify(Page.ViewOjb));
        },
        DrawStop: function () {
            Page.ViewOjb.PoolMoney = Page.GetPoolMoney();
            Page.ViewOjb.Type = Page.GetPoolMoney();
            Page.Connection.send(JSON.stringify(Page.StopObj));
        },
        NDrawInit: function () {
            $("#btn-ndraw-ready").unbind().bind("click", function () {
                Page.NDrawReady();
            });
            $("#btn-ndraw-start").unbind().bind("click", function () {
                Page.NDrawStart();
            });
            $("#btn-ndraw-stop").unbind().bind("click", function () {
                Page.NDrawStop();
            });

            $("#btn-ndraw-ready").removeClass("draw-btn-disable");
            $("#btn-ndraw-start").addClass("draw-btn-disable");
            $("#btn-ndraw-stop").addClass("draw-btn-disable");
        },
        NSetDefaultValue: function () {
            if (parseInt($("#ndraw-leader-count").val()) > 0) {
            } else {
                $("#ndraw-leader-count").val(0);
            }
        },
        NGetDrawParam: function () {
            var leaderCount_value = $("#ndraw-leader-count").val();
            var staffCount_value = $("#ndraw-staff-count").val();
            var memo_value = $("#ndraw-memo").val();
            var action_id = JSON.parse($("#input-next-award").val()).data.ID;
            var param = {
                LeaderCount: leaderCount_value ? leaderCount_value : "0",
                StaffCount: staffCount_value ? staffCount_value : "0",
                ActionID: action_id,
                Memo: memo_value
            };
            return param;
        },
        NValidateDrawParam: function (thisParam) {
            var result = true;
            var validateInfo = "";
            if (thisParam) {
                if (parseInt(thisParam.LeaderCount) > -1) {
                } else {
                    validateInfo += " 领导小于0";
                    result = false;
                }
            } else {
                validateInfo += " 前端异常";
                result = false;
            }
            if (!result) {
                Page.EXShowDrawMsg("抽奖参数验证失败：" + validateInfo, null);
            }
            return result;
        },
        NDrawReady: function () {
            Page.NSetDefaultValue();
            var nextAwardStr = $("#input-next-award").val();
            try {
                var nextAward = JSON.parse(nextAwardStr);
                if (nextAward.code == 0) {
                    Page.ReadyObj.PeopleCount = nextAward.data.PeopleCount;
                    Page.ReadyObj.AwardID = nextAward.data.AwardID;
                    Page.ReadyObj.BackMoney = nextAward.data.BackMoney;
                    Page.Connection.send(JSON.stringify(Page.ReadyObj));
                    Page.NShowDrawMsg("[" + nextAward.data.AwardID + "]准备完成...", null);
                    $("#btn-ndraw-ready").removeClass("draw-btn-disable");
                    $("#btn-ndraw-start").removeClass("draw-btn-disable");
                    $("#btn-ndraw-stop").addClass("draw-btn-disable");
                } else if (nextAward.code == 1) {
                    Page.NShowDrawMsg("奖品已经抽完..", null);
                    $("#btn-ndraw-ready").addClass("draw-btn-disable");
                    $("#btn-ndraw-start").addClass("draw-btn-disable");
                    $("#btn-ndraw-stop").addClass("draw-btn-disable");
                } else {
                    Page.NShowDrawMsg(null, "系统异常1");
                }
            } catch (ex) {
                Page.NShowDrawMsg(null, "系统异常");
            }
        },
        NDrawStart: function () {
            var drawParam = Page.NGetDrawParam();
            if (!Page.NValidateDrawParam(drawParam)) {
                return;
            } else {
                Page.NShowDrawMsg("抽奖参数校验通过", null);
            }
            $.ajax({
                url: "/api/ndraw",
                data: {
                    leaderCount: drawParam.LeaderCount,
                    staffCount: drawParam.StaffCount,
                    memo: drawParam.Memo,
                    actionID: drawParam.ActionID
                },
                beforeSend: function () {
                    $("#btn-ndraw-ready").addClass("draw-btn-disable");
                    $("#btn-ndraw-start").addClass("draw-btn-disable");
                    $("#btn-ndraw-stop").addClass("draw-btn-disable");
                },
                dataType: "json",
                type: "get",
                success: function (info) {
                    if (info.code == "0") {//成功
                        Page.NShowDrawMsg("抽奖成功...", null);
                        Page.StartObj.Data = info.data;
                        Page.Connection.send(JSON.stringify(Page.StartObj));
                        Page.NShowDrawMsg("抽奖成功，发送抽奖标识成功...", null);
                        $("#btn-ndraw-ready").addClass("draw-btn-disable");
                        $("#btn-ndraw-start").addClass("draw-btn-disable");
                        $("#btn-ndraw-stop").removeClass("draw-btn-disable");
                        Page.NShowDrawMsg("抽奖成功，等待停止...", null);
                    } else {
                        Page.NDrawInit();
                        Page.NShowDrawMsg("抽奖后台异常,确认异常信息是否重新抽奖", info.msg);
                    }
                },
                error: function (e) {
                    console.log(e);
                }
            });
        },
        NDrawStop: function () {
            $("#btn-ndraw-ready").removeClass("draw-btn-disable");
            $("#btn-ndraw-start").addClass("draw-btn-disable");
            $("#btn-ndraw-stop").addClass("draw-btn-disable");
            Page.StopObj.PoolMoney = Page.GetPoolMoney();
            Page.Connection.send(JSON.stringify(Page.StopObj));
            Page.NShowDrawMsg("停止抽奖，发送停止标识成功....", null);
            Page.NDrawInit();
            Page.NShowDrawMsg("停止抽奖成功，等待抽奖....", null);
            Page.NGetAndSetNextAward();
            Page.QueryAndShowStatics();
        },
        NShowDrawMsg: function (warnMsg, errorMsg) {
            if (warnMsg) {
                $("#ndraw-warn-msg").text(warnMsg);
            } else {
                $("#ndraw-warn-msg").text("");
            }

            if (errorMsg) {
                $("#ndraw-error-msg").text(errorMsg);
            } else {
                $("#ndraw-error-msg").text("");
            }
        },
        NGetAndSetNextAward: function () {
            $.ajax({
                url: "/api/getNextAction",
                dataType: "json",
                type: "get",
                success: function (info) {
                    $("#input-next-award").val(JSON.stringify(info));
                    if (info.code == 0) {
                        $("#now-award-info").text("当前步骤:" + info.data.ID + "---当前奖品标识:" + info.data.AwardID + "---中奖总人数:" + info.data.PeopleCount + "---返奖金额:" + info.data.BackMoney);
                    } else if (info.code == 1) {
                        $("#now-award-info").text("已抽取完成");
                        $("#btn-ndraw-ready").addClass("draw-btn-disable");
                        $("#btn-ndraw-start").addClass("draw-btn-disable");
                        $("#btn-ndraw-stop").addClass("draw-btn-disable");
                    } else {
                        $("#now-award-info").text("异常");
                    }
                },
                error: function (e) {
                    console.log(e);
                }
            });
        },

        EXDrawInit: function () {
            $("#btn-exdraw-ready").unbind().bind("click", function () {
                Page.EXDrawReady();
            });
            $("#btn-exdraw-start").unbind().bind("click", function () {
                Page.EXDrawStart();
            });
            $("#btn-exdraw-stop").unbind().bind("click", function () {
                Page.EXDrawStop();
            });

            $("#btn-exdraw-ready").removeClass("draw-btn-disable");
            $("#btn-exdraw-start").addClass("draw-btn-disable");
            $("#btn-exdraw-stop").addClass("draw-btn-disable");
        },
        EXSetDefaultValue: function () {
            if (parseInt($("#exdraw-type-value").val()) > 0) {
            } else {
                $("#exdraw-type-value").val(1);
            }
            if (parseInt($("#exdraw-leader-count").val()) > 0) {
            } else {
                $("#exdraw-leader-count").val(0);
            }
            if (parseInt($("#exdraw-mix-count").val()) > 0) {
            } else {
                $("#exdraw-mix-count").val(0);
            }
            if (parseInt($("#exdraw-backmoney").val()) > 0) {
            } else {
                $("#exdraw-backmoney").val(0);
            }
        },
        EXGetDrawParam: function () {
            return {
                Drawer: $("#exdraw-drawer").val(),
                AwardID: $("#exdraw-type").val(),
                AwardName: $("#exdraw-type").find("option:selected").attr("data-name"),
                AwardPicName: $("#exdraw-type").find("option:selected").attr("data-picname"),
                AwardCount: $("#exdraw-type-value").val(),
                BackMoney: $("#exdraw-backmoney").val(),
                LeaderCount: $("#exdraw-leader-count").val(),
                StaffCount: $("#exdraw-staff-count").val(),
                MixCount: $("#exdraw-mix-count").val(),
                Memo: $("#exdraw-memo").val()
            };
        },
        EXDrawReady: function () {
            Page.EXSetDefaultValue();
            var thisParam = Page.EXGetDrawParam();
            Page.ReadyObj.AwardID = thisParam.AwardID;
            Page.ReadyObj.AwardName = thisParam.AwardName;
            Page.ReadyObj.AwardPicName = thisParam.AwardPicName;
            Page.ReadyObj.AwardCount = thisParam.AwardCount;
            Page.ReadyObj.BackMoney = thisParam.BackMoney;
            Page.ReadyObj.LeaderCount = thisParam.LeaderCount;
            Page.ReadyObj.StaffCount = thisParam.StaffCount;
            Page.ReadyObj.MixCount = thisParam.MixCount;
            Page.ReadyObj.Memo = thisParam.Memo;
            Page.ReadyObj.PoolMoney = Page.GetPoolMoney();
            Page.ReadyObj.Drawer = thisParam.Drawer;
            var leaderCount = 0;
            if (parseInt(thisParam.LeaderCount) > 0) {
                leaderCount = parseInt(thisParam.LeaderCount);
            }
            var staffCount = 0;
            if (parseInt(thisParam.StaffCount) > 0) {
                staffCount = parseInt(thisParam.StaffCount);
            }
            var mixCount = 0;
            if (parseInt(thisParam.MixCount) > 0) {
                mixCount = parseInt(thisParam.MixCount);
            }
            Page.ReadyObj.AllPeopleCount = leaderCount + mixCount + staffCount;
            Page.Connection.send(JSON.stringify(Page.ReadyObj));

            $("#btn-exdraw-ready").removeClass("draw-btn-disable");
            $("#btn-exdraw-start").removeClass("draw-btn-disable");
            $("#btn-exdraw-stop").addClass("draw-btn-disable");

            Page.EXShowDrawMsg("参数验证通过", null);
        },
        EXDrawStart: function () {
            var drawParam = Page.EXGetDrawParam();
            $.ajax({
                url: "/api/exdraw",
                data: {
                    drawer: drawParam.Drawer,
                    memo: drawParam.Memo,
                    awardID: drawParam.AwardID,
                    awardCount: drawParam.AwardCount,
                    backMoney: drawParam.BackMoney,
                    mixPeopleCount: drawParam.MixCount,
                    leaderCount: drawParam.LeaderCount,
                    staffCount: drawParam.StaffCount,
                },
                dataType: "json",
                type: "get",
                beforeSend: function () {
                    $("#btn-exdraw-ready").addClass("draw-btn-disable");
                    $("#btn-exdraw-start").addClass("draw-btn-disable");
                    $("#btn-exdraw-stop").addClass("draw-btn-disable");
                    Page.EXShowDrawMsg("抽奖中....", null);
                },
                success: function (info) {
                    if (info.code == "0") {//成功
                        Page.EXShowDrawMsg("抽奖成功...", null);
                        Page.StartObj.Data = info.data;
                        Page.Connection.send(JSON.stringify(Page.StartObj));
                        Page.EXShowDrawMsg("抽奖成功，发送抽奖标识成功...", null);
                        $("#btn-exdraw-ready").addClass("draw-btn-disable");
                        $("#btn-exdraw-start").addClass("draw-btn-disable");
                        $("#btn-exdraw-stop").removeClass("draw-btn-disable");
                        Page.EXShowDrawMsg("抽奖成功，等待停止...", null);
                    } else {
                        Page.EXDrawInit();
                        Page.EXShowDrawMsg("抽奖后台异常,确认异常信息是否重新抽奖", info.msg);
                    }
                },
                error: function (e) {
                    $("#error-msg").text("系统异常");
                }
            });
        },
        EXDrawStop: function () {
            Page.StopObj.PoolMoney = Page.GetPoolMoney();
            Page.Connection.send(JSON.stringify(Page.StopObj));
            Page.EXShowDrawMsg("停止抽奖，发送停止标识成功....", null);
            Page.EXDrawInit();
            Page.EXShowDrawMsg("停止抽奖成功，等待抽奖....", null);
            Page.QueryAndShowStatics();
        },
        EXShowDrawMsg: function (warnMsg, errorMsg) {
            if (warnMsg) {
                $("#warn-msg").text(warnMsg);
            } else {
                $("#warn-msg").text("");
            }

            if (errorMsg) {
                $("#error-msg").text(errorMsg);
            } else {
                $("#error-msg").text("");
            }
        },

        PoolDrawInit: function () {
            $("#btn-pooldraw-ready").unbind().bind("click", function () {
                Page.PoolDrawReady();
            });
            $("#btn-pooldraw-start").unbind().bind("click", function () {
                Page.PoolDrawStart();
            });
            $("#btn-pooldraw-stop").unbind().bind("click", function () {
                Page.PoolDrawStop();
            });

            $("#btn-pooldraw-ready").removeClass("draw-btn-disable");
            $("#btn-pooldraw-start").addClass("draw-btn-disable");
            $("#btn-pooldraw-stop").addClass("draw-btn-disable");
        },
        PoolSetDefaultValue: function () {
            if (parseInt($("#pooldraw-money").val()) > 0) {
            } else {
                $("#pooldraw-money").val(0);
            }
            if (parseInt($("#pooldraw-leader-count").val()) > 0) {
            } else {
                $("#pooldraw-leader-count").val(0);
            }
            if (parseInt($("#pooldraw-staff-count").val()) > 0) {
            } else {
                $("#pooldraw-staff-count").val(0);
            }
            if (parseInt($("#pooldraw-mix-count").val()) > 0) {
            } else {
                $("#pooldraw-mix-count").val(0);
            }
        },
        PoolGetDrawParam: function () {
            return {
                DrawMoney: $("#pooldraw-money").val(),
                BackRatio: $("#pooldraw-leader-ratio").val(),
                LeaderCount: $("#pooldraw-leader-count").val(),
                StaffCount: $("#pooldraw-staff-count").val(),
                MixCount: $("#pooldraw-mix-count").val(),
                Memo: $("#pooldraw-memo").val(),
            };
        },
        PoolDrawReady: function () {
            Page.PoolSetDefaultValue();
            let thisParam = Page.PoolGetDrawParam();
            Page.ReadyObj.AwardID = 0;
            Page.ReadyObj.AwardName = "返奖池现金";
            Page.ReadyObj.DrawMoney = thisParam.DrawMoney;
            Page.ReadyObj.BackRatio = thisParam.BackRatio;
            Page.ReadyObj.LeaderCount = thisParam.LeaderCount;
            Page.ReadyObj.StaffCount = thisParam.StaffCount;
            Page.ReadyObj.MixCount = thisParam.MixCount;
            Page.ReadyObj.PoolMoney = Page.GetPoolMoney();
            Page.ReadyObj.Drawer = "pool";
            var leaderCount = 0;
            if (parseInt(thisParam.LeaderCount) > 0) {
                leaderCount = parseInt(thisParam.LeaderCount);
            }
            var mixCount = 0;
            if (parseInt(thisParam.MixCount) > 0) {
                mixCount = parseInt(thisParam.MixCount);
            }
            var staffCount = 0;
            if (parseInt(thisParam.StaffCount) > 0) {
                staffCount = parseInt(thisParam.StaffCount);
            }
            Page.ReadyObj.AllPeopleCount = leaderCount + mixCount + staffCount;
            Page.Connection.send(JSON.stringify(Page.ReadyObj));

            $("#btn-pooldraw-ready").removeClass("draw-btn-disable");
            $("#btn-pooldraw-start").removeClass("draw-btn-disable");
            $("#btn-pooldraw-stop").addClass("draw-btn-disable");

            Page.PoolShowDrawMsg("参数验证通过", null);
        },
        PoolDrawStart: function () {
            let drawParam = Page.PoolGetDrawParam();
            $.ajax({
                url: "/api/pooldraw",
                data: {
                    drawer:"返奖池",
                    memo: drawParam.Memo,
                    awardCount: drawParam.DrawMoney,
                    backRatio:drawParam.BackRatio,
                    mixPeopleCount: drawParam.MixCount,
                    leaderCount: drawParam.LeaderCount,
                    staffCount: drawParam.StaffCount,
                },
                dataType: "json",
                type: "get",
                beforeSend: function () {
                    $("#btn-pooldraw-ready").addClass("draw-btn-disable");
                    $("#btn-pooldraw-start").addClass("draw-btn-disable");
                    $("#btn-pooldraw-stop").addClass("draw-btn-disable");
                    Page.PoolShowDrawMsg("抽奖中....", null);
                },
                success: function (info) {
                    if (info.code == "0") {//成功
                        Page.PoolShowDrawMsg("抽奖成功...", null);
                        Page.StartObj.Data = info.data;
                        Page.Connection.send(JSON.stringify(Page.StartObj));
                        Page.PoolShowDrawMsg("抽奖成功，发送抽奖标识成功...", null);
                        $("#btn-pooldraw-ready").addClass("draw-btn-disable");
                        $("#btn-pooldraw-start").addClass("draw-btn-disable");
                        $("#btn-pooldraw-stop").removeClass("draw-btn-disable");
                        Page.PoolShowDrawMsg("抽奖成功，等待停止...", null);
                    } else {
                        Page.PoolDrawInit();
                        Page.PoolShowDrawMsg("抽奖后台异常,确认异常信息是否重新抽奖", info.msg);
                    }
                },
                error: function (e) {
                    Page.PoolShowDrawMsg(null, "系统异常");
                }
            });
        },
        PoolDrawStop: function () {
            Page.StopObj.PoolMoney = Page.GetPoolMoney();
            Page.Connection.send(JSON.stringify(Page.StopObj));
            Page.PoolShowDrawMsg("停止抽奖，发送停止标识成功....", null);
            Page.PoolDrawInit();
            Page.PoolShowDrawMsg("停止抽奖成功，等待抽奖....", null);
            Page.QueryAndShowStatics();
        },
        PoolShowDrawMsg: function (warnMsg, errorMsg) {
            if (warnMsg) {
                $("#pooldraw-warn-msg").text(warnMsg);
            } else {
                $("#pooldraw-warn-msg").text("");
            }

            if (errorMsg) {
                $("#pooldraw-error-msg").text(errorMsg);
            } else {
                $("#pooldraw-error-msg").text("");
            }
        },
        QueryAndShowStatics: function () {
            $.ajax({
                url: "/api/count",
                dataType: "json",
                type: "get",
                beforeSend: function () {
                    Page.ShowStaticsMsg("统计数据查询中..", null);
                },
                success: function (info) {
                    if (info.code == "0") {
                        $("#count-all").text(info.count.AllPeopleCount);
                        $("#count-poolmoney").text(info.count.PoolMoney);
                        // $("#count-bigeggmoney").text(info.Count.BigEggMoney);
                        $("#count-lucky-all").text(info.count.AllLuckyCount);
                        $("#count-lucky-leader").text(info.count.LuckyLeaderCount);
                        $("#count-nolucky-leader").text(info.count.NoLuckyLeaderCount);
                        $("#count-lucky-staff").text(info.count.LuckyStaffCount);
                        $("#count-nolucky-staff").text(info.count.NoLuckyStaffCount);
                        Page.ShowStaticsMsg("统计数据查询成功..", null);
                    } else {
                        Page.ShowStaticsMsg("统计数据查询后台异常..", info.msg);
                    }
                },
                error: function (e) {
                    Page.ShowStaticsMsg("统计数据查询异常..", null);
                    console.log(e);
                }
            });
        },
        AddPoolMoneyInit: function () {
            $("#btn-add-pool").unbind().bind("click", function () {
                Page.AddPoolMoney()
            });
        },
        AddPoolMoney: function () {
            var money_value = parseInt($("#input-pool-money").val());
            var memo_value = $("#input-pool-memo").val();
            if (money_value > 0 || money_value < 0) {
            } else {
                Page.ShowStaticsMsg("添加金额大于零或小于零", null);
                return;
            }
            if (memo_value == "") {
                Page.ShowStaticsMsg("备注不能为空", null);
                return;
            }
            $.ajax({
                url: "/api/addMoney",
                data: {
                    money: money_value, memo: memo_value
                },
                dataType: "json",
                type: "get",
                beforeSend: function () {
                    $("#btn-add-pool").unbind();
                    Page.ShowStaticsMsg("奖金池添加中...", null);
                },
                success: function (info) {
                    if (info.code == "0") {//成功
                        Page.ShowStaticsMsg("奖金池添加成功", null);
                        Page.EditPoolObj.Value = money_value;
                        Page.EditPoolObj.PoolMoney = Page.GetPoolMoney();
                        Page.Connection.send(JSON.stringify(Page.EditPoolObj));
                    } else if (info.code == "1") {//已存在
                        Page.ShowStaticsMsg("奖金池添加失败", info.msg);
                    } else {
                        Page.ShowStaticsMsg("奖金池添加后台异常", info.msg);
                    }
                },
                error: function (e) {
                    Page.ShowStaticsMsg("奖金池添加失败", "系统异常");
                },
                complete: function () {
                    Page.QueryAndShowStatics();
                    $("#btn-add-pool").unbind().bind("click", function () {
                        Page.AddPoolMoney()
                    });
                }
            });
        },
        ShowStaticsMsg: function (warnMsg, errorMsg) {
            if (warnMsg) {
                $("#count-warn-msg").text(warnMsg);
            } else {
                $("#count-warn-msg").text("");
            }

            if (errorMsg) {
                $("#count-error-msg").text(errorMsg);
            } else {
                $("#count-error-msg").text("");
            }
        },

        ExDrawSelectDataInit: function () {
            $.ajax({
                url: "/api/getAwards",
                dataType: "json",
                type: "get",
                success: function (info) {
                    $("#exdraw-type").empty()
                    if (info.code == 0) {
                        for (var i = 0; i < info.data.length; i++) {
                            $("#exdraw-type").append("<option value=\"" + info.data[i].ID + "\" data-name=\"" + info.data[i].Name + "\"\ data-picname=\"" + info.data[i].PicName + "\">" + info.data[i].Name + "</option>")
                        }
                    } else {
                        return 0;
                    }
                },
                error: function (e) {
                    console.log(e);
                }
            });
        },
        GetPoolMoney: function () {
            var poolMoney = 0;
            $.ajax({
                url: "/api/count",
                async: false,
                dataType: "json",
                type: "get",
                success: function (info) {
                    if (info.code == 0) {
                        poolMoney = info.count.PoolMoney;
                    } else {
                        return 0;
                    }
                },
                error: function (e) {
                    console.log(e);
                }
            });
            return poolMoney;
        },

        //web socket obj
        Connection: null,
    }
    Page.Init();
});
