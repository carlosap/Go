import React  from 'react';
import {makeStyles} from '@material-ui/core/styles'
import {Grid} from "@material-ui/core"

import AnalyticsWidget from '../components/AnalyticsWidget'
import Logo from "../assets/images/Honeywell-Logo-17.jpg"

const useStyles = makeStyles({
    logo: {
        maxWidth: "100%",
        height:"auto"
    },
   
    widgetsContainer: {
        marginRight:"75px",
        display: 'flex',
        width: '45%',
        justifyContent: 'space-between'
    },
    widget: {
        width:"100%",
        padding: '0 10px'
    }
})

const HomePage = () => {
    const styles = useStyles()

    return (

        <Grid style={{padding:'10px'}}container>
            <Grid style={{display:'flex'}} sm={12} item justify="space-between">
                <div>
                    <img className={styles.logo} src={Logo} alt={"logo"} width={270}></img>
                </div>

                <div className={styles.widgetsContainer}>
                    <div className={styles.widget}>
                        <AnalyticsWidget 
                            title="Current Cost" 
                            amount={100000}
                            description="This is the current cost of all of the current resources being used"
                        />
                    </div>
                    <div className={styles.widget}>
                        <AnalyticsWidget 
                            title="Savings From Optimization" 
                            amount={10000000}
                            description="This is the current cost of all of the current resources being used"
                        />
                    </div>
                    <div className={styles.widget}>
                        <AnalyticsWidget 
                            title="Cost after Optimization" 
                            amount={10000000}
                            description="This is the current cost of all of the current resources being used"
                        />
                    </div>
                </div>
            </Grid>
        </Grid>
    )
}

export default HomePage;