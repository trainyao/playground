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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//GoddessMomentSpec .
type GoddessMomentSpec struct {
	FoodDemand []FoodDemand `json:"foodDemand,omitempty"`
}

//FoodDemand .
type FoodDemand struct {
	Name string `json:"name"`
}

//GoddessMomentStatus .
type GoddessMomentStatus struct {
	FoodDemand []FoodDemandStatus `json:"foodDemand,omitempty"`
}

//FoodDemandStatus .
type FoodDemandStatus struct {
	Name        string      `json:"name,omitempty"`
	Status      FoodStatus  `json:"status,omitempty"`
	ClaimTime   metav1.Time `json:"claimTime,omitempty"`
	ArrivalTime metav1.Time `json:"arrivalTime,omitempty"`
	ClaimBy     string      `json:"claimBy,omitempty"`
}

//FoodStatus .
type FoodStatus string

const (
	//FoodStatusPending .
	FoodStatusPending = "Pending"
	//FoodStatusPendingArrival .
	FoodStatusPendingArrival = "PendingArrival"
	//FoodStatusArrived .
	FoodStatusArrived = "Arrived"
)

//+kubebuilder:object:root=true

// GoddessMoment is the Schema for the goddessmoments API
type GoddessMoment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GoddessMomentSpec   `json:"spec,omitempty"`
	Status GoddessMomentStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// GoddessMomentList contains a list of GoddessMoment
type GoddessMomentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GoddessMoment `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GoddessMoment{}, &GoddessMomentList{})
}
