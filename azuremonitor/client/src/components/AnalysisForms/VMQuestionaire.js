import React, {useState} from 'react'
import {TextField, Typography} from '@material-ui/core'

const VMQuestionaire = (props) => {
    const {questions, dispatch} = props

    const handleChange = (event) => {
        let payload = {
            id: parseInt(event.target.id),
            answer: event.target.value,
            category: questions.category
        }
        dispatch({type:"UPDATE_ANSWER", payload: payload})
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
        </div>
    )
} 

export default VMQuestionaire