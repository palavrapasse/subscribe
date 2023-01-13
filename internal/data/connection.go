package data

import (
	"os"

	"github.com/palavrapasse/damn/pkg/database"
)

const (
	leaksDbFilePathEnvKey = "leaksdb_fp"
)

var (
	leaksDbFilePath = os.Getenv(leaksDbFilePathEnvKey)
)

func Open() (database.DatabaseContext[database.Record], error) {
	return database.NewDatabaseContext[database.Record](leaksDbFilePath)
}

func Close(dbctx database.DatabaseContext[database.Record]) error {
	return dbctx.DB.Close()
}
