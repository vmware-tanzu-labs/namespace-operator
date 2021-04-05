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

package commands

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"sigs.k8s.io/yaml"
)

// openInput opens a file at the given path or captures stdin.
func openInput(path string) (io.ReadCloser, error) {
	if path == "-" {
		return os.Stdin, nil
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}

	return f, nil
}

// openOutput writes to the file at the given path or to stdout
func openOutput(path string) (io.WriteCloser, error) {
	if path == "-" {
		return os.Stdout, nil
	}

	f, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}

	return f, nil
}

// Decode decodes yaml into objects.
func Decode(r io.Reader, o interface{}) error {
	btys, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(btys, o)
}

// Encode encodes objects into yaml. Call multiple times on a single writer
// to encode multiple documents.
func Encode(w io.Writer, o interface{}) error {
	btys, err := yaml.Marshal(o)
	if err != nil {
		return err
	}

	if _, err := w.Write([]byte("---\n")); err != nil {
		return err
	}
	if _, err := w.Write(btys); err != nil {
		return err
	}

	return nil
}
