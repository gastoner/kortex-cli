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

package podman

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/openkaiden/kdn/pkg/logger"
	"github.com/openkaiden/kdn/pkg/runtime"
	"github.com/openkaiden/kdn/pkg/steplogger"
)

const (
	postgresMaxRetries    = 30
	postgresRetryInterval = 2 * time.Second
)

// Start starts all containers in the workspace pod.
// Postgres is started first and verified ready before starting the remaining
// containers so that onecli can connect to the database immediately.
func (p *podmanRuntime) Start(ctx context.Context, id string) (runtime.RuntimeInfo, error) {
	stepLogger := steplogger.FromContext(ctx)
	defer stepLogger.Complete()

	if id == "" {
		return runtime.RuntimeInfo{}, fmt.Errorf("%w: container ID is required", runtime.ErrInvalidParams)
	}

	// Resolve the pod name from the stored mapping
	podName, err := p.readPodName(id)
	if err != nil {
		return runtime.RuntimeInfo{}, fmt.Errorf("failed to resolve pod name: %w", err)
	}

	l := logger.FromContext(ctx)

	// Start the postgres container first so it is accepting connections
	// before onecli attempts its database migration.
	postgresContainer := podName + "-postgres"
	stepLogger.Start("Starting postgres", "Postgres started")
	if err := p.executor.Run(ctx, l.Stdout(), l.Stderr(), "start", postgresContainer); err != nil {
		stepLogger.Fail(err)
		return runtime.RuntimeInfo{}, fmt.Errorf("failed to start postgres container: %w", err)
	}

	// Wait until postgres is accepting connections
	stepLogger.Start("Waiting for postgres to be ready", "Postgres is ready")
	if err := p.waitForPostgres(ctx, podName); err != nil {
		stepLogger.Fail(err)
		return runtime.RuntimeInfo{}, fmt.Errorf("postgres did not become ready: %w", err)
	}

	// Start the rest of the pod (onecli + workspace); already-running
	// containers (postgres) are left untouched by pod start.
	stepLogger.Start(fmt.Sprintf("Starting pod: %s", podName), "Pod started")
	if err := p.executor.Run(ctx, l.Stdout(), l.Stderr(), "pod", "start", podName); err != nil {
		stepLogger.Fail(err)
		return runtime.RuntimeInfo{}, fmt.Errorf("failed to start pod: %w", err)
	}

	// Verify workspace container status
	stepLogger.Start("Verifying container status", "Container status verified")
	info, err := p.getContainerInfo(ctx, id)
	if err != nil {
		stepLogger.Fail(err)
		return runtime.RuntimeInfo{}, fmt.Errorf("failed to get container info after start: %w", err)
	}

	return info, nil
}

// waitForPostgres polls the postgres container inside the pod until pg_isready succeeds.
// The postgres container name follows the podman kube play convention: <podName>-postgres.
func (p *podmanRuntime) waitForPostgres(ctx context.Context, podName string) error {
	containerName := podName + "-postgres"
	var lastErr error
	for range postgresMaxRetries {
		_, err := p.executor.Output(ctx, io.Discard,
			"exec", containerName, "pg_isready", "-U", "onecli")
		if err == nil {
			return nil
		}
		lastErr = err

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(postgresRetryInterval):
		}
	}
	return fmt.Errorf("postgres not ready after %d retries: %w", postgresMaxRetries, lastErr)
}
