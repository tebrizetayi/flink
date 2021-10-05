package app

import (
	"context"
	"sync"
	"time"
)

type Location struct {
	Lat float64
	Lng float64
}
type TimedLocation struct {
	Location
	Time time.Time
}

type App struct {
	Orders map[string][]TimedLocation
	TTL    int
	mux    sync.Mutex
}

func NewApp(ttl int) App {
	return App{
		Orders: map[string][]TimedLocation{},
		TTL:    ttl,
	}
}

//Save location insert location front of the order history.
func (a *App) SaveLocation(ctx context.Context, l Location, orderId string) error {

	a.mux.Lock()
	defer a.mux.Unlock()
	if _, ok := a.Orders[orderId]; !ok {
		a.Orders[orderId] = []TimedLocation{}
	}
	temp := []TimedLocation{
		{
			Location: l,
			Time:     time.Now(),
		},
	}
	a.Orders[orderId] = append(temp, a.Orders[orderId]...)
	return nil
}

//Deletelocation remove history of order from dictionary
func (a *App) DeleteLocation(ctx context.Context, orderId string) error {

	a.mux.Lock()
	defer a.mux.Unlock()
	delete(a.Orders, orderId)
	return nil
}

//ReadLocations reads the location.
func (a *App) ReadLocations(ctx context.Context, orderId string, max int) ([]Location, error) {

	a.mux.Lock()
	defer a.mux.Unlock()
	before := time.Now().Add(time.Second * (-time.Duration(a.TTL)))
	for i, v := range a.Orders[orderId] {
		if v.Time.Before(before) {
			a.Orders[orderId] = a.Orders[orderId][0:i]
			break
		}
	}

	n := max
	if max > len(a.Orders[orderId]) || max == 0 {
		n = len(a.Orders[orderId])
	}
	tLocations := a.Orders[orderId][:n]
	if tLocations == nil {
		tLocations = []TimedLocation{}
	}
	result := []Location{}

	for _, v := range tLocations {
		result = append(result, v.Location)
	}
	return result, nil
}
