// Generated by tmpl
// https://github.com/benbjohnson/tmpl
//
// DO NOT EDIT!
// Source: response_writer.gen.go.tmpl

package reads

import (
	"github.com/influxdata/influxdb/storage/reads/datatypes"
	"github.com/influxdata/influxdb/tsdb/cursors"
)

func (w *ResponseWriter) getFloatPointsFrame() *datatypes.ReadResponse_Frame_FloatPoints {
	var res *datatypes.ReadResponse_Frame_FloatPoints
	if len(w.buffer.Float) > 0 {
		i := len(w.buffer.Float) - 1
		res = w.buffer.Float[i]
		w.buffer.Float[i] = nil
		w.buffer.Float = w.buffer.Float[:i]
	} else {
		res = &datatypes.ReadResponse_Frame_FloatPoints{
			FloatPoints: &datatypes.ReadResponse_FloatPointsFrame{
				Timestamps: make([]int64, 0, batchSize),
				Values:     make([]float64, 0, batchSize),
			},
		}
	}
	return res
}

func (w *ResponseWriter) putFloatPointsFrame(f *datatypes.ReadResponse_Frame_FloatPoints) {
	f.FloatPoints.Timestamps = f.FloatPoints.Timestamps[:0]
	f.FloatPoints.Values = f.FloatPoints.Values[:0]
	w.buffer.Float = append(w.buffer.Float, f)
}

func (w *ResponseWriter) getFloatValues() *datatypes.ReadResponse_AnyPoints_Floats {
	var res *datatypes.ReadResponse_AnyPoints_Floats
	if len(w.buffer.FloatValues) > 0 {
		i := len(w.buffer.FloatValues) - 1
		res = w.buffer.FloatValues[i]
		w.buffer.FloatValues[i] = nil
		w.buffer.FloatValues = w.buffer.FloatValues[:i]
	} else {
		res = &datatypes.ReadResponse_AnyPoints_Floats{
			Floats: &datatypes.ReadResponse_FloatValues{
				Values: make([]float64, 0, batchSize),
			},
		}
	}
	return res
}

func (w *ResponseWriter) putFloatValues(f *datatypes.ReadResponse_AnyPoints_Floats) {
	f.Floats.Values = f.Floats.Values[:0]
	w.buffer.FloatValues = append(w.buffer.FloatValues, f)
}

func (w *ResponseWriter) streamFloatArraySeries(cur cursors.FloatArrayCursor) {
	w.sf.DataType = datatypes.DataTypeFloat
	ss := len(w.res.Frames) - 1
	a := cur.Next()
	if len(a.Timestamps) == 0 {
		w.sz -= w.sf.Size()
		w.putSeriesFrame(w.res.Frames[ss].Data.(*datatypes.ReadResponse_Frame_Series))
		w.res.Frames = w.res.Frames[:ss]
	} else if w.sz > writeSize {
		w.Flush()
	}
}

func (w *ResponseWriter) streamFloatArrayPoints(cur cursors.FloatArrayCursor) {
	w.sf.DataType = datatypes.DataTypeFloat
	ss := len(w.res.Frames) - 1

	p := w.getFloatPointsFrame()
	frame := p.FloatPoints
	w.res.Frames = append(w.res.Frames, datatypes.ReadResponse_Frame{Data: p})

	var seriesValueCount = 0
	for {
		// If the number of values produced by cur > 1000,
		// cur.Next() will produce batches of values that are of
		// length ≤ 1000.
		// We attempt to limit the frame Timestamps / Values lengths
		// the same to avoid allocations. These frames are recycled
		// after flushing so that on repeated use there should be enough space
		// to append values from a into frame without additional allocations.
		a := cur.Next()

		if len(a.Timestamps) == 0 {
			break
		}

		seriesValueCount += a.Len()
		// As specified in the struct definition, w.sz is an estimated
		// size (in bytes) of the buffered data. It is therefore a
		// deliberate choice to accumulate using the array Size, which is
		// cheap to calculate. Calling frame.Size() can be expensive
		// when using varint encoding for numbers.
		w.sz += a.Size()

		frame.Timestamps = append(frame.Timestamps, a.Timestamps...)
		frame.Values = append(frame.Values, a.Values...)

		// given the expectation of cur.Next, we attempt to limit
		// the number of values appended to the frame to batchSize (1000)
		needsFrame := len(frame.Timestamps) >= batchSize

		if w.sz >= writeSize {
			needsFrame = true
			w.Flush()
			if w.err != nil {
				break
			}
		}

		if needsFrame {
			// new frames are returned with Timestamps and Values preallocated
			// to a minimum of batchSize length to reduce further allocations.
			p = w.getFloatPointsFrame()
			frame = p.FloatPoints
			w.res.Frames = append(w.res.Frames, datatypes.ReadResponse_Frame{Data: p})
		}
	}

	w.vc += seriesValueCount
	if seriesValueCount == 0 {
		w.sz -= w.sf.Size()
		w.putSeriesFrame(w.res.Frames[ss].Data.(*datatypes.ReadResponse_Frame_Series))
		w.res.Frames = w.res.Frames[:ss]
	} else if w.sz > writeSize {
		w.Flush()
	}
}

func (w *ResponseWriter) getIntegerPointsFrame() *datatypes.ReadResponse_Frame_IntegerPoints {
	var res *datatypes.ReadResponse_Frame_IntegerPoints
	if len(w.buffer.Integer) > 0 {
		i := len(w.buffer.Integer) - 1
		res = w.buffer.Integer[i]
		w.buffer.Integer[i] = nil
		w.buffer.Integer = w.buffer.Integer[:i]
	} else {
		res = &datatypes.ReadResponse_Frame_IntegerPoints{
			IntegerPoints: &datatypes.ReadResponse_IntegerPointsFrame{
				Timestamps: make([]int64, 0, batchSize),
				Values:     make([]int64, 0, batchSize),
			},
		}
	}
	return res
}

func (w *ResponseWriter) putIntegerPointsFrame(f *datatypes.ReadResponse_Frame_IntegerPoints) {
	f.IntegerPoints.Timestamps = f.IntegerPoints.Timestamps[:0]
	f.IntegerPoints.Values = f.IntegerPoints.Values[:0]
	w.buffer.Integer = append(w.buffer.Integer, f)
}

func (w *ResponseWriter) getIntegerValues() *datatypes.ReadResponse_AnyPoints_Integers {
	var res *datatypes.ReadResponse_AnyPoints_Integers
	if len(w.buffer.IntegerValues) > 0 {
		i := len(w.buffer.IntegerValues) - 1
		res = w.buffer.IntegerValues[i]
		w.buffer.IntegerValues[i] = nil
		w.buffer.IntegerValues = w.buffer.IntegerValues[:i]
	} else {
		res = &datatypes.ReadResponse_AnyPoints_Integers{
			Integers: &datatypes.ReadResponse_IntegerValues{
				Values: make([]int64, 0, batchSize),
			},
		}
	}
	return res
}

func (w *ResponseWriter) putIntegerValues(f *datatypes.ReadResponse_AnyPoints_Integers) {
	f.Integers.Values = f.Integers.Values[:0]
	w.buffer.IntegerValues = append(w.buffer.IntegerValues, f)
}

func (w *ResponseWriter) streamIntegerArraySeries(cur cursors.IntegerArrayCursor) {
	w.sf.DataType = datatypes.DataTypeInteger
	ss := len(w.res.Frames) - 1
	a := cur.Next()
	if len(a.Timestamps) == 0 {
		w.sz -= w.sf.Size()
		w.putSeriesFrame(w.res.Frames[ss].Data.(*datatypes.ReadResponse_Frame_Series))
		w.res.Frames = w.res.Frames[:ss]
	} else if w.sz > writeSize {
		w.Flush()
	}
}

func (w *ResponseWriter) streamIntegerArrayPoints(cur cursors.IntegerArrayCursor) {
	w.sf.DataType = datatypes.DataTypeInteger
	ss := len(w.res.Frames) - 1

	p := w.getIntegerPointsFrame()
	frame := p.IntegerPoints
	w.res.Frames = append(w.res.Frames, datatypes.ReadResponse_Frame{Data: p})

	var seriesValueCount = 0
	for {
		// If the number of values produced by cur > 1000,
		// cur.Next() will produce batches of values that are of
		// length ≤ 1000.
		// We attempt to limit the frame Timestamps / Values lengths
		// the same to avoid allocations. These frames are recycled
		// after flushing so that on repeated use there should be enough space
		// to append values from a into frame without additional allocations.
		a := cur.Next()

		if len(a.Timestamps) == 0 {
			break
		}

		seriesValueCount += a.Len()
		// As specified in the struct definition, w.sz is an estimated
		// size (in bytes) of the buffered data. It is therefore a
		// deliberate choice to accumulate using the array Size, which is
		// cheap to calculate. Calling frame.Size() can be expensive
		// when using varint encoding for numbers.
		w.sz += a.Size()

		frame.Timestamps = append(frame.Timestamps, a.Timestamps...)
		frame.Values = append(frame.Values, a.Values...)

		// given the expectation of cur.Next, we attempt to limit
		// the number of values appended to the frame to batchSize (1000)
		needsFrame := len(frame.Timestamps) >= batchSize

		if w.sz >= writeSize {
			needsFrame = true
			w.Flush()
			if w.err != nil {
				break
			}
		}

		if needsFrame {
			// new frames are returned with Timestamps and Values preallocated
			// to a minimum of batchSize length to reduce further allocations.
			p = w.getIntegerPointsFrame()
			frame = p.IntegerPoints
			w.res.Frames = append(w.res.Frames, datatypes.ReadResponse_Frame{Data: p})
		}
	}

	w.vc += seriesValueCount
	if seriesValueCount == 0 {
		w.sz -= w.sf.Size()
		w.putSeriesFrame(w.res.Frames[ss].Data.(*datatypes.ReadResponse_Frame_Series))
		w.res.Frames = w.res.Frames[:ss]
	} else if w.sz > writeSize {
		w.Flush()
	}
}

func (w *ResponseWriter) getUnsignedPointsFrame() *datatypes.ReadResponse_Frame_UnsignedPoints {
	var res *datatypes.ReadResponse_Frame_UnsignedPoints
	if len(w.buffer.Unsigned) > 0 {
		i := len(w.buffer.Unsigned) - 1
		res = w.buffer.Unsigned[i]
		w.buffer.Unsigned[i] = nil
		w.buffer.Unsigned = w.buffer.Unsigned[:i]
	} else {
		res = &datatypes.ReadResponse_Frame_UnsignedPoints{
			UnsignedPoints: &datatypes.ReadResponse_UnsignedPointsFrame{
				Timestamps: make([]int64, 0, batchSize),
				Values:     make([]uint64, 0, batchSize),
			},
		}
	}
	return res
}

func (w *ResponseWriter) putUnsignedPointsFrame(f *datatypes.ReadResponse_Frame_UnsignedPoints) {
	f.UnsignedPoints.Timestamps = f.UnsignedPoints.Timestamps[:0]
	f.UnsignedPoints.Values = f.UnsignedPoints.Values[:0]
	w.buffer.Unsigned = append(w.buffer.Unsigned, f)
}

func (w *ResponseWriter) getUnsignedValues() *datatypes.ReadResponse_AnyPoints_Unsigneds {
	var res *datatypes.ReadResponse_AnyPoints_Unsigneds
	if len(w.buffer.UnsignedValues) > 0 {
		i := len(w.buffer.UnsignedValues) - 1
		res = w.buffer.UnsignedValues[i]
		w.buffer.UnsignedValues[i] = nil
		w.buffer.UnsignedValues = w.buffer.UnsignedValues[:i]
	} else {
		res = &datatypes.ReadResponse_AnyPoints_Unsigneds{
			Unsigneds: &datatypes.ReadResponse_UnsignedValues{
				Values: make([]uint64, 0, batchSize),
			},
		}
	}
	return res
}

func (w *ResponseWriter) putUnsignedValues(f *datatypes.ReadResponse_AnyPoints_Unsigneds) {
	f.Unsigneds.Values = f.Unsigneds.Values[:0]
	w.buffer.UnsignedValues = append(w.buffer.UnsignedValues, f)
}

func (w *ResponseWriter) streamUnsignedArraySeries(cur cursors.UnsignedArrayCursor) {
	w.sf.DataType = datatypes.DataTypeUnsigned
	ss := len(w.res.Frames) - 1
	a := cur.Next()
	if len(a.Timestamps) == 0 {
		w.sz -= w.sf.Size()
		w.putSeriesFrame(w.res.Frames[ss].Data.(*datatypes.ReadResponse_Frame_Series))
		w.res.Frames = w.res.Frames[:ss]
	} else if w.sz > writeSize {
		w.Flush()
	}
}

func (w *ResponseWriter) streamUnsignedArrayPoints(cur cursors.UnsignedArrayCursor) {
	w.sf.DataType = datatypes.DataTypeUnsigned
	ss := len(w.res.Frames) - 1

	p := w.getUnsignedPointsFrame()
	frame := p.UnsignedPoints
	w.res.Frames = append(w.res.Frames, datatypes.ReadResponse_Frame{Data: p})

	var seriesValueCount = 0
	for {
		// If the number of values produced by cur > 1000,
		// cur.Next() will produce batches of values that are of
		// length ≤ 1000.
		// We attempt to limit the frame Timestamps / Values lengths
		// the same to avoid allocations. These frames are recycled
		// after flushing so that on repeated use there should be enough space
		// to append values from a into frame without additional allocations.
		a := cur.Next()

		if len(a.Timestamps) == 0 {
			break
		}

		seriesValueCount += a.Len()
		// As specified in the struct definition, w.sz is an estimated
		// size (in bytes) of the buffered data. It is therefore a
		// deliberate choice to accumulate using the array Size, which is
		// cheap to calculate. Calling frame.Size() can be expensive
		// when using varint encoding for numbers.
		w.sz += a.Size()

		frame.Timestamps = append(frame.Timestamps, a.Timestamps...)
		frame.Values = append(frame.Values, a.Values...)

		// given the expectation of cur.Next, we attempt to limit
		// the number of values appended to the frame to batchSize (1000)
		needsFrame := len(frame.Timestamps) >= batchSize

		if w.sz >= writeSize {
			needsFrame = true
			w.Flush()
			if w.err != nil {
				break
			}
		}

		if needsFrame {
			// new frames are returned with Timestamps and Values preallocated
			// to a minimum of batchSize length to reduce further allocations.
			p = w.getUnsignedPointsFrame()
			frame = p.UnsignedPoints
			w.res.Frames = append(w.res.Frames, datatypes.ReadResponse_Frame{Data: p})
		}
	}

	w.vc += seriesValueCount
	if seriesValueCount == 0 {
		w.sz -= w.sf.Size()
		w.putSeriesFrame(w.res.Frames[ss].Data.(*datatypes.ReadResponse_Frame_Series))
		w.res.Frames = w.res.Frames[:ss]
	} else if w.sz > writeSize {
		w.Flush()
	}
}

func (w *ResponseWriter) getStringPointsFrame() *datatypes.ReadResponse_Frame_StringPoints {
	var res *datatypes.ReadResponse_Frame_StringPoints
	if len(w.buffer.String) > 0 {
		i := len(w.buffer.String) - 1
		res = w.buffer.String[i]
		w.buffer.String[i] = nil
		w.buffer.String = w.buffer.String[:i]
	} else {
		res = &datatypes.ReadResponse_Frame_StringPoints{
			StringPoints: &datatypes.ReadResponse_StringPointsFrame{
				Timestamps: make([]int64, 0, batchSize),
				Values:     make([]string, 0, batchSize),
			},
		}
	}
	return res
}

func (w *ResponseWriter) putStringPointsFrame(f *datatypes.ReadResponse_Frame_StringPoints) {
	f.StringPoints.Timestamps = f.StringPoints.Timestamps[:0]
	f.StringPoints.Values = f.StringPoints.Values[:0]
	w.buffer.String = append(w.buffer.String, f)
}

func (w *ResponseWriter) getStringValues() *datatypes.ReadResponse_AnyPoints_Strings {
	var res *datatypes.ReadResponse_AnyPoints_Strings
	if len(w.buffer.StringValues) > 0 {
		i := len(w.buffer.StringValues) - 1
		res = w.buffer.StringValues[i]
		w.buffer.StringValues[i] = nil
		w.buffer.StringValues = w.buffer.StringValues[:i]
	} else {
		res = &datatypes.ReadResponse_AnyPoints_Strings{
			Strings: &datatypes.ReadResponse_StringValues{
				Values: make([]string, 0, batchSize),
			},
		}
	}
	return res
}

func (w *ResponseWriter) putStringValues(f *datatypes.ReadResponse_AnyPoints_Strings) {
	f.Strings.Values = f.Strings.Values[:0]
	w.buffer.StringValues = append(w.buffer.StringValues, f)
}

func (w *ResponseWriter) streamStringArraySeries(cur cursors.StringArrayCursor) {
	w.sf.DataType = datatypes.DataTypeString
	ss := len(w.res.Frames) - 1
	a := cur.Next()
	if len(a.Timestamps) == 0 {
		w.sz -= w.sf.Size()
		w.putSeriesFrame(w.res.Frames[ss].Data.(*datatypes.ReadResponse_Frame_Series))
		w.res.Frames = w.res.Frames[:ss]
	} else if w.sz > writeSize {
		w.Flush()
	}
}

func (w *ResponseWriter) streamStringArrayPoints(cur cursors.StringArrayCursor) {
	w.sf.DataType = datatypes.DataTypeString
	ss := len(w.res.Frames) - 1

	p := w.getStringPointsFrame()
	frame := p.StringPoints
	w.res.Frames = append(w.res.Frames, datatypes.ReadResponse_Frame{Data: p})

	var seriesValueCount = 0
	for {
		// If the number of values produced by cur > 1000,
		// cur.Next() will produce batches of values that are of
		// length ≤ 1000.
		// We attempt to limit the frame Timestamps / Values lengths
		// the same to avoid allocations. These frames are recycled
		// after flushing so that on repeated use there should be enough space
		// to append values from a into frame without additional allocations.
		a := cur.Next()

		if len(a.Timestamps) == 0 {
			break
		}

		seriesValueCount += a.Len()
		// As specified in the struct definition, w.sz is an estimated
		// size (in bytes) of the buffered data. It is therefore a
		// deliberate choice to accumulate using the array Size, which is
		// cheap to calculate. Calling frame.Size() can be expensive
		// when using varint encoding for numbers.
		w.sz += a.Size()

		frame.Timestamps = append(frame.Timestamps, a.Timestamps...)
		frame.Values = append(frame.Values, a.Values...)

		// given the expectation of cur.Next, we attempt to limit
		// the number of values appended to the frame to batchSize (1000)
		needsFrame := len(frame.Timestamps) >= batchSize

		if w.sz >= writeSize {
			needsFrame = true
			w.Flush()
			if w.err != nil {
				break
			}
		}

		if needsFrame {
			// new frames are returned with Timestamps and Values preallocated
			// to a minimum of batchSize length to reduce further allocations.
			p = w.getStringPointsFrame()
			frame = p.StringPoints
			w.res.Frames = append(w.res.Frames, datatypes.ReadResponse_Frame{Data: p})
		}
	}

	w.vc += seriesValueCount
	if seriesValueCount == 0 {
		w.sz -= w.sf.Size()
		w.putSeriesFrame(w.res.Frames[ss].Data.(*datatypes.ReadResponse_Frame_Series))
		w.res.Frames = w.res.Frames[:ss]
	} else if w.sz > writeSize {
		w.Flush()
	}
}

func (w *ResponseWriter) getBooleanPointsFrame() *datatypes.ReadResponse_Frame_BooleanPoints {
	var res *datatypes.ReadResponse_Frame_BooleanPoints
	if len(w.buffer.Boolean) > 0 {
		i := len(w.buffer.Boolean) - 1
		res = w.buffer.Boolean[i]
		w.buffer.Boolean[i] = nil
		w.buffer.Boolean = w.buffer.Boolean[:i]
	} else {
		res = &datatypes.ReadResponse_Frame_BooleanPoints{
			BooleanPoints: &datatypes.ReadResponse_BooleanPointsFrame{
				Timestamps: make([]int64, 0, batchSize),
				Values:     make([]bool, 0, batchSize),
			},
		}
	}
	return res
}

func (w *ResponseWriter) putBooleanPointsFrame(f *datatypes.ReadResponse_Frame_BooleanPoints) {
	f.BooleanPoints.Timestamps = f.BooleanPoints.Timestamps[:0]
	f.BooleanPoints.Values = f.BooleanPoints.Values[:0]
	w.buffer.Boolean = append(w.buffer.Boolean, f)
}

func (w *ResponseWriter) getBooleanValues() *datatypes.ReadResponse_AnyPoints_Booleans {
	var res *datatypes.ReadResponse_AnyPoints_Booleans
	if len(w.buffer.BooleanValues) > 0 {
		i := len(w.buffer.BooleanValues) - 1
		res = w.buffer.BooleanValues[i]
		w.buffer.BooleanValues[i] = nil
		w.buffer.BooleanValues = w.buffer.BooleanValues[:i]
	} else {
		res = &datatypes.ReadResponse_AnyPoints_Booleans{
			Booleans: &datatypes.ReadResponse_BooleanValues{
				Values: make([]bool, 0, batchSize),
			},
		}
	}
	return res
}

func (w *ResponseWriter) putBooleanValues(f *datatypes.ReadResponse_AnyPoints_Booleans) {
	f.Booleans.Values = f.Booleans.Values[:0]
	w.buffer.BooleanValues = append(w.buffer.BooleanValues, f)
}

func (w *ResponseWriter) streamBooleanArraySeries(cur cursors.BooleanArrayCursor) {
	w.sf.DataType = datatypes.DataTypeBoolean
	ss := len(w.res.Frames) - 1
	a := cur.Next()
	if len(a.Timestamps) == 0 {
		w.sz -= w.sf.Size()
		w.putSeriesFrame(w.res.Frames[ss].Data.(*datatypes.ReadResponse_Frame_Series))
		w.res.Frames = w.res.Frames[:ss]
	} else if w.sz > writeSize {
		w.Flush()
	}
}

func (w *ResponseWriter) streamBooleanArrayPoints(cur cursors.BooleanArrayCursor) {
	w.sf.DataType = datatypes.DataTypeBoolean
	ss := len(w.res.Frames) - 1

	p := w.getBooleanPointsFrame()
	frame := p.BooleanPoints
	w.res.Frames = append(w.res.Frames, datatypes.ReadResponse_Frame{Data: p})

	var seriesValueCount = 0
	for {
		// If the number of values produced by cur > 1000,
		// cur.Next() will produce batches of values that are of
		// length ≤ 1000.
		// We attempt to limit the frame Timestamps / Values lengths
		// the same to avoid allocations. These frames are recycled
		// after flushing so that on repeated use there should be enough space
		// to append values from a into frame without additional allocations.
		a := cur.Next()

		if len(a.Timestamps) == 0 {
			break
		}

		seriesValueCount += a.Len()
		// As specified in the struct definition, w.sz is an estimated
		// size (in bytes) of the buffered data. It is therefore a
		// deliberate choice to accumulate using the array Size, which is
		// cheap to calculate. Calling frame.Size() can be expensive
		// when using varint encoding for numbers.
		w.sz += a.Size()

		frame.Timestamps = append(frame.Timestamps, a.Timestamps...)
		frame.Values = append(frame.Values, a.Values...)

		// given the expectation of cur.Next, we attempt to limit
		// the number of values appended to the frame to batchSize (1000)
		needsFrame := len(frame.Timestamps) >= batchSize

		if w.sz >= writeSize {
			needsFrame = true
			w.Flush()
			if w.err != nil {
				break
			}
		}

		if needsFrame {
			// new frames are returned with Timestamps and Values preallocated
			// to a minimum of batchSize length to reduce further allocations.
			p = w.getBooleanPointsFrame()
			frame = p.BooleanPoints
			w.res.Frames = append(w.res.Frames, datatypes.ReadResponse_Frame{Data: p})
		}
	}

	w.vc += seriesValueCount
	if seriesValueCount == 0 {
		w.sz -= w.sf.Size()
		w.putSeriesFrame(w.res.Frames[ss].Data.(*datatypes.ReadResponse_Frame_Series))
		w.res.Frames = w.res.Frames[:ss]
	} else if w.sz > writeSize {
		w.Flush()
	}
}
