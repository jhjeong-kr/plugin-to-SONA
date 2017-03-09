package v1

import (
	"errors"
	"os/user"
	"time"

	config "plugin-to-SONA/config"
	log "plugin-to-SONA/log"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/fields"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

// Run registers CRUD watchers for Pod on kubernetes client
func Run() int {
	clientset, err := initKubeClient()
	if err != nil {
		log.Error(err.Error())
		return config.EXITKUBEINIT
	}

	stop := make(chan struct{})

	registerPODWatcher(clientset, stop)

	<-stop

	return config.EXITNORMAL
}

func initKubeClient() (*kubernetes.Clientset, error) {
	var kubeConfig *rest.Config
	var err error

	if len(config.KubeConfig) > 0 {
		if u, _ := user.Current(); u.Gid != "0" {
			return nil, errors.New("please run with root permission")
		}
		kubeConfig, err = clientcmd.BuildConfigFromFlags("", config.KubeConfig)
		if err != nil {
			return nil, err
		}
	} else {
		kubeConfig, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	}
	return kubernetes.NewForConfig(kubeConfig)
}

func podAddFunc(obj interface{}) {
	log.Info("Pod is added")
}

func podDeleteFunc(obj interface{}) {
	log.Info("Pod is added")
}

func podUpdateFunc(oldObj, newObj interface{}) {
	log.Info("Pod is added")
}

func registerPODWatcher(clientset *kubernetes.Clientset, stop chan struct{}) {
	watchlist := cache.NewListWatchFromClient(clientset.Core().RESTClient(), "pods", v1.NamespaceDefault, fields.Everything())
	_, controller := cache.NewInformer(
		watchlist,
		&v1.Pod{},
		time.Second*0,
		cache.ResourceEventHandlerFuncs{
			AddFunc:    podAddFunc,
			DeleteFunc: podDeleteFunc,
			UpdateFunc: podUpdateFunc,
		},
	)
	go controller.Run(stop)
}
