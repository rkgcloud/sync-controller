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
	"net/url"

	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/go-logr/logr"
	"github.com/google/go-containerregistry/pkg/authn/k8schain"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	corev1 "k8s.io/api/core/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"
	"reconciler.io/runtime/reconcilers"

	"carvel.dev/imgpkg/pkg/imgpkg/cmd"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"

	syncv1alpha1 "github.com/rkgcloud/sync-controller/api/sync/v1alpha1"
)

// ImageSyncReconciler reconciles a ImageSync object
const SourceSecretsStashKey reconcilers.StashKey = syncv1alpha1.Group + "/source-secret"
const SourceImageRefStashKey reconcilers.StashKey = syncv1alpha1.Group + "/source-image-ref"
const DestinationSecretsStashKey reconcilers.StashKey = syncv1alpha1.Group + "/destination-secret"

// +kubebuilder:rbac:groups=sync.controller.rkgcloud.com,resources=imagesyncs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=sync.controller.rkgcloud.com,resources=imagesyncs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core,resources=events,verbs=get;list;watch;create;update;patch;delete

// ImageSyncReconciler reconciles a ImageSync object
func ImageSyncReconciler(c reconcilers.Config) *reconcilers.ResourceReconciler[*syncv1alpha1.ImageSync] {
	return &reconcilers.ResourceReconciler[*syncv1alpha1.ImageSync]{
		Reconciler: &reconcilers.Sequence[*syncv1alpha1.ImageSync]{
			SourceSecretSyncReconciler(),
			DestinationSecretSyncReconciler(),
			// SourceImageSyncReconciler(),
			// SynchronizeImageReconciler(),
		},
		Config: c,
	}
}

// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch
// +kubebuilder:rbac:groups=core,resources=serviceaccounts,verbs=get;list;watch

func SourceSecretSyncReconciler() reconcilers.SubReconciler[*syncv1alpha1.ImageSync] {
	return &reconcilers.SyncReconciler[*syncv1alpha1.ImageSync]{
		Name: "SecretSyncReconciler",
		Sync: func(ctx context.Context, resource *syncv1alpha1.ImageSync) error {
			c := reconcilers.RetrieveConfigOrDie(ctx)
			log := logr.FromContextOrDiscard(ctx)

			if resource.Namespace == "" {
				resource.Namespace = "default"
			}

			sourceSecretNames := sets.NewString()
			for _, ps := range resource.Spec.SourceImage.SecretRef {
				sourceSecretNames.Insert(ps.Name)
			}

			// lookup source service account. NOTE: don't use Default SA.
			// Default SA is ambiguous when using as source/dest secret.
			// We will go with the defined SA only.
			if len(resource.Spec.SourceImage.ServiceAccountName) > 0 {
				sourceServiceAccountName := resource.Spec.SourceImage.ServiceAccountName
				serviceAccount := corev1.ServiceAccount{}
				err := c.TrackAndGet(ctx, types.NamespacedName{Namespace: resource.Namespace, Name: sourceServiceAccountName}, &serviceAccount)
				if err != nil {
					if apierrs.IsNotFound(err) {
						resource.ManageConditions().MarkFalse(syncv1alpha1.ImageSyncConditionSourceImageResolved, "ServiceAccountMissing", "ServiceAccount %q not found in namespace %q", sourceServiceAccountName, resource.Namespace)
						return nil
					}
					log.Error(err, "unable to track source service account", sourceServiceAccountName, fmt.Sprintf("%s-%s", resource.Namespace, resource.Name))
					return err
				}

				for _, ips := range serviceAccount.ImagePullSecrets {
					sourceSecretNames.Insert(ips.Name)
				}
			}

			imagePullSecrets := make([]corev1.Secret, len(sourceSecretNames))
			for i, secretName := range sourceSecretNames.List() {
				srcSecret := corev1.Secret{}
				err := c.TrackAndGet(ctx, types.NamespacedName{Namespace: resource.Namespace, Name: secretName}, &srcSecret)
				if err != nil {
					if apierrs.IsNotFound(err) {
						resource.ManageConditions().MarkFalse(syncv1alpha1.ImageSyncConditionSourceImageResolved, "SecretMissing", "Secret %q not found in namespace %q", secretName, resource.Namespace)
						return nil
					}
					log.Error(err, "unable to track source secret", secretName, fmt.Sprintf("%s-%s", resource.Namespace, resource.Name))
					return err
				}
				imagePullSecrets[i] = srcSecret
			}

			stashSecrets(ctx, SourceSecretsStashKey, imagePullSecrets)
			resource.ManageConditions().MarkTrue(syncv1alpha1.ImageSyncConditionSourceImageResolved, "SourceSecret", "Resolved for %s:%s", resource.Name, resource.Namespace)
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

// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch
// +kubebuilder:rbac:groups=core,resources=serviceaccounts,verbs=get;list;watch

func DestinationSecretSyncReconciler() reconcilers.SubReconciler[*syncv1alpha1.ImageSync] {
	return &reconcilers.SyncReconciler[*syncv1alpha1.ImageSync]{
		Name: "DestinationSecretSyncReconciler",
		Sync: func(ctx context.Context, resource *syncv1alpha1.ImageSync) error {
			c := reconcilers.RetrieveConfigOrDie(ctx)
			log := logr.FromContextOrDiscard(ctx)

			if resource.Namespace == "" {
				resource.Namespace = "default"
			}

			destSecretNames := sets.NewString()
			for _, ps := range resource.Spec.DestinationImage.SecretRef {
				destSecretNames.Insert(ps.Name)
			}

			imagePullSecrets := make([]corev1.Secret, len(destSecretNames))
			for i, secretName := range destSecretNames.List() {
				srcSecret := corev1.Secret{}
				err := c.TrackAndGet(ctx, types.NamespacedName{Namespace: resource.Namespace, Name: secretName}, &srcSecret)
				if err != nil {
					if apierrs.IsNotFound(err) {
						resource.ManageConditions().MarkFalse(syncv1alpha1.ImageSyncConditionDestinationImageResolved, "SecretMissing", "Secret %q not found in namespace %q", secretName, resource.Namespace)
						return nil
					}
					log.Error(err, "unable to track destination secret", secretName, fmt.Sprintf("%s-%s", resource.Namespace, resource.Name))
					return err
				}
				imagePullSecrets[i] = srcSecret
			}

			stashSecrets(ctx, DestinationSecretsStashKey, imagePullSecrets)

			resource.ManageConditions().MarkTrue(syncv1alpha1.ImageSyncConditionDestinationImageResolved, "DestinationSecret", "Resolved for %s:%s", resource.Name, resource.Namespace)

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

func SourceImageSyncReconciler() reconcilers.SubReconciler[*syncv1alpha1.ImageSync] {
	return &reconcilers.SyncReconciler[*syncv1alpha1.ImageSync]{
		Name: "SourceImageSyncReconciler",
		Sync: func(ctx context.Context, resource *syncv1alpha1.ImageSync) error {
			log := logr.FromContextOrDiscard(ctx)
			log.Info("SourceImageSyncReconciler", "image", resource.Spec.SourceImage.Image)

			_, err := name.NewDigest(resource.Spec.SourceImage.Image, name.WeakValidation)
			if err == nil {
				// image already resolved to digest
				stashImageRef(ctx, SourceImageRefStashKey, resource.Spec.SourceImage.Image)
				log.Info("SourceImageSyncReconciler", "skipping", "image already resolved to digest")
				return nil
			}

			// resolve source image to digest
			pullSecrets := retrieveSecrets(ctx, SourceSecretsStashKey)

			keychain, err := k8schain.NewFromPullSecrets(ctx, pullSecrets)
			if err != nil {
				log.Error(err, "unable to process source image secret", "secret", resource.Spec.SourceImage.Image)
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

			stashImageRef(ctx, SourceImageRefStashKey, fmt.Sprintf("%s@%s", tag.Name(), image.Digest))
			log.Info("SourceImageSyncReconciler", "resolved", fmt.Sprintf("%s@%s", tag.Name(), image.Digest))
			resource.Status.LastSyncTime = metav1.Now()
			resource.ManageConditions().MarkTrue(syncv1alpha1.ImageSyncConditionSourceImageResolved, "Available", "")

			return nil
		},
	}
}

func SynchronizeImageReconciler() reconcilers.SubReconciler[*syncv1alpha1.ImageSync] {
	return &reconcilers.SyncReconciler[*syncv1alpha1.ImageSync]{
		Name: "SynchronizeImageReconciler",
		Sync: func(ctx context.Context, resource *syncv1alpha1.ImageSync) error {
			log := logr.FromContextOrDiscard(ctx)

			// resolve tagged image to digest
			// pullSecrets := retrieveSecrets(ctx, SourceSecretsStashKey)
			// if pullSecrets == nil {
			// 	log.Info("SynchronizeImageReconciler", "pullSecrets", "no secrets found")
			// 	// return nil
			// }
			// Use carvel.dev imgpkg to copy the image to the target repository

			// Example
			// destimation repo: rkamaldocker/bundle-test
			// Source Bundle: "gcr.io/tanzucliappdev/source-controller-bundle:latest"
			confUI := ui.NewConfUI(ui.NewNoopLogger())
			defer confUI.Flush()
			cmd.NewCopyOptions(confUI)
			dest, err := getRepository(resource.Spec.DestinationImage.Image)
			if err != nil {
				return err
			}
			log.Info("SynchronizeImageReconciler", "copy to dest: ", dest)
			imgpkgCopy := cmd.CopyOptions{
				BundleFlags: cmd.BundleFlags{
					Bundle: resource.Spec.SourceImage.Image,
				},
				Concurrency: 5,
				RepoDst:     resource.Spec.DestinationImage.Image,
			}

			err = imgpkgCopy.Run()
			if err != nil {
				// log.Error(err, "error copying bundle image", "image", resource.Spec.SourceImage.Image)
				resource.Status.LastSyncTime = metav1.Now()
				resource.ManageConditions().MarkFalse(syncv1alpha1.ImageSyncConditionDestinationImageResolved, "Error", "Unable to copy bundle image %s", resource.Spec.SourceImage.Image)
				return nil
			}

			// update subreconciler status
			log.Info("SynchronizeImageReconciler", "resolved", fmt.Sprintf("%s@%s", "tag.Name()", "image.Digest"))
			resource.Status.LastSyncTime = metav1.Now()
			resource.ManageConditions().MarkTrue(syncv1alpha1.ImageSyncConditionDestinationImageResolved, "Synchronized", "")

			return nil
			// NOTE: Run is working. Need to pass auth
		},
	}
}

func stashSecrets(ctx context.Context, key reconcilers.StashKey, pullSecrets []corev1.Secret) {
	reconcilers.StashValue(ctx, key, pullSecrets)
}

func stashImageRef(ctx context.Context, key reconcilers.StashKey, image string) {
	reconcilers.StashValue(ctx, key, image)
}

func retrieveSecrets(ctx context.Context, key reconcilers.StashKey) []corev1.Secret {
	pullSecrets, ok := reconcilers.RetrieveValue(ctx, key).([]corev1.Secret)
	if !ok {
		return nil
	}
	return pullSecrets
}

func getRepository(repo string) (string, error) {
	u, err := url.Parse(repo)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}
