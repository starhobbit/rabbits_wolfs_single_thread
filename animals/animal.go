package animals

import (
	"math/rand"
	. "rabbits_wolfs/movement"
)

const SexMale = false
const SexFemale = true

const MoveRun = "Run"
const MoveStep = "Step"
const MoveStealth = "Stealth"

func (a *animals) Len() *int {
	return &a.count
}

func (a *animals) Append(newAnimal *animal) {
	if a.size <= a.count {
		a.animals = append(a.animals, newAnimal)
		a.size++
	} else {
		a.animals[a.lastIndex] = newAnimal
	}
	newAnimal.index = a.lastIndex
	a.count++
	a.lastIndex++
	a.internalIndex++
}

func (a *animals) Remove(index int) {
	lastIndex := a.lastIndex - 1
	if lastIndex <= 0 || lastIndex == index {
		a.animals[index] = nil
	} else {
		a.animals[index] = a.animals[lastIndex]
		a.animals[lastIndex] = nil
		a.animals[index].index = index
	}

	a.count--
	a.lastIndex--
}

func (a *animals) GetStartCount() int {
	return a.startCount
}

type animals struct {
	animals       []*animal
	count         int
	lastIndex     int
	internalIndex int
	startCount    int
	newAnimal     func() *animal
	animalType    string
	size          int
}

func (a *animals) NewAnimal() *animal {
	return a.newAnimal()
}

func (a *animals) GetInternalIndex() int {
	return a.internalIndex
}

func (a *animals) GetCount() int {
	return a.count
}

func (a *animals) GetAnimal(index int) *animal {
	return a.animals[index]
}

func (a *animals) GetType() string {
	return a.animalType
}

func NewAnimals(startCount int, newAnimal func() *animal, animalType string) Animals {
	return &animals{startCount: startCount, newAnimal: newAnimal, animalType: animalType}
}

func (a *animal) Dies() bool {
	return a.age == a.deathAge
}

func (a *animal) Run(direction MoveDirection, max int) Position {
	return a.position.Move(direction, a.speed.run, max)
}

func (a *animal) Step(direction MoveDirection, max int) Position {
	return a.position.Move(direction, a.speed.step, max)
}

func (a *animal) Stealth(direction MoveDirection, max int) Position {
	return a.position.Move(direction, a.speed.stealth, max)
}

func (a *animal) isMale() bool {
	return a.sex == SexMale
}

func (a *animal) isFemale() bool {
	return a.sex == SexFemale
}

func (a *animal) GetRandomBirths() uint8 {
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

func (a *animal) GetType() string {
	return a.typeOf
}

type animal struct {
	position    Position
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

func (a *animal) MakeOlder() {
	a.age++
}

func (a *animal) SetPosition(position Position) {
	a.position = position
}

func (a *animal) GetPosition() Position {
	return a.position
}

func (a *animal) GetAge() uint8 {
	return a.age
}

func (a *animal) GetBirthRadius() uint8 {
	return a.birthRadius
}

type speed struct {
	run, step, stealth uint8
}

type breedAge struct {
	min uint8
	max uint8
}

type Animals interface {
	Len() *int
	Append(*animal)
	Remove(int)
	GetStartCount() int
	GetCount() int
	GetInternalIndex() int
	GetAnimal(int) *animal
	NewAnimal() *animal
	GetType() string
}

type Animal interface {
	GetPosition() Position
	CanGiveLife() bool
	SetPosition(Position)
	GetAge() uint8
	GetRandomBirths() uint8
	GetBirthRadius() uint8
	Run(MoveDirection, int) Position
	Step(MoveDirection, int) Position
	Stealth(MoveDirection, int) Position
	MakeOlder()
	GetType() string
}
