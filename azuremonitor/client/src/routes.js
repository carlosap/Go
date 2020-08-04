import HomePage from './pages/HomePage';
import React from 'react';
import { Redirect } from "react-router-dom";

const NavbarRoutes = {
    title: 'Azure Monitor',
    routes: [
        {
            path: "/home",
            name: "Home",
            icon: "",
            component: HomePage
        },
        {
            path: "/",
            name: "",
            icon: "",
            component: () => <Redirect to="/home"/>
        }
        
    ]
}
export default NavbarRoutes;