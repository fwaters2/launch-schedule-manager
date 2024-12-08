package launches

import (
	"time"
)

type Launch struct {
	ID          string    `json:"id"`
	MissionName string    `json:"mission_name"`
	LaunchTime  time.Time `json:"launch_time"`
	VehicleName string    `json:"vehicle_name"`
	LaunchSite  string    `json:"launch_site"`
	Status      string    `json:"status,omitempty"`
}

type LaunchCreateRequest struct {
	MissionName string `json:"mission_name"`
	LaunchTime  string `json:"launch_time"` // RFC3339 format expected
	VehicleName string `json:"vehicle_name"`
	LaunchSite  string `json:"launch_site"`
	Status      string `json:"status,omitempty"`
}

func ValidateLaunchRequest(req LaunchCreateRequest) error {
	if req.MissionName == "" ||
		req.LaunchTime == "" ||
		req.VehicleName == "" ||
		req.LaunchSite == "" {
		return ErrInvalidInput
	}

	_, err := time.Parse(time.RFC3339, req.LaunchTime)
	if err != nil {
		return ErrInvalidTimeFormat
	}
	return nil
}
