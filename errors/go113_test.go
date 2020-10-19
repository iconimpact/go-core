/*
   Copyright 2020 iconmobile GmbH

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIs(t *testing.T) {
	err := E(Unprocessable, "test")

	type args struct {
		err    error
		target error
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "with kind",
			args: args{
				err:    E(err, Internal),
				target: err,
			},
			want: true,
		},
		{
			name: "std errors compatibility",
			args: args{
				err:    fmt.Errorf("wrap it: %w", err),
				target: err,
			},
			want: true,
		},
		{
			name: "std match custom",
			args: args{
				err:    fmt.Errorf("std error"),
				target: err,
			},
			want: false,
		},
		{
			name: "custom match std",
			args: args{
				err:    err,
				target: fmt.Errorf("std error"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Is(tt.args.err, tt.args.target)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAs(t *testing.T) {
	err := E(Unprocessable, "test")

	type args struct {
		err    error
		target interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "with kind",
			args: args{
				err:    E(err, Internal),
				target: &err,
			},
			want: true,
		},
		{
			name: "std errors compatibility",
			args: args{
				err:    fmt.Errorf("wrap it: %w", err),
				target: &err,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := As(tt.args.err, tt.args.target)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUnwrap(t *testing.T) {
	err := E(Unprocessable, "test")

	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "with kind",
			args: args{err: E(err, Internal)},
			want: err,
		},
		{
			name: "std errors compatibility",
			args: args{err: fmt.Errorf("wrap: %w", err)},
			want: err,
		},
		{
			name: "std match custom",
			args: args{err: fmt.Errorf("std error")},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Unwrap(tt.args.err)
			assert.Equal(t, tt.want, err)
		})
	}
}
