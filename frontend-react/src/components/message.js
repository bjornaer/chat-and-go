import React from 'react';

const Message = ({message, user}) => (
    <li className={`chat ${user === message.username ? "right" : "left"}`}>
        <strong>{message.username}</strong>
        {message.content}
        <p className="time-of-message">{message.timestamp}</p>
    </li>
);

export default Message;

/*
export default ({ name, message }) =>
  <p>
    <strong>{name}</strong> <em>{message}</em>
  </p>
*/