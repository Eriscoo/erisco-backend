package contact

import (
	"fmt"

	"github.com/eriscoo/blog-backend/internal/application"
	"github.com/eriscoo/blog-backend/internal/domain"
	"github.com/eriscoo/blog-backend/internal/infrastructure/email"
	"github.com/eriscoo/blog-backend/internal/infrastructure/turnstile"
)

type Service struct {
	repo       application.ContactRepository
	emailSvc   *email.Service
	turnstile  *turnstile.Service
	adminEmail string
}

func New(repo application.ContactRepository, emailSvc *email.Service, turnstile *turnstile.Service, adminEmail string) *Service {
	return &Service{
		repo:       repo,
		emailSvc:   emailSvc,
		turnstile:  turnstile,
		adminEmail: adminEmail,
	}
}

func (s *Service) Submit(msg *domain.ContactMessage, turnstileToken string) error {
	ok, err := s.turnstile.Verify(turnstileToken)
	if err != nil {
		return fmt.Errorf("turnstile: %w", err)
	}
	if !ok {
		return domain.ErrSpamDetected
	}

	if err := s.repo.Create(msg); err != nil {
		return err
	}

	go s.sendAdminEmail(msg)
	go s.sendAutoReply(msg)

	return nil
}

func (s *Service) sendAdminEmail(msg *domain.ContactMessage) {
	name := msg.Name
	if name == "" {
		name = "Anonymous"
	}
	phone := msg.Phone
	if phone == "" {
		phone = "-"
	}

	body := fmt.Sprintf(
		"From: %s <%s>\r\nPhone: %s\r\n\r\nMessage:\r\n%s",
		name, msg.Email, phone, msg.Message,
	)
	_ = s.emailSvc.Send(s.adminEmail, "[Eriscoo] "+msg.Subject, body)
}

func (s *Service) sendAutoReply(msg *domain.ContactMessage) {
	name := msg.Name
	if name == "" {
		name = "there"
	}

	body := fmt.Sprintf(
		"Hi %s,\r\n\r\n"+
			"Thank you for reaching out! I've received your message and will get back to you as soon as possible.\r\n\r\n"+
			"Here's a copy of your message:\r\n---\r\n%s\r\n---\r\n\r\n"+
			"Best regards,\r\n"+
			"Erisco Berto",
		name, msg.Message,
	)
	_ = s.emailSvc.Send(msg.Email, "[Eriscoo.com] Thank you for contacting me!", body)
}
