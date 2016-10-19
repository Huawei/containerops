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

export function mouseoverRelevantPipeline(thisData){
    var pathAry = d3.selectAll("#pipeline-line-view path")[0];
    pathAry.forEach(function(i){
       try{
            var _path = d3.select(i),
                _class = _path.attr("class");
            if(!!_class){
                // _path.attr("stroke-opacity","0.1");
            }
           
            if(_class.indexOf(thisData.id) == 0){
                i.parentNode.appendChild(i);
                _path.attr("stroke-opacity","1");
            }
       }catch(e){

       }
      
    })
}


export function mouseoutRelevantPipeline(){
    var pathAry = d3.selectAll("#pipeline-line-view path")[0];
    pathAry.forEach(function(i){
        var _path = d3.select(i),
             _class = _path.attr("class");
        if(!!_class){
            _path.attr("stroke-opacity","0.2");
         }
      
    })
}