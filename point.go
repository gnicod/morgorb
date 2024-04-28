package georm

import (
	"database/sql/driver"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"github.com/twpayne/go-geom/encoding/geojson"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Point struct {
	geom *geom.Point
}

func (p Point) ToGeoJson() (string, error) {
	geometry, err := geojson.Marshal(p.geom)
	return string(geometry), err
}

func NewGeormPoint(point geom.Point) Point {
	return Point{
		geom: &point,
	}
}

func NewPoint(coordinates ...float64) (Point, error) {
	switch len(coordinates) {
	case 2:
		return Point{
			geom: geom.NewPointFlat(geom.XY, []float64{coordinates[0], coordinates[1]}),
		}, nil
	case 3:
		return Point{
			geom: geom.NewPointFlat(geom.XYZ, []float64{coordinates[0], coordinates[1], coordinates[2]}),
		}, nil
	default:
		return Point{}, errors.New("point must have 2 or 3 coordinates")
	}
}

func (p *Point) Scan(value interface{}) error {
	data, err := hex.DecodeString(value.(string))
	if err != nil {
		return nil
	}

	point, err := ewkb.Unmarshal(data)
	p.geom = point.(*geom.Point)
	return err
}

func (p Point) Value() (driver.Value, error) {
	switch p.geom.Layout() {
	case geom.XY:
		return fmt.Sprintf("SRID=4326;POINT(%v %v)", p.geom.X(), p.geom.Y()), nil
	case geom.XYZ:
		return fmt.Sprintf("SRID=4326;POINT(%v %v %v)", p.geom.X(), p.geom.Y(), p.geom.Z()), nil
	default:
		return "", errors.New(fmt.Sprintf("layout %s not implemented", p.geom.Layout()))
	}
}

func (Point) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "mysql", "sqlite":
		// TODO for now there is no way it could work with mysql
		return "JSON"
	case "postgres":
		srid, exists := field.TagSettings["SRID"]
		if !exists {
			srid = "4326"
		}
		return fmt.Sprintf("geometry(Point, %s)", srid)
	}
	return ""
}

func (Point) GormDataType() string {
	return "geometry(Point, 4326)"
}
