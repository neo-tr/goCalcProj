// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.3
// source: calculator.proto

package pb

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Operand defines a value that can be either a number or a variable reference
type Operand struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Types that are valid to be assigned to Value:
	//
	//	*Operand_Number
	//	*Operand_Variable
	Value         isOperand_Value `protobuf_oneof:"value"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Operand) Reset() {
	*x = Operand{}
	mi := &file_calculator_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Operand) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Operand) ProtoMessage() {}

func (x *Operand) ProtoReflect() protoreflect.Message {
	mi := &file_calculator_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Operand.ProtoReflect.Descriptor instead.
func (*Operand) Descriptor() ([]byte, []int) {
	return file_calculator_proto_rawDescGZIP(), []int{0}
}

func (x *Operand) GetValue() isOperand_Value {
	if x != nil {
		return x.Value
	}
	return nil
}

func (x *Operand) GetNumber() int64 {
	if x != nil {
		if x, ok := x.Value.(*Operand_Number); ok {
			return x.Number
		}
	}
	return 0
}

func (x *Operand) GetVariable() string {
	if x != nil {
		if x, ok := x.Value.(*Operand_Variable); ok {
			return x.Variable
		}
	}
	return ""
}

type isOperand_Value interface {
	isOperand_Value()
}

type Operand_Number struct {
	// Numeric value for calculations (e.g., 10), minimum 0.
	Number int64 `protobuf:"varint,1,opt,name=number,proto3,oneof"`
}

type Operand_Variable struct {
	// Single-letter variable name (e.g., "x").
	Variable string `protobuf:"bytes,2,opt,name=variable,proto3,oneof"`
}

func (*Operand_Number) isOperand_Value() {}

func (*Operand_Variable) isOperand_Value() {}

// Instruction defines a calc or print operation
type Instruction struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Type          string                 `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`   // Operation type: "calc" or "print"
	Op            string                 `protobuf:"bytes,2,opt,name=op,proto3" json:"op,omitempty"`       // Operator for calc: "+", "-", "*"
	Var           string                 `protobuf:"bytes,3,opt,name=var,proto3" json:"var,omitempty"`     // Variable name (e.g., "x")
	Left          *Operand               `protobuf:"bytes,4,opt,name=left,proto3" json:"left,omitempty"`   // Left operand (number or variable)
	Right         *Operand               `protobuf:"bytes,5,opt,name=right,proto3" json:"right,omitempty"` // Right operand (number or variable)
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Instruction) Reset() {
	*x = Instruction{}
	mi := &file_calculator_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Instruction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Instruction) ProtoMessage() {}

func (x *Instruction) ProtoReflect() protoreflect.Message {
	mi := &file_calculator_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Instruction.ProtoReflect.Descriptor instead.
func (*Instruction) Descriptor() ([]byte, []int) {
	return file_calculator_proto_rawDescGZIP(), []int{1}
}

func (x *Instruction) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *Instruction) GetOp() string {
	if x != nil {
		return x.Op
	}
	return ""
}

func (x *Instruction) GetVar() string {
	if x != nil {
		return x.Var
	}
	return ""
}

func (x *Instruction) GetLeft() *Operand {
	if x != nil {
		return x.Left
	}
	return nil
}

func (x *Instruction) GetRight() *Operand {
	if x != nil {
		return x.Right
	}
	return nil
}

// ResultItem holds the output for print instructions
type ResultItem struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Var           string                 `protobuf:"bytes,1,opt,name=var,proto3" json:"var,omitempty"`      // Variable name
	Value         int64                  `protobuf:"varint,2,opt,name=value,proto3" json:"value,omitempty"` // Computed value
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ResultItem) Reset() {
	*x = ResultItem{}
	mi := &file_calculator_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ResultItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResultItem) ProtoMessage() {}

func (x *ResultItem) ProtoReflect() protoreflect.Message {
	mi := &file_calculator_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResultItem.ProtoReflect.Descriptor instead.
func (*ResultItem) Descriptor() ([]byte, []int) {
	return file_calculator_proto_rawDescGZIP(), []int{2}
}

func (x *ResultItem) GetVar() string {
	if x != nil {
		return x.Var
	}
	return ""
}

func (x *ResultItem) GetValue() int64 {
	if x != nil {
		return x.Value
	}
	return 0
}

// ProcessInstructionsRequest contains a list of instructions to process
type ProcessInstructionsRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Instructions  []*Instruction         `protobuf:"bytes,1,rep,name=instructions,proto3" json:"instructions,omitempty"` // Array of instructions
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ProcessInstructionsRequest) Reset() {
	*x = ProcessInstructionsRequest{}
	mi := &file_calculator_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ProcessInstructionsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProcessInstructionsRequest) ProtoMessage() {}

func (x *ProcessInstructionsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_calculator_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProcessInstructionsRequest.ProtoReflect.Descriptor instead.
func (*ProcessInstructionsRequest) Descriptor() ([]byte, []int) {
	return file_calculator_proto_rawDescGZIP(), []int{3}
}

func (x *ProcessInstructionsRequest) GetInstructions() []*Instruction {
	if x != nil {
		return x.Instructions
	}
	return nil
}

// ProcessInstructionsResponse contains the results of print instructions
type ProcessInstructionsResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Items         []*ResultItem          `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"` // Array of print results
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ProcessInstructionsResponse) Reset() {
	*x = ProcessInstructionsResponse{}
	mi := &file_calculator_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ProcessInstructionsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProcessInstructionsResponse) ProtoMessage() {}

func (x *ProcessInstructionsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_calculator_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProcessInstructionsResponse.ProtoReflect.Descriptor instead.
func (*ProcessInstructionsResponse) Descriptor() ([]byte, []int) {
	return file_calculator_proto_rawDescGZIP(), []int{4}
}

func (x *ProcessInstructionsResponse) GetItems() []*ResultItem {
	if x != nil {
		return x.Items
	}
	return nil
}

var File_calculator_proto protoreflect.FileDescriptor

const file_calculator_proto_rawDesc = "" +
	"\n" +
	"\x10calculator.proto\x12\rcalculator.v1\x1a\x1cgoogle/api/annotations.proto\"J\n" +
	"\aOperand\x12\x18\n" +
	"\x06number\x18\x01 \x01(\x03H\x00R\x06number\x12\x1c\n" +
	"\bvariable\x18\x02 \x01(\tH\x00R\bvariableB\a\n" +
	"\x05value\"\x9d\x01\n" +
	"\vInstruction\x12\x12\n" +
	"\x04type\x18\x01 \x01(\tR\x04type\x12\x0e\n" +
	"\x02op\x18\x02 \x01(\tR\x02op\x12\x10\n" +
	"\x03var\x18\x03 \x01(\tR\x03var\x12*\n" +
	"\x04left\x18\x04 \x01(\v2\x16.calculator.v1.OperandR\x04left\x12,\n" +
	"\x05right\x18\x05 \x01(\v2\x16.calculator.v1.OperandR\x05right\"4\n" +
	"\n" +
	"ResultItem\x12\x10\n" +
	"\x03var\x18\x01 \x01(\tR\x03var\x12\x14\n" +
	"\x05value\x18\x02 \x01(\x03R\x05value\"\\\n" +
	"\x1aProcessInstructionsRequest\x12>\n" +
	"\finstructions\x18\x01 \x03(\v2\x1a.calculator.v1.InstructionR\finstructions\"N\n" +
	"\x1bProcessInstructionsResponse\x12/\n" +
	"\x05items\x18\x01 \x03(\v2\x19.calculator.v1.ResultItemR\x05items2\xa3\x01\n" +
	"\x11CalculatorService\x12\x8d\x01\n" +
	"\x13ProcessInstructions\x12).calculator.v1.ProcessInstructionsRequest\x1a*.calculator.v1.ProcessInstructionsResponse\"\x1f\x82\xd3\xe4\x93\x02\x19:\x01*\"\x14/api/v1/instructionsB\x1eZ\x1cgithub.com/goCalcProj/gen/pbb\x06proto3"

var (
	file_calculator_proto_rawDescOnce sync.Once
	file_calculator_proto_rawDescData []byte
)

func file_calculator_proto_rawDescGZIP() []byte {
	file_calculator_proto_rawDescOnce.Do(func() {
		file_calculator_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_calculator_proto_rawDesc), len(file_calculator_proto_rawDesc)))
	})
	return file_calculator_proto_rawDescData
}

var file_calculator_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_calculator_proto_goTypes = []any{
	(*Operand)(nil),                     // 0: calculator.v1.Operand
	(*Instruction)(nil),                 // 1: calculator.v1.Instruction
	(*ResultItem)(nil),                  // 2: calculator.v1.ResultItem
	(*ProcessInstructionsRequest)(nil),  // 3: calculator.v1.ProcessInstructionsRequest
	(*ProcessInstructionsResponse)(nil), // 4: calculator.v1.ProcessInstructionsResponse
}
var file_calculator_proto_depIdxs = []int32{
	0, // 0: calculator.v1.Instruction.left:type_name -> calculator.v1.Operand
	0, // 1: calculator.v1.Instruction.right:type_name -> calculator.v1.Operand
	1, // 2: calculator.v1.ProcessInstructionsRequest.instructions:type_name -> calculator.v1.Instruction
	2, // 3: calculator.v1.ProcessInstructionsResponse.items:type_name -> calculator.v1.ResultItem
	3, // 4: calculator.v1.CalculatorService.ProcessInstructions:input_type -> calculator.v1.ProcessInstructionsRequest
	4, // 5: calculator.v1.CalculatorService.ProcessInstructions:output_type -> calculator.v1.ProcessInstructionsResponse
	5, // [5:6] is the sub-list for method output_type
	4, // [4:5] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_calculator_proto_init() }
func file_calculator_proto_init() {
	if File_calculator_proto != nil {
		return
	}
	file_calculator_proto_msgTypes[0].OneofWrappers = []any{
		(*Operand_Number)(nil),
		(*Operand_Variable)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_calculator_proto_rawDesc), len(file_calculator_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_calculator_proto_goTypes,
		DependencyIndexes: file_calculator_proto_depIdxs,
		MessageInfos:      file_calculator_proto_msgTypes,
	}.Build()
	File_calculator_proto = out.File
	file_calculator_proto_goTypes = nil
	file_calculator_proto_depIdxs = nil
}
