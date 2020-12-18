package controllers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// EndpointsReconciler reconciles a Endpoints object
type EndpointsReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=core,resources=endpoints,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=endpoints/status,verbs=get;update;patch

func (r *EndpointsReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	reqLogger := r.Log.WithValues("endpoints", req.NamespacedName)

	// your logic here
	if req.Namespace != "default" {
		// skip endpoints outside default namespace
		return ctrl.Result{}, nil
	}

	instance := &corev1.Endpoints{}
	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// log the service dns and the endpoint IPs with ports
	reqLogger.Info("Endpoint watcher triggered", "DNS:", fmt.Sprintf("%s.%s.svc.cluster.local", req.Name, req.Namespace))
	for _, subset := range instance.Subsets {
		reqLogger.Info("IP Addresses:", "Ready IP addresses", fmt.Sprintf("%#v\n", subset.Addresses))
		reqLogger.Info("IP Addresses:", "Not Ready IP addresses", fmt.Sprintf("%#v\n", subset.NotReadyAddresses))
		reqLogger.Info("Ports:", "Available port numbers", fmt.Sprintf("%#v\n", subset.Ports))
	}

	return ctrl.Result{}, nil
}

func (r *EndpointsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Endpoints{}).
		Complete(r)
}
