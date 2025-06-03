package calculator

import (
	"context"
	"fmt"
	"sync"

	pb "github.com/goCalcProj/gen/pb"
)

const debugEnabled = true

func debugLog(format string, args ...interface{}) {
	if debugEnabled {
		fmt.Printf("DEBUG: "+format+"\n", args...)
	}
}

// ProcessInstructions — выполняет список инструкций параллельно
func (cb *Builder) ProcessInstructions(ctx context.Context, req *pb.ProcessInstructionsRequest) (*pb.ProcessInstructionsResponse, error) {
	// Оборачиваем ctx в WithCancel, чтобы прерывать всех воркеров при первой ошибке
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	deps := make(map[string][]string)
	var printOrder []string
	computedVars := make(map[string]bool)

	// Первый проход: отмечаем, какие переменные вообще вычисляем
	for _, instr := range req.Instructions {
		if instr.Type == "calc" {
			computedVars[instr.Var] = true
		}
	}

	// Второй проход
	for _, instr := range req.Instructions {
		switch instr.Type {
		case "calc":
			// Если в левом или правом операнде есть переменная — добавляем зависимость
			if leftVar := instr.GetLeft().GetVariable(); leftVar != "" {
				deps[instr.Var] = append(deps[instr.Var], leftVar)
			}
			if rightVar := instr.GetRight().GetVariable(); rightVar != "" {
				deps[instr.Var] = append(deps[instr.Var], rightVar)
			}

		case "print":
			// Если пытаются печатать невычисляемую переменную — сразу ошибка
			if !computedVars[instr.Var] {
				debugLog("Неизвестная переменная %s для печати", instr.Var)
				return nil, fmt.Errorf("неизвестная переменная: %s", instr.Var)
			}
			printOrder = append(printOrder, instr.Var)
		}
	}

	// Проверяем циклические зависимости в графе deps
	if hasCycle(deps) {
		debugLog("Обнаружена циклическая зависимость")
		return nil, fmt.Errorf("обнаружена циклическая зависимость")
	}

	// Синхронизация между горутинами:
	varMu := sync.Mutex{}             // для cond
	cond := sync.NewCond(&varMu)      // чтобы уведомлять о появлении новых вычисленных переменных
	computed := make(map[string]bool) // какие переменные уже вычислены
	var wg sync.WaitGroup

	// Каналы для ошибок и результатов
	errChan := make(chan error, len(req.Instructions))
	resultsChan := make(chan *pb.ResultItem, len(req.Instructions))

	// Запускаем по одной горутине на каждую инструкцию
	for _, instr := range req.Instructions {
		wg.Add(1)
		go func(instr *pb.Instruction) {
			defer wg.Done()

			switch instr.Type {
			case "calc":
				// Проверяем, что обе стороны выражения заданы
				if instr.Left == nil || instr.Right == nil {
					debugLog("Пустой операнд у %s", instr.Var)
					// Прерываем всех
					cancel()
					varMu.Lock()
					cond.Broadcast()
					varMu.Unlock()
					errChan <- fmt.Errorf("пустой операнд у переменной %s", instr.Var)
					return
				}

				// Сразу проверяем, что все переменные в выражении вообще запланированы к вычислению.
				//    Если кто-то отсутствует в computedVars — возвращаем ошибку.
				for _, varName := range []string{
					instr.GetLeft().GetVariable(),
					instr.GetRight().GetVariable(),
				} {
					if varName == "" {
						continue
					}
					if !computedVars[varName] {
						// Прерываем всех и посылаем ошибку
						cancel()
						errChan <- fmt.Errorf("переменная %s не найдена", varName)
						return
					}
				}
				// Ждём, пока вычислятся все зависимости (computed[varName] станет true)
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
							// Если контекст отменили — выходим
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

				// Обе зависимости гарантированно есть в cb.variables — получаем значения
				leftVal, err := cb.resolveOperand(instr.Left, &varMu, cond)
				if err != nil {
					cancel()
					errChan <- fmt.Errorf("ошибка операнда слева (%s): %w", instr.Var, err)
					return
				}
				rightVal, err := cb.resolveOperand(instr.Right, &varMu, cond)
				if err != nil {
					cancel()
					errChan <- fmt.Errorf("ошибка операнда справа (%s): %w", instr.Var, err)
					return
				}

				// Выполняем операцию
				var result int64
				switch instr.Op {
				case "+":
					result = leftVal + rightVal
				case "-":
					result = leftVal - rightVal
				case "*":
					result = leftVal * rightVal
				default:
					cancel()
					errChan <- fmt.Errorf("неподдерживаемая операция %s: %s", instr.Var, instr.Op)
					return
				}

				// Записываем результат в карту переменных
				cb.mu.Lock()
				if _, exists := cb.variables[instr.Var]; exists {
					cb.mu.Unlock()
					cancel()
					errChan <- fmt.Errorf("переменная %s уже вычислена", instr.Var)
					return
				}
				cb.variables[instr.Var] = result
				cb.mu.Unlock()

				// Сообщаем остальным горутинам, что instr.Var вычислена
				varMu.Lock()
				computed[instr.Var] = true
				cond.Broadcast()
				varMu.Unlock()

			case "print":
				// Ждём, пока нужная переменная вычислится
				varName := instr.Var
				varMu.Lock()
				for !computed[varName] {
					select {
					case <-ctx.Done():
						// Если контекст отменили — выходим
						varMu.Unlock()
						return
					default:
						cond.Wait()
					}
				}
				varMu.Unlock()

				// Берём значение из карты
				cb.mu.Lock()
				value, ok := cb.variables[varName]
				cb.mu.Unlock()
				if !ok {
					cancel()
					errChan <- fmt.Errorf("переменная %s не найдена", varName)
					return
				}
				// Отправляем в результирующий канал
				resultsChan <- &pb.ResultItem{Var: varName, Value: value}
			}
		}(instr)
	}

	// Закрываем каналы после завершения всех горутин
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(resultsChan)
		close(errChan)
		close(done)
	}()

	// Обрабатываем первую появившуюся ошибку (или отмену контекста)
	for {
		select {
		case err := <-errChan:
			if err != nil {
				// Ждём, пока остальные горутины закончатся (они увидят ctx.Done и выйдут)
				<-done
				return nil, err
			}
		case <-ctx.Done():
			// Если контекст отменился снаружи — ждём завершения и возвращаем ctx.Err()
			<-done
			return nil, ctx.Err()
		case <-done:
			goto collectResults
		}
	}

collectResults:
	// Собираем результаты в порядке printOrder
	resultMap := make(map[string]*pb.ResultItem)
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
