package data

import (
	"github.com/palavrapasse/damn/pkg/database"
	"github.com/palavrapasse/damn/pkg/entity/subscribe"
)

func StoreSubscriptionDB(dbctx database.DatabaseContext[database.Record], subscription subscribe.Subscription) error {

	return dbctx.InsertSubscription(subscription)
}
