import React from 'react';

import SubmitOnEnterForm from './SubmitOnEnter.js';

const Login = ({setUsername}) => (
  <div className='login'>
    <h1>Login</h1>
    <SubmitOnEnterForm
      placeholder="Enter your username"
      onSubmit={setUsername} />
  </div>
);

/* const style = {
 margin: 15,
}; */
export default Login;
