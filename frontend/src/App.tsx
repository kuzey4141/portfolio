import React, {useState, useEffect} from "react";
import {
  Main,
  Project,
  Contact,
  Navigation,
  Footer,
  AdminLogin,
} from "./components";
import About from './components/About';
import FadeIn from './components/FadeIn';
import './index.scss';

function App() {
    const [mode, setMode] = useState<string>('dark');
    const [showAdmin, setShowAdmin] = useState<boolean>(false);
    
    const handleModeChange = () => {
        if (mode === 'dark') {
            setMode('light');
        } else {
            setMode('dark');
        }
    }
    
    useEffect(() => {
        window.scrollTo({top: 0, left: 0, behavior: 'smooth'});
      }, []);
      
    // Check for admin URL
    useEffect(() => {
        if (window.location.pathname === '/admin') {
            setShowAdmin(true);
        }
    }, []);
    
    if (showAdmin) {
        return <AdminLogin />;
    }
    
    return (
    <div className={`main-container ${mode === 'dark' ? 'dark-mode' : 'light-mode'}`}>
        <Navigation parentToChild={{mode}} modeChange={handleModeChange}/>
        <FadeIn transitionDuration={700}>
            <Main/>
            <About/>
            <Project/>
            <Contact/>
        </FadeIn>
        <Footer />
    </div>
    );
}

export default App;