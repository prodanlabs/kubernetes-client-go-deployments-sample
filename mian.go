package main

import (
	"flag"
	"fmt"
	"github.com/prodanlabs/kubernetes-client-go-deployments-sample/client"
	"github.com/prodanlabs/kubernetes-client-go-deployments-sample/harbor"
	"github.com/prodanlabs/kubernetes-client-go-deployments-sample/utils"
	"os"
	"strings"
	"time"
)

func init() {
	utils.Logs()
}
func main() {

	application := flag.String("apps", os.Getenv("APPLICATION_GROUP"), "namespace_deploy_container.Examples: kube-system.busybox.busybox,default.for.test.")
	clusterEnv := flag.String("cluster", os.Getenv("CLUSTER_ENV"), "out-of or in cluster.")
	harborUser := flag.String("user", os.Getenv("HARBOR_USER"), "username for harbor.")
	harborPass := flag.String("pass", os.Getenv("HARBOR_PASS"), "passwd for harbor.")
	registry := flag.String("registry", os.Getenv("HARBOR_URL"), "url for harbor.")
	projectName := flag.String("project", os.Getenv("HARBOR_PROJECT"), "app harbor projectName.")
	version := flag.Bool("v", false, "version")
	flag.Parse()
	if *version {
		fmt.Fprint(os.Stderr, utils.VerSion())
		os.Exit(1)
	}
	var clientset, err = client.Connect(*clusterEnv)
	if err != nil {
		utils.Error.Panic(err)
		panic(err.Error())
	}
	for {
		sep := ","
		apps := strings.Split(*application, sep)
		for _, appName := range apps {
			sep := "."
			arr := strings.Split(appName, sep)
			utils.Info.Printf("namespace: %v  deployment: %v container: %v", arr[0], arr[1], arr[2])
			repositories := *projectName + "/" + arr[2]
			imageName, err := harbor.GetImages(*harborUser, *harborPass, *registry, repositories)
			if err != nil {
				utils.Info.Println(err)
			}
			if imageName != "0" {
				client.UpdateDeployments(clientset, arr[0], arr[1], arr[2], imageName)
			}

		}
		time.Sleep(120 * time.Second)
	}
}
