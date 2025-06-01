echo "Test 1: Simple addition"
grpcurl -plaintext -d '{"instructions": [{"type":"calc","op":"+","var":"x","left":{"number":1},"right":{"number":2}},{"type":"print","var":"x"}]}' localhost:50051 calculator.v1.CalculatorService/ProcessInstructions
echo -e "\n"

echo "Test 2: Variable dependencies"
grpcurl -plaintext -d '{"instructions": [{"type":"calc","op":"+","var":"x","left":{"number":10},"right":{"number":2}},{"type":"calc","op":"-","var":"y","left":{"variable":"x"},"right":{"number":3}},{"type":"print","var":"y"}]}' localhost:50051 calculator.v1.CalculatorService/ProcessInstructions
echo -e "\n"

echo "Test 3: Immutable variable error"
grpcurl -plaintext -d '{"instructions": [{"type":"calc","op":"+","var":"x","left":{"number":1},"right":{"number":2}},{"type":"calc","op":"+","var":"x","left":{"number":3},"right":{"number":4}}]}' localhost:50051 calculator.v1.CalculatorService/ProcessInstructions
echo -e "\n"

echo "Test 4: Parallel execution"
grpcurl -plaintext -d '{"instructions": [{"type":"calc","op":"+","var":"x","left":{"number":1},"right":{"number":2}},{"type":"calc","op":"*","var":"y","left":{"number":3},"right":{"number":4}},{"type":"print","var":"x"},{"type":"print","var":"y"}]}' localhost:50051 calculator.v1.CalculatorService/ProcessInstructions
echo -e "\n"

echo "Test 5: Print order"
grpcurl -plaintext -d '{"instructions": [{"type":"calc","op":"+","var":"x","left":{"number":10},"right":{"number":2}},{"type":"calc","op":"*","var":"y","left":{"variable":"x"},"right":{"number":5}},{"type":"calc","op":"-","var":"q","left":{"variable":"y"},"right":{"number":20}},{"type":"print","var":"q"},{"type":"calc","op":"-","var":"z","left":{"variable":"x"},"right":{"number":15}},{"type":"print","var":"z"},{"type":"print","var":"x"}]}' localhost:50051 calculator.v1.CalculatorService/ProcessInstructions
echo -e "\n"

echo "Test 6: Invalid operation"
grpcurl -plaintext -d '{"instructions": [{"type":"calc","op":"/","var":"x","left":{"number":1},"right":{"number":2}}]}' localhost:50051 calculator.v1.CalculatorService/ProcessInstructions
echo -e "\n"

echo "Test 7: Unknown variable"
grpcurl -plaintext -d '{"instructions": [{"type":"print","var":"x"}]}' localhost:50051 calculator.v1.CalculatorService/ProcessInstructions
echo -e "\n"

echo "Test 8: Cyclic dependency"
grpcurl -plaintext -d '{"instructions": [{"type":"calc","op":"+","var":"x","left":{"variable":"y"},"right":{"number":1}},{"type":"calc","op":"-","var":"y","left":{"variable":"x"},"right":{"number":2}}]}' localhost:50051 calculator.v1.CalculatorService/ProcessInstructions
echo -e "\n"

echo "Test 9: Empty input"
grpcurl -plaintext -d '{"instructions": []}' localhost:50051 calculator.v1.CalculatorService/ProcessInstructions
echo -e "\n"