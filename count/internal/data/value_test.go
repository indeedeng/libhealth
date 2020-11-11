package data

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Int(t *testing.T) {
	a := NewInt(1)
	require.Equal(t, "1", a.String())

	b := NewInt(2)
	c := NewInt(3)
	require.True(t, b.Less(c))
	require.True(t, a.Less(c))
	require.False(t, c.Less(b))
	require.False(t, c.Less(a))

	d := a.Add(c)
	require.Equal(t, "4", d.String())
	require.True(t, c.Less(d))
	require.False(t, d.Less(c))

	z := Zero(d)
	require.Equal(t, "0", z.String())
	require.True(t, z.Less(a))
}

func Test_Float(t *testing.T) {
	a := NewFloat(1.1)
	require.Equal(t, "1.100", a.String())

	b := NewFloat(2.2)
	c := NewFloat(3.3)

	require.True(t, b.Less(c))
	require.False(t, c.Less(b))

	d := a.Add(c)
	require.Equal(t, "4.400", d.String())
	require.True(t, c.Less(d))
	require.False(t, d.Less(c))

	z := Zero(d)
	require.Equal(t, "0.000", z.String())
	require.True(t, z.Less(a))
}

func Test_BadConverts(t *testing.T) {
	i := NewInt(1)
	f := NewFloat(2.0)

	require.Panics(t, func() { i.Less(f) })
	require.Panics(t, func() { f.Less(i) })
	require.Panics(t, func() { i.Add(f) })
	require.Panics(t, func() { f.Add(i) })
}
