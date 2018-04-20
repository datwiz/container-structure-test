// Copyright 2018 Google Inc. All rights reserved.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package unversioned

import (
	"fmt"
)

type EnvVar struct {
	Key   string
	Value string
}

type Label struct {
	Key   string
	Value string
}

type Config struct {
	Cmd          []string
	Entrypoint   []string
	Env          map[string]string
	ExposedPorts []string
	Labels       map[string]string
	OnBuild      []string
	// Shell        []string
	StopSignal   string
	User         string
	Volumes      []string
	Workdir      string
}

type FlattenedConfig struct {
	Cmd          []string            `json:"Cmd"`
	Entrypoint   []string            `json:"Entrypoint"`
	Env          []string            `json:"Env"`
	ExposedPorts map[string][]string `json:"ExposedPorts"`
	Labels       []string            `json:"Labels"`
	OnBuild      []string            `json:"OnBuild"`
	// Shell        []string            `json: "Shell"`
	StopSignal   string              `json:"StopSignal"`
	User         string              `json:"User"`
	Volumes      map[string]string   `json:"Volumes"`
	Workdir      string              `json:"WorkingDir"`
}

type FlattenedMetadata struct {
	Config FlattenedConfig `json:"config"`
}

type TestResult struct {
	Name   string
	Pass   bool
	Stdout string
	Stderr string
	Errors []string
}

func (t *TestResult) Error(s string) {
	t.Errors = append(t.Errors, s)
}

func (t *TestResult) Errorf(s string, args ...interface{}) {
	t.Errors = append(t.Errors, fmt.Sprintf(s, args...))
}

func (t *TestResult) Fail() {
	t.Pass = false
}
