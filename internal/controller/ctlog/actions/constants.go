package actions

import "github.com/securesign/operator/internal/controller/constants"

const (
	DeploymentName     = "ctlog"
	ComponentName      = "ctlog"
	RBACName           = "ctlog"
	MonitoringRoleName = "prometheus-k8s-ctlog"

	CertCondition    = "FulcioCertAvailable"
	ServerPortName   = "http"
	ServerPort       = 80
	ServerTargetPort = 6962
	MetricsPortName  = "metrics"
	MetricsPort      = 6963
	ServerCondition  = "ServerAvailable"

	CTLPubLabel = constants.LabelNamespace + "/ctfe.pub"
)
