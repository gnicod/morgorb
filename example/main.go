package main

import (
	"fmt"
	"github.com/gnicod/morgorb"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Geom struct {
	gorm.Model
	Name  string
	Point morgorb.Point `gorm:"srid:4326"`
}

func main() {
	host := "127.0.0.1"
	user := "postgres"
	password := "mysecretpassword"
	dbname := "test"
	psqlInfo := fmt.Sprintf("host=%s port=5432 user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, user, password, dbname)
	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	// Migrate the schema
	db.AutoMigrate(&Geom{})

	// insert a point
	point , err := morgorb.NewPoint(50, 23)
	if err != nil {
		panic(err)
	}
	g := Geom{Point: point, Name: "test"}
	db.Create(&g)

	var fetched Geom
	db.First(&fetched, 23)
	gjson, _ := fetched.Point.ToGeoJson()
	print(gjson)
}
