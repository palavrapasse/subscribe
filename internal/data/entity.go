package data

import (
	"github.com/palavrapasse/damn/pkg/entity"
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

type QueryLeaksResult []entity.HSHA256

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

func (qasr QueryAllSubscriptionsResult) GetAffectUsers() []subscribe.Affected {

	aff := []subscribe.Affected{}

	for _, v := range qasr {

		if v.Affected != nil {
			aff = append(aff, v.Affected...)
		}
	}

	return aff
}

func (qasr QueryAllSubscriptionsResult) GetSubscriptionsOfAffectUsers(affectedByLeak []entity.HSHA256) []subscribe.Subscription {
	alreadyAdded := make(map[subscribe.Subscriber]bool)
	sub := []subscribe.Subscription{}

	for _, affEmail := range affectedByLeak {
		for _, v := range qasr {
			if !alreadyAdded[v.Subscriber] && containsEmail(v.Affected, affEmail) {
				sub = append(sub, v)
				alreadyAdded[v.Subscriber] = true
			}
		}
	}

	return sub
}

func (qasr QueryAllSubscriptionsResult) GetSubscriptionsToAllLeaks() []subscribe.Subscription {

	sub := []subscribe.Subscription{}

	for _, v := range qasr {

		if len(v.Affected) == 0 {
			sub = append(sub, v)
		}
	}

	return sub
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

func containsEmail(affected []subscribe.Affected, email entity.HSHA256) bool {

	for _, v := range affected {

		if v.HSHA256Email == email {
			return true
		}
	}

	return false
}
