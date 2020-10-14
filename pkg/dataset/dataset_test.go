package dataset_test

import (
	"math/rand"
	"testing"

	"github.com/bwagner5/ds/pkg/dataset"
	h "github.com/bwagner5/ds/pkg/test"
)

func TestPut(t *testing.T) {
	ds := dataset.New()
	ds.Put(9.8)
	ds.Put(10.1)
	ds.Put(1.1)
	ds.Put(0.5)
	ds.Put(2.0)

	h.Assert(t, 0.5 == ds.Get(0), "the 0th index should equal min")
	h.Assert(t, 1.1 == ds.Get(1), "the first index should equal 1.1")
	h.Assert(t, 2 == ds.Get(2), "the second index should equal 2")
	h.Assert(t, 10.1 == ds.Get(len(ds.SortedSlice())-1), "the last index should equal max")
}

func TestGetFreq(t *testing.T) {
	ds := dataset.New()
	ds.Put(9.8)
	ds.Put(10.1)
	ds.Put(10.1)
	ds.Put(1.1)
	ds.Put(1.1)
	ds.Put(0.5)
	ds.Put(2.0)
	ds.Put(1.1)
	ds.Put(1.1)

	h.Assert(t, 1 == ds.GetFreq(9.8), "only 1 9.8")
	h.Assert(t, 2 == ds.GetFreq(10.1), "should be 2 10.1's")
	h.Assert(t, 4 == ds.GetFreq(1.1), "should 4 3 1.1's")
	h.Assert(t, 0 == ds.GetFreq(999.1), "should be 0 999.1")
}

func TestCompressedLen(t *testing.T) {
	ds := dataset.New()
	ds.Put(9.8)
	ds.Put(10.1)
	ds.Put(10.1)
	ds.Put(1.1)
	ds.Put(1.1)
	ds.Put(0.5)
	ds.Put(2.0)
	ds.Put(1.1)
	ds.Put(1.1)

	h.Assert(t, 5 == ds.CompressedLen(), "should be 5 uniq vals")
}

func TestTotalCount(t *testing.T) {
	ds := dataset.New()
	ds.Put(9.8)
	ds.Put(10.1)
	ds.Put(10.1)
	ds.Put(1.1)
	ds.Put(1.1)
	ds.Put(0.5)
	ds.Put(2.0)
	ds.Put(1.1)
	ds.Put(1.1)

	h.Assert(t, 9 == ds.TotalCount(), "should be 9 vals")
}

func TestGet(t *testing.T) {
	ds := dataset.New()
	ds.Put(9.8)
	ds.Put(10.1)
	ds.Put(10.1)
	ds.Put(1.1)
	ds.Put(1.1)
	ds.Put(0.5)
	ds.Put(2.0)
	ds.Put(1.1)
	ds.Put(1.1)

	h.Assert(t, 0.5 == ds.Get(0), "should have retrieved 0th index of sorted values")
}

func TestSum(t *testing.T) {
	ds := dataset.New()
	ds.Put(9.8)
	ds.Put(10.1)
	ds.Put(1.1)
	ds.Put(0.5)
	ds.Put(2.0)

	h.Assert(t, 23.5 == ds.Sum(), "should have summed unique values to 23.5")

	ds.Put(9.8)
	ds.Put(10.1)
	ds.Put(1.1)
	ds.Put(0.5)
	ds.Put(2.0)

	h.Assert(t, 47 == ds.Sum(), "should have summed duplicate values to 47")
}

func TestCalcMean(t *testing.T) {
	ds := dataset.New()
	ds.Put(9.8)
	ds.Put(10.1)
	ds.Put(1.1)
	ds.Put(0.5)
	ds.Put(2.0)

	h.Assert(t, 4.7 == ds.CalcMean(), "should have found mean of unique values to be 4.7")

	ds.Put(9.8)
	ds.Put(10.1)
	ds.Put(1.1)
	ds.Put(0.5)
	ds.Put(2.0)

	h.Assert(t, 4.7 == ds.CalcMean(), "should have found mean of duplicate values to be 4.7")
}

func TestCalcMedian(t *testing.T) {
	ds := dataset.New()
	ds.Put(9.8)
	ds.Put(10.1)
	ds.Put(1.1)
	ds.Put(0.5)
	ds.Put(2.0)

	h.Equals(t, 2.0, ds.CalcMedian())

	ds.Put(100)

	h.Equals(t, 5.9, ds.CalcMedian())

	ds.Put(9.8)

	h.Equals(t, 9.8, ds.CalcMedian())
}

func TestCalcVariance(t *testing.T) {
	ds := dataset.New()
	ds.Put(9.8)
	ds.Put(10.1)
	ds.Put(1.1)
	ds.Put(0.5)
	ds.Put(2.0)

	h.Equals(t, 18.612000000000002, ds.CalcVariance())

	ds.Put(9.8)
	ds.Put(10.1)
	ds.Put(1.1)
	ds.Put(0.5)
	ds.Put(2.0)

	h.Equals(t, 18.612000000000002, ds.CalcVariance())
}

func TestCalcStdDev(t *testing.T) {
	ds := dataset.New()
	ds.Put(9.8)
	ds.Put(10.1)
	ds.Put(1.1)
	ds.Put(0.5)
	ds.Put(2.0)

	h.Equals(t, 4.314162722939412, ds.CalcStdDev())

	ds.Put(9.8)
	ds.Put(10.1)
	ds.Put(1.1)
	ds.Put(0.5)
	ds.Put(2.0)

	h.Equals(t, 4.314162722939412, ds.CalcStdDev())
}

func TestCalcPercentile(t *testing.T) {
	ds := dataset.New()
	ds.Put(9.8)
	ds.Put(10.1)
	ds.Put(1.1)
	ds.Put(0.5)
	ds.Put(2.0)

	h.Equals(t, 0.5, ds.CalcPercentile(0))
	h.Equals(t, 10.1, ds.CalcPercentile(100))
	h.Equals(t, 1.55, ds.CalcPercentile(50))

	ds.Put(9.8)
	ds.Put(10.1)
	ds.Put(1.1)
	ds.Put(0.5)
	ds.Put(2.0)

	h.Equals(t, 0.5, ds.CalcPercentile(10))
	h.Equals(t, 0.5, ds.CalcPercentile(20))
	h.Equals(t, 1.1, ds.CalcPercentile(30))
	h.Equals(t, 1.1, ds.CalcPercentile(40))
	h.Equals(t, 2.0, ds.CalcPercentile(50))
	h.Equals(t, 2.0, ds.CalcPercentile(60))
	h.Equals(t, 9.8, ds.CalcPercentile(70))
	h.Equals(t, 9.8, ds.CalcPercentile(80))
	h.Equals(t, 10.1, ds.CalcPercentile(90))
}

func TestMaxFrequencyNum(t *testing.T) {
	ds := dataset.New()
	ds.Put(9.8)
	ds.Put(9.8)
	ds.Put(10.1)
	ds.Put(1.1)
	ds.Put(0.5)
	ds.Put(2.0)

	h.Equals(t, 9.8, ds.MaxFrequencyNum())

	ds.Put(10.1)
	ds.Put(10.1)
	ds.Put(10.1)
	ds.Put(10.1)

	h.Equals(t, 10.1, ds.MaxFrequencyNum())
}

func TestFindNumbersGreaterThan(t *testing.T) {
	ds := dataset.New()
	ds.Put(9.8)
	ds.Put(9.8)
	ds.Put(10.1)
	ds.Put(1.1)
	ds.Put(0.5)
	ds.Put(2.0)

	h.Equals(t, int64(5), ds.FindNumbersGreaterThan(0.5))

	ds.Put(10.1)
	ds.Put(10.1)
	ds.Put(10.1)
	ds.Put(10.1)

	h.Equals(t, int64(9), ds.FindNumbersGreaterThan(0.5))
	h.Equals(t, int64(10), ds.FindNumbersGreaterThan(0))
	h.Equals(t, int64(0), ds.FindNumbersGreaterThan(100))
}

func TestFindNumbersLessThan(t *testing.T) {
	ds := dataset.New()
	ds.Put(9.8)
	ds.Put(9.8)
	ds.Put(10.1)
	ds.Put(1.1)
	ds.Put(0.5)
	ds.Put(2.0)

	h.Equals(t, int64(2), ds.FindNumbersLessThan(2.0))

	ds.Put(10.1)
	ds.Put(10.1)
	ds.Put(10.1)
	ds.Put(10.1)

	h.Equals(t, int64(5), ds.FindNumbersLessThan(10.1))
	h.Equals(t, int64(10), ds.FindNumbersLessThan(10.2))
	h.Equals(t, int64(0), ds.FindNumbersLessThan(0))
}

func BenchmarkPut(b *testing.B) {
	ds := dataset.New()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ds.Put(rand.NormFloat64())
		}
	})
}
