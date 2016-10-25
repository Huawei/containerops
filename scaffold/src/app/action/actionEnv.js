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

import * as actionEnvData from "./actionEnvData";

export function initActionEnv(action){
    actionEnvData.getActionEnvData(action);

    showActionEnvKVs();  

    $(".c-new-kv").on('click',function(){
        actionEnvData.addEnv();
        showActionEnvKVs();
    });
}

function showActionEnvKVs(){
    $("#action-envs").empty();
    _.each(actionEnvData.data,function(item,index){
        var row = '<tr data-index="'+index+'"><td>'
                    +'<input type="text" class="form-control col-md-5 c-env-key" value="'+item.key+'" required>'
                    + '</td><td>'
                    +'<input type="text" class="form-control col-md-5 c-env-value" required value='+item.value+'>'
                    + '</td><td>'
                    +'<span class="glyphicon glyphicon-minus c-rm-kv"></span>'
                    +'</td></tr>';
        $("#action-envs").append(row);
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
}