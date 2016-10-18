
export let  PIPELINE_START = "pipeline-start",
            PIPELINE_END = "pipeline-end",
            PIPELINE_ADD_STAGE = "pipeline-add-stage",
            PIPELINE_ADD_ACTION = "pipeline-add-action",
            PIPELINE_STAGE = "pipeline-stage",
            PIPELINE_ACTION = "pipeline-action",


			svgStageWidth = 45,
		    svgStageHeight = 42,
		    svgActionWidth = 30,
		    svgActionHeight = 28,

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

		    PipelineNodeSpaceSize = 200,
		    ActionNodeSpaceSize = 75,

		    pipelineNodeStartX = 0,
		    pipelineNodeStartY = 0,

		    svgWidth = 0,
		    svgHeight = 0,
		    svgMainRect = null,
		    svg = null,
		    g = null;

export	function setPipelineView(v){
	pipelineView = v;
}

export	function setActionsView(v){
	actionsView = v;
}

export	function setActionView(v){
	actionView = v;
}

export	function setButtonView(v){
	buttonView = v;
}

export	function setLinesView(v){
	linesView = v;
}

export	function setLineView(v){
	lineView = v;
}

export	function setClickNodeData(v){
	clickNodeData = v;
}

export	function setLinePathAry(v){
	linePathAry = v;
}


export	function setPipelineNodeSpaceSize(v){
	PipelineNodeSpaceSize = v;
}

export	function setActionNodeSpaceSize(v){
	ActionNodeSpaceSize = v;
}


export	function setPipelineNodeStartX(v){
	pipelineNodeStartX = v;
}

export	function setPipelineNodeStartY(v){
	pipelineNodeStartY = v;
}


export	function setSvgWidth(v){
	svgWidth = v;
}

export	function setSvgHeight(v){
	svgHeight = v;
}

export	function setSvgMainRect(v){
	svgMainRect = v;
}

export	function setSvg(v){
	svg = v;
}

export	function setG(v){
	g = v;
}
