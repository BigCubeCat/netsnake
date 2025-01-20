package model

import (
	"math/rand"
)

type Role int

const (
	Viewer Role = iota
	Normal
	Master
	Deputy
)

type SnakeState int

const (
	ALIVE SnakeState = iota
	ZOMBIE
)

// Direction представляет направление движения змейки.
type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

// Point представляет координаты на игровом поле.
type Point struct {
	X int
	Y int
}

// Snake представляет змейку.
type Snake struct {
	ID        int
	Role      Role
	Body      []Point
	Direction Direction
	Score     int
	IsAlive   bool
}

// GameState представляет состояние игры.
type GameState struct {
	Width    int
	Height   int
	Snakes   []*Snake
	Food     []Point
	StateID  int
	GameOver bool
}

// Game представляет игровую логику.
type Game struct {
	State      *GameState
	FoodStatic int
	Rand       *rand.Rand
}
