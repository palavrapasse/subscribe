package data

import (
	"fmt"

	"github.com/palavrapasse/damn/pkg/database"
	"github.com/palavrapasse/damn/pkg/entity"
)

const leakByIdQuery = `
SELECT L.*, P.name FROM Leak L, LeakPlatform LP, Platform P
WHERE L.leakid = %d and LP.leakid = L.leakid and LP.platid = P.platid
`

var leakByIdQueryMapper = func() (*QueryLeakByIdResult, []any) {
	aul := QueryLeakByIdResult{}

	return &aul, []any{&aul.LeakId, &aul.ShareDateSC, &aul.Context, &aul.Name}
}

func QueryLeakByIdDB(dbctx database.DatabaseContext[database.Record], leakid entity.AutoGenKey) ([]QueryLeakByIdResult, error) {
	ctx := database.Convert[database.Record, QueryLeakByIdResult](dbctx)

	return queryLeakById(ctx, leakid)
}

func queryLeakById(dbctx database.DatabaseContext[QueryLeakByIdResult], leakid entity.AutoGenKey) ([]QueryLeakByIdResult, error) {
	q, m, vs := prepareLeakByIdQuery(leakid)

	return dbctx.CustomQuery(q, m, vs...)
}

func prepareLeakByIdQuery(leakid entity.AutoGenKey) (string, database.TypedQueryResultMapper[QueryLeakByIdResult], []any) {
	return fmt.Sprintf(leakByIdQuery, leakid), leakByIdQueryMapper, []any{}
}
