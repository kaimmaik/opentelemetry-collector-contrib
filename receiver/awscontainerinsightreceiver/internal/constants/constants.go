// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package constants // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awscontainerinsightreceiver/internal/constants"

import "time"

// DefaultCollectionInterval is the default collection interval for most receivers
const DefaultCollectionInterval = 60 * time.Second

// EFADefaultCollectionInterval is the default collection interval specifically for EFA metrics
const EFADefaultCollectionInterval = 20 * time.Second
