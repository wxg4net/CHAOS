package device

import (
	"errors"
	"time"

	"github.com/tiagorlampert/CHAOS/entities"
	"github.com/tiagorlampert/CHAOS/internal/utils"
	"github.com/tiagorlampert/CHAOS/repositories"
	"github.com/tiagorlampert/CHAOS/repositories/device"
)

type deviceService struct {
	Repository device.Repository
}

func NewDeviceService(repository device.Repository) Service {
	return &deviceService{Repository: repository}
}

func (d deviceService) UpdateDeviceName(input entities.UDeviceName) error {
	input.UpdatedAt = time.Now().UTC()

	device, err := d.Repository.FindByMacAddress(input.MacAddress)
	if errors.Is(err, repositories.ErrNotFound) {
		return err
	}
	device.Devicename = input.Devicename
	return d.Repository.Update(*device)
}

func (d deviceService) Delete(input entities.UDeviceName) error {
	device, err := d.Repository.FindByMacAddress(input.MacAddress)
	if errors.Is(err, repositories.ErrNotFound) {
		return err
	}
	return d.Repository.Delete(*device)
}

func (d deviceService) Insert(input entities.Device) error {
	input.UpdatedAt = time.Now().UTC()

	_, err := d.Repository.FindByMacAddress(input.MacAddress)
	if errors.Is(err, repositories.ErrNotFound) {
		return d.Repository.Insert(input)
	}
	return d.Repository.Update(input)
}

func (d deviceService) FindAllConnected() ([]entities.Device, error) {
	until := time.Now().UTC().
		Add(time.Minute * time.Duration(-3))

	devices, err := d.Repository.FindAll(until)
	if err != nil {
		return nil, err
	}

	for index, entity := range devices {
		devices[index].MacAddressBase64 = utils.EncodeBase64(entity.MacAddress)
	}
	return devices, nil
}

func (d deviceService) FindAllExisted() ([]entities.Device, error) {

	devices, err := d.Repository.FindAllDevices()
	if err != nil {
		return nil, err
	}

	for index, entity := range devices {
		devices[index].MacAddressBase64 = utils.EncodeBase64(entity.MacAddress)
	}
	return devices, nil
}

func (d deviceService) FindByMacAddress(address string) (*entities.Device, error) {
	return d.Repository.FindByMacAddress(address)
}
