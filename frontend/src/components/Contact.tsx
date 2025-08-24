import React, { useState } from 'react';
import '../assets/styles/Contact.scss';
import { apiService } from '../services/api';

function Contact() {
  const [name, setName] = useState<string>('');
  const [email, setEmail] = useState<string>('');
  const [phone, setPhone] = useState<string>('');
  const [message, setMessage] = useState<string>('');
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [successMessage, setSuccessMessage] = useState<string>('');
  const [errorMessage, setErrorMessage] = useState<string>('');

  const sendEmail = async (e: any) => {
    e.preventDefault();
    
    // Clear previous messages
    setSuccessMessage('');
    setErrorMessage('');
    
    // Validation
    if (!name.trim() || !email.trim() || !phone.trim() || !message.trim()) {
      setErrorMessage('❌ Please fill all fields!');
      return;
    }
    
    setIsLoading(true);
    
    try {
      const contactData = {
        name: name.trim(),
        email: email.trim(),
        phone: phone.trim(),
        message: message.trim()
      };
      
      const response = await apiService.sendContact(contactData);
      
      if (response.message) {
        setSuccessMessage('✅ Message sent successfully! I\'ll get back to you soon.');
        // Clear form
        setName('');
        setEmail('');
        setPhone('');
        setMessage('');
      } else {
        setErrorMessage('❌ Failed to send message. Please try again.');
      }
    } catch (error) {
      console.error('Contact form error:', error);
      setErrorMessage('❌ Failed to send message. Please check your connection and try again.');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div id="contact">
      <div className="items-container">
        <div className="contact_wrapper">
          <h1>Contact Me</h1>
          
          
          {/* Success/Error Messages */}
          {successMessage && (
            <div style={{ 
              backgroundColor: '#d4edda', 
              color: '#155724', 
              padding: '15px', 
              borderRadius: '5px', 
              marginBottom: '20px',
              border: '1px solid #c3e6cb'
            }}>
              {successMessage}
            </div>
          )}
          {errorMessage && (
            <div style={{ 
              backgroundColor: '#f8d7da', 
              color: '#721c24', 
              padding: '15px', 
              borderRadius: '5px', 
              marginBottom: '20px',
              border: '1px solid #f5c6cb'
            }}>
              {errorMessage}
            </div>
          )}
          
          <form className='contact-form' onSubmit={sendEmail}>
            <div className='form-flex' style={{ display: 'flex', gap: '20px', marginBottom: '20px' }}>
              <input
                type="text"
                placeholder="Your Name *"
                value={name}
                onChange={(e) => setName(e.target.value)}
                disabled={isLoading}
                required
                style={{
                  flex: 1,
                  padding: '15px',
                  fontSize: '16px',
                  border: '1px solid #ccc',
                  borderRadius: '5px',
                  backgroundColor: isLoading ? '#f5f5f5' : 'white'
                }}
              />
              <input
                type="email"
                placeholder="Email *"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                disabled={isLoading}
                required
                style={{
                  flex: 1,
                  padding: '15px',
                  fontSize: '16px',
                  border: '1px solid #ccc',
                  borderRadius: '5px',
                  backgroundColor: isLoading ? '#f5f5f5' : 'white'
                }}
              />
            </div>
            
            <div style={{ marginBottom: '20px' }}>
              <input
                type="tel"
                placeholder="Phone *"
                value={phone}
                onChange={(e) => setPhone(e.target.value)}
                disabled={isLoading}
                required
                style={{
                  width: '100%',
                  padding: '15px',
                  fontSize: '16px',
                  border: '1px solid #ccc',
                  borderRadius: '5px',
                  backgroundColor: isLoading ? '#f5f5f5' : 'white',
                  boxSizing: 'border-box'
                }}
              />
            </div>
            
            <div style={{ marginBottom: '20px' }}>
              <textarea
                placeholder="Message *"
                value={message}
                onChange={(e) => setMessage(e.target.value)}
                disabled={isLoading}
                required
                rows={8}
                style={{
                  width: '100%',
                  padding: '15px',
                  fontSize: '16px',
                  border: '1px solid #ccc',
                  borderRadius: '5px',
                  backgroundColor: isLoading ? '#f5f5f5' : 'white',
                  boxSizing: 'border-box',
                  resize: 'vertical',
                  fontFamily: 'inherit'
                }}
              />
            </div>
            
            <button 
              type="submit"
              disabled={isLoading}
              style={{
                backgroundColor: isLoading ? '#6c757d' : '#007bff',
                color: 'white',
                padding: '12px 30px',
                fontSize: '16px',
                border: 'none',
                borderRadius: '5px',
                cursor: isLoading ? 'not-allowed' : 'pointer',
                display: 'flex',
                alignItems: 'center',
                gap: '10px'
              }}
            >
              {isLoading ? 'Sending...' : 'Send ➤'}
            </button>
          </form>
        </div>
      </div>
    </div>
  );
}

export default Contact;