package main

import (
	"fmt"
	"github.com/gnicod/georm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Geom struct {
	gorm.Model
	Name  string
	Point georm.Point `gorm:"srid:4326"`
	LineString georm.LineString `gorm:"srid:4326"`
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
	point , err := georm.NewPoint(50, 23)
	if err != nil {
		panic(err)
	}
	// insert a linestinrg
	linestring , err := georm.NewLineString([]float64{50, 23}, []float64{89, 78})
	if err != nil {
		panic(err)
	}
	g := Geom{Point: point, Name: "test", LineString: linestring}
	db.Create(&g)

	var fetched Geom
	db.First(&fetched, 23)
	gjson, _ := fetched.Point.ToGeoJson()
	print(gjson)
	gjsonl, _ := fetched.LineString.ToGeoJson()
	print(gjsonl)
	tx := db.Raw("select st_asgeojson(line_string) from geoms;")

	var geojson string
	tx.Scan(&geojson)
	print(geojson)

}
