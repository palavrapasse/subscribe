package data

import (
	"os"

	"github.com/palavrapasse/damn/pkg/database"
)

const (
	leaksDbFilePathEnvKey         = "leaksdb_fp"
	subscriptionsDbFilePathEnvKey = "subscriptionsdb_fp"
)

var (
	leaksDbFilePath         = os.Getenv(leaksDbFilePathEnvKey)
	subscriptionsDbFilePath = os.Getenv(subscriptionsDbFilePathEnvKey)
)

func OpenLeaksDB() (database.DatabaseContext[database.Record], error) {
	return database.NewDatabaseContext[database.Record](leaksDbFilePath)
}

func OpenSubscriptionsDB() (database.DatabaseContext[database.Record], error) {
	return database.NewDatabaseContext[database.Record](subscriptionsDbFilePath)
}

func Close(dbctx database.DatabaseContext[database.Record]) error {
	return dbctx.DB.Close()
}
