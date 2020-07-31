import React  from 'react';
import {makeStyles} from '@material-ui/core/styles'
import {Grid} from "@material-ui/core"

import AnalyticsWidget from '../components/AnalyticsWidget'
import Logo from "../assets/images/Honeywell-Logo-17.jpg"

const useStyles = makeStyles({
    logo: {
        maxWidth: "100%",
        height:"auto"
    }
})

const HomePage = () => {

    const styles = useStyles()
    return (

        <Grid style={{padding:'10px'}}container>
            <Grid sm={12} item>
                <div>
                    <img className={styles.logo} src={Logo} alt={"logo"} width={270}></img>
                </div>

                <div>

                </div>
            </Grid>
        </Grid>
    )
}

export default HomePage;