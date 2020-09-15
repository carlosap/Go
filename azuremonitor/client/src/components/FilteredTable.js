import React, {useState} from 'react'
import {Table, TableHead, TableBody, TableRow, TableCell} from '@material-ui/core'
import {Box, IconButton,  Paper, Typography} from '@material-ui/core'
import _ from 'lodash'

// Import Components
import Recommendations from './Recommendations'
import Edit from "./Edit"
import Usages from "./Usages"
import MoreVert from "./MoreVert"

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
	const styles = {
		noPaddingVert: {
			paddingBottom: 0,
			paddingTop: 0
		}
	}

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
							if(temp.toLowerCase().includes(filter.toLowerCase())){
								filteredData.push(Object.assign({}, 
									{subscriptionName: subscription.subscriptionName},
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
		return filteredData
	}

	return (
		<Paper elevation={3}>
			<Table>
				<TableHead>
					<TableRow>
						{headerTexts.map((text,idx) => (
							<TableCell align='center' key={idx}>
								<Typography style={{fontWeight:"bold"}} variant="subtitle2">
									{text}
								</Typography>
							</TableCell>
						))}
					</TableRow>
				</TableHead>
				<TableBody>
					{filterData().map((data,idx) => (
						<TableRow key={idx}>
							<TableCell align='center' style={styles.noPaddingVert}>{data.groupName}</TableCell>
							<TableCell align='center' style={styles.noPaddingVert}>{data.resourceName}</TableCell>
							<TableCell align='center' style={styles.noPaddingVert}>{data.type}</TableCell>
							<TableCell align='center' style={styles.noPaddingVert}>{data.product ? data.product : "-"}</TableCell>
							<TableCell align='center' style={styles.noPaddingVert}>{data.consumption}</TableCell>
							<TableCell align='center' style={styles.noPaddingVert}>{data.savings}</TableCell>
							<TableCell style={styles.noPaddingVert}>
								<div style={{display:'flex', justifyContent:'center'}}>
									<Usages resource={data}/>
									<Recommendations recommendations={data.recommendations}/>
									<Edit groupName={data.groupName} subscription={data.subscriptionName} resource={data}/>
									<MoreVert subscription={data.subscriptionName} groupName={data.groupName} resourceName={data.resourceName}/>
								</div>
							</TableCell>
						</TableRow>
					))}
				</TableBody>
			</Table>
		</Paper>
	)
}

export default FilteredTable