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

import * as actionEnvData from "./actionEnvData";
import {workflowVars} from "../workflow/workflowVar";

export function initActionEnv(action){
    actionEnvData.getActionEnvData(action);

    showActionEnvKVs();  
}

function showActionEnvKVs(){
    $("#action-envs").empty();
    _.each(actionEnvData.data,function(item,index){
        var row = '<div class="port-row"><div class="env-key-div">'
                        +'<div>'
                            +'<label for="normal-field" class="col-sm-4 control-label">'
                                +'Key'
                            +'</label>'
                            +'<div class="col-sm-7" data-index="' + index + '">'
                                +'<input type="text" value="' + item.key + '" class="form-control c-env-key allowFromVar" required>'
                            +'</div>'
                        +'</div>'
                    +'</div>'
                    +'<div class="env-value-div">'
                        +'<div>'
                            +'<label for="normal-field" class="col-sm-4 control-label">'
                                +'Value'
                            +'</label>'
                            +'<div class="col-sm-7" data-index="' + index + '">' 
                                +'<input type="text" class="form-control c-env-value allowFromVar" required>'
                            +'</div>'
                        +'</div>'
                    +'</div>'
                    +'<div class="port-remove-div c-rm-kv" data-index="' + index + '">'
                        +'<span class="glyphicon glyphicon-remove"></span>'
                    +'</div></div>';
        $("#action-envs").append(row);
        $("#action-envs").find("div[data-index="+index+"]").find(".c-env-value").val(item.value);
    });
    
    var addrow = `<button type="button" class="btn btn-success c-new-kv">
                        <i class="glyphicon glyphicon-plus" style="top:1px"></i>
                        <span style="margin-left:5px">Add Env</span>
                    </button>`;
    $("#action-envs").append(addrow);

    $(".c-new-kv").on('click',function(){
        actionEnvData.addEnv();
        showActionEnvKVs();
    });

    $(".c-env-key").on('input',function(event){
        var key = $(event.currentTarget).val();
        $(event.currentTarget).val(key.toUpperCase());
    });

    $(".c-env-key").on('blur',function(event){
        actionEnvData.setEnvKey(event);
    });

    $(".c-env-value").on('blur',function(event){
        actionEnvData.setEnvValue(event);
    });

    $(".c-rm-kv").on('click',function(event){
        actionEnvData.removeEnv(event);
        showActionEnvKVs();
    }); 

    var globalvars = _.map(workflowVars,function(item){
                        return "@"+item[0]+"@";
                    });
    $(".allowFromVar").autocomplete({
        source:[globalvars],
        limit: 100,
        visibleLimit: 5
    }); 
}