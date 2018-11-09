import React from 'react';

const Message = ({message, user, time}) => (
    <li className={`chat ${user === message.username ? "right" : "left"}`}>
        {user !== message.username
            && <img src={message.img} alt={`${message.username}'s profile pic`} />
        }
        <strong>{chat.username} - {time}</strong>
        {chat.content}
    </li>
);

export default Message;

/*
export default ({ name, message }) =>
  <p>
    <strong>{name}</strong> <em>{message}</em>
  </p>
*/