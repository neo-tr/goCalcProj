package main

import (
	"encoding/json"
	"fmt"
	"sync"
)

// Instruction определяет операцию: вычисление (calc) или вывод (print)
type Instruction struct {
	Type  string      `json:"type"`
	Op    string      `json:"op,omitempty"`
	Var   string      `json:"var"`
	Left  interface{} `json:"left,omitempty"`
	Right interface{} `json:"right,omitempty"`
}

// ResultItem представляет результат для инструкции типа print
type ResultItem struct {
	Var   string `json:"var"`
	Value int64  `json:"value"`
}

// CalculatorBuilder выполняет инструкции и сохраняет результаты
type CalculatorBuilder struct {
	variables map[string]int64
	results   []ResultItem
	mu        sync.Mutex
}

// NewCalculatorBuilder создает и инициализирует новый CalculatorBuilder
func NewCalculatorBuilder() *CalculatorBuilder {
	return &CalculatorBuilder{
		variables: make(map[string]int64),
		results:   []ResultItem{},
	}
}

// resolveOperand превращает операнд в значение int64
func (cb *CalculatorBuilder) resolveOperand(operand interface{}) (int64, error) {
	switch v := operand.(type) {
	case float64:
		return int64(v), nil
	case string:
		cb.mu.Lock()
		defer cb.mu.Unlock()
		if val, exists := cb.variables[v]; exists {
			return val, nil
		}
		return 0, fmt.Errorf("variable %s not found", v)
	default:
		return 0, fmt.Errorf("invalid operand type: %T", v)
	}
}

// ProcessInstructions обрабатывает список инструкций и возвращает результаты print-инструкций
func (cb *CalculatorBuilder) ProcessInstructions(instructions []Instruction) ([]ResultItem, error) {
	for _, instr := range instructions {
		if instr.Type == "calc" {
			leftVal, err := cb.resolveOperand(instr.Left)
			if err != nil {
				return nil, fmt.Errorf("resolving left operand for %s: %w", instr.Var, err)
			}
			rightVal, err := cb.resolveOperand(instr.Right)
			if err != nil {
				return nil, fmt.Errorf("resolving right operand for %s: %w", instr.Var, err)
			}

			var result int64
			switch instr.Op {
			case "+":
				result = leftVal + rightVal
			case "-":
				result = leftVal - rightVal
			case "*":
				result = leftVal * rightVal
			default:
				return nil, fmt.Errorf("unsupported operation for %s: %s", instr.Var, instr.Op)
			}

			cb.mu.Lock()
			if _, exists := cb.variables[instr.Var]; exists {
				cb.mu.Unlock()
				return nil, fmt.Errorf("variable %s is immutable and cannot be reassigned", instr.Var)
			}
			cb.variables[instr.Var] = result
			cb.mu.Unlock()
		} else if instr.Type == "print" {
			cb.mu.Lock()
			value, exists := cb.variables[instr.Var]
			if !exists {
				cb.mu.Unlock()
				return nil, fmt.Errorf("variable %s not found", instr.Var)
			}
			cb.results = append(cb.results, ResultItem{Var: instr.Var, Value: value})
			cb.mu.Unlock()
		}
	}
	return cb.results, nil
}

func main() {
	// Пример входных данных в формате JSON
	inputJSON := `[
    	{ "type": "calc", "op": "+", "var": "x", "left": 1,  "right": 2 },
  		{ "type": "print", "var": "x" }
	]`

	// Десериализация JSON в срез инструкций
	var instructions []Instruction
	if err := json.Unmarshal([]byte(inputJSON), &instructions); err != nil {
		fmt.Printf("Error unmarshaling JSON: %v\n", err)
		return
	}

	// Создание объекта CalculatorBuilder
	builder := NewCalculatorBuilder()

	// Выполнение инструкций
	results, err := builder.ProcessInstructions(instructions)
	if err != nil {
		fmt.Printf("Error processing instructions: %v\n", err)
		return
	}

	// Сериализация и вывод результатов выполнения инструкций типа print
	output, _ := json.MarshalIndent(map[string][]ResultItem{"items": results}, "", "  ")
	fmt.Println(string(output))
}
