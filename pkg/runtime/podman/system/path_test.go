// Copyright 2026 Red Hat, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build !windows

package system

import "testing"

func TestHostPathToMachinePath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "absolute path",
			input: "/home/user/project",
			want:  "/home/user/project",
		},
		{
			name:  "root path",
			input: "/",
			want:  "/",
		},
		{
			name:  "relative path",
			input: "relative/path",
			want:  "relative/path",
		},
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := HostPathToMachinePath(tt.input)
			if got != tt.want {
				t.Errorf("HostPathToMachinePath(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestMachinePathToHostPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "absolute path",
			input: "/home/user/project",
			want:  "/home/user/project",
		},
		{
			name:  "root path",
			input: "/",
			want:  "/",
		},
		{
			name:  "relative path",
			input: "relative/path",
			want:  "relative/path",
		},
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := MachinePathToHostPath(tt.input)
			if got != tt.want {
				t.Errorf("MachinePathToHostPath(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
