package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var kubeconfig *string

	// home是家目录，如果能取得家目录的值，就可以用来做默认值
	if home := homedir.HomeDir(); home != "" {
		// 如果输入了kubeconfig参数，该参数的值就是kubeconfig文件的绝对路径，
		// 如果没有输入kubeconfig参数，就用默认路径~/.kube/config
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		// 如果取不到当前用户的家目录，就没办法设置kubeconfig的默认目录了，只能从入参中取
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	flag.Parse()

	// 从本机加载kubeconfig配置文件，因此第一个参数为空字符串
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)

	// kubeconfig加载失败就直接退出了
	if err != nil {
		panic(err.Error())
	}

	// 实例化 ClientSet
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	namespace := "default"

	// 查询 default 下的 pods 部门资源信息
	pods, err := clientset.
		CoreV1().        // 实例化资源客户端，这里标识实例化 CoreV1Client
		Pods(namespace). // 选择 namespace，为空则表示所有 Namespace
		// Namespaces(). // 若是查询 namespace resource
		List(context.TODO(), metav1.ListOptions{}) // 查询 pods 列表

	if err != nil {
		panic(err.Error())
	}

	// 表头
	fmt.Printf("namespace\t status\t\t name\n")
	// 每个pod都打印namespace、status.Phase、name三个字段
	for _, d := range pods.Items {
		fmt.Printf("%v\t\t %v\t %v\n",
			d.Namespace,
			d.Status.Phase,
			d.Name)
	}

	deploymentsClient, err := clientset.
		AppsV1().
		Deployments(namespace).
		List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("\nnamespace\t Replicas\t name\n")
	for _, d := range deploymentsClient.Items {
		fmt.Printf("%v\t\t %v\t\t\t %v\n",
			d.Namespace,
			*d.Spec.Replicas,
			d.Name)
	}

}
