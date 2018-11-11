import React from 'react';

const Message = ({message, user, time}) => (
    <li className={`chat ${user === message.username ? "right" : "left"}`}>
        <strong>{message.username} - {time}</strong>
        {message.content}
    </li>
);

export default Message;

/*
export default ({ name, message }) =>
  <p>
    <strong>{name}</strong> <em>{message}</em>
  </p>
*/