package osm

import (
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/planar"
	"io/ioutil"
	"strings"
)

var tehranFc *geojson.FeatureCollection

func init() {
	b, _ := ioutil.ReadFile("data/tehran-districts.geojson")
	tehranFc, _ = geojson.UnmarshalFeatureCollection(b)
}

func NormalizeOsmRelationId(rid string) string {
	return strings.Replace(rid, "relation/", "", 1)
}

func FindOsmRelationID(lat, long float64) string {
	point := orb.Point{long, lat}
	for _, feature := range tehranFc.Features {
		polygon, isPoly := feature.Geometry.(orb.Polygon)
		if isPoly {
			if planar.PolygonContains(polygon, point) {
				return feature.Properties.MustString("@id")
			}
		}
	}
	return ""
}
