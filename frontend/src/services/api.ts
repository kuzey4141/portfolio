const API_BASE_URL = 'http://localhost:8081/api';

// Interface for contact form submission (frontend to backend)
export interface ContactForm {
  name: string;
  email: string;
  phone: string;
  message: string;
}

// Interface for contact data from database (backend to frontend) 
export interface Contact {
  id: number;
  name: string;
  email: string;
  phone: string;
  message: string;
  created_at: string;
}

// Interface for home page data
export interface Home {
  id: number;
  title: string;
  description: string;
}

// Interface for about page data
export interface About {
  id: number;
  content: string;
}

// Interface for project data
export interface Project {
  id: number;
  name: string;
  description: string;
  message: string;
}

// Interface for login request
export interface LoginRequest {
  username: string;
  password: string;
}

// Interface for login response
export interface LoginResponse {
  message: string;
  token: string;
  user: {
    id: number;
    username: string;
    email: string;
  };
}

// API service object with all endpoint methods
export const apiService = {
  // Authentication API - user login
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

  // Contact API - send contact form data
  async sendContact(contact: ContactForm) {
    const response = await fetch(`${API_BASE_URL}/contact`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(contact)
    });
    return response.json();
  },

  // Home API - get home page data
  async getHome() {
    const response = await fetch(`${API_BASE_URL}/home`);
    return response.json();
  },

  // About API - get about page data
  async getAbout() {
    const response = await fetch(`${API_BASE_URL}/about`);
    return response.json();
  },

  // Projects API - get projects data
  async getProjects() {
    const response = await fetch(`${API_BASE_URL}/projects`);
    return response.json();
  }
};