# Project Status: [WIP]

# Overview

This application enables a user to easily find a list of the closest AirBNBs based on the attractions they have chosen to visit.  Due to the constraints AirBNB imposes on a posting's address, the closest region that can be identified is a neighborhood within a city.  I am primarily developing this to gain a basic understanding of golang.  The  code is likely to not align very well with golang's best practices.

# Motivation

While AirBNB provides the means to select attractions near-by a listing, I personaly choose my attractions prior to picking where I sleep (so I am building a tool to do just this).

# Idea
1. ✅ Take a set of attractions and determine their coordinates
2. ✅ Relate each attraction's coordinates to a neighborhood within the same city
3. ✅ Construct a frequency table of where the key is the neighborhood name, and the value is the number of times it has appeared based on the attractions. For example:
```
{
    "Dunbar": 3,
    "Kitsilano": 2,
    ...,
    "Marpole": 0
}
```
4. ✅ Find the neighborhood which contains a majority of the attractions, and minimizes the distance to the other attractions within other neighborhoods
5. Look-up all related AirBNBs for the matched neighborhood and apply some filtering criterion based on personal preferences (i.e, "must be < $90 a night, be an apartment, and not shared.")
6. Construct a list of matches and optimize based on the following orderings (first item has highest priority):
    1. Cost per night
    2. Not shared
    3. Building type

Unfortunately, AirBNB policy does not enable us to factor in distance to this optimization as they keep exact addresses private until booking. The best we can do is optimize within the best neighborhood.

### Requirements

A file titled `init.sql` is required. This file should perform the following actions:
1. Populate a schema and table to store the neighborhood geocodings:

        CREATE EXTENSION IF NOT EXISTS postgis;
        SET CLIENT_ENCODING TO UTF8;
        SET STANDARD_CONFORMING_STRINGS TO ON;
        CREATE SCHEMA neighborhood_geocoding;
        BEGIN;
        CREATE TABLE "neighborhood_geocoding"."neighborhoods" (gid serial,
        "name" varchar(254),
        "city" varchar(80),
        "state" varchar(80),
        "country" varchar(80),

        UNIQUE (name, city, state)
        );
        ALTER TABLE "neighborhood_geocoding"."neighborhoods" ADD PRIMARY KEY (gid);
        SELECT AddGeometryColumn('neighborhood_geocoding','neighborhoods','geom','4326','MULTIPOLYGON',2);
2. Insert some neighborhood multipolygons
    - Note: you will have to resolve this yourself as insert files occupy too much space on GitHub. These are typically located within `.shp` files and can be found from a local government Open Data portal. There exists a tool, `shp2pgsql` which will convert these into valid PostgreSQL insert statements for you.
 
### Usage
1. Build and run the application:
    ```
    cd cmd
    go build <some_binary_file_name>
    ./<some_binary_file_name>
    ```

    By default, the application runs on port 8080.

2. Pass a list of attractions via a POST request to `/attractions`

    - A JSON array is expected to be passed to the `/attractions` endpoint:
    ```
    [
        {
            "name": "",
            "city": "",
            "state_or_province_name": "",
        }
    ]
    ```

    The result looks as follows:
    ```
    {
        "successful_attractions": [
            {
                "name": "",
                "city": "",
                "state_or_province_name": ""
                "latitude": 0.0,
                "longitude": 0.0
            }
        ],
        "failed_attractions": [
            {
                "name": "",
                "city": "",
                "state": "",
                "latitude": 0.0,
                "longitude": 0.0
            }
        ],
        "ClosestNeighborhood": {
            "name": "",
            "city_name": "",
            "latitude": 0.0,
            "longitude": 0.0
        }
    }
    ```

    **Note**: In the event either all attractions are unsuccessfully geocoded, or all attractions are successfully geocoded, the `*_attractions` key may be null.
