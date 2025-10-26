package snake

import (
	"snakeoop/internal/cfg"
)

type Snake struct {
	head head
	body []body
}

type head struct {
	Dir string
	P   cfg.Point
}

type body struct {
	Dir string
	P   cfg.Point
}

func New() *Snake {
	return &Snake{
		head: head{
			Dir: cfg.DirUP,
			P:   cfg.Point{X: cfg.SizeX / 2, Y: cfg.SizeY / 2},
		},
		body: []body{
			{
				P:   cfg.Point{X: cfg.SizeX / 2, Y: cfg.SizeY/2 + 1},
				Dir: cfg.DirUP,
			},
		},
	}
}

func (s *Snake) IsSnakeAt(x, y int) (string, bool) {
	if s.head.P.X == x && s.head.P.Y == y {
		return s.head.Dir, true
	}

	for _, b := range s.body {
		if b.P.X == x && b.P.Y == y {
			return b.Dir, true
		}
	}

	return "", false
}

func (s *Snake) IsDead() bool {
	for i := 2; i < len(s.body); i++ {
		if s.body[i].P.X == s.head.P.X &&
			s.body[i].P.Y == s.head.P.Y {
			return true
		}
	}
	return false
}

func (s *Snake) Create() {
	s.body = append(s.body, body{
		Dir: s.body[len(s.body)-1].Dir,
		P:   s.body[len(s.body)-1].P,
	})
}

func (s *Snake) canChangeDirection(newDir string) bool {
	opposites := map[string]string{
		cfg.DirUP: cfg.DirDown, cfg.DirDown: cfg.DirUP,
		cfg.DirLeft: cfg.DirRight, cfg.DirRight: cfg.DirLeft,
	}
	return opposites[s.head.Dir] != newDir
}
func (s *Snake) ChangeDir(newDir string) {
	if !s.canChangeDirection(newDir) {
		return
	}
	s.head.Dir = newDir
}

func (s *Snake) Move() {
	newBody := make([]body, len(s.body))

	for i := 1; i <= len(s.body)-1; i++ {
		newBody[i].P.X = s.body[i-1].P.X
		newBody[i].P.Y = s.body[i-1].P.Y
		newBody[i].Dir = s.body[i-1].Dir
	}
	newBody[0].P.X = s.head.P.X
	newBody[0].P.Y = s.head.P.Y
	newBody[0].Dir = s.head.Dir

	s.body = newBody
	switch s.head.Dir {
	case cfg.DirDown:
		s.head.P.Y++
		if s.head.P.Y >= cfg.SizeY {
			s.head.P.Y = 0
		}
	case cfg.DirUP:
		s.head.P.Y--
		if s.head.P.Y < 0 {
			s.head.P.Y = cfg.SizeY - 1
		}
	case cfg.DirLeft:
		s.head.P.X--
		if s.head.P.X < 0 {
			s.head.P.X = cfg.SizeX - 1
		}
	case cfg.DirRight:
		s.head.P.X++
		if s.head.P.X >= cfg.SizeX {
			s.head.P.X = 0
		}
	}
}
