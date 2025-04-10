import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App';

const appRootContainer = ReactDOM.createRoot(document.getElementById('root'));

appRootContainer.render(
    <React.StrictMode>
        <App />
    </React.StrictMode>
);