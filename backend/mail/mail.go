package mail

import (
	"fmt"
	"log"

	"github.com/resend/resend-go/v2"
)

// Mail configuration
const (
	RESEND_API_KEY = "re_gUY5KhWm_6HUw3WPe2ZjQd6yWGY4SzH42" // Sizin API key'iniz
	FROM_EMAIL     = "onboarding@resend.dev"                // Resend test email
	TO_EMAIL       = "dcnitro41@gmail.com"                  // Size gelecek maillerin adresi
)

// ContactMailData represents contact form data for email
type ContactMailData struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Message string `json:"message"`
}

// SendContactMail sends contact form submission via email
func SendContactMail(data ContactMailData) error {
	// Initialize Resend client
	client := resend.NewClient(RESEND_API_KEY)

	// Inline CSS HTML Template (Gmail compatible)
	htmlContent := `
		<div style="font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; line-height: 1.6; margin: 0; padding: 20px; background-color: #f4f4f4;">
			<div style="max-width: 600px; margin: 0 auto; background: white; border-radius: 10px; box-shadow: 0 0 20px rgba(0,0,0,0.1); overflow: hidden;">
				
				<!-- Header -->
				<div style="background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); color: white; padding: 30px 20px; text-align: center;">
					<h1 style="margin: 0; font-size: 24px; font-weight: 300;">📧 New Portfolio Contact</h1>
					<div style="display: inline-block; background: #28a745; color: white; padding: 5px 15px; border-radius: 20px; font-size: 12px; margin-top: 15px;">New Message Received</div>
				</div>
				
				<!-- Content -->
				<div style="padding: 30px;">
					
					<!-- Email Field -->
					<div style="margin-bottom: 20px; padding: 15px; background: #f8f9fa; border-left: 4px solid #667eea; border-radius: 5px;">
						<div style="display: flex; align-items: center; margin-bottom: 10px;">
							<div style="width: 20px; height: 20px; background: #667eea; border-radius: 50%%; display: inline-flex; align-items: center; justify-content: center; color: white; font-size: 12px; margin-right: 10px;">📧</div>
							<div>
								<div style="font-weight: bold; color: #333; font-size: 14px; text-transform: uppercase; letter-spacing: 1px; margin-bottom: 5px;">Email Address</div>
								<div style="color: #555; font-size: 16px;">%s</div>
							</div>
						</div>
					</div>

					<!-- Phone Field -->
					<div style="margin-bottom: 20px; padding: 15px; background: #f8f9fa; border-left: 4px solid #667eea; border-radius: 5px;">
						<div style="display: flex; align-items: center; margin-bottom: 10px;">
							<div style="width: 20px; height: 20px; background: #667eea; border-radius: 50%%; display: inline-flex; align-items: center; justify-content: center; color: white; font-size: 12px; margin-right: 10px;">📱</div>
							<div>
								<div style="font-weight: bold; color: #333; font-size: 14px; text-transform: uppercase; letter-spacing: 1px; margin-bottom: 5px;">Phone Number</div>
								<div style="color: #555; font-size: 16px;">%s</div>
							</div>
						</div>
					</div>

					<!-- Message Field -->
					<div style="margin-bottom: 20px; padding: 15px; background: #f8f9fa; border-left: 4px solid #667eea; border-radius: 5px;">
						<div style="font-weight: bold; color: #333; font-size: 14px; text-transform: uppercase; letter-spacing: 1px; margin-bottom: 10px;">💬 Message Content</div>
						<div style="background: #fff; border: 1px solid #e0e0e0; border-radius: 8px; padding: 20px; font-style: italic; color: #555;">%s</div>
					</div>

				</div>

				<!-- Footer -->
				<div style="background: #f8f9fa; padding: 20px; text-align: center; border-top: 1px solid #e0e0e0; color: #666; font-size: 12px;">
					<p style="margin: 5px 0;">🚀 This message was sent from your portfolio website</p>
					<p style="margin: 5px 0;">Generated automatically • Portfolio Contact System</p>
				</div>

			</div>
		</div>
	`

	textContent := fmt.Sprintf(`
		New Portfolio Contact Message
		
		Email: %s
		Phone: %s
		
		Message:
		%s
		
		---
		This message was sent from your portfolio website.
	`, data.Email, data.Phone, data.Message)

	// Prepare email parameters
	subject := "Portfolio Contact: Contact Form"
	params := &resend.SendEmailRequest{
		From:    FROM_EMAIL,
		To:      []string{TO_EMAIL},
		Subject: subject,
		Html:    htmlContent,
		Text:    textContent,
	}

	// Send email
	sent, err := client.Emails.Send(params)
	if err != nil {
		log.Printf("Mail sending error: %v", err)
		return fmt.Errorf("failed to send email: %v", err)
	}

	log.Printf("Mail sent successfully. ID: %s", sent.Id)
	return nil
}

// SendWelcomeMail sends welcome email to new users (bonus feature)
func SendWelcomeMail(userEmail, userName string) error {
	client := resend.NewClient(RESEND_API_KEY)

	subject := "Welcome to Portfolio Admin Panel!"

	htmlContent := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px; background-color: #f4f4f4;">
			<div style="background: white; padding: 30px; border-radius: 10px; box-shadow: 0 0 10px rgba(0,0,0,0.1);">
				<h2 style="color: #333; margin-bottom: 20px;">Welcome %s!</h2>
				<p style="color: #555; line-height: 1.6;">You have successfully registered to the portfolio admin panel.</p>
				<p style="color: #555; line-height: 1.6;">You can now manage your portfolio content.</p>
				<br>
				<p style="color: #667eea; font-weight: bold;">Happy coding!</p>
			</div>
		</div>
	`, userName)

	params := &resend.SendEmailRequest{
		From:    FROM_EMAIL,
		To:      []string{userEmail},
		Subject: subject,
		Html:    htmlContent,
	}

	sent, err := client.Emails.Send(params)
	if err != nil {
		log.Printf("Welcome mail error: %v", err)
		return err
	}

	log.Printf("Welcome mail sent successfully. ID: %s", sent.Id)
	return nil
}
