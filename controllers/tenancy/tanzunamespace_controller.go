// Copyright 2006-2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package tenancy

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	"github.com/vmware-tanzu-labs/namespace-operator/apis/common"
	tenancyv1alpha2 "github.com/vmware-tanzu-labs/namespace-operator/apis/tenancy/v1alpha2"
	"github.com/vmware-tanzu-labs/namespace-operator/apis/tenancy/v1alpha2/tanzunamespace"
	"github.com/vmware-tanzu-labs/namespace-operator/internal/controllers/phases"
	"github.com/vmware-tanzu-labs/namespace-operator/internal/controllers/utils"
	"github.com/vmware-tanzu-labs/namespace-operator/internal/dependencies"
	"github.com/vmware-tanzu-labs/namespace-operator/internal/mutate"
	"github.com/vmware-tanzu-labs/namespace-operator/internal/resources"
	"github.com/vmware-tanzu-labs/namespace-operator/internal/wait"
)

// TanzuNamespaceReconciler reconciles a TanzuNamespace object.
type TanzuNamespaceReconciler struct {
	client.Client
	Name       string
	Log        logr.Logger
	Scheme     *runtime.Scheme
	Context    context.Context
	Controller controller.Controller
	Watches    []client.Object
	Resources  []common.ComponentResource
	Component  *tenancyv1alpha2.TanzuNamespace
}

// +kubebuilder:rbac:groups=tenancy.platform.cnr.vmware.com,resources=tanzunamespaces,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=tenancy.platform.cnr.vmware.com,resources=tanzunamespaces/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core,resources=namespaces,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=limitranges,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=resourcequotas,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=networking.k8s.io,resources=networkpolicies,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the WebApp object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.2/pkg/reconcile
func (r *TanzuNamespaceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	r.Context = ctx
	log := r.Log.WithValues("tanzunamespace", req.NamespacedName)

	// get and store the component
	r.Component = &tenancyv1alpha2.TanzuNamespace{}
	if err := r.Get(r.Context, req.NamespacedName, r.Component); err != nil {
		log.V(0).Info("unable to fetch TanzuNamespace")

		return ctrl.Result{}, utils.IgnoreNotFound(err)
	}

	// get and store the resources
	if err := r.SetResources(); err != nil {
		return ctrl.Result{}, err
	}

	// execute the phases
	for _, phase := range utils.Phases(r.Component) {
		r.GetLogger().V(7).Info(fmt.Sprintf("enter phase: %T", phase))
		proceed, err := phase.Execute(r)
		result, err := phases.HandlePhaseExit(r, phase, proceed, err)

		// return only if we have an error or are told not to proceed
		if err != nil || !proceed {
			log.V(2).Info(fmt.Sprintf("not ready; requeuing phase: %T", phase))

			return result, err
		}

		r.GetLogger().V(5).Info(fmt.Sprintf("completed phase: %T", phase))
	}

	return phases.DefaultReconcileResult(), nil
}

// Construct resources runs the methods to properly construct the resources.
func (r *TanzuNamespaceReconciler) ConstructResources() ([]metav1.Object, error) {

	resourceObjects := make([]metav1.Object, len(tanzunamespace.CreateFuncs))

	// create resources in memory
	for i, f := range tanzunamespace.CreateFuncs {
		resource, err := f(r.Component)
		if err != nil {
			return nil, err
		}

		resourceObjects[i] = resource
	}

	return resourceObjects, nil
}

// GetResources will return the resources associated with the reconciler.
func (r *TanzuNamespaceReconciler) GetResources() []common.ComponentResource {
	return r.Resources
}

// SetResources will create and return the resources in memory.
func (r *TanzuNamespaceReconciler) SetResources() error {
	// create resources in memory
	baseResources, err := r.ConstructResources()
	if err != nil {
		return err
	}

	// loop through the in memory resources and store them on the reconciler
	for _, base := range baseResources {
		// run through the mutation functions to mutate the resources
		mutatedResources, skip, err := r.Mutate(&base)
		if err != nil {
			return err
		}
		if skip {
			continue
		}

		for _, mutated := range mutatedResources {
			resourceObject := resources.NewResourceFromClient(mutated.(client.Object))
			resourceObject.Reconciler = r

			r.SetResource(resourceObject)
		}
	}

	return nil
}

// SetResource will set a resource on the objects if the relevant object does not already exist.
func (r *TanzuNamespaceReconciler) SetResource(new common.ComponentResource) {

	// set and return immediately if nothing exists
	if len(r.Resources) == 0 {
		r.Resources = append(r.Resources, new)

		return
	}

	// loop through the resources and set or update when found
	for i, existing := range r.Resources {
		if new.EqualGVK(existing) && new.EqualNamespaceName(existing) {
			r.Resources[i] = new

			return
		}
	}

	// if we haven't returned yet, we have not found the resource and must add it
	r.Resources = append(r.Resources, new)
}

// CreateOrUpdate creates a resource if it does not already exist or updates a resource
// if it does already exist.
func (r *TanzuNamespaceReconciler) CreateOrUpdate(
	resource metav1.Object,
) error {
	// set ownership on the underlying resource being created or updated
	if err := ctrl.SetControllerReference(r.Component, resource, r.Scheme); err != nil {
		r.GetLogger().V(0).Info("unable to set owner reference on resource")

		return err
	}

	// create a stub object to store the current resource in the cluster so that we do not affect
	// the desired state of the resource object in memory
	newResource := resources.NewResourceFromClient(resource.(client.Object), r)
	resourceStub := &unstructured.Unstructured{}
	resourceStub.SetGroupVersionKind(newResource.Object.GetObjectKind().GroupVersionKind())
	oldResource := resources.NewResourceFromClient(resourceStub, r)

	if err := r.Get(
		r.Context,
		client.ObjectKeyFromObject(newResource.Object),
		oldResource.Object,
	); err != nil {
		// create the resource if we cannot find one
		if errors.IsNotFound(err) {
			if err := newResource.Create(); err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		// update the resource
		if err := newResource.Update(oldResource); err != nil {
			return err
		}
	}

	return utils.Watch(r, newResource.Object)
}

// GetLogger returns the logger from the reconciler.
func (r *TanzuNamespaceReconciler) GetLogger() logr.Logger {
	return r.Log
}

// GetClient returns the client from the reconciler.
func (r *TanzuNamespaceReconciler) GetClient() client.Client {
	return r.Client
}

// GetScheme returns the scheme from the reconciler.
func (r *TanzuNamespaceReconciler) GetScheme() *runtime.Scheme {
	return r.Scheme
}

// GetContext returns the context from the reconciler.
func (r *TanzuNamespaceReconciler) GetContext() context.Context {
	return r.Context
}

// GetName returns the name of the reconciler.
func (r *TanzuNamespaceReconciler) GetName() string {
	return r.Name
}

// GetComponent returns the component the reconciler is operating against.
func (r *TanzuNamespaceReconciler) GetComponent() common.Component {
	return r.Component
}

// GetController returns the controller object associated with the reconciler.
func (r *TanzuNamespaceReconciler) GetController() controller.Controller {
	return r.Controller
}

// GetWatches returns the objects which are current being watched by the reconciler.
func (r *TanzuNamespaceReconciler) GetWatches() []client.Object {
	return r.Watches
}

// SetWatch appends a watch to the list of currently watched objects.
func (r *TanzuNamespaceReconciler) SetWatch(watch client.Object) {
	r.Watches = append(r.Watches, watch)
}

// UpdateStatus updates the status for a component.
func (r *TanzuNamespaceReconciler) UpdateStatus() error {
	return r.Status().Update(r.Context, r.Component)
}

// CheckReady will return whether a component is ready.
func (r *TanzuNamespaceReconciler) CheckReady() (bool, error) {
	return dependencies.TanzuNamespaceCheckReady(r)
}

// Mutate will run the mutate phase of a resource.
func (r *TanzuNamespaceReconciler) Mutate(
	object *metav1.Object,
) ([]metav1.Object, bool, error) {
	return mutate.TanzuNamespaceMutate(r, object)
}

// Wait will run the wait phase of a resource.
func (r *TanzuNamespaceReconciler) Wait(
	object *metav1.Object,
) (bool, error) {
	return wait.TanzuNamespaceWait(r, object)
}

func (r *TanzuNamespaceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	options := controller.Options{
		RateLimiter: utils.NewDefaultRateLimiter(5*time.Microsecond, 5*time.Minute),
	}

	baseController, err := ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		WithEventFilter(utils.ComponentPredicates()).
		For(&tenancyv1alpha2.TanzuNamespace{}).
		Build(r)
	if err != nil {
		return err
	}

	r.Controller = baseController

	return nil
}
