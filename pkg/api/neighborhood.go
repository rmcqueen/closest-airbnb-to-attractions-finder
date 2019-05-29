package api

import (
	"fmt"
	"log"

	"../connections"

	_ "github.com/lib/pq"
)

type Neighborhood struct {
	Name                string  `json:"name"`
	City                string  `json:"city_name"`
	StateOrProvinceName string  `json:"state_or_province_name"`
	Country             string  `json:"country"`
	Latitude            float64 `json:"latitude"`
	Longitude           float64 `json:"longitude"`
}

func FindNeighborhoodContainingAttraction(attraction Attraction) (Neighborhood, error) {
	attractionInNeighborhoodQuery := `
        SELECT ST_Contains(neighborhood_poly, attr_point) as in_neighborhood, name, city, state, country
        FROM (
            SELECT ST_SetSRID(ST_Point($1, $2),4326) as attr_point, geom as neighborhood_poly, name, city, state, country
            FROM neighborhood_geocoding.neighborhoods
        ) as foo
        WHERE ST_Contains(neighborhood_poly, attr_point) is true
        `

	rows, err := connections.PostgreSqlConnector{}.Connect().Query(
		attractionInNeighborhoodQuery,
		attraction.Longitude,
		attraction.Latitude)

	if err != nil {
		return Neighborhood{}, err
	}

	defer rows.Close()

	var matchedNeighborhoods []Neighborhood
	minDistanceInMeters := 999999.0
	i := 0
	bestNeighborhoodIdx := 0

	for rows.Next() {
		var name string
		var city string
		var stateOrProvinceName string
		var country string
		var inNeighborhood bool
		if err := rows.Scan(&inNeighborhood, &name, &city, &stateOrProvinceName, &country); err != nil {
			return Neighborhood{}, err
		}

		if inNeighborhood == false {
			continue
		}

		coordinates, err := resolveNeighborhoodMultiPolygonsCentroidPoint(name, city, stateOrProvinceName)

		if err != nil {
			log.Printf("Unable to resolve coordinates for %s", name)
			continue
		}

		// TODO: make coordinates struct since there is so much re-use throughout the app
		latitude := coordinates[0]
		longitude := coordinates[1]
		attractionsCoordinates := []float64{attraction.Longitude, attraction.Latitude}
		fmt.Printf("Getting distance for: %s\n", attraction.Name)
		distanceInMeters, err := getDistanceBetweenTwoCoordinates(coordinates, attractionsCoordinates)

		if err != nil {
			log.Fatal(err)
			continue
		}

		neighborhood := Neighborhood{name, city, stateOrProvinceName, country, latitude, longitude}
		matchedNeighborhoods = append(matchedNeighborhoods, neighborhood)
		if distanceInMeters < minDistanceInMeters {
			fmt.Printf("Distance before setting %.6f\n", minDistanceInMeters)
			minDistanceInMeters = distanceInMeters
			fmt.Printf("Distance after setting %.6f\n", minDistanceInMeters)
			bestNeighborhoodIdx = i
		}
		i++
	}

	if len(matchedNeighborhoods) == 0 {
		return Neighborhood{}, err
	}

	return matchedNeighborhoods[bestNeighborhoodIdx], err
}

func resolveNeighborhoodMultiPolygonsCentroidPoint(
	neighborhoodName string,
	neighborhoodCity string,
	neighborhoodState string) ([]float64, error) {
	/**
	  Returns the coordinates of a MultiPolygon's centroid (if found).
	  idx 0 => latitude, idx 1 => longitude
	*/
	centroidQueryStr := `
    SELECT ST_X(coordinates) as longitude, ST_Y(coordinates) as latitude
    FROM (
        SELECT ST_AsText(ST_centroid(multi_poly)) as coordinates
        FROM (
            SELECT geom as multi_poly
            FROM neighborhood_geocoding.neighborhoods
            WHERE name ilike $1
                AND city ilike $2
                AND state ilike $3
            ) as coordinates
        ) as result
    `

	row := connections.PostgreSqlConnector{}.Connect().QueryRow(
		centroidQueryStr,
		neighborhoodName,
		neighborhoodCity,
		neighborhoodState)

	coordinates := make([]float64, 2)
	err := row.Scan(&coordinates[0], &coordinates[1])

	if err != nil {
		return []float64{}, err
	}

	return coordinates, err
}

func getDistanceBetweenTwoCoordinates(point1 []float64, point2 []float64) (float64, error) {
	// Get distance in meters. See: https://postgis.net/docs/manual-1.4/ST_Distance_Sphere.html
	pointDistanceQueryStr := `
    SELECT ST_Distance_Sphere(
        ST_SetSRID(ST_Point($1, $2), 4326),
        ST_SetSRID(ST_Point($3, $4), 4326)
    ) as distance_in_meters`

	fmt.Printf("Point1: %v, Point2: %v\n", point1, point2)
	row := connections.PostgreSqlConnector{}.Connect().QueryRow(
		pointDistanceQueryStr,
		point1[0],
		point1[1],
		point2[0],
		point2[1])

	var distanceInMeters float64
	err := row.Scan(&distanceInMeters)

	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	return distanceInMeters, err
}
