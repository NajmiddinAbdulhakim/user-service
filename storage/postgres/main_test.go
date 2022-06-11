package postgres

import (
	"log"
	"os"
	"testing"

	"github.com/NajmiddinAbdulhakim/user-service/config"
	"github.com/NajmiddinAbdulhakim/user-service/pkg/db"
	"github.com/NajmiddinAbdulhakim/user-service/pkg/logger"

)

var repo *UserRepo

func TestMain(m *testing.M) {
	cfg := config.Load()

	connDB, err := db.ConnectToDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to psotgres error", logger.Error(err))
	}

	repo = NewUserRepo(connDB)

	os.Exit(m.Run())
}
