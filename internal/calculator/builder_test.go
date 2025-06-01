package calculator

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"google.golang.org/protobuf/proto"

	pb "github.com/goCalcProj/gen/pb"
)

func TestBuilder_ProcessInstructions(t *testing.T) {
	tests := []struct {
		name    string
		input   *pb.ProcessInstructionsRequest
		want    *pb.ProcessInstructionsResponse
		wantErr bool
		setup   func(context.CancelFunc, *sync.WaitGroup)
	}{
		{
			name: "Simple_addition",
			input: &pb.ProcessInstructionsRequest{
				Instructions: []*pb.Instruction{
					{Type: "calc", Var: "x", Op: "+", Left: &pb.Operand{Value: &pb.Operand_Number{Number: 10}}, Right: &pb.Operand{Value: &pb.Operand_Number{Number: 2}}},
					{Type: "print", Var: "x"},
				},
			},
			want:    &pb.ProcessInstructionsResponse{Items: []*pb.ResultItem{{Var: "x", Value: 12}}},
			wantErr: false,
		},
		{
			name: "Variable_dependencies",
			input: &pb.ProcessInstructionsRequest{
				Instructions: []*pb.Instruction{
					{Type: "calc", Var: "x", Op: "+", Left: &pb.Operand{Value: &pb.Operand_Number{Number: 10}}, Right: &pb.Operand{Value: &pb.Operand_Number{Number: 2}}},
					{Type: "calc", Var: "y", Op: "-", Left: &pb.Operand{Value: &pb.Operand_Variable{Variable: "x"}}, Right: &pb.Operand{Value: &pb.Operand_Number{Number: 3}}},
					{Type: "print", Var: "y"},
				},
			},
			want:    &pb.ProcessInstructionsResponse{Items: []*pb.ResultItem{{Var: "y", Value: 9}}},
			wantErr: false,
		},
		{
			name: "Immutable_variable_error",
			input: &pb.ProcessInstructionsRequest{
				Instructions: []*pb.Instruction{
					{Type: "calc", Var: "x", Op: "+", Left: &pb.Operand{Value: &pb.Operand_Number{Number: 10}}, Right: &pb.Operand{Value: &pb.Operand_Number{Number: 2}}},
					{Type: "calc", Var: "x", Op: "+", Left: &pb.Operand{Value: &pb.Operand_Number{Number: 5}}, Right: &pb.Operand{Value: &pb.Operand_Number{Number: 1}}},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Parallel_execution",
			input: &pb.ProcessInstructionsRequest{
				Instructions: []*pb.Instruction{
					{Type: "calc", Var: "x", Op: "+", Left: &pb.Operand{Value: &pb.Operand_Number{Number: 10}}, Right: &pb.Operand{Value: &pb.Operand_Number{Number: 2}}},
					{Type: "calc", Var: "y", Op: "-", Left: &pb.Operand{Value: &pb.Operand_Variable{Variable: "x"}}, Right: &pb.Operand{Value: &pb.Operand_Number{Number: 3}}},
					{Type: "calc", Var: "z", Op: "*", Left: &pb.Operand{Value: &pb.Operand_Variable{Variable: "x"}}, Right: &pb.Operand{Value: &pb.Operand_Variable{Variable: "y"}}},
					{Type: "calc", Var: "w", Op: "*", Left: &pb.Operand{Value: &pb.Operand_Variable{Variable: "z"}}, Right: &pb.Operand{Value: &pb.Operand_Number{Number: 0}}},
					{Type: "print", Var: "w"},
				},
			},
			want:    &pb.ProcessInstructionsResponse{Items: []*pb.ResultItem{{Var: "w", Value: 0}}},
			wantErr: false,
		},
		{
			name: "Print_order",
			input: &pb.ProcessInstructionsRequest{
				Instructions: []*pb.Instruction{
					{Type: "calc", Var: "x", Op: "+", Left: &pb.Operand{Value: &pb.Operand_Number{Number: 10}}, Right: &pb.Operand{Value: &pb.Operand_Number{Number: 2}}},
					{Type: "calc", Var: "y", Op: "-", Left: &pb.Operand{Value: &pb.Operand_Variable{Variable: "x"}}, Right: &pb.Operand{Value: &pb.Operand_Number{Number: 3}}},
					{Type: "print", Var: "x"},
					{Type: "print", Var: "y"},
				},
			},
			want:    &pb.ProcessInstructionsResponse{Items: []*pb.ResultItem{{Var: "x", Value: 12}, {Var: "y", Value: 9}}},
			wantErr: false,
		},
		{
			name: "Invalid_operation",
			input: &pb.ProcessInstructionsRequest{
				Instructions: []*pb.Instruction{
					{Type: "calc", Var: "x", Op: "/", Left: &pb.Operand{Value: &pb.Operand_Number{Number: 10}}, Right: &pb.Operand{Value: &pb.Operand_Number{Number: 2}}},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Unknown_variable_in_print",
			input: &pb.ProcessInstructionsRequest{
				Instructions: []*pb.Instruction{
					{Type: "print", Var: "x"},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Cyclic_dependency",
			input: &pb.ProcessInstructionsRequest{
				Instructions: []*pb.Instruction{
					{Type: "calc", Var: "x", Op: "+", Left: &pb.Operand{Value: &pb.Operand_Variable{Variable: "y"}}, Right: &pb.Operand{Value: &pb.Operand_Number{Number: 1}}},
					{Type: "calc", Var: "y", Op: "+", Left: &pb.Operand{Value: &pb.Operand_Variable{Variable: "x"}}, Right: &pb.Operand{Value: &pb.Operand_Number{Number: 1}}},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Empty_input",
			input: &pb.ProcessInstructionsRequest{
				Instructions: []*pb.Instruction{},
			},
			want:    &pb.ProcessInstructionsResponse{},
			wantErr: false,
		},
		{
			name: "Invalid_operand_type",
			input: &pb.ProcessInstructionsRequest{
				Instructions: []*pb.Instruction{
					{Type: "calc", Var: "x", Op: "+", Left: &pb.Operand{}, Right: &pb.Operand{Value: &pb.Operand_Number{Number: 2}}},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Nil_left_operand",
			input: &pb.ProcessInstructionsRequest{
				Instructions: []*pb.Instruction{
					{Type: "calc", Var: "x", Op: "+", Left: nil, Right: &pb.Operand{Value: &pb.Operand_Number{Number: 2}}},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Nil_right_operand",
			input: &pb.ProcessInstructionsRequest{
				Instructions: []*pb.Instruction{
					{Type: "calc", Var: "x", Op: "+", Left: &pb.Operand{Value: &pb.Operand_Number{Number: 10}}, Right: nil},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Non-existing_right_variable",
			input: &pb.ProcessInstructionsRequest{
				Instructions: []*pb.Instruction{
					{Type: "calc", Var: "x", Op: "+", Left: &pb.Operand{Value: &pb.Operand_Number{Number: 10}}, Right: &pb.Operand{Value: &pb.Operand_Variable{Variable: "y"}}},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Variable_not_found_in_print",
			input: &pb.ProcessInstructionsRequest{
				Instructions: []*pb.Instruction{
					{Type: "calc", Var: "x", Op: "+", Left: &pb.Operand{Value: &pb.Operand_Number{Number: 10}}, Right: &pb.Operand{Value: &pb.Operand_Variable{Variable: "y"}}},
					{Type: "print", Var: "x"},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Right_variable_dependency_wait",
			input: &pb.ProcessInstructionsRequest{
				Instructions: []*pb.Instruction{
					{Type: "calc", Var: "y", Op: "+", Left: &pb.Operand{Value: &pb.Operand_Number{Number: 5}}, Right: &pb.Operand{Value: &pb.Operand_Number{Number: 3}}},
					{Type: "calc", Var: "x", Op: "*", Left: &pb.Operand{Value: &pb.Operand_Number{Number: 10}}, Right: &pb.Operand{Value: &pb.Operand_Variable{Variable: "y"}}},
					{Type: "print", Var: "x"},
				},
			},
			want:    &pb.ProcessInstructionsResponse{Items: []*pb.ResultItem{{Var: "x", Value: 80}}},
			wantErr: false,
		},
		{
			name: "No_result_in_print_order",
			input: &pb.ProcessInstructionsRequest{
				Instructions: []*pb.Instruction{
					{Type: "calc", Var: "x", Op: "+", Left: &pb.Operand{Value: &pb.Operand_Number{Number: 10}}, Right: &pb.Operand{Value: &pb.Operand_Number{Number: 2}}},
					{Type: "print", Var: "y"},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Cancel_during_calc",
			input: &pb.ProcessInstructionsRequest{
				Instructions: []*pb.Instruction{
					{Type: "calc", Var: "x", Op: "+", Left: &pb.Operand{Value: &pb.Operand_Number{Number: 10}}, Right: &pb.Operand{Value: &pb.Operand_Variable{Variable: "y"}}},
				},
			},
			want:    nil,
			wantErr: true,
			setup: func(cancel context.CancelFunc, wg *sync.WaitGroup) {
				go func() {
					fmt.Printf("DEBUG: Cancelling context during calc\n")
					cancel()
					wg.Done()
				}()
			},
		},
		{
			name: "Cancel_during_print",
			input: &pb.ProcessInstructionsRequest{
				Instructions: []*pb.Instruction{
					{Type: "calc", Var: "x", Op: "+", Left: &pb.Operand{Value: &pb.Operand_Number{Number: 10}}, Right: &pb.Operand{Value: &pb.Operand_Variable{Variable: "y"}}},
					{Type: "print", Var: "x"},
				},
			},
			want:    nil,
			wantErr: true,
			setup: func(cancel context.CancelFunc, wg *sync.WaitGroup) {
				go func() {
					fmt.Printf("DEBUG: Cancelling context during print\n")
					cancel()
					wg.Done()
				}()
			},
		},
		{
			name: "Variable_not_found_in_print_explicit",
			input: &pb.ProcessInstructionsRequest{
				Instructions: []*pb.Instruction{
					{Type: "calc", Var: "x", Op: "+", Left: &pb.Operand{Value: &pb.Operand_Number{Number: 10}}, Right: &pb.Operand{Value: &pb.Operand_Variable{Variable: "y"}}},
					{Type: "print", Var: "z"},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Cancel_during_leftVar_wait",
			input: &pb.ProcessInstructionsRequest{
				Instructions: func() []*pb.Instruction {
					instrs := make([]*pb.Instruction, 100)
					for i := 0; i < 100; i++ {
						varName := fmt.Sprintf("var%d", i+1)
						var leftOperand *pb.Operand
						if i == 0 {
							leftOperand = &pb.Operand{Value: &pb.Operand_Number{Number: 10}}
						} else {
							leftOperand = &pb.Operand{Value: &pb.Operand_Variable{Variable: fmt.Sprintf("var%d", i)}}
						}
						instrs[i] = &pb.Instruction{
							Type:  "calc",
							Var:   varName,
							Op:    "+",
							Left:  leftOperand,
							Right: &pb.Operand{Value: &pb.Operand_Number{Number: 1}},
						}
					}
					return instrs
				}(),
			},
			want:    nil,
			wantErr: true,
			setup: func(cancel context.CancelFunc, wg *sync.WaitGroup) {
				go func() {
					fmt.Printf("DEBUG: Cancelling context\n")
					cancel()
					wg.Done()
				}()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := NewBuilder()
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
			defer cancel()

			var wg sync.WaitGroup
			if tt.setup != nil {
				wg.Add(1)
				tt.setup(cancel, &wg)
			}

			got, err := cb.ProcessInstructions(ctx, tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ProcessInstructions() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.want != nil && !proto.Equal(got, tt.want) {
				t.Errorf("ProcessInstructions() got = %+v, want = %+v", got, tt.want)
			}
		})
	}
}

func TestBuilder_ResolveOperand(t *testing.T) {
	cb := NewBuilder()
	cb.variables["x"] = 42

	tests := []struct {
		name    string
		operand *pb.Operand
		want    int64
		wantErr bool
	}{
		{
			name:    "Nil_operand",
			operand: nil,
			want:    0,
			wantErr: true,
		},
		{
			name:    "Number_operand",
			operand: &pb.Operand{Value: &pb.Operand_Number{Number: 10}},
			want:    10,
			wantErr: false,
		},
		{
			name:    "Existing_variable",
			operand: &pb.Operand{Value: &pb.Operand_Variable{Variable: "x"}},
			want:    42,
			wantErr: false,
		},
		{
			name:    "Non-existing_variable",
			operand: &pb.Operand{Value: &pb.Operand_Variable{Variable: "y"}},
			want:    0,
			wantErr: true,
		},
		{
			name:    "Invalid_operand_type",
			operand: &pb.Operand{},
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var mu sync.Mutex
			cond := sync.NewCond(&mu)
			got, err := cb.resolveOperand(tt.operand, &mu, cond)
			if (err != nil) != tt.wantErr {
				t.Errorf("resolveOperand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("resolveOperand() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewBuilder(t *testing.T) {
	cb := NewBuilder()
	if cb.variables == nil {
		t.Error("NewBuilder() variables map is nil")
	}
	if cb.results == nil {
		t.Error("NewBuilder() results slice is nil")
	}
	if len(cb.variables) != 0 {
		t.Errorf("NewBuilder() variables map is not empty, got %v", cb.variables)
	}
	if len(cb.results) != 0 {
		t.Errorf("NewBuilder() results slice is not empty, got %v", cb.results)
	}
}
