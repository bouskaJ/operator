package utils

import (
	"fmt"

	"github.com/operator-framework/operator-lib/proxy"
	"github.com/securesign/operator/api/v1alpha1"
	"github.com/securesign/operator/internal/controller/constants"
	v1 "k8s.io/api/batch/v1"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

func CreateTufInitJob(instance *v1alpha1.Tuf, name string, sa string, labels map[string]string) *v1.Job {
	isName := "repository-" + instance.Name
	env := []core.EnvVar{
		{
			Name:  "NAMESPACE",
			Value: instance.Namespace,
		},
	}
	env = append(env, proxy.ReadProxyVarsFromEnv()...)
	job := &v1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: instance.Namespace,
			Labels:    labels,
		},
		Spec: v1.JobSpec{
			Parallelism: ptr.To[int32](1),
			Completions: ptr.To[int32](1),

			Template: core.PodTemplateSpec{
				Spec: core.PodSpec{
					ServiceAccountName: sa,
					RestartPolicy:      core.RestartPolicyNever,
					Volumes: []core.Volume{
						{
							Name: "tuf-secrets",
							VolumeSource: core.VolumeSource{
								Projected: secretsVolumeProjection(instance.Status.Keys),
							},
						},
						{
							Name: "repository",
							VolumeSource: core.VolumeSource{
								EmptyDir: &core.EmptyDirVolumeSource{},
							},
						},
						{
							Name: "dockerfile",
							VolumeSource: core.VolumeSource{
								EmptyDir: &core.EmptyDirVolumeSource{},
							},
						},
					},

					InitContainers: []core.Container{
						{
							Name:  "tuf-init",
							Image: constants.TufImage,
							Env:   env,
							Args: []string{
								"-mode", "init-no-overwrite",
								"-target-dir", "/var/run/target",
							},
							VolumeMounts: []core.VolumeMount{
								{
									Name:      "tuf-secrets",
									MountPath: "/var/run/tuf-secrets",
								},
								{
									Name:      "repository",
									MountPath: "/var/run/target",
								},
							},
						},
						{
							Name:    "dockerfile-init",
							Image:   constants.OseToolsImage,
							Env:     env,
							Command: []string{"/bin/bash", "-c"},
							Args: []string{
								"echo 'FROM registry.access.redhat.com/ubi9/ubi-minimal:9.4 \n COPY . /var/run/repository' > /var/run/dockerfile/Dockerfile ",
							},
							VolumeMounts: []core.VolumeMount{
								{
									Name:      "dockerfile",
									MountPath: "/var/run/dockerfile",
								},
							},
						},
					},
					Containers: []core.Container{
						{
							Name:    "init-build",
							Image:   constants.OseToolsImage,
							Env:     env,
							Command: []string{"/bin/bash", "-c"},
							Args: []string{
								fmt.Sprintf("oc new-build --strategy=docker --binary --name=%s --to=%s:1 && "+
									"oc set image-lookup %s &&"+
									"oc start-build %s --from-dir=/var/run/target --from-dir=/var/run/dockerfile --follow", isName, isName, isName, isName),
							},
							VolumeMounts: []core.VolumeMount{
								{
									Name:      "repository",
									MountPath: "/var/run/target",
								},
								{
									Name:      "dockerfile",
									MountPath: "/var/run/dockerfile",
								},
							},
						}},
				},
			},
		},
	}
	return job
}
