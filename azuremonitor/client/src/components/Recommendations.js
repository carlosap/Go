import React, {useState} from 'react'
import {Dialog, DialogActions, DialogTitle, DialogContent} from '@material-ui/core'
import {Divider, Tooltip, IconButton, Typography} from '@material-ui/core'
import InfoIcon from '@material-ui/icons/Info';
import CloseIcon from '@material-ui/icons/Close';

const Recommendations = (props) => {    
    const [open, setOpen] = useState(false)
    const {recommendations} = props
    return (
        <div>
            <Tooltip title="Show Recommendations">
                <IconButton onClick={() => setOpen(true)}>
                    <InfoIcon color='primary'/>
                </IconButton>
            </Tooltip>

            <Dialog maxWidth='xs' fullWidth open={open} onClose={() => setOpen(false)}>

                <DialogTitle style={{textAlign:'center', paddingBottom:'8px'}}>
                    <Typography style={{fontWeight:'bold'}} variant="h5"> Recommendations</Typography>
                    <Divider fullWidth/>
                </DialogTitle>

                <DialogContent>
                    <Typography style={{marginBottom:' 10px'}} color="primary" variant="h6">
                        {`${recommendations.length} Optimization(s) Available`}
                    </Typography>
                    {recommendations.length > 0 ? 
                        recommendations.map((rec, idx) => (
                            <Typography key={idx} style={{marginBottom:'4px'}} color="textPrimary">{`- ${rec}`}</Typography>
                        ))
                    : 
                        <Typography> This area has no optimizations available </Typography>
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

export default Recommendations