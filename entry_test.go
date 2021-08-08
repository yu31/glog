package glog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestEntry_Byte_WithText(t *testing.T) {
	var eb bytes.Buffer
	var b bytes.Buffer
	l := NewDefault().WithExporter(MatchExporter(&b, nil)).WithErrorOutput(&eb)

	l.Info().Byte("key", '=').Fire()

	require.Greater(t, b.Len(), 0)
	s := string(b.Bytes()[0 : b.Len()-1])

	// parse fields
	fields := strings.Split(s, " ")
	require.Greater(t, len(fields), 0)
	fields = fields[len(fields)-1:]
	require.Equal(t, len(fields), 1)

	info := fields[0]
	require.True(t, strings.HasPrefix(info, "key="), "%q", info)
	require.False(t, strings.Contains(info, " "), "%q", info)

	info = strings.TrimPrefix(info, "key=")
	require.Greater(t, len(info), 0)

	i, err := strconv.ParseInt(info, 10, 64)
	require.Nil(t, err, "%q", info)
	require.Equal(t, byte(i), byte('='))

	require.Equal(t, eb.Len(), 0)
}

func TestEntry_Byte_WithJSON(t *testing.T) {
	var eb bytes.Buffer
	var b bytes.Buffer
	l := NewDefault().WithExporter(MatchExporter(&b, nil)).WithEncoderFunc(JSONEncoder).WithErrorOutput(&eb)

	l.Info().Byte("key", '=').Fire()

	require.Greater(t, b.Len(), 0)
	c := b.Bytes()[:b.Len()-1]

	m := make(map[string]interface{})
	err := json.Unmarshal(c, &m)
	require.Nil(t, err, "%q", string(c))

	n, ok := m["key"].(float64)
	require.True(t, ok, "%+q", m["key"])
	require.Equal(t, byte(n), byte('='), "%q", m["key"])

	require.Equal(t, eb.Len(), 0)
}

func TestEntry_Duration_WithText(t *testing.T) {
	var eb bytes.Buffer
	var b bytes.Buffer
	l := NewDefault().WithExporter(MatchExporter(&b, nil)).WithErrorOutput(&eb)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	d := time.Second + time.Duration(r.Int63n(int64(time.Second)))

	l.Info().
		Nanosecond("Nanosecond", d).
		Microsecond("Microsecond", d).
		Millisecond("Millisecond", d).
		Second("Second", d).
		Minute("Minute", d).
		Hour("Hour", d).
		Fire()

	require.Greater(t, b.Len(), 0)
	s := string(b.Bytes()[0 : b.Len()-1])
	fmt.Println(s)

	// parse fields
	fields := strings.Split(s, " ")
	require.Greater(t, len(fields), 6)
	fields = fields[len(fields)-6:]
	require.Equal(t, len(fields), 6)

	// Nanosecond
	require.True(t, strings.HasPrefix(fields[0], "Nanosecond="), "%q", fields[0])
	require.True(t, strings.HasSuffix(fields[0], nanoSuffix), "%q", fields[0])
	require.False(t, strings.Contains(fields[0], " "), "%q", fields[0])

	// Microsecond
	require.True(t, strings.HasPrefix(fields[1], "Microsecond="), "%q", fields[1])
	require.True(t, strings.HasSuffix(fields[1], microSuffix), "%q", fields[1])
	require.False(t, strings.Contains(fields[1], " "), "%q", fields[1])

	// Millisecond
	require.True(t, strings.HasPrefix(fields[2], "Millisecond="), "%q", fields[2])
	require.True(t, strings.HasSuffix(fields[2], milliSuffix), "%q", fields[2])
	require.False(t, strings.Contains(fields[2], " "), "%q", fields[2])

	// Second
	require.True(t, strings.HasPrefix(fields[3], "Second="), "%q", fields[3])
	require.True(t, strings.HasSuffix(fields[3], secondSuffix), "%q", fields[3])
	require.False(t, strings.Contains(fields[3], " "), "%q", fields[3])

	// Minute
	require.True(t, strings.HasPrefix(fields[4], "Minute="), "%q", fields[4])
	require.True(t, strings.HasSuffix(fields[4], minuteSuffix), "%q", fields[4])
	require.False(t, strings.Contains(fields[4], " "), "%q", fields[4])

	// Hour
	require.True(t, strings.HasPrefix(fields[5], "Hour="), "%q", fields[5])
	require.True(t, strings.HasSuffix(fields[5], hourSuffix), "%q", fields[5])
	require.False(t, strings.Contains(fields[5], " "), "%q", fields[5])

	require.Equal(t, eb.Len(), 0)
}

func TestEntry_Duration_WithJSON(t *testing.T) {
	var eb bytes.Buffer
	var b bytes.Buffer
	l := NewDefault().WithExporter(MatchExporter(&b, nil)).WithEncoderFunc(JSONEncoder).WithErrorOutput(&eb)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	d := time.Second + time.Duration(r.Int63n(int64(time.Second)))

	l.Info().
		Msg("Test for JSON format").
		Nanosecond("Nanosecond", d).
		Microsecond("Microsecond", d).
		Millisecond("Millisecond", d).
		Second("Second", d).
		Minute("Minute", d).
		Hour("Hour", d).
		Fire()

	require.Greater(t, b.Len(), 0)

	c := b.Bytes()[:b.Len()-1]
	fmt.Println(string(c))

	var m map[string]string
	err := json.Unmarshal(c, &m)
	require.Nil(t, err, "%+q", string(c))

	// Nanosecond
	require.NotEqual(t, len(m["Nanosecond"]), 0, "%q", m)
	require.True(t, strings.HasSuffix(m["Nanosecond"], nanoSuffix), "%q", m["Nanosecond"])

	// Microsecond
	require.NotEqual(t, len(m["Microsecond"]), 0, "%q", m)
	require.True(t, strings.HasSuffix(m["Microsecond"], microSuffix), "%q", m["Microsecond"])

	// Millisecond
	require.NotEqual(t, len(m["Millisecond"]), 0, "%q", m)
	require.True(t, strings.HasSuffix(m["Millisecond"], milliSuffix), "%q", m["Millisecond"])

	// Second
	require.NotEqual(t, len(m["Second"]), 0, "%q", m)
	require.True(t, strings.HasSuffix(m["Second"], secondSuffix), "%q", m["Second"])

	// Minute
	require.NotEqual(t, len(m["Minute"]), 0, "%q", m)
	require.True(t, strings.HasSuffix(m["Minute"], minuteSuffix), "%q", m["Minute"])

	// Hour
	require.NotEqual(t, len(m["Hour"]), 0, "%q", m)
	require.True(t, strings.HasSuffix(m["Hour"], hourSuffix), "%q", m["Hour"])

	require.Equal(t, eb.Len(), 0)
}

func TestEntry_Raw_WithText(t *testing.T) {
	var b bytes.Buffer
	l := NewDefault().WithExporter(MatchExporter(&b, nil))

	data := struct {
		Name   string
		Number int64
	}{
		Name:   "n1",
		Number: 1,
	}
	bd, err := json.Marshal(&data)
	require.Nil(t, err)

	l.Info().String("key1", string(bd)).Fire()
	require.Contains(t, b.String(), "\\")

	b.Reset()
	l.Info().RawBytes("key2", bd).RawString("key3", string(bd)).Fire()
	require.NotContains(t, b.String(), "\\")
}

func TestEntry_Raw_WithJSON(t *testing.T) {
	var b bytes.Buffer
	l := NewDefault().WithExporter(MatchExporter(&b, nil)).WithEncoderFunc(JSONEncoder)

	data := struct {
		Name   string
		Number int64
	}{
		Name:   "n1",
		Number: 1,
	}

	bd, err := json.Marshal(&data)
	require.Nil(t, err)

	l.Info().String("key1", string(bd)).RawBytes("key2", bd).RawString("key3", string(bd)).Fire()

	require.Greater(t, b.Len(), 0)

	b1 := b.Bytes()[0 : b.Len()-1]

	var m map[string]interface{}
	err = json.Unmarshal(b1, &m)
	require.Nil(t, err)

	// Expected "key1" is string
	k1, ok := m["key1"]
	require.True(t, ok)
	_, ok = k1.(string)
	require.True(t, ok)

	// Expected "key2" is map
	k2, ok := m["key2"]
	require.True(t, ok)
	_, ok = k2.(map[string]interface{})
	require.True(t, ok)

	// Expected "key3" is map
	k3, ok := m["key3"]
	require.True(t, ok)
	_, ok = k3.(map[string]interface{})
	require.True(t, ok)
}

func TestEntry_Error(t *testing.T) {
	l := NewDefault().WithExporter(MatchExporter(ioutil.Discard, nil))
	var err error

	require.NotPanics(t, func() {
		l.Error().Error("error", err).Fire()
	})
}

func TestEntry_Errors(t *testing.T) {
	l := NewDefault().WithExporter(MatchExporter(ioutil.Discard, nil))
	var errors []error

	require.NotPanics(t, func() {
		l.Error().Errors("error", errors).Fire()
	})

	errors = append(errors, nil)

	require.NotPanics(t, func() {
		l.Error().Errors("error", errors).Fire()
	})

}
