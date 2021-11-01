// Copyright 2006-2021 VMware, Inc.
// SPDX-License-Identifier: MIT
/*

 */

package commands

import (
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	tenancyv1alpha1 "github.com/vmware-tanzu-labs/namespace-operator/apis/tenancy/v1alpha1"
	tanzunamespace "github.com/vmware-tanzu-labs/namespace-operator/apis/tenancy/v1alpha1/tanzunamespace"
)

func NewGenerateCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "generate <custom-input-manifest-file> <k8s-output-manifest-file>",
		Short: "Generate kubernetes manifests from custom resource manifest",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("requires a custom-input-manifest-file argument")
			}
			if len(args) < 2 {
				return errors.New("requires a k8s-output-manifest-file argument")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			in, err := openInput(args[0])
			if err != nil {
				log.Fatalf("opening input: %v", err)
			}
			defer in.Close()

			out, err := openOutput(args[1])
			if err != nil {
				log.Fatalf("opening output: %v", err)
			}
			defer func() {
				if err := out.Close(); err != nil {
					log.Fatal(err)
				}
			}()

			if err := runGenerate(in, out); err != nil {
				log.Fatal(err)
			}
		},
	}

	return &cmd
}

func runGenerate(in io.Reader, out io.Writer) error {

	var cr tenancyv1alpha1.TanzuNamespace
	if err := Decode(in, &cr); err != nil {
		return fmt.Errorf("decoding: %w", err)
	}

	children, err := generateChildren(&cr)
	if err != nil {
		return err
	}

	for _, child := range children {
		if err := Encode(out, child); err != nil {
			return fmt.Errorf("encoding: %w", err)
		}
	}

	return nil
}

func generateChildren(cr *tenancyv1alpha1.TanzuNamespace) ([]metav1.Object, error) {

	var children []metav1.Object
	for _, f := range tanzunamespace.CreateFuncs {

		rsrc, err := f(cr)
		if err != nil {
			return nil, err
		}

		children = append(children, rsrc)
	}

	return children, nil
}
