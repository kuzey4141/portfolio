import React, { useState, useEffect } from "react";
import { apiService, About as AboutType } from '../services/api';

function About() {
  const [aboutData, setAboutData] = useState<AboutType | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string>('');

  useEffect(() => {
    const fetchAboutData = async () => {
      try {
        setLoading(true);
        const data = await apiService.getAbout();
        
        // Backend'den array gelir, ilk elemanı al
        if (data && data.length > 0) {
          setAboutData(data[0]);
        } else {
          setError('No about data found');
        }
      } catch (err) {
        console.error('Error fetching about data:', err);
        setError('Failed to load about data');
      } finally {
        setLoading(false);
      }
    };

    fetchAboutData();
  }, []);

  return (
    <section id="about">
      <div className="container" style={{ 
        maxWidth: '800px', 
        margin: '0 auto', 
        padding: '60px 20px',
        textAlign: 'center'
      }}>
        <h2 style={{ 
          fontSize: '2.5rem', 
          marginBottom: '30px',
          color: '#333'
        }}>
          About Me
        </h2>
        
        {loading && (
          <div style={{ padding: '40px' }}>
            <p style={{ fontSize: '1.1rem', color: '#666' }}>Loading about information...</p>
          </div>
        )}
        
        {error && (
          <div style={{ padding: '40px' }}>
            <p style={{ fontSize: '1.1rem', color: '#e74c3c' }}>Error: {error}</p>
          </div>
        )}
        
        {aboutData && (
          <div style={{ 
            fontSize: '1.2rem', 
            lineHeight: '1.8', 
            color: '#555',
            textAlign: 'left',
            maxWidth: '600px',
            margin: '0 auto'
          }}>
            <p>{aboutData.content}</p>
          </div>
        )}
        
        {!loading && !error && !aboutData && (
          <div style={{ padding: '40px' }}>
            <p style={{ fontSize: '1.1rem', color: '#666' }}>No about information available.</p>
          </div>
        )}
      </div>
    </section>
  );
}

export default About;