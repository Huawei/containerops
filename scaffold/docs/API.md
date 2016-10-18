# API Document
## Pipeline
### GET: /pipelineInit.do   (Get pipeline initial data) 
#### Request Parameters
    Empty
#### Response Message 
    {  
      "stageArray" : [
        {
          "id" : "",
          "type" : ""
        },
        {
          "id": "",
          "type" : "",
          "preStageId":"",
          "actions" : [
            {
              "id": "",
              "type": ""
            }
           ]
         }
       ],

       "linkPathArray" : [
         {
           "id" : "",
           "startId" : "",
           "endId" : ""
         }
        ]
    }
<br>
## Stage
### POST: /setStartStageSetupData.do   (Set config data for start stage) 
#### Request Parameters
    {
      "id" : "",
      "setupData" : {
        "gitUrl" : "",
        "gitTag" : "",
        "gitEvent" : "",
        "callbackUrl" : ""
      }
    }
#### Response Message
    {
      "message" : "error&success"
    }
<br>
### GET: /getStartStageSetupData.do   (Get config data of start stage) 
#### Request Parameters
    {
      "id" : ""
    }
#### Response Message
    {
      "setupData" : {
        "gitUrl" : "",
        "gitTag" : "",
        "gitEvent" : "",
        "callbackUrl" : ""
      } 
    }
<br>
### POST: /stage.do   (Add new stage or save edited stage) 
#### Request Parameters
    {
      "id" : "",
      "preStageId" : "",
      "setupData" : {
        "stageID" : "",
        "stageName" : "",
        "stageTimeout" : "",
        "stageEnv" : "",
        "callbackUrl" : ""
      }
    }
#### Response Message
    {
      "message" : "error&success",
      "stage" : {}
    }
<br>
### GET: /getStageSetupData.do   (Get stage config data) 
#### Request Parameters
    {
      "id" : ""
    }
#### Response Message
    {
      "setupData" : {
        "stageID" : "",
        "stageName" : "",
        "stageTimeout" : "",
        "stageEnv" : "",
        "callbackUrl" : ""
      } 
    } 
<br>
### DELETE: /deleteStage.do   (Delete a stage) 
#### Request Parameters
    {
      "id" : ""
    }
#### Response Message
    {
      "message" : "error&success"
    } 
<br>
## Action
### POST: /addAction.do   (Add an action) 
#### Request Parameters
    {
      "parentStageId" : ""
    }
#### Response Message
    {
      "id": "",
      "type": "",
      "typeConfig" : {},
      "setupData" : {}
    } 
<br>
### DELETE: /deleteAction.do   (Delete an action) 
#### Request Parameters
    {
      "id" : ""
    }
#### Response Message
    {
      "message" : "error & success"
    } 
<br>
### POST: /setActionSetupData.do   (Set action config data) 
#### Request Parameters
    {
      "id" : "",
      "setupData" : {
        "actionId" : "",
        "actionName" :"",
        "actionTimeout" : "",
        "actionEnv" : "",
        "actionImage" : ""
      }
    }
#### Response Message
    {
      "message" : "error & success"
    } 
<br>
### GET: /getActionSetupData.do   (Get action config data) 
#### Request Parameters
    {
      "id" : ""
    }
#### Response Message
    {
      "actionId" : "",
      "actionName" :"",
      "actionTimeout" : "",
      "actionEnv" : "",
      "actionImage" : ""
    } 
<br>
### POST: /setActionInputOutput.do   (Set action input and output data) 
#### Request Parameters
    {
     "id" : "",
      "input" : {},
      "output" : {}
    }
#### Response Message
    {
      "message" : "error&success"
    } 
<br>
### GET: /getActionInputOutput.do   (Get action input and output data) 
#### Request Parameters
    {
     "id" : ""
    }
#### Response Message
    {
      "input" : {},
      "output" : {}
    } 
<br>
## Relationship
### POST: /setLinkPath.do   (Add relationship between actions) 
#### Request Parameters
    {
      "startId" : "",
      "endId" : ""
    }
#### Response Message
    {
      "id" : ""
    } 
<br>
### DELETE: /setLinkPath.do   (Delete relationship between actions) 
#### Request Parameters
    {
     "id" : ""
    }
#### Response Message
    {
      "message" : "error&success"
    } 
<br>
### POST: /setLinkPathConfig.do   (Set input and output data of relationship) 
#### Request Parameters
    {
      "id" : "",
      "config" : {},
      "input" : [
        {
          "key" : ".action",
          "type" : "object",
          "path" : ".action",
          "childNode" : [
            {
              "key" : "",
              "type" : "",
              "path" : "",
              "childNode" : []
            }
          ]
        }
      ],
      "output" : [],
      "relation" : [
        {
          "from" : "",
          "fromShow" : "",
          "to" : "",
          "toShow" : "",
          "isToEqual" : true,
          "isFromEqual" : true,
          "child" : [
            {
              "from" : "",
              "fromShow" : "",
              "to" : "",
              "toShow" : "",
              "isToEqual" : true,
              "isFromEqual" : true,
              "child" : []
            }
          ]
        }
      ]
    }
#### Response Message
    {
      "message" : "error&success"
    } 
<br>
### GET: /getLinkPathConfig.do   (Get input and output data of relationship) 
#### Request Parameters
    {
     "id" : ""
    }
#### Response Message
    {
      "config" : {},
      "input" : [
        {
          "key" : ".action",
          "type" : "object",
          "path" : ".action",
          "childNode" : [
            {
              "key" : "",
              "type" : "",
              "path" : "",
              "childNode" : []
            }
          ]
        }
      ],
      "output" : [],
      "relation" : [
        {
          "from" : "",
          "fromShow" : "",
          "to" : "",
          "toShow" : "",
          "isToEqual" : true,
          "isFromEqual" : true,
          "child" : [
            {
              "from" : "",
              "fromShow" : "",
              "to" : "",
              "toShow" : "",
              "isToEqual" : true,
              "isFromEqual" : true,
              "child" : []
            }
          ]
        }
      ]
    }
