package pricealert

import (
	"context"
	"log"
)

type MailSender interface {
	SendPriceAlert(ctx context.Context, email string, productTitle string, oldPrice string, newPrice string, currency string) error
}

type EmailResolver interface {
	FindEmailByUserID(ctx context.Context, userID string) (string, error)
}

type Notifier struct {
	finder        Finder
	mailSender    MailSender
	emailResolver EmailResolver
}

func NewNotifier(finder Finder, mailSender MailSender, emailResolver EmailResolver) *Notifier {
	return &Notifier{
		finder:        finder,
		mailSender:    mailSender,
		emailResolver: emailResolver,
	}
}

func (n *Notifier) NotifyPriceChange(ctx context.Context, productID string, productTitle string, oldPrice string, newPrice string, currency string) {
	alerts, err := n.finder.ListByProductID(ctx, productID)
	if err != nil {
		log.Printf("failed to list alerts for product %s: %v", productID, err)
		return
	}

	for _, alert := range alerts {
		if alert.Channel != "EMAIL" {
			continue
		}

		email, err := n.emailResolver.FindEmailByUserID(ctx, alert.UserID)
		if err != nil {
			log.Printf("failed to resolve email for user %s: %v", alert.UserID, err)
			continue
		}

		if err := n.mailSender.SendPriceAlert(ctx, email, productTitle, oldPrice, newPrice, currency); err != nil {
			log.Printf("failed to send price alert to %s for product %s: %v", email, productID, err)
			continue
		}
	}
}
