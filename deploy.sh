#!/usr/bin/env bash
set -e

echo "==> Pointing Docker to minikube's daemon..."
eval $(minikube docker-env)

echo "==> Building Docker image inside minikube..."
docker build -t distributed-lock-demo:latest .

echo "==> Applying Redis manifests..."
kubectl apply -f k8s/redis.yaml

echo "==> Waiting for Redis to be ready..."
kubectl rollout status deployment/redis

echo "==> Applying app manifests..."
kubectl apply -f k8s/app.yaml

echo "==> Waiting for app to be ready..."
kubectl rollout status deployment/distributed-lock-demo

echo ""
echo "==> Done. Get the service URL with:"
echo "    minikube service distributed-lock-demo --url"
echo ""
echo "==> Test concurrent drivers with:"
echo "    URL=\$(minikube service distributed-lock-demo --url)"
echo "    seq 1 30 | xargs -I{} -P 30 curl -s \"\$URL/accept?ride_id=ride-101&driver_id=driver-{}\""
