# Setup

# Initialize the Kubebuilder project with your chosen domain and repo path
kubebuilder init --domain example.com --repo github.com/mustafa-qamaruddin/multitenancy-operator

# Create a new API group/version/kind (CRD: TenantInfo under multitenancy-management.example.com/v1)
kubebuilder create api --group multitenancy-management --version v1 --kind TenantInfo

# Generate boilerplate code (types, controller, etc.)
make generate

# Generate CRD manifests and RBAC configurations
make manifests

# Install the CRD into the Kubernetes cluster
make install


# Run Demo

# Start the Kubernetes dashboard (useful for watching resources in Minikube)
minikube dashboard

# Run the operator locally (will watch TenantInfo resources and reconcile them)
make run

# Apply the sample custom resource defined in config/samples/
kubectl apply -k config/samples

# View all created TenantInfo custom resources across namespaces
kubectl get tenantinfoes -A

# Delete the sample custom resource to test cleanup logic
kubectl delete tenantinfoes/tenantinfo-sample


# Cleanup

# Uninstall the CRDs and related cluster resources installed by the operator
make uninstall
