"use strict";
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};
var core_1 = require("@angular/core");
var d3 = require("d3");
var WorkflowDesignComponent = (function () {
    function WorkflowDesignComponent() {
    }
    WorkflowDesignComponent.prototype.ngOnInit = function () {
        this.createChart();
    };
    WorkflowDesignComponent.prototype.createChart = function () {
        var element = this.designContainer.nativeElement;
        var svg = d3.select(element).append('svg');
        this.design = svg.append('g');
    };
    return WorkflowDesignComponent;
}());
__decorate([
    core_1.ViewChild('design'),
    __metadata("design:type", core_1.ElementRef)
], WorkflowDesignComponent.prototype, "designContainer", void 0);
WorkflowDesignComponent = __decorate([
    core_1.Component({
        selector: 'app-workflow-design',
        templateUrl: './app.workflow.design.html',
        styleUrls: ['./app.workflow.design.css'],
        encapsulation: core_1.ViewEncapsulation.None
    }),
    __metadata("design:paramtypes", [])
], WorkflowDesignComponent);
exports.WorkflowDesignComponent = WorkflowDesignComponent;
//# sourceMappingURL=app.workflow.design.js.map