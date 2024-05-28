//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by diegen. DO NOT EDIT.

package v1alpha1

import (
	json "encoding/json"
	fmtx "fmt"
	osx "os"
	reflectx "reflect"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	unstructured "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	runtime "k8s.io/apimachinery/pkg/runtime"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	jsonpath "k8s.io/client-go/util/jsonpath"
	v1 "reconciler.io/dies/apis/meta/v1"
	apis "reconciler.io/runtime/apis"
	yaml "sigs.k8s.io/yaml"

	syncv1alpha1 "github.com/rkgcloud/sync-controller/api/sync/v1alpha1"
)

var ImageSyncBlank = (&ImageSyncDie{}).DieFeed(syncv1alpha1.ImageSync{})

type ImageSyncDie struct {
	v1.FrozenObjectMeta
	mutable bool
	r       syncv1alpha1.ImageSync
}

// DieImmutable returns a new die for the current die's state that is either mutable (`false`) or immutable (`true`).
func (d *ImageSyncDie) DieImmutable(immutable bool) *ImageSyncDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy()
	d.mutable = !immutable
	return d
}

// DieFeed returns a new die with the provided resource.
func (d *ImageSyncDie) DieFeed(r syncv1alpha1.ImageSync) *ImageSyncDie {
	if d.mutable {
		d.FrozenObjectMeta = v1.FreezeObjectMeta(r.ObjectMeta)
		d.r = r
		return d
	}
	return &ImageSyncDie{
		FrozenObjectMeta: v1.FreezeObjectMeta(r.ObjectMeta),
		mutable:          d.mutable,
		r:                r,
	}
}

// DieFeedPtr returns a new die with the provided resource pointer. If the resource is nil, the empty value is used instead.
func (d *ImageSyncDie) DieFeedPtr(r *syncv1alpha1.ImageSync) *ImageSyncDie {
	if r == nil {
		r = &syncv1alpha1.ImageSync{}
	}
	return d.DieFeed(*r)
}

// DieFeedJSON returns a new die with the provided JSON. Panics on error.
func (d *ImageSyncDie) DieFeedJSON(j []byte) *ImageSyncDie {
	r := syncv1alpha1.ImageSync{}
	if err := json.Unmarshal(j, &r); err != nil {
		panic(err)
	}
	return d.DieFeed(r)
}

// DieFeedYAML returns a new die with the provided YAML. Panics on error.
func (d *ImageSyncDie) DieFeedYAML(y []byte) *ImageSyncDie {
	r := syncv1alpha1.ImageSync{}
	if err := yaml.Unmarshal(y, &r); err != nil {
		panic(err)
	}
	return d.DieFeed(r)
}

// DieFeedYAMLFile returns a new die loading YAML from a file path. Panics on error.
func (d *ImageSyncDie) DieFeedYAMLFile(name string) *ImageSyncDie {
	y, err := osx.ReadFile(name)
	if err != nil {
		panic(err)
	}
	return d.DieFeedYAML(y)
}

// DieFeedRawExtension returns the resource managed by the die as an raw extension. Panics on error.
func (d *ImageSyncDie) DieFeedRawExtension(raw runtime.RawExtension) *ImageSyncDie {
	j, err := json.Marshal(raw)
	if err != nil {
		panic(err)
	}
	return d.DieFeedJSON(j)
}

// DieRelease returns the resource managed by the die.
func (d *ImageSyncDie) DieRelease() syncv1alpha1.ImageSync {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

// DieReleasePtr returns a pointer to the resource managed by the die.
func (d *ImageSyncDie) DieReleasePtr() *syncv1alpha1.ImageSync {
	r := d.DieRelease()
	return &r
}

// DieReleaseUnstructured returns the resource managed by the die as an unstructured object. Panics on error.
func (d *ImageSyncDie) DieReleaseUnstructured() *unstructured.Unstructured {
	r := d.DieReleasePtr()
	u, err := runtime.DefaultUnstructuredConverter.ToUnstructured(r)
	if err != nil {
		panic(err)
	}
	return &unstructured.Unstructured{
		Object: u,
	}
}

// DieReleaseJSON returns the resource managed by the die as JSON. Panics on error.
func (d *ImageSyncDie) DieReleaseJSON() []byte {
	r := d.DieReleasePtr()
	j, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	return j
}

// DieReleaseYAML returns the resource managed by the die as YAML. Panics on error.
func (d *ImageSyncDie) DieReleaseYAML() []byte {
	r := d.DieReleasePtr()
	y, err := yaml.Marshal(r)
	if err != nil {
		panic(err)
	}
	return y
}

// DieReleaseRawExtension returns the resource managed by the die as an raw extension. Panics on error.
func (d *ImageSyncDie) DieReleaseRawExtension() runtime.RawExtension {
	j := d.DieReleaseJSON()
	raw := runtime.RawExtension{}
	if err := json.Unmarshal(j, &raw); err != nil {
		panic(err)
	}
	return raw
}

// DieStamp returns a new die with the resource passed to the callback function. The resource is mutable.
func (d *ImageSyncDie) DieStamp(fn func(r *syncv1alpha1.ImageSync)) *ImageSyncDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

// Experimental: DieStampAt uses a JSON path (http://goessner.net/articles/JsonPath/) expression to stamp portions of the resource. The callback is invoked with each JSON path match. Panics if the callback function does not accept a single argument of the same type or a pointer to that type as found on the resource at the target location.
//
// Future iterations will improve type coercion from the resource to the callback argument.
func (d *ImageSyncDie) DieStampAt(jp string, fn interface{}) *ImageSyncDie {
	return d.DieStamp(func(r *syncv1alpha1.ImageSync) {
		if ni := reflectx.ValueOf(fn).Type().NumIn(); ni != 1 {
			panic(fmtx.Errorf("callback function must have 1 input parameters, found %d", ni))
		}
		if no := reflectx.ValueOf(fn).Type().NumOut(); no != 0 {
			panic(fmtx.Errorf("callback function must have 0 output parameters, found %d", no))
		}

		cp := jsonpath.New("")
		if err := cp.Parse(fmtx.Sprintf("{%s}", jp)); err != nil {
			panic(err)
		}
		cr, err := cp.FindResults(r)
		if err != nil {
			// errors are expected if a path is not found
			return
		}
		for _, cv := range cr[0] {
			arg0t := reflectx.ValueOf(fn).Type().In(0)

			var args []reflectx.Value
			if cv.Type().AssignableTo(arg0t) {
				args = []reflectx.Value{cv}
			} else if cv.CanAddr() && cv.Addr().Type().AssignableTo(arg0t) {
				args = []reflectx.Value{cv.Addr()}
			} else {
				panic(fmtx.Errorf("callback function must accept value of type %q, found type %q", cv.Type(), arg0t))
			}

			reflectx.ValueOf(fn).Call(args)
		}
	})
}

// DieWith returns a new die after passing the current die to the callback function. The passed die is mutable.
func (d *ImageSyncDie) DieWith(fns ...func(d *ImageSyncDie)) *ImageSyncDie {
	nd := ImageSyncBlank.DieFeed(d.DieRelease()).DieImmutable(false)
	for _, fn := range fns {
		if fn != nil {
			fn(nd)
		}
	}
	return d.DieFeed(nd.DieRelease())
}

// DeepCopy returns a new die with equivalent state. Useful for snapshotting a mutable die.
func (d *ImageSyncDie) DeepCopy() *ImageSyncDie {
	r := *d.r.DeepCopy()
	return &ImageSyncDie{
		FrozenObjectMeta: v1.FreezeObjectMeta(r.ObjectMeta),
		mutable:          d.mutable,
		r:                r,
	}
}

var _ runtime.Object = (*ImageSyncDie)(nil)

func (d *ImageSyncDie) DeepCopyObject() runtime.Object {
	return d.r.DeepCopy()
}

func (d *ImageSyncDie) GetObjectKind() schema.ObjectKind {
	r := d.DieRelease()
	return r.GetObjectKind()
}

func (d *ImageSyncDie) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.r)
}

func (d *ImageSyncDie) UnmarshalJSON(b []byte) error {
	if d == ImageSyncBlank {
		return fmtx.Errorf("cannot unmarshal into the blank die, create a copy first")
	}
	if !d.mutable {
		return fmtx.Errorf("cannot unmarshal into immutable dies, create a mutable version first")
	}
	r := &syncv1alpha1.ImageSync{}
	err := json.Unmarshal(b, r)
	*d = *d.DieFeed(*r)
	return err
}

// APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
func (d *ImageSyncDie) APIVersion(v string) *ImageSyncDie {
	return d.DieStamp(func(r *syncv1alpha1.ImageSync) {
		r.APIVersion = v
	})
}

// Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
func (d *ImageSyncDie) Kind(v string) *ImageSyncDie {
	return d.DieStamp(func(r *syncv1alpha1.ImageSync) {
		r.Kind = v
	})
}

// MetadataDie stamps the resource's ObjectMeta field with a mutable die.
func (d *ImageSyncDie) MetadataDie(fn func(d *v1.ObjectMetaDie)) *ImageSyncDie {
	return d.DieStamp(func(r *syncv1alpha1.ImageSync) {
		d := v1.ObjectMetaBlank.DieImmutable(false).DieFeed(r.ObjectMeta)
		fn(d)
		r.ObjectMeta = d.DieRelease()
	})
}

// SpecDie stamps the resource's spec field with a mutable die.
func (d *ImageSyncDie) SpecDie(fn func(d *ImageSyncSpecDie)) *ImageSyncDie {
	return d.DieStamp(func(r *syncv1alpha1.ImageSync) {
		d := ImageSyncSpecBlank.DieImmutable(false).DieFeed(r.Spec)
		fn(d)
		r.Spec = d.DieRelease()
	})
}

// StatusDie stamps the resource's status field with a mutable die.
func (d *ImageSyncDie) StatusDie(fn func(d *ImageSyncStatusDie)) *ImageSyncDie {
	return d.DieStamp(func(r *syncv1alpha1.ImageSync) {
		d := ImageSyncStatusBlank.DieImmutable(false).DieFeed(r.Status)
		fn(d)
		r.Status = d.DieRelease()
	})
}

func (d *ImageSyncDie) Spec(v syncv1alpha1.ImageSyncSpec) *ImageSyncDie {
	return d.DieStamp(func(r *syncv1alpha1.ImageSync) {
		r.Spec = v
	})
}

func (d *ImageSyncDie) Status(v syncv1alpha1.ImageSyncStatus) *ImageSyncDie {
	return d.DieStamp(func(r *syncv1alpha1.ImageSync) {
		r.Status = v
	})
}

var ImageSyncSpecBlank = (&ImageSyncSpecDie{}).DieFeed(syncv1alpha1.ImageSyncSpec{})

type ImageSyncSpecDie struct {
	mutable bool
	r       syncv1alpha1.ImageSyncSpec
}

// DieImmutable returns a new die for the current die's state that is either mutable (`false`) or immutable (`true`).
func (d *ImageSyncSpecDie) DieImmutable(immutable bool) *ImageSyncSpecDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy()
	d.mutable = !immutable
	return d
}

// DieFeed returns a new die with the provided resource.
func (d *ImageSyncSpecDie) DieFeed(r syncv1alpha1.ImageSyncSpec) *ImageSyncSpecDie {
	if d.mutable {
		d.r = r
		return d
	}
	return &ImageSyncSpecDie{
		mutable: d.mutable,
		r:       r,
	}
}

// DieFeedPtr returns a new die with the provided resource pointer. If the resource is nil, the empty value is used instead.
func (d *ImageSyncSpecDie) DieFeedPtr(r *syncv1alpha1.ImageSyncSpec) *ImageSyncSpecDie {
	if r == nil {
		r = &syncv1alpha1.ImageSyncSpec{}
	}
	return d.DieFeed(*r)
}

// DieFeedJSON returns a new die with the provided JSON. Panics on error.
func (d *ImageSyncSpecDie) DieFeedJSON(j []byte) *ImageSyncSpecDie {
	r := syncv1alpha1.ImageSyncSpec{}
	if err := json.Unmarshal(j, &r); err != nil {
		panic(err)
	}
	return d.DieFeed(r)
}

// DieFeedYAML returns a new die with the provided YAML. Panics on error.
func (d *ImageSyncSpecDie) DieFeedYAML(y []byte) *ImageSyncSpecDie {
	r := syncv1alpha1.ImageSyncSpec{}
	if err := yaml.Unmarshal(y, &r); err != nil {
		panic(err)
	}
	return d.DieFeed(r)
}

// DieFeedYAMLFile returns a new die loading YAML from a file path. Panics on error.
func (d *ImageSyncSpecDie) DieFeedYAMLFile(name string) *ImageSyncSpecDie {
	y, err := osx.ReadFile(name)
	if err != nil {
		panic(err)
	}
	return d.DieFeedYAML(y)
}

// DieFeedRawExtension returns the resource managed by the die as an raw extension. Panics on error.
func (d *ImageSyncSpecDie) DieFeedRawExtension(raw runtime.RawExtension) *ImageSyncSpecDie {
	j, err := json.Marshal(raw)
	if err != nil {
		panic(err)
	}
	return d.DieFeedJSON(j)
}

// DieRelease returns the resource managed by the die.
func (d *ImageSyncSpecDie) DieRelease() syncv1alpha1.ImageSyncSpec {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

// DieReleasePtr returns a pointer to the resource managed by the die.
func (d *ImageSyncSpecDie) DieReleasePtr() *syncv1alpha1.ImageSyncSpec {
	r := d.DieRelease()
	return &r
}

// DieReleaseJSON returns the resource managed by the die as JSON. Panics on error.
func (d *ImageSyncSpecDie) DieReleaseJSON() []byte {
	r := d.DieReleasePtr()
	j, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	return j
}

// DieReleaseYAML returns the resource managed by the die as YAML. Panics on error.
func (d *ImageSyncSpecDie) DieReleaseYAML() []byte {
	r := d.DieReleasePtr()
	y, err := yaml.Marshal(r)
	if err != nil {
		panic(err)
	}
	return y
}

// DieReleaseRawExtension returns the resource managed by the die as an raw extension. Panics on error.
func (d *ImageSyncSpecDie) DieReleaseRawExtension() runtime.RawExtension {
	j := d.DieReleaseJSON()
	raw := runtime.RawExtension{}
	if err := json.Unmarshal(j, &raw); err != nil {
		panic(err)
	}
	return raw
}

// DieStamp returns a new die with the resource passed to the callback function. The resource is mutable.
func (d *ImageSyncSpecDie) DieStamp(fn func(r *syncv1alpha1.ImageSyncSpec)) *ImageSyncSpecDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

// Experimental: DieStampAt uses a JSON path (http://goessner.net/articles/JsonPath/) expression to stamp portions of the resource. The callback is invoked with each JSON path match. Panics if the callback function does not accept a single argument of the same type or a pointer to that type as found on the resource at the target location.
//
// Future iterations will improve type coercion from the resource to the callback argument.
func (d *ImageSyncSpecDie) DieStampAt(jp string, fn interface{}) *ImageSyncSpecDie {
	return d.DieStamp(func(r *syncv1alpha1.ImageSyncSpec) {
		if ni := reflectx.ValueOf(fn).Type().NumIn(); ni != 1 {
			panic(fmtx.Errorf("callback function must have 1 input parameters, found %d", ni))
		}
		if no := reflectx.ValueOf(fn).Type().NumOut(); no != 0 {
			panic(fmtx.Errorf("callback function must have 0 output parameters, found %d", no))
		}

		cp := jsonpath.New("")
		if err := cp.Parse(fmtx.Sprintf("{%s}", jp)); err != nil {
			panic(err)
		}
		cr, err := cp.FindResults(r)
		if err != nil {
			// errors are expected if a path is not found
			return
		}
		for _, cv := range cr[0] {
			arg0t := reflectx.ValueOf(fn).Type().In(0)

			var args []reflectx.Value
			if cv.Type().AssignableTo(arg0t) {
				args = []reflectx.Value{cv}
			} else if cv.CanAddr() && cv.Addr().Type().AssignableTo(arg0t) {
				args = []reflectx.Value{cv.Addr()}
			} else {
				panic(fmtx.Errorf("callback function must accept value of type %q, found type %q", cv.Type(), arg0t))
			}

			reflectx.ValueOf(fn).Call(args)
		}
	})
}

// DieWith returns a new die after passing the current die to the callback function. The passed die is mutable.
func (d *ImageSyncSpecDie) DieWith(fns ...func(d *ImageSyncSpecDie)) *ImageSyncSpecDie {
	nd := ImageSyncSpecBlank.DieFeed(d.DieRelease()).DieImmutable(false)
	for _, fn := range fns {
		if fn != nil {
			fn(nd)
		}
	}
	return d.DieFeed(nd.DieRelease())
}

// DeepCopy returns a new die with equivalent state. Useful for snapshotting a mutable die.
func (d *ImageSyncSpecDie) DeepCopy() *ImageSyncSpecDie {
	r := *d.r.DeepCopy()
	return &ImageSyncSpecDie{
		mutable: d.mutable,
		r:       r,
	}
}

func (d *ImageSyncSpecDie) SourceImage(v syncv1alpha1.Image) *ImageSyncSpecDie {
	return d.DieStamp(func(r *syncv1alpha1.ImageSyncSpec) {
		r.SourceImage = v
	})
}

func (d *ImageSyncSpecDie) DestinationImage(v syncv1alpha1.Image) *ImageSyncSpecDie {
	return d.DieStamp(func(r *syncv1alpha1.ImageSyncSpec) {
		r.DestinationImage = v
	})
}

// The timeout for remote OCI Repository operations like pulling, defaults to 60s.
func (d *ImageSyncSpecDie) Timeout(v *metav1.Duration) *ImageSyncSpecDie {
	return d.DieStamp(func(r *syncv1alpha1.ImageSyncSpec) {
		r.Timeout = v
	})
}

var ImageSyncStatusBlank = (&ImageSyncStatusDie{}).DieFeed(syncv1alpha1.ImageSyncStatus{})

type ImageSyncStatusDie struct {
	mutable bool
	r       syncv1alpha1.ImageSyncStatus
}

// DieImmutable returns a new die for the current die's state that is either mutable (`false`) or immutable (`true`).
func (d *ImageSyncStatusDie) DieImmutable(immutable bool) *ImageSyncStatusDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy()
	d.mutable = !immutable
	return d
}

// DieFeed returns a new die with the provided resource.
func (d *ImageSyncStatusDie) DieFeed(r syncv1alpha1.ImageSyncStatus) *ImageSyncStatusDie {
	if d.mutable {
		d.r = r
		return d
	}
	return &ImageSyncStatusDie{
		mutable: d.mutable,
		r:       r,
	}
}

// DieFeedPtr returns a new die with the provided resource pointer. If the resource is nil, the empty value is used instead.
func (d *ImageSyncStatusDie) DieFeedPtr(r *syncv1alpha1.ImageSyncStatus) *ImageSyncStatusDie {
	if r == nil {
		r = &syncv1alpha1.ImageSyncStatus{}
	}
	return d.DieFeed(*r)
}

// DieFeedJSON returns a new die with the provided JSON. Panics on error.
func (d *ImageSyncStatusDie) DieFeedJSON(j []byte) *ImageSyncStatusDie {
	r := syncv1alpha1.ImageSyncStatus{}
	if err := json.Unmarshal(j, &r); err != nil {
		panic(err)
	}
	return d.DieFeed(r)
}

// DieFeedYAML returns a new die with the provided YAML. Panics on error.
func (d *ImageSyncStatusDie) DieFeedYAML(y []byte) *ImageSyncStatusDie {
	r := syncv1alpha1.ImageSyncStatus{}
	if err := yaml.Unmarshal(y, &r); err != nil {
		panic(err)
	}
	return d.DieFeed(r)
}

// DieFeedYAMLFile returns a new die loading YAML from a file path. Panics on error.
func (d *ImageSyncStatusDie) DieFeedYAMLFile(name string) *ImageSyncStatusDie {
	y, err := osx.ReadFile(name)
	if err != nil {
		panic(err)
	}
	return d.DieFeedYAML(y)
}

// DieFeedRawExtension returns the resource managed by the die as an raw extension. Panics on error.
func (d *ImageSyncStatusDie) DieFeedRawExtension(raw runtime.RawExtension) *ImageSyncStatusDie {
	j, err := json.Marshal(raw)
	if err != nil {
		panic(err)
	}
	return d.DieFeedJSON(j)
}

// DieRelease returns the resource managed by the die.
func (d *ImageSyncStatusDie) DieRelease() syncv1alpha1.ImageSyncStatus {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

// DieReleasePtr returns a pointer to the resource managed by the die.
func (d *ImageSyncStatusDie) DieReleasePtr() *syncv1alpha1.ImageSyncStatus {
	r := d.DieRelease()
	return &r
}

// DieReleaseJSON returns the resource managed by the die as JSON. Panics on error.
func (d *ImageSyncStatusDie) DieReleaseJSON() []byte {
	r := d.DieReleasePtr()
	j, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	return j
}

// DieReleaseYAML returns the resource managed by the die as YAML. Panics on error.
func (d *ImageSyncStatusDie) DieReleaseYAML() []byte {
	r := d.DieReleasePtr()
	y, err := yaml.Marshal(r)
	if err != nil {
		panic(err)
	}
	return y
}

// DieReleaseRawExtension returns the resource managed by the die as an raw extension. Panics on error.
func (d *ImageSyncStatusDie) DieReleaseRawExtension() runtime.RawExtension {
	j := d.DieReleaseJSON()
	raw := runtime.RawExtension{}
	if err := json.Unmarshal(j, &raw); err != nil {
		panic(err)
	}
	return raw
}

// DieStamp returns a new die with the resource passed to the callback function. The resource is mutable.
func (d *ImageSyncStatusDie) DieStamp(fn func(r *syncv1alpha1.ImageSyncStatus)) *ImageSyncStatusDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

// Experimental: DieStampAt uses a JSON path (http://goessner.net/articles/JsonPath/) expression to stamp portions of the resource. The callback is invoked with each JSON path match. Panics if the callback function does not accept a single argument of the same type or a pointer to that type as found on the resource at the target location.
//
// Future iterations will improve type coercion from the resource to the callback argument.
func (d *ImageSyncStatusDie) DieStampAt(jp string, fn interface{}) *ImageSyncStatusDie {
	return d.DieStamp(func(r *syncv1alpha1.ImageSyncStatus) {
		if ni := reflectx.ValueOf(fn).Type().NumIn(); ni != 1 {
			panic(fmtx.Errorf("callback function must have 1 input parameters, found %d", ni))
		}
		if no := reflectx.ValueOf(fn).Type().NumOut(); no != 0 {
			panic(fmtx.Errorf("callback function must have 0 output parameters, found %d", no))
		}

		cp := jsonpath.New("")
		if err := cp.Parse(fmtx.Sprintf("{%s}", jp)); err != nil {
			panic(err)
		}
		cr, err := cp.FindResults(r)
		if err != nil {
			// errors are expected if a path is not found
			return
		}
		for _, cv := range cr[0] {
			arg0t := reflectx.ValueOf(fn).Type().In(0)

			var args []reflectx.Value
			if cv.Type().AssignableTo(arg0t) {
				args = []reflectx.Value{cv}
			} else if cv.CanAddr() && cv.Addr().Type().AssignableTo(arg0t) {
				args = []reflectx.Value{cv.Addr()}
			} else {
				panic(fmtx.Errorf("callback function must accept value of type %q, found type %q", cv.Type(), arg0t))
			}

			reflectx.ValueOf(fn).Call(args)
		}
	})
}

// DieWith returns a new die after passing the current die to the callback function. The passed die is mutable.
func (d *ImageSyncStatusDie) DieWith(fns ...func(d *ImageSyncStatusDie)) *ImageSyncStatusDie {
	nd := ImageSyncStatusBlank.DieFeed(d.DieRelease()).DieImmutable(false)
	for _, fn := range fns {
		if fn != nil {
			fn(nd)
		}
	}
	return d.DieFeed(nd.DieRelease())
}

// DeepCopy returns a new die with equivalent state. Useful for snapshotting a mutable die.
func (d *ImageSyncStatusDie) DeepCopy() *ImageSyncStatusDie {
	r := *d.r.DeepCopy()
	return &ImageSyncStatusDie{
		mutable: d.mutable,
		r:       r,
	}
}

func (d *ImageSyncStatusDie) Status(v apis.Status) *ImageSyncStatusDie {
	return d.DieStamp(func(r *syncv1alpha1.ImageSyncStatus) {
		r.Status = v
	})
}

// URL is the destination link for the latest Artifact.
func (d *ImageSyncStatusDie) URL(v string) *ImageSyncStatusDie {
	return d.DieStamp(func(r *syncv1alpha1.ImageSyncStatus) {
		r.URL = v
	})
}

// LastSyncTime to the destination repository
func (d *ImageSyncStatusDie) LastSyncTime(v metav1.Time) *ImageSyncStatusDie {
	return d.DieStamp(func(r *syncv1alpha1.ImageSyncStatus) {
		r.LastSyncTime = v
	})
}

var ImageBlank = (&ImageDie{}).DieFeed(syncv1alpha1.Image{})

type ImageDie struct {
	mutable bool
	r       syncv1alpha1.Image
}

// DieImmutable returns a new die for the current die's state that is either mutable (`false`) or immutable (`true`).
func (d *ImageDie) DieImmutable(immutable bool) *ImageDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy()
	d.mutable = !immutable
	return d
}

// DieFeed returns a new die with the provided resource.
func (d *ImageDie) DieFeed(r syncv1alpha1.Image) *ImageDie {
	if d.mutable {
		d.r = r
		return d
	}
	return &ImageDie{
		mutable: d.mutable,
		r:       r,
	}
}

// DieFeedPtr returns a new die with the provided resource pointer. If the resource is nil, the empty value is used instead.
func (d *ImageDie) DieFeedPtr(r *syncv1alpha1.Image) *ImageDie {
	if r == nil {
		r = &syncv1alpha1.Image{}
	}
	return d.DieFeed(*r)
}

// DieFeedJSON returns a new die with the provided JSON. Panics on error.
func (d *ImageDie) DieFeedJSON(j []byte) *ImageDie {
	r := syncv1alpha1.Image{}
	if err := json.Unmarshal(j, &r); err != nil {
		panic(err)
	}
	return d.DieFeed(r)
}

// DieFeedYAML returns a new die with the provided YAML. Panics on error.
func (d *ImageDie) DieFeedYAML(y []byte) *ImageDie {
	r := syncv1alpha1.Image{}
	if err := yaml.Unmarshal(y, &r); err != nil {
		panic(err)
	}
	return d.DieFeed(r)
}

// DieFeedYAMLFile returns a new die loading YAML from a file path. Panics on error.
func (d *ImageDie) DieFeedYAMLFile(name string) *ImageDie {
	y, err := osx.ReadFile(name)
	if err != nil {
		panic(err)
	}
	return d.DieFeedYAML(y)
}

// DieFeedRawExtension returns the resource managed by the die as an raw extension. Panics on error.
func (d *ImageDie) DieFeedRawExtension(raw runtime.RawExtension) *ImageDie {
	j, err := json.Marshal(raw)
	if err != nil {
		panic(err)
	}
	return d.DieFeedJSON(j)
}

// DieRelease returns the resource managed by the die.
func (d *ImageDie) DieRelease() syncv1alpha1.Image {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

// DieReleasePtr returns a pointer to the resource managed by the die.
func (d *ImageDie) DieReleasePtr() *syncv1alpha1.Image {
	r := d.DieRelease()
	return &r
}

// DieReleaseJSON returns the resource managed by the die as JSON. Panics on error.
func (d *ImageDie) DieReleaseJSON() []byte {
	r := d.DieReleasePtr()
	j, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	return j
}

// DieReleaseYAML returns the resource managed by the die as YAML. Panics on error.
func (d *ImageDie) DieReleaseYAML() []byte {
	r := d.DieReleasePtr()
	y, err := yaml.Marshal(r)
	if err != nil {
		panic(err)
	}
	return y
}

// DieReleaseRawExtension returns the resource managed by the die as an raw extension. Panics on error.
func (d *ImageDie) DieReleaseRawExtension() runtime.RawExtension {
	j := d.DieReleaseJSON()
	raw := runtime.RawExtension{}
	if err := json.Unmarshal(j, &raw); err != nil {
		panic(err)
	}
	return raw
}

// DieStamp returns a new die with the resource passed to the callback function. The resource is mutable.
func (d *ImageDie) DieStamp(fn func(r *syncv1alpha1.Image)) *ImageDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

// Experimental: DieStampAt uses a JSON path (http://goessner.net/articles/JsonPath/) expression to stamp portions of the resource. The callback is invoked with each JSON path match. Panics if the callback function does not accept a single argument of the same type or a pointer to that type as found on the resource at the target location.
//
// Future iterations will improve type coercion from the resource to the callback argument.
func (d *ImageDie) DieStampAt(jp string, fn interface{}) *ImageDie {
	return d.DieStamp(func(r *syncv1alpha1.Image) {
		if ni := reflectx.ValueOf(fn).Type().NumIn(); ni != 1 {
			panic(fmtx.Errorf("callback function must have 1 input parameters, found %d", ni))
		}
		if no := reflectx.ValueOf(fn).Type().NumOut(); no != 0 {
			panic(fmtx.Errorf("callback function must have 0 output parameters, found %d", no))
		}

		cp := jsonpath.New("")
		if err := cp.Parse(fmtx.Sprintf("{%s}", jp)); err != nil {
			panic(err)
		}
		cr, err := cp.FindResults(r)
		if err != nil {
			// errors are expected if a path is not found
			return
		}
		for _, cv := range cr[0] {
			arg0t := reflectx.ValueOf(fn).Type().In(0)

			var args []reflectx.Value
			if cv.Type().AssignableTo(arg0t) {
				args = []reflectx.Value{cv}
			} else if cv.CanAddr() && cv.Addr().Type().AssignableTo(arg0t) {
				args = []reflectx.Value{cv.Addr()}
			} else {
				panic(fmtx.Errorf("callback function must accept value of type %q, found type %q", cv.Type(), arg0t))
			}

			reflectx.ValueOf(fn).Call(args)
		}
	})
}

// DieWith returns a new die after passing the current die to the callback function. The passed die is mutable.
func (d *ImageDie) DieWith(fns ...func(d *ImageDie)) *ImageDie {
	nd := ImageBlank.DieFeed(d.DieRelease()).DieImmutable(false)
	for _, fn := range fns {
		if fn != nil {
			fn(nd)
		}
	}
	return d.DieFeed(nd.DieRelease())
}

// DeepCopy returns a new die with equivalent state. Useful for snapshotting a mutable die.
func (d *ImageDie) DeepCopy() *ImageDie {
	r := *d.r.DeepCopy()
	return &ImageDie{
		mutable: d.mutable,
		r:       r,
	}
}

// Image is a reference to an image in a remote repository
func (d *ImageDie) Image(v string) *ImageDie {
	return d.DieStamp(func(r *syncv1alpha1.Image) {
		r.Image = v
	})
}

// SecretRef contains the names of the Kubernetes Secrets containing registry login
//
// information to resolve image metadata.
func (d *ImageDie) SecretRef(v ...corev1.LocalObjectReference) *ImageDie {
	return d.DieStamp(func(r *syncv1alpha1.Image) {
		r.SecretRef = v
	})
}

// ServiceAccountName is the name of the Kubernetes ServiceAccount used to authenticate
//
// the image pull if the service account has attached pull secrets. For more information:
//
// https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/#add-imagepullsecrets-to-a-service-account
func (d *ImageDie) ServiceAccountName(v string) *ImageDie {
	return d.DieStamp(func(r *syncv1alpha1.Image) {
		r.ServiceAccountName = v
	})
}

// Insecure allows connecting to a non-TLS HTTP container registry.
func (d *ImageDie) Insecure(v bool) *ImageDie {
	return d.DieStamp(func(r *syncv1alpha1.Image) {
		r.Insecure = v
	})
}

// IsBundleImage allows synchronizing bundle images.
func (d *ImageDie) IsBundleImage(v bool) *ImageDie {
	return d.DieStamp(func(r *syncv1alpha1.Image) {
		r.IsBundleImage = v
	})
}
