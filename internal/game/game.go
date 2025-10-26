package game

import (
	"fmt"
	"github.com/eiannone/keyboard"
	"math/rand"
	"os"
	"os/signal"
	"snakeoop/internal/cfg"
	"snakeoop/internal/exec"
	"snakeoop/internal/snake"
	"syscall"
	"time"
)

type Game struct {
	apple      cfg.Point
	freeFields []cfg.Point
	snake      *snake.Snake
	score      int
	stop       chan bool
	click      chan keyboard.Key
	t          *time.Ticker
}

func Start() {
	g := Game{
		apple:      cfg.Point{X: rand.Intn(cfg.SizeX-2) + 1, Y: rand.Intn(cfg.SizeY-2) + 1},
		freeFields: nil,
		snake:      snake.New(),
		score:      0,
		stop:       make(chan bool, 1),
		click:      make(chan keyboard.Key, 1),
		t:          time.NewTicker(500 * time.Millisecond),
	}

	defer g.t.Stop()

	go func() {
		if err := keyboard.Open(); err != nil {
			panic(err)
		}
		defer func() {
			_ = keyboard.Close()
		}()

		for {
			_, key, err := keyboard.GetSingleKey()
			if err != nil {
				panic(err)
			}
			g.click <- key
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	for {
		select {
		case <-g.t.C:
			g.render()
		case key := <-g.click:
			switch key {
			case keyboard.KeyArrowUp:
				g.snake.ChangeDir(cfg.DirUP)
			case keyboard.KeyArrowDown:
				g.snake.ChangeDir(cfg.DirDown)
			case keyboard.KeyArrowLeft:
				g.snake.ChangeDir(cfg.DirLeft)
			case keyboard.KeyArrowRight:
				g.snake.ChangeDir(cfg.DirRight)
			}
			g.render()
		case <-sig:
			End("int")
		case <-g.stop:
			break
		}
	}

}

func (g *Game) toRestart() {
	g.stop <- true
	exec.Clean()
	fmt.Printf("You dead! \nScore: %d\n", g.score)
	time.Sleep(time.Second)
	fmt.Println("To restart press enter")

	for c := range g.click {
		if c == keyboard.KeyEnter {
			break
		}
	}
	exec.Clean()
	Start()
}

func (g *Game) render() {
	g.snake.Move()

	if g.snake.IsDead() {
		g.toRestart()
	}
	g.checkEaten()
	g.draw()
}

func End(rec string) {
	switch rec {
	case "recover":
		fmt.Println("ERROR!!")
		os.Exit(1)
	case "int":
		fmt.Println("goodbye")
		os.Exit(0)
	default:
		fmt.Println("Undefined error!!")
		os.Exit(1)
	}
}

func (g *Game) newApple() {
	g.apple = g.freeFields[rand.Intn(len(g.freeFields))]
}

func (g *Game) checkEaten() {
	if _, at := g.snake.IsSnakeAt(g.apple.X, g.apple.Y); at {
		g.score++
		g.snake.Create()
		g.newApple()
	}
}

func (g *Game) draw() {
	g.freeFields = []cfg.Point{}

	exec.Clean()
	var field [cfg.SizeY][cfg.SizeX]rune
	//field := make([][]rune, cfg.SizeY, cfg.SizeX)

	fmt.Printf("Score: %d\n", g.score)

	for y := range field {
		for x := range field[y] {
			switch {
			case y == 0 || y == cfg.SizeY-1 || x == 0 || x == cfg.SizeX-1:
				field[y][x] = '#'
			case y == g.apple.Y && x == g.apple.X:
				field[y][x] = '%'
			default:
				if dir, at := g.snake.IsSnakeAt(x, y); at {
					field[y][x] = getSnakeSymbol(dir)
					continue
				}
				field[y][x] = ' '
				g.freeFields = append(g.freeFields, cfg.Point{X: x, Y: y})
			}
		}
	}

	for _, row := range field {
		for _, s := range row {
			fmt.Print(string(s))
		}
		fmt.Print("\n")
	}
}

func getSnakeSymbol(dir string) rune {
	switch dir {
	case cfg.DirUP, cfg.DirDown:
		return '|'
	case cfg.DirLeft, cfg.DirRight:
		return '-'
	default:
		return ' '
	}
}
