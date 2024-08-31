package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	diemetav1 "reconciler.io/dies/apis/meta/v1"

	syncv1alpha1 "github.com/rkgcloud/sync-controller/api/sync/v1alpha1"
)

// +die:object=true
type _ = syncv1alpha1.ImageSync

// +die
type _ = syncv1alpha1.ImageSyncSpec

// +die
type _ = syncv1alpha1.ImageSyncStatus

// +die
type _ = syncv1alpha1.ImageSource

func (d *ImageSyncStatusDie) ConditionsDie(conditions ...*diemetav1.ConditionDie) *ImageSyncStatusDie {
	return d.DieStamp(func(r *syncv1alpha1.ImageSyncStatus) {
		r.Conditions = make([]metav1.Condition, len(conditions))
		for i := range conditions {
			r.Conditions[i] = conditions[i].DieRelease()
		}
	})
}

func (d *ImageSyncStatusDie) ObservedGeneration(v int64) *ImageSyncStatusDie {
	return d.DieStamp(func(r *syncv1alpha1.ImageSyncStatus) {
		r.ObservedGeneration = v
	})
}
