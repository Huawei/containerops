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

export let WORKFLOW_START = "workflow-start",
    WORKFLOW_END = "workflow-end",
    WORKFLOW_ADD_STAGE = "workflow-add-stage",
    WORKFLOW_ADD_ACTION = "workflow-add-action",
    WORKFLOW_STAGE = "workflow-stage",
    WORKFLOW_ACTION = "workflow-action",


    svgStageWidth = 45,
    svgStageHeight = 52,
    svgActionWidth = 38,
    svgActionHeight = 38,

    svgButtonWidth = 30,
    svgButtonHeight = 30,


    workflowView = null,
    actionsView = null,
    actionView = [],
    buttonView = null,
    linesView = null,
    lineView = [],
    clickNodeData = {},
    linePathAry = [],

    sequenceWorkflowView = null,
    sequenceActionsView = null,
    sequenceActionLinkView = null,
    sequenceActionView = [],
    sequenceButtonView = null,
    sequenceLinesView = null,
    sequenceLineView = [],
    sequenceLinePathArray = [],
    sequenceRunData = [],
    sequenceRunStatus = null,
    refreshSequenceRunData = [],

    WorkflowNodeSpaceSize = 200,
    ActionNodeSpaceSize = 75,

    workflowNodeStartX = 0,
    workflowNodeStartY = 0,

    svgWidth = 0,
    svgHeight = 0,
    svgMainRect = null,
    svg = null,
    g = null,

    popupWidth = 110,
    popupHeight = 22,
    currentSelectedItem = null,
    
    zoomFactor = 0.2,
    zoomMinimum = 0.1,
    zoomMaximum = 3,
    zoomDuration = 300,

    buttonWidth = 23,
    buttonHeight = 23,
    buttonVerticalSpace = 6,
    toolTipBackground = "#555",
    buttonHorizonSpace = 20,
    rectBackgroundY = 15;


export function setWorkflowView(v) {
    workflowView = v;
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


export function setWorkflowNodeSpaceSize(v) {
    WorkflowNodeSpaceSize = v;
}

export function setActionNodeSpaceSize(v) {
    ActionNodeSpaceSize = v;
}


export function setWorkflowNodeStartX(v) {
    workflowNodeStartX = v;
}

export function setWorkflowNodeStartY(v) {
    workflowNodeStartY = v;
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
