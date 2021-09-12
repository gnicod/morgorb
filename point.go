package morgorb

import (
	"database/sql/driver"
	"fmt"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/wkb"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Point struct {
	orb orb.Point
}

func NewPoint(x float64, y float64) Point {
	return Point{
		orb: orb.Point{x, y},
	}
}

func (p *Point) Scan(value interface{}) error {
	bytes, _ := value.([]byte)
	g, err := wkb.Unmarshal(bytes)
	if err != nil {
		return err
	}
	p.orb = g.(orb.Point)
	return err
}

func (p Point) Value() (driver.Value, error) {
	return fmt.Sprintf("POINT(%v %v)", p.orb.Lat(), p.orb.Lon()), nil
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
