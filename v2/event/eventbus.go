package event

import (
	evbus "github.com/asaskevich/EventBus"
)

type Event evbus.Bus
type BusPublisher evbus.BusPublisher
type BusSubscriber evbus.BusSubscriber
type BusController evbus.BusController

func NewEvent() Event {
	return evbus.New()
}
