import React, { useState } from 'react';
import { apiService } from '../services/api';

function AdminLogin() {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [message, setMessage] = useState('');

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    
    try {
      const response = await apiService.login({ username, password });
      setMessage(`Success: ${response.message}`);
    } catch (error) {
      setMessage('Error: Login failed');
    }
  };

  return (
    <div style={{ padding: '50px', textAlign: 'center' }}>
      <h2>Admin Login</h2>
      <form onSubmit={handleLogin}>
        <div>
          <input
            type="text"
            placeholder="Username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            style={{ margin: '10px', padding: '10px' }}
          />
        </div>
        <div>
          <input
            type="password"
            placeholder="Password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            style={{ margin: '10px', padding: '10px' }}
          />
        </div>
        <button type="submit" style={{ margin: '10px', padding: '10px 20px' }}>
          Login
        </button>
      </form>
      {message && <p>{message}</p>}
      <p>Default: admin / admin123</p>
    </div>
  );
}

export default AdminLogin;