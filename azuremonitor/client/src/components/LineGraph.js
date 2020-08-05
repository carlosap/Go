import React  from 'react';
import {Line} from 'react-chartjs-2';


const data = {
	labels: ['Year 1', 'Year 2', 'Year 3'],
	datasets: [
		{
			label: 'Current Cost',
			fill: false,
			lineTension: 0.3,
			backgroundColor: 'rgba(75,192,192,0.4)',
			borderColor: '#DC143C',
			borderCapStyle: 'butt',
			borderDash: [],
			borderDashOffset: 0.0,
			borderJoinStyle: 'miter',
			pointBorderColor: '#DC143C',
			pointBackgroundColor: '#fff',
			pointBorderWidth: 5,
			pointHoverRadius: 5,
			pointHoverBackgroundColor: '#DC143C',
			pointHoverBorderColor: '#DC143C',
			pointHoverBorderWidth: 2,
			pointRadius: 1,
			pointHitRadius: 10,
			data: [20000000, 20000000, 20000000]
		},
		{
			label: 'Cost After Optimization',
			fill: false,
			lineTension: 0.3,
			backgroundColor: 'rgba(75,192,192,0.4)',
			borderColor: 'rgba(75,192,192,1)',
			borderCapStyle: 'butt',
			borderDash: [],
			borderDashOffset: 0.0,
			borderJoinStyle: 'miter',
			pointBorderColor: 'rgba(75,192,192,1)',
			pointBackgroundColor: '#fff',
			pointBorderWidth: 5,
			pointHoverRadius: 5,
			pointHoverBackgroundColor: 'rgba(75,192,192,1)',
			pointHoverBorderColor: 'rgba(220,220,220,1)',
			pointHoverBorderWidth: 2,
			pointRadius: 1,
			pointHitRadius: 10,
			data: [19000000, 18000000, 15000000]
		},
		{
			label: 'Savings',
			fill: false,
			lineTension: 0.3,
			backgroundColor: '#00FF7F',
			borderColor: '#00FF7F',
			borderCapStyle: 'butt',
			borderDash: [],
			borderDashOffset: 0.0,
			borderJoinStyle: 'miter',
			pointBorderColor: '#00FF7F',
			pointBackgroundColor: '#fff',
			pointBorderWidth: 5,
			pointHoverRadius: 5,
			pointHoverBackgroundColor: 'rgba(75,192,192,1)',
			pointHoverBorderColor: 'rgba(220,220,220,1)',
			pointHoverBorderWidth: 2,
			pointRadius: 1,
			pointHitRadius: 10,
			data: [1000000, 2000000, 5000000]
		}
	]
};

const options= {
	tooltips: {
		callbacks: {
			label: function(t, d) {
				var xLabel = d.datasets[t.datasetIndex].label;
				var yLabel = t.yLabel >= 1000 ? '$' + t.yLabel.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",") : '$' + t.yLabel;
				return xLabel + ': ' + yLabel;
			}
		}
	},
	scales: {
		yAxes: [{
			ticks: {
				callback: function(value, index, values) {
					if (parseInt(value) >= 1000) {
						return '$' + value.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
					} else {
						return '$' + value;
					}
				},
				suggestedMax: 25000000,
			}
		}]
	}
}

const LineGraph = (props) => {

    return (
        <Line data={data} options={options} />
    )
}

export default LineGraph