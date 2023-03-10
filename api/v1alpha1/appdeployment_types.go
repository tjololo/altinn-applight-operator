/*
Copyright 2022.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AppDeploymentSpec defines the desired state of AppDeployment
type AppDeploymentSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Image defines the image to deploy
	Image string `json:"foo,omitempty"`

	// LayoutSets defines layoutsets for the app
	LayoutSets []LayoutSet `json:"layouts,omitempty"`
	// Process xml string defining the apps process
	Process       string          `json:"process"`
	LanguageTexts map[string]Text `json:"languageTexts"`
	Policy        string          `json:"policy"`
}

type LayoutSet struct {
	Id    string   `json:"id"`
	Model string   `json:"model"`
	Tasks []string `json:"tasks"`
	Pages []Page   `json:"pages"`
}

type Page struct {
	Name   string `json:"name"`
	Config string `json:"config"`
}

type Text struct {
	Text map[string]string `json:",inline"`
}

// AppDeploymentStatus defines the observed state of AppDeployment
type AppDeploymentStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// AppDeployment is the Schema for the appdeployments API
type AppDeployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AppDeploymentSpec   `json:"spec,omitempty"`
	Status AppDeploymentStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AppDeploymentList contains a list of AppDeployment
type AppDeploymentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AppDeployment `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AppDeployment{}, &AppDeploymentList{})
}
