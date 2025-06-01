package calculator

// hasCycle находит циклические зависимости
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
