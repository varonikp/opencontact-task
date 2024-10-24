package httpserver

import (
	"time"
)

type Date time.Time

func (d *Date) String() string {
	return time.Time(*d).String()
}

func (d *Date) UnmarshalParam(param string) error {
	t, err := time.Parse(`2006-01-02`, param)
	if err != nil {
		return err
	}
	*d = Date(t)
	return nil
}
