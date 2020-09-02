package agent

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/mattbaird/jsonpatch"
	corev1 "k8s.io/api/core/v1"
)

// TODO swap out 'github.com/mattbaird/jsonpatch' for 'github.com/evanphx/json-patch'

// Agent is the top level structure holding all the
// configurations for the Vault Agent container.
type Agent struct {
	// Annotations are the current pod annotations used to
	// configure the Vault Agent container.
	Annotations map[string]string

	// ImageName is the name of the Vault image to use for the
	// sidecar container.
	ImageName string

	// GoogleSecretEntryName is the name of GoogleSecretEntry
	GoogleSecretEntryName string

	// Inject is the flag used to determine if a container should be requested
	// in a pod request.
	Inject bool

	// InitFirst controls whether an init container is first to run.
	InitFirst bool

	// LimitsCPU is the upper CPU limit the sidecar container is allowed to consume.
	LimitsCPU string

	// LimitsMem is the upper memory limit the sidecar container is allowed to consume.
	LimitsMem string

	// Namespace is the Kubernetes namespace the request originated from.
	Namespace string

	// Patches are all the mutations we will make to the pod request.
	Patches []*jsonpatch.JsonPatchOperation

	// Pod is the original Kubernetes pod spec.
	Pod *corev1.Pod

	// PrePopulate controls whether an init container is added to the request.
	PrePopulate bool

	// PrePopulateOnly controls whether an init container is the _only_ container
	// added to the request.
	PrePopulateOnly bool

	// RequestsCPU is the requested minimum CPU amount required  when being scheduled to deploy.
	RequestsCPU string

	// RequestsMem is the requested minimum memory amount required when being scheduled to deploy.
	RequestsMem string

	// ServiceAccountName is the Kubernetes service account name for the pod.
	// This is used when we mount the service account to the  Vault Agent container(s).
	ServiceAccountName string

	// ServiceAccountPath is the path on disk where the service account JWT
	// can be located.  This is used when we mount the service account to the
	// Vault Agent container(s).
	ServiceAccountPath string

	// Status is the current injection status.  The only status considered is "injected",
	// which prevents further mutations.  A user can patch this annotation to force a new
	// mutation.
	Status string

	// ConfigMapName is the name of the configmap a user wants to mount to Vault Agent
	// container(s).
	ConfigMapName string
}

// New creates a new instance of Agent by parsing all the Kubernetes annotations.
func New(pod *corev1.Pod, patches []*jsonpatch.JsonPatchOperation) (*Agent, error) {
	saName, saPath := serviceaccount(pod)

	agent := &Agent{
		Annotations:        pod.Annotations,
		ConfigMapName:      pod.Annotations[AnnotationAgentConfigMap],
		ImageName:          pod.Annotations[AnnotationAgentImage],
		LimitsCPU:          pod.Annotations[AnnotationAgentLimitsCPU],
		LimitsMem:          pod.Annotations[AnnotationAgentLimitsMem],
		Namespace:          pod.Annotations[AnnotationAgentRequestNamespace],
		Patches:            patches,
		Pod:                pod,
		RequestsCPU:        pod.Annotations[AnnotationAgentRequestsCPU],
		RequestsMem:        pod.Annotations[AnnotationAgentRequestsMem],
		ServiceAccountName: saName,
		ServiceAccountPath: saPath,
		Status:             pod.Annotations[AnnotationAgentStatus],
	}

	var err error
	agent.GoogleSecretEntryName, err = agent.secretEntry()
	if err != nil {
		return agent, err
	}

	agent.Inject, err = agent.inject()
	if err != nil {
		return agent, err
	}

	agent.InitFirst, err = agent.initFirst()
	if err != nil {
		return agent, err
	}

	agent.PrePopulate, err = agent.prePopulate()
	if err != nil {
		return agent, err
	}

	agent.PrePopulateOnly, err = agent.prePopulateOnly()
	if err != nil {
		return agent, err
	}

	return agent, nil
}

// ShouldInject checks whether the pod in question should be injected
// with Vault Agent containers.
func ShouldInject(pod *corev1.Pod) (bool, error) {
	raw, ok := pod.Annotations[AnnotationAgentInject]
	if !ok {
		return false, nil
	}

	inject, err := strconv.ParseBool(raw)
	if err != nil {
		return false, err
	}

	if !inject {
		return false, nil
	}

	// This shouldn't happen so bail.
	raw, ok = pod.Annotations[AnnotationAgentStatus]
	if !ok {
		return true, nil
	}

	// "injected" is the only status we care about.  Don't do
	// anything if it's set.  The user can update the status
	// to force a new mutation.
	if raw == "injected" {
		return false, nil
	}

	return true, nil
}

// Patch creates the necessary pod patches to inject the Google Secret Agent
// containers.
func (a *Agent) Patch() ([]byte, error) {
	var patches []byte

	// TODO add volume for storing google cloud secret from

	// Add our volume that will be shared by the containers
	// for passing data in the pod.
	a.Patches = append(a.Patches, addVolumes(
		a.Pod.Spec.Volumes,
		a.ContainerVolumes(),
		"/spec/volumes")...)

	// Add ConfigMap if one was provided
	if a.ConfigMapName != "" {
		a.Patches = append(a.Patches, addVolumes(
			a.Pod.Spec.Volumes,
			[]corev1.Volume{a.ContainerConfigMapVolume()},
			"/spec/volumes")...)
	}

	//Add Volume Mounts
	for i, container := range a.Pod.Spec.Containers {
		a.Patches = append(a.Patches, addVolumeMounts(
			container.VolumeMounts,
			a.ContainerVolumeMounts(),
			fmt.Sprintf("/spec/containers/%d/volumeMounts", i))...)
	}

	// Init Container
	if a.PrePopulate {
		container, err := a.ContainerInitSidecar()
		if err != nil {
			return patches, err
		}

		containers := a.Pod.Spec.InitContainers

		if a.InitFirst {

			// Remove all init containers from the document so we can re-add them after the agent.
			if len(a.Pod.Spec.InitContainers) != 0 {
				a.Patches = append(a.Patches, removeContainers("/spec/initContainers")...)
			}

			containers = []corev1.Container{container}
			containers = append(containers, a.Pod.Spec.InitContainers...)

			a.Patches = append(a.Patches, addContainers(
				[]corev1.Container{},
				containers,
				"/spec/initContainers")...)
		} else {
			a.Patches = append(a.Patches, addContainers(
				a.Pod.Spec.InitContainers,
				[]corev1.Container{container},
				"/spec/initContainers")...)
		}

		//Add Volume Mounts
		for i, container := range containers {
			if container.Name == "vault-agent-init" {
				continue
			}
			a.Patches = append(a.Patches, addVolumeMounts(
				container.VolumeMounts,
				a.ContainerVolumeMounts(),
				fmt.Sprintf("/spec/initContainers/%d/volumeMounts", i))...)
		}
	}

	// Sidecar Container
	if !a.PrePopulateOnly {
		container, err := a.ContainerSidecar()
		if err != nil {
			return patches, err
		}
		a.Patches = append(a.Patches, addContainers(
			a.Pod.Spec.Containers,
			[]corev1.Container{container},
			"/spec/containers")...)
	}

	// Add annotations so that we know we're injected
	a.Patches = append(a.Patches, updateAnnotations(
		a.Pod.Annotations,
		map[string]string{AnnotationAgentStatus: "injected"})...)

	// Generate the patch
	if len(a.Patches) > 0 {
		var err error
		patches, err = json.Marshal(a.Patches)
		if err != nil {
			return patches, err
		}
	}
	return patches, nil
}

// Validate the instance of Agent to ensure we have everything needed
// for basic functionality.
func (a *Agent) Validate() error {
	if a.Namespace == "" {
		return errors.New("namespace missing from request")
	}

	if a.ServiceAccountName == "" || a.ServiceAccountPath == "" {
		return errors.New("no service account name or path found")
	}

	if a.ImageName == "" {
		return errors.New("no sidecar image found")
	}

	return nil
}

func serviceaccount(pod *corev1.Pod) (string, string) {
	var serviceAccountName, serviceAccountPath string
	for _, container := range pod.Spec.Containers {
		for _, volumes := range container.VolumeMounts {
			if strings.Contains(volumes.MountPath, "serviceaccount") {
				return volumes.Name, volumes.MountPath
			}
		}
	}
	return serviceAccountName, serviceAccountPath
}
