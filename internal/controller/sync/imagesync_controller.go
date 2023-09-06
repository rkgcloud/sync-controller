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

	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"

	syncv1alpha1 "github.com/rkgcloud/sync-controller/api/sync/v1alpha1"
)

// ImageSyncReconciler reconciles a ImageSync object
const SourceImagePullSecretsStashKey reconcilers.StashKey = syncv1alpha1.Group + "/source-image-pull-secrets"
const SourceImageRefStashKey reconcilers.StashKey = syncv1alpha1.Group + "/source-image-ref"
const DestinationSecretsStashKey reconcilers.StashKey = syncv1alpha1.Group + "/destination-secret"

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
			log := logr.FromContextOrDiscard(ctx)
			log.Info("ImageSyncSourceSecretSyncReconciler", "sub-reconciler", "ImageSyncSourceSecretSyncReconciler")

			pullSecretNames := sets.NewString()

			for _, ps := range resource.Spec.SourceImage.SecretRef {
				pullSecretNames.Insert(ps.Name)
			}

			// lookup source service account
			sourceServiceAccountName := resource.Spec.SourceImage.ServiceAccountName
			if sourceServiceAccountName == "" {
				sourceServiceAccountName = "default"
			}
			serviceAccount := corev1.ServiceAccount{}
			err := c.TrackAndGet(ctx, types.NamespacedName{Namespace: resource.Namespace, Name: sourceServiceAccountName}, &serviceAccount)
			if err != nil {
				if apierrs.IsNotFound(err) {
					resource.ManageConditions().MarkFalse(syncv1alpha1.ImageSyncConditionSourceImageResolved, "ServiceAccountMissing", "ServiceAccount %q not found in namespace %q", sourceServiceAccountName, resource.Namespace)
					return nil
				}
				return err
			}

			for _, ips := range serviceAccount.ImagePullSecrets {
				pullSecretNames.Insert(ips.Name)
			}

			// lookup destination service acount
			destServiceAccountName := resource.Spec.DestinationImage.ServiceAccountName
			if destServiceAccountName == "" && sourceServiceAccountName != "default" {
				destServiceAccountName = "default"
			}

			destServiceAccount := corev1.ServiceAccount{}
			errx := c.TrackAndGet(ctx, types.NamespacedName{Namespace: resource.Namespace, Name: destServiceAccountName}, &destServiceAccount)
			if errx != nil {
				if apierrs.IsNotFound(err) {
					resource.ManageConditions().MarkFalse(syncv1alpha1.ImageSyncConditionSourceImageResolved, "ServiceAccountMissing", "ServiceAccount %q not found in namespace %q", destServiceAccountName, resource.Namespace)
					return nil
				}
				return err
			}

			for _, ds := range destServiceAccount.ImagePullSecrets {
				pullSecretNames.Insert(ds.Name)
			}

			log.Info("ImageSyncSourceSecretSyncReconciler", "pullSecretNames len", len(pullSecretNames))

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

			stashSecrets(ctx, imagePullSecrets)

			return nil
		},

		Setup: func(ctx context.Context, mgr controllerruntime.Manager, bldr *builder.Builder) error {
			// register an informer to watch Secret's metadata only. This reduces the cache size in memory.
			bldr.Watches(&corev1.Secret{}, reconcilers.EnqueueTracked(ctx), builder.OnlyMetadata)
			// register an informer to watch ServiceAccounts
			bldr.Watches(&corev1.ServiceAccount{}, reconcilers.EnqueueTracked(ctx))

			return nil
		},
	}
}

func ImageSyncSourceImageReconciler() reconcilers.SubReconciler[*syncv1alpha1.ImageSync] {
	return &reconcilers.SyncReconciler[*syncv1alpha1.ImageSync]{
		Name: "ImageSyncSourceImageReconciler",
		Sync: func(ctx context.Context, resource *syncv1alpha1.ImageSync) error {
			log := logr.FromContextOrDiscard(ctx)
			log.Info("ImageSyncSourceImageReconciler", "image", resource.Spec.SourceImage.Image)

			_, err := name.NewDigest(resource.Spec.SourceImage.Image, name.WeakValidation)
			if err == nil {
				// image already resolved to digest
				stashImageRef(ctx, resource.Spec.SourceImage.Image)
				log.Info("ImageSyncSourceImageReconciler", "skipping", "image already resolved to digest")
				return nil
			}

			// resolve tagged image to digest
			pullSecrets := retrieveSecrets(ctx)
			if pullSecrets == nil {
				log.Info("ImageSyncSourceImageReconciler", "pullSecrets", "no pull secrets")
				// return nil
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

			stashImageRef(ctx, fmt.Sprintf("%s@%s", tag.Name(), image.Digest))
			log.Info("ImageSyncSourceImageReconciler", "resolved", fmt.Sprintf("%s@%s", tag.Name(), image.Digest))
			resource.Status.LastSyncTime = metav1.Now()
			resource.ManageConditions().MarkTrue(syncv1alpha1.ImageSyncConditionSourceImageResolved, "Available", "")

			return nil
		},
	}
}

func stashSecrets(ctx context.Context, pullSecrets []corev1.Secret) {
	reconcilers.StashValue(ctx, SourceImagePullSecretsStashKey, pullSecrets)
}

func stashImageRef(ctx context.Context, image string) {
	reconcilers.StashValue(ctx, SourceImageRefStashKey, image)
}

func retrieveSecrets(ctx context.Context) []corev1.Secret {
	pullSecrets, ok := reconcilers.RetrieveValue(ctx, SourceImageRefStashKey).([]corev1.Secret)
	if !ok {
		return nil
	}
	return pullSecrets
}
