// Copyright 2006-2021 VMware, Inc.
// SPDX-License-Identifier: MIT
/*

 */

package resources

import (
	"bytes"
	"fmt"
	"text/template"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	tenancyv1alpha1 "github.com/vmware-tanzu-labs/namespace-operator/api/v1alpha1"
)

// CreateArrayFuncs is a variable containing functions which return an array of resources
var CreateArrayFuncs = []func(*tenancyv1alpha1.TanzuNamespace) ([]metav1.Object, error){
	CreateNetworkPolicies,
}

// CreateFuncs is a variable containing functions which return a single resource
var CreateFuncs = []func(*tenancyv1alpha1.TanzuNamespace) (metav1.Object, error){
	CreateNamespace,
	CreateLimitRange,
	CreateResourceQuota,
}

// genericResource is a generic type to allow multiple different objects to be templated
// in the same manner
type genericResource interface{}

// runTemplate renders a template for a child object to the custom resource
func runTemplate(templateName, templateValue string, data genericResource,
	funcMap template.FuncMap) (string, error) {

	t, err := template.New(templateName).Funcs(funcMap).Parse(templateValue)
	if err != nil {
		return "", fmt.Errorf("error parsing template %s: %v", templateName, err)
	}

	var b bytes.Buffer
	if err := t.Execute(&b, &data); err != nil {
		return "", fmt.Errorf("error rendering template %s: %v", templateName, err)
	}

	return b.String(), nil
}

const (
	limitRangeDefaultCPULimitDefault      = "125m"
	limitRangeDefaultMemoryLimitDefault   = "64Mi"
	limitRangeDefaultCPURequestDefault    = "125m"
	limitRangeDefaultMemoryRequestDefault = "64Mi"
	limitRangeMaxCPULimitDefault          = "1000m"
	limitRangeMaxMemoryLimitDefault       = "2Gi"
	resourceQuotaCPURequestsDefault       = "2000m"
	resourceQuotaMemoryRequestsDefault    = "4Gi"
	resourceQuotaCPULimitsDefault         = "2000m"
	resourceQuotaMemoryLimitsDefault      = "4Gi"
)

const namespaceAdminPerms = `
  - apiGroups:
      - "*"
    resources:
      - "*"
    verbs:
      - "*"
`

const developerPerms = `
  - apiGroups:
      - ""
    resources:
      - configmaps
      - secrets
      - services
    verbs:
      - "*"
  - apiGroups:
      - apps
    resources:
      - pods
      - deployments
      - statefulsets
      - daemonsets
    verbs:
      - "*"
  - apiGroups:
      - extensions
    resources:
      - ingress
      - ingresses
    verbs:
      - "*"
  - apiGroups:
      - networking.k8s.io
    resources:
      - networkpolicies
    verbs:
      - "*"
  - apiGroups:
      - batch
    resources:
      - jobs
      - cronjobs
    verbs:
      - "*"
`

const readOnlyPerms = `
  - apiGroups:
      - "*"
    resources:
      - "*"
    verbs:
      - get
      - list
      - watch
`

var rbacDefaults = map[string]map[string]string{
	"namespace-admin": {
		"user":        "tanzu-namespace-admin",
		"role":        "tanzu-namespace-admin-role",
		"rolebinding": "tanzu-namespace-admin-rolebinding",
		"permissions": namespaceAdminPerms,
	},
	"developer": {
		"user":        "tanzu-developer",
		"role":        "tanzu-developer-role",
		"rolebinding": "tanzu-developer-rolebinding",
		"permissions": developerPerms,
	},
	"read-only": {
		"user":        "tanzu-read-only",
		"role":        "tanzu-read-only-role",
		"rolebinding": "tanzu-read-only-rolebinding",
		"permissions": readOnlyPerms,
	},
}

func defaultLimitRangeDefaultCPULimit(spec *tenancyv1alpha1.TanzuNamespaceSpec) string {
	if spec.LimitRange.DefaultCPULimit != "" {
		return spec.LimitRange.DefaultCPULimit
	} else if spec.TanzuLimitRangeDefaultCpuLimit != "" {
		return spec.TanzuLimitRangeDefaultCpuLimit
	}

	return limitRangeDefaultCPULimitDefault
}

func defaultLimitRangeDefaultMemoryLimit(spec *tenancyv1alpha1.TanzuNamespaceSpec) string {
	if spec.LimitRange.DefaultMemoryLimit != "" {
		return spec.LimitRange.DefaultMemoryLimit
	} else if spec.TanzuLimitRangeDefaultMemoryLimit != "" {
		return spec.TanzuLimitRangeDefaultMemoryLimit
	}

	return limitRangeDefaultMemoryLimitDefault
}

func defaultLimitRangeDefaultCPURequest(spec *tenancyv1alpha1.TanzuNamespaceSpec) string {
	if spec.LimitRange.DefaultCPURequest != "" {
		return spec.LimitRange.DefaultCPURequest
	} else if spec.TanzuLimitRangeDefaultCpuRequest != "" {
		return spec.TanzuLimitRangeDefaultCpuRequest
	}

	return limitRangeDefaultCPURequestDefault
}

func defaultLimitRangeDefaultMemoryRequest(spec *tenancyv1alpha1.TanzuNamespaceSpec) string {
	if spec.LimitRange.DefaultMemoryRequest != "" {
		return spec.LimitRange.DefaultMemoryRequest
	} else if spec.TanzuLimitRangeDefaultMemoryRequest != "" {
		return spec.TanzuLimitRangeDefaultMemoryRequest
	}

	return limitRangeDefaultMemoryRequestDefault
}

func defaultLimitRangeMaxCPULimit(spec *tenancyv1alpha1.TanzuNamespaceSpec) string {
	if spec.LimitRange.MaxCPULimit != "" {
		return spec.LimitRange.MaxCPULimit
	} else if spec.TanzuLimitRangeMaxCpuLimit != "" {
		return spec.TanzuLimitRangeMaxCpuLimit
	}

	return limitRangeMaxCPULimitDefault
}

func defaultLimitRangeMaxMemoryLimit(spec *tenancyv1alpha1.TanzuNamespaceSpec) string {
	if spec.LimitRange.MaxMemoryLimit != "" {
		return spec.LimitRange.MaxMemoryLimit
	} else if spec.TanzuLimitRangeMaxMemoryLimit != "" {
		return spec.TanzuLimitRangeMaxMemoryLimit
	}

	return limitRangeMaxMemoryLimitDefault
}

func defaultResourceQuotaCPURequests(spec *tenancyv1alpha1.TanzuNamespaceSpec) string {
	if spec.ResourceQuota.RequestsCPU != "" {
		return spec.ResourceQuota.RequestsCPU
	} else if spec.TanzuResourceQuotaCpuRequests != "" {
		return spec.TanzuResourceQuotaCpuRequests
	}

	return resourceQuotaCPURequestsDefault
}

func defaultResourceQuotaMemoryRequests(spec *tenancyv1alpha1.TanzuNamespaceSpec) string {
	if spec.ResourceQuota.RequestsMemory != "" {
		return spec.ResourceQuota.RequestsMemory
	} else if spec.TanzuResourceQuotaMemoryRequests != "" {
		return spec.TanzuResourceQuotaMemoryRequests
	}

	return resourceQuotaMemoryRequestsDefault
}

func defaultResourceQuotaCPULimits(spec *tenancyv1alpha1.TanzuNamespaceSpec) string {
	if spec.ResourceQuota.LimitsCPU != "" {
		return spec.ResourceQuota.LimitsCPU
	} else if spec.TanzuResourceQuotaCpuLimits != "" {
		return spec.TanzuResourceQuotaCpuLimits
	}

	return resourceQuotaCPULimitsDefault
}

func defaultResourceQuotaMemoryLimits(spec *tenancyv1alpha1.TanzuNamespaceSpec) string {
	if spec.ResourceQuota.LimitsMemory != "" {
		return spec.ResourceQuota.LimitsMemory
	} else if spec.TanzuResourceQuotaMemoryLimits != "" {
		return spec.TanzuResourceQuotaMemoryLimits
	}

	return resourceQuotaMemoryLimitsDefault
}

func defaultLabels(labels map[string]string) map[string]string {
	if len(labels) == 0 {
		return labels
	}

	return labels
}

func defaultIngressTCPPorts(ingressTCPPorts []int) (ports []int) {
	if len(ingressTCPPorts) == 0 {
		return ports
	}

	return ingressTCPPorts
}

func defaultIngressUDPPorts(ingressUDPPorts []int) (ports []int) {
	if len(ingressUDPPorts) == 0 {
		return ports
	}

	return ingressUDPPorts
}

func defaultEgressTCPPorts(egressTCPPorts []int) (ports []int) {
	if len(egressTCPPorts) == 0 {
		return ports
	}

	return egressTCPPorts
}

func defaultEgressUDPPorts(egressUDPPorts []int) (ports []int) {
	if len(egressUDPPorts) == 0 {
		return ports
	}

	return egressUDPPorts
}

func defaultIngressPorts(ingressTCPPorts []int, ingressUDPPorts []int) (ports []tenancyv1alpha1.NetworkPort) {
	for _, port := range defaultIngressTCPPorts(ingressTCPPorts) {
		ports = append(ports, tenancyv1alpha1.NetworkPort{Protocol: "TCP", Port: port})
	}

	for _, port := range defaultIngressUDPPorts(ingressUDPPorts) {
		ports = append(ports, tenancyv1alpha1.NetworkPort{Protocol: "UDP", Port: port})
	}

	return ports
}

func defaultEgressPorts(egressTCPPorts []int, egressUDPPorts []int) (ports []tenancyv1alpha1.NetworkPort) {
	for _, port := range defaultEgressTCPPorts(egressTCPPorts) {
		ports = append(ports, tenancyv1alpha1.NetworkPort{Protocol: "TCP", Port: port})
	}

	for _, port := range defaultEgressUDPPorts(egressUDPPorts) {
		ports = append(ports, tenancyv1alpha1.NetworkPort{Protocol: "UDP", Port: port})
	}

	return ports
}

func defaultNamespace(parentName string, spec *tenancyv1alpha1.TanzuNamespaceSpec) string {
	if spec.Name != "" {
		return spec.Name
	} else if spec.TanzuNamespaceName != "" {
		return spec.TanzuNamespaceName
	}

	return parentName
}

func setRBAC(parent *tenancyv1alpha1.TanzuNamespace) (rbac []tenancyv1alpha1.RBAC) {
OuterLoop:
	// loop through the defaults and ensure they have values set
	for rbacType, defaults := range rbacDefaults {
		// loop through the input from the crd and select items which have been selected
		matchedTypes := make([]tenancyv1alpha1.RBAC, 0)
		for _, rbacObject := range parent.Spec.RBAC {
			if rbacObject.Type == rbacType {
				matchedTypes = append(matchedTypes, rbacObject)
			}
		}

		// set namespace
		namespace := defaultNamespace(parent.Name, &parent.Spec)

		if len(matchedTypes) == 0 {
			// use the object defaults
			rbac = append(rbac, tenancyv1alpha1.RBAC{
				Type:        rbacType,
				User:        defaults["user"],
				Permissions: defaults["permissions"],
				Role:        defaults["role"],
				RoleBinding: defaults["rolebinding"],
				Create:      true,
				Namespace:   namespace,
			})
		} else {
			// loop through selected elements and append them to the return values
			// NOTE: we will only use the first value if multiple are specified
			for _, match := range matchedTypes {
				if match.Create != false {
					match.Namespace = namespace
					match.Permissions = defaults["permissions"]
					rbac = append(rbac, match)
					continue OuterLoop
				}
			}
		}
	}

	return rbac
}
