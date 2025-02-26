package config

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MigrationControllerArgs holds arguments used to configure the MigrationController
type MigrationControllerArgs struct {
	metav1.TypeMeta
}
