package glog

import "time"

type timeArray []time.Time

func (vv timeArray) MarshalGLogArray(ae ArrayEncoder) error {
	for i := range vv {
		ae.AppendTime(vv[i], time.RFC3339Nano)
	}
	return nil
}

type info struct {
	Name  string    `json:"name"`
	Sex   string    `json:"sex"`
	Age   int       `json:"age"`
	Times timeArray `json:"times"`
}

func (o *info) MarshalGLogObject(oe ObjectEncoder) error {
	oe.AddString("name", o.Name)
	oe.AddString("sex", o.Sex)
	oe.AddInt64("age", int64(o.Age))
	err := oe.AddArray("times", o.Times)
	return err
}

type infos []*info

func (vv infos) MarshalGLogArray(ae ArrayEncoder) error {
	for i := range vv {
		if err := ae.AppendObject(vv[i]); err != nil {
			return err
		}
	}
	return nil
}
