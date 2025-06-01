package calculator

import (
	"fmt"
	"sync"

	pb "github.com/goCalcProj/gen/pb"
)

// resolveOperand — преобразует операнд в значение int64, извлекая число или переменную
func (cb *Builder) resolveOperand(operand *pb.Operand, varMu *sync.Mutex, cond *sync.Cond) (int64, error) {
	if operand == nil {
		return 0, fmt.Errorf("нулевой операнд")
	}

	switch v := operand.Value.(type) {
	case *pb.Operand_Number:
		return v.Number, nil

	case *pb.Operand_Variable:
		cb.mu.Lock()
		defer cb.mu.Unlock()

		if val, ok := cb.variables[v.Variable]; ok {
			return val, nil
		}

		varMu.Lock()
		defer varMu.Unlock()
		cond.Broadcast()
		return 0, fmt.Errorf("переменная %s не найдена", v.Variable)

	default:
		varMu.Lock()
		defer varMu.Unlock()
		cond.Broadcast()
		return 0, fmt.Errorf("неизвестный тип операнда: %v", v)
	}
}
