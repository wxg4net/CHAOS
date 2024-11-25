package device

import (
	"time"

	"github.com/tiagorlampert/CHAOS/entities"
)

type Repository interface {
	Insert(device entities.Device) error
	Update(device entities.Device) error
	Delete(device entities.Device) error
	FindByMacAddress(address string) (*entities.Device, error)
	FindAll(fetchedAt time.Time) ([]entities.Device, error)
	FindAllDevices() ([]entities.Device, error)
}
