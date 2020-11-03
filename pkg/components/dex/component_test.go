// Copyright 2020 The Lokomotive Authors
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

package dex_test

import (
	"testing"

	"github.com/kinvolk/lokomotive/pkg/components/dex"
	"github.com/kinvolk/lokomotive/pkg/components/util"
)

const name = "dex"

var tests = []struct {
	desc    string
	hcl     string
	wantErr bool
}{
	{
		desc: "first valid config",
		hcl: `
component "dex" {
  ingress_host = "foo"
  issuer_host  = "bar"
  connector "github" {
    id   = "github"
    name = "Github"
    config {
      client_id       = "clientid"
      client_secret   = "clientsecret"
      redirect_uri    = "redirecturi"
      team_name_field = "slug"
      org {
        name = "kinvolk"
        teams = [
          "lokomotive-developers-1",
        ]
      }
    }
  }

  static_client {
    name          = "gangway"
    id            = "gangway id"
    secret        = "gangway secret"
    redirect_uris = ["redirecturis"]
  }
}
			`,
	},
	{
		desc: "second valid config",
		hcl: `
component "dex" {
  ingress_host = "foo"
  issuer_host  = "bar"
  connector "github" {
    id   = "github"
    name = "Github"
    config {
      client_id       = "clientid"
      client_secret   = "clientsecret"
      redirect_uri    = "redirecturi"
      team_name_field = "slug"
      org {
        name = "kinvolk"
        teams = [
          "lokomotive-developers-2",
        ]
      }
    }
  }

  static_client {
    name          = "gangway"
    id            = "gangway id"
    secret        = "gangway secret"
    redirect_uris = ["redirecturis"]
  }
}
			`,
	},
	{
		desc: "invalid config",
		hcl: `
component "dex" {
  ingress_host = "NodePort"
}
			`,
		wantErr: true,
	},
}

func TestRenderManifest(t *testing.T) {
	for _, tc := range tests {
		b, d := util.GetComponentBody(tc.hcl, name)
		if d != nil {
			t.Fatalf("%s - Error getting component body: %v", tc.desc, d)
		}

		c := dex.NewConfig()

		d = c.LoadConfig(b, nil)

		if !tc.wantErr && d.HasErrors() {
			t.Fatalf("%s - Valid config should not return error, got: %s", tc.desc, d)
		}

		if tc.wantErr && !d.HasErrors() {
			t.Fatalf("%s - Wrong config should have returned an error", tc.desc)
		}

		m, err := c.RenderManifests()
		if err != nil {
			t.Fatalf("%s - Rendering manifests with valid config should succeed, got: %s", tc.desc, err)
		}

		if len(m) == 0 {
			t.Fatalf("%s - Rendered manifests shouldn't be empty", tc.desc)
		}
	}
}

func TestDeploymentAnnotationHashChange(t *testing.T) {
	deploymentArr := []string{}

	for _, tc := range tests {
		b, d := util.GetComponentBody(tc.hcl, name)
		if d != nil {
			t.Fatalf("%s - Error getting component body: %v", tc.desc, d)
		}

		c := dex.NewConfig()

		d = c.LoadConfig(b, nil)
		if !tc.wantErr && d.HasErrors() {
			t.Fatalf("%s - Valid config should not return error, got: %s", tc.desc, d)
		}

		m, err := c.RenderManifests()
		if err != nil {
			t.Fatalf("%s - Rendering manifests with valid config should succeed, got: %s", tc.desc, err)
		}

		deploymentArr = append(deploymentArr, m["dex/templates/deployment.yaml"])
	}

	if deploymentArr[0] == deploymentArr[1] {
		t.Fatalf("expected checksum/configmap hash to be different.")
	}
}
