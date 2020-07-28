import React from 'react';
import NavbarRoutes from '../routes';
import { makeStyles } from '@material-ui/core/styles';
import {Link} from "react-router-dom"
import {AppBar, Toolbar, Typography} from '@material-ui/core'

const useStyles = makeStyles((theme) => ({
    title:{
        flexGrow: 1
    },
    linksContainer: {
        display: "flex",
    },
    menuItems: {
        margin: theme.spacing(1),
        color:"white"
    }
}))

const Navbar = () => {
    const styles = useStyles()

    return (
        <AppBar position="static">
            <Toolbar>
                <Typography color="inherit" variant="h5" className={styles.title}>
                    {NavbarRoutes.title}
                </Typography>
                
                <div className={styles.linksContainer}>
                    {NavbarRoutes.routes.map((menuItem, index) => (
                        <Typography key={index} variant="subtitle1" >
                            <Link to={menuItem.path} className={styles.menuItems}>
                                {menuItem.name}
                            </Link>
                        </Typography>
                    ))}
                </div>
                
            </Toolbar>
        </AppBar>
    )
}

export default Navbar;