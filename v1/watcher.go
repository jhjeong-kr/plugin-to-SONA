package v1

import (
	"encoding/json"
	"errors"
	"os/user"
	"time"

	config "plugin-to-SONA/config"
	log "plugin-to-SONA/log"
	util "plugin-to-SONA/util"

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

	registerPodWatcher(clientset, stop)
	registerServiceWatcher(clientset, stop)

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

func toPod(obj interface{}) *v1.Pod {
	var pod v1.Pod

	json.Unmarshal([]byte(util.InterfaceToString(obj)), &pod)
	return &pod
}

func toService(obj interface{}) *v1.Service {
	var service v1.Service

	json.Unmarshal([]byte(util.InterfaceToString(obj)), &service)
	return &service
}

func podAddFunc(obj interface{}) {
	pod := toPod(obj)
	log.Infof("Pod(\"%s\") is added", pod.Name)
	//	log.Info(util.InterfaceToIndenttedString(obj))
}

func podDeleteFunc(obj interface{}) {
	pod := toPod(obj)
	log.Infof("Pod(\"%s\") is deleted", pod.Name)
}

func podUpdateFunc(oldObj, newObj interface{}) {
	oldPod := toPod(oldObj)
	log.Infof("Pod(\"%s\") is updated", oldPod.Name)
}

func registerPodWatcher(clientset *kubernetes.Clientset, stop chan struct{}) {
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

func serviceAddFunc(obj interface{}) {
	service := toService(obj)
	log.Infof("Service(\"%s\") is added", service.Name)
}

func serviceDeleteFunc(obj interface{}) {
	service := toService(obj)
	log.Infof("Service(\"%s\") is deleted", service.Name)
}

func serviceUpdateFunc(oldObj, newObj interface{}) {
	oldService := toService(oldObj)
	log.Infof("Service(\"%s\") is updated", oldService.Name)
}

func registerServiceWatcher(clientset *kubernetes.Clientset, stop chan struct{}) {
	watchlist := cache.NewListWatchFromClient(clientset.Core().RESTClient(), "services", v1.NamespaceDefault, fields.Everything())
	_, controller := cache.NewInformer(
		watchlist,
		&v1.Service{},
		time.Second*0,
		cache.ResourceEventHandlerFuncs{
			AddFunc:    serviceAddFunc,
			DeleteFunc: serviceDeleteFunc,
			UpdateFunc: serviceUpdateFunc,
		},
	)
	go controller.Run(stop)
}
