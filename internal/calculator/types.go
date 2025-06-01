package calculator

import (
	"github.com/goCalcProj/gen/pb"
	"sync"
)

// Builder обрабатывает инструкции
type Builder struct {
	variables map[string]int64 // Хранит значения переменных
	results   []*pb.ResultItem // Хранит результаты вывода
	mu        sync.Mutex
}

func NewBuilder() *Builder {
	return &Builder{
		variables: make(map[string]int64),
		results:   []*pb.ResultItem{},
	}
}
