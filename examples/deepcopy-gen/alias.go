package deepcopy_gen

// +gen:deepcopy=true
type Slice []int

// +gen:deepcopy=true
type Map map[string]int

// +gen:deepcopy=true
type Struct Builtin

// +gen:deepcopy=true
type StructPointer *Builtin
