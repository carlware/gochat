# Chat written in golang with bot command

This is chat application composed with three micro-services, one handle the user authentication, another chatrooms, and the third one get stock prices.

## How to Run
```
# install rabbitmq
# on Mac
brew install rabbitmq

# with docker
docker run -it --rm --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3-management

# run
go run main.go

# visit
http://localhost:8080/

```

## Design
The following diagram shows the pieces a high level. Where there are three microservices.  

![high](https://github.com/carlware/gochat/blob/crl/dev3/design/high.svg "High")

In the sequence diagram below is shown what happen when a user send a message and a command by the websocket. The sequence is different.  

![sequence](https://github.com/carlware/gochat/blob/crl/dev3/design/sequence.svg "Sequence")

Once the user makes login a WS connection is established between the browser and the client. This connection is closed until the user close the browser tab. Then, the user can send any type the data to the server but the microservice handle a specific request and response types, any attempt to send data that does not satisfy the request definition will return an error response.  

Until now, there are two type of requests. Message request only broadcast the message to all the clients connected. Command request first computed the result and then broadcast the result to all the clients too.  

By the design, there are two types of command, `fast` and `queue`. Fast commands are executed immediately and queue and send to a MQ broker. Then another microservice listen for the command and send the result back.  

This diagram show the package architecture. Some interfaces where created in order to abstract implementation details and use the inversion of control pattern. That means that we can change the implementation details of the MQ broker and the database.  

![arch](https://github.com/carlware/gochat/blob/crl/dev3/design/arch.svg "Architecture")


## How to test
There are three users available: `carlos`, `john` and `gerard`, the password for the three users is `1234`.  

Once the user make a successfully login a session key is stored in the browser, for this reason a new window must open in private mode or in other browser. There is a small bug that when a user creates new room all the users need to refresh the browser, this will be solved.  


## Future work
I noticed that allow `fast` commands is not maintainable. If this service scale maintain the code is going to result in a bottleneck. If new commands are added a new version of the service will be deployed. This characteristic will be removed. It is better only have commands that are queued and another microservice listen and returns the result. However a nice feature would be wait certain time and if nobody respond return an error to client.  


## TODO
0. When a room is created notify other users
1. Update channels with directions
2. Wait N ms for a queued command to response and if the result is not returned send an error to the user
3. Add unit test
4. Refactor web app
5. Show an error if the room cannot be created in the frontend
6. Create a link struct N-to-N where a user can be linked to any rooms
7. Add a button to start a private conversation with someone in the frontend (on the backend this will create room)
8. Support broadcast to clients that are in a specific room (this can support private rooms too)
9. Gracefully shutdown microservices and gorutines
