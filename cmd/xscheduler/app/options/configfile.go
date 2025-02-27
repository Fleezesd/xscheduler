package options

import (
	"fmt"
	"os"

	xschedulerconfig "github.com/fleezesd/xscheduler/pkg/xscheduler/apis/config"
	xschedulerconfigscheme "github.com/fleezesd/xscheduler/pkg/xscheduler/apis/config/scheme"
)

func loadConfigFromFile(file string) (*xschedulerconfig.XschedulerConfiguration, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return loadConfig(data)
}

func loadConfig(data []byte) (*xschedulerconfig.XschedulerConfiguration, error) {
	// The UniversalDecoder runs defaulting and returns the internal type by default.
	// Let the decoder judge automatically gvk and version
	obj, gvk, err := xschedulerconfigscheme.Codecs.UniversalDecoder().Decode(data, nil, nil)
	if err != nil {
		return nil, err
	}

	if cfgObj, ok := obj.(*xschedulerconfig.XschedulerConfiguration); ok {
		// the field will be cleared later by API machinery during
		// conversion. See DeschedulerConfiguration internal type definition for
		// more details.
		cfgObj.TypeMeta.APIVersion = gvk.GroupVersion().String()
		return cfgObj, nil
	}
	return nil, fmt.Errorf("couldn't decode as DeschedulerConfiguration, got %s: ", gvk)
}
