package trillian

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	v1 "github.com/openshift/api/route/v1"
	rhtasv1alpha1 "github.com/securesign/operator/api/v1alpha1"
	"github.com/securesign/operator/controllers/common"
	"github.com/securesign/operator/controllers/common/utils/kubernetes"
	"github.com/securesign/operator/controllers/constants"
	trillianUtils "github.com/securesign/operator/controllers/trillian/utils"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	dbDeploymentName        = "trillian-db"
	logserverDeploymentName = "trillian-logserver"
	logsignerDeploymentName = "trillian-logsigner"

	ComponentName = "trillian"
)

func NewCreateAction() Action {
	return &createAction{}
}

type createAction struct {
	common.BaseAction
}

func (i createAction) Name() string {
	return "create"
}

func (i createAction) CanHandle(trillian *rhtasv1alpha1.Trillian) bool {
	return trillian.Status.Phase == rhtasv1alpha1.PhaseNone
}

func (i createAction) Handle(ctx context.Context, instance *rhtasv1alpha1.Trillian) (*rhtasv1alpha1.Trillian, error) {
	//log := ctrllog.FromContext(ctx)
	var err error

	if instance.Spec.TrillianDB.Create && instance.Spec.TrillianDB.Url != "" {
		err = common.InvalidConfigurationError("db.create and db.url configuration can't be used at once")
		i.Logger.Error(err, "Can't use db.create and db.url at once")
		instance.Status.Phase = rhtasv1alpha1.PhaseError
		return instance, err
	}

	if instance.Spec.TrillianDB.Url != "" && instance.Spec.TrillianDB.ConnectionSecret == "" {
		err = common.InvalidConfigurationError("db.connectionSecret is not set")
		i.Logger.Error(err, "You need to configure connection secret together with provided DB")
		instance.Status.Phase = rhtasv1alpha1.PhaseError
		return instance, err
	}

	logSignerLabels := kubernetes.FilterCommonLabels(instance.Labels)
	logSignerLabels["app.kubernetes.io/component"] = ComponentName
	logSignerLabels["app.kubernetes.io/name"] = logsignerDeploymentName

	logServerLabels := kubernetes.FilterCommonLabels(instance.Labels)
	logServerLabels["app.kubernetes.io/component"] = ComponentName
	logServerLabels["app.kubernetes.io/name"] = logserverDeploymentName

	if instance.Spec.TrillianDB.Create {
		dbLabels := kubernetes.FilterCommonLabels(instance.Labels)
		dbLabels["app.kubernetes.io/component"] = ComponentName
		dbLabels["app.kubernetes.io/name"] = dbDeploymentName

		if instance.Spec.TrillianDB.ConnectionSecret != "" {
			i.createDbSecret(ctx, instance, dbLabels)
		}
		dbPVC := instance.Spec.TrillianDB.PvcName
		if dbPVC == "" {
			i.createDbPVC(ctx, instance, dbLabels)
		}

	}

	db := trillianUtils.CreateTrillDb(instance.Namespace, constants.TrillianDbImage, dbDeploymentName, trillPVC, dbSecret.Name, dbLabels)
	controllerutil.SetControllerReference(instance, db, i.Client.Scheme())
	if err = i.Client.Create(ctx, db); err != nil {
		instance.Status.Phase = rhtasv1alpha1.PhaseError
		return instance, fmt.Errorf("could not create trillian DB: %w", err)
	}

	mysqlPort := 3306
	mysql := kubernetes.CreateService(instance.Namespace, "trillian-mysql", mysqlPort, dbLabels)
	controllerutil.SetControllerReference(instance, mysql, i.Client.Scheme())
	if err = i.Client.Create(ctx, mysql); err != nil {
		instance.Status.Phase = rhtasv1alpha1.PhaseError
		return instance, fmt.Errorf("could not create service: %w", err)
	}
	instance.Status.DBUrl = fmt.Sprintf("%s.%s.svc:%d", mysql.Name, mysql.Namespace, mysqlPort)

	// Log Server
	svcName := "trillian-logserver"
	tlsSecretName := ""
	serverPort := 8091

	logserverService := kubernetes.CreateService(instance.Namespace, svcName, serverPort, logServerLabels)

	if instance.Spec.External {
		tlsSecretName = svcName + "-" + uuid.New().String()
		// generate signed certificate
		logserverService.Annotations = map[string]string{"service.beta.openshift.io/serving-cert-secret-name": tlsSecretName}

		route := kubernetes.CreateRoute(*logserverService, svcName, logServerLabels)
		route.Spec.TLS.Termination = v1.TLSTerminationReencrypt
		controllerutil.SetControllerReference(instance, route, i.Client.Scheme())
		if err = i.Client.Create(ctx, route); err != nil {
			instance.Status.Phase = rhtasv1alpha1.PhaseError
			return instance, fmt.Errorf("could not create route: %w", err)
		}
		instance.Status.Url = route.Spec.Host + ":443"
	} else {
		instance.Status.Url = fmt.Sprintf("%s.%s.svc:%d", logserverService.Name, logserverService.Namespace, serverPort)
	}
	controllerutil.SetControllerReference(instance, logserverService, i.Client.Scheme())
	if err = i.Client.Create(ctx, logserverService); err != nil {
		instance.Status.Phase = rhtasv1alpha1.PhaseError
		return instance, fmt.Errorf("could not create service: %w", err)
	}

	server := trillianUtils.CreateTrillDeployment(instance.Namespace, constants.TrillianServerImage, logserverDeploymentName, dbSecret.Name, tlsSecretName, logServerLabels)
	controllerutil.SetControllerReference(instance, server, i.Client.Scheme())
	server.Spec.Template.Spec.Containers[0].Ports = append(server.Spec.Template.Spec.Containers[0].Ports, corev1.ContainerPort{
		Protocol:      corev1.ProtocolTCP,
		ContainerPort: 8090,
	})
	if err = i.Client.Create(ctx, server); err != nil {
		instance.Status.Phase = rhtasv1alpha1.PhaseError
		return instance, fmt.Errorf("could not create job: %w", err)
	}

	// Log Signer
	signer := trillianUtils.CreateTrillDeployment(instance.Namespace, constants.TrillianLogSignerImage, logsignerDeploymentName, dbSecret.Name, "", logSignerLabels)
	controllerutil.SetControllerReference(instance, signer, i.Client.Scheme())
	signer.Spec.Template.Spec.Containers[0].Args = append(signer.Spec.Template.Spec.Containers[0].Args, "--force_master=true")
	if err = i.Client.Create(ctx, signer); err != nil {
		instance.Status.Phase = rhtasv1alpha1.PhaseError
		return instance, fmt.Errorf("could not create job: %w", err)
	}

	instance.Status.Phase = rhtasv1alpha1.PhaseCreating
	return instance, nil

}

func (i createAction) createDbSecret(ctx context.Context, instance *rhtasv1alpha1.Trillian, labels map[string]string) (*corev1.Secret, error) {
	secret := kubernetes.CreateSecret(dbDeploymentName, instance.Namespace, map[string][]byte{
		// generate a random password for the mysql root user and the mysql password
		// TODO - use a random password generator
		"mysql-root-password": []byte("password"),
		"mysql-password":      []byte("password"),
		"mysql-database":      []byte("trillian"),
		"mysql-user":          []byte("mysql"),
	}, labels)
	controllerutil.SetControllerReference(instance, secret, i.Client.Scheme())
	err := i.Client.Create(ctx, secret)
	return secret, err
}

func (i createAction) createDbPVC(ctx context.Context, instance *rhtasv1alpha1.Trillian, labels map[string]string) (*corev1.PersistentVolumeClaim, error) {
	pvc := kubernetes.CreatePVC(instance.Namespace, "trillian-mysql", "5Gi", labels)
	controllerutil.SetControllerReference(instance, pvc, i.Client.Scheme())
	if err := i.Client.Create(ctx, pvc); err != nil {
		instance.Status.Phase = rhtasv1alpha1.PhaseError
		return nil, fmt.Errorf("could not create pvc: %w", err)
	}
	return pvc, nil
}
