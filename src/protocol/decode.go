package protocol

import (
	"math"
)

func decode_panic() {
	panic("decode error!")
}

func decode_uint8(Data []byte) (uint8, []byte) {
	size := len(Data)
	rs := 1
	if size < rs {
		decode_panic()
	}

	ret := uint8(Data[rs-1])
	return ret, Data[rs:]
}

func decode_uint16(Data []byte) (uint16, []byte) {
	size := len(Data)
	rs := 2
	if size < rs {
		decode_panic()
	}

	ret := uint16(Data[rs-2])
	ret |= uint16(Data[rs-1]) << 8
	return ret, Data[rs:]
}

func decode_uint32(Data []byte) (uint32, []byte) {
	size := len(Data)
	rs := 4
	if size < rs {
		decode_panic()
	}

	ret := uint32(Data[rs-4])
	ret |= uint32(Data[rs-3]) << 8
	ret |= uint32(Data[rs-2]) << 16
	ret |= uint32(Data[rs-1]) << 24
	return ret, Data[rs:]
}

func decode_uint64(Data []byte) (uint64, []byte) {
	size := len(Data)
	rs := 8
	if size < rs {
		decode_panic()
	}

	ret := uint64(Data[rs-8])
	ret |= uint64(Data[rs-7]) << 8
	ret |= uint64(Data[rs-6]) << 16
	ret |= uint64(Data[rs-5]) << 24
	ret |= uint64(Data[rs-4]) << 32
	ret |= uint64(Data[rs-3]) << 40
	ret |= uint64(Data[rs-2]) << 48
	ret |= uint64(Data[rs-1]) << 56
	return ret, Data[rs:]
}

func decode_float32(Data []byte) (float32, []byte) {
	var tmp uint32
	tmp, Data = decode_uint32(Data)
	ret := math.Float32frombits(tmp)
	return ret, Data
}

func decode_float64(Data []byte) (float64, []byte) {
	var tmp uint64
	tmp, Data = decode_uint64(Data)
	ret := math.Float64frombits(tmp)
	return ret, Data
}

func decode_string(Data []byte) (string, []byte) {
	var rs uint16
	rs, Data = decode_uint16(Data)
	size := len(Data)
	if size < int(rs) {
		decode_panic()
	}

	return string(Data[:rs]), Data[rs:]
}

func decode_array_uint8(Data []byte) ([]uint8, []byte) {
	var rs uint16
	rs, Data = decode_uint16(Data)
	ret := make([]uint8, 0, rs)
	var tmp uint8
	for i := uint16(0); i < rs; i++ {
		tmp, Data = decode_uint8(Data)
		ret = append(ret, tmp)
	}
	return ret, Data
}

func decode_array_uint16(Data []byte) ([]uint16, []byte) {
	var rs uint16
	rs, Data = decode_uint16(Data)
	ret := make([]uint16, 0, rs)
	var tmp uint16
	for i := uint16(0); i < rs; i++ {
		tmp, Data = decode_uint16(Data)
		ret = append(ret, tmp)
	}
	return ret, Data
}

func decode_array_uint32(Data []byte) ([]uint32, []byte) {
	var rs uint16
	rs, Data = decode_uint16(Data)
	ret := make([]uint32, 0, rs)
	var tmp uint32
	for i := uint16(0); i < rs; i++ {
		tmp, Data = decode_uint32(Data)
		ret = append(ret, tmp)
	}
	return ret, Data
}

func decode_array_uint64(Data []byte) ([]uint64, []byte) {
	var rs uint16
	rs, Data = decode_uint16(Data)
	ret := make([]uint64, 0, rs)
	var tmp uint64
	for i := uint16(0); i < rs; i++ {
		tmp, Data = decode_uint64(Data)
		ret = append(ret, tmp)
	}
	return ret, Data
}

func decode_array_float32(Data []byte) ([]float32, []byte) {
	var rs uint16
	rs, Data = decode_uint16(Data)
	ret := make([]float32, 0, rs)
	var tmp float32
	for i := uint16(0); i < rs; i++ {
		tmp, Data = decode_float32(Data)
		ret = append(ret, tmp)
	}
	return ret, Data
}

func decode_array_float64(Data []byte) ([]float64, []byte) {
	var rs uint16
	rs, Data = decode_uint16(Data)
	ret := make([]float64, 0, rs)
	var tmp float64
	for i := uint16(0); i < rs; i++ {
		tmp, Data = decode_float64(Data)
		ret = append(ret, tmp)
	}
	return ret, Data
}

func decode_array_string(Data []byte) ([]string, []byte) {
	var rs uint16
	rs, Data = decode_uint16(Data)
	ret := make([]string, 0, rs)
	var tmp string
	for i := uint16(0); i < rs; i++ {
		tmp, Data = decode_string(Data)
		ret = append(ret, tmp)
	}
	return ret, Data
}
