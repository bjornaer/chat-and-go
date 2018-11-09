import React from 'react';
import ReactDOM from 'react-dom';
import './App.css';

import Message from './Message.js';

class Chatroom extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            messages: [{
                username: "Chatengo",
                content: <p>Welcome!</p>
            }]
        };
        this.ws;
        this.initSocket = this.initSocket.bind(this);

        this.submitMessage = this.submitMessage.bind(this);
    }

    initSocket () {
        this.ws = new WebSocket("ws://" + window.location.host + "/ws");
        this.ws.onmessage = (msg) => {
          this.state.messages.push(msg.data);
          this.setState({ messages: this.state.messages });
        }
    }

    generateTimestamp () {
        var iso = new Date().toISOString();
        return iso.split("T")[1].split(".")[0];
    }

    componentDidMount() {
        var _this = this;
        fetch('/history')
        .then( response => {
            if(response.status !== 200) {
                console.log('Whoops! Not the expected status! Status:' + response.status);
                return
            }
            response.json()
            .then( data => {
                _this.setState({
                    messages: _this.state.messages.concat(data.messages) // data should look like: {messages: []}
                });
            })
        })
        .catch(err => {
            console.log('Error retrieving history: -S', err)
        })
        this.initSocket()
        this.scrollToBttm();
    }

    componentDidUpdate() {
        this.scrollToBttm();
    }

    scrollToBttm() {
        ReactDOM.findDOMNode(this.refs.messages).scrollTop = ReactDOM.findDOMNode(this.refs.messages).scrollHeight;
    }

    submitMessage(e) {
        e.preventDefault();

        this.setState({
            messages: this.state.messages.concat([{
                username: props.username,
                content: <p>{ReactDOM.findDOMNode(this.refs.msg).value}</p>
            }])
        }, () => {
            ReactDOM.findDOMNode(this.refs.msg).value = "";
        });
    }

    render() {
        const { messages } = this.state;

        return (
            <div className="chatroom">
                <h3>Chat-And-Go</h3>
                <ul className="chats" ref="chats">
                    {
                        messages.map((message) => 
                            <Message chat={message} user={props.username} time={this.generateTimestamp()}/>
                        )
                    }
                </ul>
                <form className="input" onSubmit={(e) => this.submitMessage(e)}>
                    <input type="text" ref="msg" />
                    <input type="submit" value="Submit" />
                </form>
            </div>
        );
    }
}

export default Chatroom;