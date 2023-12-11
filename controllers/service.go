package controllers

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"

	rhtasv1alpha1 "github.com/securesign/operator/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// ensureServiceCluster Returns the existing IPFS cluster service object or an error.
func (r *SecuresignReconciler) ensureServiceCluster(
	ctx context.Context,
	m *rhtasv1alpha1.Securesign,
) (*corev1.Service, error) {
	logger := log.FromContext(ctx)
	svcName := "todo-" + m.Name
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      svcName,
			Namespace: m.Namespace,
			// TODO: annotations for external dns
		},
	}

	logger.Info("creating or updating svc")
	op, err := ctrl.CreateOrUpdate(ctx, r.Client, svc, func() error {
		svc.Spec = corev1.ServiceSpec{}
		svc.Spec.Ports = []corev1.ServicePort{
			{
				Name:       "swarm",
				Protocol:   corev1.ProtocolTCP,
				Port:       8080,
				TargetPort: intstr.FromString("swarm"),
			},
		}
		svc.Spec.Selector = map[string]string{
			"app.kubernetes.io/name": "rhats-" + m.Name,
		}
		if err := ctrl.SetControllerReference(m, svc, r.Scheme); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		logger.Error(err, "failed on operation", "operation", op)
		return nil, fmt.Errorf("failed to create service: %w", err)
	}
	logger.Info("completed operation", "operation", op)
	return svc, nil
}
