/*
 * Copyright 2019 gosoon.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package pointer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInt32Ptr(t *testing.T) {
	testCases := []struct {
		int32Value int32
	}{
		{
			int32Value: int32(1),
		},
		{
			int32Value: int32(-1),
		},
	}

	for _, test := range testCases {
		ptr := Int32Ptr(test.int32Value)
		if !assert.Equal(t, &test.int32Value, ptr) {
			t.Fatalf("expected: %v but get %v", &test.int32Value, ptr)
		}
	}
}

func TestInt64Ptr(t *testing.T) {
	testCases := []struct {
		int64Value int64
	}{
		{
			int64Value: int64(1),
		},
		{
			int64Value: int64(-1),
		},
	}

	for _, test := range testCases {
		ptr := Int64Ptr(test.int64Value)
		if !assert.Equal(t, &test.int64Value, ptr) {
			t.Fatalf("expected: %v but get %v", &test.int64Value, ptr)
		}
	}
}

func TestBoolPtr(t *testing.T) {
	testCases := []struct {
		boolValue bool
	}{
		{
			boolValue: true,
		},
		{
			boolValue: false,
		},
	}

	for _, test := range testCases {
		ptr := BoolPtr(test.boolValue)
		if !assert.Equal(t, &test.boolValue, ptr) {
			t.Fatalf("expected: %v but get %v", &test.boolValue, ptr)
		}
	}
}
