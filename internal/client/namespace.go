package client

import (
	"context"
	"flag"
	"log/slog"
	"path/filepath"
	"time"

	"github.com/adalbertjnr/kcmgr/internal/models"
	"github.com/adalbertjnr/kcmgr/internal/util"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

const TimeoutDuration = time.Second * 5

func GetKubeConfigFile() string {
	var kubeconfig string

	if home := homedir.HomeDir(); home != "" {
		path := filepath.Join(home, ".kube", "config")
		flag.StringVar(&kubeconfig, "kubeconfig", path, "kubeconfig file path")
	} else {
		flag.StringVar(&kubeconfig, "kubeconfig", "", "kubeconfig file path")
	}
	flag.Parse()

	return kubeconfig
}

func GetNamespacesByContext(kubeconfig, kubernetesContext string) ([]models.Namespace, error) {
	configLoader := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig},
		&clientcmd.ConfigOverrides{CurrentContext: kubernetesContext},
	)

	config, err := configLoader.ClientConfig()
	if err != nil {
		return nil, err
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), TimeoutDuration)
	defer cancel()

	list, err := clientSet.CoreV1().Namespaces().List(ctx, v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	namespaces := make([]models.Namespace, len(list.Items))
	for i, namespace := range list.Items {
		slog.Info("namespace", "parsing", namespace)
		timeStamp, err := util.ParseTime(namespace.CreationTimestamp.String())
		if err != nil {
			timeStamp = namespace.CreationTimestamp.String()
		}
		namespaces[i] = models.Namespace{
			Name: namespace.Name,
			Age:  timeStamp,
		}
	}

	return namespaces, nil
}
