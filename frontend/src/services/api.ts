const API_BASE_URL = 'http://localhost:8081/api';

export interface Contact {
  name: string;
  email: string;
  phone: string;
  message: string;
}

export interface Home {
  id: number;
  title: string;
  description: string;
}

export interface About {
  id: number;
  content: string;
}

export interface Project {
  id: number;
  name: string;
  description: string;
  message: string;
}

export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  message: string;
  token: string;
  user: {
    id: number;
    username: string;
    email: string;
  };
}

export const apiService = {
  // Auth API
  async login(credentials: LoginRequest) {
    const response = await fetch(`${API_BASE_URL}/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(credentials)
    });
    return response.json();
  },

  // Contact API
  async sendContact(contact: Contact) {
    const response = await fetch(`${API_BASE_URL}/contact`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(contact)
    });
    return response.json();
  },

  // Home API
  async getHome() {
    const response = await fetch(`${API_BASE_URL}/home`);
    return response.json();
  },

  // About API
  async getAbout() {
    const response = await fetch(`${API_BASE_URL}/about`);
    return response.json();
  },

  // Projects API
  async getProjects() {
    const response = await fetch(`${API_BASE_URL}/projects`);
    return response.json();
  }
};