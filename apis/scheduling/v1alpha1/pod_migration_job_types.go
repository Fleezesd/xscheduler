package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PodMigrationJob is the Schema for the PodMigrationJob API
// +genclient
// +genclient:nonNamespaced
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:scope=Cluster,shortName=pmj
// +kubebuilder:object:root=true

type PodMigrationJob struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PodMigrationJobSpec   `json:"spec,omitempty"`
	Status PodMigrationJobStatus `json:"status,omitempty"`
}

type PodMigrationJobSpec struct {
}

type PodMigrationJobStatus struct {
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PodMigrationJobList contains a list of PodMigrationJob
type PodMigrationJobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PodMigrationJob `json:"items"`
}
