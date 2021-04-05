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

package resources

import (
	"strconv"
	"text/template"

	core_k8s_api "k8s.io/api/core/v1"
	rbac_k8s_api "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"

	tenancyv1alpha1 "github.com/vmware-tanzu-labs/namespace-operator/api/v1alpha1"
)

const resourceServiceAccount = `
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: {{ .User }}
  namespace: {{ .Namespace }}
`

const resourceRole = `
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: {{ .Role }}
  namespace: {{ .Namespace }}
rules:
{{ .Permissions }}
`

const resourceRoleBinding = `
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: {{ .RoleBinding }}
  namespace: {{ .Namespace }}
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: {{ .Role }}
subjects:
  - kind: ServiceAccount
    name: {{ .User }}
    namespace: {{ .Namespace }}
`

type rbacObject struct {
	Template   string
	ObjectType string
}

// CreateRBAC creates the appropriate RBAC policies
func CreateRBAC(parent *tenancyv1alpha1.TanzuNamespace) ([]metav1.Object, error) {
	var resourceObjs []metav1.Object

	fmap := template.FuncMap{}

	var rbacObjects = []rbacObject{
		{
			Template:   resourceServiceAccount,
			ObjectType: "ServiceAccount",
		},
		{
			Template:   resourceRole,
			ObjectType: "Role",
		},
		{
			Template:   resourceRoleBinding,
			ObjectType: "RoleBinding",
		},
	}

	for i, rbacItem := range rbacObjects {
		for _, rbac := range setRBAC(parent) {
			templateName := "tanzu-rbac-" + rbac.User + "-" + strconv.Itoa(i)

			childContent, err := runTemplate(templateName, rbacItem.Template, rbac, fmap)
			if err != nil {
				return nil, err
			}

			decode := scheme.Codecs.UniversalDeserializer().Decode
			obj, _, err := decode([]byte(childContent), nil, nil)
			if err != nil {
				return nil, err
			}

			var resourceObj metav1.Object
			if rbacItem.ObjectType == "ServiceAccount" {
				resourceObj = obj.(*core_k8s_api.ServiceAccount)
			} else if rbacItem.ObjectType == "Role" {
				resourceObj = obj.(*rbac_k8s_api.Role)
			} else if rbacItem.ObjectType == "RoleBinding" {
				resourceObj = obj.(*rbac_k8s_api.RoleBinding)
			}
			resourceObjs = append(resourceObjs, resourceObj)
		}
	}

	return resourceObjs, nil
}
