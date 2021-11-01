// Copyright 2006-2021 VMware, Inc.
// SPDX-License-Identifier: MIT
/*

 */

package resources

import (
	"strconv"
	"text/template"

	k8s_api "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"

	tenancyv1alpha1 "github.com/vmware-tanzu-labs/namespace-operator/apis/v1alpha1"
)

const resourceNetworkPolicy = `
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: tanzu-network-policy-template
spec:
  #
  # target
  #
  podSelector:

  # targetPodLabels
{{ $targetPodLabels := defaultLabels .TargetPodLabels }}
{{ if eq (len $targetPodLabels) 0 }}
    matchLabels: {}
{{ else }}
    matchLabels:
{{ range $labelKey, $labelValue := $targetPodLabels }}
      {{ $labelKey }}: {{ $labelValue }}
{{ end }}
{{ end }}

  #
  # allow ingress
  #
{{ $ingressNamespaceLabels := defaultLabels .IngressNamespaceLabels }}
{{ $ingressPodLabels := defaultLabels .IngressPodLabels }}
  ingress:

{{ if and (eq (len $ingressNamespaceLabels) 0) (eq (len $ingressPodLabels) 0) }}
    - from: []
{{ else }}
    - from:
      # ingressNamespaceLabels
      - namespaceSelector:
{{ if eq (len $ingressNamespaceLabels) 0 }}
          matchLabels: {}
{{ else }}
          matchLabels:
{{ range $labelKey, $labelValue := $ingressNamespaceLabels }}
            {{ $labelKey }}: {{ $labelValue }}
{{ end }}
{{ end }}

      # ingressPodLabels
      - podSelector:
{{ if eq (len $ingressPodLabels) 0 }}
          matchLabels: {}
{{ else }}
          matchLabels:
{{ range $labelKey, $labelValue := $ingressPodLabels }}
            {{ $labelKey }}: {{ $labelValue }}
{{ end }}
{{ end }}
{{ end }}

      # ingressPorts
{{ $ingressPorts := ingressPorts .IngressTCPPorts .IngressUDPPorts }}
{{ if eq (len $ingressPorts) 0 }}
      ports: []
{{ else }}
      ports:
{{ range $ingressPort := $ingressPorts }}
        - protocol: {{ $ingressPort.Protocol }}
          port: {{ $ingressPort.Port }}
{{ end }}
{{ end }}

  #
  # allow egress
  #
{{ $egressNamespaceLabels := defaultLabels .EgressNamespaceLabels }}
{{ $egressPodLabels := defaultLabels .EgressPodLabels }}
  egress:

{{ if and (eq (len $egressNamespaceLabels) 0) (eq (len $egressPodLabels) 0) }}
    - to: []
{{ else }}
    - to:

      # egressNamespaceLabels
      - namespaceSelector:
{{ if eq (len $egressNamespaceLabels) 0 }}
          matchLabels: {}
{{ else }}
          matchLabels:
{{ range $labelKey, $labelValue := $egressNamespaceLabels }}
            {{ $labelKey }}: {{ $labelValue }}
{{ end }}
{{ end }}

      # egressPodLabels
      - podSelector:
{{ if eq (len $egressPodLabels) 0 }}
          matchLabels: {}
{{ else }}
          matchLabels:
{{ range $labelKey, $labelValue := $egressPodLabels }}
            {{ $labelKey }}: {{ $labelValue }}
{{ end }}
{{ end }}
{{ end }}

      # egressPorts
{{ $egressPorts := egressPorts .EgressTCPPorts .EgressUDPPorts }}
{{ if eq (len $egressPorts) 0 }}
      ports: []
{{ else }}
      ports:
{{ range $egressPort := $egressPorts }}
        - protocol: {{ $egressPort.Protocol }}
          port: {{ $egressPort.Port }}
{{ end }}
{{ end }}
`

// CreateNetworkPolicies creates the tanzu-network-policy-n resource
func CreateNetworkPolicies(parent *tenancyv1alpha1.TanzuNamespace) ([]metav1.Object, error) {
	var resourceObjs []metav1.Object

	for i, networkPolicy := range parent.Spec.NetworkPolicies {
		templateName := "tanzu-network-policy-" + strconv.Itoa(i)

		fmap := template.FuncMap{
			"defaultLabels": defaultLabels,
			"ingressPorts":  defaultIngressPorts,
			"egressPorts":   defaultEgressPorts,
		}

		childContent, err := runTemplate(templateName, resourceNetworkPolicy, networkPolicy, fmap)
		if err != nil {
			return nil, err
		}

		decode := scheme.Codecs.UniversalDeserializer().Decode
		obj, _, err := decode([]byte(childContent), nil, nil)
		if err != nil {
			return nil, err
		}

		resourceObj := obj.(*k8s_api.NetworkPolicy)
		resourceObj.Namespace = defaultNamespace(parent.Name, &parent.Spec)
		resourceObj.Name = templateName
		resourceObjs = append(resourceObjs, resourceObj)
	}

	return resourceObjs, nil
}
