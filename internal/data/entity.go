package data

import (
	"github.com/palavrapasse/damn/pkg/entity"
	"github.com/palavrapasse/damn/pkg/entity/subscribe"
)

type Subscriptions []subscribe.Subscription

type QuerySubscriptionsResult []QuerySubscriptionResult

type QuerySubscriptionResult struct {
	subscribe.Subscriber
	subscribe.Affected
}

type QuerySubscriptionWithoutAffectedResult struct {
	subscribe.Subscriber
	subscribe.Affected
}

type QueryLeaksResult []entity.HSHA256

func (qsr QuerySubscriptionsResult) ConvertToSubscriptions() Subscriptions {
	var r Subscriptions

	for _, v := range qsr {
		r = r.addSubscription(v.Subscriber, v.Affected)
	}

	return r
}

func (qasr Subscriptions) addSubscription(sub subscribe.Subscriber, aff subscribe.Affected) Subscriptions {

	for i, v := range qasr {

		if v.Subscriber.B64Email == sub.B64Email {
			qasr[i].Affected = append(qasr[i].Affected, aff)
			return qasr
		}
	}

	subscription := subscribe.Subscription{
		Subscriber: sub,
	}

	if len(aff.HSHA256Email) != 0 {
		subscription.Affected = []subscribe.Affected{aff}
	}

	qasr = append(qasr, subscription)

	return qasr
}

func (qsr QuerySubscriptionsResult) GetAffectUsers() []subscribe.Affected {

	aff := []subscribe.Affected{}

	for _, v := range qsr {

		if len(v.Affected.HSHA256Email) != 0 {
			aff = append(aff, v.Affected)
		}
	}

	return aff
}

func (qsr QuerySubscriptionsResult) RemoveNotAffected(usersAffectedByLeak QueryLeaksResult) QuerySubscriptionsResult {

	alreadyAdded := make(map[QuerySubscriptionResult]bool)
	result := QuerySubscriptionsResult{}

	for _, userAffected := range usersAffectedByLeak {
		for _, sub := range qsr {

			if !alreadyAdded[sub] && sub.Affected.HSHA256Email == userAffected {
				result = append(result, sub)
				alreadyAdded[sub] = true
			}
		}
	}

	return result
}
