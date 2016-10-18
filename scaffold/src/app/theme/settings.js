
$(function(){
    /**
     * A global object containing theme specific colors, screen variables & color functions.
     * @type Object
     */
    window.Scaffold = {
        colors: {
            'white': '#fff',
            'black': '#000',
            'gray-light': '#999',
            'gray-lighter': '#eee',
            'gray': '#666',
            'gray-dark': '#343434',
            'gray-darker': '#222',
            'gray-semi-light': '#777',
            'gray-semi-lighter': '#ddd',
            'brand-primary': '#5d8fc2',
            'brand-success': '#64bd63',
            'brand-warning': '#f0b518',
            'brand-danger': '#dd5826',
            'brand-info': '#5dc4bf'
        },

        screens: {
            'xs-max': 543,
            'sm-min': 544,
            'sm-max': 767,
            'md-min': 768,
            'md-max': 991,
            'lg-min': 992,
            'lg-max': 1199,
            'xl-min': 1200
        },

        isScreen: function(size){
            var screenPx = window.innerWidth;
            return (screenPx >= this.screens[size + '-min'] || size == 'xs') && (screenPx <= this.screens[size + '-max'] || size == 'xl');
        },

        getScreenSize: function(){
            var screenPx = window.innerWidth;
            if (screenPx <= this.screens['xs-max']) return 'xs';
            if ((screenPx >= this.screens['sm-min']) && (screenPx <= this.screens['sm-max'])) return 'sm';
            if ((screenPx >= this.screens['md-min']) && (screenPx <= this.screens['md-max'])) return 'md';
            if ((screenPx >= this.screens['lg-min']) && (screenPx <= this.screens['lg-max'])) return 'lg';
            if (screenPx >= this.screens['xl-min']) return 'xl';
        },

        //credit http://stackoverflow.com/questions/1507931/generate-lighter-darker-color-in-css-using-javascript
        changeColor: function(color, ratio, darker) {
            var pad = function(num, totalChars) {
                var pad = '0';
                num = num + '';
                while (num.length < totalChars) {
                    num = pad + num;
                }
                return num;
            };
            // Trim trailing/leading whitespace
            color = color.replace(/^\s*|\s*$/, '');

            // Expand three-digit hex
            color = color.replace(
                /^#?([a-f0-9])([a-f0-9])([a-f0-9])$/i,
                '#$1$1$2$2$3$3'
            );

            // Calculate ratio
            var difference = Math.round(ratio * 256) * (darker ? -1 : 1),
            // Determine if input is RGB(A)
                rgb = color.match(new RegExp('^rgba?\\(\\s*' +
                    '(\\d|[1-9]\\d|1\\d{2}|2[0-4][0-9]|25[0-5])' +
                    '\\s*,\\s*' +
                    '(\\d|[1-9]\\d|1\\d{2}|2[0-4][0-9]|25[0-5])' +
                    '\\s*,\\s*' +
                    '(\\d|[1-9]\\d|1\\d{2}|2[0-4][0-9]|25[0-5])' +
                    '(?:\\s*,\\s*' +
                    '(0|1|0?\\.\\d+))?' +
                    '\\s*\\)$'
                    , 'i')),
                alpha = !!rgb && rgb[4] != null ? rgb[4] : null,

            // Convert hex to decimal
                decimal = !!rgb? [rgb[1], rgb[2], rgb[3]] : color.replace(
                    /^#?([a-f0-9][a-f0-9])([a-f0-9][a-f0-9])([a-f0-9][a-f0-9])/i,
                    function() {
                        return parseInt(arguments[1], 16) + ',' +
                            parseInt(arguments[2], 16) + ',' +
                            parseInt(arguments[3], 16);
                    }
                ).split(/,/),
                returnValue;

            // Return RGB(A)
            return !!rgb ?
                'rgb' + (alpha !== null ? 'a' : '') + '(' +
                    Math[darker ? 'max' : 'min'](
                        parseInt(decimal[0], 10) + difference, darker ? 0 : 255
                    ) + ', ' +
                    Math[darker ? 'max' : 'min'](
                        parseInt(decimal[1], 10) + difference, darker ? 0 : 255
                    ) + ', ' +
                    Math[darker ? 'max' : 'min'](
                        parseInt(decimal[2], 10) + difference, darker ? 0 : 255
                    ) +
                    (alpha !== null ? ', ' + alpha : '') +
                    ')' :
                // Return hex
                [
                    '#',
                    pad(Math[darker ? 'max' : 'min'](
                        parseInt(decimal[0], 10) + difference, darker ? 0 : 255
                    ).toString(16), 2),
                    pad(Math[darker ? 'max' : 'min'](
                        parseInt(decimal[1], 10) + difference, darker ? 0 : 255
                    ).toString(16), 2),
                    pad(Math[darker ? 'max' : 'min'](
                        parseInt(decimal[2], 10) + difference, darker ? 0 : 255
                    ).toString(16), 2)
                ].join('');
        },
        lighten: function(color, ratio) {
            return this.changeColor(color, ratio, false);
        },
        darken: function(color, ratio) {
            return this.changeColor(color, ratio, true);
        }
    };

    var ScaffoldSettingsBundle = function(){
        var defaultSettings =  {
            /**
             * whether navigation is static (does not collapse automatically)
             */
            'nav-static': false
        };
        this.settingName = 'scaffold-app-settings';
        this._settings = JSON.parse(localStorage.getItem(this.settingName)) || defaultSettings;
    };

    ScaffoldSettingsBundle.prototype.save = function(){
        localStorage.setItem(this.settingName, JSON.stringify(this._settings));
        return this;
    };

    ScaffoldSettingsBundle.prototype.get = function(key){
        return this._settings[key];
    };

    ScaffoldSettingsBundle.prototype.set = function(key, value){
        this._settings[key] = value;
        return this;
    };

    window.ScaffoldSettings = new ScaffoldSettingsBundle();
});