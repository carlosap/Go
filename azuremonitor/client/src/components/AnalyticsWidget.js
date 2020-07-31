import React from 'react'
import {makeStyles} from '@material-ui/core/styles'
import {Divider, Typography, Paper} from '@material-ui/core';

const useStyles = makeStyles({
    title:{

    },
})

const AnalyticsWidget = (props) => {
    const styles = useStyles()

    return (
        <Paper style={{padding:'10px'}} elevation={3}>
            <div>
                <Typography color='primary' variant='h6'>
                    {props.title}
                </Typography>
            </div>
            <Divider style={{marginBottom: '10px'}}/>
            <div>
                <Typography>
                    ${props.amount}
                </Typography>
            </div>

            <div style={{marginBottom: '10px', paddingLeft: '-10px'}} >
                <Divider/>
            </div>

            <div>
                <Typography>
                    {props.description}
                </Typography>
            </div>
                
        </Paper>
    )
}

export default AnalyticsWidget