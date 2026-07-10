package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"turtle/core/lgr"
)

type DbCompany struct {
	ID    int    `gorm:"primaryKey"`
	Name  string `gorm:"type:text"`
	Email string `gorm:"type:text"`
	TaxID string `gorm:"type:text"`
	VatID string `gorm:"type:text"`
}

type PositionsStream struct {
	ID     int     `gorm:"primaryKey;autoIncrement;columnt:id"`
	At     int64   `gorm:"index;columnt:id"`      // Indexed for efficient queries
	AreaId int8    `gorm:"index;columnt:area_id"` // Indexed
	TwinId int64   `gorm:"index;columnt:twin_id"` // Indexed
	PosX   float32 `gorm:"not null;columnt:posX"` // Ensuring it is not null
	PosY   float32 `gorm:"not null;columnt:posY"`
	PosZ   float32 `gorm:"not null;columnt:posZ"`
}

func InitGorm() {
	db, err := gorm.Open(sqlite.Open("turtle.db"), &gorm.Config{})

	if err != nil {
		lgr.Error(err.Error())
	}

	err = db.AutoMigrate(
		&DbCompany{},
		&PositionsStream{},
	)

	x := 0

	for i := 0; i < 1000000; i++ {
		feed := PositionsStream{}
		feed.PosX = float32(i)
		db.Create(&feed)

		if x > 10000 {
			db.Commit()
			x = 0
			lgr.Ok("Done commit")
		}

		x += 1

	}

	db.Commit()

	lgr.Error(err.Error())

}
