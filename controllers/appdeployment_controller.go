/*
Copyright 2022.

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

package controllers

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	runtimev1alpha1 "github.com/tjololo/altinn-applight-operator/api/v1alpha1"
)

// AppDeploymentReconciler reconciles a AppDeployment object
type AppDeploymentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=runtime.applight.418.cloud,resources=appdeployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=runtime.applight.418.cloud,resources=appdeployments/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=runtime.applight.418.cloud,resources=appdeployments/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the AppDeployment object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.1/pkg/reconcile
func (r *AppDeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var appDeployment runtimev1alpha1.AppDeployment
	if err := r.Get(ctx, req.NamespacedName, &appDeployment); err != nil {
		if !errors.IsNotFound(err) {
			logger.Error(err, "unable to fetch AppDeployment")
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	var policyConfigMap corev1.ConfigMap
	err := r.Get(ctx, types.NamespacedName{Namespace: appDeployment.Namespace, Name: fmt.Sprintf("policy-%s", appDeployment.Name)}, &policyConfigMap)
	if client.IgnoreNotFound(err) != nil {
		logger.Error(err, "unable to fetch policy configmap")
		return ctrl.Result{}, err
	}
	if errors.IsNotFound(err) {
		policyConfigMap = getConfigMapForApp(&appDeployment, "policy", map[string]string{"policy.xml": appDeployment.Spec.Policy})
		err = r.Create(ctx, &policyConfigMap)
		if err != nil {
			logger.Error(err, "failed to create policy configmap")
			return ctrl.Result{}, err
		}
	}
	if policyConfigMap.Data["policy.xml"] != appDeployment.Spec.Policy {
		logger.Info("Updating policy configmap")
		policyConfigMap.Data["policy.xml"] = appDeployment.Spec.Policy
		err = r.Update(ctx, &policyConfigMap)
		if err != nil {
			logger.Error(err, "failed to update policy configmap")
			return ctrl.Result{}, err
		}
	}
	var processConfigmap corev1.ConfigMap
	err = r.Get(ctx, types.NamespacedName{Namespace: appDeployment.Namespace, Name: fmt.Sprintf("bpmn-%s", appDeployment.Name)}, &processConfigmap)
	if client.IgnoreNotFound(err) != nil {
		logger.Error(err, "unable to fetch bpmn configmap")
		return ctrl.Result{}, err
	}
	if errors.IsNotFound(err) {
		processConfigmap = getConfigMapForApp(&appDeployment, "bpmn", map[string]string{"process.bpmn": appDeployment.Spec.Process})
		err = r.Create(ctx, &processConfigmap)
		if err != nil {
			logger.Error(err, "failed to create bpmn configmap")
			return ctrl.Result{}, err
		}
	}
	if processConfigmap.Data["process.bpmn"] != appDeployment.Spec.Process {
		logger.Info("Updating bpmn configmap")
		processConfigmap.Data["process.bpmn"] = appDeployment.Spec.Process
		err = r.Update(ctx, &processConfigmap)
		if err != nil {
			logger.Error(err, "failed to update bpmn configmap")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AppDeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&runtimev1alpha1.AppDeployment{}).
		Owns(&corev1.ConfigMap{}).
		Complete(r)
}

func getConfigMapForApp(appDeployment *runtimev1alpha1.AppDeployment, prefix string, data map[string]string) corev1.ConfigMap {
	return corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s", prefix, appDeployment.Name),
			Namespace: appDeployment.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion:         appDeployment.APIVersion,
					Kind:               appDeployment.Kind,
					Name:               appDeployment.Name,
					UID:                appDeployment.UID,
					BlockOwnerDeletion: BoolPtr(true),
					Controller:         BoolPtr(true),
				},
			},
		},
		Immutable: BoolPtr(false),
		Data:      data,
	}
}

func BoolPtr(b bool) *bool {
	return &b
}
