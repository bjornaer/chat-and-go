/* import React from 'react';
//import ReactDOM from 'react-dom';

//import SubmitOnEnterForm from './components/SubmitOnEnterForm.jsx';
import Chatroom from './components/Chatroom.js';
import Login from './components/Login.js';

class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = {user: null};
    this.setUsername = this.setUsername.bind(this);
  }

  componentDidMount() {
    // Show Login -> if exists show Chatroom History 
    // else store username+password then show Chatroom
    /* fetch('/test').then((res) => {
      return res.json();
    }).then((res) => {
      this.setState({res});
    }).catch((err) => {
      this.setState({err});
    });
  }

  setUsername(username) {
    // CHECK THE SERVER URL CORRECTLY DUDE
    // http://backend:8000
    fetch('/login', {
      method: 'POST',
      headers: {
        "Content-type": "application/json; charset=UTF-8"
      },
      body: `"Username": "${username}"`
    })
      .then((response) => { return response.json(); })
      .then((json) => { this.setState({ user: json.Username }); });
  }

  render () {
    let loginForm = <Login setUsername={this.setUsername}/>
    let sharedChatroom = <Chatroom user={this.state.user} />

    if (this.state.user) {
      return sharedChatroom;
    } else {
      return loginForm;
    }
  }
}

export default App;
*/
import React from 'react';
import ReactDOM from 'react-dom';

import SubmitOnEnterForm from './components/SubmitOnEnter.js';
import ChatRoom from './components/Chatroom.js';

class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = {user: null};
    this.setUsername = this.setUsername.bind(this);
  }

  setUsername(username) {
    /* fetch('/login', {
      method: 'POST',
      headers: {
        "Content-type": "application/x-www-form-urlencoded; charset=UTF-8"
      },
      body: 'username='+username
    }) */
    fetch('/login', {
      method: 'POST',
      headers: {
        "Content-type": "application/json; charset=UTF-8"
      },
      body: `"Username": "${username}"`
    })
      .then((response) => { return response.json(); })
      .then((json) => { this.setState({ user: json.user }); });
  }

  render () {
    let loginForm = (
      // set propTypes in this component
      <div>
        <h1>Login</h1>
        <SubmitOnEnterForm
          placeholder="Enter your username"
          onSubmit={this.setUsername} />
      </div>
    )

    let sharedChatRoom = <ChatRoom user={this.state.user} />

    if (this.state.user) {
      return sharedChatRoom;
    } else {
      return loginForm;
    }
  }
}

export default App;