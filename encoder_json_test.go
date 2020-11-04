package glog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLoggerWithJSONEncoder(t *testing.T) {
	var b bytes.Buffer
	l := NewDefault().WithEncoderFunc(JSONEncoder).WithExecutor(MatchExecutor(&b, nil)).WithCaller(true)

	l.Info().
		Msg("test logger out").
		String("String", "Value").
		Strings("Strings", []string{"a", "b", "c"}).
		Byte("Byte", 'a').
		Bytes("Bytes", []byte("abc")).
		Int64("Int64", 64).
		Int64s("Int64s", []int64{641, 642, 643}).
		Complex128("Complex128", complex(1, 2)).
		Float64("Float64", 99.99).
		Strings("Strings", []string{"a", "b", "c"}).
		Bytes("Bytes", []byte("HEllO")).
		Bool("Bool", true).
		Duration("Duration", time.Duration(1)).
		Time("Time", time.Now(), time.RFC3339Nano).
		Any("Interface1", []string{"i1", "i2", "i3"}).
		Any("Interface2", nil).
		Fire()

	c := b.Bytes()
	c = c[:len(c)-1]
	fmt.Println(string(c))

	m := make(map[string]interface{})
	err := json.Unmarshal(c, &m)
	require.Nil(t, err)
}
