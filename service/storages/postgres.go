package storages

import (
	"strings"

	"github.com/JUSSIAR/Golang-gRPC-blog/service/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectPostgres() *gorm.DB {
	creds := []string{
		"host=localhost",
		"port=5432",
		"dbname=grpc-blog-storage",
		"sslmode=disable",
		"user=jussiar",
		"password=admin",
	}

	postgresDB, err := gorm.Open(
		postgres.Open(strings.Join(creds, " ")),
		&gorm.Config{},
	)
	if err != nil {
		panic(err)
	}

	err = postgresDB.AutoMigrate(&models.Post{}, &models.Comment{})
	if err != nil {
		panic(err)
	}

	return postgresDB
}
