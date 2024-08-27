package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "./internal/infrastructure/db/query",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})

	dsn := "host=localhost user=admin password=admin dbname=auth port=5432 sslmode=disable TimeZone=Asia/Tokyo"
	gormdb, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	g.UseDB(gormdb) // reuse your gorm db

	all := g.GenerateAllTable() // database to table model.

	g.ApplyBasic(all...)

	// Generate the code
	g.Execute()
}
