import React  from 'react';
import {makeStyles} from '@material-ui/core/styles'
import {Grid, Typography} from '@material-ui/core'

import AnalyticsWidget from '../components/AnalyticsWidget'
import ResourcesTable from '../components/ResourcesTable'
import LineGraph from '../components/LineGraph'
import PieChart from '../components/PieChart'
import Logo from '../assets/images/logo.png'

import data from "../MockData/resources.json"

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

        <Grid style={{padding:'10px'}} container>
            <Grid style={{display:'flex', justifyContent:'space-between'}} sm={12} item>
                <div>
                    <img className={styles.logo} src={Logo} alt={'logo'} width={270}></img>
                </div>

                <div className={styles.widgetsContainer}>
                    <div className={styles.widget}>
                        <AnalyticsWidget 
                            title='Current Cost' 
                            amount={20000000}
                            description='This is the current cost of all of the current resources being used'
                        />
                    </div>

                    <div className={styles.widget}>
                        <AnalyticsWidget 
                            title='Savings w/ Optimization'
                            amount={4000000}
                            description='Estimated savings after optimization and modernizing'
                        />
                    </div>
                    
                    <div className={styles.widget}>
                        <AnalyticsWidget 
                            title='Cost after Optimization'
                            amount={16000000}
                            description={'Total cost with optimization and modernization'}
                        />
                    </div>
                </div>
            </Grid>

            <Grid style={{padding:'30px'}} item sm={12}>
                <ResourcesTable data={data}/>
            </Grid>

            <Grid style={{paddingLeft:'30px', paddingRight:'30px'}} container item sm={12} spacing={2}>
                <Grid item sm={6} style={{display:'flex-column'}}>
                    <div style={{display:'flex', flexDirection:'column', alignItems:'center'}}>
                        <Typography 
                            variant='h6'
                            color="primary"
                        > 
                            Area Of Savings
                        </Typography>
                        <PieChart/>
                    </div>
                </Grid>
                <Grid item sm={6}>
                    <div style={{display:'flex', flexDirection:'column', alignItems:'center'}}>
                        <Typography 
                       
                            variant='h6'
                            color="primary"
                        > 
                            Cost Projection
                        </Typography>
                        <LineGraph/>
                    </div>                
                </Grid>
            </Grid>
        </Grid>
    )
}

export default HomePage;