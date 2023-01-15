package data

import (
	"github.com/palavrapasse/damn/pkg/entity/subscribe"
)

type SubscriptionRequest struct {
	NotifyEmail    string   `json:"notifyEmail"`
	AffectedEmails []string `json:"affectedEmails"`
}

func SubscriptionRequestToSubscription(request SubscriptionRequest) subscribe.Subscription {
	sub := subscribe.NewSubscriber(request.NotifyEmail)

	affectedEmails := request.AffectedEmails
	laff := len(affectedEmails)
	aff := make([]subscribe.Affected, laff)

	for i := 0; i < laff; i++ {
		aff[i] = subscribe.NewAffected(affectedEmails[i])
	}

	return subscribe.Subscription{
		Subscriber: sub,
		Affected:   aff,
	}
}
