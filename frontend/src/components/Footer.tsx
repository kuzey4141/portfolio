import React from "react";
import GitHubIcon from '@mui/icons-material/GitHub';
import LinkedInIcon from '@mui/icons-material/LinkedIn';
import '../assets/styles/Footer.scss'

function Footer() {
  return (
    <footer>
      <div>
        <a href="https://github.com/kuzey4141" target="_blank" rel="noreferrer"><GitHubIcon/></a>
        <a href="" target="_blank" rel="noreferrer"><LinkedInIcon/></a>
      </div>
      <p>A portfolio designed & built by <a href="https://github.com/kuzey4141" target="_blank" rel="noreferrer">Kuzey</a> </p>
    </footer>
  );
}

export default Footer;