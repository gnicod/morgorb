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

type LineString struct {
	geom *geom.LineString
}

func (p LineString) ToGeoJson() (string, error) {
	geometry, err := geojson.Marshal(p.geom)
	return string(geometry), err
}

func NewLineString(coordinates ...float64) (LineString, error) {
	switch len(coordinates) {
	case 2:
		return LineString{
			// TODO iterate
			geom: geom.NewLineStringFlat(geom.XY, []float64{coordinates[0], coordinates[1]} ,[]int{10}), //.SetSRID(4326),
		}, nil
	case 3:
		return LineString{
			// TODO iterate
			geom: geom.NewLineStringFlat(geom.XYZ , []float64{coordinates[0], coordinates[1], coordinates[2]}, []int{10}),
		}, nil
	default:
		return LineString{}, errors.New("point must have 2 or 3 coordinates")
	}
}

func (p *LineString) Scan(value interface{}) error {
	data, err := hex.DecodeString(value.(string))
	if err != nil {
		return nil
	}

	point, err := ewkb.Unmarshal(data)
	p.geom = point.(*geom.LineString)
	return err
}

func (p LineString) Value() (driver.Value, error) {
		// TODO LINESTRING(0 0, 1 1, 2 1, 2 2)
		switch p.geom.Layout() {
		case geom.XY:
			return fmt.Sprintf("LINESTRING(%v %v)", p.geom.X(), p.geom.Y()), nil
		case geom.XYZ:
			return fmt.Sprintf("LINESTRING(%v %v %v)", p.geom.X(), p.geom.Y(), p.geom.Z()), nil
		default:
			return "", errors.New(fmt.Sprintf("layout %s not implemented", p.geom.Layout()))
		}
}

func (LineString) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "mysql", "sqlite":
		// TODO for now there is no way it could work with mysql
		return "JSON"
	case "postgres":
		srid, exists := field.TagSettings["SRID"]
		if !exists {
			srid = "4326"
		}
		return fmt.Sprintf("geometry(LINESTRING, %s)", srid)
	}
	return ""
}

func (LineString) GormDataType() string {
	return "geometry(LineString, 4326)"
}
