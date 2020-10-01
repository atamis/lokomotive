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

package flatcarlinuxupdateoperator

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"

	"github.com/kinvolk/lokomotive/pkg/components"
	"github.com/kinvolk/lokomotive/pkg/components/util"
	"github.com/kinvolk/lokomotive/pkg/k8sutil"
)

const name = "flatcar-linux-update-operator"

func init() {
	components.Register(name, func() components.Component { return &component{} })
}

type component struct{}

func (c *component) LoadConfig(configBody *hcl.Body, evalContext *hcl.EvalContext) hcl.Diagnostics {
	if configBody == nil {
		// This component has no configuration, so don't complain when there is no configuration defined.
		return nil
	}
	return gohcl.DecodeBody(*configBody, evalContext, c)
}

func (c *component) RenderManifests() (map[string]string, error) {
	helmChart, err := components.Chart(name)
	if err != nil {
		return nil, fmt.Errorf("loading chart from assets: %w", err)
	}

	return util.RenderChart(helmChart, name, c.Metadata().Namespace.Name, "")
}

func (c *component) Metadata() components.Metadata {
	return components.Metadata{
		Name: name,
		Namespace: k8sutil.Namespace{
			Name: "reboot-coordinator",
		},
	}
}
