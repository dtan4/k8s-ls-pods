package main

import (
	"fmt"
	"os"

	k8s "github.com/dtan4/k8s-ls-pods/kubernetes"
	"github.com/k0kubun/pp"
	flag "github.com/spf13/pflag"
)

func main() {
	var (
		allNamespaces bool
		kubeContext   string
		inCluster     bool
		kubeconfig    string
		labels        string
		namespace     string
		verbose       bool
	)

	flags := flag.NewFlagSet("k8s-pod-notifier", flag.ExitOnError)
	flags.Usage = func() {
		flags.PrintDefaults()
	}

	flags.BoolVar(&allNamespaces, "all-namespaces", false, "List pods across all namespaces")
	flags.StringVar(&kubeContext, "context", "", "Kubernetes context")
	flags.BoolVar(&inCluster, "in-cluster", false, "Execute in Kubernetes cluster")
	flags.StringVar(&kubeconfig, "kubeconfig", "", "Path of kubeconfig")
	flags.StringVarP(&labels, "labels", "l", "", "Label filter query")
	flags.StringVarP(&namespace, "namespace", "n", "", "Kubernetes namespace")
	flags.BoolVarP(&verbose, "verbose", "v", false, "Verbose output")

	if err := flags.Parse(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if kubeconfig == "" {
		if os.Getenv("KUBECONFIG") != "" {
			kubeconfig = os.Getenv("KUBECONFIG")
		} else {
			kubeconfig = k8s.DefaultConfigFile()
		}
	}

	var k8sClient *k8s.Client

	if inCluster {
		c, err := k8s.NewClientInCluster()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if namespace == "" {
			if allNamespaces {
				namespace = k8s.AllNamespaces()
			} else {
				namespace = k8s.DefaultNamespace()
			}
		}

		k8sClient = c
	} else {
		c, err := k8s.NewClient(kubeconfig, kubeContext)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if namespace == "" {
			namespaceInConfig, err := c.NamespaceInConfig()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			if namespaceInConfig == "" {
				if allNamespaces {
					namespace = k8s.AllNamespaces()
				} else {
					namespace = k8s.DefaultNamespace()
				}
			} else {
				namespace = namespaceInConfig
			}
		}

		k8sClient = c
	}

	pods, err := k8sClient.ListPods(namespace, labels)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if verbose {
		for _, pod := range pods {
			pp.Println(pod)
		}
	} else {
		for _, pod := range pods {
			fmt.Println(pod.Name)
		}
	}
}
