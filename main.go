package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"image/color"
	"math/rand/v2"
	"time"
)

var Up fyne.Position = fyne.Position{X: 0, Y: -step}
var Down fyne.Position = fyne.Position{X: 0, Y: step}
var Left fyne.Position = fyne.Position{X: -step, Y: 0}
var Right fyne.Position = fyne.Position{X: step, Y: 0}

const winX float32 = 720
const winY float32 = 720
const winPadding float32 = 10
const step float32 = 10
const speed time.Duration = 60
const gameTitle string = "Snake The Game"

var green = color.NRGBA{R: 0, G: 180, B: 0, A: 255}
var direction fyne.Position = Up

type Snake struct {
	body      []fyne.Position
	direction fyne.Position
	segments  []*canvas.Rectangle
}

type GameCtx struct {
	food    canvas.Rectangle
	window  fyne.Window
	canvas  fyne.Canvas
	content *fyne.Container
	snake   Snake
}

func (ctx *GameCtx) Quit() {
	ctx.window.Close()
}

func (ctx *GameCtx) Init(app fyne.App) {
	ctx.window = app.NewWindow("Snake The Game")
	ctx.window.Resize(fyne.NewSize(winX+winPadding, winY+winPadding))
	ctx.window.SetFixedSize(true)
	ctx.canvas = ctx.window.Canvas()

	centerX := winX / 2
	centerY := winY / 2
	ctx.snake = Snake{
		body: []fyne.Position{
			{X: centerX, Y: centerY - step*3},
			{X: centerX, Y: centerY - step*2},
			{X: centerX, Y: centerY - step},
			{X: centerX, Y: centerY},
		},
		direction: direction,
	}

	ctx.content = container.NewWithoutLayout()
	ctx.snake.Init(ctx.content)

	ctx.canvas.SetContent(ctx.content)
	ctx.setKeys()
}

func (ctx *GameCtx) setKeys() {

	ctx.canvas.SetOnTypedKey(func(k *fyne.KeyEvent) {
		handleMovementKeys(k)
		if k.Name == "Q" {
			ctx.Quit()
		}
	})

}

func (ctx *GameCtx) Run() {
	apple := getRandomRect(ctx.content)
	var collision bool = false
	for {
		time.Sleep(time.Millisecond * speed)
		err := ctx.snake.Render(ctx.content)
		if err != nil {
			GameOver(ctx.window)
		}

		collision = isCollision(ctx.snake.body[0], apple.Position())
		if collision {
			ctx.snake.Eat(apple)
			apple = getRandomRect(ctx.content)
		}
	}
}

func (s *Snake) startMoving() ([]fyne.Position, error) {
	var newPos fyne.Position = NewMovePos(s.body[0], direction)
	for _, pos := range s.body {
		if isCollision(newPos, pos) {
			return nil, fmt.Errorf("Game Over")
		}
	}
	result := append([]fyne.Position{newPos}, s.body[:len(s.body)-1]...)
	s.body = result
	return result, nil
}

func (s *Snake) addTail(rect *canvas.Rectangle) {
	rect.Move(s.body[len(s.body)-1])
	rect.FillColor = green
	s.body = append(s.body, rect.Position())
}

func (s *Snake) Init(content *fyne.Container) {
	for _, pos := range s.body {
		rect := canvas.NewRectangle(green)
		rect.Move(pos)
		rect.Resize(fyne.NewSize(step, step))
		s.segments = append(s.segments, rect)
		content.Add(rect)
	}
}

func (s *Snake) Render(content *fyne.Container) error {
	newPos, err := s.startMoving()
	if err != nil {
		return err
	}
	for i, rect := range s.segments {
		rect.Move(newPos[i])
		rect.Refresh()
	}
	return nil
}

func (s *Snake) Eat(food *canvas.Rectangle) {
	s.addTail(food)
	s.segments = append(s.segments, food)
}

func main() {
	myApp := app.New()
	ctx := GameCtx{}
	ctx.Init(myApp)

	go ctx.Run()

	ctx.window.ShowAndRun()
}

func getRandomRect(c *fyne.Container) *canvas.Rectangle {
	red := color.NRGBA{R: 180, G: 0, B: 0, A: 255}
	rect := canvas.NewRectangle(red)
	pos := float32((rand.IntN(int(winX)) / 10) * 10)
	rect.Move(fyne.NewPos(pos, pos))
	fmt.Println(rect.Position())
	rect.Resize(fyne.NewSize(step, step))
	(*c).Add(rect)
	return rect
}

func handleMovementKeys(k *fyne.KeyEvent) {
	switch k.Name {
	case "J":
		if direction != Up {
			direction = Down
		}
	case "K":
		if direction != Down {
			direction = Up
		}
	case "H":
		if direction != Right {
			direction = Left
		}
	case "L":
		if direction != Left {
			direction = Right
		}
	}
}

func GameOver(parent fyne.Window) {
	dialog.ShowConfirm("Game Over", "Restart?", Restart, parent)
}

func Restart(confirm bool) {
	// TODO: Need to be invistigated
	fmt.Println("Game restarted")
}

func isCollision(pos1 fyne.Position, pos2 fyne.Position) bool {
	return pos1.X == pos2.X && pos1.Y == pos2.Y
}

func NewMovePos(prev_position fyne.Position, move fyne.Position) fyne.Position {
	newPosition := prev_position.Add(move)
	next_y := newPosition.Y
	next_x := newPosition.X

	if next_y+step > winY {
		next_y = next_y - winY
	} else if next_y < 0 {
		next_y = next_y + winY
	}

	if next_x+step > winX {
		next_x -= winX
	} else if next_x < 0 {
		next_x += winY
	}
	newPosition = fyne.Position{X: next_x, Y: next_y}
	fmt.Println(newPosition)
	return newPosition
}
