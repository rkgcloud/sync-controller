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
	"reconciler.io/runtime/apis"
)

var (
	ImageSyncLabelKey = GroupVersion.Group + "/image-sync"
)

const (
	ImageSyncConditionReady                    = apis.ConditionReady
	ImageSyncConditionSourceImageResolved      = "SourceImageResolved"
	ImageSyncConditionDestinationImageResolved = "DestinationImageResolved"
)

var imagesyncConditionSet = apis.NewLivingConditionSet(
	ImageSyncConditionSourceImageResolved,
	ImageSyncConditionDestinationImageResolved,
)

func (r *ImageSync) ManageConditions() apis.ConditionManager {
	return r.GetConditionSet().Manage(r.GetConditionsAccessor())
}

func (r *ImageSync) GetConditionsAccessor() apis.ConditionsAccessor {
	return &r.Status
}

func (r *ImageSync) GetConditionSet() apis.ConditionSet {
	return imagesyncConditionSet
}

func (r *ImageSyncStatus) InitializeConditions() {
	// reset conditions
	r.Conditions = nil
	imagesyncConditionSet.Manage(r).InitializeConditions()
}
