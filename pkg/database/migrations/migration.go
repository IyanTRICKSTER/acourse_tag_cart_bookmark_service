package migrations

import "acourse_tag_cart_bookmark_service/pkg/database"

type Migration struct {
	DB *database.Database
}

func Construct(db *database.Database) *Migration {
	return &Migration{DB: db}
}
