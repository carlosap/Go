import React, {useState} from 'react'
import {Table, TableHead, TableBody, TableRow, TableCell} from '@material-ui/core'
import {Box, IconButton,  Paper, Typography} from '@material-ui/core'
import _ from 'lodash'

const FilteredTable = ({filter, data}) => {
	const headerTexts = [
		'Resource Group',
		'Resource Name', 
		'Type', 
		'Product', 
		'Consumption', 
		'Savings', 
		'Actions'
	]
	
	const filterData = () => {
		let filteredData = []
		for(var subscription of data) {
			for(var resourceGroup of subscription.resourceGroups){
				for(var resource of resourceGroup.resources){
					let Arr = Object.values(resource)
					for(var x of Arr) {
						if(!Array.isArray(x)) {
							let temp
							if(typeof(x) === 'number') {
								temp = x.toString()
							} else {
								temp = x
							}
							if(temp.includes(filter)){
								filteredData.push(Object.assign({}, 
									{groupName: resourceGroup.groupName},
									resource
								)) 
								break
							}
						}
					}
				}
			}
		}
		console.log(filteredData)
	}

	filterData()
	return (
		<Paper elevation={3}>
			<Table>
				<TableHead>
					<TableRow>
						{headerTexts.map((text,idx) => (
							<TableCell key={idx}>
								<Typography style={{fontWeight:"bold"}} variant="subtitle2">
									{text}
								</Typography>
							</TableCell>
						))}
					</TableRow>
				</TableHead>
			</Table>
		</Paper>
	)
}

export default FilteredTable