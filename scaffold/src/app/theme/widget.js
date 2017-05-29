/*
Copyright 2016 - 2017 Huawei Technologies Co., Ltd. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

export function resizeWidget(){
	var $widgets = $('.widget');

    $widgets.on("fullscreen.widgster", function(){
    	$('.content-wrap').css({
         	'-webkit-transform': 'none',
            '-ms-transform': 'none',
            'transform': 'none',
            'margin': 0,
            'z-index': 2
        });
        // $(".treeview").css("max-height",window.screen.height * 2 / 3);
        // $(".treeview").css("overflow","auto");
        $(".widget").css("overflow","auto");
    }).on("restore.widgster closed.widgster", function(){
        $('.content-wrap').css({
        	'-webkit-transform': '',
            '-ms-transform': '',
            'transform': '',
            'margin': '',
            'z-index': ''
        });
        // $(".treeview").css("max-height","");
        // $(".treeview").css("overflow","");
        // $(".widget").css("overflow","hidden");
    });

    $widgets.widgster();
}