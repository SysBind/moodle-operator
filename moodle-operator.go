package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/kubernetes"
	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	api "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var clientset *kubernetes.Clientset

func main() {
	log.Info("Starting")
	kubeconfig := flag.String("kubeconfig", "", "Path to a kube config. Only required if out-of-cluster.")
	flag.Parse()	
	
	log.Info("kubeconfig is " + *kubeconfig)
	// Create the client config. Use kubeconfig if given, otherwise assume in-cluster.
	config, err := buildConfig(*kubeconfig)
	if err != nil {
		panic(err)
	}

	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Info("error creating Clientset")
		panic(err.Error())
	}

	if pods, err := clientset.Core().Pods("default").List(api.ListOptions{}); err != nil {
		log.Errorf("List Pods of  error:%v", err)
	} else {
		for _, v := range pods.Items {
			log.Infoln(v.GetName(), v.Spec.NodeName, v.Spec.Containers)
		}
	}

}

func buildConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	return rest.InClusterConfig()
}
