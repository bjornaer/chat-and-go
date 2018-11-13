# chat-and-go   
[![Go Report Card](https://goreportcard.com/badge/github.com/Draska/chat-and-go)](https://goreportcard.com/report/github.com/Draska/chat-and-go)

### What is this?
This is a chat web app running on Golang and VueJS. This is **_Chat-And-Go!_** 

To login enter an email and choose an username(and remember them!)

Don't worry, the email is just to get you a nice gravatar icon! 
_(So_ use an actual email _, if ye already configured this stuff, it's gonna get your image - and that rocks.)_

**This project is dockerized, so you don't need to suffer configuring an environment!**

The server handles the following endpoints:
- `/` --> entrypoint to the chat.
- `/login` --> logs in if you are already registered, otherwise registers you.
- `/history` --> brings recent history to the chat
- `/ws` --> establishes a websocket connection!

To get the project up and running:
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


##### Bibliography (sort of)
- [WebSockets](https://github.com/gorilla/websocket/tree/master/examples/chat)
- [VueJS important thingies](https://vuejs.org/v2/guide/events.html)
- [How to use Gravatar](https://en.gravatar.com/site/implement/)
