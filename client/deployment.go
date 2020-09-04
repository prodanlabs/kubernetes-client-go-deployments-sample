package client

import (
	"context"
	"github.com/prodanlabs/kubernetes-client-go-deployments-sample/utils"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func UpdateDeployments(client *kubernetes.Clientset, namespace, deploymentName, appName, imageName string) {

	deployment, err := client.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		utils.Error.Println(err)
	}
	if errors.IsNotFound(err) {
		utils.Error.Println("Deployment not found")
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		utils.Error.Printf("Error get deployment %v", statusError.ErrStatus.Message)
	} else if err != nil {
		utils.Error.Println(err)
	} else {
		utils.Info.Println("Found deployment")
		name := deployment.GetName()
		utils.Info.Println("container name ->", name)
		containers := &deployment.Spec.Template.Spec.Containers
		found := false
		for i := range *containers {

			c := *containers
			if c[i].Name == appName {
				found = true
				utils.Info.Println("Current container image ->", c[i].Image)
				utils.Info.Println("New image of harbor registry ->", imageName)
				if c[i].Image == imageName {
					utils.Info.Println("The application container image is new")
					continue
				} else {
					c[i].Image = imageName
					if found == false {
						utils.Error.Println("The application container not exist in the deployment pods.")
						utils.Info.Println("test")
						_, err := client.AppsV1().Deployments(namespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
						if err != nil {
							utils.Error.Panic(err)
						}
					}
				}

			}
		}

	}
}
