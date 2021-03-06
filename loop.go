package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"rabbits_wolfs/animals"
	"rabbits_wolfs/field"
	"rabbits_wolfs/movement"
)

type (
	Field interface {
		Move(*animals.Animal, movement.Position, movement.Position)
		SpawnAnimals(field.Animals)
		Spawn(*animals.Animal, movement.Position)
		GetRandomEmpty(int) []movement.Position
		GetRandomEmptyNear(movement.Position, uint8, int) []movement.Position
		SetCellByPos(movement.Position, *animals.Animal)
		GetSize() int
		Draw() chan bool
		Each(movement.Position, movement.Position, func(*animals.Animal))
	}
)

func move(animal *animals.Animal, field Field, direction movement.MoveDirection, moveType string) {
	var newPos movement.Position
	oldPos := animal.Position()
	switch moveType {
	case animals.MoveStep:
		newPos = animal.Step(direction, field.GetSize())
	case animals.MoveRun:
		newPos = animal.Run(direction, field.GetSize())
	case animals.MoveStealth:
		newPos = animal.Stealth(direction, field.GetSize())
	default:
		panic("Unknown moveType")
	}
	field.Move(animal, oldPos, newPos)
}

func getMoveTypes() []string {
	return []string{animals.MoveRun, animals.MoveStep, animals.MoveStealth}
}

func getRandomMoveType() string {
	return getMoveTypes()[rand.Intn(len(getMoveTypes()))]
}

func moveRand(animal *animals.Animal, field Field) {
	var direction movement.MoveDirection
	direction.Random()
	move(animal, field, direction, getRandomMoveType())
}

func breed(creatures animals.Animals, index int, fld Field) {
	father := creatures.Animal(index)
	if !father.CanMakeLife() {
		return
	}
	fld.Each(movement.Position{
		X: father.Position().X - 1,
		Y: father.Position().Y - 1,
	}, movement.Position{
		X: father.Position().X + 1,
		Y: father.Position().Y + 1,
	},
		func(mother *animals.Animal) {
			if mother == nil || !mother.CanGiveLife() || rand.Intn(99) < 49 {
				return
			}
			animalCount := int(mother.RandomBirths())
			for _, newPos := range fld.GetRandomEmptyNear(mother.Position(), mother.BirthRadius(), animalCount) {
				animal := creatures.NewAnimal()
				creatures.Append(animal)
				fld.Spawn(animal, newPos)
			}
		})
}

func checkLife(creatures animals.Animals, index int, field Field) bool {
	if creatures.Animal(index).Dies() {
		kill(creatures, index, field)
		return false
	}
	creatures.Animal(index).MakeOlder()
	return true
}

func kill(creatures animals.Animals, index int, field Field) {
	field.SetCellByPos(creatures.Animal(index).Position(), nil)
	creatures.Remove(index)
}

func workAnimal(creatures animals.Animals, index int, field Field) {
	breed(creatures, index, field)
	moveRand(creatures.Animal(index), field)
}

func workAnimals(creatures animals.Animals, field Field) {
	// use pointer because length can become shorter inside for loop
	length := creatures.Len()
	for i := 0; i < *length; i++ {
		if !checkLife(creatures, i, field) {
			i--
			if i >= *length-1 {
				return
			}
			continue
		}
		workAnimal(creatures, i, field)
	}
}

func workOfLoop(data *loopData) {
	tickFinished := false
	drawFinished := false
	for _, creatures := range data.animals {
		workAnimals(creatures, data.field)
	}
	if data.field.GetSize() <= field.ScreenSize() {
		drawCh := data.field.Draw()
		timer := time.NewTimer(time.Second / 30)
		for !(tickFinished && drawFinished) {
			select {
			case <-drawCh:
				drawFinished = true
			case <-timer.C:
				tickFinished = true
			}
		}
	}
}

func loop(data *loopData) {
	maxCount := make(map[string]int)
	startCount := make(map[string]int)
	totalCount := make(map[string]int)

	for _, creatures := range data.animals {
		start := time.Now()
		data.field.SpawnAnimals(creatures)
		maxCount[creatures.Type()] = creatures.Count()
		startCount[creatures.Type()] = creatures.Count()
		since := time.Since(start)
		if data.field.GetSize() > field.ScreenSize() && since > 0 {
			fmt.Println("spawn ", creatures.Type(), " : ", since)
		}
	}
	finishMessage := "Epoch Finished!!! Max count: "
	i := 0
	for ; i < data.time2Live; i++ {
		start := time.Now()
		workOfLoop(data)
		since := time.Since(start)
		if data.field.GetSize() > field.ScreenSize() && since > 0 {
			fmt.Println("loop: ", since)
		}
		curTotalCount := 0
		for _, creatures := range data.animals {
			curCount := creatures.Count()
			curTotalCount += curCount
			maxCount[creatures.Type()] = int(math.Max(float64(maxCount[creatures.Type()]), float64(curCount)))
		}
		if curTotalCount == 0 {
			finishMessage = "All died!!! Max count: "
			break
		}
	}

	for _, creatures := range data.animals {
		totalCount[creatures.Type()] = creatures.InternalIndex()
	}

	fmt.Println(
		finishMessage,
		maxCount,
		" Start count:",
		startCount,
		" Total count: ",
		totalCount,
		" Iterations: ",
		i,
	)
}

type loopData struct {
	animals   []animals.Animals
	field     Field
	time2Live int
}
