import React from 'react';
import { Switch, Route } from "react-router-dom";
import AppContextProvider from '../contexts/AppContext';
import Navbar from '../components/Navbar';
import NavbarRoutes from '../routes';


const App = ({props}) => {

  const getRoutes = (routes) => {
    return NavbarRoutes.routes.map((prop, key) => {
        return (
          <Route
            path={prop.path}
            key={key}
            component={prop.component}
          />
        );
    });
  }

  return (
    <div className="App">
      <AppContextProvider>
        <Navbar />
        <Switch>{getRoutes(NavbarRoutes)}</Switch>
      </AppContextProvider>
    </div>
  );
}

export default App;
