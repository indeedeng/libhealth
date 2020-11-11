package data

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func testBucketSize() (Ticker, *int) {
	tickTime := 0
	return func() int {
		return tickTime
	}, &tickTime
}

func Test_IntBuckets2(t *testing.T) {
	size, tickTime := testBucketSize()
	buckets := NewBuckets(size, 2, NewInt(0))

	buckets.Increment(NewInt(4))
	{
		require.Equal(t, NewInt(4), buckets.Sum())
		require.Equal(t, "[4, 0]", buckets.String())
	}

	*tickTime = 1
	buckets.Increment(NewInt(1))
	{
		require.Equal(t, NewInt(5), buckets.Sum())
		require.Equal(t, "[1, 4]", buckets.String())
	}

	*tickTime = 2
	buckets.Increment(NewInt(7))
	{
		require.Equal(t, NewInt(8), buckets.Sum())
		require.Equal(t, "[7, 1]", buckets.String())
	}

	*tickTime = 6
	buckets.Increment(NewInt(11))
	{
		require.Equal(t, NewInt(11), buckets.Sum())
		require.Equal(t, "[11, 0]", buckets.String())
	}
}

func Test_IntBuckets4(t *testing.T) {
	size, tickTime := testBucketSize()
	buckets := NewBuckets(size, 4, NewInt(0))

	buckets.Increment(NewInt(4))
	{
		require.Equal(t, NewInt(4), buckets.Sum())
		require.Equal(t, "[4, 0, 0, 0]", buckets.String())
	}

	*tickTime = 1
	buckets.Increment(NewInt(1))
	{
		require.Equal(t, NewInt(5), buckets.Sum())
		require.Equal(t, "[1, 4, 0, 0]", buckets.String())
	}

	*tickTime = 2
	buckets.Increment(NewInt(7))
	{
		require.Equal(t, NewInt(12), buckets.Sum())
		require.Equal(t, "[7, 1, 4, 0]", buckets.String())
	}

	*tickTime = 3
	buckets.Increment(NewInt(13))
	{
		require.Equal(t, NewInt(25), buckets.Sum())
		require.Equal(t, "[13, 7, 1, 4]", buckets.String())
	}

	*tickTime = 4
	buckets.Increment(NewInt(5))
	{
		require.Equal(t, NewInt(26), buckets.Sum())
		require.Equal(t, "[5, 13, 7, 1]", buckets.String())
	}

	*tickTime = 6
	buckets.Increment(NewInt(11))
	{
		require.Equal(t, NewInt(29), buckets.Sum())
		require.Equal(t, "[11, 0, 5, 13]", buckets.String())
	}

	*tickTime = 7
	buckets.Increment(NewInt(4))
	{
		require.Equal(t, NewInt(20), buckets.Sum())
		require.Equal(t, "[4, 11, 0, 5]", buckets.String())

		// oldest complete == 11
		require.Equal(t, true, buckets.Compare(GreaterEq, NewInt(10)))
		require.Equal(t, true, buckets.Compare(GreaterEq, NewInt(11)))
		require.Equal(t, false, buckets.Compare(GreaterEq, NewInt(12)))

		require.Equal(t, false, buckets.Compare(LessEq, NewInt(10)))
		require.Equal(t, true, buckets.Compare(LessEq, NewInt(11)))
		require.Equal(t, true, buckets.Compare(LessEq, NewInt(12)))
	}

	*tickTime = 10
	buckets.Increment(NewInt(2))
	{
		require.Equal(t, NewInt(6), buckets.Sum())
		require.Equal(t, "[2, 0, 0, 4]", buckets.String())

		// oldest complete == 0
		require.Equal(t, true, buckets.Compare(GreaterEq, NewInt(-1)))
		require.Equal(t, true, buckets.Compare(GreaterEq, NewInt(0)))
		require.Equal(t, false, buckets.Compare(GreaterEq, NewInt(1)))

		require.Equal(t, false, buckets.Compare(LessEq, NewInt(-1)))
		require.Equal(t, true, buckets.Compare(LessEq, NewInt(0)))
		require.Equal(t, true, buckets.Compare(LessEq, NewInt(1)))
	}

	*tickTime = 15
	buckets.Increment(NewInt(99))
	{
		require.Equal(t, NewInt(99), buckets.Sum())
		require.Equal(t, "[99, 0, 0, 0]", buckets.String())
	}

	// no tick
	buckets.Increment(NewInt(14))
	{
		require.Equal(t, NewInt(113), buckets.Sum())
		require.Equal(t, "[113, 0, 0, 0]", buckets.String())
	}
}

func Test_FloatBuckets2(t *testing.T) {
	size, tickTime := testBucketSize()
	buckets := NewBuckets(size, 2, NewFloat(0))

	buckets.Increment(NewFloat(4))
	{
		require.Equal(t, NewFloat(4), buckets.Sum())
		require.Equal(t, "[4.000, 0.000]", buckets.String())
	}

	*tickTime = 1
	buckets.Increment(NewFloat(1))
	{
		require.Equal(t, NewFloat(5), buckets.Sum())
		require.Equal(t, "[1.000, 4.000]", buckets.String())
	}

	*tickTime = 2
	buckets.Increment(NewFloat(7))
	{
		require.Equal(t, NewFloat(8), buckets.Sum())
		require.Equal(t, "[7.000, 1.000]", buckets.String())
	}

	*tickTime = 6
	buckets.Increment(NewFloat(11))
	{
		require.Equal(t, NewFloat(11), buckets.Sum())
		require.Equal(t, "[11.000, 0.000]", buckets.String())
	}
}

func Test_FloatBuckets4(t *testing.T) {
	size, tickTime := testBucketSize()
	buckets := NewBuckets(size, 4, NewFloat(0))

	buckets.Increment(NewFloat(4.4))
	{
		require.Equal(t, NewFloat(4.4), buckets.Sum())
		require.Equal(t, "[4.400, 0.000, 0.000, 0.000]", buckets.String())
	}

	*tickTime = 1
	buckets.Increment(NewFloat(1.1))
	{
		require.Equal(t, NewFloat(5.5), buckets.Sum())
		require.Equal(t, "[1.100, 4.400, 0.000, 0.000]", buckets.String())
	}

	*tickTime = 2
	buckets.Increment(NewFloat(7.7))
	{
		require.Equal(t, NewFloat(13.2), buckets.Sum())
		require.Equal(t, "[7.700, 1.100, 4.400, 0.000]", buckets.String())
	}

	*tickTime = 3
	buckets.Increment(NewFloat(13.13))
	{
		require.Equal(t, NewFloat(26.33), buckets.Sum())
		require.Equal(t, "[13.130, 7.700, 1.100, 4.400]", buckets.String())
	}

	*tickTime = 4
	buckets.Increment(NewFloat(5.5))
	{
		require.Equal(t, NewFloat(27.43), buckets.Sum())
		require.Equal(t, "[5.500, 13.130, 7.700, 1.100]", buckets.String())
	}

	*tickTime = 6
	buckets.Increment(NewFloat(-1.1))
	{
		require.Equal(t, NewFloat(17.53), buckets.Sum())
		require.Equal(t, "[-1.100, 0.000, 5.500, 13.130]", buckets.String())
	}

	*tickTime = 7
	buckets.Increment(NewFloat(4.4))
	{
		require.Equal(t, NewFloat(8.8), buckets.Sum())
		require.Equal(t, "[4.400, -1.100, 0.000, 5.500]", buckets.String())

		// oldest complete is -1.1
		require.Equal(t, true, buckets.Compare(GreaterEq, NewFloat(-1.2)))
		require.Equal(t, true, buckets.Compare(GreaterEq, NewFloat(-1.1)))
		require.Equal(t, false, buckets.Compare(GreaterEq, NewFloat(-1.0)))

		require.Equal(t, false, buckets.Compare(LessEq, NewFloat(-1.2)))
		require.Equal(t, true, buckets.Compare(LessEq, NewFloat(-1.1)))
		require.Equal(t, true, buckets.Compare(LessEq, NewFloat(-1.0)))
	}

	*tickTime = 10
	buckets.Increment(NewFloat(2.3))
	{
		require.Equal(t, NewFloat(6.7), buckets.Sum())
		require.Equal(t, "[2.300, 0.000, 0.000, 4.400]", buckets.String())
	}

	*tickTime = 15
	buckets.Increment(NewFloat(99.99))
	{
		require.Equal(t, NewFloat(99.99), buckets.Sum())
		require.Equal(t, "[99.990, 0.000, 0.000, 0.000]", buckets.String())
	}

	// no tick
	buckets.Increment(NewFloat(14.14))
	{
		require.Equal(t, NewFloat(114.13), buckets.Sum())
		require.Equal(t, "[114.130, 0.000, 0.000, 0.000]", buckets.String())
	}
}

func Test_CompareFloat(t *testing.T) {
	size, tickTime := testBucketSize()
	buckets := NewBuckets(size, 4, NewFloat(0))

	*tickTime = 0 // 0.000, 0.000, 0.000, 0.000
	require.Equal(t, false, buckets.Compare(LessEq, NewFloat(-1)))
	require.Equal(t, true, buckets.Compare(LessEq, NewFloat(0)))
	require.Equal(t, true, buckets.Compare(LessEq, NewFloat(1)))
	require.Equal(t, true, buckets.Compare(GreaterEq, NewFloat(-1)))
	require.Equal(t, true, buckets.Compare(GreaterEq, NewFloat(0)))
	require.Equal(t, false, buckets.Compare(GreaterEq, NewFloat(1)))

	*tickTime = 1                    // 4.400, 0.000, 0.000, 0.000
	buckets.Increment(NewFloat(4.4)) // last complete is still zero
	require.Equal(t, false, buckets.Compare(LessEq, NewFloat(-1)))
	require.Equal(t, true, buckets.Compare(LessEq, NewFloat(0)))
	require.Equal(t, true, buckets.Compare(LessEq, NewFloat(1)))
	require.Equal(t, true, buckets.Compare(GreaterEq, NewFloat(-1)))
	require.Equal(t, true, buckets.Compare(GreaterEq, NewFloat(0)))
	require.Equal(t, false, buckets.Compare(GreaterEq, NewFloat(1)))

	*tickTime = 2                    // 2.200, 4.400, 0.000, 0.000
	buckets.Increment(NewFloat(2.2)) // last complete is still zero
	require.Equal(t, false, buckets.Compare(LessEq, NewFloat(4.3)))
	require.Equal(t, true, buckets.Compare(LessEq, NewFloat(4.4)))
	require.Equal(t, true, buckets.Compare(LessEq, NewFloat(4.5)))
	require.Equal(t, true, buckets.Compare(GreaterEq, NewFloat(4.3)))
	require.Equal(t, true, buckets.Compare(GreaterEq, NewFloat(4.4)))
	require.Equal(t, false, buckets.Compare(GreaterEq, NewFloat(4.5)))

	*tickTime = 3 // 0.000 2.200, 4.400, 0.000
	require.Equal(t, false, buckets.Compare(LessEq, NewFloat(2.1)))
	require.Equal(t, true, buckets.Compare(LessEq, NewFloat(2.2)))
	require.Equal(t, true, buckets.Compare(LessEq, NewFloat(2.3)))
	require.Equal(t, true, buckets.Compare(GreaterEq, NewFloat(2.1)))
	require.Equal(t, true, buckets.Compare(GreaterEq, NewFloat(2.2)))
	require.Equal(t, false, buckets.Compare(GreaterEq, NewFloat(2.3)))
}

func Test_CompareInt(t *testing.T) {
	size, tickTime := testBucketSize()
	buckets := NewBuckets(size, 4, NewInt(0))

	*tickTime = 0 // 0, 0, 0, 0
	require.Equal(t, false, buckets.Compare(LessEq, NewInt(-1)))
	require.Equal(t, true, buckets.Compare(LessEq, NewInt(0)))
	require.Equal(t, true, buckets.Compare(LessEq, NewInt(1)))
	require.Equal(t, true, buckets.Compare(GreaterEq, NewInt(-1)))
	require.Equal(t, true, buckets.Compare(GreaterEq, NewInt(0)))
	require.Equal(t, false, buckets.Compare(GreaterEq, NewInt(1)))

	*tickTime = 1                // 4, 0, 0, 0
	buckets.Increment(NewInt(4)) // last complete is still zero
	require.Equal(t, false, buckets.Compare(LessEq, NewInt(-1)))
	require.Equal(t, true, buckets.Compare(LessEq, NewInt(0)))
	require.Equal(t, true, buckets.Compare(LessEq, NewInt(1)))
	require.Equal(t, true, buckets.Compare(GreaterEq, NewInt(-1)))
	require.Equal(t, true, buckets.Compare(GreaterEq, NewInt(0)))
	require.Equal(t, false, buckets.Compare(GreaterEq, NewInt(1)))

	*tickTime = 2                // 2, 4, 0, 0
	buckets.Increment(NewInt(2)) // last complete is still zero
	require.Equal(t, false, buckets.Compare(LessEq, NewInt(3)))
	require.Equal(t, true, buckets.Compare(LessEq, NewInt(4)))
	require.Equal(t, true, buckets.Compare(LessEq, NewInt(5)))
	require.Equal(t, true, buckets.Compare(GreaterEq, NewInt(3)))
	require.Equal(t, true, buckets.Compare(GreaterEq, NewInt(4)))
	require.Equal(t, false, buckets.Compare(GreaterEq, NewInt(5)))

	*tickTime = 3 // 0, 2, 4, 0
	require.Equal(t, false, buckets.Compare(LessEq, NewInt(1)))
	require.Equal(t, true, buckets.Compare(LessEq, NewInt(2)))
	require.Equal(t, true, buckets.Compare(LessEq, NewInt(3)))
	require.Equal(t, true, buckets.Compare(GreaterEq, NewInt(1)))
	require.Equal(t, true, buckets.Compare(GreaterEq, NewInt(2)))
	require.Equal(t, false, buckets.Compare(GreaterEq, NewInt(3)))
}
