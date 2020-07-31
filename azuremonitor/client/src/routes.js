import HomePage from './pages/HomePage';
// import LocationPage from './pages/LocationPage';
// import NewsPage from './pages/NewsPage';
// import WeatherPage from './pages/WeatherPage';

const NavbarRoutes = {
    title: 'Azure Monitor',
    routes: [
        {
            path: "/home",
            name: "Home",
            icon: "",
            component: HomePage
        },
        // {
        //     path: "/location",
        //     name: "Location",
        //     icon: "",
        //     component: LocationPage
        // },
        // {
        //     path: "/weather",
        //     name: "Weather",
        //     icon: "",
        //     component: WeatherPage
        // },
        // {
        //     path: "/news",
        //     name: "News",
        //     icon: "",
        //     component: NewsPage
        // },
    ]
}
export default NavbarRoutes;