// Copyright 2006-2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package resources

import (
	"fmt"

	"github.com/banzaicloud/k8s-objectmatcher/patch"
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"github.com/imdario/mergo"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/vmware-tanzu-labs/namespace-operator/apis/common"
)

const (
	FieldManager = "reconciler"
)

// Create creates a resource.
func (resource *Resource) Create() error {
	resource.Reconciler.GetLogger().V(0).Info(fmt.Sprintf("creating resource; kind: [%s], name: [%s], namespace: [%s]",
		resource.Kind, resource.Name, resource.Namespace))

	if err := resource.Reconciler.Create(
		resource.Reconciler.GetContext(),
		resource.Object,
		&client.CreateOptions{FieldManager: FieldManager},
	); err != nil {
		return fmt.Errorf("unable to create resource; %v", err)
	}

	return nil
}

// Update updates a resource.
func (resource *Resource) Update(oldResource *Resource) error {
	needsUpdate, err := NeedsUpdate(*resource, *oldResource)
	if err != nil {
		return err
	}

	if needsUpdate {
		resource.Reconciler.GetLogger().V(0).Info(fmt.Sprintf("updating resource; kind: [%s], name: [%s], namespace: [%s]",
			resource.Kind, resource.Name, resource.Namespace))

		if err := resource.Reconciler.Patch(
			resource.Reconciler.GetContext(),
			resource.Object,
			client.Merge,
			&client.PatchOptions{FieldManager: FieldManager},
		); err != nil {
			return fmt.Errorf("unable to update resource; %v", err)
		}
	}

	return nil
}

// NewResourceFromClient returns a new resource given a client object.  It optionally will take in
// a reconciler and set it.
func NewResourceFromClient(resource client.Object, reconciler ...common.ComponentReconciler) *Resource {
	newResource := &Resource{
		Object: resource,
	}

	// set the inherited fields
	newResource.Group = resource.GetObjectKind().GroupVersionKind().Group
	newResource.Version = resource.GetObjectKind().GroupVersionKind().Version
	newResource.Kind = resource.GetObjectKind().GroupVersionKind().Kind
	newResource.Name = resource.GetName()
	newResource.Namespace = resource.GetNamespace()

	if len(reconciler) > 0 {
		newResource.Reconciler = reconciler[0]
	}

	return newResource
}

// ToUnstructured returns an unstructured representation of a Resource.
func (resource *Resource) ToUnstructured() (*unstructured.Unstructured, error) {
	innerObject, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&resource.Object)
	if err != nil {
		return nil, err
	}

	return &unstructured.Unstructured{Object: innerObject}, nil
}

// ToCommonResource converts a resources.Resource into a common API resource.
func (resource *Resource) ToCommonResource() *common.Resource {
	commonResource := &common.Resource{}

	// set the inherited fields
	commonResource.Group = resource.Group
	commonResource.Version = resource.Version
	commonResource.Kind = resource.Kind
	commonResource.Name = resource.Name
	commonResource.Namespace = resource.Namespace

	return commonResource
}

// IsReady returns whether a specific known resource is ready.  Always returns true for unknown resources
// so that dependency checks will not fail and reconciliation of resources can happen with errors rather
// than stopping entirely.
func (resource *Resource) IsReady() (bool, error) {
	switch resource.Kind {
	case NamespaceKind:
		return NamespaceIsReady(resource)
	case CustomResourceDefinitionKind:
		return CustomResourceDefinitionIsReady(resource)
	case SecretKind:
		return SecretIsReady(resource)
	case ConfigMapKind:
		return ConfigMapIsReady(resource)
	case DeploymentKind:
		return DeploymentIsReady(resource)
	case DaemonSetKind:
		return DaemonSetIsReady(resource)
	case StatefulSetKind:
		return StatefulSetIsReady(resource)
	case JobKind:
		return JobIsReady(resource)
	case ServiceKind:
		return ServiceIsReady(resource)
	}

	return true, nil
}

// AreReady returns whether resources are ready.  All resources must be ready in order
// to satisfy the requirement that resources are ready.
func AreReady(resources ...common.ComponentResource) (bool, error) {
	for _, resource := range resources {
		ready, err := resource.IsReady()
		if !ready || err != nil {
			return false, err
		}
	}

	return true, nil
}

// AreEqual determines if two resources are equal.
func AreEqual(desired, actual Resource) (bool, error) {
	mergedResource, err := actual.ToUnstructured()
	if err != nil {
		return false, err
	}

	actualResource, err := actual.ToUnstructured()
	if err != nil {
		return false, err
	}

	desiredResource, err := desired.ToUnstructured()
	if err != nil {
		return false, err
	}

	// ensure that resource versions and observed generation do not interfere
	// with calculating equality
	desiredResource.SetResourceVersion(actualResource.GetResourceVersion())
	desiredResource.SetGeneration(actualResource.GetGeneration())

	// ensure that a current cluster-scoped resource is not evaluated against
	// a manifest which may include a namespace
	if actualResource.GetNamespace() == "" {
		desiredResource.SetNamespace(actualResource.GetNamespace())
	}

	// merge the overrides from the desired resource into the actual resource
	mergo.Merge(
		&mergedResource.Object,
		desiredResource.Object,
		mergo.WithOverride,
		mergo.WithSliceDeepCopy,
	)

	// calculate the actual differences
	diffOptions := []patch.CalculateOption{
		reconciler.IgnoreManagedFields(),
		patch.IgnoreStatusFields(),
		patch.IgnoreVolumeClaimTemplateTypeMetaAndStatus(),
		patch.IgnorePDBSelector(),
	}

	diffResults, err := patch.DefaultPatchMaker.Calculate(
		actualResource,
		mergedResource,
		diffOptions...,
	)
	if err != nil {
		return false, err
	}

	return diffResults.IsEmpty(), nil
}

// NeedsUpdate determines if a resource needs to be updated.
func NeedsUpdate(desired, actual Resource) (bool, error) {
	// check for equality first as this will let us avoid spamming user logs
	// when resources that need to be skipped explicitly (e.g. CRDs) are seen
	// as equal anyway
	equal, err := AreEqual(desired, actual)
	if equal || err != nil {
		return !equal, err
	}

	// always skip custom resource updates as they are sensitive to modification
	// e.g. resources provisioned by the resource definition would not
	// understand the update to a spec
	if desired.Kind == "CustomResourceDefinition" {
		message := fmt.Sprintf("skipping update of CustomResourceDefinition "+
			"[%s]", desired.Name)
		messageVerbose := fmt.Sprintf("if updates to CustomResourceDefinition "+
			"[%s] are desired, consider re-deploying the parent "+
			"resource or generating a new api version with the desired "+
			"changes", desired.Name)
		desired.Reconciler.GetLogger().V(4).Info(message)
		desired.Reconciler.GetLogger().V(7).Info(messageVerbose)

		return false, nil
	}

	return true, nil
}

// EqualNamespaceName will compare the namespace and name of two resource objects for equality.
func (resource *Resource) EqualNamespaceName(compared common.ComponentResource) bool {
	comparedResource := compared.(*Resource)
	return (resource.Name == comparedResource.Name) && (resource.Namespace == comparedResource.Namespace)
}

// EqualGVK will compare the GVK of two resource objects for equality.
func (resource *Resource) EqualGVK(compared common.ComponentResource) bool {
	comparedResource := compared.(*Resource)
	return resource.Group == comparedResource.Group &&
		resource.Version == comparedResource.Version &&
		resource.Kind == comparedResource.Kind
}

// GetObject returns the Object field of a Resource.
func (resource *Resource) GetObject() client.Object {
	return resource.Object
}

// GetReconciler returns the Reconciler field of a Resource.
func (resource *Resource) GetReconciler() common.ComponentReconciler {
	return resource.Reconciler
}

// GetGroup returns the Group field of a Resource.
func (resource *Resource) GetGroup() string {
	return resource.Group
}

// GetVersion returns the Version field of a Resource.
func (resource *Resource) GetVersion() string {
	return resource.Version
}

// GetKind returns the Kind field of a Resource.
func (resource *Resource) GetKind() string {
	return resource.Kind
}

// GetName returns the Name field of a Resource.
func (resource *Resource) GetName() string {
	return resource.Name
}

// GetNamespace returns the Name field of a Resource.
func (resource *Resource) GetNamespace() string {
	return resource.Namespace
}

// getObject returns an object based on an input object, and a destination object.
// TODO: move to controller utils as this is not specific to resources.
func getObject(source common.ComponentResource, destination client.Object, allowMissing bool) error {
	namespacedName := types.NamespacedName{
		Name:      source.GetName(),
		Namespace: source.GetNamespace(),
	}
	if err := source.GetReconciler().Get(source.GetReconciler().GetContext(), namespacedName, destination); err != nil {
		if allowMissing {
			if errors.IsNotFound(err) {
				return nil
			}
		} else {
			return err
		}
	}

	return nil
}
