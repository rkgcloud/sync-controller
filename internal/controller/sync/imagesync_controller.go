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

package sync

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/google/go-containerregistry/pkg/authn/k8schain"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/vmware-labs/reconciler-runtime/reconcilers"
	corev1 "k8s.io/api/core/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"

	syncv1alpha1 "github.com/rkgcloud/sync-controller/api/sync/v1alpha1"
)

// ImageSyncReconciler reconciles a ImageSync object
const SourceImagePullSecretsStashKey reconcilers.StashKey = syncv1alpha1.Group + "/source-image-pull-secrets"
const SourceImageRefStashKey reconcilers.StashKey = syncv1alpha1.Group + "/source-image-ref"

// +kubebuilder:rbac:groups=sync.controller.rkgcloud.com,resources=imagesyncs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=sync.controller.rkgcloud.com,resources=imagesyncs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core,resources=events,verbs=get;list;watch;create;update;patch;delete

// ImageSyncReconciler reconciles a ImageSync object
func ImageSyncReconciler(c reconcilers.Config) *reconcilers.ResourceReconciler[*syncv1alpha1.ImageSync] {
	return &reconcilers.ResourceReconciler[*syncv1alpha1.ImageSync]{
		Reconciler: &reconcilers.Sequence[*syncv1alpha1.ImageSync]{
			ImageSyncSourceSecretSyncReconciler(),
			ImageSyncSourceImageReconciler(),
		},
		Config: c,
	}
}

// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch
// +kubebuilder:rbac:groups=core,resources=serviceaccounts,verbs=get;list;watch

func ImageSyncSourceSecretSyncReconciler() reconcilers.SubReconciler[*syncv1alpha1.ImageSync] {
	return &reconcilers.SyncReconciler[*syncv1alpha1.ImageSync]{
		Name: "ImageSyncSourceSecretSyncReconciler",
		Sync: func(ctx context.Context, resource *syncv1alpha1.ImageSync) error {
			c := reconcilers.RetrieveConfigOrDie(ctx)

			pullSecretNames := sets.NewString()

			for _, ps := range resource.Spec.SourceImage.SecretRef {
				pullSecretNames.Insert(ps.Name)
			}

			// lookup service account
			serviceAccountName := resource.Spec.SourceImage.ServiceAccountName
			serviceAccount := corev1.ServiceAccount{}
			err := c.TrackAndGet(ctx, types.NamespacedName{Namespace: resource.Namespace, Name: serviceAccountName}, &serviceAccount)
			if err != nil {
				if apierrs.IsNotFound(err) {
					resource.ManageConditions().MarkFalse(syncv1alpha1.ImageSyncConditionSourceImageResolved, "ServiceAccountMissing", "ServiceAccount %q not found in namespace %q", serviceAccountName, resource.Namespace)
					return nil
				}
				return err
			}

			for _, ips := range serviceAccount.ImagePullSecrets {
				pullSecretNames.Insert(ips.Name)
			}

			// lookup image pull secrets
			imagePullSecrets := make([]corev1.Secret, len(pullSecretNames))
			for i, imagePullSecretName := range pullSecretNames.List() {
				imagePullSecret := corev1.Secret{}
				err := c.TrackAndGet(ctx, types.NamespacedName{Namespace: resource.Namespace, Name: imagePullSecretName}, &imagePullSecret)
				if err != nil {
					if apierrs.IsNotFound(err) {
						resource.ManageConditions().MarkFalse(syncv1alpha1.ImageSyncConditionSourceImageResolved, "SecretMissing", "Secret %q not found in namespace %q", imagePullSecretName, resource.Namespace)
						return nil
					}
					return err
				}
				imagePullSecrets[i] = imagePullSecret
			}

			StashSourceImagePullSecrets(ctx, imagePullSecrets)

			return nil
		},
	}
}

func ImageSyncSourceImageReconciler() reconcilers.SubReconciler[*syncv1alpha1.ImageSync] {
	return &reconcilers.SyncReconciler[*syncv1alpha1.ImageSync]{
		Name: "ImageSyncSourceImageReconciler",
		Sync: func(ctx context.Context, resource *syncv1alpha1.ImageSync) error {
			// c := reconcilers.RetrieveConfigOrDie(ctx)
			log := logr.FromContextOrDiscard(ctx)

			_, err := name.NewDigest(resource.Spec.SourceImage.Image, name.WeakValidation)
			if err == nil {
				// image already resolved to digest
				StashImageRef(ctx, resource.Spec.SourceImage.Image)
				return nil
			}

			// resolve tagged image to digest
			pullSecrets := RetrieveImagePullSecrets(ctx)
			if pullSecrets == nil {
				return nil
			}
			keychain, err := k8schain.NewFromPullSecrets(ctx, pullSecrets)
			if err != nil {
				return err
			}
			tag, err := name.NewTag(resource.Spec.SourceImage.Image, name.WeakValidation)
			if err != nil {
				resource.ManageConditions().MarkFalse(syncv1alpha1.ImageSyncConditionSourceImageResolved, "MalformedRepository", "Image name %q failed validation: %s", resource.Spec.SourceImage.Image, err)
				return nil
			}
			image, err := remote.Head(tag, remote.WithContext(ctx), remote.WithAuthFromKeychain(keychain))
			if err != nil {
				log.Error(err, "unable to resolve image tag to a digest", "image", resource.Spec.SourceImage.Image)
				resource.ManageConditions().MarkFalse(syncv1alpha1.ImageSyncConditionSourceImageResolved, "RemoteError", "Unable to resolve image with tag %q to a digest: %s", resource.Spec.SourceImage.Image, err)
				return nil
			}

			StashImageRef(ctx, fmt.Sprintf("%s@%s", tag.Name(), image.Digest))
			resource.Status.LastSyncTime = metav1.Now()
			resource.ManageConditions().MarkTrue(syncv1alpha1.ImageSyncConditionSourceImageResolved, "Available", "")

			return nil
		},
	}
}

func StashSourceImagePullSecrets(ctx context.Context, pullSecrets []corev1.Secret) {
	reconcilers.StashValue(ctx, SourceImagePullSecretsStashKey, pullSecrets)
}

func StashImageRef(ctx context.Context, image string) {
	reconcilers.StashValue(ctx, SourceImageRefStashKey, image)
}

func RetrieveImagePullSecrets(ctx context.Context) []corev1.Secret {
	pullSecrets, ok := reconcilers.RetrieveValue(ctx, SourceImageRefStashKey).([]corev1.Secret)
	if !ok {
		return nil
	}
	return pullSecrets
}
