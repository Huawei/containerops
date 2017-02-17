define(["app"], function(app) {
    app.controllerProvider.register('WorkflowCreateController', ['$scope', '$state', 'notifyService', function($scope, $state, notifyService) {
        $scope.backToList = function() {
            $state.go("workflow");
        };
        $scope.saveWorkflow = function(){
            $state.go("workflow");
        };
        // $scope.drawWorkflow = function() {
        //     var svg = d3.select("#div-d3-main-svg")
        //         .append("svg")
        //         .attr("width", "100%")
        //         .attr("height", "100%");
        //     var g = svg.append("g");
        //     var svgMainRect = g.append("rect")
        //         .attr("width", "100%")
        //         .attr("height", "100%")
        //         .attr("fill", "white");
        //     var svgMainRect = g.append("circle")
        //         .attr("cx", 0) 
        //         .attr("cy", 0)
        //         .attr("r", 20)
        //         .attr("fill", "green")
        //         .attr("transform", function(d, i) {
        //             return "translate(" +150 + "," + 150 + ")";
        //         })
        //         .attr("cursor","pointer")
        //         .on("click", function() {
        //            notifyService.notify("click stage","success");
        //         });
        // };

        $scope.svgpad = 50;
        $scope.stagePad = 140;
        $scope.stageWidth = 30;
        $scope.stageHeight = 30;
        $scope.actionPad = 15;
        $scope.actionToTop = 0; 
        $scope.componentWidth = 15;
        $scope.componentHeight = 15;
        $scope.componentPad = 4;
        $scope.gatherWidth = 40;
        $scope.gatherHeight = 50;
        $scope.rowActionNum = 4;
        $scope.runModeWidth = 10;
        $scope.addComponentWidth = 10;
        $scope.arcPadBig = 10;
        $scope.arcPadSmall = 3;
        $scope.addIconPad = 3;
        $scope.lineWidth = 10;

        $scope.dataset = [
            {
                "name":"stage0",
                "id":"s0",
                "type":"start-stage",
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
                "type":"end-stage",
                "runMode":"",
                "actions":[]
            }
        ];
        $scope.stageLength = $scope.dataset.length;
        $scope.drawWorkflow = function() {
            var currentDragIndex;
            var endPointIndex;
            var elementType = '';
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

            var _this = $scope;
            var svg = d3.select('#div-d3-main-svg')
                .append('svg')
                .attr('width','100%')
                .attr('height','100%')
                .attr('fill','#fff');

            var lines = svg.append('g')
                .attr('id','lines');
                
            var stagelines = lines.append('g')
                .attr('class','stagelines')

            var itemStage = svg.selectAll('.item-stage')
                .data(this.dataset)
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

            // add stage circle
            itemStage.append('circle')
                .attr('cx',_this.stageWidth/2)
                .attr('cy',_this.stageHeight/2)
                .attr('r',_this.stageHeight/2)
                .attr('stroke','#d7d7d7')
                .attr('fill',function(d){
                    return d.type === 'end-stage' ? '#fff' : '#d7d7d7';
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
                            .attr('stroke','#d7d7d7')
                            .attr('stroke-width','2');
                    };

                    // add actions
                    // action start y point
                    var currentActionY = _this.stageHeight+_this.actionPad;
                    var itemAction = d3.select(this)
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
                        .attr('transform',function(a,ai){ // action的y轴起点
                            var y = currentActionY;
                            var padding = _this.componentPad * 3;
                            var componentRows = Math.ceil(a.components.length/_this.rowActionNum);
                            currentActionY += componentRows*(_this.componentHeight+_this.componentPad) - _this.componentPad + padding * 2 + _this.actionPad;

                            y = i%2===0 ? y + _this.actionToTop : y; 

                            return 'translate(0,'+y+')';

                        });

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
                                return '#d7d7d7';
                            });

                        // action borders
                        perAction.append('path')
                            .attr('class','gatherLine')
                            .attr('fill','none')
                            .attr('stroke','#d7d7d7')
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
                                var x3 = x2 + _this.stageWidth/2 + _this.stagePad/2; //弧线在stage-line的中点
                                var x4 = x2 - _this.stageWidth/2 - _this.stagePad/2; //弧线在stage-line的中点

                                var x5 = x1 + _this.lineWidth;
                                var x6 = x5 + _this.arcPadBig;
                                var x7 = x0 - _this.lineWidth;
                                var x8 = x7 - _this.arcPadBig;
                                var x9 = x0 + (x1 - x0)/2; //action x轴中点
                                var y3 = y0 + (y1 - y0)/2; //action y轴中点

                                var y4 = y3 - _this.arcPadBig;

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
                                    return bottom + lineToLeft + arcToLeft +'L'+x8+' '+(y2 - _this.actionPad - _this.arcPadBig*2) + top + lineToRight + arcToRight + 'L'+x6+' '+(y2 - _this.actionPad - _this.arcPadBig*2);
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
                                    return 'assets/images/parallel.jpg'
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
                                    return 'assets/images/serial.jpg'
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
                                    return 'assets/images/addComponent.jpg'
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

                    });
                })
        };

        $scope.drawWorkflow();




    }]);
})
