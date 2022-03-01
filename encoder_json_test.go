package glog

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestJsonEncoder_AddArray(t *testing.T) {
	enc := JSONEncoder()
	defer func() {
		_ = enc.Close()
	}()

	n := 3
	ts := make(timeArray, 0, n)
	for i := 0; i < n; i++ {
		ts = append(ts, time.Now())
	}

	enc.AddBeginMarker()
	require.Nil(t, enc.AddArray("time1", ts))
	require.Nil(t, enc.AddArray("time2", ts))
	enc.AddEndMarker()

	var err error

	var m1 map[string][]time.Time
	err = json.Unmarshal(enc.Bytes(), &m1)
	require.Nil(t, err)
	require.Equal(t, len(m1["time1"]), n)
	require.Equal(t, len(m1["time2"]), n)

	var m2 map[string][]string
	err = json.Unmarshal(enc.Bytes(), &m2)
	require.Nil(t, err)
	require.Equal(t, len(m2["time1"]), n)
	require.Equal(t, len(m2["time2"]), n)
}

func TestJsonEncoder_AddObject(t *testing.T) {
	enc := JSONEncoder()
	defer func() {
		_ = enc.Close()
	}()

	var infos infos
	infos = append(infos, &info{
		Name:  "aa",
		Sex:   "man",
		Age:   999,
		Times: timeArray{time.Now(), time.Now()},
	})
	infos = append(infos, &info{
		Name:  "bb",
		Sex:   "man",
		Age:   999,
		Times: timeArray{time.Now(), time.Now()},
	})
	infos = append(infos, &info{
		Name:  "cc",
		Sex:   "man",
		Age:   999,
		Times: timeArray{time.Now(), time.Now()},
	})

	enc.AddBeginMarker()
	require.Nil(t, enc.AddArray("infos", infos))
	enc.AddEndMarker()

	var err error
	var m1 map[string][]*info

	err = json.Unmarshal(enc.Bytes(), &m1)
	require.Nil(t, err)
	require.Equal(t, len(m1), 1)
	require.Equal(t, len(m1["infos"]), len(infos))
}

func TestJSONEncoder_AddFields(t *testing.T) {
	var eb bytes.Buffer
	var b bytes.Buffer
	l := NewDefault().WithEncoderFunc(JSONEncoder).WithExporter(StandardExporter(&b)).
		WithCaller(true).WithErrorOutput(&eb)

	l.WithFields().AddString("rid", "xxxxxx01")
	l.WithFields().AddString("tid", "yyyyyy02")
	l.WithFields().AddString("xid", "zzzzzz03")

	l.Info().
		Msg("test logger out").
		String("String", "Value").
		Fire()

	require.Greater(t, b.Len(), 0)
	c := b.Bytes()[:b.Len()-1]

	m := make(map[string]interface{})
	err := json.Unmarshal(c, &m)
	require.Nil(t, err, "%q", string(c))
	require.Equal(t, "xxxxxx01", m["rid"])
	require.Equal(t, "yyyyyy02", m["tid"])
	require.Equal(t, "zzzzzz03", m["xid"])
	require.Equal(t, "Value", m["String"])
	require.Equal(t, "test logger out", m["message"])

	require.Equal(t, eb.Len(), 0)
}
