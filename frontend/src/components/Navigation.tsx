import React, { useEffect, useState } from "react";
import AppBar from '@mui/material/AppBar';
import Box from '@mui/material/Box';
import Button from '@mui/material/Button';
import CssBaseline from '@mui/material/CssBaseline';
import DarkModeIcon from '@mui/icons-material/DarkMode';
import Divider from '@mui/material/Divider';
import Drawer from '@mui/material/Drawer';
import IconButton from '@mui/material/IconButton';
import LightModeIcon from '@mui/icons-material/LightMode';
import List from '@mui/material/List';
import ListIcon from '@mui/icons-material/List';
import ListItem from '@mui/material/ListItem';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemText from '@mui/material/ListItemText';
import MenuIcon from '@mui/icons-material/Menu';
import Toolbar from '@mui/material/Toolbar';
import AdminPanelSettingsIcon from '@mui/icons-material/AdminPanelSettings';

const drawerWidth = 240;
const navItems = [['Home', 'home'], ['About', 'about'], ['Projects', 'projects'], ['Contact', 'contact']];

function Navigation({parentToChild, modeChange}: any) {
  const {mode} = parentToChild;
  const [mobileOpen, setMobileOpen] = useState<boolean>(false);
  const [scrolled, setScrolled] = useState<boolean>(false);

  const handleDrawerToggle = () => {
    setMobileOpen((prevState) => !prevState);
  };

  const handleAdminClick = () => {
    window.open('/admin', '_blank');
  };

  useEffect(() => {
    const handleScroll = () => {
      const navbar = document.getElementById("navigation");
      if (navbar) {
        const scrolled = window.scrollY > navbar.clientHeight;
        setScrolled(scrolled);
      }
    };

    window.addEventListener('scroll', handleScroll);

    return () => {
      window.removeEventListener('scroll', handleScroll);
    };
  }, []);

  const scrollToSection = (section: string) => {
    console.log(section)
    const element = document.getElementById(section);
    if (element) {
      element.scrollIntoView({ behavior: 'smooth' });
      console.log('Scrolling to:', element);
    } else {
      console.error(`Element with id "${section}" not found`);
    }
  };

  const drawer = (
    <Box className="navigation-bar-responsive" onClick={handleDrawerToggle} sx={{ textAlign: 'center' }}>
      <p className="mobile-menu-top"><ListIcon/>Menu</p>
      <Divider />
      <List>
        {navItems.map((item) => (
          <ListItem key={item[0]} disablePadding>
            <ListItemButton sx={{ textAlign: 'center' }} onClick={() => scrollToSection(item[1])}>
              <ListItemText primary={item[0]} />
            </ListItemButton>
          </ListItem>
        ))}
        {/* Mobile Admin Button */}
        <ListItem disablePadding>
          <ListItemButton sx={{ textAlign: 'center' }} onClick={handleAdminClick}>
            <ListItemText primary="Admin" />
          </ListItemButton>
        </ListItem>
      </List>
    </Box>
  );

  return (
    <Box sx={{ display: 'flex' }}>
      <CssBaseline />
      <AppBar component="nav" id="navigation" className={`navbar-fixed-top${scrolled ? ' scrolled' : ''}`}>
        <Toolbar className='navigation-bar' sx={{ justifyContent: 'center', position: 'relative' }}>
          {/* Mobile Menu Button - Left */}
          <IconButton
            color="inherit"
            aria-label="open drawer"
            edge="start"
            onClick={handleDrawerToggle}
            sx={{ 
              position: 'absolute',
              left: 0,
              display: { sm: 'none' }
            }}
          >
            <MenuIcon />
          </IconButton>
          
          {/* Dark/Light Mode Toggle - Left Side */}
          <Box sx={{ 
            position: 'absolute', 
            left: { xs: '50px', sm: '20px' },
            display: 'flex',
            alignItems: 'center'
          }}>
            {mode === 'dark' ? (
              <LightModeIcon 
                onClick={() => modeChange()}
                sx={{ cursor: 'pointer' }}
              />
            ) : (
              <DarkModeIcon 
                onClick={() => modeChange()}
                sx={{ cursor: 'pointer' }}
              />
            )}
          </Box>
          
          {/* Center Navigation */}
          <Box sx={{ 
            display: { xs: 'none', sm: 'flex' },
            alignItems: 'center',
            gap: 1
          }}>
            {navItems.map((item) => (
              <Button 
                key={item[0]} 
                onClick={() => scrollToSection(item[1])} 
                sx={{ 
                  color: '#fff',
                  fontSize: '0.95rem',
                  fontWeight: '500',
                  textTransform: 'none',
                  padding: '8px 16px',
                  borderRadius: '8px',
                  transition: 'all 0.3s ease',
                  '&:hover': {
                    backgroundColor: 'rgba(255, 255, 255, 0.1)',
                    transform: 'translateY(-1px)'
                  }
                }}
              >
                {item[0]}
              </Button>
            ))}
          </Box>

          {/* Admin Button - Right */}
          <Box sx={{ 
            position: 'absolute', 
            right: '20px',
            display: { xs: 'none', sm: 'block' }
          }}>
            <Button
              onClick={handleAdminClick}
              sx={{
                color: '#fff',
                border: '1px solid rgba(255, 255, 255, 0.3)',
                borderRadius: '20px',
                padding: '6px 16px',
                fontSize: '0.875rem',
                fontWeight: '500',
                textTransform: 'none',
                background: 'linear-gradient(45deg, rgba(255, 255, 255, 0.1), rgba(255, 255, 255, 0.05))',
                backdropFilter: 'blur(10px)',
                transition: 'all 0.3s ease',
                '&:hover': {
                  background: 'linear-gradient(45deg, rgba(255, 255, 255, 0.2), rgba(255, 255, 255, 0.1))',
                  transform: 'translateY(-1px)',
                  boxShadow: '0 4px 12px rgba(0, 0, 0, 0.15)',
                  borderColor: 'rgba(255, 255, 255, 0.5)'
                }
              }}
              startIcon={<AdminPanelSettingsIcon sx={{ fontSize: '18px' }} />}
            >
              Admin
            </Button>
          </Box>
        </Toolbar>
      </AppBar>
      <nav>
        <Drawer
          variant="temporary"
          open={mobileOpen}
          onClose={handleDrawerToggle}
          ModalProps={{
            keepMounted: true,
          }}
          sx={{
            display: { xs: 'block', sm: 'none' },
            '& .MuiDrawer-paper': { boxSizing: 'border-box', width: drawerWidth },
          }}
        >
          {drawer}
        </Drawer>
      </nav>
    </Box>
  );
}

export default Navigation;