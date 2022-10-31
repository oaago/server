package v2

import (
	evbus "github.com/asaskevich/EventBus"
)

type Event evbus.Bus
type BusPublisher evbus.BusPublisher
type BusSubscriber evbus.BusSubscriber
type BusController evbus.BusController

var EventBus Event

func init() {
	EventBus = evbus.New()
}
