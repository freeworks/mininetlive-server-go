if (!_){
   var _ = {};
}

(function(){
    var browser = _;

    /**
     * [parseParams 解析参数]
     * @param  {[type]} stringParams [string]
     * @return {[type]}              [object]
     */
    function parseParams(stringParams){
        return stringParams ? JSON.parse('{"' + stringParams.replace(/&/g, '","').replace(/=/g, '":"') + '"}', function(key, value){
            return key === "" ? value : decodeURIComponent(value)
        }) : {} 
    }

    /**
     * [parseUrlParams 解析url参数]
     * @return {[type]} [object]
     */
    function parseUrlParams(){
        var search = location.search.substring(1);
        return parseParams(search);
        // return search ? JSON.parse('{"' + search.replace(/&/g, '","').replace(/=/g, '":"') + '"}', function(key, value){
        //     return key === "" ? value : decodeURIComponent(value)
        // }) : {}
    }

    /**
     * [stringifyUrlParams 格式化url参数]
     * @param  {[type]} params [object]
     * @return {[type]}        [description]
     */
    function stringifyUrlParams(params){
        var string = "";
        for (var i in params){
            if (typeof i !== "function"){
                string = i + "=" + params[i] + "&"
            }
            
        }
        return string.slice(0, -1);
    }


    browser.parseParams = parseParams;
    browser.parseUrlParams = parseUrlParams;
    browser.stringifyUrlParams = stringifyUrlParams;

})();