// Copyright 2017 Google Inc. All rights reserved.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v2

import (
	"fmt"

	"github.com/GoogleCloudPlatform/runtimes-common/ctc_lib"
	"github.com/GoogleContainerTools/container-structure-test/pkg/drivers"
	types "github.com/GoogleContainerTools/container-structure-test/pkg/types/unversioned"
	"github.com/GoogleContainerTools/container-structure-test/pkg/utils"
)

type MetadataTest struct {
	Cmd          *[]string          `yaml:"cmd"`
	Entrypoint   *[]string          `yaml:"entrypoint"`
	Env          []types.KeyValPair `yaml:"env"`
	ExposedPorts []string           `yaml:"exposedPorts"`
	Labels       []types.KeyValPair `yaml:"labels"`
	OnBuild      []string           `yaml:"onBuild"`
	// Shell        []string           `yaml:"shell"`
	StopSignal   string             `yaml:"stopSignal"`
	User         string             `yaml:"user"`
	Volumes      []string           `yaml:"volumes"`
	Workdir      string             `yaml:"workdir"`
}

func (mt MetadataTest) LogName() string {
	return "Metadata Test"
}

func (mt MetadataTest) Validate() error {
	if mt.hasMissingKeys(mt.Env) {
		return fmt.Errorf("Environment variable key cannot be empty")
	}
	if mt.hasMissingKeys(mt.Labels) {
		return fmt.Errorf("Label key cannot be empty")
	}
	if mt.hasBlankStrings(mt.ExposedPorts) {
		return fmt.Errorf("Port cannot be empty")
	}
	if mt.hasBlankStrings(mt.Volumes) {
		return fmt.Errorf("Volume cannot be empty")
	}
	return nil
}

func (mt MetadataTest) hasBlankStrings(strings []string) bool {
	for _, s := range strings {
		if s == "" {
			return true
		}
	}
	return false
}

func (mt MetadataTest) hasMissingKeys(pairs []types.KeyValPair) bool {
	for _, pair := range pairs {
		if pair.Key == "" {
			return true
		}
	}
	return false
}

func (mt MetadataTest) testStringRegex(
	fieldName string,
	expectedPattern string,
	actualValue string,
	result *types.TestResult) {
	if !utils.CompileAndRunRegex(expectedPattern, actualValue, true) {
		result.Errorf("Image %s \"%s\" does not match expected value: \"%s\"", fieldName, actualValue, expectedPattern)
		result.Fail()
	}
}

func (mt MetadataTest) testListOfValues(
	fieldName string,
	expectedValues []string,
	actualValues []string,
	result *types.TestResult) {

	for _,  value := range expectedValues {
		if !utils.ValueInList(value, actualValues) {
			result.Errorf("Image %s \"%s\" not found in config", fieldName, value)
			result.Fail()
		}
	}
}

func (mt MetadataTest) testKeyValPairs(
	fieldName string,
	expectedKeyValPairs []types.KeyValPair,
	actualKeyValPairs map[string]string,
	result *types.TestResult) {
	for _, expectedPair := range expectedKeyValPairs {
		if actualValue, hasKey := actualKeyValPairs[expectedPair.Key]; hasKey {
			if !utils.CompileAndRunRegex(expectedPair.Value, actualValue, true) {
				result.Errorf("Image %s \"%s\"=\"%s\" does not match expected value: \"%s\"",
					fieldName, expectedPair.Key, actualValue, expectedPair.Value)
				result.Fail()
			}
		} else {
			result.Errorf("Image %s \"%s\" not found", fieldName, expectedPair.Key)
			result.Fail()
		}
	}
}

func (mt MetadataTest) testArgsArray(
	fieldName string, 
	expectedArgsArray *[]string,
	actualArgsArray *[]string,
	result *types.TestResult) {

	if expectedArgsArray != nil {
		if len(*expectedArgsArray) != len(*actualArgsArray) {
			result.Errorf("Image %s %v does not match expected value: %v",
				fieldName, *actualArgsArray, *expectedArgsArray)
			result.Fail()
		} else {
			for i := range *expectedArgsArray {
				if (*expectedArgsArray)[i] != (*actualArgsArray)[i] {
					result.Errorf("Image %s[%d] \"%s\" does not match expected value: \"%s\"",
						fieldName, i, (*actualArgsArray)[i], (*expectedArgsArray)[i])
					result.Fail()
				}
			}
		}
	}
}

func (mt MetadataTest) Run(driver drivers.Driver) *types.TestResult {
	result := &types.TestResult{
		Name: mt.LogName(),
		Pass: true,
	}
	ctc_lib.Log.Debug(mt.LogName())
	imageConfig, err := driver.GetConfig()
	if err != nil {
		result.Errorf("Error retrieving image config: %s", err.Error())
		result.Fail()
		return result
	}
	ctc_lib.Log.Debugf("actual imageConfig: %+v", imageConfig)

	// []KeyValuePair tests
	mt.testKeyValPairs("envvar", mt.Env, imageConfig.Env, result)
	mt.testKeyValPairs("label", mt.Labels, imageConfig.Labels, result)

	// []string tests - argument array
	mt.testArgsArray("cmd", mt.Cmd, &imageConfig.Cmd, result)
	mt.testArgsArray("entrypoint", mt.Entrypoint, &imageConfig.Entrypoint, result)

	// string tests
	mt.testStringRegex("workdir", mt.Workdir, imageConfig.Workdir, result)
	mt.testStringRegex("user", mt.User, imageConfig.User, result)
	mt.testStringRegex("stopSignal", mt.StopSignal, imageConfig.StopSignal, result)

	// []string tests - list of values
	mt.testListOfValues("ExposedPort", mt.ExposedPorts, imageConfig.ExposedPorts, result)
	mt.testListOfValues("Volume", mt.Volumes, imageConfig.Volumes, result)

	return result
}
