package main

import (
	"github.com/xxarupkaxx/anke-two/infrastructure"
	"os"
)

func main() {
	env, ok := os.LookupEnv("ANKE-TWO_ENV")
	if !ok {
		env = "production"
	}

	logOn := env == "pprof" || env == "dev"

	sqlHandler, err := infrastructure.EstablishConnection(logOn)
	if err != nil {
		panic(err)
	}

	err = infrastructure.Migrate(sqlHandler.Db)
	if err != nil {
		panic(err)
	}

	port, ok := os.LookupEnv("PORT")
	if !ok {
		panic("no PORT")
	}

	SetRouting(port, sqlHandler.Db)
}
