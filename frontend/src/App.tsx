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
        const checkAdminRoute = () => {
            if (window.location.pathname === '/admin' || window.location.hash === '#admin') {
                setShowAdmin(true);
            }
        };
        
        checkAdminRoute();
        
        // Listen for hash changes
        const handleHashChange = () => {
            checkAdminRoute();
        };
        
        window.addEventListener('hashchange', handleHashChange);
        window.addEventListener('popstate', handleHashChange);
        
        return () => {
            window.removeEventListener('hashchange', handleHashChange);
            window.removeEventListener('popstate', handleHashChange);
        };
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