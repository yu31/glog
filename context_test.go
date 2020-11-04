package glog

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWithContext(t *testing.T) {
	l := NewDefault()
	ctx := WithContext(context.Background(), l)
	nl := FromContext(ctx)
	require.True(t, reflect.DeepEqual(l, nl))

	// update logger
	l = l.WithLevel(InfoLevel)
	ctx = WithContext(ctx, l)
	nl = FromContext(ctx)
	require.True(t, reflect.DeepEqual(l, nl))

	// new logger
	l1 := NewDefault()
	ctx = WithContext(ctx, l1)
	nl = FromContext(ctx)
	require.False(t, reflect.DeepEqual(l, nl))
}

func TestFromContext(t *testing.T) {
	// get a default logger
	l := FromContext(context.Background())
	require.NotNil(t, l)
}
