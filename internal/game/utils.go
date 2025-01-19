package game

import "math/rand/v2"

func RandRange(min, max int) int {
	return rand.IntN(max-min) + min
}

func randomWithPercent(successProb uint) bool {
	return RandRange(0, 100) <= int(successProb)
}

func findSquares(matrix [][]ItemType, k int) []GameCoords {
	n := len(matrix)    // Число строк
	m := len(matrix[0]) // Число столбцов

	if k > n || k > m {
		return nil // Квадрат такого размера не может существовать
	}

	// Префиксная сумма для быстрой проверки суммы элементов в квадратах
	prefixSum := make([][]int, n+1)
	for i := range prefixSum {
		prefixSum[i] = make([]int, m+1)
	}

	// Заполняем префиксную сумму
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			prefixSum[i+1][j+1] = int(
				matrix[i][j],
			) + prefixSum[i+1][j] + prefixSum[i][j+1] - prefixSum[i][j]
		}
	}

	// Результат: список координат верхнего левого угла всех подходящих квадратов
	var result []GameCoords
	// Проверяем все квадраты k x k
	for i := 0; i <= n-k; i++ {
		for j := 0; j <= m-k; j++ {
			// Сумма элементов в квадрате (i, j) -> (i+k-1, j+k-1)
			areaSum := prefixSum[i+k][j+k] - prefixSum[i+k][j] - prefixSum[i][j+k] + prefixSum[i][j]

			// Если сумма элементов <= k*k*2, добавляем координаты в результат
			if areaSum == 0 {
				result = append(result, GameCoords{i, j})
			}
		}
	}

	return result
}
