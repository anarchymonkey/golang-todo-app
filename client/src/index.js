import React from 'react';
import { createRoot } from 'react-dom/client';
import App from './App.jsx'
import './index.css';

const renderRoot = () => {
    const reactRoot = createRoot(document.getElementById("todo-app-root"))
    reactRoot.render(<App />)
}

renderRoot()