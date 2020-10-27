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

```

## How to test
There are three users available: `carlos`, `john` and `gerard`, the password for the three users is `1234`.
Once the user make a successfully login a session key is stored in the browser.
For this reason a new window must open in private mode or in other browser.

## TODO
1. Refactor command procesor
2. Add unit test
3. Refactor web app
4. Add documentation and design diagrams
