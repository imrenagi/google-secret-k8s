package agent

import (
	corev1 "k8s.io/api/core/v1"
)

// ContainerInitSidecar creates a new init container to be added
// to the pod being mutated.  After Vault 1.4 is released, this can
// be removed because an exit_after_auth environment variable is
// available for the agent.  This means we won't need to generate
// two config files.
func (a *Agent) ContainerInitSidecar() (corev1.Container, error) {
	volumeMounts := []corev1.VolumeMount{
		{
			Name:      a.ServiceAccountName,
			MountPath: a.ServiceAccountPath,
			ReadOnly:  true,
		},
	}
	volumeMounts = append(volumeMounts, a.ContainerVolumeMounts()...)

	arg := DefaultContainerArg

	if a.ConfigMapName != "" {
		volumeMounts = append(volumeMounts, corev1.VolumeMount{
			Name:      configVolumeName,
			MountPath: configVolumePath,
			ReadOnly:  true,
		})
	}

	envs, err := a.ContainerEnvVars(true)
	if err != nil {
		return corev1.Container{}, err
	}

	resources, err := a.parseResources()
	if err != nil {
		return corev1.Container{}, err
	}

	newContainer := corev1.Container{
		Name:         "vault-agent-init",
		Image:        a.ImageName,
		Env:          envs,
		Resources:    resources,
		VolumeMounts: volumeMounts,
		Command:      []string{"/bin/sh", "-ec"},
		Args:         []string{arg},
	}

	return newContainer, nil
}
