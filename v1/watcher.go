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

// Run registers CRUD watchers for Pod on kubernetes client.
func Run() int {
	clientset, err := initKubeClient()
	if err != nil {
		log.Error(err.Error())
		return config.EXITKUBEINIT
	}

	stop := make(chan struct{})

	registerPodWatcher(clientset, stop)
	registerServiceWatcher(clientset, stop)
	//	registerNodeWatcher(clientset, stop)

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

func addFunc(obj interface{}) {
	log.Infof("event \"ADD\" on resource \"%T\" is fired", obj)
	switch obj.(type) {
	case *v1.Pod:
		event := NewPodAsyncEvent(AddEvent, obj)
		GetAsyncHandler().Run(event)
	default:
		log.Info("\tnot implemented object")
	}
}

func deleteFunc(obj interface{}) {
	log.Infof("event \"DELETE\" on resource \"%T\" is fired", obj)
	switch obj.(type) {
	case *v1.Pod:
		event := NewPodAsyncEvent(DeleteEvent, obj)
		log.Info(event.String())
	default:
		log.Info("\tnot implemented object")
	}
}

func updateFunc(oldObj, newObj interface{}) {
	log.Infof("event \"UPDATE\" on resource \"%T\" is fired", oldObj)
	switch oldObj.(type) {
	case *v1.Pod:
		event := NewPodAsyncEvent(UpdateEvent, oldObj, newObj)
		log.Info(event.String())
	default:
		log.Info("\tnot implemented object")
	}
}

func registerPodWatcher(clientset *kubernetes.Clientset, stop chan struct{}) {
	watchlist := cache.NewListWatchFromClient(clientset.Core().RESTClient(), string(v1.ResourcePods), v1.NamespaceAll, fields.Everything())
	_, controller := cache.NewInformer(
		watchlist,
		&v1.Pod{},
		time.Second*0,
		cache.ResourceEventHandlerFuncs{
			AddFunc:    addFunc,
			DeleteFunc: deleteFunc,
			UpdateFunc: updateFunc,
		},
	)
	go controller.Run(stop)
}

func registerServiceWatcher(clientset *kubernetes.Clientset, stop chan struct{}) {
	watchlist := cache.NewListWatchFromClient(clientset.Core().RESTClient(), string(v1.ResourceServices), v1.NamespaceAll, fields.Everything())
	_, controller := cache.NewInformer(
		watchlist,
		&v1.Service{},
		time.Second*0,
		cache.ResourceEventHandlerFuncs{
			AddFunc:    addFunc,
			DeleteFunc: deleteFunc,
			UpdateFunc: updateFunc,
		},
	)
	go controller.Run(stop)
}

func registerNodeWatcher(clientset *kubernetes.Clientset, stop chan struct{}) {
	watchlist := cache.NewListWatchFromClient(clientset.Core().RESTClient(), "nodes", v1.NamespaceAll, fields.Everything())
	_, controller := cache.NewInformer(
		watchlist,
		&v1.Node{},
		time.Second*0,
		cache.ResourceEventHandlerFuncs{
			AddFunc:    addFunc,
			DeleteFunc: deleteFunc,
			UpdateFunc: updateFunc,
		},
	)
	go controller.Run(stop)
}
