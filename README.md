# Distributed Lock Demo

A Go HTTP service that demonstrates distributed locking using Redis `SetNX` to prevent multiple drivers from accepting the same ride simultaneously.

## How it works

1. When a driver calls `/accept`, the app tries to acquire a Redis lock on `ride_lock:<ride_id>` using `SetNX`.
2. Only the first driver to set the key wins the ride.
3. All other concurrent drivers get a `409 Conflict` immediately.
4. Once the ride is accepted, the lock is released and in-memory state prevents re-acceptance.
5. TTL (10s) ensures the lock auto-expires if the process crashes before releasing.

```
Driver-1 ──┐
Driver-2 ──┤──> SetNX ride_lock:ride-101 ──> Only one wins ──> ride accepted
Driver-3 ──┘                                  Others rejected
```

## Project structure

```
cmd/main.go              — HTTP server entry point
internal/model.go        — Ride struct
internal/redis_lock.go   — RedisClient with AcquireLock / ReleaseLock
rideshare/rides.go       — AcceptRide handler with global map mutex
k8s/redis.yaml           — Redis Deployment + Service
k8s/app.yaml             — App Deployment (2 replicas) + NodePort Service
Dockerfile               — Multi-stage build
deploy.sh                — Build and deploy to minikube
```

## Run locally

**Prerequisites:** Docker running with a Redis container.

```bash
docker run --name redis-dev -p 6379:6379 -d redis:7
go run cmd/main.go
```

Test with 30 concurrent drivers:
```bash
seq 1 30 | xargs -I{} -P 30 curl -s "http://localhost:8080/accept?ride_id=ride-101&driver_id=driver-{}"
```

Expected output: exactly one success, all others conflict.

## Run on minikube

**Prerequisites:** minikube running (`minikube start`).

```bash
./deploy.sh
```

Get the service URL and test:
```bash
URL=$(minikube service distributed-lock-demo --url)
seq 1 30 | xargs -I{} -P 30 curl -s "$URL/accept?ride_id=ride-101&driver_id=driver-{}"
```

With 2 replicas, both pods share Redis so the distributed lock is the only thing preventing a double-accept across instances.

## Environment variables

| Variable     | Default         | Description              |
|--------------|-----------------|--------------------------|
| `REDIS_ADDR` | `localhost:6379` | Redis host and port       |
