package agent

import (
	corev1 "k8s.io/api/core/v1"
)

const (
	// 	tokenVolumeNameInit    = "home-init"
	// 	tokenVolumeNameSidecar = "home-sidecar"
	// 	tokenVolumePath        = "/home/vault"
	configVolumeName = "google-agent-config"
	configVolumePath = "/google/configs"
	secretVolumeName = "google-secrets"
	secretVolumePath = "/google/secrets"
)

// ContainerVolumes returns the volume data to add to the pod. This volumes
// are used for shared data between containers.
func (a *Agent) ContainerVolumes() []corev1.Volume {
	containerVolumes := []corev1.Volume{
		corev1.Volume{
			Name: secretVolumeName,
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{
					Medium: "Memory",
				},
			},
		},
	}
	return containerVolumes
}

// ContainerConfigMapVolume returns a volume to mount a config map
// if the user supplied any.
func (a *Agent) ContainerConfigMapVolume() corev1.Volume {
	return corev1.Volume{
		Name: configVolumeName,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: a.ConfigMapName,
				},
			},
		},
	}
}

// ContainerVolumeMounts mounts the shared memory volume where secrets
// will be rendered.
func (a *Agent) ContainerVolumeMounts() []corev1.VolumeMount {
	volumeMounts := []corev1.VolumeMount{
		corev1.VolumeMount{
			Name:      secretVolumeName,
			MountPath: a.Annotations[AnnotationAgentSecretVolumePath],
			ReadOnly:  false,
		},
	}
	return volumeMounts
}
