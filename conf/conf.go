package conf

type Config struct {
	WeatherAPIUrl         string // the URL or the API of weatherapi, used to load forecast data
	MusementAPIUrl        string // the URL of musement API, used to load cities
	WeatherAPIKey         string // the API key for weatherapi
	APIConcurrentRequests int    // the number of concurrent requests sent to weatherapi API to retrieve forecasts
	MaxConcurrent         int    // the number of max concurrent the service will allow, requests over this threshold will be dropped with a 429 code
}
