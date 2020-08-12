import React, {useState} from 'react'
import {Dialog, DialogActions, DialogTitle, DialogContent} from '@material-ui/core'
import {Divider, Tooltip, IconButton, Typography} from '@material-ui/core'
import InfoIcon from '@material-ui/icons/Info';
import CloseIcon from '@material-ui/icons/Close';

const Recommendations = (props) => {    
    const [open, setOpen] = useState(false)
    const {recommendations} = props

    // Check if any questions were answered in this category
    const checkIfAnswered = (questions) => {
        for(var question of questions) {
            if(question.answer !== ''){
                return true
            }
        }
        return false
    }

    return (
        <div>
            <Tooltip title="Show Recommendations">
                <IconButton onClick={() => setOpen(true)}>
                    <InfoIcon color='primary'/>
                </IconButton>
            </Tooltip>

            <Dialog scroll='paper' maxWidth='md' fullWidth open={open} onClose={() => setOpen(false)}>

                <DialogTitle style={{textAlign:'center', paddingBottom:'8px'}}>
                    <div>
                        <Typography style={{fontWeight:'bold'}} variant="h5"> Recommendations</Typography>
                        <Divider/>
                    </div>
                </DialogTitle>

                <DialogContent>
                    {recommendations.length > 0 ? 
                        recommendations.map((rec, idx) => (
                            <div key={idx} style={{textAlign:'center'}}>
                                {checkIfAnswered(rec.questions) &&  (
                                        <Typography style={{marginBottom:' 15px', marginTop: '15px', fontWeight:'bold'}} color="primary" variant="h6">
                                            {`${rec.category}`}
                                        </Typography>
                                    )
                                }
                                
                                {rec.questions.map((question) => {
                                    if(question.answer !== "") {
                                        return (
                                            <div key={idx} >
                                                <Typography variant='h6' color="textPrimary">{question.question}</Typography>
                                                <Typography style={{marginBottom:'8px'}} color="textPrimary">{`- ${question.answer}`}</Typography>
                                            </div>
                                        )
                                    }
                                })}
                            </div>
                           
                        ))
                    : 
                        <Typography> Analysis form has not been filled out yet. You can find this form by clicking on the "Analysis Form" in the actions section</Typography>
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