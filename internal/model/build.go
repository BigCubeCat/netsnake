package model

func (g *Game) Rebuild(width, height, food int) {
	g.State.Width = width
	g.State.Height = height
	g.FoodStatic = food
}

func (g *Game) ResetSnakes() {
	g.State.Snakes = []*Snake{}
}

func (g *Game) SetSnake(s *Snake) {
	g.State.Snakes = append(g.State.Snakes, s)
}
