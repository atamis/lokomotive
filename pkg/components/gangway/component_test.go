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

package gangway_test

import (
	"testing"

	"github.com/kinvolk/lokomotive/pkg/components/gangway"
	"github.com/kinvolk/lokomotive/pkg/components/util"
)

const name = "gangway"

func TestRenderManifest(t *testing.T) { //nolint:funlen
	tests := []struct {
		desc    string
		hcl     string
		wantErr bool
	}{
		{
			desc: "Valid config",
			hcl: `
component "gangway" {
  cluster_name   = "lokomotive"
  ingress_host   = "foo"
  session_key    = "sessionkey"
  api_server_url = "apiserverurl"
  authorize_url  = "authorizeurl"
  token_url      = "tokenurl"
  client_id      = "clientid"
  client_secret  = "clientsecret"
  redirect_url   = "redirecturl"
}
			`,
		},
		{
			desc: "invalid config",
			hcl: `
component "gangway" {
  certmanager_cluster_issuer = "letsencrypt-staging"
}
			`,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		b, d := util.GetComponentBody(tc.hcl, name)
		if d != nil {
			t.Errorf("%s - Error getting component body: %v", tc.desc, d)
		}

		c := gangway.NewConfig()

		d = c.LoadConfig(b, nil)

		if !tc.wantErr && d.HasErrors() {
			t.Errorf("%s - Valid config should not return error, got: %s", tc.desc, d)
		}

		if tc.wantErr && !d.HasErrors() {
			t.Errorf("%s - Wrong config should have returned an error", tc.desc)
		}

		m, err := c.RenderManifests()
		if err != nil {
			t.Errorf("%s - Rendering manifests with valid config should succeed, got: %s", tc.desc, err)
		}

		if len(m) == 0 {
			t.Errorf("%s - Rendered manifests shouldn't be empty", tc.desc)
		}
	}
}
