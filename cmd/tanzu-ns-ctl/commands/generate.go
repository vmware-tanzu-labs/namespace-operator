// Copyright 2006-2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"sigs.k8s.io/yaml"

	tenancyv1alpha2 "github.com/vmware-tanzu-labs/namespace-operator/apis/tenancy/v1alpha2"
	"github.com/vmware-tanzu-labs/namespace-operator/apis/tenancy/v1alpha2/tanzunamespace"
)

type generateCommand struct {
	*cobra.Command
	workloadManifest string
}

// newGenerateCommand creates a new instance of the generate subcommand.
func (c *TanzuNsCtlCommand) newGenerateCommand() {
	g := &generateCommand{}
	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate child resource manifests from a workload's custom resource",
		Long:  "Generate child resource manifests from a workload's custom resource",
		RunE:  g.generate,
	}

	generateCmd.Flags().StringVarP(
		&g.workloadManifest,
		"workload-manifest",
		"w",
		"",
		"Filepath to the workload manifest to generate child resources for.",
	)
	generateCmd.MarkFlagRequired("workload-manifest")

	c.AddCommand(generateCmd)
}

// generate creates child resource manifests from a workload's custom resource.
func (g *generateCommand) generate(cmd *cobra.Command, args []string) error {

	filename, _ := filepath.Abs(g.workloadManifest)

	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to open file %s, %w", filename, err)
	}

	var workload tenancyv1alpha2.TanzuNamespace

	err = yaml.Unmarshal(yamlFile, &workload)
	if err != nil {
		return fmt.Errorf("failed to unmarshal yaml %s into workload, %w", filename, err)
	}

	err = validateWorkload(&workload)
	if err != nil {
		return fmt.Errorf("error validating yaml %s, %w", filename, err)
	}

	resourceObjects := make([]metav1.Object, len(tanzunamespace.CreateFuncs))

	for i, f := range tanzunamespace.CreateFuncs {
		resource, err := f(&workload)
		if err != nil {
			return err
		}

		resourceObjects[i] = resource
	}

	e := json.NewYAMLSerializer(json.DefaultMetaFactory, nil, nil)

	outputStream := os.Stdout

	for _, o := range resourceObjects {
		if _, err := outputStream.WriteString("---\n"); err != nil {
			return fmt.Errorf("failed to write output, %w", err)
		}

		if err := e.Encode(o.(runtime.Object), os.Stdout); err != nil {
			return fmt.Errorf("failed to write output, %w", err)
		}
	}

	return nil
}
