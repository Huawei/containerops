/* 
Copyright 2014 Huawei Technologies Co., Ltd. All rights reserved.

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

export let PIPELINE_START = "pipeline-start",
    PIPELINE_END = "pipeline-end",
    PIPELINE_ADD_STAGE = "pipeline-add-stage",
    PIPELINE_ADD_ACTION = "pipeline-add-action",
    PIPELINE_STAGE = "pipeline-stage",
    PIPELINE_ACTION = "pipeline-action",


    svgStageWidth = 45,
    svgStageHeight = 52,
    svgActionWidth = 38,
    svgActionHeight = 38,

    svgButtonWidth = 30,
    svgButtonHeight = 30,


    pipelineView = null,
    actionsView = null,
    actionView = [],
    buttonView = null,
    linesView = null,
    lineView = [],
    clickNodeData = {},
    linePathAry = [],

    sequencePipelineView = null,
    sequenceActionsView = null,
    sequenceActionLinkView = null,
    sequenceActionView = [],
    sequenceLinesView = null,
    sequenceLineView = [],
    sequenceLinePathArray = [],
    sequenceRunData = [],

    PipelineNodeSpaceSize = 200,
    ActionNodeSpaceSize = 75,

    pipelineNodeStartX = 0,
    pipelineNodeStartY = 0,

    svgWidth = 0,
    svgHeight = 0,
    svgMainRect = null,
    svg = null,
    g = null,

    popupWidth = 110,
    popupHeight = 25,
    currentSelectedItem = null,
    
    zoomScale = 1,
    zoomTargetScale = 1,
    zoomFactor = 0.2,
    zoomMinimum = 0.1,
    zoomMaximum = 3,
    zoomDuration = 300;

export function setPipelineView(v) {
    pipelineView = v;
}

export function setActionsView(v) {
    actionsView = v;
}

export function setActionView(v) {
    actionView = v;
}

export function setButtonView(v) {
    buttonView = v;
}

export function setLinesView(v) {
    linesView = v;
}

export function setLineView(v) {
    lineView = v;
}

export function setClickNodeData(v) {
    clickNodeData = v;
}

export function setLinePathAry(v) {
    linePathAry = v;
}


export function setPipelineNodeSpaceSize(v) {
    PipelineNodeSpaceSize = v;
}

export function setActionNodeSpaceSize(v) {
    ActionNodeSpaceSize = v;
}


export function setPipelineNodeStartX(v) {
    pipelineNodeStartX = v;
}

export function setPipelineNodeStartY(v) {
    pipelineNodeStartY = v;
}


export function setSvgWidth(v) {
    svgWidth = v;
}

export function setSvgHeight(v) {
    svgHeight = v;
}

export function setSvgMainRect(v) {
    svgMainRect = v;
}

export function setSvg(v) {
    svg = v;
}

export function setG(v) {
    g = v;
}

export function setCurrentSelectedItem(v) {
    currentSelectedItem = v;
}
