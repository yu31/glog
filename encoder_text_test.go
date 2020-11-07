package glog

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTextEncoder_AddArray(t *testing.T) {
	enc := TextEncoder()
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

	s := string(enc.Bytes())

	require.Contains(t, s, "time1=")
	require.Contains(t, s, " time2=")
}

func TestTextEncoder_AddObject(t *testing.T) {
	enc := TextEncoder()
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

	s := string(enc.Bytes())
	require.Equal(t, strings.Count(s, "infos="), 1)
	require.Equal(t, strings.Count(s, "times="), 3)
	require.Equal(t, strings.Count(s, "name="), 3)
	require.Equal(t, strings.Count(s, "sex="), 3)
	require.Equal(t, strings.Count(s, "age="), 3)
}
