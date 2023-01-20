package data

import (
	"github.com/palavrapasse/damn/pkg/entity"
	"github.com/palavrapasse/damn/pkg/entity/query"
	"github.com/palavrapasse/damn/pkg/entity/subscribe"
)

type EmailInfo struct {
	UsersAffected     AllAffectedsInfo
	Leak              query.Leak
	PlatformsAffected query.Platform
}

type AllAffectedsInfo []AffectedInfo

type AffectedInfo struct {
	DestinationB64Email entity.Base64
	AffectedsEmail      []string
}

type QuerySubscriptionsResult []QuerySubscriptionResult

type QuerySubscriptionResult struct {
	subscribe.Subscriber
	subscribe.Affected
}

type QuerySubscriptionWithoutAffectedResult struct {
	subscribe.Subscriber
	subscribe.Affected
}

type QueryAffectedByLeakResult struct {
	entity.HSHA256
	Email string
}

type QueryLeakByIdResult struct {
	query.Leak
	query.Platform
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

func (qsr QuerySubscriptionsResult) GetAffectedsInfo(usersAffectedByLeak []QueryAffectedByLeakResult) AllAffectedsInfo {

	alreadyAdded := make(map[QuerySubscriptionResult]bool)
	result := AllAffectedsInfo{}

	for _, userAffected := range usersAffectedByLeak {
		for _, sub := range qsr {

			if !alreadyAdded[sub] && sub.Affected.HSHA256Email == userAffected.HSHA256 {
				result = result.addAffectedInfo(sub.Subscriber, userAffected.Email)
				alreadyAdded[sub] = true
			}
		}
	}

	return result
}

func (aainfo AllAffectedsInfo) addAffectedInfo(sub subscribe.Subscriber, aff string) AllAffectedsInfo {

	for i, v := range aainfo {

		if v.DestinationB64Email == sub.B64Email {
			aainfo[i].AffectedsEmail = append(aainfo[i].AffectedsEmail, aff)
			return aainfo
		}
	}

	affInfo := AffectedInfo{
		DestinationB64Email: sub.B64Email,
	}

	if len(aff) != 0 {
		affInfo.AffectedsEmail = []string{aff}
	}

	aainfo = append(aainfo, affInfo)

	return aainfo
}
