package main

import (
	"flag"
	"fmt"
	"sort"

	"path/filepath"

	termbox "github.com/nsf/termbox-go"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var kubeconfig *string

func init() {
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
}

func runPodList() {
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

	pods, err := clientset.CoreV1().Pods("").List( /*context.TODO(),*/ metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	namespace := "default"

	pods, _ = clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{})
	items := pods.DeepCopy().Items
	sort.Sort(byAge(items))
	y := 0
	for _, pod := range items {
		tbprint(0, y, coldef, coldef, pod.Name)
		w, _ := termbox.Size()

		status := string(pod.Status.Phase)
		if pod.Status.Reason != "" {
			status = pod.Status.Reason
		}

		desc := getDescription(pod)
		if desc != "" {
			status += fmt.Sprintf(" (%s)", desc)
		}

		x := w - 10 - len(status)
		tbprint(x, y, coldef, coldef, status)

		tm := pod.Status.StartTime
		if tm != nil {
			ago := timeSince(tm.Time)
			x = w - len(ago)
			tbprint(x, y, coldef, coldef, ago)
		}

		y++
	}
}
