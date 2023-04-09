import React, { useEffect, useState } from 'react';
import {Flex, Avatar, Button, Input, Stack, Text, useDisclosure, Modal, ModalOverlay, ModalContent, ModalHeader, ModalBody} from '@chakra-ui/core';
import './App.css';
import moment from 'moment'

const callAPI = (url: string, options: any) : Promise<any> => {
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
    code: string
    message: string
}

interface MessageResponse {
    rid: string
    uid: string
    type: string
    message: string
    error: Error
    id: string
    created: string
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

interface Profile {
  id: string
  name: string
}

interface Room {
  id: string
  name: string
}


interface Session {
  name: string
  token: string
}


function App() {
    const ws = React.useRef<WebSocket>(null)
    const [rooms, SetRooms] = useState<Room[]>([])
    const [messages, SetMessages] = useState<Message[]>([])
    const [profiles, SetProfiles] = useState<Profile[]>([])

    const [err, setErr] = useState<string>("")
    const [session, setSession] = useState<Session>({name: "", token: ""})
    const [currentRoom, SetCurrenRoom] = useState<string>("general")
    
    const [roomVal, setRoomVal] = React.useState("");
    const handleRoomValChange = (e: any) => setRoomVal(e.target.value);
    
    const [msgVal, setMsgVal] = React.useState("");
    const handleMsgValChange = (e: any) => setMsgVal(e.target.value);

    const [userVal, setUserVal] = React.useState("");
    const handleUserValChange = (e: any) => setUserVal(e.target.value);
    const [passVal, setPassVal] = React.useState("");
    const handlePassValChange = (e: any) => setPassVal(e.target.value);

    const { isOpen, onOpen, onClose } = useDisclosure(true);
    
    useEffect(() => {
      const s = localStorage.getItem("session")
      if (s === null) return
      const session = JSON.parse(s) as Session
      setSession(session)
      onClose()
    }, [onClose])

    useEffect(() => {
        // @ts-ignore
        ws.current = new WebSocket("ws://127.0.0.1:8080/ws")
        ws.current.onclose = function(evt) {console.log(evt)};
        ws.current.onopen = function(evt) {console.log(evt)};
        return () => {
          // @ts-ignore
          ws.current.close();
        };
    }, [])

    useEffect(() => {
      if (!ws.current) return;
      ws.current.onmessage = function(evt) {
        const response = JSON.parse(evt.data) as MessageResponse
        if (response.rid === currentRoom && response.error === null) {
          const msg = {
            id: response.id,
            rid: response.rid,
            uid: response.uid,
            message: response.message,
            created: response.created,
          }
          SetMessages(prev => [...prev, msg])
        } else if (response.error !== null) {
          const msg = {
            id: response.id,
            rid: response.rid,
            uid: response.uid,
            message: response.error.message,
            created: response.created,
          }
          SetMessages(prev => [...prev, msg])
        } 
    };
    }, [currentRoom])

    useEffect( () => {
      async function Retrieve() {
        const header = {
          method: "GET",
          headers: {
            'Authorization': `Bearer ${session.token}`, 
          },    
        };
        const rooms = await callAPI("http://127.0.0.1:8080/rooms", header)
        SetRooms(rooms)
      }

      Retrieve()
    }, [session])

    useEffect( () => {
      async function Retrieve() {
        const header = {
          method: "GET",
          headers: {
            'Authorization': `Bearer ${session.token}`, 
          },    
        };
        const profiles = await callAPI("http://127.0.0.1:8080/profiles", header)
        SetProfiles(profiles)
      }

      Retrieve()
    }, [session])

    useEffect(() => {
      async function Retrieve() {
        const header = {
          method: "GET",
          headers: {
            'Authorization': `Bearer ${session.token}`, 
          },    
        };
        const messages = await callAPI(`http://127.0.0.1:8080/messages/${currentRoom}`, header)
        SetMessages(messages)
      }
      Retrieve()
    }, [currentRoom, session])

    const createRoom = () => {
      async function Create() {        
        const body = {
          name: roomVal
        };

        const header = {
          method: "POST",
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${session.token}`, 
          },
          body: JSON.stringify(body),
          };
          const resp = await callAPI(`http://127.0.0.1:8080/rooms`, header)
          SetRooms(prev => [...prev, resp])
        }
        Create()
    }

    const login = () => {
      async function Login() {        
        const body = {
          username: userVal,
          password: passVal
        };

        const header = {
          method: "POST",
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(body),
          };
          try {
            const resp = await callAPI(`http://127.0.0.1:8080/login`, header)
            const s = {name: userVal, token: resp.access_token}
            setSession(s)
            setErr("")
            localStorage.setItem("session", JSON.stringify(s))
            onClose()
            setPassVal("")
            setUserVal("")           
          } catch(err) {
            setErr(err.message)
          }
        }
        Login()
    }

    const getMyID = (): string => {
      const val = profiles.find((p) => p.name === session.name)
      return val?.id || ""
    }

    const getProfile = (id: string) : Profile => {
      if (id === "bot") return {id: "bot", name: "BOT"}
      const p = profiles.find((p) => p.id === id)
      return p || {name: "", id: ""}
    }


    const sendMessage = (msg: string = "") => {
        const message: MessageRequest = {
            rid: currentRoom,
            uid: getMyID(),
            type: "message",
            message: msg,
        }
        // @ts-ignore
        ws.current.send(JSON.stringify(message))
    }
    
    const sendCommand = (cmd: string = "") => {
        const message: MessageRequest = {
            rid: currentRoom,
            uid: getMyID(),
            type: "command",
            message: cmd.substring(1),
        }
        // @ts-ignore
        ws.current.send(JSON.stringify(message))
    }

    const send = (): void =>  {
      const msg = msgVal;
      if (msg === "") return
      if (msg.startsWith("/stock")) {
        sendCommand(msg)
      } else {
        sendMessage(msg)
      }
      setMsgVal("")
    }
    const onKeyUp = (event: any) => {
      if (event.charCode === 13) {
        send()
      }
    }

    const logout = () => {
      localStorage.removeItem("session")
      setSession({name: "", token: ""})
      onOpen()
    }

    return (
      <Flex flexDirection="column" width="100%">
        <Flex justifyContent="space-between" width="100%" px="2rem" backgroundColor="#34495e" paddingTop="1rem" paddingBottom="1rem">
          <Avatar name={session.name} />
          <Text fontSize="2rem" color="#FFF">{currentRoom}</Text>
          <Button onClick={logout}>Logout</Button>
        </Flex>
        <Flex width="100%" >
          <Flex flex="3" backgroundColor="#ecf0f1" height="90vh" flexDirection="column">
            <Flex margin=".5rem" justifyContent="center">
              <Input placeholder="New Room" value={roomVal} onChange={handleRoomValChange}/>
              <Button onClick={createRoom}>Create</Button>
            </Flex>
            {
              rooms.map((r) => (
                <Stack cursor="pointer" key={r.id} isInline backgroundColor="#2980b9" margin="0.2rem" onClick={() => SetCurrenRoom(r.name)}>
                  <Text>{r.name}</Text>
                </Stack>
              ))
            }
          </Flex>
          
          <Flex flex="7" height="100%" flexDirection="column" px="1rem" >
            <Flex flexDirection="column" height="85vh" maxHeight="85vh" overflowY="scroll">
              {
                messages.map((m) => (
                  <Flex key={m.id} flexDirection="column" my="0.3rem">
                    <Stack isInline alignItems="center">
                      <Avatar size="md" name={getProfile(m.uid).name} />
                      <Flex flexDirection="column">
                        <Text fontSize=".8rem">{moment(m.created).format("YYYY-MM-DD HH:mm:ss")}</Text>
                        <Text fontSize=".8rem">{getProfile(m.uid).name}</Text>
                      </Flex>
                    </Stack>
                    <Flex marginLeft="3rem" my="0.5rem">
                    <Text>{m.message}</Text>
                    </Flex>
                  </Flex>
                ))
              }
            </Flex>
            <Flex>
              <Stack isInline width="100%">
                <Input placeholder="" value={msgVal} onChange={handleMsgValChange} onKeyPress={onKeyUp}/>
                <Button onClick={send}>Send</Button>
              </Stack>
            </Flex>

          </Flex>

        </Flex>
        <Modal isOpen={isOpen} onClose={onClose} closeOnEsc={false} closeOnOverlayClick={false}>
        <ModalOverlay />
        <ModalContent>
          <ModalHeader>Login</ModalHeader>
          <ModalBody>
              <Stack width="100%">
                <Input placeholder="username" value={userVal} onChange={handleUserValChange}/>
                <Input placeholder="password" type="password" value={passVal} onChange={handlePassValChange}/>
                <Button onClick={login}>Login</Button>
                {
                  err !== "" && <Text>Error: {err}</Text>
                }
              </Stack>
          </ModalBody>
        </ModalContent>
      </Modal>
      </Flex>
    );
}

export default App;
