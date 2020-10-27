import React, { useEffect, useState } from 'react';
import {Flex, Avatar, Button} from '@chakra-ui/core';
import './App.css';

const callAPI = (url: string, options: any): Promise<any> => {
  return new Promise(async (resolve, reject) => {
    fetch(url, options)
      .then(async (response) => {
        const text = await response.text()
        if (!response.ok) {
          reject(JSON.parse(text));
        }
        return text;
      })
      .then(response => resolve(JSON.parse(response)))
      .catch(error => reject(error));
  });
}

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

interface Message {
   id: string
   rid: string
   uid: string
   created: string
   message: string
}

interface Room {
  id: string
  name: string
}


interface RoomListProps {
  messages: Message[]
}

function RoomList({ messages=[] }: RoomListProps) {
  return(
    <span>room list</span>
  )
}

function Room() {
  
  return(
  <>
    <span>room</span>
  </>
  ) 
}

function Rooms() {

}

function NewRoom() {

}

function SendMessage() {

}

function LoginModal() {

}

function App() {
    let conn: WebSocket
    const [rooms, SetRooms] = useState([])

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
      <Flex flexDirection="column" width="100%">
        <Flex justifyContent="space-between" width="100%" px="2rem" backgroundColor="#34495e" paddingTop="1rem" paddingBottom="1rem">
          <Avatar name="Carlos" />
          <Button>Logout</Button>
        </Flex>
        <Flex width="100%" >
          <Flex flex="3" backgroundColor="#ecf0f1" height="90vh">
            
          </Flex>
          
          <Flex flex="7" height="100%">

          </Flex>

        </Flex>

      </Flex>
    );
}

export default App;
