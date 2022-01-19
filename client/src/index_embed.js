import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import App from './App';
import reportWebVitals from './reportWebVitals';

const APP_ID = 'surfboard-embed';
const RUNNER_KEY = `${APP_ID}_runner`;
const DESTROYER_KEY = `${APP_ID}_destroyer`;

window[RUNNER_KEY] = (opts, targetElementId) => {
  ReactDOM.render(
    <React.StrictMode>
      <App {...opts} />
    </React.StrictMode>,
    document.getElementById(targetElementId)
  );
};

window[DESTROYER_KEY] = (opts, targetElementId) => {
  React.DOM.unmountComponentAtNode(document.getElementById(targetElementId));
};

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
