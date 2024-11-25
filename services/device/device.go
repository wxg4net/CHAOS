package device

import (
	"github.com/tiagorlampert/CHAOS/entities"
)

type Service interface {
	UpdateDeviceName(entities.UDeviceName) error
	Delete(entities.UDeviceName) error
	Insert(entities.Device) error
	FindAllConnected() ([]entities.Device, error)
	FindAllExisted() ([]entities.Device, error)
	FindByMacAddress(address string) (*entities.Device, error)
}
