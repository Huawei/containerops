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

import {PIPELINE_ACTION} from "..common/constant";
import {svgActionWidth, svgActionHeight, svgButtonWidth, svgButtonHeight} from "..common/constant";



export function stageButtonView(sd,si,buttonView,initAction,initPipeline,pipelineData){
	
	buttonView.selectAll("image").remove();

    //show stage pop button
    buttonView.append("image")
        .attr("xlink:href", function (d, i) {
            return "../../assets/svg/actionAdd.svg";
        })
        .attr("id", function (d, i) {
            return "button" + "-" + uuid.v1();
        })
        .attr("width", function (d, i) {
            return svgButtonWidth;
        })
        .attr("height", function (d, i) {
            return svgButtonHeight;
        })
        .attr("translateX", function (d, i) {
            return sd.translateX - (svgButtonWidth * 2);
        })
        .attr("translateY", function (d, i) {
            return sd.translateY;
        })
        .attr("transform", function (d, i) {
            return "translate(" + this.attributes["translateX"].value + "," + this.attributes["translateY"].value + ")";
        })
        .on("mouseover", function (d, i) {
            d3.select(this)
                .attr("transform",
                    "translate("
                    + (this.attributes["translateX"].value - svgButtonWidth / 2) + ","
                    + (this.attributes["translateY"].value - svgButtonHeight / 2) + ") scale(2)");
        })
        .on("mouseout", function (d, i) {
            d3.select(this)
                .attr("transform",
                    "translate("
                    + this.attributes["translateX"].value + ","
                    + this.attributes["translateY"].value + ") scale(1)");
        })
        .on("click", function (d, i) {
            sd.actions.splice(
                sd.actions.length,
                0,
                {
                    id: PIPELINE_ACTION + "-" + uuid.v1(),
                    type: PIPELINE_ACTION,
                    class: PIPELINE_ACTION,
                    drawX: 0,
                    drawY: 0,
                    width: 0,
                    height: 0,
                    translateX: 0,
                    translateY: 0,
                    setupData: {}
                });
            buttonView.selectAll("image").remove();
            initAction();
        });

    //show del stage button
    buttonView.append("image")
        .attr("xlink:href", function (d, i) {
            return "../../assets/svg/stageDel.svg";
        })
        .attr("id", function (d, i) {
            return "button" + "-" + uuid.v1();
        })
        .attr("width", function (d, i) {
            return svgButtonWidth;
        })
        .attr("height", function (d, i) {
            return svgButtonHeight;
        })
        .attr("translateX", function (d, i) {
            return sd.translateX + (svgButtonWidth / 3);
        })
        .attr("translateY", function (d, i) {
            return sd.translateY - (svgButtonHeight * 2);
        })
        .attr("transform", function (d, i) {
            return "translate("
                + this.attributes["translateX"].value + ","
                + this.attributes["translateY"].value + ")";
        })
        .on("mouseover", function (d, i) {
            d3.select(this)
                .attr("transform",
                    "translate("
                    + (this.attributes["translateX"].value - svgButtonWidth / 2) + ","
                    + (this.attributes["translateY"].value - svgButtonHeight / 2) + ") scale(2)");
        })
        .on("mouseout", function (d, i) {
            d3.select(this)
                .attr("transform",
                    "translate("
                    + this.attributes["translateX"].value + ","
                    + this.attributes["translateY"].value + ") scale(1)");
        })
        .on("click", function (d, i) {
            buttonView.selectAll("image").remove();
            pipelineData.splice(si, 1);
            initPipeline();
            initAction();
        });


    //show close stage pop button
    buttonView.append("image")
        .attr("xlink:href", function (d, i) {
            return "../../assets/svg/stageClosePop.svg";
        })
        .attr("id", function (d, i) {
            return "button" + "-" + uuid.v1();
        })
        .attr("width", function (d, i) {
            return svgButtonWidth;
        })
        .attr("height", function (d, i) {
            return svgButtonHeight;
        })
        .attr("translateX", function (d, i) {
            return sd.translateX + (svgButtonWidth * 2.6);
        })
        .attr("translateY", function (d, i) {
            return sd.translateY;
        })
        .attr("transform", function (d, i) {
            return "translate("
                + this.attributes["translateX"].value + ","
                + this.attributes["translateY"].value + ")";
        })
        .on("mouseover", function (d, i) {
            d3.select(this)
                .attr("transform",
                    "translate("
                    + (this.attributes["translateX"].value - svgButtonWidth / 2) + ","
                    + (this.attributes["translateY"].value - svgButtonHeight / 2) + ") scale(2)");
        })
        .on("mouseout", function (d, i) {
            d3.select(this)
                .attr("transform",
                    "translate("
                    + this.attributes["translateX"].value + ","
                    + this.attributes["translateY"].value + ") scale(1)");
        })
        .on("click", function (d, i) {
            buttonView.selectAll("image").remove();
        });

}