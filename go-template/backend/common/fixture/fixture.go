package fixture

import (
	"context"
	"os"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dbfixture"

	"github.com/hayrat/go-template2/backend/common/model"
)

func Load(db *bun.DB) error {
	db.RegisterModel((*model.User)(nil))
	db.RegisterModel((*model.Post)(nil))
	db.RegisterModel((*model.Comment)(nil))
	db.RegisterModel((*model.Like)(nil))

	fixture := dbfixture.New(db, dbfixture.WithTruncateTables())
	err := fixture.Load(context.Background(), os.DirFS("../../common/fixture"), "fixture.yml")

	return err
}
