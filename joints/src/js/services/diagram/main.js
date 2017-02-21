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
define(['app','services/diagram/api'], function(app) {
    app.provide.factory("diagramService", ["diagramApiService","$state", function(diagramApiService,$state) {
        var dataset = [
            {
                "name":"stage0",
                "id":"s0",
                "type":"edit-stage",
                "runMode":"serial", //serial串行，parallel并行
                // "runResult":3,
                "actions":[
                    {
                        "components":[
                            {
                                "name":"action0",
                                "id":"s0-at0",
                                "type":"action",
                                "inputData":"",
                                "outputData":""
                            }
                        ]
                    },
                    {
                        "components":[
                            {
                                "name":"action0",
                                "id":"s0-at0",
                                "type":"action",
                                "inputData":"",
                                "outputData":""
                            }
                        ]
                    }
                ]
            },
            {
                "name":"stage1",
                "id":"s1",
                "type":"edit-stage",
                "runMode":"parallel",
                "actions":[
                    {
                        "components":[
                            {
                                "name":"action0",
                                "id":"s1-at0",
                                "type":"action",
                                "inputData":"",
                                "outputData":""
                            },
                            {
                                "name":"action1",
                                "id":"s1-at1",
                                "type":"action",
                                "inputData":"",
                                "outputData":""
                            },
                            {
                                "name":"action2",
                                "id":"s1-at2",
                                "type":"action",
                                "inputData":"",
                                "outputData":""
                            },
                            {
                                "name":"action3",
                                "id":"s1-at3",
                                "type":"action",
                                "inputData":"",
                                "outputData":""
                            },
                            {
                                "name":"action4",
                                "id":"s1-at4",
                                "type":"action",
                                "inputData":"",
                                "outputData":""
                            },
                            {
                                "name":"action0",
                                "id":"s1-at0",
                                "type":"action",
                                "inputData":"",
                                "outputData":""
                            },
                            {
                                "name":"action1",
                                "id":"s1-at1",
                                "type":"action",
                                "inputData":"",
                                "outputData":""
                            },
                            {
                                "name":"action2",
                                "id":"s1-at2",
                                "type":"action",
                                "inputData":"",
                                "outputData":""
                            },
                            {
                                "name":"action3",
                                "id":"s1-at3",
                                "type":"action",
                                "inputData":"",
                                "outputData":""
                            },
                            {
                                "name":"action4",
                                "id":"s1-at4",
                                "type":"action",
                                "inputData":"",
                                "outputData":""
                            }
                        ]
                    },
                    {
                        "components":[
                            {
                                "name":"action0",
                                "id":"s1-at0",
                                "type":"action",
                                "inputData":"",
                                "outputData":""
                            },
                            {
                                "name":"action3",
                                "id":"s1-at3",
                                "type":"action",
                                "inputData":"",
                                "outputData":""
                            },
                            {
                                "name":"action4",
                                "id":"s1-at4",
                                "type":"action",
                                "inputData":"",
                                "outputData":""
                            }
                        ]
                    }
                ]
            },
            {
                "name":"stage2",
                "id":"s2",
                "type":"edit-stage",
                "runMode":"parallel",
                "actions":[
                    {
                        "components":[
                            {
                                "name":"action0",
                                "id":"s2-at0",
                                "type":"action",
                                "inputData":"",
                                "outputData":""
                            },
                            {
                                "name":"action1",
                                "id":"s2-at1",
                                "type":"action",
                                "inputData":"",
                                "outputData":""
                            }
                        ]
                    }
                ]
            },
            {
                "name":"stage3",
                "id":"s3",
                "type":"add-stage",
                "runMode":"",
                "actions":[]
            },
            {
                "name":"stage3",
                "id":"s3",
                "type":"end-stage",
                "runMode":"",
                "actions":[]
            }
        ];

        function drawWorkflow(selector,dataset) {
            var baseSize = {
                svgpad : 50,
                stagePad : 110,
                stageWidth : 26,
                stageHeight : 26,
                actionPad : 15,
                actionToTop : 0, 
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
            }
            
            var currentDragIndex;
            var endPointIndex;
            var elementType = '';

            var baseColor = {
                stageBorder: '#d7d7d7',
                stageOrigin: '#65baf7',
                stageChosed: '#ff3333',
                componentOrigin: '#43b594'
            };

            var chosedStageIndex = '';
            var chosedActionIndex = '';
            // var drag = d3.behavior.drag()
            //     .on("drag", dragmove)
            //     .on('dragstart',function(d,i){
            //         currentDragIndex = i;
            //     })
            //     .on('dragend',sort); 
                    
            // function dragmove(d,i) {
            //     elementType = d3.select(this).attr('data-type');
            //     if(elementType !== 'end-stage'){
            //         d3.select(this)
            //           .attr("translateX", d.translateX = d3.event.x )
            //           .attr("transform", 'translate('+d.translateX+','+d.translateY+')');
            //     }   
            // };

            // function sort(d,i){
            //     if(currentDragIndex&&elementType!=='end-stage'){

            //         var dragTranslateX = d3.select(this)
            //             .attr("translateX");

            //         var stages = d3.selectAll('.item-stage')
            //             .each(function(d,i){
            //                 // origin translateX
            //                 var preTranslateX = i*(_this.stageWidth+_this.stagePad)+_this.svgpad;
            //                 var nextTranslateX = (i+1)*(_this.stageWidth+_this.stagePad)+_this.svgpad;

            //                 if(currentDragIndex !== i){
            //                     var stageCenterX = preTranslateX+_this.stageWidth/2;

            //                     if(dragTranslateX>=stageCenterX&&(dragTranslateX)<nextTranslateX){
            //                         dealSortData(currentDragIndex,i);
            //                         return;
            //                     }
            //                 }
            //                 // if(currentDragIndex&&currentX>d.translateX&&currentX<)
            //                 // console.log(d.translateX,(d.translateX+_this.stageWidth+_this.stagePad)+_this.svgpad)
            //                 // console.log(i*(_this.stageWidth+_this.stagePad)+_this.svgpad)
            //             })
            //      // console.log(currentX)
                    
            //     }
            // };

            // function dealSortData(currentDragIndex,endPointIndex){
            //     console.log(currentDragIndex,endPointIndex)
            // }

            var newStage = {
                "name":"",
                "id":"",
                "type":"edit-stage",
                "runMode":"",
                "actions":[]
            };

            var newComponent = {
                "name":"action1",
                "id":"s2-at1",
                "type":"action",
                "inputData":"",
                "outputData":""
            };

            var newAction = {
                "components":[
                    {
                        "name":"action0",
                        "id":"s2-at0",
                        "type":"action",
                        "inputData":"",
                        "outputData":""
                    }
                ]
            };

            var _this = baseSize;

            function chosedStage(d,i){
                clearChosedStageIndex();
                clearChosedStageColor();

                if(d.type === 'add-stage'){
                    var stage = angular.copy(newStage);
                    var addstage = angular.copy(dataset[i]);
                    var endstage = angular.copy(dataset[i+1]);
                    dataset[i] = stage;
                    dataset[i+1] = addstage;
                    dataset[i+2] = endstage;
                    // dataset.splice(dataset[i-1],0,stage)
                    drawWorkflow(selector,dataset); 
                };

                if(d.type === 'edit-stage'){
                    d3.select(this).select('circle').attr('fill','#e43937');
                    chosedStageIndex = i;
                    $state.go("workflow.create.stage",{"id": d.id});
                };
            };

            function clearChosedStageColor(){
                d3.selectAll('circle')
                    .attr('fill',function(d){
                        if(d.type === 'end-stage' || d.type === 'add-stage'){
                            return '#fff';
                        }
                        return baseColor.stageOrigin;
                    })
            };

            function clearChosedStageIndex(){
                chosedStageIndex = '';
            };

            function addComponent(){
                addElement(d3.select(this),'component');
            };

            function addAction(){
                addElement(d3.select(this),'action');
            };

            function addElement(currentElement,type){
                var chosedStageIndex = currentElement.attr('data-stageIndex');
                if(type === 'component'){
                    var chosedActionIndex = currentElement.attr('data-actionIndex');
                    var component = angular.copy(newComponent);
                    dataset[chosedStageIndex]['actions'][chosedActionIndex]['components'].push(component);

                }else{
                    var action = angular.copy(newAction);
                    dataset[chosedStageIndex]['actions'].push(action);
                };

                drawWorkflow(selector,dataset); 
            }; 


            d3.select(selector)
                .selectAll('svg')
                .remove();

            var svg = d3.select(selector)
                .append('svg')
                .attr('width','100%')
                .attr('height','100%')
                .attr('fill','#fff');


            var lines = svg.append('g')
                .attr('id','lines');
                
            var stagelines = lines.append('g')
                .attr('class','stagelines')

            var itemStage = svg.selectAll('.item-stage')
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
                .on('click',chosedStage)
                // .call(drag);

            // add stage circle
            itemStage.append('circle')
                .attr('cx',_this.stageWidth/2)
                .attr('cy',_this.stageHeight/2)
                .attr('r',_this.stageHeight/2)
                .attr('stroke',baseColor.stageOrigin)
                .attr('fill',function(d){
                    if(d.type === 'end-stage' || d.type === 'add-stage'){
                        return '#fff';
                    }
                    return baseColor.stageOrigin;
                })
                .attr('class','stage-pic')
                .attr('data-name',function(d){
                    return d.name
                })
                .attr('data-id',function(d){
                    return d.id
                })
                .attr('data-type',function(d){
                    return d.type
                })
                .each(function(d,i){
                    d.translateX = i*(_this.stageWidth+_this.stagePad)+_this.svgpad;
                    d.translateY = _this.stageHeight*2;
                });
               

            // add add-stage icon
            itemStage.each(function(d){
                if(d.type === 'add-stage'){
                    d3.select(this).append('svg:image')
                        .attr('width',_this.addStageWidth)
                        .attr('height',_this.addStageHeight)
                        .attr('href','assets/images/icon-add-stage.svg')
                        .attr('x',(_this.stageWidth/2 - _this.addStageWidth / 2))
                        .attr('y',(_this.stageHeight/2 - _this.addStageHeight / 2));
                }
            });

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
                            var componentRows = Math.ceil(a.components.length/_this.rowActionNum);
                            currentActionY += componentRows*(_this.componentHeight+_this.componentPad) - _this.componentPad + padding * 2 + _this.actionPad;

                            y = i%2===0 ? y + _this.actionToTop : y; 

                            return y;
                        })
                        .attr('transform',function(a,ai){ // action的y轴起点
                            var translateY = d3.select(this).attr('translateY');
                            return 'translate(0,'+translateY+')';
                        });

                    // action-add icon
                    if(d.type === 'edit-stage'){
                        var addActionGroup = stageElement.append('g')
                            .attr('class','add-action')
                            .attr('transform','translate(0,'+currentActionY+')');

                        addActionGroup.append('svg:image')
                            .attr('class','add-action-img')
                            .attr('width',_this.addActionWidth)
                            .attr('height',_this.addActionHeight)
                            .attr('href','assets/images/icon-add-action-old.png')
                            .attr('data-stageIndex',function(){
                                return i;
                            })
                            .on('click',addAction);

                        addActionGroup.append('path')
                            .attr('stroke',baseColor.stageBorder)
                            .attr('d',function(){
                                var centerPoint = 0 + _this.addActionWidth/2;
                                return 'M'+centerPoint+' '+0+ 'L'+centerPoint+' '+(0 - _this.actionPad);
                            });
                    };
                    

                    // add components
                    itemAction.each(function(a,ai){
                        // component start y point
                        var currentComponentY = 0;
                        var perAction = d3.select(this); 

                        // action item-component
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
                                var moveright = (_this.stageWidth - (_this.componentWidth+_this.componentPad)*_this.rowActionNum)/2;
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
                            });

                        // action borders
                        perAction.append('path')
                            .attr('class','gatherLine')
                            .attr('fill','none')
                            .attr('stroke',baseColor.stageBorder)
                            .attr('stroke-width','1')
                            .attr('d',function(){
                                var length = a.components.length;
                                var x = a.components[0].x;
                                var padding = _this.componentPad * 3;
                                var x0 = x - padding;
                                var y0 = a.components[0].y - padding;
                                var x1 = x + _this.rowActionNum * (_this.componentWidth + _this.componentPad) - _this.componentPad + padding;
                                var y1 = a.components[length-1].y + _this.componentHeight + padding;
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
                                    return bottom + top + lineToRight + arcToRight + 'L'+x6+' '+(y2 - _this.actionPad - _this.arcPadBig);
                                }else{

                                    if(ai === 0){
                                        return bottom + lineToLeft + arcToLeft + commonArcToLeft + top + lineToRight + arcToRight + commonArcToRight;
                                    }
                                    return bottom + lineToLeft + arcToLeft +'L'+x8+' '+ y5 + top + lineToRight + arcToRight + 'L'+x6+' '+y5;
                                }
                                
                            });

                        // action runMode
                        var actionLength = d.actions.length;
                        if(d.runMode === 'parallel'){
                            perAction.append('svg:image')
                                .attr('width',_this.runModeWidth)
                                .attr('height',_this.runModeWidth)
                                .attr('class','run-mode')
                                .attr('href',function(d,i){
                                    return 'assets/images/icon-action-parallel.svg'
                                })
                                .attr('x',function(){
                                    var x = a.components[0].x;
                                    var padding = _this.componentPad*3;
                                    var x0 = x - padding;
                                    return x0 -_this.runModeWidth/2;
                                })
                                .attr('y',function(){
                                    var length = a.components.length;
                                    var padding = _this.componentPad * 3;
                                    var y0 = a.components[0].y - padding;
                                    var y1 = a.components[length-1].y + _this.componentHeight + padding;
                                    return y0 + (y1 - y0 - _this.runModeWidth) / 2;
                                });
                        }else if(d.runMode === 'serial'&&ai!==(actionLength-1)){
                            perAction.append('svg:image')
                                .attr('width',_this.runModeWidth)
                                .attr('height',_this.runModeWidth)
                                .attr('class','run-mode')
                                .attr('href',function(d,i){
                                    return 'assets/images/icon-action-serial.svg'
                                })
                                .attr('x',function(){
                                    var length = a.components.length;
                                    var x = a.components[0].x;
                                    var padding = _this.componentPad*3;
                                    var x0 = x - padding;
                                    var x1 = x + _this.rowActionNum * (_this.componentWidth + _this.componentPad) - _this.componentPad + padding;
                                    return x0 + (x1 - x0 - _this.runModeWidth) / 2;
                                })
                                .attr('y',function(){
                                    var length = a.components.length;
                                    var padding = _this.componentPad * 3;
                                    var y1 = a.components[length-1].y + _this.componentHeight + padding;
                                    return y1 - _this.runModeWidth/ 2;
                                })
                        }

                        // component add icon 
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
                                var x = a.components[0].x;
                                var padding = _this.componentPad * 3;
                                var x1 = x + _this.rowActionNum * (_this.componentWidth + _this.componentPad) - _this.componentPad + padding;
                                return x1 - _this.addComponentWidth/2;
                            })
                            .attr('y',function(){
                                var padding = _this.componentPad * 3;
                                var y0 = a.components[0].y - padding;
                                return y0 - _this.addComponentWidth/2 + _this.addIconPad;
                            })
                            .on('click',addComponent);

                    });
                })

            

        }; 



        return {
            "dataset": dataset,
            "drawWorkflow": drawWorkflow
        }
    }])
})
