package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"

	"k8s.io/client-go/util/workqueue"
)

type handler struct {
	q workqueue.RateLimitingInterface
}

func (h *handler) OnAdd(obj interface{}) {
	fmt.Println("add")
}
func (h *handler) OnUpdate(o, n interface{}) {
	fmt.Println("update")
}
func (h *handler) OnDelete(obj interface{}) {
	fmt.Println("delete")
}

func newHandler() *handler {
	return &handler{
		q: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "test"),
	}
}

func main() {
	url := "https://192.168.0.4:6443"
	fmt.Println(url)
	cfg, err := clientcmd.BuildConfigFromFlags(url, "/home/trainyao/.kube/config")
	if err != nil {
		panic(err)
	}

	client, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		panic(err)
	}

	i := informers.NewSharedInformerFactory(client, 30*time.Second)
	informer := i.Core().V1().Pods().Informer()
	informer.AddEventHandler(newHandler())

	stopChan := make(chan struct{})
	sigChan := make(chan os.Signal)

	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("notified")

		stopChan <- struct{}{}
		close(stopChan)
	}()

	i.Start(stopChan)

	if ok := cache.WaitForCacheSync(stopChan, informer.HasSynced); !ok {
		fmt.Println("fail to sync")

		return
	}
	fmt.Println("synced")

	wait.Until(func() {}, 1*time.Second, stopChan)
}
