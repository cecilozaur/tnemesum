# tnemesum

###Alexandru Mihalcea

### Project Overview
The project is a Go application. When the app runs it preloads all cities and weather info from the two APIs and stores the items in an inmemory map
### Prerequisites

- Go

### Run with Docker
```
docker build -t tnemesum .
docker run -it --rm -p 8000:8000 --name tnemesum-app tnemesum
```
or use the makefile
```
make run
```

### API endpoint
The API service runs on `localhost:8000`
```
GET /api/cities
```
Retrieves the list of all cities with weather information.
```
GET /api/cities/{cityId}
```
Retrieves a city with weather information for all the dates.

## API design

### Get the forecast for a city
```
GET /api/forecast/{cityId}
```
Retrieves the forecast for the `cityId` for the next 2 days

#####Responses
An object of the following format
```
{
    "forecast": [
        {
            "day": "2021-07-09",
            "condition": "Partly cloudy"
        },
        {
            "day": "2021-07-10",
            "condition": "Partly cloudy"
        }
    ]
}
```

- success
  Returns 200 status code, and the JSON body.
- error
    - Returns 400 status code if the `cityId` parameter is not a number.
    - Returns 404 status code if the city `cityId` is not found.

### Get the forecast for a city
```
GET /api/forecast/{cityId}/{day}
```
Retrieves the forecast for the `cityId` for the specified `day`

#####Response
An object of the following format
```
{
    "forecast": {
        "day": "2021-07-09",
        "condition": "Partly cloudy"
    }
}
```

- success
  Returns 200 status code, and the JSON body.
- error
    - Returns 400 status code if the `cityId` parameter is not a number.
    - Returns 400 status code if the `day` parameter is not a valid date in the format `YYYY-mm-dd`.
    - Returns 404 status code if the forecast for the specified `day` is not found.

### Set the forecast for a range of days or a single date for a given city.
```
POST /api/forecast/{cityId}
```
Create or update the forecast in the city identified by `cityId`.
#####Input JSON
```
[
    {
        "day": "2021-07-09",
        "condition": "Sunny"
    },
    {
        "day": "2021-07-10",
        "condition": "Still sunny"
    }
]
```

#####Response codes

- success
  Returns 204 (No content) status code along with a Location header for the updated resource.
- error
    - Returns 400 status code if the provided JSON input is invalid, or the cityId parameter is not a number
    - Returns 404 status code if the city is not found.
    - Returns 502 status code if the update failed