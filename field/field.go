package field

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	. "rabbits_wolfs/animals"
	. "rabbits_wolfs/movement"
	"strconv"
)

const Size = 35

var drawCh chan field
var finishDraw chan bool

type Field interface {
	Move(animal Animal, oldPos Position, newPos Position)
	SpawnAnimals(animals Animals)
	Spawn(animal Animal, position Position)
	GetRandomEmpty(count int) (pos []Position)
	GetRandomEmptyNear(parentPos Position, radius uint8, count int) []Position
	SetCellByPos(pos Position, animal Animal)
	GetSize() int
	Draw() chan bool
	Each(Position, Position, func(Animal))
}

func (f field) Move(animal Animal, oldPos Position, newPos Position) {
	if f.empty(newPos) && oldPos != newPos {
		f.SetCellByPos(oldPos, nil)
		f.Spawn(animal, newPos)
	}
}

func (f field) SpawnAnimals(anmls Animals) {
	var animalCount int
	animalCount = anmls.GetStartCount()
	positions := f.GetRandomEmpty(animalCount)
	for _, pos := range positions {
		animal := anmls.NewAnimal()
		anmls.Append(animal)
		f.Spawn(animal, pos)
	}
}

func (f field) Spawn(animal Animal, position Position) {
	animal.SetPosition(position)
	f.SetCellByPos(animal.GetPosition(), animal)
}

func (f field) GetRandomEmpty(count int) []Position {
	emptyPositions := make([]Position, 0, f.GetSize()*f.GetSize())
	for i := 0; i < f.GetSize()*f.GetSize(); i++ {
		x := i % f.GetSize()
		y := i / f.GetSize()
		pos := Position{X: x, Y: y}
		if f.empty(pos) {
			emptyPositions = append(emptyPositions, pos)
		}
	}
	arrayLen := len(emptyPositions)
	resSize := int(math.Min(float64(arrayLen), float64(count)))
	resPositions := make([]Position, resSize)
	j := 0
	for _, i := range rand.Perm(arrayLen) {
		if j == resSize {
			break
		}
		resPositions[j] = emptyPositions[i]
		j++
	}
	return resPositions
}

func (f field) GetRandomEmptyNear(parentPos Position, radius uint8, count int) []Position {
	radiusInternal := int(radius)
	diametr := radiusInternal*2 + 1
	emptyPositions := make([]Position, 0, diametr)
	for i := 0; i < diametr*diametr; i++ {
		xDelta := i%diametr - radiusInternal
		yDelta := i/diametr - radiusInternal
		pos := Position{
			X: parentPos.X + xDelta,
			Y: parentPos.Y + yDelta,
		}
		if pos.IsValid(f.GetSize()) && f.empty(pos) {
			emptyPositions = append(emptyPositions, pos)
		}
	}
	arrayLen := len(emptyPositions)

	resSize := int(math.Min(float64(arrayLen), float64(count)))
	resPositions := make([]Position, resSize)
	j := 0
	for _, i := range rand.Perm(arrayLen) {
		if j == resSize {
			break
		}
		resPositions[j] = emptyPositions[i]
		j++
	}
	return resPositions
}

func (f field) getCellByPos(pos Position) Animal {
	return f[pos.X][pos.Y]
}

func (f field) SetCellByPos(pos Position, animal Animal) {
	f[pos.X][pos.Y] = animal
}

func (f field) empty(position Position) bool {
	return f.getCellByPos(position) == nil
}

func NewField(size int) Field {
	f := make(field, size)
	for i := range f {
		f[i] = make([]Animal, size)
	}
	drawCh = make(chan field)
	finishDraw = make(chan bool)
	go draw(drawCh, finishDraw)
	return f
}

func (f field) GetSize() int {
	return len(f)
}

func (f field) Draw() chan bool {
	drawCh <- f
	return finishDraw
}

func (f field) String() string {
	var buf bytes.Buffer
	var b [3]byte
	for x, row := range f {
		for y, animal := range row {
			b[0] = byte('-')
			b[1] = byte('-')
			b[2] = byte('-')
			if animal != nil {
				if animal.GetPosition().X != x || animal.GetPosition().Y != y {
					panic("GHOST!!!!")
				}
				switch animal.GetType() {
				case WolfType:
					b[0] = byte('[')
					b[2] = byte(']')
				case RabbitType:
					b[0] = byte('{')
					b[2] = byte('}')
				}
				b[1] = strconv.Itoa(int(animal.GetAge()))[0]
			}
			buf.WriteByte(b[0])
			buf.WriteByte(b[1])
			buf.WriteByte(b[2])
		}
		buf.WriteByte('\n')
	}
	return buf.String()
}

func (f field) Each(upLeftPos Position, botRightPos Position, fu func(animal Animal)) {
	for _, row := range f[FixPos(upLeftPos.X, f.GetSize()):FixPos(botRightPos.X, f.GetSize())] {
		for _, animal := range row[FixPos(upLeftPos.Y, f.GetSize()):FixPos(botRightPos.Y, f.GetSize())] {
			fu(animal)
		}
	}
}

func draw(drawCh chan field, finishDraw chan bool) {
	for fld := range drawCh {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
		fmt.Print(fld)
		finishDraw <- true
	}
}

type field [][]Animal
