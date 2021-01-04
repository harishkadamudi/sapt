package mutating

import (
	"context"
	"encoding/json"
	"net/http"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// +kubebuilder:webhook:path=/mutate-v1-pod,mutating=true,failurePolicy=fail,groups="",resources=pods,verbs=create;update,versions=v1,name=mpod.kb.io

type PodSideCarInjector struct {
	Client  client.Client
	decoder *admission.Decoder
}

// Handle yields a response to an AdmissionRequest.
//
// The supplied context is extracted from the received http.Request, allowing wrapping
// http.Handlers to inject values into and control cancelation of downstream request processing.
func (p *PodSideCarInjector) Handle(ctx context.Context, request admission.Request) admission.Response {
	pod := &corev1.Pod{}

	err := p.decoder.Decode(request, pod)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	v, ok := pod.Annotations["sapt-inject"]
	if !ok || v != "enabled" {
		return admission.Allowed("Injection annotation missing")
	}

	// inject envoy sidecar into the pod
	InjectEnvoySidecar(pod)

	marshaledPod, err := json.Marshal(pod)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}

	return admission.PatchResponseFromRaw(request.Object.Raw, marshaledPod)
}

func (p *PodSideCarInjector) InjectDecoder(decoder *admission.Decoder) error {
	p.decoder = decoder
	return nil
}

func InjectEnvoySidecar(pod *corev1.Pod) {
	pod.Spec.Volumes = append(pod.Spec.Volumes, corev1.Volume{
		Name: "startup-config",
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: v1.LocalObjectReference{
					Name: "envoy-startup",
				},
			},
		},
	})

	pod.Spec.Containers = append(pod.Spec.Containers, corev1.Container{
		Name:  "envoy",
		Image: "getenvoy/envoy:stable",
		Args:  []string{"-c", "/etc/envoy/envoy.yaml"},
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      "startup-config",
				MountPath: "/etc/envoy",
			},
		},
	})
}
