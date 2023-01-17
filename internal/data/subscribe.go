package data

import (
	"fmt"

	"github.com/palavrapasse/damn/pkg/database"
	"github.com/palavrapasse/damn/pkg/entity"
	"github.com/palavrapasse/damn/pkg/entity/subscribe"
)

const subscriptionsWithoutAffectedQuery = `
SELECT S.b64email FROM  Subscriber S
WHERE S.subid NOT IN
    (SELECT SA.subid
    FROM SubscriberAffected SA) 
`

const subscriptionsQuery = `
SELECT A.hsha256email, S.b64email FROM Affected A, Subscriber S, SubscriberAffected SA
WHERE SA.affid = A.affid and SA.subid = S.subid
`

const leaksByLeakIdUserHashPreparedQuery = `
SELECT HU.hsha256 FROM Leak L, LeakUser LU, User U, HashUser HU
    WHERE L.leakid = ? and LU.leakid = L.leakid and LU.userid = U.userid and U.userid = HU.userid 
    AND HU.hsha256 IN (%s)
`

var subscriptionsQueryMapper = func() (*QuerySubscriptionResult, []any) {
	aul := QuerySubscriptionResult{}

	return &aul, []any{&aul.HSHA256Email, &aul.B64Email}
}

var subscriptionsWithoutAffectedQueryMapper = func() (*QuerySubscriptionWithoutAffectedResult, []any) {
	aul := QuerySubscriptionWithoutAffectedResult{}

	return &aul, []any{&aul.B64Email}
}

var leaksQueryMapper = func() (*entity.HSHA256, []any) {
	aul := entity.HSHA256("")

	return &aul, []any{&aul}
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

func QueryLeaksDB(dbctx database.DatabaseContext[database.Record], leakid entity.AutoGenKey, affected []subscribe.Affected) (QueryLeaksResult, error) {
	ctx := database.Convert[database.Record, entity.HSHA256](dbctx)

	return queryLeaksThatAffectUsers(ctx, leakid, affected)
}

func queryAllSubscriptions(dbctx database.DatabaseContext[QuerySubscriptionWithoutAffectedResult]) ([]QuerySubscriptionWithoutAffectedResult, error) {
	q, m, vs := prepareAllSubscriptionsQuery()

	return dbctx.CustomQuery(q, m, vs...)
}

func querySubscriptions(dbctx database.DatabaseContext[QuerySubscriptionResult]) (QuerySubscriptionsResult, error) {
	q, m, vs := prepareSubscriptionsQuery()

	return dbctx.CustomQuery(q, m, vs...)
}

func queryLeaksThatAffectUsers(dbctx database.DatabaseContext[entity.HSHA256], leakid entity.AutoGenKey, affected []subscribe.Affected) (QueryLeaksResult, error) {
	q, m, vs := prepareAffectedUsersQuery(leakid, affected)

	return dbctx.CustomQuery(q, m, vs...)
}

func prepareAllSubscriptionsQuery() (string, database.TypedQueryResultMapper[QuerySubscriptionWithoutAffectedResult], []any) {
	return subscriptionsWithoutAffectedQuery, subscriptionsWithoutAffectedQueryMapper, []any{}
}

func prepareSubscriptionsQuery() (string, database.TypedQueryResultMapper[QuerySubscriptionResult], []any) {
	return subscriptionsQuery, subscriptionsQueryMapper, []any{}
}

func prepareAffectedUsersQuery(leakid entity.AutoGenKey, affected []subscribe.Affected) (string, database.TypedQueryResultMapper[entity.HSHA256], []any) {
	values := []any{}

	values = append(values, leakid)
	for _, v := range affected {
		values = append(values, string(v.HSHA256Email))
	}

	return fmt.Sprintf(leaksByLeakIdUserHashPreparedQuery, database.MultiplePlaceholder(len(values)-1)), leaksQueryMapper, values
}
