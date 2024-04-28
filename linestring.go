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
	"strings"
)

type LineString struct {
	geom *geom.LineString
}

func (p LineString) ToGeoJson() (string, error) {
	geometry, err := geojson.Marshal(p.geom)
	return string(geometry), err
}

func (p *LineString) ToLineString() *geom.LineString {
	return p.geom
}

func flatten(m [][]float64) []float64 {
	res := []float64{}
	for i := range m {
		res = append(res, m[i]...)
	}
	return res
}

func NewGeormLineString(lineString geom.LineString) LineString {
	return LineString{
		geom: &lineString,
	}
}

func NewLineString(coordinates ...[]float64) (LineString, error) {
	flattenCoordinate := flatten(coordinates)
	switch len(coordinates[0]) {
	case 2:
		return LineString{
			// TODO iterate
			geom: geom.NewLineStringFlat(geom.XY, flattenCoordinate).SetSRID(4326), // TODO SRID should be configurable
		}, nil
	case 3:
		return LineString{
			// TODO iterate
			geom: geom.NewLineStringFlat(geom.XYZ, flattenCoordinate),
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
	strLineString := ""
	switch p.geom.Layout() {
	case geom.XY:
		for _, coord := range p.geom.Coords() {
			strLineString += fmt.Sprintf("%v %v,", coord[0], coord[1])
		}
	case geom.XYZ:
		for _, coord := range p.geom.Coords() {
			strLineString += fmt.Sprintf("%v %v %v,", coord[0], coord[1], coord[2])
		}
	default:
		return "", errors.New(fmt.Sprintf("layout %s not implemented", p.geom.Layout()))
	}
	strLineString = strings.TrimSuffix(strLineString, ",")
	return fmt.Sprintf("SRID=4326;LINESTRING(%v)", strLineString), nil
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
