package data

import (
	"github.com/palavrapasse/damn/pkg/entity/subscribe"
)

type QueryAllSubscriptionsResult []subscribe.Subscription

type QuerySubscriptionsResult []QuerySubscriptionResult

type QuerySubscriptionResult struct {
	subscribe.Subscriber
	subscribe.Affected
}

type QuerySubscriptionWithoutAffectedResult struct {
	subscribe.Subscriber
}

func (qsr QuerySubscriptionsResult) ConvertToQueryAllSubscriptionsResult() QueryAllSubscriptionsResult {
	var r QueryAllSubscriptionsResult

	for _, v := range qsr {
		r = r.AddSubscription(v.Subscriber, v.Affected)
	}

	return r
}

func (qsar QuerySubscriptionWithoutAffectedResult) ConvertToSubscription() subscribe.Subscription {
	return subscribe.Subscription{
		Subscriber: qsar.Subscriber,
	}
}

func (qasr QueryAllSubscriptionsResult) AddSubscription(sub subscribe.Subscriber, aff subscribe.Affected) QueryAllSubscriptionsResult {

	for i, v := range qasr {

		if v.Subscriber.B64Email == sub.B64Email {
			qasr[i].Affected = append(qasr[i].Affected, aff)
			return qasr
		}
	}

	subscription := subscribe.Subscription{
		Subscriber: sub,
		Affected:   []subscribe.Affected{aff},
	}

	qasr = append(qasr, subscription)

	return qasr
}
