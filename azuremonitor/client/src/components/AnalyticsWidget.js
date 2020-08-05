import React from 'react'
import {Divider, Typography, Paper} from '@material-ui/core';

const AnalyticsWidget = (props) => {
    const numberWithCommas = (num) => {
        return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
    }

    return (
        <Paper elevation={3}>
            <div style={{padding:'10px'}}>
                <Typography style={{marginBottom: '10px'}} color='primary' variant='h6'>
                    {props.title}
                </Typography>

                <Divider style={{marginBottom: '10px'}}/>

                <div style={{textAlign:'center'}}>  
                    <Typography style={{fontWeight:'bold'}} color='textPrimary' variant='h4'>
                        ${numberWithCommas(props.amount)}
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