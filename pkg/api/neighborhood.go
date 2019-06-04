package api

import (
	"container/heap"
	"crypto/md5"
	"encoding/hex"
	"log"
	"sort"
	"strings"

	"../connections"

	_ "github.com/lib/pq" // Used to interact with PostgreSQL/PostGIS
)

// Neighborhood is defined as a localised community within a larger city (i.e, 'Downtown')
// TODO: make lat/lng a struct
type Neighborhood struct {
	Name                string  `json:"name"`
	City                string  `json:"city_name"`
	StateOrProvinceName string  `json:"state_or_province_name"`
	Country             string  `json:"country"`
	Latitude            float64 `json:"latitude"`
	Longitude           float64 `json:"longitude"`
}

// FindNeighborhoodContainingAttraction resolves the neighborhood of the given attraction via geocoding.
func FindNeighborhoodContainingAttraction(attraction Attraction) (Neighborhood, error) {
	attractionInNeighborhoodQuery := `
        SELECT ST_Contains(neighborhood_poly, attr_point) as in_neighborhood, name, city, state, country
        FROM (
            SELECT ST_SetSRID(ST_Point($1, $2),4326) as attr_point, geom as neighborhood_poly, name, city, state, country
            FROM neighborhood_geocoding.neighborhoods
        ) as foo
        WHERE ST_Contains(neighborhood_poly, attr_point) is true
        `

	rows, err := connections.Init().Query(
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
		distanceInMeters, err := getDistanceBetweenTwoCoordinates(coordinates, attractionsCoordinates)

		if err != nil {
			log.Printf("Unable to get distance between two coordinates having error: %v\n", err)
			continue
		}

		neighborhood := Neighborhood{name, city, stateOrProvinceName, country, latitude, longitude}
		matchedNeighborhoods = append(matchedNeighborhoods, neighborhood)
		if distanceInMeters < minDistanceInMeters {
			minDistanceInMeters = distanceInMeters
			bestNeighborhoodIdx = i
		}
		i++
	}

	if len(matchedNeighborhoods) == 0 {
		return Neighborhood{}, err
	}

	return matchedNeighborhoods[bestNeighborhoodIdx], err
}

// Returns the coordinates of a MultiPolygon's centroid (if found). idx 0 => latitude, idx 1 => longitude
func resolveNeighborhoodMultiPolygonsCentroidPoint(
	neighborhoodName string,
	neighborhoodCity string,
	neighborhoodState string) ([]float64, error) {
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

	row := connections.Init().QueryRow(
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

// Get distance between two coordinate in meters.
// See: https://postgis.net/docs/manual-1.4/ST_Distance_Sphere.html
func getDistanceBetweenTwoCoordinates(point1 []float64, point2 []float64) (float64, error) {
	pointDistanceQueryStr := `
    SELECT ST_Distance_Sphere(
        ST_SetSRID(ST_Point($1, $2), 4326),
        ST_SetSRID(ST_Point($3, $4), 4326)
    ) as distance_in_meters`

	row := connections.Init().QueryRow(
		pointDistanceQueryStr,
		point1[0],
		point1[1],
		point2[0],
		point2[1])

	var distanceInMeters float64
	err := row.Scan(&distanceInMeters)

	if err != nil {
		log.Print(err)
		return 0.0, err
	}

	return distanceInMeters, err
}

// FindBestNeighborhood resolves the "best" neighborhood from the given array of neighborhoods.
// Best is defined here as:
// 	a) Having the highest occurrence (frequency)
//	b) Minimized distance between all other neighborhoods in the list
func FindBestNeighborhood(neighborhoods []Neighborhood) (Neighborhood, error) {
	neighborhoodNames, err := findNeighborhoodWithHighestOccurrence(neighborhoods)
	if err != nil {
		log.Printf("Unable to resolve neighborhoods with highest occurrence having error: %v\n", err)
		return Neighborhood{}, err
	}

	var highestOccurrenceNeighborhoods []Neighborhood
	for _, neighborhoodName := range neighborhoodNames {
		for _, neighborhood := range neighborhoods {
			if neighborhoodName == neighborhood.Name {
				highestOccurrenceNeighborhoods = append(highestOccurrenceNeighborhoods, neighborhood)
			}
		}
	}

	optimalNeighborhoodName, err := findNeighborhoodWithLeastDistanceToAllOtherNeighborhoods(highestOccurrenceNeighborhoods)

	if err == nil {
		return optimalNeighborhoodName, nil
	}

	return Neighborhood{}, &NoNeighborhoodFoundError{"Unable to resolve neighborhood after attempting to find best match."}
}

// NoNeighborhoodFoundError indicates a neighborhood was not resolved
type NoNeighborhoodFoundError struct {
	message string
}

func (e *NoNeighborhoodFoundError) Error() string {
	return e.message
}

// Finds the neighborhood which has the highest occurrence based on its name.
func findNeighborhoodWithHighestOccurrence(neighborhoods []Neighborhood) ([]string, error) {
	neighborhoodFrequency := make(map[string]int)

	// Construct frequency table
	for _, neighborhood := range neighborhoods {
		_, keyExists := neighborhoodFrequency[neighborhood.Name]

		if keyExists {
			neighborhoodFrequency[neighborhood.Name]++
		} else {
			neighborhoodFrequency[neighborhood.Name] = 1
		}
	}

	// Build a max-heap: O(n log(n)). We choose a heap to easily find all neighborhoods tying for the max
	// occurrence.
	h := getMaxHeap(neighborhoodFrequency)

	neighborhoodNames, err := findNeighborhoodsWithSameFrequency(h)
	if err != nil {
		log.Printf("Unable to resolve neighborhoods with the same frequency; having error: %v\n", err)
		return []string{}, err
	}

	return neighborhoodNames, nil
}

// findNeighborhoodsWithSameFrequency returns all neighborhoods that have the same number of entries.
// Example: {"Downtown": 4, "Southside": 4, "East Bay": 4}
func findNeighborhoodsWithSameFrequency(h *neighborhoodNameFrequencyMaxHeap) ([]string, error) {
	if h.Len() == 0 {
		return []string{}, nil
	}

	if h.Len() == 1 {
		v := heap.Pop(h)
		return []string{v.(neighorboodNameFrequency).name}, nil
	}

	maxCount := 0
	var neighborhoodNames []string
	for i := 0; i < h.Len(); i++ {
		v := heap.Pop(h).(neighorboodNameFrequency)
		if v.count < maxCount {
			break
		} else {
			maxCount = v.count
			neighborhoodNames = append(neighborhoodNames, v.name)
		}
	}

	return neighborhoodNames, nil
}

// Caching PostGIS calculations on geometric objects is desired as they're computationally, and time expensive.
func generateNeighborhoodCacheKey(neighborhoodName string, otherNeighborhoodName string) string {
	// sort to ensure we always get the same hash value for the same two keys.
	cacheKeyElements := []string{neighborhoodName, otherNeighborhoodName}
	sort.Strings(cacheKeyElements)
	cacheKey := strings.Join(cacheKeyElements, "")

	hasher := md5.New()
	hasher.Write([]byte(cacheKey))
	hashedString := hex.EncodeToString(hasher.Sum(nil))

	return hashedString
}
