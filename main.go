package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type PodDetails struct {
	Image     string `yaml:"image"`
	PodName   string `yaml:"name"`
	Namespace string `yaml:"namespace"`
	Success   bool
}

func main() {
	tmpl := template.Must(template.ParseFiles("index.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			tmpl.Execute(w, nil)
			return
		}

		details := PodDetails{
			Image:     r.FormValue("image"),
			PodName:   r.FormValue("podName"),
			Namespace: r.FormValue("namespace"),
			Success:   true,
		}

		tmpl.Execute(w, details)
		podTemplate := template.Must((template.ParseFiles(("pod-template.tmpl"))))
		buf := &bytes.Buffer{}

		fmt.Println("Pod definition\n--------------")
		podTemplate.Execute(buf, details)
		fmt.Println(buf.String())
		createPod()
	})

	http.ListenAndServe(":8080", nil)
}

func createPod() {
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
	for {
		pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

		time.Sleep(10 * time.Second)
	}
}
