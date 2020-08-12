import React, {useState, useContext} from 'react'
import { AppContext } from '../contexts/AppContext'
import {useParams} from 'react-router-dom'
import {makeStyles, withStyles} from '@material-ui/core/styles'
import {Tab, Tabs, Typography, Paper} from '@material-ui/core'
import {VMQuestions} from "../MockData/questions.json"
import VMQuestionaire from '../components/AnalysisForms/VMQuestionaire'

const StyledTabs = withStyles({
    indicator: {
        display: 'flex',
        justifyContent: 'center',
        backgroundColor: 'transparent',

        // define and inject spans to further customize indicator
        '& > span': {   
            maxWidth: '60%',
            width:'100%',
            backgroundColor: '#87CEFA'
        }
    }
})((props) => <Tabs {...props} TabIndicatorProps={{children: <span/>}}/>);

const StyledTab = withStyles((theme) => ({
    root: {
        textTransform:'none',
        fontSize: theme.typography.pxToRem(16),
        fontWeight: theme.typography.fontWeightBold,
        color: 'black',
        paddingLeft: '0',
        paddingRight: '0'
    }
}))((props) => <Tab {...props}/>)

const useStyles = makeStyles({
    root: {
        width: '70%',
        display: 'flex',
        flexDirection: 'column',
        margin: '20px auto',
    },
    questions: {
        display: 'flex',
        flexDirection: 'column',
        alignItems:'center',
        maxHeight: 550,
        overflow:'auto'
    }
}) 

const AnalysisPage = () => {
    const { state, dispatch } = useContext(AppContext)
    const styles = useStyles()
    const {type, id} = useParams()

    // Used to keep track of tabs
    const [tabIndex, setIndex] = useState(0)

    const handleTabSwitch = (e, index) => {
        setIndex(index)
    }

    return (
        <div className={styles.root}>
            <div>
                <Typography style={{fontWeight: 'bold'}} variant='h4' color='textPrimary'>
                    {id} - {type}
                </Typography>

                <Typography style={{marginTop:'10px'}} variant='body1'>
                    Collect, analyze, and act on telemetry data from your Azure Subscription. 
                    Capgimin Monitor helps you optimize performance and availability of your applications 
                    and proactively identify problems in seconds.
                </Typography>
            </div>

            <Paper style={{margin: '35px auto', height:600, paddingBottom: '15px'}} elevation={3}>
                <div>
                    <StyledTabs value={tabIndex} onChange={handleTabSwitch}>
                        {VMQuestions.map((question, idx) => (
                            <StyledTab id={question.category} key={idx} label={question.category}/>
                        ))}
                    </StyledTabs>
                </div>

                <div className={styles.questions}>
                    <VMQuestionaire resourceName={id} state={state} dispatch={dispatch} questions={state.Questions[tabIndex]}/>
                </div>

            </Paper>
                
           
        </div>
    )
}

export default AnalysisPage