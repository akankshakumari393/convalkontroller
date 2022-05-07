package main

import (
	"os"

	"github.com/spf13/pflag"
	"k8s.io/apiserver/pkg/server"
	"k8s.io/apiserver/pkg/server/options"
	"k8s.io/component-base/cli/globalflag"
)

type Options struct {
	SecureServingOptions options.SecureServingOptions
}

type Config struct {
	SecureServingInfo server.SecureServingInfo
}

const (
	controller = "con-valkontroller"
)

func (options *Options) AddFlagSet(fs *pflag.FlagSet) {
	options.SecureServingOptions.AddFlags(fs)
}

func NewDefaultOption() *Options {
	options := &Options{
		SecureServingOptions: *options.NewSecureServingOptions(),
	}
	options.SecureServingOptions.BindPort = 8443
	options.SecureServingOptions.ServerCert.PairName = controller
	return options
}

func main() {
	// initialize default option
	options := NewDefaultOption()
	// create a new flag set
	fs := pflag.NewFlagSet(controller, pflag.ExitOnError)
	// Add global flag like --help to the flag set
	globalflag.AddGlobalFlags(fs, controller)
	// add the flagset to the options
	options.AddFlagSet(fs)
	// parse flagset
	if err := fs.Parse(os.Args); err != nil {
		panic(err)
	}

	// create config from options

	// create channel that can be passed to .Serve

	// create new http handler

	// register validation function to http handler

	// run the https server by calling .Server on config Info

}
