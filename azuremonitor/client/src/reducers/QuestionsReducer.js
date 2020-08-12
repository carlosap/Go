
const QuestionsReducer = (state, action) => {
    switch (action.type) {
        case 'UPDATE_ANSWER':
            const {id, answer, category} = action.payload
            return {
                ...state,
                Questions: state.Questions.map((question) => {
                    if(question.category === category){
                        // Find question and update answer
                        let temp = question 
                        for(var i = 0; i < temp.questions.length; i++){
                            if(temp.questions[i].id === id ){
                                temp.questions[i].answer = answer
                                return temp
                            }
                        }
                    } 
                    return question
                })
            }
        default:
            return state;
    }
}

export default QuestionsReducer