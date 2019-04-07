package migrations

import (
	"github.com/jinzhu/gorm"
	"github.com/smartcontractkit/chainlink/store/migrations/migration0"
	"github.com/smartcontractkit/chainlink/store/migrations/migration1549496047"
	"github.com/smartcontractkit/chainlink/store/migrations/migration1551816486"
	"github.com/smartcontractkit/chainlink/store/migrations/migration1551895034"
	"github.com/smartcontractkit/chainlink/store/migrations/migration1552418531"
	"github.com/smartcontractkit/chainlink/store/migrations/migration1553029703"
	"github.com/smartcontractkit/chainlink/store/migrations/migration1554131520"
	gormigrate "gopkg.in/gormigrate.v1"
)

type migration interface {
	Migrate(tx *gorm.DB) error
}

// Migrate iterates through available migrations, running and tracking
// migrations that have not been run.
func Migrate(db *gorm.DB) error {
	err := upgradeOldMigrationSchema(db)
	if err != nil {
		return err
	}

	options := gormigrate.DefaultOptions
	options.IDColumnSize = 12

	m := gormigrate.New(db, options, []*gormigrate.Migration{
		{
			ID:      "0",
			Migrate: migration0.Migration{}.Migrate,
		},
		{
			ID:      "1549496047",
			Migrate: migration1549496047.Migration{}.Migrate,
		},
		{
			ID:      "1551816486",
			Migrate: migration1551816486.Migration{}.Migrate,
		},
		{
			ID:      "1551895034",
			Migrate: migration1551895034.Migration{}.Migrate,
		},
		{
			ID:      "1552418531",
			Migrate: migration1552418531.Migration{}.Migrate,
		},
		{
			ID:      "1553029703",
			Migrate: migration1553029703.Migration{}.Migrate,
		},
		{
			ID:      "1554131520",
			Migrate: migration1554131520.Migration{}.Migrate,
		},
	})

	return m.Migrate()
}

func upgradeOldMigrationSchema(db *gorm.DB) error {
	if !db.HasTable("migration_timestamps") {
		return nil
	}

	tx := db.Begin()
	err := tx.Exec(`
ALTER TABLE migration_timestamps RENAME TO migrations;
ALTER TABLE migrations RENAME COLUMN timestamp to id;
`).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
