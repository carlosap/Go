import React from 'react';
import NavbarRoutes from '../routes';
import { makeStyles } from '@material-ui/core/styles';
import {Link} from "react-router-dom"
import {AppBar, Toolbar, Typography} from '@material-ui/core'
import SearchBar from './SearchBar'

const useStyles = makeStyles((theme) => ({
    toolBarContainer: {
        display:'flex',
        justifyContent: 'space-between'
    },
    SearchBarContainer: {
        width: 600
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
            <Toolbar className={styles.toolBarContainer}>
                <Typography color="inherit" variant="h5">
                    {NavbarRoutes.title}
                </Typography>
                
                <div className={styles.SearchBarContainer}>
                    <SearchBar/>
                </div>

                <div className={styles.linksContainer}>
                    {NavbarRoutes.routes.map((menuItem, index) => (
                        !menuItem.hidden && (
                            <Typography key={index} variant="subtitle1" >
                                <Link to={menuItem.path} className={styles.menuItems}>
                                    {menuItem.name}
                                </Link>
                            </Typography>
                        )
                    ))}
                </div>
                
            </Toolbar>
        </AppBar>
    )
}

export default Navbar;