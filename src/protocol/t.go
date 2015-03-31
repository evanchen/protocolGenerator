package protocol

import (
	"fmt"
)

func TEncodeDecode() {
	t1 := uint8(11)
	t11 := encode_uint8(t1)
	t111, _ := decode_uint8(t11)
	fmt.Printf("uint8_encode_decode: %v \n%v \n%v\n\n", t1, t11, t111)
	arr_t1 := []uint8{t1, t1 + 1, t1 + 2}
	arr_t11 := encode_array_uint8(arr_t1)
	arr_t111, _ := decode_array_uint8(arr_t11)
	fmt.Printf("array_uint8_encode_decode: %v \n%v \n%v\n\n", arr_t1, arr_t11, arr_t111)

	t2 := uint16(22)
	t22 := encode_uint16(t2)
	t222, _ := decode_uint16(t22)
	fmt.Printf("uint16_encode_decode: %v \n%v \n%v\n\n", t2, t22, t222)
	arr_t2 := []uint16{t2, t2 + 1, t2 + 2}
	arr_t22 := encode_array_uint16(arr_t2)
	arr_t222, _ := decode_array_uint16(arr_t22)
	fmt.Printf("array_uint16_encode_decode: %v \n%v \n%v\n\n", arr_t2, arr_t22, arr_t222)

	t3 := uint32(33)
	t33 := encode_uint32(t3)
	t333, _ := decode_uint32(t33)
	fmt.Printf("uint32_encode_decode: %v \n%v \n%v\n\n", t3, t33, t333)
	arr_t3 := []uint32{t3, t3 + 1, t3 + 2}
	arr_t33 := encode_array_uint32(arr_t3)
	arr_t333, _ := decode_array_uint32(arr_t33)
	fmt.Printf("array_uint32_encode_decode: %v \n%v \n%v\n\n", arr_t3, arr_t33, arr_t333)

	t4 := uint64(44)
	t44 := encode_uint64(t4)
	t444, _ := decode_uint64(t44)
	fmt.Printf("uint64_encode_decode: %v \n%v \n%v\n\n", t4, t44, t444)
	arr_t4 := []uint64{t4, t4 + 1, t4 + 2}
	arr_t44 := encode_array_uint64(arr_t4)
	arr_t444, _ := decode_array_uint64(arr_t44)
	fmt.Printf("array_uint64_encode_decode: %v \n%v \n%v\n\n", arr_t4, arr_t44, arr_t444)

	t5 := float32(55.555)
	t55 := encode_float32(t5)
	t555, _ := decode_float32(t55)
	fmt.Printf("float32_encode_decode: %v \n%v \n%v\n\n", t5, t55, t555)
	arr_t5 := []float32{t5, t5 + 1, t5 + 2}
	arr_t55 := encode_array_float32(arr_t5)
	arr_t555, _ := decode_array_float32(arr_t55)
	fmt.Printf("array_float32_encode_decode: %v \n%v \n%v\n\n", arr_t5, arr_t55, arr_t555)

	t6 := float64(6666.666)
	t66 := encode_float64(t6)
	t666, _ := decode_float64(t66)
	fmt.Printf("float64_encode_decode: %v \n%v \n%v\n\n", t6, t66, t666)
	arr_t6 := []float64{t6, t6 + 1, t6 + 2}
	arr_t66 := encode_array_float64(arr_t6)
	arr_t666, _ := decode_array_float64(arr_t66)
	fmt.Printf("array_float64_encode_decode: %v \n%v \n%v\n\n", arr_t6, arr_t66, arr_t666)

	t7 := "aha,eehah"
	t77 := encode_string(t7)
	t777, _ := decode_string(t77)
	fmt.Printf("string_encode_decode: %v \n%v \n%v\n\n", t7, t77, t777)
	arr_t7 := []string{"asda", "56tyt", "casd45fsfsbn"}
	arr_t77 := encode_array_string(arr_t7)
	arr_t777, _ := decode_array_string(arr_t77)
	fmt.Printf("array_string_encode_decode: %v \n%v \n%v\n\n", arr_t7, arr_t77, arr_t777)

	obj := CreateRoleInfo()
	obj.Id = uint64(11)
	obj.Name = "ohmyname"
	obj.Age = uint8(22)
	obj.Lv = uint16(45)
	obj.Money = float32(919.22)
	obj.Exp = uint32(81)
	obj.Soul = float64(555.16)
	obj.Unarray = []uint8{1, 2, 3}
	obj.Skill = []string{"s1", "s2", "s3"}

	arr_obj := CreateRoleInfoList()
	for i := 0; i < 3; i++ {
		arr_obj.RoleList = append(arr_obj.RoleList, *obj)
	}

	EncodeData := arr_obj.Marshal()
	decodeData := CreateRoleInfoList()
	decodeData, EncodeData = decodeData.Unmarshal(EncodeData)
	for k, v := range decodeData.RoleList {
		fmt.Printf("obj%d: %v\n", k, v)
	}
}
