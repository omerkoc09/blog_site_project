package migrations

import (
	"context"

	"github.com/uptrace/bun"

	migration "github.com/hayrat/go-template2/backend/common/migrations/20220902105002_bismillah"
)

func init() {
	m := []interface{}{
		&migration.User{},
		&migration.AuthRefreshToken{},
		&migration.Post{},
		&migration.Comment{},
		&migration.Like{},
		&migration.Follow{},
		&migration.Saved{},
	}

	up := func(ctx context.Context, db *bun.DB) error {
		return db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
			for _, i := range m {
				if _, err := tx.NewCreateTable().Model(i).IfNotExists().WithForeignKeys().Exec(ctx); err != nil {
					return err
				}
			}
			return nil
		})
	}

	down := func(ctx context.Context, db *bun.DB) error {
		return db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
			for _, i := range m {
				if _, err := tx.NewDropTable().Model(i).IfExists().Cascade().Exec(ctx); err != nil {
					return err
				}
			}
			return nil
		})
	}

	if err := Migrations.Register(up, down); err != nil {
		panic(err)
	}
}
