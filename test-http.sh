echo "Test 1: Simple addition"
curl -X POST http://localhost:8080/api/v1/instructions -H "Content-Type: application/json" -d '{"instructions": [{"type":"calc","op":"+","var":"x","left":{"number":1},"right":{"number":2}},{"type":"print","var":"x"}]}'
echo -e "\n"

echo "Test 2: Variable dependencies"
curl -X POST http://localhost:8080/api/v1/instructions -H "Content-Type: application/json" -d '{"instructions": [{"type":"calc","op":"+","var":"x","left":{"number":10},"right":{"number":2}},{"type":"calc","op":"-","var":"y","left":{"variable":"x"},"right":{"number":3}},{"type":"print","var":"y"}]}'
echo -e "\n"

echo "Test 3: Immutable variable error"
curl -X POST http://localhost:8080/api/v1/instructions -H "Content-Type: application/json" -d '{"instructions": [{"type":"calc","op":"+","var":"x","left":{"number":1},"right":{"number":2}},{"type":"calc","op":"+","var":"x","left":{"number":3},"right":{"number":4}}]}'
echo -e "\n"

echo "Test 4: Parallel execution"
curl -X POST http://localhost:8080/api/v1/instructions -H "Content-Type: application/json" -d '{"instructions": [{"type":"calc","op":"+","var":"x","left":{"number":1},"right":{"number":2}},{"type":"calc","op":"*","var":"y","left":{"number":3},"right":{"number":4}},{"type":"print","var":"x"},{"type":"print","var":"y"}]}'
echo -e "\n"

echo "Test 5: Print order"
curl -X POST http://localhost:8080/api/v1/instructions -H "Content-Type: application/json" -d '{"instructions": [{"type":"calc","op":"+","var":"x","left":{"number":10},"right":{"number":2}},{"type":"calc","op":"*","var":"y","left":{"variable":"x"},"right":{"number":5}},{"type":"calc","op":"-","var":"q","left":{"variable":"y"},"right":{"number":20}},{"type":"print","var":"q"},{"type":"calc","op":"-","var":"z","left":{"variable":"x"},"right":{"number":15}},{"type":"print","var":"z"},{"type":"print","var":"x"}]}'
echo -e "\n"

echo "Test 6: Invalid operation"
curl -X POST http://localhost:8080/api/v1/instructions -H "Content-Type: application/json" -d '{"instructions": [{"type":"calc","op":"/","var":"x","left":{"number":1},"right":{"number":2}}]}'
echo -e "\n"

echo "Test 7: Unknown variable"
curl -X POST http://localhost:8080/api/v1/instructions -H "Content-Type: application/json" -d '{"instructions": [{"type":"print","var":"x"}]}'
echo -e "\n"

echo "Test 8: Cyclic dependency"
curl -X POST http://localhost:8080/api/v1/instructions -H "Content-Type: application/json" -d '{"instructions": [{"type":"calc","op":"+","var":"x","left":{"variable":"y"},"right":{"number":1}},{"type":"calc","op":"-","var":"y","left":{"variable":"x"},"right":{"number":2}}]}'
echo -e "\n"

echo "Test 9: Empty input"
curl -X POST http://localhost:8080/api/v1/instructions -H "Content-Type: application/json" -d '{"instructions": []}'
echo -e "\n"