package metrics

import (
	"sort"

	"go.uber.org/zap/zapcore"
)

type (
	// Stat is the type of statistical metrics.
	Stat int
	// Stats contains all metrics packed as an array.
	Stats [NumStats]uint64
	// PerTargetStats is a map of Stats per target.
	PerTargetStats map[string]Stats
)

const (
	RequestsAttemptedStat Stat = iota
	RequestsSentStat
	ResponsesReceivedStat
	BytesSentStat
	BytesReceivedStat

	NumStats
)

func (ts PerTargetStats) sortedTargets() []string {
	res := make([]string, 0, len(ts))
	for k := range ts {
		res = append(res, k)
	}

	sort.Strings(res)

	return res
}

func Diff(lhs, rhs Stats) Stats {
	var res Stats
	for i := range res {
		res[i] = lhs[i] - rhs[i]
	}

	return res
}

func (ts PerTargetStats) Diff(other PerTargetStats) PerTargetStats {
	if other == nil {
		return ts
	}

	res := make(PerTargetStats)
	for k := range ts {
		res[k] = Diff(ts[k], other[k])
	}

	return res
}

// MarshalLogObject is required to log PerTargetStats objects to zap
func (ts PerTargetStats) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	for _, tgt := range ts.sortedTargets() {
		tgtStats := ts[tgt]

		if err := enc.AddObject(tgt, &tgtStats); err != nil {
			return err
		}
	}

	return nil
}

// MarshalLogObject is required to log Stats objects to zap
func (stats *Stats) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddUint64("requests_attempted", stats[RequestsAttemptedStat])
	enc.AddUint64("requests_sent", stats[RequestsSentStat])
	enc.AddUint64("responses_received", stats[ResponsesReceivedStat])
	enc.AddUint64("bytes_sent", stats[BytesSentStat])
	enc.AddUint64("bytes_received", stats[BytesReceivedStat])

	return nil
}
