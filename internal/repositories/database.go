package repositories

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // This need to be here, so bun knows which sql driver to use
	"github.com/spf13/viper"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func MustSetupDatabase() *bun.DB {
	connectionString, err := getConnectionString()
	if err != nil {
		panic(err)
	}

	sqlClient := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(connectionString)))

	client := bun.NewDB(sqlClient, pgdialect.New())

	return client
}

func getConnectionString() (string, error) {
	var err error
	username := viper.GetString("DB_USER")
	password := viper.GetString("DB_PASSWORD")
	host := viper.GetString("DB_HOST")
	port := viper.GetString("DB_PORT")
	dbname := viper.GetString("DB_NAME")

	if username == "" || password == "" || host == "" || port == "" || dbname == "" {
		err = fmt.Errorf("missing required environment variables")
	}

	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, dbname)

	return connectionString, err
}
