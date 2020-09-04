package client

import (
	"flag"
	"github.com/prodanlabs/kubernetes-client-go-deployments-sample/utils"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
)

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
func Connect(env string) (*kubernetes.Clientset, error) {
	if env == "out-of-cluster" {
		var kubeconfig *string
		if home := homeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			utils.Error.Panic(err)
			panic(err.Error())
		}

		// create the clientset
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			utils.Error.Panic(err)
			panic(err.Error())
		}
		return clientset, err
	} else {
		config, err := rest.InClusterConfig()
		if err != nil {
			utils.Error.Panic(err)
			panic(err.Error())
		}
		// creates the clientset
		clientset, err := kubernetes.NewForConfig(config)
		return clientset, err
	}
}
