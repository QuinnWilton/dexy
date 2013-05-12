package dexy

import (
	"bytes"
	"encoding/binary"
	"os"
)

type DexFile struct {
	Header    HeaderItem
	StringIds []StringIdItem
	TypeIds   []TypeIdItem
	ProtoIds  []ProtoIdItem
	FieldIds  []FieldIdItem
	MethodIds []MethodIdItem
	ClassDefs []ClassDefItem
	Data      Data
	LinkData  []uint8
}

type HeaderItem struct {
	Magic         [8]uint8
	Checksum      uint32
	Signature     [20]uint8
	FileSize      uint32
	HeaderSize    uint32
	EndianTag     uint32
	LinkSize      uint32
	LinkOff       uint32
	MapOff        uint32
	StringIdsSize uint32
	StringIdsOff  uint32
	TypeIdsSize   uint32
	TypeIdsOff    uint32
	ProtoIdsSize  uint32
	ProtoIdsOff   uint32
	FieldIdsSize  uint32
	FieldIdsOff   uint32
	MethodIdsSize uint32
	MethodIdsOff  uint32
	ClassDefsSize uint32
	ClassDefsOff  uint32
	DataSize      uint32
	DataOff       uint32
}

type Data struct {
	Map             MapList
	StringDataItems []StringDataItem
	ClassDataItems  []ClassDataItem
}

type MapList struct {
	Size uint32
	List []MapItem
}

type MapItem struct {
	Type   uint16
	Unused uint16
	Size   uint32
	Offset uint32
}

type StringIdItem struct {
	StringDataOff uint32
}

type StringDataItem struct {
	Utf16Size uint32
	Raw       []byte
	Decoded   string
}

type TypeIdItem struct {
	DescriptorIdx uint32
}

type ProtoIdItem struct {
	ShortyIdx     uint32
	ReturnTypeIdx uint32
	ParametersOff uint32
}

type FieldIdItem struct {
	ClassIdx uint16
	TypeIdx  uint16
	NameIdx  uint32
}

type MethodIdItem struct {
	ClassIdx uint16
	ProtoIdx uint16
	NameIdx  uint32
}

type ClassDefItem struct {
	ClassIdx        uint32
	AccessFlags     uint32
	SuperclassIdx   uint32
	InterfacesOff   uint32
	SourceFileIdx   uint32
	AnnotationsOff  uint32
	ClassDataOff    uint32
	StaticValuesOff uint32
}

type ClassDataItem struct {
	StaticFieldsSize   uint32
	InstanceFieldsSize uint32
	DirectMethodsSize  uint32
	VirtualMethodsSize uint32
	//StaticFields []encodedField
	//InstanceFields []encodedField
	//DirectMethods []encodedMethod
	//VirtualMethods []encodedMethod
}

func NewDex(b []byte) DexFile {
	r := bytes.NewReader(b)
	var dexFile DexFile
	binary.Read(r, binary.LittleEndian, &dexFile.Header)

	dexFile.StringIds = make([]StringIdItem, dexFile.Header.StringIdsSize)
	binary.Read(r, binary.LittleEndian, &dexFile.StringIds)

	dexFile.TypeIds = make([]TypeIdItem, dexFile.Header.TypeIdsSize)
	binary.Read(r, binary.LittleEndian, &dexFile.TypeIds)

	dexFile.ProtoIds = make([]ProtoIdItem, dexFile.Header.ProtoIdsSize)
	binary.Read(r, binary.LittleEndian, &dexFile.ProtoIds)

	dexFile.FieldIds = make([]FieldIdItem, dexFile.Header.FieldIdsSize)
	binary.Read(r, binary.LittleEndian, &dexFile.FieldIds)

	dexFile.MethodIds = make([]MethodIdItem, dexFile.Header.MethodIdsSize)
	binary.Read(r, binary.LittleEndian, &dexFile.MethodIds)

	dexFile.ClassDefs = make([]ClassDefItem, dexFile.Header.ClassDefsSize)
	binary.Read(r, binary.LittleEndian, &dexFile.ClassDefs)

	readMap(r, &dexFile)
	readStrings(r, &dexFile)

	return dexFile
}

func readMap(r *bytes.Reader, dexFile *DexFile) {
	r.Seek(int64(dexFile.Header.MapOff), os.SEEK_SET)
	binary.Read(r, binary.LittleEndian, &dexFile.Data.Map.Size)

	dexFile.Data.Map.List = make([]MapItem, dexFile.Data.Map.Size)
	binary.Read(r, binary.LittleEndian, &dexFile.Data.Map.List)
}

func readStrings(r *bytes.Reader, dexFile *DexFile) {
	dexFile.Data.StringDataItems = make([]StringDataItem, dexFile.Header.StringIdsSize)

	for i, stringIdItem := range dexFile.StringIds {
		r.Seek(int64(stringIdItem.StringDataOff), os.SEEK_SET)
		size, _ := binary.ReadUvarint(r)
		dexFile.Data.StringDataItems[i].Utf16Size = uint32(size)

		var buf bytes.Buffer
		var b byte = 0x1
		for b != 0x0 {
			b, _ = r.ReadByte()
			buf.WriteByte(b)
		}
		dexFile.Data.StringDataItems[i].Raw = buf.Bytes()

		dexFile.Data.StringDataItems[i].Decoded, _ = Mutf8(dexFile.Data.StringDataItems[i].Raw)
	}
}
