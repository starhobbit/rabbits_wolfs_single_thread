package animals

import (
	"math/rand"

	"rabbits_wolfs/movement"
)

const (
	SexMale     = false
	SexFemale   = true
	MoveRun     = "Run"
	MoveStep    = "Step"
	MoveStealth = "Stealth"
)

type (
	speed struct {
		run, step, stealth uint8
	}
	breedAge struct {
		min uint8
		max uint8
	}
	animal struct {
		position    movement.Position
		speed       speed
		typeOf      string
		index       int
		sex         bool
		age         uint8
		breedAge    breedAge
		deathAge    uint8
		maxBirths   uint8
		minBirths   uint8
		birthRadius uint8
		newFunc     func()
	}
)

func (a *animal) Dies() bool {
	return a.age == a.deathAge
}

func (a *animal) Run(direction movement.MoveDirection, max int) movement.Position {
	return a.position.Move(direction, a.speed.run, max)
}

func (a *animal) Step(direction movement.MoveDirection, max int) movement.Position {
	return a.position.Move(direction, a.speed.step, max)
}

func (a *animal) Stealth(direction movement.MoveDirection, max int) movement.Position {
	return a.position.Move(direction, a.speed.stealth, max)
}

func (a *animal) isMale() bool {
	return a.sex == SexMale
}

func (a *animal) isFemale() bool {
	return a.sex == SexFemale
}

func (a *animal) RandomBirths() uint8 {
	if a.isMale() {
		panic("Male can't ")
	}
	return uint8(rand.Intn(int(a.maxBirths-a.minBirths)) + int(a.minBirths))
}

func (a *animal) CanMakeLife() bool {
	return a.isMale() && a.inBreedAge()
}

func (a *animal) CanGiveLife() bool {
	return a.isFemale() && a.inBreedAge()
}

func (a *animal) inBreedAge() bool {
	return a.age >= a.breedAge.min && a.age <= a.breedAge.max
}

func (a *animal) Type() string {
	return a.typeOf
}

func (a *animal) MakeOlder() {
	a.age++
}

func (a *animal) SetPosition(position movement.Position) {
	a.position = position
}

func (a *animal) Position() movement.Position {
	return a.position
}

func (a *animal) Age() uint8 {
	return a.age
}

func (a *animal) BirthRadius() uint8 {
	return a.birthRadius
}

type Animal interface {
	Position() movement.Position
	CanGiveLife() bool
	SetPosition(movement.Position)
	Age() uint8
	RandomBirths() uint8
	BirthRadius() uint8
	Run(movement.MoveDirection, int) movement.Position
	Step(movement.MoveDirection, int) movement.Position
	Stealth(movement.MoveDirection, int) movement.Position
	MakeOlder()
	Type() string
	Dies() bool
}
