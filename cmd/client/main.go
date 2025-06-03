package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/goCalcProj/gen/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

// clientInstruction — упрощённое представление JSON-формата
type clientInstruction struct {
	Type  string      `json:"type"`
	Op    string      `json:"op,omitempty"`
	Var   string      `json:"var"`
	Left  interface{} `json:"left,omitempty"`
	Right interface{} `json:"right,omitempty"`
}

// clientInstruction в pb.Instruction
func (ci *clientInstruction) toProto() (*pb.Instruction, error) {
	instr := &pb.Instruction{
		Type: ci.Type,
		Op:   ci.Op,
		Var:  ci.Var,
	}

	if ci.Left != nil {
		operand, err := toOperand(ci.Left)
		if err != nil {
			return nil, fmt.Errorf("некорректный левый операнд: %v", err)
		}
		instr.Left = operand
	}

	if ci.Right != nil {
		operand, err := toOperand(ci.Right)
		if err != nil {
			return nil, fmt.Errorf("некорректный правый операнд: %v", err)
		}
		instr.Right = operand
	}

	return instr, nil
}

// интерфейсное значение в pb.Operand
func toOperand(v interface{}) (*pb.Operand, error) {
	switch val := v.(type) {
	case float64:
		if val != float64(int64(val)) {
			return nil, fmt.Errorf("операнд должен быть целым числом")
		}
		return &pb.Operand{
			Value: &pb.Operand_Number{Number: int64(val)},
		}, nil
	case string:
		return &pb.Operand{
			Value: &pb.Operand_Variable{Variable: val},
		}, nil
	case map[string]interface{}:
		// Новый блок: пришёл JSON-объект, например {"number": 1} или {"variable": "x"}
		if numRaw, ok := val["number"]; ok {
			// проверяем, что тип у numRaw — число
			if numF, ok2 := numRaw.(float64); ok2 {
				if numF != float64(int64(numF)) {
					return nil, fmt.Errorf("операнд.number должен быть целым числом")
				}
				return &pb.Operand{
					Value: &pb.Operand_Number{Number: int64(numF)},
				}, nil
			}
			return nil, fmt.Errorf("значение number должно быть числом")
		}
		if varRaw, ok := val["variable"]; ok {
			if varStr, ok2 := varRaw.(string); ok2 {
				return &pb.Operand{
					Value: &pb.Operand_Variable{Variable: varStr},
				}, nil
			}
			return nil, fmt.Errorf("значение variable должно быть строкой")
		}
		return nil, fmt.Errorf("в объекте ожидался ключ \"number\" или \"variable\"")
	default:
		return nil, fmt.Errorf("неподдерживаемый тип операнда: %T", v)
	}
}

func main() {
	// Флаг -debug
	debug := flag.Bool("debug", false, "включить вывод отладочной информации")
	flag.Parse()

	// Подключаемся к gRPC-серверу
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Ошибка подключения: %v", err)
	}
	defer conn.Close()

	client := pb.NewCalculatorServiceClient(conn)
	reader := bufio.NewReader(os.Stdin)

	jsonOpts := protojson.MarshalOptions{
		Indent:          "  ",
		EmitUnpopulated: true,
	}

	for {
		fmt.Println("Введите JSON-запрос (двойной Enter — отправка, 'exit' — выход):")
		var input strings.Builder
		emptyLine := false

		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				log.Fatalf("Ошибка чтения ввода: %v", err)
			}
			line = strings.TrimSpace(line)

			if line == "exit" {
				fmt.Println("Выход...")
				return
			}

			if line == "" {
				if emptyLine || input.Len() == 0 {
					break
				}
				emptyLine = true
				continue
			}

			emptyLine = false
			input.WriteString(line + "\n")
		}

		if strings.TrimSpace(input.String()) == "" {
			continue
		}

		var clientInstrs []clientInstruction
		if err := json.Unmarshal([]byte(input.String()), &clientInstrs); err != nil {
			fmt.Printf("Ошибка разбора JSON: %v\n", err)
			continue
		}

		req := &pb.ProcessInstructionsRequest{}
		for _, ci := range clientInstrs {
			protoInstr, err := ci.toProto()
			if err != nil {
				fmt.Printf("Ошибка в инструкции: %v\n", err)
				continue
			}
			req.Instructions = append(req.Instructions, protoInstr)
		}

		resp, err := client.ProcessInstructions(context.Background(), req)
		if err != nil {
			fmt.Printf("Ошибка выполнения запроса: %v\n", err)
			continue
		}

		// -debug
		if *debug {
			fmt.Println("Отладка: полученные значения переменных:")
			for _, item := range resp.Items {
				fmt.Printf("Переменная: %s, значение: %d\n", item.Var, item.Value)
			}
		}

		respJSON, err := jsonOpts.Marshal(resp)
		if err != nil {
			fmt.Printf("Ошибка сериализации ответа: %v\n", err)
			continue
		}

		fmt.Println("Ответ:")
		fmt.Println(string(respJSON))
	}
}
