package main

import (
	"context"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	apiv1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type PodDetails struct {
	Image     string
	PodName   string
	Namespace string
	Success   bool
	Failure   bool
}

func main() {
	tmpl := template.Must(template.ParseFiles("views/index.html"))
	clientset, err := createKubeClient()

	if err != nil {
		panic(err.Error())
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Print("Starting request")
		if r.Method != http.MethodPost {
			tmpl.Execute(w, nil)
			return
		}

		details := PodDetails{
			Image:     r.FormValue("image"),
			PodName:   r.FormValue("podName"),
			Namespace: r.FormValue("namespace"),
		}

		status, _ := createPod(clientset, r.FormValue("image"), r.FormValue("podName"), r.FormValue("namespace"))

		if status == "Error" {
			details.Failure = true
		} else {
			details.Success = true
		}
		log.Print("Rendering page")
		tmpl.Execute(w, details)
	})

	log.Print("Application is available at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func createKubeClient() (kubernetes.Clientset, error) {
	log.Print("Creating kube client")
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return *clientset, err
}

func createPodObject(image, podName, namespace string) *apiv1.Pod {
	log.Print("Creating Pod Object")
	return &apiv1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: namespace,
			Labels: map[string]string{
				"app": podName,
			},
		},
		Spec: apiv1.PodSpec{
			Containers: []apiv1.Container{
				{
					Name:            image,
					Image:           image,
					ImagePullPolicy: apiv1.PullIfNotPresent,
				},
			},
		},
	}
}

func createPod(clientset kubernetes.Clientset, image, podName, namespace string) (string, error) {
	var err error
	pod := createPodObject(image, podName, namespace)
	log.Print("Creating Pod")
	_, err = clientset.CoreV1().Pods(pod.Namespace).Create(context.TODO(), pod, metav1.CreateOptions{})

	if err != nil {
		log.Print(err)
		return "Error", err
	}
	fmt.Println("Pod created successfully...")
	return "Deployed", nil
}
