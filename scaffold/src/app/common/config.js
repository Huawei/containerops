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

export let 
    PNG_BASEDIR = "../../assets/images/",
    PNG_SUFFIX = ".png",
    SVG_BASEDIR = "../../assets/svg/",
    SVG_SUFFIX = ".svg",
    SVG_WORKFLOW_SET = "workflow-setting",
    SVG_WORKFLOW_SET_SELECTED = "workflow-setting-selected",
    SVG_START = "start-latest",
    SVG_START_SELECTED = "start-selected-latest",
    SVG_STAGE = "stage-latest",
    SVG_STAGE_SELECTED = "stage-selected-latest",
    SVG_ADD_STAGE = "add-stage-latest",
    SVG_ADD_STAGE_SELECTED = "add-stage-selected-latest",
    SVG_END = "end-latest",
    SVG_ACTION = "action-latest",
    SVG_ACTION_SELECTED = "action-selected-latest",

    SVG_ADD_ACTION = "add-action-latest",
    SVG_ADD_ACTION_SELECTED = "add-action-selected-latest",
    SVG_DELETE = "delete-latest",
    SVG_DELETE_SELECTED = "delete-selected-latest",
    SVG_REMOVE_LINK = "remove-link-latest",
    SVG_REMOVE_LINK_SELECTED = "remove-link-selected-latest",

    SVG_ZOOMIN = "zoomin",
    SVG_ZOOMOUT = "zoomout",

    SVG_CONFLICT = "conflict",
    SVG_REMOVE_CONFLICT = "remove-conflict",
    SVG_HIGHLIGHT_CONFLICT = "highlight-conflict";
    

export function getSVG(name) {
    return SVG_BASEDIR + name + SVG_SUFFIX;
}

export function getImage(name) {
    return PNG_BASEDIR + name + PNG_SUFFIX;
}

