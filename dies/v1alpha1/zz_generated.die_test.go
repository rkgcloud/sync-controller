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
	testingx "testing"

	testing "reconciler.io/dies/testing"
)

func TestImageSyncDie_MissingMethods(t *testingx.T) {
	die := ImageSyncBlank
	ignore := []string{"TypeMeta", "ObjectMeta"}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for ImageSyncDie: %s", diff.List())
	}
}

func TestImageSyncSpecDie_MissingMethods(t *testingx.T) {
	die := ImageSyncSpecBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for ImageSyncSpecDie: %s", diff.List())
	}
}

func TestImageSyncStatusDie_MissingMethods(t *testingx.T) {
	die := ImageSyncStatusBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for ImageSyncStatusDie: %s", diff.List())
	}
}

func TestImageSourceDie_MissingMethods(t *testingx.T) {
	die := ImageSourceBlank
	ignore := []string{}
	diff := testing.DieFieldDiff(die).Delete(ignore...)
	if diff.Len() != 0 {
		t.Errorf("found missing fields for ImageSourceDie: %s", diff.List())
	}
}
