// Copyright 2026 Red Hat, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build windows

package system

import (
	"path/filepath"
	"strings"
)

func HostPathToMachinePath(path string) string {
	posixPath := filepath.ToSlash(path)

	// If the path starts with a drive letter (e.g., C:), convert it to a POSIX-like path
	if len(posixPath) > 2 && posixPath[1] == ':' {
		posixPath = "/mnt/" + strings.ToLower(string(posixPath[0])) + "/" + posixPath[2:]
	}

	return posixPath
}

func MachinePathToHostPath(path string) string {
	hostPath := filepath.FromSlash(path)
	if strings.HasPrefix(path, "/mnt/") {
		hostPath = strings.ToUpper(string(hostPath[5])) + ":" + hostPath[6:]
	}
	return hostPath
}
