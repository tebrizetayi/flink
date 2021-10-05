package app

import "context"

type Location struct {
	Lat float64
	Lng float64
}

type App struct {
	Orders map[string][]Location
}

func NewApp() App {
	return App{
		Orders: map[string][]Location{},
	}
}

//Save location insert location front of the order history.
func (a *App) SaveLocation(ctx context.Context, l Location, orderId string) error {
	if _, ok := a.Orders[orderId]; !ok {
		a.Orders[orderId] = []Location{}
	}
	temp := []Location{l}
	a.Orders[orderId] = append(temp, a.Orders[orderId]...)
	return nil
}

//Deletelocation remove history of order from dictionary
func (a *App) DeleteLocation(ctx context.Context, orderId string) error {
	delete(a.Orders, orderId)
	return nil
}

//ReadLocations reads the location.
func (a *App) ReadLocations(ctx context.Context, orderId string, max int) ([]Location, error) {

	n := max
	if max > len(a.Orders[orderId]) || max == 0 {
		n = len(a.Orders[orderId])
	}
	locations := a.Orders[orderId][:n]
	if locations == nil {
		locations = []Location{}
	}
	return locations, nil
}
