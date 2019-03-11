import { Component, AfterViewInit, ElementRef,ViewChild} from '@angular/core';
import ResizeObserver from 'resize-observer-polyfill';
import { Chart } from 'chart.js'


@Component({
  selector: 'app-climate-chart',
  templateUrl: './climate-chart.component.html',
  styleUrls: ['./climate-chart.component.css']
})
export class ClimateChartComponent implements AfterViewInit {

    canvas: any;
	ctx: any;
	chart : any;

	@ViewChild('climateChartMonitor')
	public monitor: ElementRef


	constructor() { }

	ngAfterViewInit() {
		let ro = new ResizeObserver(entries => {
			for ( let e of entries) {
				const cr = e.contentRect;
				this.chart.options.width = cr.width;
				this.chart.options.height = cr.height;
				this.chart.resize();
			}
		});
		ro.observe(this.monitor.nativeElement);
		this.canvas = document.getElementById('climateChart');
		this.ctx = this.canvas.getContext('2d');
		this.chart = new Chart(this.ctx,{
			type: 'line',
			data: {
				labels: ['0','1','2','3','4'],
				datasets: [
					{
						label: 'Humidity',
						fill: false,
						lineTension: 0,
						data: [
							{x:0,y:40},
							{x:1,y:42},
							{x:2,y:38},
							{x:3,y:50},
							{x:4,y:50.2}
						]
					},
					{
						label: 'Temperature',
						fill: false,
						lineTension: 0,
						data: [
							{x:0,y:20},
							{x:1,y:20.3},
							{x:2,y:19.8},
							{x:3,y:21.2},
							{x:4,y:20.8}
						]
					}
				],
			},
			options: {
				responsive: false,
			}

		});
        //todo display a chart
    }



}
