import React, { useState, useEffect } from 'react';
import { Upload, Github, ExternalLink, Edit, Trash2, Plus, Save, X, Image as ImageIcon } from 'lucide-react';
import { API_BASE_URL } from '../config'; 

interface Project {
  id?: number;
  name: string;
  description: string;
  message: string;
  image_url?: string;
  technologies?: string;
  github_url?: string;
  demo_url?: string;
}

const AdminProjects: React.FC = () => {
  const [projects, setProjects] = useState<Project[]>([]);
  const [currentProject, setCurrentProject] = useState<Project>({
    name: '',
    description: '',
    message: '',
    image_url: '',
    technologies: '',
    github_url: '',
    demo_url: ''
  });
  const [isEditing, setIsEditing] = useState<boolean>(false);
  const [showForm, setShowForm] = useState<boolean>(false);
  const [loading, setLoading] = useState<boolean>(false);
  const [message, setMessage] = useState<string>('');

  // API base'ini son slash'ı atarak güvene al
  const API = (API_BASE_URL || '').replace(/\/$/, '');

  const getToken = () => (typeof window !== 'undefined' ? localStorage.getItem('authToken') : null);

  const safeJson = async (res: Response) => {
    try {
      return await res.json();
    } catch {
      return null;
    }
  };

  const fetchProjects = async () => {
    try {
      const response = await fetch(`${API}/projects`, { credentials: 'include' });
      if (!response.ok) {
        const errData = await safeJson(response);
        throw new Error(errData?.error || `Fetch failed with ${response.status}`);
      }
      const data = await safeJson(response);
      console.log('Fetched projects:', data);
      setProjects(Array.isArray(data) ? data : []);
    } catch (error) {
      console.error('Failed to fetch projects:', error);
      setMessage('Error loading projects');
      setProjects([]);
    }
  };

  useEffect(() => {
    fetchProjects();
    
  }, []);

  const resetForm = () => {
    setCurrentProject({
      name: '',
      description: '',
      message: '',
      image_url: '',
      technologies: '',
      github_url: '',
      demo_url: ''
    });
    setIsEditing(false);
    setShowForm(false);
  };

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setMessage('');

    const token = getToken();
    if (!token) {
      setMessage('Please login first');
      setLoading(false);
      return;
    }

    try {
      const url = isEditing 
        ? `${API}/admin/projects/${currentProject.id}`
        : `${API}/admin/projects`;
      
      const method = isEditing ? 'PUT' : 'POST';

      const response = await fetch(url, {
        method,
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify(currentProject)
      });

      const result = await safeJson(response);
      
      if (response.ok) {
        setMessage(isEditing ? '✅ Project updated successfully!' : '✅ Project added successfully!');
        fetchProjects();
        resetForm();
        setTimeout(() => setMessage(''), 3000);
      } else {
        setMessage(`❌ Error: ${result?.error || `Operation failed (${response.status})`}`);
      }
    } catch (error) {
      console.error('Save error:', error);
      setMessage('❌ Error occurred while saving');
    } finally {
      setLoading(false);
    }
  };

  const handleEdit = (project: Project) => {
    setCurrentProject(project);
    setIsEditing(true);
    setShowForm(true);
  };

  const handleDelete = async (id: number) => {
    if (!window.confirm('🗑️ Are you sure you want to delete this project? This action cannot be undone.')) {
      return;
    }

    const token = getToken();
    if (!token) {
      setMessage('Please login first');
      return;
    }

    try {
      const response = await fetch(`${API}/admin/projects/${id}`, {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });

      const result = await safeJson(response);
      
      if (response.ok) {
        setMessage('✅ Project deleted successfully!');
        fetchProjects();
        setTimeout(() => setMessage(''), 3000);
      } else {
        setMessage(`❌ Error: ${result?.error || `Delete failed (${response.status})`}`);
      }
    } catch (error) {
      console.error('Delete error:', error);
      setMessage('❌ Error occurred while deleting');
    }
  };

  const handleImageUpload = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      if (file.size > 2 * 1024 * 1024) {
        setMessage('❌ Image size should be less than 2MB');
        return;
      }

      const reader = new FileReader();
      reader.onload = (event) => {
        const base64 = event.target?.result as string;
        setCurrentProject(prev => ({
          ...prev,
          image_url: base64
        }));
      };
      reader.readAsDataURL(file);
    }
  };

  return (
    <div style={{ padding: '24px', maxWidth: '1400px', margin: '0 auto' }}>
      {/* Header */}
      <div style={{ 
        display: 'flex', 
        justifyContent: 'space-between', 
        alignItems: 'center', 
        marginBottom: '32px',
        background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
        padding: '24px',
        borderRadius: '16px',
        color: 'white'
      }}>
        <div>
          <h1 style={{ margin: '0 0 8px 0', fontSize: '28px', fontWeight: '700' }}>🚀 Project Management</h1>
          <p style={{ margin: 0, opacity: 0.9 }}>Manage your portfolio projects with ease</p>
        </div>
        <button
          onClick={() => {
            resetForm();
            setShowForm(true);
          }}
          style={{
            display: 'flex',
            alignItems: 'center',
            gap: '8px',
            background: 'rgba(255, 255, 255, 0.2)',
            color: 'white',
            border: '1px solid rgba(255, 255, 255, 0.3)',
            borderRadius: '12px',
            padding: '12px 20px',
            fontSize: '14px',
            fontWeight: '600',
            cursor: 'pointer',
            transition: 'all 0.3s ease'
          }}
        >
          <Plus size={18} />
          Add New Project
        </button>
      </div>
      
      {/* Message Display */}
      {message && (
        <div style={{
          padding: '16px 20px',
          marginBottom: '24px',
          backgroundColor: message.includes('❌') ? '#fee2e2' : '#dcfce7',
          color: message.includes('❌') ? '#dc2626' : '#16a34a',
          border: `1px solid ${message.includes('❌') ? '#fecaca' : '#bbf7d0'}`,
          borderRadius: '12px',
          fontSize: '14px',
          fontWeight: '500',
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center'
        }}>
          {message}
          <button 
            onClick={() => setMessage('')}
            style={{
              background: 'none',
              border: 'none',
              color: 'inherit',
              cursor: 'pointer',
              fontSize: '18px'
            }}
          >
            ×
          </button>
        </div>
      )}

      {/* Project Form */}
      {showForm && (
        <div style={{ 
          background: 'white', 
          borderRadius: '16px', 
          padding: '32px', 
          marginBottom: '32px',
          boxShadow: '0 10px 25px -5px rgba(0, 0, 0, 0.1)',
          border: '1px solid #e5e7eb'
        }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '24px' }}>
            <h2 style={{ margin: 0, fontSize: '24px', fontWeight: '600', color: '#1f2937' }}>
              {isEditing ? '✏️ Edit Project' : '➕ Add New Project'}
            </h2>
            <button
              onClick={resetForm}
              style={{
                background: '#f3f4f6',
                border: 'none',
                borderRadius: '8px',
                padding: '8px',
                cursor: 'pointer',
                color: '#6b7280'
              }}
            >
              <X size={20} />
            </button>
          </div>
          
          <form onSubmit={handleSave}>
            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(300px, 1fr))', gap: '24px', marginBottom: '24px' }}>
              {/* Project Name */}
              <div>
                <label style={{ display: 'block', fontSize: '14px', fontWeight: '600', color: '#374151', marginBottom: '8px' }}>
                  Project Name *
                </label>
                <input
                  type="text"
                  value={currentProject.name}
                  onChange={(e) => setCurrentProject(prev => ({ ...prev, name: e.target.value }))}
                  required
                  placeholder="Enter project name"
                  style={{
                    width: '100%',
                    padding: '12px 16px',
                    border: '2px solid #e5e7eb',
                    borderRadius: '8px',
                    fontSize: '16px',
                    outline: 'none',
                    transition: 'border-color 0.3s ease',
                    boxSizing: 'border-box'
                  }}
                />
              </div>

              {/* Technologies */}
              <div>
                <label style={{ display: 'block', fontSize: '14px', fontWeight: '600', color: '#374151', marginBottom: '8px' }}>
                  Technologies
                </label>
                <input
                  type="text"
                  value={currentProject.technologies || ''}
                  onChange={(e) => setCurrentProject(prev => ({ ...prev, technologies: e.target.value }))}
                  placeholder="React, Node.js, PostgreSQL"
                  style={{
                    width: '100%',
                    padding: '12px 16px',
                    border: '2px solid #e5e7eb',
                    borderRadius: '8px',
                    fontSize: '16px',
                    outline: 'none',
                    boxSizing: 'border-box'
                  }}
                />
              </div>

              {/* GitHub URL */}
              <div>
                <label style={{ display: 'block', fontSize: '14px', fontWeight: '600', color: '#374151', marginBottom: '8px' }}>
                  GitHub URL
                </label>
                <input
                  type="url"
                  value={currentProject.github_url || ''}
                  onChange={(e) => setCurrentProject(prev => ({ ...prev, github_url: e.target.value }))}
                  placeholder="https://github.com/username/project"
                  style={{
                    width: '100%',
                    padding: '12px 16px',
                    border: '2px solid #e5e7eb',
                    borderRadius: '8px',
                    fontSize: '16px',
                    outline: 'none',
                    boxSizing: 'border-box'
                  }}
                />
              </div>

              {/* Demo URL */}
              <div>
                <label style={{ display: 'block', fontSize: '14px', fontWeight: '600', color: '#374151', marginBottom: '8px' }}>
                  Demo URL
                </label>
                <input
                  type="url"
                  value={currentProject.demo_url || ''}
                  onChange={(e) => setCurrentProject(prev => ({ ...prev, demo_url: e.target.value }))}
                  placeholder="https://project-demo.com"
                  style={{
                    width: '100%',
                    padding: '12px 16px',
                    border: '2px solid #e5e7eb',
                    borderRadius: '8px',
                    fontSize: '16px',
                    outline: 'none',
                    boxSizing: 'border-box'
                  }}
                />
              </div>
            </div>

            {/* Short Description */}
            <div style={{ marginBottom: '24px' }}>
              <label style={{ display: 'block', fontSize: '14px', fontWeight: '600', color: '#374151', marginBottom: '8px' }}>
                Short Description *
              </label>
              <input
                type="text"
                value={currentProject.description}
                onChange={(e) => setCurrentProject(prev => ({ ...prev, description: e.target.value }))}
                required
                placeholder="Brief description for project cards"
                style={{
                  width: '100%',
                  padding: '12px 16px',
                  border: '2px solid #e5e7eb',
                  borderRadius: '8px',
                  fontSize: '16px',
                  outline: 'none',
                  boxSizing: 'border-box'
                }}
              />
            </div>

            {/* Detailed Description */}
            <div style={{ marginBottom: '24px' }}>
              <label style={{ display: 'block', fontSize: '14px', fontWeight: '600', color: '#374151', marginBottom: '8px' }}>
                Detailed Description *
              </label>
              <textarea
                value={currentProject.message}
                onChange={(e) => setCurrentProject(prev => ({ ...prev, message: e.target.value }))}
                required
                rows={4}
                placeholder="Detailed project description"
                style={{
                  width: '100%',
                  padding: '12px 16px',
                  border: '2px solid #e5e7eb',
                  borderRadius: '8px',
                  fontSize: '16px',
                  outline: 'none',
                  resize: 'vertical',
                  fontFamily: 'inherit',
                  boxSizing: 'border-box'
                }}
              />
            </div>

            {/* Image Upload */}
            <div style={{ marginBottom: '32px' }}>
              <label style={{ display: 'block', fontSize: '14px', fontWeight: '600', color: '#374151', marginBottom: '8px' }}>
                Project Image
              </label>
              <div style={{
                border: '2px dashed #d1d5db',
                borderRadius: '12px',
                padding: '24px',
                textAlign: 'center',
                background: '#f9fafb'
              }}>
                <input
                  type="file"
                  accept="image/*"
                  onChange={handleImageUpload}
                  style={{ display: 'none' }}
                  id="image-upload"
                />
                <label 
                  htmlFor="image-upload" 
                  style={{
                    cursor: 'pointer',
                    display: 'flex',
                    flexDirection: 'column',
                    alignItems: 'center',
                    gap: '12px'
                  }}
                >
                  {currentProject.image_url ? (
                    <img 
                      src={currentProject.image_url} 
                      alt="Preview" 
                      style={{ 
                        maxWidth: '200px', 
                        maxHeight: '150px',
                        objectFit: 'cover',
                        borderRadius: '8px'
                      }} 
                    />
                  ) : (
                    <>
                      <ImageIcon size={48} style={{ color: '#9ca3af' }} />
                      <div>
                        <p style={{ margin: '0 0 4px 0', color: '#374151', fontWeight: '500' }}>Upload Project Image</p>
                        <p style={{ margin: 0, color: '#6b7280', fontSize: '14px' }}>PNG, JPG up to 2MB</p>
                      </div>
                    </>
                  )}
                </label>
              </div>
            </div>

            <div style={{ display: 'flex', gap: '12px' }}>
              <button
                type="submit"
                disabled={loading}
                style={{
                  display: 'flex',
                  alignItems: 'center',
                  gap: '8px',
                  background: loading ? '#9ca3af' : 'linear-gradient(135deg, #10b981 0%, #059669 100%)',
                  color: 'white',
                  border: 'none',
                  borderRadius: '12px',
                  padding: '12px 24px',
                  fontSize: '16px',
                  fontWeight: '600',
                  cursor: loading ? 'not-allowed' : 'pointer'
                }}
              >
                <Save size={16} />
                {loading ? 'Saving...' : (isEditing ? 'Update Project' : 'Add Project')}
              </button>
              
              <button
                type="button"
                onClick={resetForm}
                style={{
                  background: '#6b7280',
                  color: 'white',
                  border: 'none',
                  borderRadius: '12px',
                  padding: '12px 24px',
                  fontSize: '16px',
                  fontWeight: '600',
                  cursor: 'pointer'
                }}
              >
                Cancel
              </button>
            </div>
          </form>
        </div>
      )}

      {/* Projects Grid */}
      <div>
        <h2 style={{ fontSize: '24px', fontWeight: '600', color: '#1f2937', marginBottom: '24px' }}>
          Current Projects ({projects.length})
        </h2>
        
        {projects.length === 0 ? (
          <div style={{ 
            textAlign: 'center', 
            padding: '60px 20px',
            background: 'white',
            borderRadius: '16px',
            border: '2px dashed #d1d5db'
          }}>
            <p style={{ fontSize: '18px', color: '#6b7280', margin: 0 }}>
              No projects yet. Add your first project! 🚀
            </p>
          </div>
        ) : (
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(350px, 1fr))', gap: '24px' }}>
            {projects.map((project) => (
              <div
                key={project.id}
                style={{
                  background: 'white',
                  border: '1px solid #e5e7eb',
                  borderRadius: '16px',
                  overflow: 'hidden',
                  boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.1)',
                  transition: 'transform 0.3s ease, box-shadow 0.3s ease'
                }}
                onMouseEnter={(e) => {
                  e.currentTarget.style.transform = 'translateY(-4px)';
                  e.currentTarget.style.boxShadow = '0 10px 25px -5px rgba(0, 0, 0, 0.1)';
                }}
                onMouseLeave={(e) => {
                  e.currentTarget.style.transform = 'translateY(0)';
                  e.currentTarget.style.boxShadow = '0 4px 6px -1px rgba(0, 0, 0, 0.1)';
                }}
              >
                {project.image_url && (
                  <img 
                    src={project.image_url} 
                    alt={project.name}
                    style={{ 
                      width: '100%', 
                      height: '200px', 
                      objectFit: 'cover'
                    }} 
                  />
                )}
                
                <div style={{ padding: '24px' }}>
                  <h3 style={{ fontSize: '20px', fontWeight: '700', color: '#1f2937', marginBottom: '12px' }}>
                    {project.name}
                  </h3>
                  
                  <p style={{ color: '#6b7280', marginBottom: '16px', lineHeight: '1.5' }}>
                    {project.description}
                  </p>
                  
                  {project.technologies && (
                    <div style={{ marginBottom: '16px' }}>
                      <div style={{ display: 'flex', flexWrap: 'wrap', gap: '6px' }}>
                        {project.technologies.split(',').map((tech, index) => (
                          <span 
                            key={index}
                            style={{
                              background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
                              color: 'white',
                              padding: '4px 8px',
                              borderRadius: '6px',
                              fontSize: '12px',
                              fontWeight: '500'
                            }}
                          >
                            {tech.trim()}
                          </span>
                        ))}
                      </div>
                    </div>
                  )}
                  
                  <div style={{ display: 'flex', gap: '8px', marginBottom: '20px' }}>
                    {project.github_url && (
                      <a 
                        href={project.github_url} 
                        target="_blank" 
                        rel="noopener noreferrer"
                        style={{
                          display: 'flex',
                          alignItems: 'center',
                          gap: '6px',
                          background: '#333',
                          color: 'white',
                          padding: '6px 12px',
                          borderRadius: '8px',
                          fontSize: '14px',
                          textDecoration: 'none',
                          fontWeight: '500'
                        }}
                      >
                        <Github size={16} />
                        Code
                      </a>
                    )}
                    {project.demo_url && (
                      <a 
                        href={project.demo_url} 
                        target="_blank" 
                        rel="noopener noreferrer"
                        style={{
                          display: 'flex',
                          alignItems: 'center',
                          gap: '6px',
                          background: '#10b981',
                          color: 'white',
                          padding: '6px 12px',
                          borderRadius: '8px',
                          fontSize: '14px',
                          textDecoration: 'none',
                          fontWeight: '500'
                        }}
                      >
                        <ExternalLink size={16} />
                        Demo
                      </a>
                    )}
                  </div>
                  
                  <div style={{ display: 'flex', gap: '12px' }}>
                    <button
                      onClick={() => handleEdit(project)}
                      style={{
                        display: 'flex',
                        alignItems: 'center',
                        gap: '6px',
                        background: '#3b82f6',
                        color: 'white',
                        border: 'none',
                        borderRadius: '8px',
                        padding: '8px 16px',
                        fontSize: '14px',
                        fontWeight: '500',
                        cursor: 'pointer'
                      }}
                    >
                      <Edit size={16} />
                      Edit
                    </button>
                    <button
                      onClick={() => handleDelete(project.id!)}
                      style={{
                        display: 'flex',
                        alignItems: 'center',
                        gap: '6px',
                        background: '#ef4444',
                        color: 'white',
                        border: 'none',
                        borderRadius: '8px',
                        padding: '8px 16px',
                        fontSize: '14px',
                        fontWeight: '500',
                        cursor: 'pointer'
                      }}
                    >
                      <Trash2 size={16} />
                      Delete
                    </button>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
};

export default AdminProjects;
