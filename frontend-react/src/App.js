import React, { Component } from 'react';
import Chatroom from './components/chatroom';

class App extends Component {
  constructor(props) {
    // bind methods to whole component context
    this.storeUser = this.storeUser.bind(this);
  }

  render() {
    return (
      <div>
        <Login storeUser={storeUser}/>
        ${// <Chatroom username={} profPic={}/>
        }
      </div>
    );
  }

  storeUser() {

  }

  componentDidMount() {
    // Show Login -> if exists show Chatroom History 
    // else store username+password then show Chatroom
    fetch('/test').then((res) => {
      return res.json();
    }).then((res) => {
      this.setState({res});
    }).catch((err) => {
      this.setState({err});
    });
  }
}

export default App;
