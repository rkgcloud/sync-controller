/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"github.com/vmware-labs/reconciler-runtime/apis"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type Image struct {
	// Image is a reference to an image in a remote repository
	// +required
	Image string `json:"image"`

	// SecretRef contains the names of the Kubernetes Secrets containing registry login
	// information to resolve image metadata.
	// +optional
	SecretRef []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`

	// ServiceAccountName is the name of the Kubernetes ServiceAccount used to authenticate
	// the image pull if the service account has attached pull secrets. For more information:
	// https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/#add-imagepullsecrets-to-a-service-account
	// +optional
	ServiceAccountName string `json:"serviceAccountName,omitempty"`

	// Insecure allows connecting to a non-TLS HTTP container registry.
	// +optional
	Insecure bool `json:"insecure,omitempty"`
}

// ImageSyncSpec defines the desired state of ImageSync
type ImageSyncSpec struct {
	// +required
	SourceImage Image `json:"sourceImage,omitempty"`

	// IsBundleImage allows synchronizing bundle images.
	// +optional
	IsBundleImage bool `json:"isBundleImage,omitempty"`

	// +required
	DestinationImage Image `json:"destinationImage,omitempty"`

	// The timeout for remote OCI Repository operations like pulling, defaults to 60s.
	// +kubebuilder:default="60s"
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="^([0-9]+(\\.[0-9]+)?(ms|s|m))+$"
	// +optional
	Timeout *metav1.Duration `json:"timeout,omitempty"`
}

// ImageSyncStatus defines the observed state of ImageSync
type ImageSyncStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	apis.Status `json:",inline"`

	// URL is the destination link for the latest Artifact.
	// +optional
	URL string `json:"url,omitempty"`

	// LastSyncTime to the destination repository
	LastSyncTime metav1.Time `json:"lastSyncTime,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Image",type=string,JSONPath=`.spec.sourceImage.image`
//+kubebuilder:printcolumn:name="URL",type=string,JSONPath=`.spec.destinationImage.image`
//+kubebuilder:printcolumn:name="Bundle",type=boolean,JSONPath=`.spec.isBundleImage`
//+kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`
//+kubebuilder:printcolumn:name="Reason",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].reason`
//+kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// ImageSync is the Schema for the imagesyncs API
type ImageSync struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ImageSyncSpec   `json:"spec,omitempty"`
	Status ImageSyncStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ImageSyncList contains a list of ImageSync
type ImageSyncList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ImageSync `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ImageSync{}, &ImageSyncList{})
}
