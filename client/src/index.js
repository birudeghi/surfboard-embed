import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import App from './App';
import reportWebVitals from './reportWebVitals';

import { Client as Styletron } from 'styletron-engine-atomic';
import { Provider as StyletronProvider } from "styletron-react";
import { BaseProvider, createTheme } from 'baseui';

const engine = new Styletron();

const primitives = {
    primaryFontFamily: 'Work Sans',
};

const theme = createTheme(primitives);

const widgetDivs = document.querySelectorAll('.surfboard-app-board');

widgetDivs.forEach(div => {
  ReactDOM.render(
    <StyletronProvider value={engine}>
        <BaseProvider theme={theme}>
            <React.StrictMode>
                <App appId={div.dataset.appid}/>
            </React.StrictMode>
        </BaseProvider>
    </StyletronProvider>,
    div
  );
  });

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
