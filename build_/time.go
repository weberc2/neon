package build

import "time"

var layout = "2006-01-02"

type Time time.Time

func (t Time) After(t2 Time) bool {
	return time.Time(t).After(time.Time(t2))
}

func (t *Time) UnmarshalJSON(p []byte) error {
	t2, err := time.Parse(layout, string(p[1:len(p)-1]))
	if err != nil {
		return err
	}
	*t = Time(t2)
	return nil
}

func (t Time) String() string {
	return time.Time(t).Format(layout)
}
