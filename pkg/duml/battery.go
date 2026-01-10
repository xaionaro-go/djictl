package duml

import (
	"context"
	"fmt"
	"math"
)

type BatteryCapacity int8

const (
	UndefinedBatteryCapacity = BatteryCapacity(math.MinInt8)
)

func (b BatteryCapacity) String() string {
	if b == UndefinedBatteryCapacity {
		return "<undefined>"
	}
	return fmt.Sprintf("%d%%", int8(b))
}

type BatteryStatus struct {
	Capacity BatteryCapacity
}

func ParseBatteryStatus(
	ctx context.Context,
	payload []byte,
) (*BatteryStatus, error) {
	if len(payload) < 13 {
		return nil, fmt.Errorf("payload is too short: %d < 13", len(payload))
	}

	var capacity BatteryCapacity
	switch {
	case len(payload) >= 21:
		capacity = BatteryCapacity(payload[20])
	case len(payload) >= 13:
		capacity = BatteryCapacity(payload[12])
	default:
		return nil, fmt.Errorf("unable to determine battery status; the payload is too short: %d", len(payload))
	}

	return &BatteryStatus{
		Capacity: capacity,
	}, nil
}
