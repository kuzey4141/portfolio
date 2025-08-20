import React, { useState, useEffect } from "react";
import GitHubIcon from '@mui/icons-material/GitHub';
import LinkedInIcon from '@mui/icons-material/LinkedIn';
import '../assets/styles/Main.scss';
import { apiService, Home } from '../services/api';
import profileImage from '../assets/images/bugrakuzey.png';

function Main() {
  const [homeData, setHomeData] = useState<Home | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string>('');

  useEffect(() => {
    const fetchHomeData = async () => {
      try {
        setLoading(true);
        const data = await apiService.getHome();
        
        // Backend'den array gelir, ilk elemanı al
        if (data && data.length > 0) {
          setHomeData(data[0]);
        } else {
          setError('No home data found');
        }
      } catch (err) {
        console.error('Error fetching home data:', err);
        setError('Failed to load home data');
      } finally {
        setLoading(false);
      }
    };

    fetchHomeData();
  }, []);

  if (loading) {
    return (
      <div className="container" id="home">
        <div className="about-section">
          <div style={{ textAlign: 'center', padding: '50px' }}>
            <p>Loading...</p>
          </div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="container" id="home">
        <div className="about-section">
          <div style={{ textAlign: 'center', padding: '50px', color: 'red' }}>
            <p>Error: {error}</p>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="container" id="home">
      <div className="about-section">
        <div className="image-wrapper">
          <img 
            src={profileImage}
            alt="Avatar" 
            style={{
              width: '200px',
              height: '200px',
              borderRadius: '50%',
              objectFit: 'cover'
            }}
          />
        </div>
        <div className="content">
          <div className="social_icons">
            <a href="https://github.com/kuzey4141" target="_blank" rel="noreferrer"><GitHubIcon/></a>
            <a href="https://www.linkedin.com/in/yujisato/" target="_blank" rel="noreferrer"><LinkedInIcon/></a>
          </div>
          <h1>{homeData?.title || 'Loading...'}</h1>
          <p>{homeData?.description || 'Loading...'}</p>
          <div className="mobile_social_icons">
            <a href="https://github.com/kuzey4141" target="_blank" rel="noreferrer"><GitHubIcon/></a>
            <a href="https://www.linkedin.com/in/yujisato/" target="_blank" rel="noreferrer"><LinkedInIcon/></a>
          </div>
        </div>
      </div>
    </div>
  );
}

export default Main;