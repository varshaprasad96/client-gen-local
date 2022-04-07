package controller

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
type Foo struct {
	// TypeMeta comments should NOT appear in the CRD spec
	metav1.TypeMeta `json:",inline"`
	// ObjectMeta comments should NOT appear in the CRD spec
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec comments SHOULD appear in the CRD spec
	Spec FooSpec `json:"spec,omitempty"`
	// Status comments SHOULD appear in the CRD spec
	Status FooStatus `json:"status,omitempty"`
}

type FooStatus struct{}

type FooSpec struct {
	// This tests that defaulted fields are stripped for v1beta1,
	// but not for v1
	// +kubebuilder:default=fooDefaultString
	DefaultedString string `json:"defaultedString"`
}
