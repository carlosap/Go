import React, {useState} from 'react'
import {useParams} from 'react-router-dom'
import {makeStyles, withStyles} from '@material-ui/core/styles'
import {Tab, Tabs, Typography} from '@material-ui/core'

const StyledTabs = withStyles({
    indicator: {
        display: 'flex',
        justifyContent: 'center',
        backgroundColor: 'transparent',

        // define and inject spans to further customize indicator
        '& > span': {   
            maxWidth: 70,
            width:'100%',
            backgroundColor: 'blue'
        }
    }
})((props) => <Tabs {...props} TabIndicatorProps={{children: <span/>}}/>);

const StyledTab = withStyles((theme) => ({
    root: {
        fontSize: theme.typography.pxToRem(15),
        fontWeight: theme.typography.fontWeightBold,
        
    }
}))((props) => <Tab {...props}/>)

const useStyles = makeStyles({
    root: {
        width: '70%',
        display: 'flex',
        flexDirection: 'column',
        margin: '20px auto'
    }
}) 

const AnalysisPage = () => {
    const styles = useStyles()
    const {type, id} = useParams()
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
                    Lorem ipsum dolor sit amet, consectetur adipiscing elit,
                    sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. 
                    Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi
                    ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit 
                    in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur 
                    sint occaecat cupidatat non proident, 
                </Typography>
            </div>

            <div style={{backgroundColor:'#696969', margin: '15px auto'}}>
                <StyledTabs value={tabIndex} onChange={handleTabSwitch}>
                    <StyledTab label="Computation"/>
                    <StyledTab label="Storage"/>
                    <StyledTab label="Network Data"/>
                    <StyledTab label="Migration and Agreements"/>
                    <Tab label="Optimizations"/>
                    <Tab label="Accountability"/>
                    <Tab label="Summary"/>
                </StyledTabs>
            </div>
        </div>
    )
}

export default AnalysisPage