package dataset

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
	"text/tabwriter"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type DataSet struct {
	sync.RWMutex
	numsLookup map[float64]int
	nums       []float64
	freqs      []int64
	cache      *statsCache
}

type statsCache struct {
	totalCount    *int64
	sum           *float64
	indexOfBucket map[int64]int
	getFreq       map[float64]int64
}

func newStatsCache() *statsCache {
	s := statsCache{}
	s.indexOfBucket = map[int64]int{}
	s.getFreq = map[float64]int64{}
	return &s
}

func (s *statsCache) isCacheEmpty() bool {
	return s.totalCount == nil &&
		s.sum == nil &&
		len(s.indexOfBucket) == 0 &&
		len(s.getFreq) == 0
}

func (s *statsCache) Clear() {
	if !s.isCacheEmpty() {
		s.totalCount = nil
		s.sum = nil
		s.indexOfBucket = map[int64]int{}
		s.getFreq = map[float64]int64{}
	}
}

// New creates an instance of a dataset to derive statistics from
func New() *DataSet {
	s := DataSet{}
	s.numsLookup = make(map[float64]int)
	s.nums = []float64{}
	s.freqs = []int64{}
	s.cache = newStatsCache()
	return &s
}

// Put adds items to the dataset and maintains a value sort based on the num
func (s *DataSet) Put(num float64) {
	s.RLock()
	if _, ok := s.numsLookup[num]; ok {
		s.cache.Clear()
		s.RUnlock()
		for i, n := range s.nums {
			if n == num {
				s.freqs[i]++
				break
			}
		}
	} else {
		s.RUnlock()
		insertIndex := len(s.nums)
		for i, n := range s.nums {
			if n > num {
				insertIndex = i
				break
			}
		}
		s.Lock()
		defer s.Unlock()
		s.cache.Clear()
		s.nums = append(s.nums, 0)
		copy(s.nums[insertIndex+1:], s.nums[insertIndex:])
		s.nums[insertIndex] = num
		s.freqs = append(s.freqs, 0)
		copy(s.freqs[insertIndex+1:], s.freqs[insertIndex:])
		s.freqs[insertIndex] = 1
		s.numsLookup[num] = 1
	}
}

// Get looks up a number in a sorted unique values slice
func (s *DataSet) Get(index int) float64 {
	return s.nums[index]
}

// GetFreq looks up the frequency of a number that was Put in the dataset
func (s *DataSet) GetFreq(num float64) int64 {
	if f, ok := s.cache.getFreq[num]; ok {
		return f
	}
	for i, n := range s.nums {
		if n == num {
			s.cache.getFreq[num] = s.freqs[i]
			return s.freqs[i]
		}
	}
	return 0
}

// CompressedLen returns the length of the unique values in the dataset
func (s *DataSet) CompressedLen() int {
	return len(s.nums)
}

// TotalCount returns the length of all values (including duplicates) in the dataset
func (s *DataSet) TotalCount() int64 {
	if s.cache.totalCount != nil {
		return *s.cache.totalCount
	}
	fSum := int64(0)
	for _, f := range s.freqs {
		fSum += f
	}
	s.cache.totalCount = &fSum
	return fSum
}

// SortedSlice returns a sorted slice of the unique values in the dataset
func (s *DataSet) SortedSlice() []float64 {
	return s.nums
}

// Sum returns the sum of the values in the dataset
func (s *DataSet) Sum() float64 {
	if s.cache.sum != nil {
		return *s.cache.sum
	}
	sum := float64(0)
	for _, n := range s.nums {
		sum += (n * float64(s.GetFreq(n)))
	}
	s.cache.sum = &sum
	return sum
}

// CalcMean returns the mean of the dataset
func (s *DataSet) CalcMean() float64 {
	return float64(s.Sum()) / float64(s.TotalCount())
}

// CalcMedian returns the median of the dataset
func (s *DataSet) CalcMedian() float64 {
	n := s.TotalCount()
	medianIndex := float64(n) / float64(2)
	if n%2 != 0 {
		medianIndex = float64(n+1) / float64(2)
		return s.Get(s.getIndexOfBucket(int64(medianIndex)))
	}
	first := s.Get(s.getIndexOfBucket(int64(medianIndex)))
	second := s.Get(s.getIndexOfBucket(int64(medianIndex + 1)))
	return (first + second) / 2
}

// CalcVariance returns the variance of the dataset
func (s *DataSet) CalcVariance() float64 {
	mu := s.CalcMean()
	sumOfDistanceToMu := float64(0)
	for _, num := range s.SortedSlice() {
		count := s.GetFreq(num)
		for i := int64(0); i < count; i++ {
			sumOfDistanceToMu += math.Pow(num-mu, 2)
		}
	}
	return sumOfDistanceToMu / float64(s.TotalCount())
}

// CalcStdDev returns the standard deviation of the dataset
func (s *DataSet) CalcStdDev() float64 {
	return math.Sqrt(s.CalcVariance())
}

// CalcPercentile returns the specified percentile of the dataset
func (s *DataSet) CalcPercentile(p float64) float64 {
	index := (p / 100) * float64(s.TotalCount())
	// whole number index
	if index == float64(int64(index)) {
		return s.Get(s.getIndexOfBucket(int64(index)))
	}
	bucketIndex := s.getIndexOfBucket(int64(index))
	adjSum := float64(0)
	if bucketIndex >= s.CompressedLen()-1 {
		adjSum = s.Get(bucketIndex-1) + s.Get(bucketIndex)
	} else {
		adjSum = s.Get(bucketIndex+1) + s.Get(bucketIndex)

	}
	return float64(adjSum) / 2
}

// getIndexOfBucket takes an index based on total values in the dataset
// and returns the index of the bucket for unique value look ups
func (s *DataSet) getIndexOfBucket(i int64) int {
	if index, ok := s.cache.indexOfBucket[i]; ok {
		return index
	}
	index := int64(0)
	for k, j := range s.freqs {
		index += j
		if index >= i {
			s.cache.indexOfBucket[i] = int(k)
			return int(k)
		}
	}
	return -1
}

// MaxFrequencyNum returns the value in the dataset with the highest modality
func (s *DataSet) MaxFrequencyNum() float64 {
	max := int64(0)
	maxNum := s.Get(0)
	for _, num := range s.SortedSlice() {
		freq := s.GetFreq(num)
		if freq > max {
			max = freq
			maxNum = num
		}
	}
	return maxNum
}

// FindNumbersGreaterThan returns the count of numbers that are greater than the specified number
func (s *DataSet) FindNumbersGreaterThan(n float64) int64 {
	greaterThanSum := int64(0)
	for _, num := range s.SortedSlice() {
		if num > n {
			greaterThanSum += s.GetFreq(num)
		}
	}
	return greaterThanSum
}

// FindNumbersLessThan returns the count of numbers that are less than the specified number
func (s *DataSet) FindNumbersLessThan(n float64) int64 {
	lessThanSum := int64(0)
	for _, num := range s.SortedSlice() {
		if num < n {
			lessThanSum += s.GetFreq(num)
		}
	}
	return lessThanSum
}

// SummaryString returns a string of summary statistics on the dataset
func (s *DataSet) SummaryString(vertical bool) string {
	var out bytes.Buffer
	headers := []string{
		"n",
		"mean",
		"median",
		"std dev",
		"min",
		"max",
		"P99.99",
		"P99",
		"P95",
		"P75",
		"P25",
		"P5",
		"P1",
		"P0.01",
		"Top Freq Num",
		"Top Freq",
		"> Top Freq",
		"< Top Freq",
	}
	p := message.NewPrinter(language.English)
	maxFreqNum := s.MaxFrequencyNum()
	stats := []string{
		p.Sprintf("%d", s.TotalCount()),
		fmt.Sprintf("%s", formatFloat(s.CalcMean())),
		fmt.Sprintf("%s", formatFloat(s.CalcMedian())),
		fmt.Sprintf("%s", formatFloat(s.CalcStdDev())),
		fmt.Sprintf("%s", formatFloat(s.CalcPercentile(0))),
		fmt.Sprintf("%s", formatFloat(s.CalcPercentile(100))),
		fmt.Sprintf("%s", formatFloat(s.CalcPercentile(99.99))),
		fmt.Sprintf("%s", formatFloat(s.CalcPercentile(99))),
		fmt.Sprintf("%s", formatFloat(s.CalcPercentile(95))),
		fmt.Sprintf("%s", formatFloat(s.CalcPercentile(75))),
		fmt.Sprintf("%s", formatFloat(s.CalcPercentile(25))),
		fmt.Sprintf("%s", formatFloat(s.CalcPercentile(5))),
		fmt.Sprintf("%s", formatFloat(s.CalcPercentile(1))),
		fmt.Sprintf("%s", formatFloat(s.CalcPercentile(0.01))),
		fmt.Sprintf("%s", formatFloat(maxFreqNum)),
		p.Sprintf("%d", s.GetFreq(maxFreqNum)),
		p.Sprintf("%d", s.FindNumbersGreaterThan(maxFreqNum)),
		p.Sprintf("%d", s.FindNumbersLessThan(maxFreqNum)),
	}
	w := tabwriter.NewWriter(&out, 8, 8, 8, ' ', 0)
	delim := "\t"
	if vertical {
		for i, h := range headers {
			fmt.Fprintf(w, "%s:\t%s\n", h, stats[i])
		}
	} else {
		fmt.Fprintln(w, strings.Join(headers, delim)+delim)
		fmt.Fprintf(w, strings.Join(stats, delim)+delim)
	}
	w.Flush()
	return out.String()
}

func formatFloat(f float64) string {
	s := strconv.FormatFloat(f, 'f', 5, 64)
	parts := strings.Split(s, ".")
	reversed := reverse(parts[0])
	withCommas := ""
	for i, p := range reversed {
		if i%3 == 0 && i != 0 {
			withCommas += ","
		}
		withCommas += string(p)
	}
	s = strings.Join([]string{reverse(withCommas), parts[1]}, ".")
	return strings.TrimRight(strings.TrimRight(s, "0"), ".")
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
