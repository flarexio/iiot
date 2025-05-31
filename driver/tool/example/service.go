package example

import (
	"context"
	"errors"
	"sync"

	"github.com/flarexio/iiot/driver"
	"github.com/flarexio/iiot/machine"
)

type Controller struct {
	Points map[string]any `json:"points"`
}

type Service interface {
	driver.Service
}

func NewService() Service {
	return &service{
		controllers: make(map[string]*Controller),
	}
}

type service struct {
	controllers map[string]*Controller
	sync.RWMutex
}

func (svc *service) AddControllers(controllers ...*machine.Controller) error {
	svc.Lock()
	defer svc.Unlock()

	for _, controller := range controllers {
		points := make(map[string]any)
		for _, point := range controller.Points {
			value, ok := point.Options["value"]
			if !ok {
				return errors.New("point value not found in options")
			}

			points[point.Name] = value
		}

		c := &Controller{points}
		svc.controllers[controller.ControllerID] = c
	}

	return nil
}

func (svc *service) ReadPoints(ctx context.Context, id string, pointNames []string) ([]any, error) {
	svc.RLock()
	defer svc.RUnlock()

	c, ok := svc.controllers[id]
	if !ok {
		return nil, driver.ErrControllerNotFound
	}

	points := make([]any, len(pointNames))
	for i, name := range pointNames {
		value, ok := c.Points[name]
		if !ok {
			return nil, driver.ErrPointNotFound
		}

		points[i] = value
	}

	return points, nil
}
