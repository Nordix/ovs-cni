// Copyright (c) 2021 Nordix Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Go version 1.10 or greater is required. Before that, switching namespaces in
// long running processes in go did not work in a reliable way.
// +build go1.10

package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	// DefaultCNIDir used for caching hostIFName
	DefaultCNIDir = "/var/lib/cni/ovs"
)

// SaveConf takes in container ID, data dir and Pod interface name as string and a json encoded struct Conf
// and save this Conf in data dir
func SaveConf(cid, podIfName string, conf interface{}) error {
	confBytes, err := json.Marshal(conf)
	if err != nil {
		return fmt.Errorf("error serializing delegate conf: %v", err)
	}

	s := []string{cid, podIfName}
	cRef := strings.Join(s, "-")

	// save the rendered conf for cmdDel
	if err = saveScratchConf(cRef, DefaultCNIDir, confBytes); err != nil {
		return err
	}

	return nil
}

func saveScratchConf(containerID, dataDir string, conf []byte) error {
	if err := os.MkdirAll(dataDir, 0700); err != nil {
		return fmt.Errorf("failed to create the sriov data directory(%q): %v", dataDir, err)
	}

	path := filepath.Join(dataDir, containerID)

	err := ioutil.WriteFile(path, conf, 0600)
	if err != nil {
		return fmt.Errorf("failed to write container data in the path(%q): %v", path, err)
	}

	return err
}

// ReadScratchConf read conf from disk and returns data in byte array
func ReadScratchConf(cRefPath string) ([]byte, error) {
	data, err := ioutil.ReadFile(cRefPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read container data in the path(%q): %v", cRefPath, err)
	}
	return data, err
}

// CleanCachedConf removed cached Conf from disk
func CleanCachedConf(cRefPath string) error {
	if err := os.Remove(cRefPath); err != nil {
		return fmt.Errorf("error removing Conf file %s: %q", cRefPath, err)
	}
	return nil
}
