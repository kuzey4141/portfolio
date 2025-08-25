import React, { useState, useEffect } from 'react';
import { LogOut, MessageSquare, Home, FolderOpen, User, Edit, Trash2, Save, Plus, TrendingUp, BarChart3, Activity } from 'lucide-react';
import { apiService, Contact as ApiContact } from '../services/api';
import AdminProjects from './AdminProjects';

interface Contact {
  id: number;
  name: string;
  email: string;
  phone: string;
  message: string;
  created_at: string;
}

interface Project {
  id: number;
  name: string;
  description: string;
  message: string;
}

interface HomeData {
  id: number;
  title: string;
  description: string;
}

interface AboutData {
  id: number;
  content: string;
}

interface UpdateFormData {
  id: number;
  title?: string;
  description?: string;
  content?: string;
}

interface AdminDashboardProps {
  onLogout?: () => void;
}

const AdminDashboard = ({ onLogout }: AdminDashboardProps): JSX.Element => {
  const [activeTab, setActiveTab] = useState<string>('dashboard');
  const [contacts, setContacts] = useState<Contact[]>([]);
  const [projects, setProjects] = useState<Project[]>([]);
  const [homeData, setHomeData] = useState<HomeData | null>(null);
  const [aboutData, setAboutData] = useState<AboutData | null>(null);
  const [loading, setLoading] = useState<boolean>(false);
  const [message, setMessage] = useState<string>('');
  const [editingHome, setEditingHome] = useState<boolean>(false);
  const [editingAbout, setEditingAbout] = useState<boolean>(false);

  // Fixed API base URL
  const API_BASE_URL = "http://3.78.181.203:8081/api";

  useEffect(() => {
    loadData();
  }, [activeTab]);

  const loadData = async (): Promise<void> => {
    setLoading(true);
    try {
      const token = localStorage.getItem('authToken');
      
      console.log(`Loading data for tab: ${activeTab}`);
      
      if (activeTab === 'contacts') {
        console.log('Loading contacts with token:', token ? 'Token exists' : 'No token');
        
        const response = await fetch(`${API_BASE_URL}/admin/contact`, { 
          method: 'GET',
          headers: { 
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json'
          }
        });
        
        console.log('Contacts response status:', response.status);
        
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        const data = await response.json();
        console.log('Contacts data received:', data);
        
        if (Array.isArray(data)) {
          setContacts(data);
          console.log('Contacts set to state:', data.length, 'items');
        } else {
          console.error('Contacts data is not an array:', data);
          setContacts([]);
        }
      } else if (activeTab === 'projects') {
        console.log('Loading projects from public endpoint...');
        
        // Use public projects endpoint - no token needed
        const response = await fetch(`${API_BASE_URL}/projects`, { 
          method: 'GET',
          headers: { 
            'Content-Type': 'application/json'
          }
        });
        
        console.log('Projects response status:', response.status);
        
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        const data = await response.json();
        console.log('Projects data received:', data);
        
        if (Array.isArray(data)) {
          setProjects(data);
          console.log('Projects set to state:', data.length, 'items');
        } else {
          console.error('Projects data is not an array:', data);
          setProjects([]);
        }
      } else if (activeTab === 'home') {
        console.log('Loading home data...');
        const data = await apiService.getHome();
        console.log('Home data received:', data);
        setHomeData(data[0] || null);
      } else if (activeTab === 'about') {
        console.log('Loading about data...');
        const data = await apiService.getAbout();
        console.log('About data received:', data);
        setAboutData(data[0] || null);
      }
    } catch (error) {
      console.error('Data loading error:', error);
      setMessage('Error loading data: ' + (error instanceof Error ? error.message : 'Unknown error'));
    } finally {
      setLoading(false);
    }
  };

  const handleLogout = (): void => {
    console.log('Logging out...');
    localStorage.removeItem('authToken');
    if (onLogout) {
      onLogout();
    } else {
      window.location.href = '/';
    }
  };

  const deleteContact = async (id: number): Promise<void> => {
    if (!window.confirm('Are you sure you want to delete this contact?')) {
      return;
    }

    try {
      console.log('Deleting contact:', id);
      const token = localStorage.getItem('authToken');
      const response = await fetch(`${API_BASE_URL}/admin/contact/${id}`, { 
        method: 'DELETE',
        headers: { 'Authorization': `Bearer ${token}` }
      });
      
      console.log('Delete contact response status:', response.status);
      
      if (response.ok) {
        setMessage('Contact deleted successfully');
        loadData();
        setTimeout(() => setMessage(''), 3000);
      } else {
        const errorData = await response.json();
        setMessage(`Error deleting contact: ${errorData.error || 'Unknown error'}`);
      }
    } catch (error) {
      console.error('Delete contact error:', error);
      setMessage('Error deleting contact');
    }
  };

  const updateHome = async (formData: UpdateFormData): Promise<void> => {
    try {
      console.log('Updating home data:', formData);
      const token = localStorage.getItem('authToken');
      const response = await fetch(`${API_BASE_URL}/admin/home`, { 
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify(formData)
      });
      
      console.log('Update home response status:', response.status);
      
      if (response.ok) {
        setMessage('Home page updated successfully');
        setEditingHome(false);
        loadData();
        setTimeout(() => setMessage(''), 3000);
      } else {
        const errorData = await response.json();
        setMessage(`Error updating home page: ${errorData.error || 'Unknown error'}`);
      }
    } catch (error) {
      console.error('Update home error:', error);
      setMessage('Error updating home page');
    }
  };

  const updateAbout = async (formData: UpdateFormData): Promise<void> => {
    try {
      console.log('Updating about data:', formData);
      const token = localStorage.getItem('authToken');
      const response = await fetch(`${API_BASE_URL}/admin/about`, { 
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify(formData)
      });
      
      console.log('Update about response status:', response.status);
      
      if (response.ok) {
        setMessage('About page updated successfully');
        setEditingAbout(false);
        loadData();
        setTimeout(() => setMessage(''), 3000);
      } else {
        const errorData = await response.json();
        setMessage(`Error updating about page: ${errorData.error || 'Unknown error'}`);
      }
    } catch (error) {
      console.error('Update about error:', error);
      setMessage('Error updating about page');
    }
  };

  const DashboardContent = (): JSX.Element => (
    <div style={{ padding: '32px' }}>
      <h1 style={{ fontSize: '32px', fontWeight: '800', color: '#1a1a1a', marginBottom: '8px' }}>
        Welcome back, Admin! 👋
      </h1>
      <p style={{ color: '#6b7280', fontSize: '16px', marginBottom: '40px' }}>
        Here is what is happening with your portfolio today
      </p>

      {message && (
        <div style={{
          padding: '16px 20px',
          borderRadius: '12px',
          marginBottom: '24px',
          backgroundColor: message.includes('Error') ? '#fef2f2' : '#f0f9ff',
          border: `1px solid ${message.includes('Error') ? '#fecaca' : '#bae6fd'}`,
          color: message.includes('Error') ? '#dc2626' : '#0369a1'
        }}>
          {message}
          <button onClick={() => setMessage('')} style={{ marginLeft: '16px', fontSize: '14px', textDecoration: 'underline', background: 'none', border: 'none', cursor: 'pointer', color: 'inherit' }}>
            Dismiss
          </button>
        </div>
      )}
      
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(280px, 1fr))', gap: '24px', marginBottom: '40px' }}>
        <div style={{ background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)', borderRadius: '20px', padding: '28px', color: 'white' }}>
          <div style={{ display: 'flex', alignItems: 'center', marginBottom: '16px' }}>
            <MessageSquare size={24} style={{ marginRight: '12px' }} />
            <span style={{ fontSize: '14px', fontWeight: '500' }}>Total Messages</span>
          </div>
          <div style={{ fontSize: '36px', fontWeight: '700', marginBottom: '8px' }}>
            {contacts.length}
          </div>
        </div>
        
        <div style={{ background: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)', borderRadius: '20px', padding: '28px', color: 'white' }}>
          <div style={{ display: 'flex', alignItems: 'center', marginBottom: '16px' }}>
            <FolderOpen size={24} style={{ marginRight: '12px' }} />
            <span style={{ fontSize: '14px', fontWeight: '500' }}>Active Projects</span>
          </div>
          <div style={{ fontSize: '36px', fontWeight: '700', marginBottom: '8px' }}>
            {projects.length}
          </div>
        </div>
        
        <div style={{ background: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)', borderRadius: '20px', padding: '28px', color: 'white' }}>
          <div style={{ display: 'flex', alignItems: 'center', marginBottom: '16px' }}>
            <BarChart3 size={24} style={{ marginRight: '12px' }} />
            <span style={{ fontSize: '14px', fontWeight: '500' }}>Page Views</span>
          </div>
          <div style={{ fontSize: '36px', fontWeight: '700', marginBottom: '8px' }}>
            2.4k
          </div>
        </div>
      </div>
      
      <div style={{ background: 'white', borderRadius: '20px', padding: '32px', boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.1)' }}>
        <h3 style={{ fontSize: '24px', fontWeight: '700', marginBottom: '24px', color: '#1f2937' }}>
          Quick Actions
        </h3>
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(300px, 1fr))', gap: '20px' }}>
          <button onClick={() => setActiveTab('contacts')} style={{ padding: '24px', background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)', color: 'white', borderRadius: '16px', border: 'none', cursor: 'pointer', textAlign: 'left' }}>
            <h4 style={{ fontSize: '18px', fontWeight: '600', marginBottom: '8px' }}>📧 View Messages</h4>
            <p style={{ fontSize: '14px', opacity: 0.9, margin: 0 }}>Check and manage contact form submissions</p>
          </button>
          
          <button onClick={() => setActiveTab('projects')} style={{ padding: '24px', background: 'linear-gradient(135deg, #10b981 0%, #059669 100%)', color: 'white', borderRadius: '16px', border: 'none', cursor: 'pointer', textAlign: 'left' }}>
            <h4 style={{ fontSize: '18px', fontWeight: '600', marginBottom: '8px' }}>🚀 Manage Projects</h4>
            <p style={{ fontSize: '14px', opacity: 0.9, margin: 0 }}>Add, edit, and delete portfolio projects</p>
          </button>
          
          <button onClick={() => setActiveTab('home')} style={{ padding: '24px', background: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)', color: 'white', borderRadius: '16px', border: 'none', cursor: 'pointer', textAlign: 'left' }}>
            <h4 style={{ fontSize: '18px', fontWeight: '600', marginBottom: '8px' }}>🏠 Edit Home Page</h4>
            <p style={{ fontSize: '14px', opacity: 0.9, margin: 0 }}>Update your main page title and description</p>
          </button>
        </div>
      </div>

      {/* Debug Info */}
      <div style={{ background: '#f8f9fa', borderRadius: '12px', padding: '20px', marginTop: '20px', fontSize: '14px', color: '#6b7280' }}>
        <h4 style={{ margin: '0 0 10px 0', color: '#374151' }}>Debug Info:</h4>
        <p style={{ margin: '5px 0' }}>Contacts loaded: {contacts.length}</p>
        <p style={{ margin: '5px 0' }}>Projects loaded: {projects.length}</p>
        <p style={{ margin: '5px 0' }}>Home data: {homeData ? 'Loaded' : 'Not loaded'}</p>
        <p style={{ margin: '5px 0' }}>About data: {aboutData ? 'Loaded' : 'Not loaded'}</p>
        <p style={{ margin: '5px 0' }}>Auth token: {localStorage.getItem('authToken') ? 'Present' : 'Missing'}</p>
      </div>
    </div>
  );

  const ContactsContent = (): JSX.Element => (
    <div style={{ padding: '32px' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '32px' }}>
        <h2 style={{ fontSize: '28px', fontWeight: '700', color: '#1f2937' }}>
          💬 Contact Messages
        </h2>
        <button 
          onClick={loadData}
          style={{
            background: '#3b82f6',
            color: 'white',
            border: 'none',
            borderRadius: '8px',
            padding: '10px 20px',
            fontSize: '14px',
            cursor: 'pointer'
          }}
        >
          🔄 Refresh
        </button>
      </div>
      
      {contacts.length === 0 ? (
        <div style={{ textAlign: 'center', padding: '80px 20px', background: 'white', borderRadius: '20px', boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.1)' }}>
          <MessageSquare size={64} style={{ color: '#d1d5db', marginBottom: '16px' }} />
          <p style={{ color: '#6b7280', fontSize: '18px' }}>No messages yet. Check back later!</p>
        </div>
      ) : (
        <div style={{ display: 'flex', flexDirection: 'column', gap: '20px' }}>
          {contacts.map((contact: Contact) => (
            <div key={contact.id} style={{ background: 'white', borderRadius: '16px', padding: '24px', boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.1)', borderLeft: '4px solid #667eea' }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', marginBottom: '16px' }}>
                <div>
                  <h3 style={{ fontSize: '18px', fontWeight: '600', color: '#1f2937', marginBottom: '8px' }}>
                    {contact.name}
                  </h3>
                  <p style={{ color: '#6b7280', marginBottom: '4px' }}>
                    📧 {contact.email} • 📱 {contact.phone}
                  </p>
                  <p style={{ fontSize: '14px', color: '#9ca3af' }}>
                    📅 {contact.created_at ? new Date(contact.created_at).toLocaleDateString() : 'No date'}
                  </p>
                </div>
                <button onClick={() => deleteContact(contact.id)} style={{ background: '#fee2e2', color: '#dc2626', border: 'none', borderRadius: '8px', padding: '8px', cursor: 'pointer' }}>
                  <Trash2 size={16} />
                </button>
              </div>
              <div style={{ background: '#f8fafc', padding: '16px', borderRadius: '12px', color: '#374151', lineHeight: '1.6' }}>
                {contact.message}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );

  const HomeEditContent = (): JSX.Element => (
    <div style={{ padding: '32px' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '32px' }}>
        <h2 style={{ fontSize: '28px', fontWeight: '700', color: '#1f2937' }}>🏠 Home Page Settings</h2>
        <div style={{ display: 'flex', gap: '12px' }}>
          <button 
            onClick={loadData}
            style={{
              background: '#3b82f6',
              color: 'white',
              border: 'none',
              borderRadius: '8px',
              padding: '10px 20px',
              fontSize: '14px',
              cursor: 'pointer'
            }}
          >
            🔄 Refresh
          </button>
          <button onClick={() => setEditingHome(!editingHome)} style={{ display: 'flex', alignItems: 'center', gap: '8px', padding: '12px 20px', background: editingHome ? '#f59e0b' : '#3b82f6', color: 'white', border: 'none', borderRadius: '12px', cursor: 'pointer', fontSize: '14px', fontWeight: '600' }}>
            <Edit size={16} />
            {editingHome ? 'Cancel' : 'Edit'}
          </button>
        </div>
      </div>
      
      {homeData && (
        <div style={{ background: 'white', borderRadius: '20px', padding: '32px', boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.1)' }}>
          {editingHome ? (
            <form onSubmit={(e: React.FormEvent<HTMLFormElement>) => {
              e.preventDefault();
              const formData = new FormData(e.currentTarget);
              updateHome({
                id: homeData.id,
                title: formData.get('title') as string,
                description: formData.get('description') as string
              });
            }}>
              <div style={{ display: 'flex', flexDirection: 'column', gap: '24px' }}>
                <div>
                  <label style={{ display: 'block', fontSize: '14px', fontWeight: '600', color: '#374151', marginBottom: '8px' }}>Page Title</label>
                  <input name="title" type="text" defaultValue={homeData.title} style={{ width: '100%', padding: '16px', border: '2px solid #e5e7eb', borderRadius: '12px', fontSize: '16px', outline: 'none', boxSizing: 'border-box' }} required />
                </div>
                <div>
                  <label style={{ display: 'block', fontSize: '14px', fontWeight: '600', color: '#374151', marginBottom: '8px' }}>Description</label>
                  <textarea name="description" rows={4} defaultValue={homeData.description} style={{ width: '100%', padding: '16px', border: '2px solid #e5e7eb', borderRadius: '12px', fontSize: '16px', outline: 'none', resize: 'vertical', boxSizing: 'border-box' }} required />
                </div>
                <button type="submit" style={{ display: 'flex', alignItems: 'center', gap: '8px', padding: '16px 24px', background: 'linear-gradient(135deg, #10b981 0%, #059669 100%)', color: 'white', border: 'none', borderRadius: '12px', cursor: 'pointer', fontSize: '16px', fontWeight: '600', alignSelf: 'flex-start' }}>
                  <Save size={16} />
                  Save Changes
                </button>
              </div>
            </form>
          ) : (
            <div style={{ display: 'flex', flexDirection: 'column', gap: '24px' }}>
              <div>
                <h3 style={{ fontSize: '18px', fontWeight: '600', color: '#374151', marginBottom: '8px' }}>Current Title</h3>
                <p style={{ fontSize: '20px', color: '#1f2937', fontWeight: '500' }}>{homeData.title}</p>
              </div>
              <div>
                <h3 style={{ fontSize: '18px', fontWeight: '600', color: '#374151', marginBottom: '8px' }}>Current Description</h3>
                <p style={{ fontSize: '16px', color: '#1f2937', lineHeight: '1.6' }}>{homeData.description}</p>
              </div>
            </div>
          )}
        </div>
      )}
      
      {!homeData && !loading && (
        <div style={{ background: 'white', borderRadius: '20px', padding: '32px', boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.1)', textAlign: 'center' }}>
          <p style={{ color: '#6b7280', fontSize: '18px' }}>No home data found. Please refresh or check your database.</p>
        </div>
      )}
    </div>
  );

  const AboutEditContent = (): JSX.Element => (
    <div style={{ padding: '32px' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '32px' }}>
        <h2 style={{ fontSize: '28px', fontWeight: '700', color: '#1f2937' }}>👤 About Page Settings</h2>
        <div style={{ display: 'flex', gap: '12px' }}>
          <button 
            onClick={loadData}
            style={{
              background: '#3b82f6',
              color: 'white',
              border: 'none',
              borderRadius: '8px',
              padding: '10px 20px',
              fontSize: '14px',
              cursor: 'pointer'
            }}
          >
            🔄 Refresh
          </button>
          <button onClick={() => setEditingAbout(!editingAbout)} style={{ display: 'flex', alignItems: 'center', gap: '8px', padding: '12px 20px', background: editingAbout ? '#f59e0b' : '#3b82f6', color: 'white', border: 'none', borderRadius: '12px', cursor: 'pointer', fontSize: '14px', fontWeight: '600' }}>
            <Edit size={16} />
            {editingAbout ? 'Cancel' : 'Edit'}
          </button>
        </div>
      </div>
      
      {aboutData && (
        <div style={{ background: 'white', borderRadius: '20px', padding: '32px', boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.1)' }}>
          {editingAbout ? (
            <form onSubmit={(e: React.FormEvent<HTMLFormElement>) => {
              e.preventDefault();
              const formData = new FormData(e.currentTarget);
              updateAbout({
                id: aboutData.id,
                content: formData.get('content') as string
              });
            }}>
              <div style={{ display: 'flex', flexDirection: 'column', gap: '24px' }}>
                <div>
                  <label style={{ display: 'block', fontSize: '14px', fontWeight: '600', color: '#374151', marginBottom: '8px' }}>About Content</label>
                  <textarea name="content" rows={8} defaultValue={aboutData.content} style={{ width: '100%', padding: '16px', border: '2px solid #e5e7eb', borderRadius: '12px', fontSize: '16px', outline: 'none', resize: 'vertical', boxSizing: 'border-box' }} required />
                </div>
                <button type="submit" style={{ display: 'flex', alignItems: 'center', gap: '8px', padding: '16px 24px', background: 'linear-gradient(135deg, #10b981 0%, #059669 100%)', color: 'white', border: 'none', borderRadius: '12px', cursor: 'pointer', fontSize: '16px', fontWeight: '600', alignSelf: 'flex-start' }}>
                  <Save size={16} />
                  Save Changes
                </button>
              </div>
            </form>
          ) : (
            <div>
              <h3 style={{ fontSize: '18px', fontWeight: '600', color: '#374151', marginBottom: '16px' }}>Current About Content</h3>
              <div style={{ background: '#f8fafc', padding: '20px', borderRadius: '12px', fontSize: '16px', color: '#1f2937', lineHeight: '1.7', whiteSpace: 'pre-wrap' }}>
                {aboutData.content}
              </div>
            </div>
          )}
        </div>
      )}
      
      {!aboutData && !loading && (
        <div style={{ background: 'white', borderRadius: '20px', padding: '32px', boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.1)', textAlign: 'center' }}>
          <p style={{ color: '#6b7280', fontSize: '18px' }}>No about data found. Please refresh or check your database.</p>
        </div>
      )}
    </div>
  );

  const renderContent = (): JSX.Element => {
    if (loading) {
      return (
        <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', height: '400px', flexDirection: 'column' }}>
          <div style={{ width: '48px', height: '48px', border: '4px solid #f3f4f6', borderTop: '4px solid #3b82f6', borderRadius: '50%', animation: 'spin 1s linear infinite', marginBottom: '16px' }} />
          <p style={{ color: '#6b7280', fontSize: '16px' }}>Loading...</p>
          <style>{`
            @keyframes spin {
              0% { transform: rotate(0deg); }
              100% { transform: rotate(360deg); }
            }
          `}</style>
        </div>
      );
    }

    switch (activeTab) {
      case 'dashboard':
        return <DashboardContent />;
      case 'contacts':
        return <ContactsContent />;
      case 'home':
        return <HomeEditContent />;
      case 'about':
        return <AboutEditContent />;
      case 'projects':
        return <AdminProjects />;
      default:
        return <DashboardContent />;
    }
  };

  return (
    <div style={{ display: 'flex', height: '100vh', backgroundColor: '#f9fafb' }}>
      <div style={{ width: '280px', background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)', color: 'white', display: 'flex', flexDirection: 'column', padding: '24px' }}>
        <div style={{ marginBottom: '40px' }}>
          <h1 style={{ fontSize: '24px', fontWeight: '700', marginBottom: '8px', background: 'rgba(255, 255, 255, 0.2)', padding: '12px 16px', borderRadius: '12px', textAlign: 'center' }}>
            🎯 Admin Panel
          </h1>
          <p style={{ fontSize: '14px', opacity: 0.8, textAlign: 'center' }}>Portfolio Management</p>
        </div>

        <nav style={{ flex: 1 }}>
          <ul style={{ listStyle: 'none', padding: 0, margin: 0 }}>
            {[
              { id: 'dashboard', icon: '📊', label: 'Dashboard' },
              { id: 'contacts', icon: '💬', label: 'Messages' },
              { id: 'home', icon: '🏠', label: 'Home Page' },
              { id: 'about', icon: '👤', label: 'About Page' },
              { id: 'projects', icon: '🚀', label: 'Projects' }
            ].map((item) => (
              <li key={item.id} style={{ marginBottom: '8px' }}>
                <button onClick={() => setActiveTab(item.id)} style={{ width: '100%', padding: '16px 20px', background: activeTab === item.id ? 'rgba(255, 255, 255, 0.2)' : 'transparent', color: 'white', border: 'none', borderRadius: '12px', cursor: 'pointer', fontSize: '16px', fontWeight: '500', textAlign: 'left', display: 'flex', alignItems: 'center', gap: '12px' }}>
                  <span style={{ fontSize: '18px' }}>{item.icon}</span>
                  {item.label}
                </button>
              </li>
            ))}
          </ul>
        </nav>

        <button onClick={handleLogout} style={{ padding: '16px 20px', background: 'rgba(239, 68, 68, 0.2)', color: 'white', border: '2px solid rgba(239, 68, 68, 0.3)', borderRadius: '12px', cursor: 'pointer', fontSize: '16px', fontWeight: '600', display: 'flex', alignItems: 'center', justifyContent: 'center', gap: '8px' }}>
          <LogOut size={20} />
          Logout
        </button>
      </div>

      <div style={{ flex: 1, overflow: 'auto', backgroundColor: '#f9fafb' }}>
        {renderContent()}
      </div>
    </div>
  );
};

export default AdminDashboard;