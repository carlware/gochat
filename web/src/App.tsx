import React, { useEffect } from 'react';
import logo from './logo.svg';
import './App.css';

function App() {
    let conn: WebSocket

    useEffect(() => {
        conn = new WebSocket("ws://127.0.0.1:8080/ws");
        conn.onclose = function(evt) {
            console.log(evt)
        };
        conn.onopen = function(evt) {
            console.log(evt)
        };
        conn.onmessage = function(evt) {
            console.log(evt)
        };
    })

    const sendMessage = (message: string = "") => {
        conn.send(message)
    }


    return (
        <div className="App">
            <header className="App-header">
                <img src={logo} className="App-logo" alt="logo" />
                <p>
                    Edit <code>src/App.tsx</code> and save to reload.
                </p>
                <a
                    className="App-link"
                    href="https://reactjs.org"
                    target="_blank"
                    rel="noopener noreferrer"
                >
                    Learn React
                </a>
                <button onClick={() => sendMessage("ping")}>Send</button>
            </header>
        </div>
    );
}

export default App;
