package app

import (
	"context"
	"fmt"
	"os"

	xschedulerappconfig "github.com/fleezesd/xscheduler/cmd/xscheduler/app/config"
	"github.com/fleezesd/xscheduler/cmd/xscheduler/app/options"
	"github.com/fleezesd/xscheduler/pkg/xscheduler"
	xschedulercontrollers "github.com/fleezesd/xscheduler/pkg/xscheduler/controllers"
	frameworkruntime "github.com/fleezesd/xscheduler/pkg/xscheduler/framework/runtime"
	"github.com/spf13/cobra"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apiserver/pkg/server"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/component-base/cli/globalflag"
	logsapi "k8s.io/component-base/logs/api/v1"
	metricsfeatures "k8s.io/component-base/metrics/features"
	"k8s.io/component-base/term"
	"k8s.io/component-base/version/verflag"
	"k8s.io/klog/v2"
)

func init() {
	utilruntime.Must(logsapi.AddFeatureGates(utilfeature.DefaultMutableFeatureGate))
	utilruntime.Must(metricsfeatures.AddFeatureGates(utilfeature.DefaultMutableFeatureGate))
}

// Option configures a framework registry
type Option func(frameworkruntime.Registry) error

func NewXschedulerCmd(registryOptions ...Option) *cobra.Command {
	opts := options.NewOptions()

	cmd := &cobra.Command{
		Use:   "xscheduler",
		Short: "xscheduler is a scheduler for cloud native scheduler",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCommand(cmd, opts, registryOptions...)
		},
		Args: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}
			return nil
		},
	}

	nfs := opts.Flags
	verflag.AddFlags(nfs.FlagSet("global"))
	globalflag.AddGlobalFlags(nfs.FlagSet("global"), cmd.Name())
	fs := cmd.Flags()
	// setup kubernetes global flagsets for cmd
	for _, f := range nfs.FlagSets {
		fs.AddFlagSet(f)
	}

	// set usage size for terminal
	cols, _, _ := term.TerminalSize(cmd.OutOrStdout())
	cliflag.SetUsageAndHelpFunc(cmd, *nfs, cols)

	// mark config file
	if err := cmd.MarkFlagFilename("config", "yaml", "yml", "json"); err != nil {
		klog.ErrorS(err, "Failed to mark flag filename")
	}
	return cmd
}

func runCommand(cmd *cobra.Command, opts *options.Options, registryOptions ...Option) error {
	if err := logsapi.ValidateAndApply(opts.Logs, utilfeature.DefaultFeatureGate); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	cliflag.PrintFlags(cmd.Flags())

	_, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		stopCh := server.SetupSignalHandler()
		<-stopCh
		cancel()
	}()

	return nil
}

// Setup creates a completed config and a scheduler based on the command args and options
func Setup(ctx context.Context, opts *options.Options, outOfTreeRegistryOptions ...Option) (
	*xschedulerappconfig.CompletedConfig, *xscheduler.Xscheduler, error) {
	if errs := opts.Validate(); len(errs) > 0 {
		return nil, nil, utilerrors.NewAggregate(errs)
	}

	c, err := opts.Config(ctx)
	if err != nil {
		return nil, nil, err
	}

	// Get the completed config
	_ = c.Complete()

	// todo: make transformers for informers and indexers

	// make out of tree registry
	outOfTreeRegistry := make(frameworkruntime.Registry)
	for _, option := range outOfTreeRegistryOptions {
		if err := option(outOfTreeRegistry); err != nil {
			return nil, nil, err
		}
	}
	if err := outOfTreeRegistry.Merge(xschedulercontrollers.NewControllerRegistry()); err != nil {
		return nil, nil, err
	}

	return nil, nil, nil
}
