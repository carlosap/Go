import React, {useState} from 'react'
import {makeStyles} from '@material-ui/core/styles'
import {Grid, Divider, Tooltip, IconButton, Typography} from '@material-ui/core'
import {Dialog, DialogActions, DialogTitle, DialogContent} from '@material-ui/core'
import DataUsageIcon from '@material-ui/icons/DataUsage'
import CloseIcon from '@material-ui/icons/Close';
import ComputerIcon from '@material-ui/icons/Computer';
import InfoIcon from '@material-ui/icons/Info';

const useStyles = makeStyles({
	titleContainer: {
		textAlign:'center', 
		paddingBottom: 8
	},
	usageContainer: {
		display: 'flex',
		overflow: 'hidden'
	},
	infoHeaderContainer:{
		display: 'flex',
		alignItems: 'center',
	},

})

const InfoHeader = (props) => {
	const styles = useStyles()
	const {icon, title} = props 

	return (
		<div className={styles.infoHeaderContainer}>
			{icon}
			<Typography style={{fontWeight:'bold'}} variant='subtitle1'>
				{title}
			</Typography>
		</div>
	)
}

const InfoLine = (props) => {
	const {name, value} = props
	return (
		<div style={{display:'flex', marginTop: 5}}>
			<div style={{width:'50%'}}>
				<Typography 
					color='textSecondary'
				> 
					{name} 
				</Typography>
			</div>
			<div style={{width: '50%'}}>
				<Typography> {value} </Typography>
			</div>
		</div>
	)
}

const Usages = (props) => {
	const styles = useStyles()
	const [open, setOpen] = useState(false)

	const handleClick = () => {
			setOpen(true)
	}

	const {resource} = props

	return (
		<div>
			<Tooltip title='Usages'>
				<IconButton onClick={handleClick}>
					<DataUsageIcon color="primary"/>
				</IconButton>
			</Tooltip>

			<Dialog maxWidth='md' fullWidth open={open} onClose={() => setOpen(false)}>

			<DialogTitle className={styles.titleContainer}>
					<Typography style={{fontWeight:'bold'}} variant='h4'> Usages </Typography>
					<Divider/>
			</DialogTitle>

			<DialogContent className={styles.usageContainer}>
				<div style={{width: '50%'}}>
					<InfoHeader 
						title={resource.type} 
						icon={<ComputerIcon color="primary" style={{marginRight: 5}}/>}
					/>

					{resource.resourceInfo.map((info) => (
						<InfoLine name={info.displayName} value={info.displayValue}/>
					))}
				</div>

				<div style={{width: '50%'}}>
					<InfoHeader 
						title="Usage" 
						icon={<InfoIcon color="primary" style={{marginRight: 5}}/>}
					/>

					{resource.resourceUsage.map((info) => (
						<InfoLine name={info.displayName} value={info.displayValue}/>
					))}
				</div>
			</DialogContent>

			<DialogActions style={{alignSelf:'center'}}>
				<Tooltip title="Close">
					<IconButton onClick={() => setOpen(false)}>
						<CloseIcon/>
					</IconButton>
				</Tooltip>
			</DialogActions>

			</Dialog>
		</div>
	)
}

export default Usages