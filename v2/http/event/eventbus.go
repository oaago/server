package event

import (
	evbus "github.com/asaskevich/EventBus"
	"github.com/oaago/server/v2/types"
)

func NewEvent() types.Event {
	return evbus.New()
}
