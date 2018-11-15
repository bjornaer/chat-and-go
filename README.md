# chat-and-go   
[![Go Report Card](https://goreportcard.com/badge/github.com/Draska/chat-and-go)](https://goreportcard.com/report/github.com/Draska/chat-and-go)

- [What is this?](#What-is-this?)
- [Setup](#Setup)
- [Endpoints](#Endpoints)
- [TO-DO](#TO-DO)
- [Bibliography](#Bibliography)

### What is this?
This is a chat web app running on Golang and VueJS. This is **_Chat-And-Go!_** 

To login enter an email and choose an username(and remember them!)

Don't worry, the email is just to get you a nice gravatar icon! 

If the preloaded history is not enough, click `Load Older Messages` until you reach the very beggining of the chat history _(which is ChantenGo saying hi basically)_

Try sending an emoji to your pals! :wink: :wink: `:wink:`

### Setup:
1. [Install Docker](https://docs.docker.com/engine/installation/)
2. Run the dockerized chat project
    1. `docker-compose up gochat`
    2. Test that it's [running](http://localhost:18000/test)

To restart the project

    docker-compose down
    docker-compose up gochat

To see schema changes, remove the old db volume by adding `-v` when stopping

    docker-compose down -v

To see code changes, rebuild by adding `--build` when starting

    docker-compose up --build gochat

If you run into issues connecting to the db on startup, try restarting (without the `-v` flag).

**[This project is dockerized](#setup), so you don't need to suffer configuring an environment!**
### Endpoints
The server handles the following endpoints:
- `/` --> entrypoint to the chat.
- `/login` --> logs in if you are already registered, otherwise registers you. Why? We don't have much auth on this bad boy.
- `/newMessages?id={user_id}`--> gets unread messages for the specified user --> this endpoint is not used by the chat itself. My criteria for _unread_ is very debatable*.
- `/history?oldest={oldest_message_displayed_id}` --> brings history to the chat from the oldest message already loaded going backwards.
- `/ws` --> establishes a websocket connection! Registers the clients, and sends new messages to the broadcast unit. This lil' guy has all the _mojo_
- `/test` --> just spits out a json response, not used by the chat. Doesn't ping the DB or anything nice. Could be a health check endpoint though, will change it some day.

*I am deciding that "unread" is all the messages that were sent since the last interaction from the user.

### TO-DO
- [ ] Normalize DB --> kill the data redundancies
- [ ] Add tests!!!
- [ ] Split the frontend in a microservice in itslef
- [ ] Encrypt Message content
- [ ] Add Auth - and also encrypt passwords when they are added


### Bibliography
- [WebSockets](https://github.com/gorilla/websocket/tree/master/examples/chat)
- [GORM](http://gorm.io/docs/index.html)
- [Gorilla/Mux router](https://github.com/gorilla/mux)
- [VueJS important thingies](https://vuejs.org/v2/guide/events.html)
- [How to use Gravatar](https://en.gravatar.com/site/implement/)
- [Materialize CSS](https://materializecss.com/getting-started.html)
