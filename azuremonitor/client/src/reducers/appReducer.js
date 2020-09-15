import {VMQuestions} from '../MockData/questions.json'
import { FormHelperText } from '@material-ui/core'

const appReducer = (state, action) => {
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

    case 'SAVE_QUESTIONS': 
      const { subscriptionName, recommendations, resourceName, groupName } = action.payload
      return {
        ...state,
        Resources: state.Resources.map((subscription) => {
          if(subscription.subscriptionName === subscriptionName) {
            let temp = subscription
            // Find Resource Group
            for(let i = 0; i < temp.resourceGroups.length; i++) { 
              if(temp.resourceGroups[i].groupName === groupName) {
                  // Find Resource within group and update recommendations
                for(let j = 0; j < temp.resourceGroups[i].resources.length; j++) {
                  if(temp.resourceGroups[i].resources[j].resourceName === resourceName) {
                    temp.resourceGroups[i].resources[j].recommendations = recommendations
                    
                    // Check if savings was entered. If so, inject into resource
                    let summary = recommendations[recommendations.length - 1]
                    let estimate = summary.questions[summary.questions.length - 1].answer
                    if(estimate !== '') {
                      temp.resourceGroups[i].resources[j].savings = estimate
                    }
                    return temp
                  }
                }
              }
            }
          } 
          return subscription
        })
      }

    case 'RESET_QUESTIONS':
      return {
        ...state,
        Questions: state.Questions.map((question) => {
          var temp = {}
          temp.category = question.category
          temp.questions = question.questions.map(q => {
            return {
              id: q.id,
              question: q.question,
              answer: ''
            }
          })
          return temp
        })
      }
      
    case 'SET_SAVINGS': 
      const {subscription, group, resource, amount} = action.payload
      return {
        ...state,
        Resources: state.Resources.map(sub => {
          if(sub.subscriptionName === subscription){
            let temp = sub
            for(var i = 0; i < temp.resourceGroups.length; i++){
              if(temp.resourceGroups[i].groupName === group) {
                for(let j = 0; j < temp.resourceGroups[i].resources.length; j++) {
                  if(temp.resourceGroups[i].resources[j].resourceName === resource) {
                    if(amount === '-'){
                      temp.resourceGroups[i].resources[j].savings = amount
                    } else {
                      temp.resourceGroups[i].resources[j].savings = `$${amount}`
                    }
                    return temp
                  }
                }
              }
            }
          } else {
            return sub
          }
        })
      }

    case 'ADD_TO_TABLE_STATE':
      return {
        ...state,
        tableState: [...state.tableState, action.payload] 
      }
    
    case 'REMOVE_FROM_TABLE_STATE':
      return {
        ...state,
        tableState: state.tableState.filter(name => name !== action.payload)
      }

    case 'UPDATE_SEARCH_FILTER' :
      return {
        ...state,
        searchFilter: action.payload
      }
    default:
      return state;
  }
}

export default appReducer