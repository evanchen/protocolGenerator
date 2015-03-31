package protocol

import (
	_ "errors"
)

///////////////////////////////////////////////////////////
type RoleInfo struct {
	Id      uint64
	Name    string
	Age     uint8
	Lv      uint16
	Money   float32
	Exp     uint32
	Soul    float64
	Unarray []uint8
	Skill   []string
}

func CreateRoleInfo() *RoleInfo {
	obj := &RoleInfo{}
	return obj
}

func (this *RoleInfo) Marshal() []byte {
	buf := make([]byte, 0, 16)
	buf = append(buf, encode_uint64(this.Id)...)
	buf = append(buf, encode_string(this.Name)...)
	buf = append(buf, encode_uint8(this.Age)...)
	buf = append(buf, encode_uint16(this.Lv)...)
	buf = append(buf, encode_float32(this.Money)...)
	buf = append(buf, encode_uint32(this.Exp)...)
	buf = append(buf, encode_float64(this.Soul)...)
	buf = append(buf, encode_array_uint8(this.Unarray)...)
	buf = append(buf, encode_array_string(this.Skill)...)
	return buf
}

func (this *RoleInfo) Unmarshal(Data []byte) (*RoleInfo, []byte) {
	this.Id, Data = decode_uint64(Data)
	this.Name, Data = decode_string(Data)
	this.Age, Data = decode_uint8(Data)
	this.Lv, Data = decode_uint16(Data)
	this.Money, Data = decode_float32(Data)
	this.Exp, Data = decode_uint32(Data)
	this.Soul, Data = decode_float64(Data)
	this.Unarray, Data = decode_array_uint8(Data)
	this.Skill, Data = decode_array_string(Data)
	return this, Data
}

///////////////////////////////////////////////////////////

type RoleInfoList struct {
	RoleList []RoleInfo
}

func CreateRoleInfoList() *RoleInfoList {
	obj := &RoleInfoList{}
	return obj
}

func (this *RoleInfoList) Marshal() []byte {
	buf := make([]byte, 0, 16)
	buf = append(buf, encode_array_roleInfo(this.RoleList)...)
	return buf
}

func (this *RoleInfoList) Unmarshal(Data []byte) (*RoleInfoList, []byte) {
	this.RoleList, Data = decode_array_roleInfo(Data)
	return this, Data
}

func encode_array_roleInfo(RoleList []RoleInfo) []byte {
	buf := make([]byte, 0, 16)
	size := uint16(len(RoleList))
	buf = append(buf, encode_uint16(size)...)
	for _, j := range RoleList {
		buf = append(buf, j.Marshal()...)
	}
	return buf
}

func decode_array_roleInfo(Data []byte) ([]RoleInfo, []byte) {
	var size uint16
	size, Data = decode_uint16(Data)
	RoleList := make([]RoleInfo, 0, size)
	var obj *RoleInfo
	for i := uint16(0); i < size; i++ {
		obj = &RoleInfo{}
		obj, Data = obj.Unmarshal(Data)
		RoleList = append(RoleList, *obj)
	}
	return RoleList, Data
}

///////////////////////////////////////////////////////////
