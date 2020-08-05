import React, {useEffect} from 'react';
import {Chart, Pie} from 'react-chartjs-2';

const BarChart = (props) => {
	const data = {
		labels: [
			'VMs/App Services',
			'Virtual Disks',
			'Storage & Database',
			'Networking',
			'Other'
		],
		datasets: [{
			data: ['0.35', '0.19', '0.22', '0.21', '0.03'],
			backgroundColor: [
			'#FF6384',
			'#36A2EB',
			'#FFCE56',
			'#DC143C',
			'#9932CC'
			],
			hoverBackgroundColor: [
			'#FF6384',
			'#36A2EB',
			'#FFCE56',
			'#DC143C',
			'#9932CC'
			]
		}]
	}

	const options = {
		title: {
			display:true,
			text: 'Cost Comparison'    
		},
	}
	return (
		<Pie
			height={130}
			data={data}     
		/>
	)
}

export default BarChart
 


