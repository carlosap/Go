import React, {useState} from 'react'
import {withStyles} from '@material-ui/core/styles'
import {Table, TableHead, TableBody, TableRow, TableCell} from '@material-ui/core'
import {Box, IconButton, Collapse, Paper, Typography} from '@material-ui/core'
import Recommendations from './Recommendations'
import KeyboardArrowDownIcon from '@material-ui/icons/KeyboardArrowDown';
import KeyboardArrowUpIcon from '@material-ui/icons/KeyboardArrowUp';
import SubdirectoryArrowRightIcon from '@material-ui/icons/SubdirectoryArrowRight';

const BorderlessCell = withStyles({
	root: {
		borderBottom: 'none'
	}
})(TableCell)

// Component for row with collapse
const CollapseRow = (props) => {
	const {name, open, setOpen, styles} = props
	return (
		<TableRow>
			<BorderlessCell style={styles}>
					{name}
					<IconButton size='small' onClick={() => setOpen(!open)}>
						{open ? <KeyboardArrowUpIcon/> : <KeyboardArrowDownIcon/> }
					</IconButton>
			</BorderlessCell>
		</TableRow>
	)
}

// Wrapper Component for anything within resource groups
const ResourceGroups = (props) => {
	const {group} = props
	const [open, setOpen] = useState(false)
	const styles = {
		padding: {
			paddingBottom: '8px', 
			paddingTop:'8px'
		},
		collapseRow: {
			paddingBottom: '8px', 
			paddingTop:'8px',
			paddingLeft:'36px'
		}
	}
	const headerTexts = ['Resource Name', 'Type', 'Consumption', 'Usage', 'Savings', 'Recommendations']
	
	return (
		<React.Fragment>
			<CollapseRow style={{paddingLeft:'36px'}}name={group.groupName} open={open} setOpen={setOpen} styles={styles.collapseRow}/>
			<TableRow>
				<TableCell style={{ paddingBottom: 0, paddingTop: 0, paddingLeft:'36px' }}>
					<Collapse in={open}> 
						<Box>
							{group.resources.length > 0 ? 
								<Table>
									<TableHead>
											<TableRow>
												{headerTexts.map((text, idx) => {
													// Add Icon and make make flex display for first header column
													if(idx === 0) { 
														return (
															<TableCell key={text} style={{ display:'flex', paddingBottom: '8px', paddingTop: 0}}>
																<SubdirectoryArrowRightIcon fontSize="small"/>
																<Typography style={{fontWeight:"bold"}} variant="subtitle2">
																	{text}
																</Typography>
															</TableCell>
														)
													} else return (
														<TableCell align='center' key={text} style={{paddingBottom: '8px', paddingTop: 0}}>
															<Typography style={{fontWeight:"bold"}} variant="subtitle2">
																{text}
															</Typography>
														</TableCell>
													)
												})}
											</TableRow>
									</TableHead>
									<TableBody>
											{group.resources.map((resource, idx) => {
												if(idx === group.resources.length - 1){
													return (
															<TableRow key={idx}>
																<BorderlessCell style={styles.padding}>{resource.resourceName}</BorderlessCell>
																<BorderlessCell align='center' style={styles.padding}>{resource.type}</BorderlessCell>
																<BorderlessCell align='center' style={styles.padding}>{resource.consumption}</BorderlessCell>
																<BorderlessCell align='center' style={styles.padding}>${resource.usage}</BorderlessCell>
																<BorderlessCell align='center' style={styles.padding}>${resource.savings}</BorderlessCell>
																<BorderlessCell align='center' style={styles.padding}>
																	<Recommendations recommendations={resource.recommendations}/>
																</BorderlessCell>
															</TableRow>
													)
												} else return (
													<TableRow key={idx}>
															<TableCell style={styles.padding}>{resource.resourceName}</TableCell>
															<TableCell align='center' style={styles.padding}>{resource.type}</TableCell>
															<TableCell align='center' style={styles.padding}>{resource.consumption}</TableCell>
															<TableCell align='center' style={styles.padding}>${resource.usage}</TableCell>
															<TableCell align='center' style={styles.padding}>${resource.savings}</TableCell>
															<TableCell align='center' style={styles.padding}>
																<Recommendations recommendations={resource.recommendations}/>
															</TableCell>
													</TableRow>
												)
											})}
									</TableBody>
								</Table> : 
								<Typography> No Resources Found</Typography>
							}
						</Box>
					</Collapse>
				</TableCell>
			</TableRow>
		</React.Fragment>
	)
}

// Wrapper Component for anything within subscriptions
const Subscriptions = (props) => {
	const {subscription} = props
	const [open, setOpen] = useState(false)
	const styles = {paddingBottom: '8px', paddingTop: '8px'}
	return (
		<React.Fragment>
			<CollapseRow name={subscription.subscriptionName} open={open} setOpen={setOpen} styles={styles}/>
			<TableRow>
					<TableCell style={{ paddingBottom: 0, paddingTop: 0 }}>
						<Collapse in={open}> 
							<Box>
								<Table>
									<TableHead>
										<TableRow>
												<TableCell style={{display:'flex', paddingBottom: '8px', paddingTop:'0'}}>
													<SubdirectoryArrowRightIcon fontSize="small"/>
													<Typography style={{fontWeight:"bold"}} variant="subtitle2">
														Resource Groups
													</Typography>
												</TableCell>
										</TableRow>
									</TableHead>
									<TableBody>
										{subscription.resourceGroups.map((group,i) => (
											<ResourceGroups key={i} group={group}/>
										))}
									</TableBody>
								</Table>
							</Box>
						</Collapse>
					</TableCell>
			</TableRow>
		</React.Fragment>
	)
}

const ResourcesTable = (props) => {
	const {data} = props
	return (
		<Paper elevation={3}>
			<Table>
				<TableHead>
					<TableRow>
						<TableCell> 
								<Typography style={{fontWeight:"bold"}} variant="h6">
									Subscriptions
								</Typography> 
						</TableCell>
					</TableRow>
				</TableHead>
				<TableBody>
					{data.map((subscription, idx) => (
						<Subscriptions 
							key={subscription.subscriptionName} 
							subscription={subscription}
						/>
					))}
				</TableBody>
			</Table>
		</Paper>
   )
}

export default ResourcesTable