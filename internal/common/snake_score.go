package common

import "fmt"

func (s *SnakeScore) String() string {
	result := fmt.Sprintf("user[%d]: %d points", s.UserId, s.Score)
	if !s.Alive {
		result += " died"
	}
	return result + "\n"
}
