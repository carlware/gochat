import React, { useEffect } from 'react';
import logo from './logo.svg';
import './App.css';


interface Error {
    Code: string
    Message: string
}

interface MessageResponse {
    rid: string
    uid: string
    type: string
    message: string
    error: Error
}

interface MessageRequest {
    rid: string
    uid: string
    type: string
    message: string
}

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

    const sendMessage = (msg: string = "") => {
        const message: MessageRequest = {
            rid: "general",
            uid: "1",
            type: "message",
            message: msg,
        }
        conn.send(JSON.stringify(message))
    }
    const sendCommand = (cmd: string = "") => {
        const message: MessageRequest = {
            rid: "general",
            uid: "1",
            type: "command",
            message: "stock=" + cmd,
        }
        conn.send(JSON.stringify(message))
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
                <button onClick={() => sendCommand("AAPL.US")}>Command</button>
            </header>
        </div>
    );
}

export default App;
