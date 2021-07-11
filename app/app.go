package app

import (
	"encoding/json"
	"fmt"
	"github.com/cecilozaur/tnemesum/conf"
	"github.com/cecilozaur/tnemesum/domain"
	"github.com/cecilozaur/tnemesum/repository"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type App struct {
	cfg       conf.Config
	status    uint32
	startTime time.Time
	repo      repository.Repository

	semaphore chan bool // barrier we'll use to throttle the number of concurrent request we do to the weather API
}

func New(config conf.Config) *App {
	return &App{
		cfg:       config,
		startTime: time.Now().UTC(),
		status:    0,
		repo:      repository.NewInMem(),
		semaphore: make(chan bool, config.APIConcurrentRequests),
	}
}

// Run the function will initialize the repo by loading all city and forecast data and saving the info
// inside the provided repository (in this case the inmem one)
func (a *App) Run() {
	var items []domain.City
	var err error

	for i := 0; i < 3; i++ {
		items, err = a.loadCities()

		if err == nil {
			break
		}

		// wait 1 second between attempts?
		time.Sleep(time.Second)
	}

	if err != nil || len(items) == 0 {
		log.Fatal("unable to load cities from API")
	}

	// store items in repository
	for _, item := range items {
		a.repo.Store(item.Id, item)
	}

	a.loadAndStoreForecast()

	// print forecast after loading
	for _, city := range a.repo.GetItems() {
		// only print the message if we know we got the forecast for at least 2 days
		if len(city.Forecast) >= 2 {
			log.Printf("Processed city %s | %s - %s", city.Name, city.Forecast[0].Condition, city.Forecast[1].Condition)
		}
	}
}

func (a *App) Healthy() bool {
	return atomic.LoadUint32(&a.status) == 1
}

func (a *App) SetHealthy() {
	atomic.StoreUint32(&a.status, 1)
}

func (a *App) Pause() {
	atomic.StoreUint32(&a.status, 0)
}

func (a *App) Uptime() time.Duration {
	return time.Since(a.startTime)
}

// loadCities loads all cities returned by API and returns a slice of items or an error is the API call failed
func (a *App) loadCities() ([]domain.City, error) {
	// get the list of cities from musement API
	result, err := http.Get(a.cfg.MusementAPIUrl)
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()

	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}
	items := make([]domain.City, 0)
	err = json.Unmarshal(body, &items)
	if err != nil {
		return nil, err
	}

	return items, nil
}

// loadAndStoreForecast will start a new goroutine for every city we have in the repository
// and do a call to weatherapi API to get the forecast for the next 2 days
// it will also save the info inside the repository for each associated city
// the function uses a semaphore to limit the number of concurrent calls we make to the weatherapi API
func (a *App) loadAndStoreForecast() {
	wg := &sync.WaitGroup{}
	for _, city := range a.repo.GetItems() {
		go func(c domain.City) {
			// throttle calls to weather API?
			a.semaphore <- true
			wg.Add(1)
			defer func() {
				<-a.semaphore
				wg.Done()
			}()

			var forecast *ForecastResult
			var err error
			for i := 0; i < 3; i++ {
				forecast, err = a.getWeatherForCity(c.Lat, c.Lng)
				if err == nil {
					break
				}

				time.Sleep(time.Second)
			}

			if err != nil {
				log.Println("unable to load forecast for city " + err.Error())
				return
			}

			// build new forecast from retrieved data
			forecastDTO := make(domain.Forecast, 0)
			for _, f := range forecast.Forecast.ForecastDay {
				newF := domain.ForecastDay{
					Day:       f.Date,
					Condition: f.Day.Condition.Text,
				}
				forecastDTO = append(forecastDTO, newF)
			}

			a.repo.UpdateForecast(c.Id, forecastDTO)
		}(city)
	}

	wg.Wait()
}

// getWeatherForCity does a GET request for a specific lat/lng combination
// and loads the weather forecast for the next 2 days
func (a *App) getWeatherForCity(lat, lng float64) (*ForecastResult, error) {
	url := fmt.Sprintf(a.cfg.WeatherAPIUrl, a.cfg.WeatherAPIKey, lat, lng)
	result, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()

	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}

	f := ForecastResult{}
	err = json.Unmarshal(body, &f)
	if err != nil {
		return nil, err
	}

	return &f, nil
}
