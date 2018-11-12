import React from 'react';
import ReactDOM from 'react-dom';

import Message from './Message.js'
import SubmitOnEnterForm from './SubmitOnEnter.js';

class Chatroom extends React.Component {
  constructor(props) {
    super(props);
    this.state = { messages: [{
        username: "Chatengo",
        content: <p>Welcome!</p>
    }] };
    // eslint-disable-next-line
    this.ws;
    this.initSocket = this.initSocket.bind(this);
    this.sendMessage = this.sendMessage.bind(this);
  }

  initSocket () {
    this.ws = new WebSocket("ws://" + window.location.host + "/ws");
    this.ws.onmessage = (msg) => {
      this.state.messages.push(msg.data);
      this.setState({ messages: this.state.messages });
    }
  }

  componentDidMount () {
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
        });
        this.initSocket();
        this.scrollToBttm();
  }

  componentDidUpdate() {
    this.scrollToBttm();
  }

  scrollToBttm() {
    ReactDOM.findDOMNode(this.refs.messages).scrollTop = ReactDOM.findDOMNode(this.refs.messages).scrollHeight;
  }

  generateTimestamp () {
    var iso = new Date().toISOString();
    return iso.split("T")[1].split(".")[0];
  }

  sendMessage (message) {
    this.ws.send(
      JSON.stringify({
        username: this.props.user.Username,
        content: (this.generateTimestamp() + " <" + this.props.user.Username + "> " + message)
      })
    );
  }

  render () {
    const { messages } = this.state;  
    return (
      <div>
        <h3>Chat-And-Go</h3>
        <pre className="chat-room">
          {/*this.state.messages.join('\n')*/}
          <ul className="chat-room" ref="chats">
                {
                    messages.map((message) => 
                        <Message chat={message} user={this.props.user.Username}/>
                    )
                }
            </ul>
        </pre>
        <SubmitOnEnterForm
          placeholder="press enter to send"
          // eslint-disable-next-line
          onSubmit={this.sendMessage} />
        <form className="input" onSubmit={(e) => {e.preventDefault();this.sendMessage}}>
            <input type="text" ref="msg" />
            <input type="submit" value="Submit" />
        </form>
      </div>
    )
  }
}

export default Chatroom;
//module.exports = Chatroom;