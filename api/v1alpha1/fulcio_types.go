package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// FulcioSpec defines the desired state of Fulcio
type FulcioSpec struct {
	External    bool                  `json:"external,omitempty"`
	KeySecret   string                `json:"keySecret,omitempty"`
	OidcIssuers map[string]OidcIssuer `json:"oidcIssuers,omitempty"`
}

type OidcIssuer struct {
	ClientID  string `json:"ClientID"`
	IssuerURL string `json:"IssuerURL"`
	Type      string `json:"Type"`
}

// FulcioStatus defines the observed state of Fulcio
type FulcioStatus struct {
	Url   string `json:"url,omitempty"`
	Phase Phase  `json:"Phase,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Fulcio is the Schema for the fulcios API
type Fulcio struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FulcioSpec   `json:"spec,omitempty"`
	Status FulcioStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// FulcioList contains a list of Fulcio
type FulcioList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Fulcio `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Fulcio{}, &FulcioList{})
}
