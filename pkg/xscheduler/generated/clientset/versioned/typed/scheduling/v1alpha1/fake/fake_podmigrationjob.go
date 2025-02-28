/*
Copyright 2025.

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
// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1alpha1 "github.com/fleezesd/xscheduler/apis/scheduling/v1alpha1"
	schedulingv1alpha1 "github.com/fleezesd/xscheduler/pkg/xscheduler/generated/clientset/versioned/typed/scheduling/v1alpha1"
	gentype "k8s.io/client-go/gentype"
)

// fakePodMigrationJobs implements PodMigrationJobInterface
type fakePodMigrationJobs struct {
	*gentype.FakeClientWithList[*v1alpha1.PodMigrationJob, *v1alpha1.PodMigrationJobList]
	Fake *FakeSchedulingV1alpha1
}

func newFakePodMigrationJobs(fake *FakeSchedulingV1alpha1) schedulingv1alpha1.PodMigrationJobInterface {
	return &fakePodMigrationJobs{
		gentype.NewFakeClientWithList[*v1alpha1.PodMigrationJob, *v1alpha1.PodMigrationJobList](
			fake.Fake,
			"",
			v1alpha1.SchemeGroupVersion.WithResource("podmigrationjobs"),
			v1alpha1.SchemeGroupVersion.WithKind("PodMigrationJob"),
			func() *v1alpha1.PodMigrationJob { return &v1alpha1.PodMigrationJob{} },
			func() *v1alpha1.PodMigrationJobList { return &v1alpha1.PodMigrationJobList{} },
			func(dst, src *v1alpha1.PodMigrationJobList) { dst.ListMeta = src.ListMeta },
			func(list *v1alpha1.PodMigrationJobList) []*v1alpha1.PodMigrationJob {
				return gentype.ToPointerSlice(list.Items)
			},
			func(list *v1alpha1.PodMigrationJobList, items []*v1alpha1.PodMigrationJob) {
				list.Items = gentype.FromPointerSlice(items)
			},
		),
		fake,
	}
}
