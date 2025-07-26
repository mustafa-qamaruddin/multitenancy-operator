# Setup
kubebuilder init --domain example.com --repo github.com/mustafa-qamaruddin/multitenancy-operator
kubebuilder create api --group multitenancy-management --version v1 --kind TenantInfo
make generate
make manifests
make install


# Run Demo
minikube  dashboard
make run
kubectl apply -k config/samples

kubectl get tenantinfoes -A
kubectl delete tenantinfoes/tenantinfo-sample

# Cleanup
make uninstall
