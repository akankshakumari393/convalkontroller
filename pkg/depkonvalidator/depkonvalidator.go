package depkonvalidator

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	depkonv1alpha1 "github.com/akankshakumari393/depkon/pkg/apis/akankshakumari393.dev/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func CheckIfDepkonValid(depkon depkonv1alpha1.Depkon) (bool, error) {
	homeDir := os.Getenv("HOME")
	kubeconfigFile := homeDir + "/.kube/config"
	kubeconfig := flag.String("kubeconfig", kubeconfigFile, "Kubeconfig File location")
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		// handle error
		log.Printf("erorr %s building config from flags\n", err.Error())
		config, err = rest.InClusterConfig()
		if err != nil {
			log.Printf("error %s, getting inclusterconfig", err.Error())
			return false, err
		}
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		// handle error
		log.Printf("error %s, kubernetes clientset\n", err.Error())
		return false, err
	}
	// check for configMap
	_, err = clientset.CoreV1().ConfigMaps(depkon.Namespace).Get(context.Background(), depkon.Spec.ConfigmapRef, metav1.GetOptions{})
	if err != nil {
		// handle error
		fmt.Printf("error %s, configmap not present \n", err.Error())
		return false, err
	}
	for _, deployment := range depkon.Spec.DeploymentRef {
		// check for Deployment
		_, err = clientset.AppsV1().Deployments(depkon.Namespace).Get(context.Background(), deployment, metav1.GetOptions{})
		if err != nil {
			// handle error
			fmt.Printf("error %s, Deployment not present \n", err.Error())
			return false, err
		}
	}
	return true, nil
}
