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

// ProcessInstructions обрабатывает список инструкций и возвращает результаты print-инструкций
func (cb *CalculatorBuilder) ProcessInstructions(instructions []Instruction) ([]ResultItem, error) {
	return nil, nil
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
