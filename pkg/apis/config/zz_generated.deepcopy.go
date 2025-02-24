//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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
// Code generated by deepcopy-gen. DO NOT EDIT.

package config

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *XschedulerConfiguration) DeepCopyInto(out *XschedulerConfiguration) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.LeaderElection = in.LeaderElection
	out.ClientConnection = in.ClientConnection
	out.DebuggingConfiguration = in.DebuggingConfiguration
	out.XschedulingInterval = in.XschedulingInterval
	if in.Profiles != nil {
		in, out := &in.Profiles, &out.Profiles
		*out = make([]XschedulerProfile, len(*in))
		copy(*out, *in)
	}
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = new(v1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
	if in.MaxNoOfPodsToEvictPerNode != nil {
		in, out := &in.MaxNoOfPodsToEvictPerNode, &out.MaxNoOfPodsToEvictPerNode
		*out = new(uint)
		**out = **in
	}
	if in.MaxNoOfPodsToEvictPerNamespace != nil {
		in, out := &in.MaxNoOfPodsToEvictPerNamespace, &out.MaxNoOfPodsToEvictPerNamespace
		*out = new(uint)
		**out = **in
	}
	if in.MaxNoOfPodsToEvictTotal != nil {
		in, out := &in.MaxNoOfPodsToEvictTotal, &out.MaxNoOfPodsToEvictTotal
		*out = new(uint)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new XschedulerConfiguration.
func (in *XschedulerConfiguration) DeepCopy() *XschedulerConfiguration {
	if in == nil {
		return nil
	}
	out := new(XschedulerConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *XschedulerConfiguration) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *XschedulerProfile) DeepCopyInto(out *XschedulerProfile) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new XschedulerProfile.
func (in *XschedulerProfile) DeepCopy() *XschedulerProfile {
	if in == nil {
		return nil
	}
	out := new(XschedulerProfile)
	in.DeepCopyInto(out)
	return out
}
