package options

import "k8s.io/apimachinery/pkg/runtime"

var scheme = runtime.NewScheme()

func init() {
	// todo: add CRD scheme to this scheme
}
