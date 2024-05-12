minikube start --cpus=2 --memory=2040 --disk-size "10 GB" --vm-driver=virtualbox 
kubectl apply -f ./k8s/.
minikube dashboard