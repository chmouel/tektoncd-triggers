/*
Copyright 2019 The Tekton Authors

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

package v1alpha1_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	"github.com/tektoncd/triggers/pkg/apis/triggers/v1alpha1"
)

func TestEventListenerSetDefaults(t *testing.T) {
	tests := []struct {
		name string
		in   *v1alpha1.EventListener
		want *v1alpha1.EventListener
		wc   func(context.Context) context.Context
	}{{
		name: "default binding",
		in: &v1alpha1.EventListener{
			Spec: v1alpha1.EventListenerSpec{
				Triggers: []v1alpha1.EventListenerTrigger{{
					Bindings: []*v1alpha1.EventListenerBinding{
						{
							Name: "binding",
						},
						{
							Name: "namespace-binding",
							Kind: v1alpha1.NamespacedTriggerBindingKind,
						},
						{
							Name: "cluster-binding",
							Kind: v1alpha1.ClusterTriggerBindingKind,
						},
					},
				}},
			},
		},
		want: &v1alpha1.EventListener{
			Spec: v1alpha1.EventListenerSpec{
				Triggers: []v1alpha1.EventListenerTrigger{{
					Bindings: []*v1alpha1.EventListenerBinding{
						{
							Name: "binding",
							Kind: v1alpha1.NamespacedTriggerBindingKind,
						},
						{
							Name: "namespace-binding",
							Kind: v1alpha1.NamespacedTriggerBindingKind,
						},
						{
							Name: "cluster-binding",
							Kind: v1alpha1.ClusterTriggerBindingKind,
						},
					},
				}},
			},
		},
	}, {
		name: "with upgrade context - deprecated params",
		in: &v1alpha1.EventListener{
			Spec: v1alpha1.EventListenerSpec{
				Triggers: []v1alpha1.EventListenerTrigger{{
					DeprecatedParams: []pipelinev1.Param{{
						Name: "param-name",
						Value: pipelinev1.ArrayOrString{
							Type:      "string",
							StringVal: "static",
						},
					},
					}},
				}},
		},
		want: &v1alpha1.EventListener{
			Spec: v1alpha1.EventListenerSpec{
				Triggers: []v1alpha1.EventListenerTrigger{{}},
			},
		},
		wc: v1alpha1.WithUpgradeViaDefaulting,
	}}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.in
			ctx := context.Background()
			if tc.wc != nil {
				ctx = tc.wc(ctx)
			}
			got.SetDefaults(ctx)

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("SetDefaults (-want, +got) = %v", diff)
			}
		})
	}
}
