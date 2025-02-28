package database

import (
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"

	//"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"github.com/uptrace/bun/schema"

	"github.com/hayrat/go-template2/backend/pkg/config"
)

func New(config config.DbConfig) *bun.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=prefer",
		config.Host, config.Port, config.Name, config.Username, config.Password,
	)

	//sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn) /*, pgdriver.WithTLSConfig(&tls.Config{InsecureSkipVerify: true})*/))

	c, err := pgx.ParseConfig(dsn)
	if err != nil {
		panic(err)
	}
	c.PreferSimpleProtocol = true

	sqldb := stdlib.OpenDB(*c)

	sqldb.SetMaxIdleConns(config.MaxIdleConn)
	sqldb.SetMaxOpenConns(config.MaxPoolSize)
	sqldb.SetConnMaxLifetime(time.Duration(config.MaxLifetime) * time.Second)

	// tablo isimlerini plural yapmasın.
	schema.SetTableNameInflector(func(s string) string {
		return s
	})

	db := bun.NewDB(sqldb, pgdialect.New())
	if err := db.Ping(); err != nil {
		fmt.Println(err)
	}

	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(config.Debug)))

	return db
}
