package geo

import (
	"context"
	"testing"

	"github.com/semi-technologies/weaviate/entities/filters"
	"github.com/semi-technologies/weaviate/entities/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGeoJourney(t *testing.T) {
	elements := []models.GeoCoordinates{
		{ // coordinates of munich
			Latitude:  ptFloat32(48.13743),
			Longitude: ptFloat32(11.57549),
		},
		{ // coordinates of stuttgart
			Latitude:  ptFloat32(48.78232),
			Longitude: ptFloat32(9.17702),
		},
	}

	getCoordinates := func(ctx context.Context, id int32) (models.GeoCoordinates, error) {
		return elements[id], nil
	}

	geoIndex, err := NewIndex(Config{
		ID:                 "unit-test",
		CoordinatesForID:   getCoordinates,
		DisablePersistence: true,
		RootPath:           "doesnt-matter-persistence-is-off",
	})
	require.Nil(t, err)

	t.Run("importing all", func(t *testing.T) {
		for id, coordinates := range elements {
			err := geoIndex.Add(id, coordinates)
			require.Nil(t, err)
		}
	})

	t.Run("importing an invalid object", func(t *testing.T) {
		err := geoIndex.Add(9000, models.GeoCoordinates{})
		assert.Equal(t, "invalid arguments: latitude must be set", err.Error())
	})

	km := float32(1000)
	t.Run("searching missing longitude", func(t *testing.T) {
		_, err := geoIndex.WithinRange(context.Background(), filters.GeoRange{
			GeoCoordinates: &models.GeoCoordinates{
				Latitude: ptFloat32(48.13743),
			},
			Distance: 300 * km,
		})
		assert.Equal(t, "invalid arguments: longitude must be set", err.Error())
	})

	t.Run("searching missing latitude", func(t *testing.T) {
		_, err := geoIndex.WithinRange(context.Background(), filters.GeoRange{
			GeoCoordinates: &models.GeoCoordinates{
				Longitude: ptFloat32(11.57549),
			},
			Distance: 300 * km,
		})
		assert.Equal(t, "invalid arguments: latitude must be set", err.Error())
	})

	t.Run("searching within 300km of munich", func(t *testing.T) {
		// should return both cities, with munich first and stuttgart second
		results, err := geoIndex.WithinRange(context.Background(), filters.GeoRange{
			GeoCoordinates: &models.GeoCoordinates{
				Latitude:  ptFloat32(48.13743),
				Longitude: ptFloat32(11.57549),
			},
			Distance: 300 * km,
		})
		require.Nil(t, err)

		expectedResults := []int{0, 1}
		assert.Equal(t, expectedResults, results)
	})
}

func ptFloat32(in float32) *float32 {
	return &in
}