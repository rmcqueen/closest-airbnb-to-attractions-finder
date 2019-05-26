# Overview

This application enables a user to easily find a list of the closest AirBNBs based on the attractions they have chosen to visit.  Due to the constraints AirBNB imposes on a posting's address, the closest region that can be identified is a neighborhood within a city.

# Motivation

While AirBNB provides the means to select attractions near-by a listing, I personaly choose my attractions prior to picking where I sleep (so I am building a tool to do just this).

# Idea
1. ✅ Take a set of attractions and determine their coordinates
2. Relate each attraction's coordinates to a neighborhood within the same city
3. Construct a frequency table of where the key is the neighborhood name, and the value is the number of times it has appeared based on the attractions. For example:
```
[
    "Dunbar": 3,
    "Kitsilano": 2,
    ...,
    "Marpole": 0
]
```
4. Find the neighborhood which contains a majority of the attractions, and minimizes the distance to the other attractions within other neighborhoods
5. Look-up all related AirBNBs for the matched neighborhood and apply some filtering criterion based on personal preferences (i.e, "must be < $90 a night, be an apartment, and not shared.")
6. Construct a list of matches and optimize based on the following orderings (first item has highest priority):
    1. Cost per night
    2. Not shared
    3. Building type

Unfortunately, AirBNB policy does not enable us to factor in distance to this optimization as they keep exact addresses private until booking. The best we can do is optimize within the best neighborhood.

### Usage
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
