package timestat

import "time"

type TimeStat struct {
	times []time.Duration
	cur   int
}

func New(size int) TimeStat {
	return TimeStat{
		times: make([]time.Duration, size),
		cur:   size - 1,
	}
}

func (t *TimeStat) Add(v time.Duration) {
	t.cur = (t.cur + 1) % len(t.times)
	t.times[t.cur] = v
}

func (t *TimeStat) Values() []time.Duration {
	r := make([]time.Duration, len(t.times))

	for i := 0; i < len(t.times); i++ {
		idx := (t.cur + 1 + i) % len(t.times)
		r[i] = t.times[idx]
	}

	return r
}

func (t *TimeStat) Len() int {
	return len(t.times)
}
