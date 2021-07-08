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

	slower chan bool // barrier we'll use to throttle the number of concurrent request we do to the weather API
}

func New(config conf.Config) *App {
	return &App{
		cfg:       config,
		startTime: time.Now().UTC(),
		status:    0,
		repo:      repository.NewInMem(),
		slower:    make(chan bool, 2), // TODO this should probably be in the config
	}
}

func (a *App) Run() {
	a.loadCities()
	a.printWeatherInfo()
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

func (a *App) GetAllCities() []domain.City {
	return a.repo.GetItems()
}

func (a *App) loadCities() {
	// get the list of cities from musement API
	result, err := http.Get(a.cfg.MusementAPIUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer result.Body.Close()

	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		log.Fatal(err)
	}
	items := make([]domain.City, 0)
	err = json.Unmarshal(body, &items)
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range items {
		a.repo.Store(item.Id, item)
	}
}

func (a *App) printWeatherInfo() {
	wg := &sync.WaitGroup{}
	for _, city := range a.repo.GetItems() {
		go func(c domain.City) {
			// throttle calls to weather API?
			a.slower <- true
			wg.Add(1)
			defer func() {
				<-a.slower
				wg.Done()
			}()
			forecast := a.getWeatherForCity(c.Lat, c.Lng)

			// only print the message if we know we got the forecast for at least 2 days
			if len(forecast.Forecast.ForecastDay) >= 2 {
				log.Printf("Processed city %s | %s - %s", c.Name, forecast.Forecast.ForecastDay[0].Day.Condition.Text, forecast.Forecast.ForecastDay[1].Day.Condition.Text)
			}
			a.repo.UpdateForecast(c.Id, forecast.Forecast)
		}(city)
	}

	wg.Wait()
}

func (a *App) getWeatherForCity(lat, lng float64) domain.ForecastResult {
	url := fmt.Sprintf(a.cfg.WeatherAPIUrl, a.cfg.WeatherAPIKey, lat, lng)
	result, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer result.Body.Close()

	f := domain.ForecastResult{}
	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(body, &f)
	if err != nil {
		log.Fatal(err)
	}

	return f
}
