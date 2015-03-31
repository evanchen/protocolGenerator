package protocol

import (
	"math"
)

func encode_uint8(val uint8) []byte {
	buf := make([]byte, 0, 1)
	buf = append(buf, val)
	return buf
}

func encode_uint16(val uint16) []byte {
	buf := make([]byte, 0, 2)
	buf = append(buf,
		uint8(val),
		uint8(val>>8))
	return buf
}

func encode_uint32(val uint32) []byte {
	buf := make([]byte, 0, 4)
	buf = append(buf,
		uint8(val),
		uint8(val>>8),
		uint8(val>>16),
		uint8(val>>24))
	return buf
}

func encode_uint64(val uint64) []byte {
	buf := make([]byte, 0, 8)
	buf = append(buf,
		uint8(val),
		uint8(val>>8),
		uint8(val>>16),
		uint8(val>>24),
		uint8(val>>32),
		uint8(val>>40),
		uint8(val>>48),
		uint8(val>>56))
	return buf
}

func encode_float32(val float32) []byte {
	tmp := math.Float32bits(val)
	return encode_uint32(tmp)
}

func encode_float64(val float64) []byte {
	tmp := math.Float64bits(val)
	return encode_uint64(tmp)
}

func encode_string(val string) []byte {
	size := uint16(len(val))
	buf := encode_uint16(size)
	buf = append(buf, val...)
	return buf
}

func encode_array_uint8(val []uint8) []byte {
	size := uint16(len(val))
	buf := encode_uint16(size)
	for _, v := range val {
		buf = append(buf, v)
	}
	return buf
}

func encode_array_uint16(val []uint16) []byte {
	size := uint16(len(val))
	buf := encode_uint16(size)
	var en []byte
	for _, v := range val {
		en = encode_uint16(v)
		buf = append(buf, en...)
	}
	return buf
}

func encode_array_uint32(val []uint32) []byte {
	size := uint16(len(val))
	buf := encode_uint16(size)
	var en []byte
	for _, v := range val {
		en = encode_uint32(v)
		buf = append(buf, en...)
	}
	return buf
}

func encode_array_uint64(val []uint64) []byte {
	size := uint16(len(val))
	buf := encode_uint16(size)
	var en []byte
	for _, v := range val {
		en = encode_uint64(v)
		buf = append(buf, en...)
	}
	return buf
}

func encode_array_float32(val []float32) []byte {
	size := uint16(len(val))
	buf := encode_uint16(size)
	var en []byte
	for _, v := range val {
		en = encode_float32(v)
		buf = append(buf, en...)
	}
	return buf
}

func encode_array_float64(val []float64) []byte {
	size := uint16(len(val))
	buf := encode_uint16(size)
	var en []byte
	for _, v := range val {
		en = encode_float64(v)
		buf = append(buf, en...)
	}
	return buf
}

func encode_array_string(val []string) []byte {
	size := uint16(len(val))
	buf := encode_uint16(size)
	var en []byte
	for _, v := range val {
		en = encode_string(v)
		buf = append(buf, en...)
	}
	return buf
}
