package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// RayClusterSpec defines the desired state of RayCluster
type RayClusterSpec struct {
	// EnableInTreeAutoscaling indicates whether operator should create in tree autoscaling configs
	EnableInTreeAutoscaling *bool `json:"enableInTreeAutoscaling,omitempty"`
	// AutoscalerOptions specifies optional configuration for the Ray autoscaler.
	AutoscalerOptions *AutoscalerOptions `json:"autoscalerOptions,omitempty"`
	// Suspend indicates whether a RayCluster should be suspended.
	// A suspended RayCluster will have head pods and worker pods deleted.
	Suspend                *bool             `json:"suspend,omitempty"`
	HeadServiceAnnotations map[string]string `json:"headServiceAnnotations,omitempty"`
	// HeadGroupSpec is the spec for the head pod
	HeadGroupSpec HeadGroupSpec `json:"headGroupSpec"`
	// RayVersion is used to determine the command for the Kubernetes Job managed by RayJob
	RayVersion string `json:"rayVersion,omitempty"`
	// WorkerGroupSpecs are the specs for the worker pods
	WorkerGroupSpecs []WorkerGroupSpec `json:"workerGroupSpecs,omitempty"`
}

// HeadGroupSpec are the spec for the head pod
type HeadGroupSpec struct {
	// Template is the exact pod template used in K8s deployments, statefulsets, etc.
	Template corev1.PodTemplateSpec `json:"template"`
	// HeadService is the Kubernetes service of the head pod.
	HeadService *corev1.Service `json:"headService,omitempty"`
	// EnableIngress indicates whether operator should create ingress object for head service or not.
	EnableIngress *bool `json:"enableIngress,omitempty"`
	// RayStartParams are the params of the start command: node-manager-port, object-store-memory, ...
	RayStartParams map[string]string `json:"rayStartParams"`
	// ServiceType is Kubernetes service type of the head service. it will be used by the workers to connect to the head pod
	ServiceType corev1.ServiceType `json:"serviceType,omitempty"`
}

// WorkerGroupSpec are the specs for the worker pods
type WorkerGroupSpec struct {
	// we can have multiple worker groups, we distinguish them by name
	GroupName string `json:"groupName"`
	// Replicas is the number of desired Pods for this worker group. See https://github.com/ray-project/kuberay/pull/1443 for more details about the reason for making this field optional.
	// +kubebuilder:default:=0
	Replicas *int32 `json:"replicas,omitempty"`
	// MinReplicas denotes the minimum number of desired Pods for this worker group.
	// +kubebuilder:default:=0
	MinReplicas *int32 `json:"minReplicas"`
	// MaxReplicas denotes the maximum number of desired Pods for this worker group, and the default value is maxInt32.
	// +kubebuilder:default:=2147483647
	MaxReplicas *int32 `json:"maxReplicas"`
	// RayStartParams are the params of the start command: address, object-store-memory, ...
	RayStartParams map[string]string `json:"rayStartParams"`
	// Template is a pod template for the worker
	Template corev1.PodTemplateSpec `json:"template"`
	// ScaleStrategy defines which pods to remove
	ScaleStrategy ScaleStrategy `json:"scaleStrategy,omitempty"`
}

// ScaleStrategy to remove workers
type ScaleStrategy struct {
	// WorkersToDelete workers to be deleted
	WorkersToDelete []string `json:"workersToDelete,omitempty"`
}

// AutoscalerOptions specifies optional configuration for the Ray autoscaler.
type AutoscalerOptions struct {
	// Resources specifies optional resource request and limit overrides for the autoscaler container.
	// Default values: 500m CPU request and limit. 512Mi memory request and limit.
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`
	// Image optionally overrides the autoscaler's container image. This override is provided for autoscaler testing and development.
	Image *string `json:"image,omitempty"`
	// ImagePullPolicy optionally overrides the autoscaler container's image pull policy. This override is provided for autoscaler testing and development.
	ImagePullPolicy *corev1.PullPolicy `json:"imagePullPolicy,omitempty"`
	// SecurityContext defines the security options the container should be run with.
	// If set, the fields of SecurityContext override the equivalent fields of PodSecurityContext.
	// More info: https://kubernetes.io/docs/tasks/configure-pod-container/security-context/
	SecurityContext *corev1.SecurityContext `json:"securityContext,omitempty"`
	// IdleTimeoutSeconds is the number of seconds to wait before scaling down a worker pod which is not using Ray resources.
	// Defaults to 60 (one minute). It is not read by the KubeRay operator but by the Ray autoscaler.
	IdleTimeoutSeconds *int32 `json:"idleTimeoutSeconds,omitempty"`
	// UpscalingMode is "Conservative", "Default", or "Aggressive."
	// Conservative: Upscaling is rate-limited; the number of pending worker pods is at most the size of the Ray cluster.
	// Default: Upscaling is not rate-limited.
	// Aggressive: An alias for Default; upscaling is not rate-limited.
	// It is not read by the KubeRay operator but by the Ray autoscaler.
	UpscalingMode *UpscalingMode `json:"upscalingMode,omitempty"`
	// Optional list of environment variables to set in the autoscaler container.
	Env []corev1.EnvVar `json:"env,omitempty"`
	// Optional list of sources to populate environment variables in the autoscaler container.
	EnvFrom []corev1.EnvFromSource `json:"envFrom,omitempty"`
	// Optional list of volumeMounts.  This is needed for enabling TLS for the autoscaler container.
	VolumeMounts []corev1.VolumeMount `json:"volumeMounts,omitempty"`
}

// +kubebuilder:validation:Enum=Default;Aggressive;Conservative
type UpscalingMode string

// The overall state of the Ray cluster.
type ClusterState string

const (
	Ready     ClusterState = "ready"
	Unhealthy ClusterState = "unhealthy"
	Failed    ClusterState = "failed"
	Suspended ClusterState = "suspended"
)

// RayClusterStatus defines the observed state of RayCluster
type RayClusterStatus struct {
	// LastUpdateTime indicates last update timestamp for this cluster status.
	// +nullable
	LastUpdateTime *metav1.Time `json:"lastUpdateTime,omitempty"`
	// Service Endpoints
	Endpoints map[string]string `json:"endpoints,omitempty"`
	// DesiredCPU indicates total desired CPUs for the cluster
	DesiredCPU resource.Quantity `json:"desiredCPU,omitempty"`
	// DesiredMemory indicates total desired memory for the cluster
	DesiredMemory resource.Quantity `json:"desiredMemory,omitempty"`
	// DesiredGPU indicates total desired GPUs for the cluster
	DesiredGPU resource.Quantity `json:"desiredGPU,omitempty"`
	// DesiredTPU indicates total desired TPUs for the cluster
	DesiredTPU resource.Quantity `json:"desiredTPU,omitempty"`
	// Head info
	Head HeadInfo `json:"head,omitempty"`
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// Status reflects the status of the cluster
	//
	// Deprecated: the State field is replaced by the Conditions field.
	State ClusterState `json:"state,omitempty"`
	// Reason provides more information about current State
	Reason string `json:"reason,omitempty"`
	// ReadyWorkerReplicas indicates how many worker replicas are ready in the cluster
	ReadyWorkerReplicas int32 `json:"readyWorkerReplicas,omitempty"`
	// AvailableWorkerReplicas indicates how many replicas are available in the cluster
	AvailableWorkerReplicas int32 `json:"availableWorkerReplicas,omitempty"`
	// DesiredWorkerReplicas indicates overall desired replicas claimed by the user at the cluster level.
	DesiredWorkerReplicas int32 `json:"desiredWorkerReplicas,omitempty"`
	// MinWorkerReplicas indicates sum of minimum replicas of each node group.
	MinWorkerReplicas int32 `json:"minWorkerReplicas,omitempty"`
	// MaxWorkerReplicas indicates sum of maximum replicas of each node group.
	MaxWorkerReplicas int32 `json:"maxWorkerReplicas,omitempty"`
	// observedGeneration is the most recent generation observed for this RayCluster. It corresponds to the
	// RayCluster's generation, which is updated on mutation by the API Server.
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}

// HeadInfo gives info about head
type HeadInfo struct {
	PodIP     string `json:"podIP,omitempty"`
	ServiceIP string `json:"serviceIP,omitempty"`
}

// RayNodeType  the type of a ray node: head/worker
type RayNodeType string

const (
	HeadNode   RayNodeType = "head"
	WorkerNode RayNodeType = "worker"
	// RedisCleanupNode is a Pod managed by a Kubernetes Job that cleans up Redis data after
	// a RayCluster with GCS fault tolerance enabled is deleted.
	RedisCleanupNode RayNodeType = "redis-cleanup"
)

// RayCluster is the Schema for the RayClusters API
// +kubebuilder:object:root=true
// +kubebuilder:resource:categories=all
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="desired workers",type=integer,JSONPath=".status.desiredWorkerReplicas",priority=0
// +kubebuilder:printcolumn:name="available workers",type=integer,JSONPath=".status.availableWorkerReplicas",priority=0
// +kubebuilder:printcolumn:name="cpus",type=string,JSONPath=".status.desiredCPU",priority=0
// +kubebuilder:printcolumn:name="memory",type=string,JSONPath=".status.desiredMemory",priority=0
// +kubebuilder:printcolumn:name="gpus",type=string,JSONPath=".status.desiredGPU",priority=0
// +kubebuilder:printcolumn:name="tpus",type=string,JSONPath=".status.desiredTPU",priority=1
// +kubebuilder:printcolumn:name="status",type="string",JSONPath=".status.state",priority=0
// +kubebuilder:printcolumn:name="age",type="date",JSONPath=".metadata.creationTimestamp",priority=0
// +kubebuilder:printcolumn:name="head pod IP",type="string",JSONPath=".status.head.podIP",priority=1
// +kubebuilder:printcolumn:name="head service IP",type="string",JSONPath=".status.head.serviceIP",priority=1
// +genclient
type RayCluster struct {
	// Standard object metadata.
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Specification of the desired behavior of the RayCluster.
	Spec   RayClusterSpec   `json:"spec,omitempty"`
	Status RayClusterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RayClusterList contains a list of RayCluster
type RayClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RayCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RayCluster{}, &RayClusterList{})
}

type EventReason string

const (
	RayConfigError         EventReason = "RayConfigError"
	PodReconciliationError EventReason = "PodReconciliationError"
)
