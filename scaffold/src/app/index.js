import {initPipelinePage} from "./pipeline/main";
import {initComponentPage} from "./component/main";
// import {historyRecord} from "./historyRecord";

// let $a = d3.select("#showHistory").on("click",historyRecord);

$._messengerDefaults = {
    extraClasses: 'messenger-fixed messenger-theme-future messenger-on-top'
}

initPipelinePage();

$(".menu-pipeline").on('click',function(event){
    initPipelinePage();
    $(event.currentTarget).parent().parent().children().removeClass("active");
    $(event.currentTarget).parent().addClass("active");
})

$(".menu-component").on('click',function(){
    initComponentPage();
    $(event.currentTarget).parent().parent().children().removeClass("active");
    $(event.currentTarget).parent().addClass("active");
})
// initActionLinkView();

function initActionLinkView() {
    actionLinkView.append("rect")
        .attr("x",10)
        .attr("y",10)
        .attr("rx",10)
        .attr("ry",10)
        .attr("width",120)
        .attr("height",40)
        .attr("stroke","red")
        .attr("fill","red")
    ;
}

