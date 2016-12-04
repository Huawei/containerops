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

import * as systemSettingDataService from "./settingData";
import { notify, confirm } from "../common/notify";
import { loading } from "../common/loading";

export let systemSettings;

export function initSystemSettingPage() {
    var promise = systemSettingDataService.getAllSystemSettings();
    promise.done(function(data) {
        loading.hide();
        systemSettings = _.pairs(data.setting);
        showSystemSettings();
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

function showSystemSettings(){
    $.ajax({
        url: "../../templates/setting/settingList.html",
        type: "GET",
        cache: false,
        success: function(data) {
            $("#main").html($(data));
            $("#settinglist").show("slow");

            $(".savesystemsetting").on('click', function() {
                var promise = systemSettingDataService.saveSystemSettings(_.object(systemSettings));
                promise.done(function(data) {
                    loading.hide();
                    notify(data.message, "success");
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
            _.each(systemSettings,function(item,index){
                 var row = `<tr data-index="`+ index +`">
                                <td style="width:50%">`+ item[0] +`</td>
                                <td style="width:50%"><input type="text" value="` + item[1] + `" class="form-control system-setting-value" required></td>
                            </tr>`;
                $(".settinglist_body").append(row);
            });

            $(".system-setting-value").on('blur',function(event){
                var index = $(event.currentTarget).parent().parent().data("index");
                systemSettings[index][1] = $(event.currentTarget).val();
            });
        }
    });
}