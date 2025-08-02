// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package multierr_test

import (
	"errors"
	"testing"

	"github.com/kopexa-grc/x/multierr"
	"github.com/stretchr/testify/assert"
)

// Static test errors to avoid dynamic error creation
var (
	errTest1 = errors.New("1")
)

func TestMultiErr(t *testing.T) {
	t.Run("add nil errors", func(t *testing.T) {
		var e multierr.Errors
		e.Add(nil)
		e.Add(nil, nil, nil)
		assert.Nil(t, e.Deduplicate())
	})

	t.Run("add mixed errors", func(t *testing.T) {
		var e multierr.Errors
		e.Add(errTest1, nil, errTest1)
		var b multierr.Errors
		b.Add(errTest1)
		assert.Equal(t, b.Deduplicate(), e.Deduplicate())
	})

	t.Run("test nil error deduplicate", func(t *testing.T) {
		var e multierr.Errors
		err := e.Deduplicate()
		assert.Nil(t, err)
	})
}
