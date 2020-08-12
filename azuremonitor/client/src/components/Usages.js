import React from 'react'
import {IconButton, Tooltip} from '@material-ui/core'
import DataUsageIcon from '@material-ui/icons/DataUsage'

const Usages = (props) => {
    const handleClick = () => {
        console.log("ckik")
    }
    
    return (
        <Tooltip title='Usages'>
            <IconButton onClick={handleClick}>
                <DataUsageIcon color="primary"/>
            </IconButton>
        </Tooltip>
    )
}

export default Usages