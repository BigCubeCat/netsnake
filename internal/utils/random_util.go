package utils

import "math/rand/v2"

const MIN_ID = 2        // чтобы не путать с пустотой и едой
const MAX_ID = 16581375 // чтобы можно было однозначно перевести ID в цвет

func RandRange(min, max int) int {
	return rand.IntN(max-min) + min
}

func RandomId() int {
	return RandRange(MIN_ID, MAX_ID)
}
