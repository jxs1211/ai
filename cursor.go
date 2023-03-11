package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type K8sClient struct {
	clientset *kubernetes.Clientset
}

func NewClientSet() (*K8sClient, error) {
	flag.Parse()
	home := homedir.HomeDir()
	kubeconfig := flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &K8sClient{clientset: clientset}, nil
}

func (k8s *K8sClient) CreatePod(name, namespace string) (*corev1.Pod, error) {
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "sample-container",
					Image: "nginx",
				},
			},
		},
	}

	createdPod, err := k8s.clientset.CoreV1().Pods(namespace).Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	return createdPod, nil
}

func (k8s *K8sClient) DeletePod(name, namespace string) error {
	err := k8s.clientset.CoreV1().Pods(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	return nil
}

func main() {
	c, err := NewClientSet()
	if err != nil {
		log.Fatalf("Failed to create pod: %v", err)
	}

	namespace := "dev"
	// pod, err := c.CreatePod("my-pod", namespace)
	// if err != nil {
	// 	log.Fatalf("Failed to create pod: %v", err)
	// }

	// fmt.Printf("Created pod %s in namespace %s\n", pod.Name, namespace)

	err = c.DeletePod("my-pod", namespace)
	if err != nil {
		log.Fatalf("Failed to create pod: %v", err)
	}
	fmt.Printf("Delete pod in namespace %s\n", namespace)
}
