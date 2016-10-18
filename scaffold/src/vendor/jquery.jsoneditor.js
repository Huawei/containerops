// Simple yet flexible JSON editor plugin.
// Turns any element into a stylable interactive JSON editor.

// Copyright (c) 2013 David Durman

// Licensed under the MIT license (http://www.opensource.org/licenses/mit-license.php).

// Dependencies:

// * jQuery
// * JSON (use json2 library for browsers that do not support JSON natively)

// Example:

//     var myjson = { any: { json: { value: 1 } } };
//     var opt = { change: function() { /* called on every change */ } };
//     /* opt.propertyElement = '<textarea>'; */ // element of the property field, <input> is default
//     /* opt.valueElement = '<textarea>'; */  // element of the value field, <input> is default
//     $('#mydiv').jsonEditor(myjson, opt);
import {ContextMenu} from "./jquery.jsoneditor.menu";

var items = [];

items.push({
      text: 'Type',
      title: 'Change the type of this field',
      className: 'jsoneditor-type-',
      submenu: [
        {
          text: 'Array',
          className: 'jsoneditor-type-array',
          title: "titles.array",
          click: function (button,opt) {
            changeType('array',button,opt);
          }
        },
        {
          text: 'Object',
          className: 'jsoneditor-type-object',
          title: "titles.object",
          click: function (button,opt) {
            changeType('object',button,opt);
          }
        },
        {
          text: 'String',
          className: 'jsoneditor-type-string',
          title: "titles.string",
          click: function (button,opt) {
            changeType('string',button,opt);
          }
        }
    ]
});



items.push({
      text: 'ValueType',
      title: 'Change the type of value',
      className: 'jsoneditor-type-',
      submenu: [
        {
          text: 'Changeable',
          className: 'jsoneditor-value-change',
          title: "titles.changeable",
          click: function (button,opt) {
            changeValueType("changeable",button,opt);
          }
        },
        {
          text: 'Unchangeable',
          className: 'jsoneditor-value-unchange',
          title: "titles.unchangeable",
          click: function (button,opt) {
            changeValueType("unchangeable",button,opt);
          }
        }
    ]
});



items.push({
    text: 'Remove',
    title: 'Remove this field (Ctrl+Del)',
    className: 'jsoneditor-remove',
    click: function (button,opt) {
      removeItem(button,opt);
    }
});







export function jsonEditor (container,json, options) {
    options = options || {};
    // Make sure functions or other non-JSON data types are stripped down.
    json = parse(stringify(json));
    
    var K = function() {};
    var onchange = options.change || K;
    var onpropertyclick = options.propertyclick || K;

    return container.each(function() {
        JSONEditorInit(container, json, onchange, onpropertyclick, options.propertyElement, options.valueElement);
    });
    
};

function JSONEditorInit(target, json, onchange, onpropertyclick, propertyElement, valueElement) {
    var opt = {
        target: target,
        onchange: onchange,
        onpropertyclick: onpropertyclick,
        original: json,
        propertyElement: propertyElement,
        valueElement: valueElement
    };
    construct(opt, json, opt.target);
    $(opt.target).on('blur focus', '.property, .value', function() {
        $(this).toggleClass('editing');
    });
}

function isObject(o) { return Object.prototype.toString.call(o) == '[object Object]'; }
function isArray(o) { return Object.prototype.toString.call(o) == '[object Array]'; }
function isBoolean(o) { return Object.prototype.toString.call(o) == '[object Boolean]'; }
function isNumber(o) { return Object.prototype.toString.call(o) == '[object Number]'; }
function isString(o) { return Object.prototype.toString.call(o) == '[object String]'; }
var types = 'object array boolean number string null';

// Feeds object `o` with `value` at `path`. If value argument is omitted,
// object at `path` will be deleted from `o`.
// Example:
//      feed({}, 'foo.bar.baz', 10);    // returns { foo: { bar: { baz: 10 } } }
function feed(o, path, value) {
    var del = arguments.length == 2;
    
    if (path.indexOf('.') > -1) {
        var diver = o,
            i = 0,
            parts = path.split('.');
        for (var len = parts.length; i < len - 1; i++) {
            diver = diver[parts[i]];
        }
        if (del) delete diver[parts[len - 1]];
        else diver[parts[len - 1]] = value;
    } else {
        if (del) delete o[path];
        else o[path] = value;
    }
    return o;
}

// Get a property by path from object o if it exists. If not, return defaultValue.
// Example:
//     def({ foo: { bar: 5 } }, 'foo.bar', 100);   // returns 5
//     def({ foo: { bar: 5 } }, 'foo.baz', 100);   // returns 100
function def(o, path, defaultValue) {
    path = path.split('.');
    var i = 0;
    while (i < path.length) {
        if ((o = o[path[i++]]) == undefined) return defaultValue;
    }
    return o;
}

function error(reason) { if (window.console) { console.error(reason); } }

function parse(str) {
    var res;
    try { res = JSON.parse(str); }
    catch (e) { res = null; error('JSON parse failed.'); }
    return res;
}

function stringify(obj) {
    var res;
    try { res = JSON.stringify(obj); }
    catch (e) { res = 'null'; error('JSON stringify failed.'); }
    return res;
}

function addMenu(item,opt){
    if (item.children('.menuer').length == 0) {
        var menuer =   $('<span>',  { 'class': 'menuer fa fa-navicon' });
        menuer.bind('click', function() {
            window.event.stopPropagation();
            ContextMenu($(this),items,opt);
        });
        item.prepend(menuer);
    }
}

function removeItem (button,opt){
    // var menuButton = button.parents(".item").find(".menuer");
    var item = button.parents(".item:eq(0)");
    
    item.find(">.property").val("").change();
    item.remove();
    
}

function addExpander(item) {
    if (item.children('.expander').length == 0) {
        var expander =   $('<span>',  { 'class': 'expander' });
        expander.bind('click', function() {
            var item = $(this).parent();
            item.toggleClass('expanded');
        });
        item.prepend(expander);
    }
}

function addListAppender(item, handler) {
    var appender = $('<div>', { 'class': 'item appender' }),
        btn      = $('<button></button>', { 'class': 'property' });

    btn.text('Add New Value');

    appender.append(btn);
    item.append(appender);

    btn.click(handler);

    return appender;
}

function addNewValue(json) {
    if (isArray(json)) {
        json.push(null);
        return true;
    }

    if (isObject(json)) {
        var i = 1, newName = "newKey";

        while (json.hasOwnProperty(newName)) {
            newName = "newKey" + i;
            i++;
        }

        json[newName] = null;
        return true;
    }

    return false;
}

function construct(opt, json, root, path) {
    path = path || '';
    
    root.children('.item').remove();
    
    for (var key in json) {
        if (!json.hasOwnProperty(key)) continue;

        var item     = $('<div>',   { 'class': 'item', 'data-path': path }),
            property =   $(opt.propertyElement || '<input>', { 'class': 'property' }),
            value    =   $(opt.valueElement || '<input>', { 'class': 'value'    });

        if (isObject(json[key]) || isArray(json[key])) {
            addExpander(item);
        }
        

        addMenu(item,opt);


        item.append(property).append(value);
        root.append(item);
        
        property.val(key).attr('title', key);
        var val = stringify(json[key]);
        value.val(val).attr('title', val);

        assignType(item, json[key]);

        property.change(propertyChanged(opt));
        value.change(valueChanged(opt));
        property.click(propertyClicked(opt));
        
        if (isObject(json[key]) || isArray(json[key])) {
            construct(opt, json[key], item, (path ? path + '.' : '') + key);
        }
    }

    if (isObject(json) || isArray(json)) {
        addListAppender(root, function () {
            addNewValue(json);
            construct(opt, json, root, path);
            opt.onchange(parse(stringify(opt.original)));
        })
    }
}

function updateParents(el, opt) {
    $(el).parentsUntil(opt.target).each(function() {
        var path = $(this).data('path');
        path = (path ? path + '.' : path) + $(this).children('.property').val();
        var val = stringify(def(opt.original, path, null));
        $(this).children('.value').val(val).attr('title', val);
    });
}

function propertyClicked(opt) {
    return function() {
        var path = $(this).parent().data('path');            
        var key = $(this).attr('title');

        var safePath = path ? path.split('.').concat([key]).join('\'][\'') : key;
        
        opt.onpropertyclick('[\'' + safePath + '\']');
    };
}

function propertyChanged(opt) {
    return function() {
        var path = $(this).parent().data('path'),
            val = parse($(this).next().val()),
            newKey = $(this).val(),
            oldKey = $(this).attr('title');

        $(this).attr('title', newKey);

        feed(opt.original, (path ? path + '.' : '') + oldKey);
        if (newKey) feed(opt.original, (path ? path + '.' : '') + newKey, val);

        updateParents(this, opt);

        if (!newKey) $(this).parent().remove();
        
        opt.onchange(parse(stringify(opt.original)));
    };
}

function valueChanged(opt) {
    return function() {
        var key = $(this).prev().val(),
            val = parse($(this).val() || 'null'),
            item = $(this).parent(),
            path = item.data('path');

        feed(opt.original, (path ? path + '.' : '') + key, val);
        if ((isObject(val) || isArray(val)) && !$.isEmptyObject(val)) {
            construct(opt, val, item, (path ? path + '.' : '') + key);
            addExpander(item);
        } else {
            item.find('.expander, .item').remove();
        }

        assignType(item, val);

        updateParents(this, opt);
        
        opt.onchange(parse(stringify(opt.original)));
    };
}


function changeType(type,button,opt){
    var item = button.parents(".item:eq(0)");
    var value = item.find("input.value");
    

    if(item.hasClass(type.toLowerCase())){
        return false;
    }
    if(type == "string")value.val('""');
    else if(type == "array") value.val("[]");
    else if(type == "object") value.val("{}");

    value.change();
}

function changeValueType(type,button,opt){
    var item = button.parents(".item:eq(0)");
    var input = item.find(">input.value");

    if(type == "changeable"){
        input.addClass("show");
    }else{
        input.removeClass("show");
    }
    
}

function assignType(item, val) {
    var className = 'null';
    
    if (isObject(val)) className = 'object';
    else if (isArray(val)) className = 'array';
    else if (isBoolean(val)) className = 'boolean';
    else if (isString(val)) className = 'string';
    else if (isNumber(val)) className = 'number';

    item.removeClass(types);
    item.addClass(className);
}


