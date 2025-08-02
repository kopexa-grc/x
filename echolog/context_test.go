// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package echolog_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/kopexa-grc/x/echolog"
	"github.com/stretchr/testify/assert"
)

func TestCtx(t *testing.T) {
	b := &bytes.Buffer{}
	l := echolog.New(b)
	zerologger := l.Unwrap()
	ctx := l.WithContext(context.Background())

	assert.Equal(t, echolog.Ctx(ctx), &zerologger)
}
