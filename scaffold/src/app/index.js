/* Copyright 2014 Huawei Technologies Co., Ltd. All rights reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License. */

import {initPipelinePage} from "./pipeline/main";
import {initComponentPage} from "./component/main";
// import {historyRecord} from "./historyRecord";

// let $a = d3.select("#showHistory").on("click",historyRecord);

$._messengerDefaults = {
    extraClasses: 'messenger-fixed messenger-theme-future messenger-on-top'
}

initPipelinePage();

$(".menu-pipeline").on('click',function(event){
    initPipelinePage();
    $(event.currentTarget).parent().parent().children().removeClass("active");
    $(event.currentTarget).parent().addClass("active");
})

$(".menu-component").on('click',function(){
    initComponentPage();
    $(event.currentTarget).parent().parent().children().removeClass("active");
    $(event.currentTarget).parent().addClass("active");
})
// initActionLinkView();

function initActionLinkView() {
    actionLinkView.append("rect")
        .attr("x",10)
        .attr("y",10)
        .attr("rx",10)
        .attr("ry",10)
        .attr("width",120)
        .attr("height",40)
        .attr("stroke","red")
        .attr("fill","red")
    ;
}

