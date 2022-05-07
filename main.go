package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/akankshakumari393/convalkontroller/pkg/depkonvalidator"
	depkonv1alpha1 "github.com/akankshakumari393/depkon/pkg/apis/akankshakumari393.dev/v1alpha1"
	"github.com/spf13/pflag"
	admv1beta1 "k8s.io/api/admission/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apiserver/pkg/server"
	"k8s.io/apiserver/pkg/server/options"
	"k8s.io/component-base/cli/globalflag"
)

type Options struct {
	SecureServingOptions options.SecureServingOptions
}

type Config struct {
	SecureServingInfo *server.SecureServingInfo
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

func (o *Options) Config() *Config {
	if err := o.SecureServingOptions.MaybeDefaultWithSelfSignedCerts("0.0.0.0", nil, nil); err != nil {
		panic(err)
	}
	c := Config{}
	o.SecureServingOptions.ApplyTo(&c.SecureServingInfo)
	return &c
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
	c := options.Config()
	// create channel that can be passed to .Serve
	stopCh := server.SetupSignalHandler()
	// create new http handler
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(serveConfigValidation))

	// register validation function to http handler

	// run the https server by calling .Server on config Info
	if serverShutdownCh, _, err := c.SecureServingInfo.Serve(mux, 30*time.Second, stopCh); err != nil {
		panic(err)
	} else {
		<-serverShutdownCh
	}
}

var (
	scheme = runtime.NewScheme()
	codecs = serializer.NewCodecFactory(scheme)
)

func serveConfigValidation(w http.ResponseWriter, r *http.Request) {
	fmt.Println("called serveConfigValidation")
	// read all input into byte
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	gvk := admv1beta1.SchemeGroupVersion.WithKind("AdmissionReview")
	var admissionReview admv1beta1.AdmissionReview
	_, _, err = codecs.UniversalDeserializer().Decode(body, &gvk, &admissionReview)
	if err != nil {
		fmt.Printf("Error %s, converting req body to admission review type", err.Error())
	}
	// get depkon spec from admission review object
	gvkDepkon := depkonv1alpha1.SchemeGroupVersion.WithKind("Depkon")
	var depkon depkonv1alpha1.Depkon
	_, _, err = codecs.UniversalDeserializer().Decode(admissionReview.Request.Object.Raw, &gvkDepkon, &depkon)
	if err != nil {
		fmt.Printf("Error %s, while getting depkon type from admission review", err.Error())
	}

	var response admv1beta1.AdmissionResponse
	allow, err := validateDepkonResource(depkon)
	if !allow || err != nil {
		response = admv1beta1.AdmissionResponse{
			UID:     admissionReview.Request.UID,
			Allowed: false,
			Result: &v1.Status{
				Status:  "Failure",
				Message: fmt.Sprintf("The specified resource %s is not valid", depkon.Name),
				Reason:  v1.StatusReason(err.Error()),
			},
		}
	} else {
		response = admv1beta1.AdmissionResponse{
			UID:     admissionReview.Request.UID,
			Allowed: allow,
		}
	}
	// write the response to response writer
	fmt.Printf("response that we return %+v\n", response)
	// Write Admission Review object to httpResponse Object
	admissionReview.Response = &response
	// convert the Admission Review Object into slice of byte
	res, err := json.Marshal(admissionReview)
	if err != nil {
		fmt.Printf("error %s, while converting response to byte slice", err.Error())
	}
	// write to the response Writer
	_, err = w.Write(res)
	if err != nil {
		fmt.Printf("error %s, writing respnse to responsewriter", err.Error())
	}
}

func validateDepkonResource(depkon depkonv1alpha1.Depkon) (bool, error) {
	return depkonvalidator.CheckIfDepkonValid(depkon)
}
