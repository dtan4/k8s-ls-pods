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
		failed        bool
		kubeContext   string
		kubeconfig    string
		labels        string
		namespace     string
		succeeded     bool
		verbose       bool
	)

	flags := flag.NewFlagSet("k8s-pod-notifier", flag.ExitOnError)
	flags.Usage = func() {
		flags.PrintDefaults()
	}

	flags.BoolVar(&allNamespaces, "all-namespaces", false, "List pods across all namespaces")
	flags.StringVar(&kubeContext, "context", "", "Kubernetes context")
	flags.BoolVar(&failed, "failed", false, "Show failed Pods only")
	flags.StringVar(&kubeconfig, "kubeconfig", "", "Path of kubeconfig")
	flags.StringVarP(&labels, "labels", "l", "", "Label filter query")
	flags.StringVarP(&namespace, "namespace", "n", "", "Kubernetes namespace")
	flags.BoolVar(&succeeded, "succeeded", false, "Show succeeded Pods only")
	flags.BoolVarP(&verbose, "verbose", "v", false, "Verbose output")

	if err := flags.Parse(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	printAll := !failed && !succeeded

	if kubeconfig == "" {
		if os.Getenv("KUBECONFIG") != "" {
			kubeconfig = os.Getenv("KUBECONFIG")
		} else {
			kubeconfig = k8s.DefaultConfigFile()
		}
	}

	k8sClient, err := k8s.NewClient(kubeconfig, kubeContext)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if namespace == "" {
		namespaceInConfig, err := k8sClient.NamespaceInConfig()
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

	pods, err := k8sClient.ListPods(namespace, labels)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if verbose {
		for _, pod := range pods {
			if !printAll {
				if (failed && !k8s.IsPodFailed(pod)) || (succeeded && !k8s.IsPodSucceeded(pod)) {
					continue
				}
			}

			pp.Println(pod)
			fmt.Println("-----")
		}
	} else {
		for _, pod := range pods {
			if !printAll {
				if (failed && !k8s.IsPodFailed(pod)) || (succeeded && !k8s.IsPodSucceeded(pod)) {
					continue
				}
			}

			fmt.Println(pod.Name)
		}
	}
}
