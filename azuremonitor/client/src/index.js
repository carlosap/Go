import './index.css';
import React from 'react';
import ReactDOM from 'react-dom';
import { Redirect, Route, Switch, BrowserRouter } from "react-router-dom";
import App from './layouts/App';


ReactDOM.render(
  <BrowserRouter>
      <Switch>
        <Route path="/" render={props => <App {...props} />} />
        <Redirect from="/" to="/home" />
      </Switch>
  </BrowserRouter>,
  document.getElementById('root')
);


