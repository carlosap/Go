import HomePage from './pages/HomePage';
import AnalysisPage from './pages/AnalysisPage'
import UsagePage from './pages/UsagePage'
import React from 'react';
import { Redirect } from "react-router-dom";

const NavbarRoutes = {
    title: 'Azure Resource Tracker',
    routes: [
        {
            path: "/home",
            name: "Cost",
            icon: "",
            component: HomePage
        },
        {
            path: "/home",
            name: "Security",
            icon: "",
            component: HomePage
        },
        {
            path: "/home",
            name: "High Availability",
            icon: "",
            component: HomePage
        },
        {
            path: "/analysis/:type/:id",
            name: "Analysis",
            icon: "",
            hidden: true,
            component: AnalysisPage
        },
        {
            path: "/usage/:id",
            name: "Usage",
            icon: "",
            component: UsagePage
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