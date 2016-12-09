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

export let cron = {

    entries : [],

    $el : "",

    period : "",

    initCronEntry : function(targetdom,cronEntry){
        this.entries = cronEntry.split(" ");
        setPeriod();
        this.$el = targetdom;
        this.showCronEntry();
    },

    showCronEntry : function(){
        var self = this;
        $.ajax({
            url: "../../templates/cron/cron.html",
            type: "GET",
            cache: false,
            success: function (data) {
                self.$el.empty();
                self.$el.append($(data));    
                self.blockDisplayControl();
                self.setSelections();
            }
        });
    },

    setPeriod : function(){
        if(entries[0] == "*"){
            this.period = "minute";
        }else if(entries[1] == "*"){
            this.period = "hour";
        }else if(entries[2] == "*" && entries[4] == "*"){
            this.period = "day";
        }else if(entries[2] == "*" && entries[3] == "*"){
            this.period = "week";
        }else if(entries[3] == "*" && entries[4] == "*"){
            this.period = "month";
        }else if(entries[4] == "*"){
            this.period = "year";
        }                   
    },

    blockDisplayControl : function(){
        this.findDomElement(".cron-block").hide();

        switch(this.period){
            case "minute": 
                break;
            case "hour" : 
                this.findDomElement(".cron-block-mins").show();
                break;
            case "day" :
                this.findDomElement(".cron-block-time").show();
                break;
            case "week" :
                this.findDomElement(".cron-block-dow").show();
                this.findDomElement(".cron-block-time").show();
                break;
            case "month" :
                this.findDomElement(".cron-block-dom").show();
                this.findDomElement(".cron-block-time").show();
                break;
            case "year" :
                this.findDomElement(".cron-block-dom").show();
                this.findDomElement(".cron-block-month").show();
                this.findDomElement(".cron-block-time").show();
                break;
        }
    },

    setSelections : function(){
        var self = this;

        self.findDomElement(".cron-period-type").val(self.period);
        self.findDomElement(".cron-period-type").select2({
            minimumResultsForSearch: Infinity
        });
        self.findDomElement(".cron-period-type").on("change",function(){
            self.period = self.findDomElement(".cron-period-type").val();
            self.blockDisplayControl();
            self.setSelections();
        });

        switch(self.period){
            case "minute": 
                break;
            case "hour" : 
                self.selectionControl(".cron-mins",0);
                break;
            case "day" :
                self.selectionControl(".cron-time-min",0);
                self.selectionControl(".cron-time-hour",1);
                break;
            case "week" :
                self.selectionControl(".cron-dow",4);
                self.selectionControl(".cron-time-min",0);
                self.selectionControl(".cron-time-hour",1);
                break;
            case "month" :
                self.selectionControl(".cron-dom",2);
                self.selectionControl(".cron-time-min",0);
                self.selectionControl(".cron-time-hour",1);
                break;
            case "year" :
                self.selectionControl(".cron-dom",2);
                self.selectionControl(".cron-month",3);
                self.selectionControl(".cron-time-min",0);
                self.selectionControl(".cron-time-hour",1);
                break;
        }
    },

    findDomElement : function(selector){
        return this.$el.find(selector);
    },

    selectionControl : function(selector,index){
        var self = this;
        self.findDomElement(selector).val(self.entries[index]);
        self.findDomElement(selector).select2({
            minimumResultsForSearch: Infinity
        });
        self.findDomElement(selector).off("change");
        self.findDomElement(selector).on("change",function(){
            self.entries[index] = self.findDomElement(selector).val();
            self.showEntry();
        });

        self.showEntry();
    },

    showEntry : function(){
        this.findDomElement(".cron-val").text(this.getEntry());
    },

    getEntry : function(){
        return this.entries.join(" ");
    }
};