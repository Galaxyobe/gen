package deepcopy_gen

// +gen:deepcopy=true
type Builtin struct {
	Byte    byte
	Int16   int16
	Int32   int32
	Int64   int64
	Uint8   uint8
	Uint16  uint16
	Uint32  uint32
	Uint64  uint64
	Float32 float32
	Float64 float64
	String  string
}

type Builtins struct {
	Builtin1 Builtin
	Builtin2 *Builtin
}
