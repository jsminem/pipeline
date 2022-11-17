/*
Copyright 2021 The Tekton Authors

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

package config_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tektoncd/pipeline/pkg/apis/config"
	test "github.com/tektoncd/pipeline/pkg/reconciler/testing"
	"github.com/tektoncd/pipeline/test/diff"
	"k8s.io/apimachinery/pkg/util/sets"
)

func TestNewTrustedResourcesFromConfigMap(t *testing.T) {
	type testCase struct {
		expectedConfig *config.TrustedResources
		fileName       string
	}

	testCases := []testCase{
		{
			expectedConfig: &config.TrustedResources{
				Keys: sets.NewString("/etc/verification-secrets/cosign.pub", "/etc/verification-secrets/cosign2.pub"),
			},
			fileName: config.GetTrustedResourcesConfigName(),
		},
		{
			expectedConfig: &config.TrustedResources{
				Keys: sets.NewString("/etc/verification-secrets/cosign.pub", "/etc/verification-secrets/cosign2.pub"),
			},
			fileName: "config-trusted-resources",
		},
	}

	for _, tc := range testCases {
		verifyConfigFileWithExpectedTrustedResourcesConfig(t, tc.fileName, tc.expectedConfig)
	}
}

func TestNewTrustedResourcesFromEmptyConfigMap(t *testing.T) {
	MetricsConfigEmptyName := "config-trusted-resources-empty"
	expectedConfig := &config.TrustedResources{
		Keys: sets.NewString(config.DefaultPublicKeyPath),
	}
	verifyConfigFileWithExpectedTrustedResourcesConfig(t, MetricsConfigEmptyName, expectedConfig)
}

func verifyConfigFileWithExpectedTrustedResourcesConfig(t *testing.T, fileName string, expectedConfig *config.TrustedResources) {
	cm := test.ConfigMapFromTestFile(t, fileName)
	if ab, err := config.NewTrustedResourcesConfigFromConfigMap(cm); err == nil {
		if d := cmp.Diff(ab, expectedConfig); d != "" {
			t.Errorf("Diff:\n%s", diff.PrintWantGot(d))
		}
	} else {
		t.Errorf("NewTrustedResourcesFromConfigMap(actual) = %v", err)
	}
}