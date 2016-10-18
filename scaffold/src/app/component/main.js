import {getAllComponents,getComponent,addComponent,addComponentVersion,saveComponent} from "./componentData";
import {initComponentIO} from "./componentIO";
import {initComponentSetup} from "./componentSetup";

export let allComponents;

export let componentData;
let componentName, componentVersion;

export function initComponentPage(){
    // handle promise

    // to be removed
    allComponents = getAllComponents();
    if(allComponents.length>0){
        showComponentList();
    }else{
        showNoComponent();
    }
}

function showComponentList(){
    $.ajax({
        url: "../../templates/component/componentList.html",
        type: "GET",
        cache: false,
        success: function (data) {
            $("#main").html($(data));    
            $("#componentlist").show("slow");

            $(".newcomponent").on('click',function(){
                showNewComponent();
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
                            +'<td><button type="button" class="btn btn-primary ppview">View</button></td></tr>';
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

            $(".ppview").on("click",function(event){
                var target = $(event.currentTarget);
                componentName = target.parent().parent().data("pname");
                componentVersion = target.parent().parent().data("version");
                showComponentDesigner();
            })
        }
    });
}

function showNoComponent(){
    $.ajax({
        url: "../../templates/component/noComponent.html",
        type: "GET",
        cache: false,
        success: function (data) {
            $("#main").html($(data));    
            $("#nocomponent").show("slow");
            $(".newcomponent").on('click',function(){
                showNewComponent();
            })  
        }
    });
}

export function showNewComponent(fromPipeline){
    $.ajax({
        url: "../../templates/component/newComponent.html",
        type: "GET",
        cache: false,
        success: function (data) {
            $("#main").children().hide();
            $("#main").append($(data));    
            $("#newcomponent").show("slow");
            $("#newComponentBtn").on('click',function(){
                // addPipeline();

                // to be removed below
                if(addComponent()){
                    initComponentPage();
                }  
            })
            $("#cancelNewComponentBtn").on('click',function(){
                if(fromPipeline){
                    $(".menu-component").parent().removeClass("active");
                    $(".menu-pipeline").parent().addClass("active");
                }
                cancelNewComponentPage();
            })
        }
    });
}

function showComponentDesigner(){  
    $.ajax({
        url: "../../templates/component/componentDesign.html",
        type: "GET",
        cache: false,
        success: function (data) {
            $("#main").html($(data));    
            $("#componentdesign").show("slow"); 

            componentData = getComponent(componentName,componentVersion);

            $("#selected_component").text(componentName + " / " + componentVersion); 

            $(".backtolist").on('click',function(){
                initComponentPage();
            });

            $(".savecomponent").on('click',function(){
                saveComponent(componentName, componentVersion, componentData);
            });

            $(".newcomponentversion").on('click',function(){
                showNewComponentVersion();
            });

            $(".newcomponent").on('click',function(){
                showNewComponent();
            });

            initComponentEdit();
        }
    }); 
}

function initComponentEdit(){
    $.ajax({
        url: "../../templates/component/componentEdit.html",
        type: "GET",
        cache: false,
        success: function (data) {
            $("#componentDesigner").html($(data));

            initComponentSetup(componentData);

            initComponentIO(componentData);

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

function showNewComponentVersion(){
    $.ajax({
        url: "../../templates/component/newComponentVersion.html",
        type: "GET",
        cache: false,
        success: function (data) {
            $("#main").children().hide();
            $("#main").append($(data));    
            $("#newcomponentversion").show("slow"); 

            $("#c-name-newversion").val(componentName);

            $("#newComponentVersionBtn").on('click',function(){
                // addPipelineVersion(pipelineVersion);

                // to be removed below
                if(addComponentVersion(componentVersion)){
                    initComponentPage();
                } 
            })
            $("#cancelNewComponentVersionBtn").on('click',function(){
                cancelNewComponentVersionPage();
            })      
        }
    }); 
    
    $("#content").hide();
    $("#nocomponent").hide();
    $("#newcomponent").hide();
    $("#newcomponentversion").show("slow");
}

function cancelNewComponentPage(){
    $("#newcomponent").remove();
    $("#main").children().show("slow");
}

function cancelNewComponentVersionPage(){
    $("#newcomponentversion").remove();
    $("#main").children().show("slow");
}

// $("#pipeline-select").on('change',function(){
//     showVersionList();
// })
// $("#version-select").on('change',function(){
//     showPipeline();
// })

// function showPipelineList(){
//     $("#pipeline-select").empty();
//     d3.select("#pipeline-select")
//         .selectAll("option")
//         .data(allPipelines)
//         .enter()
//         .append("option")
//         .attr("value",function(d,i){
//             return d.name;
//         })
//         .text(function(d,i){
//             return d.name;
//         }); 
//      $("#pipeline-select").select2({
//        minimumResultsForSearch: Infinity
//      });   
//     showVersionList();
// }

// function showVersionList(){
//     var pipeline = $("#pipeline-select").val();
//     var versions = _.find(allPipelines,function(item){
//         return item.name == pipeline;
//     }).versions;

//     $("#version-select").empty();
//     d3.select("#version-select")
//         .selectAll("option")
//         .data(versions)
//         .enter()
//         .append("option")
//         .attr("value",function(d,i){
//             return d.version;
//         })
//         .text(function(d,i){
//             return d.version;
//         }); 
//     $("#version-select").select2({
//        minimumResultsForSearch: Infinity
//      });
    
//     versions_shown = versions;

//     showPipeline(); 
// }
