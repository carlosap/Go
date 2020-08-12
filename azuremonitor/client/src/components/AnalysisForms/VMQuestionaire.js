import React, {useState} from 'react'
import {withRouter} from 'react-router'
import {IconButton, TextField, Typography, Tooltip} from '@material-ui/core'
import SaveIcon from '@material-ui/icons/Save'
import ArrowBackIcon from '@material-ui/icons/ArrowBack'
import Swal from 'sweetalert2'

const VMQuestionaire = (props) => {
    const {resourceName, questions, state, dispatch} = props

    const handleChange = (event) => {
        let payload = {
            id: parseInt(event.target.id),
            answer: event.target.value,
            category: questions.category
        }
        dispatch({type:'UPDATE_ANSWER', payload: payload})
    }

    const handleBack = () => {
        Swal.fire({
            title: 'Are you sure?',
            text: "If you leave before saving, all changes will be lost.",
            icon: 'warning',
            showCancelButton: true,
            confirmButtonColor: '#3085d6',
            cancelButtonColor: '#d33',
            confirmButtonText: 'Yes, go back!'
        }).then((result) => {
            if (result.value) {
                dispatch({type:'RESET_STATE'})
                props.history.push("/")
            }
        })
    }

    const handleSave = () => {
        Swal.fire({
            title: 'Save all answers?',
            icon: 'warning',
            showCancelButton: true,
            confirmButtonColor: '#3085d6',
            cancelButtonColor: '#d33',
            confirmButtonText: 'Yes, save my answers!'
        }).then((result) => {
            if (result.value) {
                let payload = {
                    subscription: props.location.state.subscription,
                    recommendations: state.Questions,
                    resourceName: resourceName
                }

                dispatch({type: 'SAVE_QUESTIONS', payload: payload})

                Swal.fire(
                'Saved!',
                'Answers have been updated on home page.',
                'success'
                ).then(() => {
                    props.history.push("/")
                })
            }
        })
    }
    
    return (
        <div>
            {questions.questions.map((question,idx) => (
                <div key={idx} style={{marginTop:'20px'}}>
                    <Typography variant='h6' color='primary'>
                        Question {idx + 1}: {question.question}
                    </Typography>

                    <TextField
                        id={question.id.toString()}
                        value={question.answer}
                        onChange={handleChange}
                        style={{width:'100%'}}
                        label="Answer"
                        multiline
                        rows={4}
                    />
                </div>
            ))}

            <div style={{marginTop:'15px', display:'flex', justifyContent: 'center'}}>
                <Tooltip title='Back To Home Page'>
                    <IconButton onClick={handleBack} color='primary'>
                        <ArrowBackIcon fontSize='large'/>
                    </IconButton>
                </Tooltip>

                <Tooltip title='Save Answers'>
                    <IconButton onClick={handleSave} color='primary'>
                        <SaveIcon fontSize='large'/>
                    </IconButton>
                </Tooltip>
            </div>
            
        </div>
    )
} 

export default withRouter(VMQuestionaire)