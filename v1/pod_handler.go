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

func addFunc(obj interface{}) {
	log.Infof("event \"ADD\" on \"%T\" resource", obj)
	switch obj.(type) {
	case *v1.Pod:
		pod := obj.(*v1.Pod)
		log.Infof("\tthe pod is \"%s\" in \"%s\" namespace", pod.Name, pod.GetNamespace())
		log.Infof("\tthe pod(\"%s\") is %s on \"%s\" host", pod.Status.PodIP, pod.Status.Phase, pod.Status.HostIP)
		for i, c := range pod.Status.ContainerStatuses {
			log.Infof("\t\t%d: %s(%s)", i, c.Name, c.ContainerID)
		}
	default:
		log.Info("\tnot implemented object")
	}
}

func deleteFunc(obj interface{}) {
	log.Infof("event \"DELETE\" on \"%T\" resource", obj)
	switch obj.(type) {
	case *v1.Pod:
		pod := obj.(*v1.Pod)
		log.Infof("\tthe pod is \"%s\" in \"%s\" namespace", pod.Name, pod.GetNamespace())
		log.Infof("\tthe pod(\"%s\") is %s on \"%s\" host", pod.Status.PodIP, pod.Status.Phase, pod.Status.HostIP)
	default:
		log.Info("\tnot implemented object")
	}
}

func updateFunc(oldObj, newObj interface{}) {
	log.Infof("event \"UPDATE\" on \"%T\" resource", oldObj)
	switch oldObj.(type) {
	case *v1.Pod:
		oldPod := oldObj.(*v1.Pod)
		newPod := newObj.(*v1.Pod)
		log.Infof("\tthe pod is \"%s\" in \"%s\" namespace", oldPod.Name, oldPod.GetNamespace())
		log.Infof("\told one(\"%s\") is %s on \"%s\" host", oldPod.Status.PodIP, oldPod.Status.Phase, oldPod.Status.HostIP)
		log.Infof("\tnew one(\"%s\") is %s on \"%s\" host", newPod.Status.PodIP, newPod.Status.Phase, newPod.Status.HostIP)
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
