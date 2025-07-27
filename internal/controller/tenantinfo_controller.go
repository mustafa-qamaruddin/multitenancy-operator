/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	multitenancymanagementv1 "github.com/mustafa-qamaruddin/multitenancy-operator/api/v1"
	v1 "github.com/mustafa-qamaruddin/multitenancy-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TenantInfoReconciler reconciles a TenantInfo object
type TenantInfoReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=multitenancy-management.example.com,resources=tenantinfoes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=multitenancy-management.example.com,resources=tenantinfoes/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=multitenancy-management.example.com,resources=tenantinfoes/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the TenantInfo object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile
func (r *TenantInfoReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	tenantInfo := &v1.TenantInfo{}
	err := r.Get(ctx, req.NamespacedName, tenantInfo)
	if err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("TenantInfo resource not found. Assuming it was deleted.")
			// Cleanup is handled by Kubernetes garbage collection if owner references are set.
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get TenantInfo")
		return ctrl.Result{}, err
	}

	// Step 1: Track current tenants
	currentTenantIDs := make(map[string]bool)
	for _, tenant := range tenantInfo.Spec.Tenants {
		currentTenantIDs[tenant.TenantID] = true

		// Create or update ConfigMap
		cmName := fmt.Sprintf("tenant-%s-config", tenant.TenantID)
		cm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      cmName,
				Namespace: req.Namespace,
			},
			Data: map[string]string{
				"tenantID":      tenant.TenantID,
				"webserviceURL": tenant.WebserviceURL,
			},
		}
		if err := ctrl.SetControllerReference(tenantInfo, cm, r.Scheme); err != nil {
			log.Error(err, "Failed to set controller reference")
			return ctrl.Result{}, err
		}

		found := &corev1.ConfigMap{}
		err := r.Get(ctx, types.NamespacedName{Name: cm.Name, Namespace: cm.Namespace}, found)
		if err != nil && apierrors.IsNotFound(err) {
			log.Info("Creating ConfigMap for tenant", "tenantID", tenant.TenantID)
			if err := r.Create(ctx, cm); err != nil {
				log.Error(err, "Failed to create ConfigMap", "tenantID", tenant.TenantID)
				return ctrl.Result{}, err
			}
		} else if err == nil && !reflect.DeepEqual(found.Data, cm.Data) {
			found.Data = cm.Data
			log.Info("Updating ConfigMap for tenant", "tenantID", tenant.TenantID)
			if err := r.Update(ctx, found); err != nil {
				log.Error(err, "Failed to update ConfigMap", "tenantID", tenant.TenantID)
				return ctrl.Result{}, err
			}
		} else if err != nil {
			log.Error(err, "Failed to get ConfigMap")
			return ctrl.Result{}, err
		}
	}

	// Step 2: List all ConfigMaps owned by this TenantInfo and delete the ones no longer needed
	var childCMs corev1.ConfigMapList
	if err := r.List(ctx, &childCMs, client.InNamespace(req.Namespace)); err != nil {
		log.Error(err, "Failed to list child ConfigMaps")
		return ctrl.Result{}, err
	}

	// Manually filter by ownerReference UID since fake client doesn't support field selectors
	var filteredChildCMs []corev1.ConfigMap
	for _, cm := range childCMs.Items {
		for _, ref := range cm.OwnerReferences {
			if ref.UID == tenantInfo.UID {
				filteredChildCMs = append(filteredChildCMs, cm)
				break
			}
		}
	}

	for _, cm := range filteredChildCMs {
		tenantID := strings.TrimPrefix(cm.Name, "tenant-")
		tenantID = strings.TrimSuffix(tenantID, "-config")
		if _, exists := currentTenantIDs[tenantID]; !exists {
			log.Info("Deleting orphaned ConfigMap", "name", cm.Name)
			if err := r.Delete(ctx, &cm); err != nil {
				log.Error(err, "Failed to delete ConfigMap", "name", cm.Name)
				return ctrl.Result{}, err
			}
		}
	}

	return ctrl.Result{}, nil
}

func (r *TenantInfoReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// Register field index for ownerReferences.uid
	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &corev1.ConfigMap{}, "metadata.ownerReferences.uid", func(rawObj client.Object) []string {
		ownerRefs := rawObj.GetOwnerReferences()
		for _, ref := range ownerRefs {
			return []string{string(ref.UID)}
		}
		return nil
	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&multitenancymanagementv1.TenantInfo{}).
		Named("tenantinfo").
		Complete(r)
}
