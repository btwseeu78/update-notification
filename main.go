package main

import (
	"context"
	"fmt"
	"github.com/alecthomas/kingpin/v2"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
	"time"
)

type NodeVersions struct {
	NodeVersions map[string]string `json:"nodeVersions"`
}

var (
	dynUrl         = kingpin.Flag("dynamic-url", "Dynatrace Url").Default("").Envar("DYNA_URL").String()
	projectId      = kingpin.Flag("Project Id", "GCP Project ID").Default("").Envar("GKE_PROJECT_ID").String()
	apiToken       = kingpin.Flag("dynamic-token", "Dynatrace Token").Default("").Envar("DT_API_TOKEN").String()
	oneAgentIngest = kingpin.Flag("use-oneagent", "One Agent injest mode").Default("true").Envar("ENABLE_ONE_AGENT").Bool()
	checkInterval  = kingpin.Flag("check-interval", "interval in which this function check cluster").Default("20").Envar("CHECK_INTERVAL").Duration()
)

func main() {
	config, err := rest.InClusterConfig()
	kingpin.Parse()

	//home := homedir.HomeDir()
	//kubeconfig := filepath.Join(home, ".kube", "config")
	//config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	//if err != nil {
	//	log.Fatalf("Failed to build config from default kubeconfig file: %v", err)
	//}

	//if err != nil {
	//	log.Fatalf("Failed to get in-cluster config: %v", err)
	//}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create clientset: %v", err)
	}
	for {
		status, err := getNodeVersions(clientset)
		if err != nil {
			log.Printf("Failed to get node versions: %v", err)

		}
		if status {
			log.Printf("Node versions: %v", status)
		} else {
			log.Printf("Node versions: %v", status)
		}
		time.Sleep(*checkInterval)
	}

}

func getNodeVersions(clientset *kubernetes.Clientset) (bool, error) {

	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), v1.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to list nodes: %v", err)
	}
	var currVersion string
	var tempVersion string
	for _, node := range nodes.Items {
		tempVersion = node.Status.NodeInfo.KubeletVersion
		if currVersion == "" {
			currVersion = tempVersion
		} else {
			if tempVersion != currVersion {
				fmt.Println("Version Mismatch found")
				return true, nil
				break
			}
		}
		fmt.Println("Node Version not started")

	}
	return false, nil

}
