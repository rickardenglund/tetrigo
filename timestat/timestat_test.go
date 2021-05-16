package timestat_test

import (
	"testing"
	"time"

	"github.com/rickardenglund/tetrigo/timestat"

	"github.com/stretchr/testify/require"
)

func TestAddAndGet(t *testing.T) {
	const n = 3

	var (
		i      time.Duration
		buffer = timestat.New(n)
	)

	for i = 0; i < n; i++ {
		buffer.Add(time.Second * i)
	}

	require.Equal(t, []time.Duration{time.Second * 0, time.Second * 1, time.Second * 2}, buffer.Values())

	buffer.Add(time.Second * 4)

	require.Equal(t, []time.Duration{time.Second * 1, time.Second * 2, time.Second * 4}, buffer.Values())
}

func TestGetWithoutAdd(t *testing.T) {
	buffer := timestat.New(2)
	require.Equal(t, []time.Duration{0, 0}, buffer.Values())
}
