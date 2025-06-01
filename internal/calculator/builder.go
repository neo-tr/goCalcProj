package calculator

import (
	"context"
	"fmt"
	"sync"

	pb "github.com/goCalcProj/gen/pb"
)

var debugEnabled = true

func debugLog(format string, args ...interface{}) {
	if debugEnabled {
		fmt.Printf("DEBUG: "+format+"\n", args...)
	}
}

// ProcessInstructions — выполняет список инструкций параллельно
func (cb *Builder) ProcessInstructions(ctx context.Context, req *pb.ProcessInstructionsRequest) (*pb.ProcessInstructionsResponse, error) {
	// Зависимости переменных и порядок вывода
	deps := make(map[string][]string)
	var printOrder []string
	computedVars := make(map[string]bool)

	// Первый проход — какие переменные будут вычислены
	for _, instr := range req.Instructions {
		if instr.Type == "calc" {
			computedVars[instr.Var] = true
		}
	}

	// Второй проход — зависимости и команды print
	for _, instr := range req.Instructions {
		if instr.Type == "calc" {
			if leftVar := instr.GetLeft().GetVariable(); leftVar != "" {
				deps[instr.Var] = append(deps[instr.Var], leftVar)
			}
			if rightVar := instr.GetRight().GetVariable(); rightVar != "" {
				deps[instr.Var] = append(deps[instr.Var], rightVar)
			}
		} else if instr.Type == "print" {
			if !computedVars[instr.Var] {
				debugLog("Неизвестная переменная %s для печати", instr.Var)
				return nil, fmt.Errorf("неизвестная переменная: %s", instr.Var)
			}
			printOrder = append(printOrder, instr.Var)
		}
	}

	// Проверка на циклические зависимости
	if hasCycle(deps) {
		debugLog("Обнаружена циклическая зависимость")
		return nil, fmt.Errorf("обнаружена циклическая зависимость")
	}

	// Синхронизация
	varMu := sync.Mutex{}
	cond := sync.NewCond(&varMu)
	computed := make(map[string]bool)
	var wg sync.WaitGroup
	errChan := make(chan error, len(req.Instructions))
	resultsChan := make(chan *pb.ResultItem, len(req.Instructions))

	// Каждая инструкция — в отдельной горутине
	for _, instr := range req.Instructions {
		wg.Add(1)
		go func(instr *pb.Instruction) {
			defer wg.Done()

			switch instr.Type {
			case "calc":
				if instr.Left == nil || instr.Right == nil {
					debugLog("Пустой операнд у %s", instr.Var)
					varMu.Lock()
					cond.Broadcast()
					varMu.Unlock()
					errChan <- fmt.Errorf("пустой операнд у переменной %s", instr.Var)
					return
				}

				// Ожидание зависимых переменных
				for _, varName := range []string{
					instr.GetLeft().GetVariable(),
					instr.GetRight().GetVariable(),
				} {
					if varName == "" {
						continue
					}
					varMu.Lock()
					for !computed[varName] {
						select {
						case <-ctx.Done():
							debugLog("Контекст отменён при ожидании %s", varName)
							cond.Broadcast()
							varMu.Unlock()
							return
						default:
							cond.Wait()
						}
					}
					varMu.Unlock()
				}

				// Получаем значения
				leftVal, err := cb.resolveOperand(instr.Left, &varMu, cond)
				if err != nil {
					errChan <- fmt.Errorf("ошибка операнда слева (%s): %w", instr.Var, err)
					return
				}
				rightVal, err := cb.resolveOperand(instr.Right, &varMu, cond)
				if err != nil {
					errChan <- fmt.Errorf("ошибка операнда справа (%s): %w", instr.Var, err)
					return
				}

				// Вычисление
				var result int64
				switch instr.Op {
				case "+":
					result = leftVal + rightVal
				case "-":
					result = leftVal - rightVal
				case "*":
					result = leftVal * rightVal
				default:
					errChan <- fmt.Errorf("неподдерживаемая операция %s: %s", instr.Var, instr.Op)
					return
				}

				// Сохраняем результат
				cb.mu.Lock()
				if _, exists := cb.variables[instr.Var]; exists {
					errChan <- fmt.Errorf("переменная %s уже вычислена", instr.Var)
					cb.mu.Unlock()
					return
				}
				cb.variables[instr.Var] = result
				cb.mu.Unlock()

				varMu.Lock()
				computed[instr.Var] = true
				cond.Broadcast()
				varMu.Unlock()

			case "print":
				varName := instr.Var
				varMu.Lock()
				for !computed[varName] {
					select {
					case <-ctx.Done():
						varMu.Unlock()
						return
					default:
						cond.Wait()
					}
				}
				varMu.Unlock()

				cb.mu.Lock()
				value, ok := cb.variables[varName]
				cb.mu.Unlock()
				if !ok {
					errChan <- fmt.Errorf("переменная %s не найдена", varName)
					return
				}

				resultsChan <- &pb.ResultItem{Var: varName, Value: value}
			}
		}(instr)
	}

	// Горутина завершения
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(resultsChan)
		close(errChan)
		close(done)
	}()

	// Обработка ошибок
	for {
		select {
		case err := <-errChan:
			if err != nil {
				<-done
				return nil, err
			}
		case <-ctx.Done():
			<-done
			return nil, ctx.Err()
		case <-done:
			goto collectResults
		}
	}

collectResults:
	// Собираем результат в порядке печати
	resultMap := map[string]*pb.ResultItem{}
	for item := range resultsChan {
		resultMap[item.Var] = item
	}

	resp := &pb.ProcessInstructionsResponse{}
	for _, varName := range printOrder {
		if item, ok := resultMap[varName]; ok {
			resp.Items = append(resp.Items, item)
		}
	}

	return resp, nil
}
