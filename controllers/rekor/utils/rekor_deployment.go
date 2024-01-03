package utils

import (
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateRekorDeployment(namespace string, dpName string, pvc string) *apps.Deployment {
	replicas := int32(1)
	return &apps.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      dpName,
			Namespace: namespace,
			Labels: map[string]string{
				"app.kubernetes.io/component": dpName,
				"app.kubernetes.io/name":      dpName,
				"app.kubernetes.io/instance":  "trusted-artifact-signer",
			},
		},
		Spec: apps.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app.kubernetes.io/component": dpName,
					"app.kubernetes.io/name":      dpName,
					"app.kubernetes.io/instance":  "trusted-artifact-signer",
				},
			},
			Template: core.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app.kubernetes.io/component": dpName,
						"app.kubernetes.io/name":      dpName,
						"app.kubernetes.io/instance":  "trusted-artifact-signer",
					},
				},
				Spec: core.PodSpec{
					ServiceAccountName: "sigstore-sa",
					Volumes: []core.Volume{
						{
							Name: "rekor-sharding-config",
							VolumeSource: core.VolumeSource{
								ConfigMap: &core.ConfigMapVolumeSource{
									LocalObjectReference: core.LocalObjectReference{
										Name: "rekor-sharding-config",
									},
								},
							},
						},
						{
							Name: "storage",
							VolumeSource: core.VolumeSource{
								PersistentVolumeClaim: &core.PersistentVolumeClaimVolumeSource{
									ClaimName: pvc,
								},
							},
						},
						{
							Name: "rekor-private-key-volume",
							VolumeSource: core.VolumeSource{
								Secret: &core.SecretVolumeSource{
									SecretName: "rekor-private-key",
									Items: []core.KeyToPath{
										{
											Key:  "private",
											Path: "private",
										},
									},
								},
							},
						},
					},
					Containers: []core.Container{
						{
							Name: dpName,
							// TODO add probe
							//LivenessProbe: &core.Probe{
							//	},
							//	InitialDelaySeconds: 30,
							//	TimeoutSeconds:      1,
							//	PeriodSeconds:       10,
							//	SuccessThreshold:    1,
							//	FailureThreshold:    3,
							//},
							Env: []core.EnvVar{
								{
									Name: "TREE_ID",
									ValueFrom: &core.EnvVarSource{
										ConfigMapKeyRef: &core.ConfigMapKeySelector{
											LocalObjectReference: core.LocalObjectReference{
												Name: "rekor-config",
											},
											Key: "treeID",
										},
									},
								},
							},
							Image: "registry.redhat.io/rhtas-tech-preview/rekor-server-rhel9@sha256:8ee7d5dd2fa1c955d64ab83d716d482a3feda8e029b861241b5b5dfc6f1b258e",
							Ports: []core.ContainerPort{
								{
									ContainerPort: 3000,
									Name:          "rekor-server",
								},
								{
									ContainerPort: 2112,
									Protocol:      "TCP",
								},
							},
							Args: []string{
								"serve",
								"--trillian_log_server.address=trillian-logserver." + namespace,
								"--trillian_log_server.port=8091",
								"--trillian_log_server.sharding_config=/sharding/sharding-config.yaml",
								"--redis_server.address=rekor-redis",
								"--redis_server.port=6379",
								"--rekor_server.address=0.0.0.0",
								"--rekor_server.signer=/key/private",
								"--enable_retrieve_api=true",
								"--trillian_log_server.tlog_id=$(TREE_ID)",
								"--enable_attestation_storage",
								"--attestation_storage_bucket=file:///var/run/attestations",
							},
							VolumeMounts: []core.VolumeMount{
								{
									Name:      "rekor-sharding-config",
									MountPath: "/sharding",
								},
								{
									Name:      "storage",
									MountPath: "/var/run/attestations",
								},
								{
									Name:      "rekor-private-key-volume",
									MountPath: "/key",
									ReadOnly:  true,
								},
							},
						},
					},
				},
			},
		},
	}
}
