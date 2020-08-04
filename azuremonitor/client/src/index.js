import './index.css';
import React from 'react';
import ReactDOM from 'react-dom';
import {  Route, Switch, BrowserRouter } from "react-router-dom";
import App from './layouts/App';


ReactDOM.render(
  <BrowserRouter>
      <App/>
  </BrowserRouter>,
  document.getElementById('root')
);


