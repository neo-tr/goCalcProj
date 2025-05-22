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

// NewCalculatorBuilder создает и инициализирует новый CalculatorBuilderJ
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

// hasCycle ищет циклические зависимости поиском в глубину dfs
func hasCycle(deps map[string][]string) bool {
	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	var dfs func(node string) bool
	dfs = func(node string) bool {
		visited[node] = true
		recStack[node] = true

		for _, neighbor := range deps[node] {
			if !visited[neighbor] {
				if dfs(neighbor) {
					return true
				}
			} else if recStack[neighbor] {
				return true
			}
		}

		recStack[node] = false
		return false
	}

	for node := range deps {
		if !visited[node] {
			if dfs(node) {
				return true
			}
		}
	}
	return false
}

/*
ProcessInstructions выполняет список инструкций конкурентно, обрабатывая операции calc и print.
Строит граф зависимостей, проверяет наличие цикличности и обрабатывает инструкции при помощи горутин.
Результаты собираются из print, а ошибки возвращаются если будут обнаружены.
*/
func (cb *CalculatorBuilder) ProcessInstructions(instructions []Instruction) ([]ResultItem, error) {
	//построение мапы зависимостей и определить, какие будут вычислены
	deps := make(map[string][]string)
	var printOrder []string
	computedVars := make(map[string]bool)
	for _, instr := range instructions {
		if instr.Type == "calc" {
			if leftStr, ok := instr.Left.(string); ok && leftStr != "" {
				deps[instr.Var] = append(deps[instr.Var], leftStr)
			}
			if rightStr, ok := instr.Right.(string); ok && rightStr != "" {
				deps[instr.Var] = append(deps[instr.Var], rightStr)
			}
			computedVars[instr.Var] = true
		} else if instr.Type == "print" {
			printOrder = append(printOrder, instr.Var)
		}
	}

	if hasCycle(deps) {
		return nil, fmt.Errorf("cyclic dependency detected")
	}

	varMu := sync.Mutex{}
	cond := sync.NewCond(&varMu)
	computed := make(map[string]bool)
	var wg sync.WaitGroup
	errChan := make(chan error, len(instructions))
	resultsChan := make(chan ResultItem, len(instructions))

	// Запуск горутины для каждой инструкции чтобы выполнить calc или print
	for _, instr := range instructions {
		wg.Add(1)
		go func(instr Instruction) {
			defer wg.Done()
			if instr.Type == "calc" {
				//Определить зависимые переменные
				var leftVar, rightVar string
				if leftStr, ok := instr.Left.(string); ok {
					leftVar = leftStr
				}
				if rightStr, ok := instr.Right.(string); ok {
					rightVar = rightStr
				}

				// Ждём, пока зависимые переменные не будут вычислены
				if leftVar != "" {
					varMu.Lock()
					for !computed[leftVar] {
						cond.Wait()
					}
					varMu.Unlock()
				}
				if rightVar != "" {
					varMu.Lock()
					for !computed[rightVar] {
						cond.Wait()
					}
					varMu.Unlock()
				}

				// Получаем значения операндов
				leftVal, err := cb.resolveOperand(instr.Left)
				if err != nil {
					errChan <- fmt.Errorf("resolving left operand for %s: %w", instr.Var, err)
					return
				}
				rightVal, err := cb.resolveOperand(instr.Right)
				if err != nil {
					errChan <- fmt.Errorf("resolving right operand for %s: %w", instr.Var, err)
					return
				}

				// Выполнение вычисления в зависимости от оператора
				var result int64
				switch instr.Op {
				case "+":
					result = leftVal + rightVal
				case "-":
					result = leftVal - rightVal
				case "*":
					result = leftVal * rightVal
				default:
					errChan <- fmt.Errorf("unsupported operation for %s: %s", instr.Var, instr.Op)
					return
				}

				// Сохранить результат с проверкой на иммутабельность переменной
				cb.mu.Lock()
				if _, exists := cb.variables[instr.Var]; exists {
					errChan <- fmt.Errorf("variable %s is immutable and cannot be reassigned", instr.Var)
					cb.mu.Unlock()
					return
				}
				cb.variables[instr.Var] = result
				cb.mu.Unlock()

				// Отметить переменную как вычисленную - true и оповестить другие горутины через Broadcast
				varMu.Lock()
				computed[instr.Var] = true
				cond.Broadcast()
				varMu.Unlock()
			} else if instr.Type == "print" {
				varName := instr.Var
				// Проверить, будет ли переменная вычислена
				if !computedVars[varName] {
					errChan <- fmt.Errorf("variable %s not found", varName)
					return
				}
				// Ждём, пока переменная будет вычислена
				varMu.Lock()
				for !computed[varName] {
					cond.Wait()
				}
				varMu.Unlock()

				// Получаем значение переменной и отправляем его в канал - resultsChan или хэндлим ошибку в errChan
				cb.mu.Lock()
				value, exists := cb.variables[varName]
				if !exists {
					errChan <- fmt.Errorf("variable %s not found", varName)
					cb.mu.Unlock()
					return
				}
				resultsChan <- ResultItem{Var: varName, Value: value}
				cb.mu.Unlock()
			}
		}(instr)
	}

	// Ждём пока всё выполнится и никто не будет слать, закрываем каналы
	go func() {
		wg.Wait()
		close(resultsChan)
		close(errChan)
	}()

	// Проверка на наличие ошибок
	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}

	// Результат в порядке, указанном в print
	resultMap := make(map[string]ResultItem)
	for res := range resultsChan {
		resultMap[res.Var] = res
	}

	cb.results = nil
	for _, varName := range printOrder {
		if res, exists := resultMap[varName]; exists {
			cb.results = append(cb.results, res)
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
