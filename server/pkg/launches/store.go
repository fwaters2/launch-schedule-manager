package launches

import (
	"errors"
	"fmt"
	"sync"
)

var (
	ErrNotFound         = errors.New("launch not found")
	ErrInvalidInput     = errors.New("invalid input data")
	ErrInvalidTimeFormat = errors.New("invalid time format, must be RFC3339")
)

// Store interface for managing launches
type Store interface {
	Create(l Launch) (Launch, error)
	Get(id string) (Launch, error)
	List() ([]Launch, error)
	Update(id string, l Launch) (Launch, error)
	Delete(id string) error
}

// InMemoryStore implementation
type InMemoryStore struct {
	mu      sync.Mutex
	launches map[string]Launch
	nextID  int
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		launches: make(map[string]Launch),
		nextID:   1,
	}
}

func (s *InMemoryStore) Create(l Launch) (Launch, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	l.ID = fmt.Sprintf("%d", s.nextID)
	s.nextID++
	s.launches[l.ID] = l
	return l, nil
}

func (s *InMemoryStore) Get(id string) (Launch, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	l, ok := s.launches[id]
	if !ok {
		return Launch{}, ErrNotFound
	}
	return l, nil
}

func (s *InMemoryStore) List() ([]Launch, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := make([]Launch, 0, len(s.launches))
	for _, l := range s.launches {
		result = append(result, l)
	}
	return result, nil
}

func (s *InMemoryStore) Update(id string, updated Launch) (Launch, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.launches[id]
	if !ok {
		return Launch{}, ErrNotFound
	}
	// Update the record
	launch := s.launches[id]
	if updated.MissionName != "" {
		launch.MissionName = updated.MissionName
	}
	if !updated.LaunchTime.IsZero() {
		launch.LaunchTime = updated.LaunchTime
	}
	if updated.VehicleName != "" {
		launch.VehicleName = updated.VehicleName
	}
	if updated.LaunchSite != "" {
		launch.LaunchSite = updated.LaunchSite
	}
	if updated.Status != "" {
		launch.Status = updated.Status
	}
	s.launches[id] = launch

	// Retrieve the updated record to return
	final := s.launches[id]
	return final, nil
}

func (s *InMemoryStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.launches[id]
	if !ok {
		return ErrNotFound
	}
	delete(s.launches, id)
	return nil
}
