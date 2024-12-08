package seed

import (
	"time"

	"github.com/fwaters2/launch-schedule-manager/server/pkg/launches"
)

var Launches = []launches.Launch{
	{
		ID:          "1",
		MissionName: "Test Mission A",
		LaunchTime:  time.Date(2024, 5, 20, 14, 0, 0, 0, time.UTC),
		VehicleName: "Falcon 9",
		LaunchSite:  "LC-39A",
		Status:      "scheduled",
	},
	{
		ID:          "2",
		MissionName: "Test Mission B",
		LaunchTime:  time.Now().Add(24 * time.Hour),
		VehicleName: "Starship",
		LaunchSite:  "Boca Chica",
		Status:      "pending",
	},
}
