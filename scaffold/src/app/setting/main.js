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

import * as systemSettingDataService from "./settingData";
import { notify, confirm } from "../common/notify";
import { loading } from "../common/loading";

export let systemSettings;
let settings,tempSettings;

export function initSystemSettings(next){
    var promise = systemSettingDataService.getAllSystemSettings();
    promise.done(function(data) {
        loading.hide();
        settings = data.setting;
        systemSettings = _.pairs(settings);       
        if(next){
            next();
        }  
    });
    promise.fail(function(xhr, status, error) {
        loading.hide();
        if (!_.isUndefined(xhr.respaonseJSON) && xhr.responseJSON.errMsg) {
            notify(xhr.responseJSON.errMsg, "error");
        } else if(xhr.statusText != "abort") {
            notify("Server is unreachable", "error");
        }
    });
}

export function initSystemSettingPage(){
    tempSettings = _.pairs($.extend(true,{},settings));
    $.ajax({
        url: "../../templates/setting/settingList.html",
        type: "GET",
        cache: false,
        success: function(data) {
            $("#main").html($(data));
            $("#settinglist").show("slow");

            $(".savesystemsetting").on('click', function() {
                var promise = systemSettingDataService.saveSystemSettings(_.object(tempSettings));
                promise.done(function(data) {
                    loading.hide();
                    notify(data.message, "success");
                    initSystemSettings(initSystemSettingPage);
                });
                promise.fail(function(xhr, status, error) {
                    loading.hide();
                    if (!_.isUndefined(xhr.responseJSON) && xhr.responseJSON.errMsg) {
                        notify(xhr.responseJSON.errMsg, "error");
                    } else if(xhr.statusText != "abort") {
                        notify("Server is unreachable", "error");
                    }
                });
            })

            $(".settinglist_body").empty();
            _.each(tempSettings,function(item,index){
                 var row = `<tr data-index="`+ index +`">
                                <td style="width:50%">`+ item[0] +`</td>
                                <td style="width:50%"><input type="text" value="` + item[1] + `" class="form-control system-setting-value" required></td>
                            </tr>`;
                $(".settinglist_body").append(row);
            });

            $(".system-setting-value").on('blur',function(event){
                var index = $(event.currentTarget).parent().parent().data("index");
                tempSettings[index][1] = $(event.currentTarget).val();
            });
        }
    });
}

let settingMappings = {
    "KUBE_NODE_IP" : "setupData.action.ip",
    "KUBE_APISERVER_IP" : "setupData.action.apiserver"
}

export function doMapping(action){
    var mappings = _.pairs(settingMappings);
    _.each(mappings,function(mapping){
        var paths = mapping[1].split(".");
        var target = action[paths[0]];
        for(var i =1; i< paths.length-1; i++){
            target = target[paths[i]];
        }

        var value = _.find(systemSettings,function(setting){
            return setting[0] == mapping[0];
        })[1];

        if(_.isEmpty(target[paths[paths.length-1]])){
            target[paths[paths.length-1]] = value;
        }
    })
}

