import React, { useState, useEffect } from "react";
import { apiService } from '../services/api';
import '../assets/styles/Project.scss';

interface Project {
  id: number;
  name: string;
  description: string;
  message: string;
  image_url?: string;
  technologies?: string;
  github_url?: string;
  demo_url?: string;
}

function Project() {
  const [projects, setProjects] = useState<Project[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string>('');

  useEffect(() => {
    const fetchProjects = async () => {
      try {
        setLoading(true);
        const data = await apiService.getProjects();
        setProjects(data);
      } catch (err) {
        console.error('Error fetching projects:', err);
        setError('Failed to load projects');
      } finally {
        setLoading(false);
      }
    };

    fetchProjects();
  }, []);

  if (loading) {
    return (
      <div className="projects-container" id="projects">
        <h1>Personal Projects</h1>
        <div style={{ textAlign: 'center', padding: '50px' }}>
          <p>Loading projects...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="projects-container" id="projects">
        <h1>Personal Projects</h1>
        <div style={{ textAlign: 'center', padding: '50px', color: 'red' }}>
          <p>Error: {error}</p>
        </div>
      </div>
    );
  }

  return (
    <div className="projects-container" id="projects">
      <h1>Personal Projects</h1>
      <div className="projects-grid">
        {projects.length > 0 ? (
          projects.map((project) => (
            <div key={project.id} className="project">
              {/* Project Image */}
              <div className="project-image">
                {project.image_url ? (
                  <img 
                    src={project.image_url} 
                    className="zoom" 
                    alt={project.name} 
                    width="100%"
                    style={{ 
                      height: '200px', 
                      objectFit: 'cover',
                      borderRadius: '8px 8px 0 0'
                    }}
                  />
                ) : (
                  <div 
                    style={{ 
                      height: '200px', 
                      backgroundColor: '#f0f0f0', 
                      display: 'flex', 
                      alignItems: 'center', 
                      justifyContent: 'center',
                      borderRadius: '8px 8px 0 0',
                      color: '#999'
                    }}
                  >
                    No Image
                  </div>
                )}
              </div>

              {/* Project Content */}
              <div style={{ padding: '20px' }}>
                {/* Project Title */}
                <h2 style={{ marginBottom: '10px' }}>{project.name}</h2>

                {/* Technologies */}
                {project.technologies && (
                  <div style={{ marginBottom: '15px' }}>
                    <div style={{ 
                      display: 'flex', 
                      flexWrap: 'wrap', 
                      gap: '5px' 
                    }}>
                      {project.technologies.split(',').map((tech, index) => (
                        <span 
                          key={index}
                          style={{
                            backgroundColor: '#007bff',
                            color: 'white',
                            padding: '2px 8px',
                            borderRadius: '12px',
                            fontSize: '12px',
                            fontWeight: 'bold'
                          }}
                        >
                          {tech.trim()}
                        </span>
                      ))}
                    </div>
                  </div>
                )}

                {/* Project Description */}
                <p style={{ 
                  marginBottom: '15px',
                  lineHeight: '1.6',
                  color: '#666'
                }}>
                  {project.description}
                </p>

                {/* Detailed Message */}
                <p style={{ 
                  marginBottom: '20px',
                  fontSize: '14px',
                  color: '#888',
                  fontStyle: 'italic'
                }}>
                  {project.message}
                </p>

                {/* Action Buttons */}
                <div style={{ 
                  display: 'flex', 
                  gap: '10px',
                  flexWrap: 'wrap'
                }}>
                  {project.github_url && (
                    <a 
                      href={project.github_url} 
                      target="_blank" 
                      rel="noreferrer"
                      style={{
                        backgroundColor: '#333',
                        color: 'white',
                        padding: '8px 16px',
                        textDecoration: 'none',
                        borderRadius: '5px',
                        fontSize: '14px',
                        fontWeight: 'bold',
                        display: 'inline-flex',
                        alignItems: 'center',
                        gap: '5px'
                      }}
                    >
                      🔗 GitHub
                    </a>
                  )}
                  
                  {project.demo_url && (
                    <a 
                      href={project.demo_url} 
                      target="_blank" 
                      rel="noreferrer"
                      style={{
                        backgroundColor: '#28a745',
                        color: 'white',
                        padding: '8px 16px',
                        textDecoration: 'none',
                        borderRadius: '5px',
                        fontSize: '14px',
                        fontWeight: 'bold',
                        display: 'inline-flex',
                        alignItems: 'center',
                        gap: '5px'
                      }}
                    >
                      🚀 Live Demo
                    </a>
                  )}
                </div>
              </div>
            </div>
          ))
        ) : (
          <div style={{ 
            textAlign: 'center', 
            padding: '50px',
            gridColumn: '1 / -1'
          }}>
            <p>No projects found. Add some projects through the admin panel!</p>
          </div>
        )}
      </div>
    </div>
  );
}

export default Project;