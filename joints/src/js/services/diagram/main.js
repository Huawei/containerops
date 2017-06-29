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

define(['app','services/diagram/api'], function(app) {
    app.provide.factory("diagramService", ["diagramApiService","$state", function(diagramApiService,$state) {
        var workflowData = [
            {
                "name":"",
                "id":"",
                "type":"edit-stage",
                "runMode":"parallel", //serial串行，parallel并行
                "timeout":0,
                // "runResult":3,
                "actions":[
                    {
                        "isChosed":false,
                        "name":"",
                        "id":"action-",
                        "type":"action",
                        "timeout":0,
                        "components":[]
                    }
                ]
            },
            {
                "name":"addStage",
                "id":"stage-",
                "type":"add-stage",
                "runMode":"",
                "timeout":0,
                "actions":[]
            },
            {
                "name":"endStage",
                "id":"stage-",
                "type":"end-stage",
                "runMode":"",
                "timeout":0,
                "actions":[]
            }
        ];

        // var currentStageInfo = '';
        var currentStageIndex = '';
        var currentActionIndex = '';
        var currentComponentIndex = '';

        function drawWorkflow(scope,selector,dataset) {
            var baseSize = {
                svgpad : 50,
                stagePad : 110,
                stageWidth : 24,
                stageHeight : 24,
                actionPad : 20,
                actionToTop : 10, 
                componentWidth : 15,
                componentHeight : 15,
                componentPad : 3,
                gatherWidth : 40,
                gatherHeight : 50,
                rowActionNum : 4,
                runModeWidth : 10,
                addComponentWidth : 10,
                arcPadBig : 9,
                arcPadSmall : 3,
                addIconPad : 3,
                lineWidth : 6,
                addStageWidth: 20,
                addStageHeight: 20,
                addActionWidth: 25,
                addActionHeight: 25,
                stageLength : dataset.length
            };
            
            var currentDragIndex;
            var endPointIndex;
            var elementType = '';

            var baseColor = {
                stageBorder: '#d7d7d7',
                stageOrigin: '#65baf7',
                stageChosed: '#ff3333',
                componentOrigin: '#43b594'
            };

            

            var _this = baseSize;

            var zoom = d3.behavior.zoom()
                .scaleExtent([0.5, 2])
                .on("zoom", zoomed);

            function zoomed() {
                d3.select(this).attr("transform", 
                    "translate(" + d3.event.translate + ")scale(" + d3.event.scale + ")");
            };


            d3.select(selector)
                .selectAll('svg')
                .remove();

            var allGroups = d3.select(selector)
                .append('svg')
                .attr('width','100%')
                .attr('height','100%')
                .attr('fill','#fff')
                .append('g')
                .attr('transform','translate(0,0)')
                .call(zoom)
                .on('dblclick.zoom', null);


            var lines = allGroups.append('g')
                .attr('id','lines');
                
            var stagelines = lines.append('g')
                .attr('class','stagelines')

            var itemStage = allGroups.selectAll('.item-stage')
                .data(dataset)
                .enter()
                .append('g')
                .attr('class','item-stage')
                .attr('data-stageIndex',function(d,i){
                    return i;
                })
                .attr('data-type',function(d){
                    return d.type;
                })
                .attr('translateX',function(d,i){
                    return i*(_this.stageWidth+_this.stagePad)+_this.svgpad;
                })
                .attr('translateY',function(d,i){
                    return _this.stageHeight*2
                })
                .attr('transform',function(d,i){
                    return 'translate('+(i*(_this.stageWidth+_this.stagePad)+_this.svgpad)+','+(_this.stageHeight*2)+')';
                })
                // .call(drag);

            // add stage image
            itemStage.append('svg:image')
                .attr('width',_this.stageWidth)
                .attr('height',_this.stageHeight)
                .attr('class','stage-pic')
                .attr('href',function(d){
                    if(d.type === 'add-stage'){
                        return 'assets/images/icon-add-stage.svg';
                    }else if(d.type === 'end-stage'){
                        return 'assets/images/icon-stage-empty.svg';
                    }else if(d.type === 'edit-stage'){
                        return d.runMode === 'parallel' ? 'assets/images/icon-action-parallel.svg' : 'assets/images/icon-action-serial.svg';
                    }
                })
                .attr('data-name',function(d){
                    return d.name;
                })
                .attr('data-id',function(d){
                    return d.id;
                })
                .attr('data-type',function(d){
                    return d.type;
                })
                .each(function(d,i){
                    d.translateX = i*(_this.stageWidth+_this.stagePad)+_this.svgpad;
                    d.translateY = _this.stageHeight*2;
                })
                .on('click',scope.chosedStage);

            // add stage line & actions & components
            d3.selectAll('.item-stage')
                .each(function(d,i){
                    // add stage line
                    if(i !==0 ){
                        stagelines.append('path')
                            .attr('class','stage-line')
                            .attr('d',function(){
                                var x = d.translateX;
                                var y = d.translateY + _this.stageHeight / 2;
                                return 'M'+(x+_this.stageWidth)+' '+y+'L'+(x - _this.stagePad)+' '+y;
                            })
                            .attr('fill','none')
                            .attr('stroke',baseColor.stageBorder)
                            .attr('stroke-width','2');
                    };

                    // add actions
                    // action start y point
                    var currentActionY = _this.stageHeight+_this.actionPad;
                    var stageElement = d3.select(this);
                    var itemAction = stageElement
                        .selectAll('.item-action')
                        .data(d.actions)
                        .enter()
                        .append('g')
                        .attr('class','item-action')
                        .attr('data-stageIndex',function(){
                            return i;
                        })
                        .attr('data-actionIndex',function(a,ai){
                            return ai;
                        })
                        .attr('data-runMode',d.runMode)
                        .attr('translateY',function(a,ai){
                            var y = currentActionY;
                            var padding = _this.componentPad * 3;
                            var componentRows = a.components.length>0 ? Math.ceil(a.components.length/_this.rowActionNum): 1;
                            currentActionY += componentRows*(_this.componentHeight+_this.componentPad) - _this.componentPad + padding * 2 + _this.actionPad;

                            y = i%2===0 ? y + _this.actionToTop : y; 

                            return y;
                        })
                        .attr('transform',function(a,ai){ // action的y轴起点
                            var translateY = d3.select(this).attr('translateY');
                            return 'translate(0,'+translateY+')';
                        })
                        .on('click',scope.chosedAction);


                    // add components
                    var isChosedPreAction = '';
                    itemAction.each(function(a,ai){
                        // component start y point
                        var currentComponentY = 0;
                        var perAction = d3.select(this); 

                        // every component x point and y point
                        a.components.map(function(c,r){
                            var remain = r % _this.rowActionNum;
                            var componentNum = a.components.length>=_this.rowActionNum ? _this.rowActionNum : a.components.length;
                            var moveright = (_this.stageWidth - (_this.componentWidth+_this.componentPad)*_this.rowActionNum)/2 + _this.componentPad/2;
                            c.x = remain * (_this.componentWidth + _this.componentPad) + moveright ;

                            if(r%_this.rowActionNum===0){
                                currentComponentY += 1;
                            }

                            if(r===0){
                                currentComponentY = 0;
                            }

                            c.y = currentComponentY * (_this.componentPad + _this.componentHeight) + _this.componentPad*3;
                        });

                        // action borders
                        perAction.append('path')
                            .attr('class','borderLine')
                            .attr('fill','#fff')
                            .attr('stroke',baseColor.stageBorder)
                            .attr('stroke-width','1')
                            .attr('data-stageIndex',i)
                            .attr('data-actionIndex',ai)
                            .attr('d',function(){
                                var length = a.components.length;
                                var padding = _this.componentPad * 3;

                                var x = 0 + _this.stageWidth/2 - (_this.componentWidth + _this.componentPad) * 2 + _this.componentPad/2;
                                var y0 = 0;
                                var y1 = 0 + padding + _this.componentHeight + padding;

                                if(length>0){
                                    // x = a.components[0].x;
                                    // y0 = a.components[0].y - padding;
                                    y1 = a.components[length-1].y + _this.componentHeight + padding;
                                };


                                var x0 = x - padding;
                                var x1 = x + _this.rowActionNum * (_this.componentWidth + _this.componentPad) - _this.componentPad + padding;
                                var x2 = x0 + (x1 - x0) / 2; //每个stage的中心点
                                var y2 = i%2===0 ? y0 - _this.stageHeight/2 - _this.actionPad - _this.actionToTop : y0 - _this.stageHeight/2 - _this.actionPad; //弧线控制点
                                var x3 = x2 + _this.stageWidth/2 + _this.stagePad/2; //stage-line arc center point
                                var x4 = x2 - _this.stageWidth/2 - _this.stagePad/2; //stage-line arc center point

                                var x5 = x1 + _this.lineWidth;
                                var x6 = x5 + _this.arcPadBig;
                                var x7 = x0 - _this.lineWidth;
                                var x8 = x7 - _this.arcPadBig;
                                var x9 = x0 + (x1 - x0)/2; //action center x point
                                var y3 = y0 + (y1 - y0)/2; //action center y point

                                var y4 = y3 - _this.arcPadBig;

                                var translateY = perAction.attr('translateY');
                                var y5 = y0 - translateY + _this.stageHeight/2 + _this.arcPadBig; // stage-line center y point

                                var lineToRight = 'L'+x5+' '+y3;
                                var lineToLeft = 'L'+x7+' '+y3;
                                var arcToRight = 'Q'+x6+' '+y3+' '+x6+' '+y4;
                                var arcToLeft = 'Q'+x8+' '+y3+' '+x8+' '+y4;

                                var bordRightBottom = 'M'+x1+' '+y3+'L'+x1+' '+(y1 - _this.arcPadSmall)+'Q'+x1+' '+y1+' '+(x1 - _this.arcPadSmall)+' '+y1;
                                var bordBottom = 'L'+(x0 + _this.arcPadSmall)+' '+y1+'Q'+x0+' '+y1+' '+x0+' '+(y1 - _this.arcPadSmall);
                                var bordLeftBottom = 'L'+x0+' '+y3;
                                var bottom = bordRightBottom + bordBottom + bordLeftBottom;

                                var bordLeftTop = 'M'+x0+' '+y3+'L'+x0+' '+(y0 + _this.arcPadSmall)+'Q'+x0+' '+y0+' '+(x0 + _this.arcPadSmall)+' '+y0;
                                var bordTop = 'L'+x9+' '+y0+'L'+(x1 - _this.arcPadSmall)+' '+y0+'Q'+x1+' '+y0+' '+x1+' '+(y0 + _this.arcPadSmall);
                                var bordRightTop = 'L'+x1+' '+y3;
                                var top = bordLeftTop + bordTop + bordRightTop;
                                return bottom + top;
                            });
                        
                        // arc lines
                        perAction.append('path')
                            .attr('class','arcLine')
                            .attr('fill','none')
                            .attr('stroke',baseColor.stageBorder)
                            .attr('stroke-width','1')
                            .attr('data-stageIndex',i)
                            .attr('data-actionIndex',ai)
                            .attr('d',function(){
                                var length = a.components.length;
                                var padding = _this.componentPad * 3;

                                var x = 0 + _this.stageWidth/2 - (_this.componentWidth + _this.componentPad) * 2 + _this.componentPad/2;
                                var y0 = 0;
                                var y1 = 0 + padding + _this.componentHeight + padding;

                                if(length>0){
                                    // x = a.components[0].x;
                                    // y0 = a.components[0].y - padding;
                                    y1 = a.components[length-1].y + _this.componentHeight + padding;
                                };

                                var x0 = x - padding;
                                var x1 = x + _this.rowActionNum * (_this.componentWidth + _this.componentPad) - _this.componentPad + padding;
                                var x2 = x0 + (x1 - x0) / 2; //每个stage的中心点
                                var y2 = i%2===0 ? y0 - _this.stageHeight/2 - _this.actionPad - _this.actionToTop : y0 - _this.stageHeight/2 - _this.actionPad; //弧线控制点
                                var x3 = x2 + _this.stageWidth/2 + _this.stagePad/2; //stage-line arc center point
                                var x4 = x2 - _this.stageWidth/2 - _this.stagePad/2; //stage-line arc center point

                                var x5 = x1 + _this.lineWidth;
                                var x6 = x5 + _this.arcPadBig;
                                var x7 = x0 - _this.lineWidth;
                                var x8 = x7 - _this.arcPadBig;
                                var x9 = x0 + (x1 - x0)/2; //action center x point
                                var y3 = y0 + (y1 - y0)/2; //action center y point

                                var y4 = y3 - _this.arcPadBig;

                                var translateY = perAction.attr('translateY');
                                var y5 = y0 - translateY + _this.stageHeight/2 + _this.arcPadBig; // stage-line center y point

                                var lineToRight = 'L'+x5+' '+y3;
                                var lineToLeft = 'L'+x7+' '+y3;
                                var arcToRight = 'Q'+x6+' '+y3+' '+x6+' '+y4;
                                var arcToLeft = 'Q'+x8+' '+y3+' '+x8+' '+y4;

                                var bordRightBottom = 'M'+x1+' '+y3+'L'+x1+' '+(y1 - _this.arcPadSmall)+'Q'+x1+' '+y1+' '+(x1 - _this.arcPadSmall)+' '+y1;
                                var bordBottom = 'L'+(x0 + _this.arcPadSmall)+' '+y1+'Q'+x0+' '+y1+' '+x0+' '+(y1 - _this.arcPadSmall);
                                var bordLeftBottom = 'L'+x0+' '+y3;
                                var bottom = bordRightBottom + bordBottom + bordLeftBottom;

                                var dottedEnd = i%2===0&&ai===0 ? y0 - _this.actionPad - _this.actionToTop : y0 - _this.actionPad;
                                var dottedHeight = Math.abs(dottedEnd / 4) - 2;
                                var dottedToTop = 'L'+x9+' '+(y0 - dottedHeight)+'M'+x9+' '+(y0 - dottedHeight - 2)+'L'+x9+' '+(y0 - dottedHeight*2 - 4)+'M'+x9+' '+(y0 - dottedHeight*2 - 6)+'L'+x9+' '+(y0 - dottedHeight*3 - 6);
                                var bordLeftTop = 'M'+x0+' '+y3+'L'+x0+' '+(y0 + _this.arcPadSmall)+'Q'+x0+' '+y0+' '+(x0 + _this.arcPadSmall)+' '+y0;
                                // var toDottedEnd = ai === 0 ? dottedEnd : dottedEnd + _this.addIconPad;
                                if(isChosedPreAction && ai !== 0){
                                    dottedEnd = dottedEnd + _this.addIconPad
                                }
                                isChosedPreAction = a.isChosed;
                                // else if(!a.isChosed){
                                //     dottedEnd = dottedEnd
                                // }
                                var bordTop = 'L'+x9+' '+y0+'L'+x9+' '+dottedEnd+'M'+x9+' '+y0+'L'+(x1 - _this.arcPadSmall)+' '+y0+'Q'+x1+' '+y0+' '+x1+' '+(y0 + _this.arcPadSmall);
                                
                                if(d.runMode === 'parallel'&&ai!==0){
                                    bordTop = 'L'+x9+' '+y0+'L'+(x1 - _this.arcPadSmall)+' '+y0+'Q'+x1+' '+y0+' '+x1+' '+(y0 + _this.arcPadSmall);
                                }

                                var bordRightTop = 'L'+x1+' '+y3;
                                var top = bordLeftTop + bordTop + bordRightTop;

                                var commonArcToRight = 'L'+x6+' '+(y2 + _this.arcPadBig)+'Q'+x6+' '+y2+' '+(x6 + _this.arcPadBig)+' '+y2+'L'+x3+' '+y2;
                                var commonArcToLeft = 'L'+x8+' '+(y2 + _this.arcPadBig)+'Q'+x8+' '+y2+' '+(x8 - _this.arcPadBig)+' '+y2+'L'+x4+' '+y2;
                                

                                if(i===0){
                                    if(ai === 0){
                                        return bottom + top + lineToRight + arcToRight + commonArcToRight;
                                    }
                                    return bottom + top + lineToRight + arcToRight + 'L'+x6+' '+y5;
                                }else{

                                    if(ai === 0){
                                        return bottom + lineToLeft + arcToLeft + commonArcToLeft + top + lineToRight + arcToRight + commonArcToRight;
                                    }
                                    return bottom + lineToLeft + arcToLeft +'L'+x8+' '+ y5 + top + lineToRight + arcToRight + 'L'+x6+' '+y5;
                                }
                            });
    
                        // action components
                        perAction.selectAll('.item-component')
                            .data(a.components)
                            .enter()
                            .append('rect')
                            .attr('class','item-component')
                            .attr('width',_this.componentWidth)
                            .attr('height',_this.componentHeight)
                            .attr('data-stageIndex',function(){
                                return i;
                            })
                            .attr('data-actionIndex',function(){
                                return ai;
                            })
                            .attr('data-componentIndex',function(c,ci){
                                return ci;
                            })
                            .attr('data-name',function(c){
                                return c.name;
                            })
                            .attr('data-id',function(c){
                                return c.id;
                            })
                            .attr('data-type',function(c){
                                return c.type;
                            })
                            .attr('x',function(c,r){
                                var remain = r % _this.rowActionNum;
                                var componentNum = a.components.length>=_this.rowActionNum ? _this.rowActionNum : a.components.length;
                                var moveright = (_this.stageWidth - (_this.componentWidth+_this.componentPad)*_this.rowActionNum)/2 + _this.componentPad/2;
                                c.x = remain * (_this.componentWidth + _this.componentPad) + moveright;
                                return c.x;
                            })
                            .attr('y',function(c,r){
                                if(r%_this.rowActionNum===0){
                                    currentComponentY += 1;
                                }

                                if(r===0){
                                    currentComponentY = 0;
                                }

                                c.y = currentComponentY * (_this.componentPad + _this.componentHeight) + _this.componentPad*3;
                                return c.y;
                            })
                            .attr('fill',function(a,r){
                                return baseColor.componentOrigin;
                            })

                        // action-no-components x point & y point
                        var actionLength = d.actions.length;
                        var length = a.components.length;
                        var padding = _this.componentPad * 3;
                        var x = 0 + _this.stageWidth/2 - (_this.componentWidth + _this.componentPad) * 2 + _this.componentPad/2;
                        var y0 = 0;
                        var y1 = 0 + padding + _this.componentHeight + padding;

                        if(length>0){
                            // x = a.components[0].x;
                            // y0 = a.components[0].y - padding;
                            y1 = a.components[length-1].y + _this.componentHeight + padding;
                        };

                        var x0 = x - padding;
                        var x1 = x + _this.rowActionNum * (_this.componentWidth + _this.componentPad) - _this.componentPad + padding;


                        // action-add icon 
                        if(a.isChosed){
                            // top icon
                            perAction.append('svg:image')
                                .attr('width',_this.addComponentWidth)
                                .attr('height',_this.addComponentWidth)
                                .attr('class','add-component')
                                .attr('href',function(d,i){
                                    return 'assets/images/icon-add-action.svg'
                                })
                                .attr('data-stageIndex',function(){
                                    return i;
                                })
                                .attr('data-actionIndex',function(){
                                    return ai;
                                })
                                .attr('x',function(){
                                    return x0 + (x1 - x0 - _this.addComponentWidth) / 2;
                                })
                                .attr('y',function(){
                                    return y0 - _this.addComponentWidth/2 - _this.arcPadSmall + _this.addIconPad;
                                })
                                .on('click',scope.addTopAction);

                            // bottom icon 
                            perAction.append('svg:image')
                                .attr('width',_this.addComponentWidth)
                                .attr('height',_this.addComponentWidth)
                                .attr('class','add-component')
                                .attr('href',function(d,i){
                                    return 'assets/images/icon-add-action.svg'
                                })
                                .attr('data-stageIndex',function(){
                                    return i;
                                })
                                .attr('data-actionIndex',function(){
                                    return ai;
                                })
                                .attr('x',function(){
                                    return x0 + (x1 - x0 - _this.addComponentWidth) / 2;
                                })
                                .attr('y',function(){
                                    return y1 - _this.addComponentWidth/2;
                                })
                                .on('click',scope.addBottomAction);
                        };                        

                    });
                })
        }; 

        function resetWorkflowData(newData){
            this.workflowData = newData;
        };

        return {
            "workflowData": workflowData,
            // "currentStageInfo": currentStageInfo,
            "currentStageIndex": currentStageIndex,
            "currentActionIndex": currentActionIndex,
            "currentComponentIndex": currentComponentIndex,
            "drawWorkflow": drawWorkflow,
            "resetWorkflowData": resetWorkflowData
        }
    }])
})
