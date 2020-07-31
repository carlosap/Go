import React from 'react'
import {makeStyles} from '@material-ui/core/styles'
import {Divider, Typography, Paper} from '@material-ui/core';

const useStyles = makeStyles({
   amountText: {

   }
})

const AnalyticsWidget = (props) => {
    const styles = useStyles()

    return (
        <Paper elevation={3}>
            <div style={{padding:'10px'}}>
                <Typography style={{marginBottom: '10px'}} color='primary' variant='h6'>
                    {props.title}
                </Typography>

                <Divider style={{marginBottom: '10px'}}/>

                <div style={{textAlign:'center'}}>  
                    <Typography color='textPrimary' variant='h4'>
                        ${props.amount}
                    </Typography>
                </div>
            </div>

            <Divider />

            <div style={{padding:'10px'}}>
                <Typography color='textSecondary' variant="subtitle2">
                    {props.description}
                </Typography>
            </div>
                
        </Paper>
    )
}

export default AnalyticsWidget