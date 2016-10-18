import * as constant from "../common/constant";
import {initPipeline} from "../pipeline/initPipeline";
import {initAction} from "../pipeline/initAction";
import {initLine} from "../pipeline/initLine";
import {pipelineData} from "../pipeline/main";
import {resizeWidget} from "../theme/widget";
import {pipelineEdit} from "../relation/pipelineEdit";
import {removeLinkArray} from "../relation/removeLinkArray";
import {initActionIO} from "./actionIO";
import {initActionSetup} from "./actionSetup";
import {getAllComponents,getComponent} from "../component/componentData";
import {showNewComponent} from "../component/main";
import {notify} from "../common/notify";

export function clickAction(sd, si) {
    $.ajax({
        url: "../../templates/action/actionMain.html",
        type: "GET",
        cache: false,
        success: function (data) {
            $("#pipeline-info-edit").html($(data));

            $(".actionfromcomponent").on('click',function(){
                getComponents(sd);
            });

            $(".actionnocomponent").on('click',function(){
                showActionEditor(sd);
            });

            resizeWidget();       
        }
    });


}

function showActionEditor(action){
    $.ajax({
        url: "../../templates/action/actionEdit.html",
        type: "GET",
        cache: false,
        success: function (data) {
            $("#actionMain").html($(data));

            initActionSetup(action);

            initActionIO(action);

            $("#uuid").attr("value", action.id);

            // view select init
            $("#action-component-select").select2({
                minimumResultsForSearch: Infinity
            });
            $("#k8s-service-protocol").select2({
                minimumResultsForSearch: Infinity
            });     
        }
    });
}

let allComponents;
function getComponents(action){
    allComponents = getAllComponents();

    if(allComponents.length > 0){
        $.ajax({
            url: "../../templates/action/actionComponentList.html",
            type: "GET",
            cache: false,
            success: function (data) {
                $("#actionMain").html($(data));    

                $(".newcomponent").on('click',function(){
                    $(".menu-component").parent().addClass("active");
                    $(".menu-pipeline").parent().removeClass("active");
                    showNewComponent(true);
                })

                $(".componentlist_body").empty();
                _.each(allComponents,function(item){
                    var pprow = '<tr style="height:50px"><td class="pptd">'
                    +'<span class="glyphicon glyphicon-menu-down treeclose" data-name="'+item.name+'"></span>&nbsp;'
                    +'<span class="glyphicon glyphicon-menu-right treeopen" data-name="'+item.name+'"></span>&nbsp;' 
                    + item.name + '</td><td></td><td></td></tr>';
                    $(".componentlist_body").append(pprow);
                    _.each(item.versions,function(version){
                        var vrow = '<tr data-pname="' + item.name + '" data-version="' + version.version + '" style="height:50px">'
                        +'<td></td><td class="pptd">' + version.version + '</td>'
                        +'<td><button type="button" class="btn btn-primary cload">Load</button></td></tr>';
                        $(".componentlist_body").append(vrow);
                    })
                }) ;

                $(".treeclose").on("click",function(event){
                    var target = $(event.currentTarget);
                    target.hide();
                    target.next().show();
                    var name = target.data("name");
                    $('*[data-pname='+name+']').hide();
                });

                $(".treeopen").on("click",function(event){
                    var target = $(event.currentTarget);
                    target.hide();
                    target.prev().show();
                    var name = target.data("name");
                    $('*[data-pname='+name+']').show();
                });

                $(".cload").on("click",function(event){
                    var target = $(event.currentTarget);
                    var componentName = target.parent().parent().data("pname");
                    var componentVersion = target.parent().parent().data("version");
                    LoadComponentToAction(componentName,componentVersion,action);
                })
            }
        });      
    }else{
        notify("You have no components to reuse, please go to 'Component' to create one.","warning");
        showActionEditor(action);
    }
    
}

function LoadComponentToAction(componentName,componentVersion,action){
    var component = getComponent(componentName,componentVersion);
    action.setupData = $.extend(true,{},component.setupData);
    action.inputJson = $.extend(true,{},component.inputJson);
    action.outputJson = $.extend(true,{},component.outputJson);
    showActionEditor(action);
}

function jsonChanged(root,json){
    root.val(JSON.stringify(json));
}