package conf

type Config struct {
	WeatherAPIUrl         string
	MusementAPIUrl        string
	WeatherAPIKey         string
	APIConcurrentRequests int
	MaxConcurrent         int
}
