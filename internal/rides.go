package internal

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

var rides = make(map[string]*Ride)
var ridesMu sync.RWMutex

// AcceptRide accepts a ride for a driver using distributed and in-process locks.
func AcceptRide(w http.ResponseWriter, r *http.Request, rclient *RedisClient) {
	ctx := context.Background()
	rideID := r.URL.Query().Get("ride_id")
	driverID := r.URL.Query().Get("driver_id")

	if rideID == "" || driverID == "" {
		http.Error(w, "ride_id and driver_id are required", http.StatusBadRequest)
		return
	}

	lockKey := "ride_lock:" + rideID
	locked, err := rclient.AcquireLock(ctx, lockKey, driverID, 10*time.Second)
	if err != nil {
		http.Error(w, "could not acquire distributed lock", http.StatusInternalServerError)
		return
	}
	if !locked {
		http.Error(w, fmt.Sprintf("ride %s is already being processed", rideID), http.StatusConflict)
		return
	}

	ridesMu.Lock()
	ride, ok := rides[rideID]
	if !ok {
		ride = &Ride{ID: rideID, Status: "pending"}
		rides[rideID] = ride
	}
	if ride.Status == "accepted" {
		ridesMu.Unlock()
		http.Error(w, fmt.Sprintf("ride %s already accepted by driver %s", rideID, ride.DriverID), http.StatusConflict)
		return
	}

	ride.DriverID = driverID
	ride.Status = "accepted"
	ridesMu.Unlock()

	rclient.ReleaseLock(ctx, lockKey)

	fmt.Fprintf(w, "ride %s accepted by driver %s", rideID, driverID)
}
