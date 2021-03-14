package field

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"strconv"

	"github.com/inancgumus/screen"

	"rabbits_wolfs/animals"
	"rabbits_wolfs/movement"
)

var (
	drawCh     chan field
	finishDraw chan bool
	screenSize = -1
)

type (
	Animals interface {
		StartCount() int
		NewAnimal() *animals.Animal
		Append(animal *animals.Animal)
	}
)

func (f field) Move(animal *animals.Animal, oldPos movement.Position, newPos movement.Position) {
	if f.empty(newPos) && oldPos != newPos {
		f.SetCellByPos(oldPos, nil)
		f.Spawn(animal, newPos)
	}
}

func (f field) SpawnAnimals(creatures Animals) {
	var animalCount int
	animalCount = creatures.StartCount()
	positions := f.GetRandomEmpty(animalCount)
	for _, pos := range positions {
		animal := creatures.NewAnimal()
		creatures.Append(animal)
		f.Spawn(animal, pos)
	}
}

func (f field) Spawn(animal *animals.Animal, position movement.Position) {
	animal.SetPosition(position)
	f.SetCellByPos(animal.Position(), animal)
}

func (f field) GetRandomEmpty(count int) []movement.Position {
	emptyPositions := make([]movement.Position, 0, f.GetSize()*f.GetSize())
	for i := 0; i < f.GetSize()*f.GetSize(); i++ {
		x := i % f.GetSize()
		y := i / f.GetSize()
		pos := movement.Position{X: x, Y: y}
		if f.empty(pos) {
			emptyPositions = append(emptyPositions, pos)
		}
	}
	arrayLen := len(emptyPositions)
	resSize := int(math.Min(float64(arrayLen), float64(count)))
	resPositions := make([]movement.Position, resSize)
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

func (f field) GetRandomEmptyNear(parentPos movement.Position, radius uint8, count int) []movement.Position {
	radiusInternal := int(radius)
	diameter := radiusInternal*2 + 1
	emptyPositions := make([]movement.Position, 0, diameter)
	for i := 0; i < diameter*diameter; i++ {
		xDelta := i%diameter - radiusInternal
		yDelta := i/diameter - radiusInternal
		pos := movement.Position{
			X: parentPos.X + xDelta,
			Y: parentPos.Y + yDelta,
		}
		if pos.IsValid(f.GetSize()) && f.empty(pos) {
			emptyPositions = append(emptyPositions, pos)
		}
	}
	arrayLen := len(emptyPositions)

	resSize := int(math.Min(float64(arrayLen), float64(count)))
	resPositions := make([]movement.Position, resSize)
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

func (f field) getCellByPos(pos movement.Position) *animals.Animal {
	return f[pos.X][pos.Y]
}

func (f field) SetCellByPos(pos movement.Position, animal *animals.Animal) {
	f[pos.X][pos.Y] = animal
}

func (f field) empty(position movement.Position) bool {
	return f.getCellByPos(position) == nil
}

func NewField(size int) *field {
	f := make(field, size)
	for i := range f {
		f[i] = make([]*animals.Animal, size)
	}
	if size <= ScreenSize() {
		drawCh = make(chan field)
		finishDraw = make(chan bool)
		go draw(drawCh, finishDraw)
	}
	return &f
}

func (f field) GetSize() int {
	return len(f)
}

func (f field) Draw() chan bool {
	drawCh <- f
	return finishDraw
}

func (f field) String() string {
	var (
		buf bytes.Buffer
		b   [3]byte
	)
	for x, row := range f {
		for y, animal := range row {
			b[0] = byte('-')
			b[1] = byte('-')
			b[2] = byte('-')
			if animal != nil {
				if animal.Position().X != x || animal.Position().Y != y {
					panic("GHOST!!!!")
				}
				switch animal.Type() {
				case animals.WolfType:
					b[0] = byte('[')
					b[2] = byte(']')
				case animals.RabbitType:
					b[0] = byte('{')
					b[2] = byte('}')
				}
				b[1] = strconv.Itoa(int(animal.Age()))[0]
			}
			buf.WriteByte(b[0])
			buf.WriteByte(b[1])
			buf.WriteByte(b[2])
		}
		buf.WriteByte('\n')
	}
	return buf.String()
}

func (f field) Each(upLeftPos movement.Position, botRightPos movement.Position, fu func(animal *animals.Animal)) {
	for _, row := range f[movement.FixPos(upLeftPos.X, f.GetSize()):movement.FixPos(botRightPos.X, f.GetSize())] {
		for _, animal := range row[movement.FixPos(upLeftPos.Y, f.GetSize()):movement.FixPos(botRightPos.Y, f.GetSize())] {
			fu(animal)
		}
	}
}

func draw(drawCh chan field, finishDraw chan bool) {
	screen.Clear()
	for fld := range drawCh {
		screen.MoveTopLeft()
		fmt.Print(fld)
		finishDraw <- true
	}
}

func ScreenSize() int {
	if screenSize == -1 {
		xSize, ySize := screen.Size()
		screenSize = int(math.Min(float64(xSize/3), float64(ySize-1)))
	}
	return screenSize
}

type field [][]*animals.Animal
