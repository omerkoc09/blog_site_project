package main

import (
	"fmt"
	"os"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
	"github.com/urfave/cli/v2"

	"github.com/hayrat/go-template2/backend/common/fixture"
	"github.com/hayrat/go-template2/backend/common/migrations"
	"github.com/hayrat/go-template2/backend/pkg/config"
	"github.com/hayrat/go-template2/backend/pkg/database"
)

func main() {
	_, err := config.Setup()
	if err != nil {
		panic(err)
	}

	app := &cli.App{
		Name:     "bun",
		Commands: commands,
	}

	err = app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

var commands = []*cli.Command{
	{
		Name:  "init",
		Usage: "create migration tables",
		Action: func(c *cli.Context) error {
			migrator := getMigrator()
			return migrator.Init(c.Context)
		},
	},
	{
		Name:  "migrate",
		Usage: "migrate database",
		Action: func(c *cli.Context) error {
			migrator := getMigrator()

			group, err := migrator.Migrate(c.Context)
			if err != nil {
				return err
			}

			if group.ID == 0 {
				fmt.Printf("there are no new migrations to run\n")
				return nil
			}

			fmt.Printf("migrated to %s\n", group)
			return nil
		},
	},
	{
		Name:  "rollback",
		Usage: "rollback the last migration group",
		Action: func(c *cli.Context) error {
			migrator := getMigrator()

			group, err := migrator.Rollback(c.Context)
			if err != nil {
				return err
			}

			if group.ID == 0 {
				fmt.Printf("there are no groups to roll back\n")
				return nil
			}

			fmt.Printf("rolled back %s\n", group)
			return nil
		},
	}, {
		Name:  "status",
		Usage: "print migrations status",
		Action: func(c *cli.Context) error {
			migrator := getMigrator()

			ms, err := migrator.MigrationsWithStatus(c.Context)
			if err != nil {
				return err
			}
			fmt.Printf("migrations: %s\n", ms)
			fmt.Printf("unapplied migrations: %s\n", ms.Unapplied())
			fmt.Printf("last migration group: %s\n", ms.LastGroup())

			return nil
		},
	}, {
		Name:  "seed",
		Usage: "seed the database with fixtures",
		Action: func(c *cli.Context) error {
			db := getDB()
			err := fixture.Load(db)
			if err != nil {
				return err
			}
			fmt.Println("fixtures loaded")
			return nil
		},
	},
}

func getMigrator() *migrate.Migrator {
	db := getDB()

	return migrate.NewMigrator(db, migrations.Migrations, migrate.WithMarkAppliedOnSuccess(true))
}

func getDB() *bun.DB {
	cfg := config.Get()
	return database.New(cfg.Database)
}
