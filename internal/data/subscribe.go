package data

import (
	"github.com/palavrapasse/damn/pkg/database"
	"github.com/palavrapasse/damn/pkg/entity/subscribe"
)

const subscriptionsWithoutAffectedQuery = `
SELECT S.b64email
 FROM  Subscriber S
 WHERE S.subid not in
    (SELECT SA.subid
    FROM SubscriberAffected SA) 
`

const subscriptionsQuery = `
SELECT A.hsha256email, S.b64email
 FROM Affected A, Subscriber S, SubscriberAffected SA
 WHERE SA.affid = A.affid and SA.subid = S.subid
`

var subscriptionsQueryMapper = func() (*QuerySubscriptionResult, []any) {
	aul := QuerySubscriptionResult{}

	return &aul, []any{&aul.HSHA256Email, &aul.B64Email}
}

var subscriptionsWithoutAffectedQueryMapper = func() (*QuerySubscriptionWithoutAffectedResult, []any) {
	aul := QuerySubscriptionWithoutAffectedResult{}

	return &aul, []any{&aul.B64Email}
}

func StoreSubscriptionDB(dbctx database.DatabaseContext[database.Record], subscription subscribe.Subscription) error {

	return dbctx.InsertSubscription(subscription)
}

func QuerySubscriptionsDB(dbctx database.DatabaseContext[database.Record]) (QueryAllSubscriptionsResult, error) {
	ctx := database.Convert[database.Record, QuerySubscriptionResult](dbctx)

	subscriptions, err := querySubscriptions(ctx)

	if err != nil {
		return nil, err
	}

	ctxAll := database.Convert[database.Record, QuerySubscriptionWithoutAffectedResult](dbctx)
	subscriptionsWithoutAffected, err := queryAllSubscriptions(ctxAll)

	if err != nil {
		return nil, err
	}

	sub := subscriptions.ConvertToQueryAllSubscriptionsResult()

	for _, v := range subscriptionsWithoutAffected {
		sub = append(sub, v.ConvertToSubscription())
	}

	return sub, nil
}

func queryAllSubscriptions(dbctx database.DatabaseContext[QuerySubscriptionWithoutAffectedResult]) ([]QuerySubscriptionWithoutAffectedResult, error) {
	q, m, vs := prepareAllSubscriptionsQuery()

	return dbctx.CustomQuery(q, m, vs...)
}

func querySubscriptions(dbctx database.DatabaseContext[QuerySubscriptionResult]) (QuerySubscriptionsResult, error) {
	q, m, vs := prepareSubscriptionsQuery()

	return dbctx.CustomQuery(q, m, vs...)
}

func prepareSubscriptionsQuery() (string, database.TypedQueryResultMapper[QuerySubscriptionResult], []any) {
	return subscriptionsQuery, subscriptionsQueryMapper, []any{}
}

func prepareAllSubscriptionsQuery() (string, database.TypedQueryResultMapper[QuerySubscriptionWithoutAffectedResult], []any) {
	return subscriptionsWithoutAffectedQuery, subscriptionsWithoutAffectedQueryMapper, []any{}
}
