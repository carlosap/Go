import React, {useState} from 'react'
import {Divider, Tooltip, IconButton, Typography} from '@material-ui/core'
import {Dialog, DialogActions, DialogTitle, DialogContent} from '@material-ui/core'
import DataUsageIcon from '@material-ui/icons/DataUsage'
import CloseIcon from '@material-ui/icons/Close';

const Usages = (props) => {
    const [open, setOpen] = useState(false)

    const handleClick = () => {
        setOpen(true)
    }

    const {usage} = props

    return (
        <div>
            <Tooltip title='Usages'>
                <IconButton onClick={handleClick}>
                    <DataUsageIcon color="primary"/>
                </IconButton>
            </Tooltip>

            <Dialog maxWidth='xs' fullWidth open={open} onClose={() => setOpen(false)}>

            <DialogTitle style={{textAlign:'center', paddingBottom:'8px'}}>
                <div>
                    <Typography variant='h4' style={{fontWeight:'bold'}}> Usages </Typography>
                    <Divider/>
                </div>
            </DialogTitle>

            <DialogContent>
                <Typography style={{marginBottom:' 10px'}} color="primary" variant="h6">
                    {`${usage.length} detected`}
                </Typography>
                {usage.length > 0 ? 
                    usage.map((u, idx) => (
                        <Typography key={idx} style={{marginBottom:'4px'}} color="textPrimary">{`- ${u}`}</Typography>
                    ))
                : 
                    <Typography> No Usage Detected </Typography>
                }
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