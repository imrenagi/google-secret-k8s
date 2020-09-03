package agent

import (
	corev1 "k8s.io/api/core/v1"
)

// ContainerEnvVars adds the applicable environment vars
// for the Vault Agent sidecar.
func (a *Agent) ContainerEnvVars(init bool) ([]corev1.EnvVar, error) {
	var envs []corev1.EnvVar
	return envs, nil
}
