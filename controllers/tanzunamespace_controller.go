// Copyright 2006-2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0
/*
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

	"strings"

	"github.com/go-logr/logr"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	tenancyv1alpha1 "github.com/vmware-tanzu-labs/namespace-operator/api/v1alpha1"
	"github.com/vmware-tanzu-labs/namespace-operator/resources"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TanzuNamespaceReconciler reconciles a TanzuNamespace object
type TanzuNamespaceReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

func ignoreNotFound(err error) error {
	if apierrs.IsNotFound(err) {
		return nil
	}
	return err
}

// +kubebuilder:rbac:groups=tenancy.platform.cnr.vmware.com,resources=tanzunamespaces,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=tenancy.platform.cnr.vmware.com,resources=tanzunamespaces/status,verbs=get;update;patch

func (r *TanzuNamespaceReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("tanzunamespace", req.NamespacedName)

	// get the custom resource, which has been persisted to etcd by now, from the cluster
	var customResource tenancyv1alpha1.TanzuNamespace
	if err := r.Get(ctx, req.NamespacedName, &customResource); err != nil {
		log.V(0).Info("unable to fetch TanzuNamespace")
		return ctrl.Result{}, ignoreNotFound(err)
	}

	// logic for a resource that has not yet been created
	if customResource.Status.Created == false {

		for _, f := range resources.CreateFuncs {

			rsrc, err := f(&customResource)

			if err != nil {
				return ctrl.Result{}, err
			}

			_, err = ctrl.CreateOrUpdate(ctx, r.Client, rsrc.(runtime.Object), func() error {
				object := rsrc.(metav1.Object)
				if err := ctrl.SetControllerReference(&customResource, object, r.Scheme); err != nil {
					log.Error(err, "unable to set owner reference on resource")
					return err
				}
				return nil
			})
			if err != nil {
				log.Error(err, "unable to create or update resource")
				return ctrl.Result{}, err
			}
		}

		for _, f := range resources.CreateArrayFuncs {

			rsrcs, err := f(&customResource)
			if err != nil {
				return ctrl.Result{}, err
			}

			for _, rsrc := range rsrcs {

				_, err = ctrl.CreateOrUpdate(ctx, r.Client, rsrc.(runtime.Object), func() error {
					object := rsrc.(metav1.Object)
					if err := ctrl.SetControllerReference(&customResource, object, r.Scheme); err != nil {
						log.Error(err, "unable to set owner reference on resource")
						return err
					}
					return nil
				})
				if err != nil {
					log.Error(err, "unable to create or update resource")
					return ctrl.Result{}, err
				}
			}
		}

		customResource.Status.Created = true
		customResource.Status.Conditions = []tenancyv1alpha1.Condition{
			{
				Type:    "Ready",
				Status:  "True",
				Reason:  "InitialCreate",
				Message: "TanzuNamespace Created",
			},
		}

		err := r.Status().Update(ctx, &customResource)
		if err != nil {
			if strings.Contains(err.Error(), "the object has been modified") {
				log.V(0).Info("unable to inject status field; retrying reconciliation")
				return ctrl.Result{}, nil
			}
			log.Error(err, "unable to update custom resource status")
			return ctrl.Result{}, err
		}

		log.V(0).Info("new resources created")

	} else {
		log.V(0).Info("resources exist - no update logic implemented")
	}

	return ctrl.Result{}, nil
}

func (r *TanzuNamespaceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&tenancyv1alpha1.TanzuNamespace{}).
		Complete(r)
}
