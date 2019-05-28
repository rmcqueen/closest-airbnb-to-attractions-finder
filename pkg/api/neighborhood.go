package api

import (
    "database/sql"
    "fmt"
    "log"
    _ "github.com/lib/pq"
)

type Neighborhood struct {
    Name string `json:"name"`
    City string `json:"city_name"`
    Latitude float64 `json:"latitude"`
    Longitude float64 `json:"longitude"`
}

func FindNeighborhoodContainingAttraction(attraction Attraction, db sql.DB) {
    attractionInNeighborhoodQuery := `
        SELECT ST_Contains(neighborhood_poly, attr_point) as in_neighborhood, name
        FROM (
            SELECT ST_SetSRID(ST_Point($1, $2),4326) as attr_point, geom as neighborhood_poly, name
            FROM neighborhood_geocoding.neighborhoods
        ) as foo
        WHERE ST_Contains(neighborhood_poly, attr_point) is true
        `

    rows, err := db.Query(attractionInNeighborhoodQuery, attraction.Longitude, attraction.Latitude)
    if err != nil {
        log.Fatal(err)
    }

    defer rows.Close()
    for rows.Next() {
        var name string
        var in_neighborhood bool
        if err := rows.Scan(&in_neighborhood, &name); err != nil {

        }
        if len(name) == 0 {
            fmt.Println("Name is empty\n")
            fmt.Printf("Is empty: %t\n", in_neighborhood)
        } else {
            fmt.Printf("Name: %s\n", name)
        }
    }


    // TODO: Get the neighborhood geom's centroid point's lat/lng and construct a Neighborhood struct.

}
