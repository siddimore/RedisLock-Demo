package internal

// Ride represents a ride request and its current assignment state.
type Ride struct {
	ID       string
	DriverID string
	Status   string
}
