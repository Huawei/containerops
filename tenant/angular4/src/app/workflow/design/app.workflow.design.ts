import { Component, ViewChild, ElementRef, OnInit, ViewEncapsulation } from '@angular/core';
import * as d3 from 'd3';

@Component({
  selector: 'app-workflow-design',
  templateUrl: './app.workflow.design.html',
  styleUrls: ['./app.workflow.design.css'],
  encapsulation: ViewEncapsulation.None


})
export class WorkflowDesignComponent implements OnInit  {
  @ViewChild('design') private designContainer: ElementRef;
  private design: any;

  constructor() { }

  ngOnInit() {
    this.createChart();
  }
  createChart() {

    let element = this.designContainer.nativeElement;
    let svg = d3.select(element).append('svg');

    this.design = svg.append('g');
  }
}
