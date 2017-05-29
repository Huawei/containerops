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

import {initWorkflowPage,hideWorkflowEnv} from "./workflow/main";
import {initComponentPage} from "./component/main";
import {initHistoryPage} from "./history/main";
import {initSystemSettings,initSystemSettingPage} from "./setting/main";
import {initApi} from "./common/api";

initApi("demo","demo");

$._messengerDefaults = {
    extraClasses: 'messenger-fixed messenger-theme-future messenger-on-bottom messenger-on-right'
}

initSystemSettings(initWorkflowPage);

$(".menu-workflow").on('click',function(event){
    initWorkflowPage();
    $(event.currentTarget).parent().parent().children().removeClass("active");
    $(event.currentTarget).parent().addClass("active");
})

$(".menu-component").on('click',function(event){
    initComponentPage();
    $(event.currentTarget).parent().parent().children().removeClass("active");
    $(event.currentTarget).parent().addClass("active");
})

$(".menu-history").on('click',function(event){
    initHistoryPage();
    $(event.currentTarget).parent().parent().children().removeClass("active");
    $(event.currentTarget).parent().addClass("active");
})

$(".menu-setting").on('click',function(event){
    initSystemSettingPage();
    $(event.currentTarget).parent().parent().children().removeClass("active");
    $(event.currentTarget).parent().addClass("active");
})
// initActionLinkView();

// sidebar nav control
$(".nav-control").on("click",function(event){
    var target = $(event.currentTarget);
    if(target.hasClass("sidebar-close")){
        target.removeClass("sidebar-close").addClass("sidebar-open");
        target.removeClass("fa-chevron-circle-left").addClass("fa-chevron-circle-right");
        $("body").removeClass("nav-static").addClass("nav-collapsed");
    }else if(target.hasClass("sidebar-open")){
        target.removeClass("sidebar-open").addClass("sidebar-close");
        target.removeClass("fa-chevron-circle-right").addClass("fa-chevron-circle-left");
        $("body").removeClass("nav-collapsed").addClass("nav-static");
    }
})

$(".workflow-close-env").on('click', function() {
    hideWorkflowEnv();
});

