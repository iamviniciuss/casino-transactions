docker build -t vsantos98/cassino-api:latest .
docker push vsantos98/cassino-api:latest

kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/hpa.yaml
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml

kubectl port-forward svc/cassino-api-service 8080:80
