package main

import (
	"fmt"
	"math"
	"math/rand"
	. "rabbits_wolfs/animals"
	. "rabbits_wolfs/field"
	. "rabbits_wolfs/movement"
	"time"
)

func move(anml Animal, field Field, direction MoveDirection, moveType string) {
	oldPos := anml.GetPosition()
	var newPos Position
	switch moveType {
	case MoveStep:
		newPos = anml.Step(direction, field.GetSize())
	case MoveRun:
		newPos = anml.Run(direction, field.GetSize())
	case MoveStealth:
		newPos = anml.Stealth(direction, field.GetSize())
	default:
		panic("Unknown moveType")
	}
	field.Move(anml, oldPos, newPos)
}

func getMoveTypes() []string {
	return []string{MoveRun, MoveStep, MoveStealth}
}

func getRandomMoveType() string {
	return getMoveTypes()[rand.Intn(len(getMoveTypes()))]
}

func moveRand(animal Animal, field Field) {
	var direction MoveDirection
	direction.Random()
	move(animal, field, direction, getRandomMoveType())
}

func breed(animals Animals, index int, fld Field) {
	father := animals.GetAnimal(index)
	if !father.CanMakeLife() {
		return
	}
	fld.Each(Position{
		X: father.GetPosition().X - 1,
		Y: father.GetPosition().Y - 1,
	}, Position{
		X: father.GetPosition().X + 1,
		Y: father.GetPosition().Y + 1,
	},
		func(mother Animal) {
			if mother == nil || !mother.CanGiveLife() || rand.Intn(10000) < 5000 {
				return
			}
			animalCount := int(mother.GetRandomBirths())
			for _, newPos := range fld.GetRandomEmptyNear(mother.GetPosition(), mother.GetBirthRadius(), animalCount) {
				animal := animals.NewAnimal()
				animals.Append(animal)
				fld.Spawn(animal, newPos)
			}
		})
}

func checkLife(animals Animals, index int, field Field) bool {
	if animals.GetAnimal(index).Dies() {
		kill(animals, index, field)
		return false
	}
	animals.GetAnimal(index).MakeOlder()
	return true
}

func kill(animals Animals, index int, field Field) {
	field.SetCellByPos(animals.GetAnimal(index).GetPosition(), nil)
	animals.Remove(index)
}

func workAnimal(animals Animals, index int, field Field) {
	breed(animals, index, field)
	moveRand(animals.GetAnimal(index), field)
	//ch <- true
}

func workAnimals(animals Animals, field Field) {
	lngth := animals.Len()
	for i := 0; i < *lngth; i++ {
		if !checkLife(animals, i, field) {
			i--
			if i >= *lngth-1 {
				return
			}
			continue
		}
		workAnimal(animals, i, field)
	}
}

func workOfLoop(data *loopData) {
	tickFinished := false
	drawFinished := false
	for _, animals := range data.animals {
		workAnimals(animals, data.field)
	}
	if data.field.GetSize() <= 35 {
		drawCh := data.field.Draw()
		timer := time.NewTimer(time.Second / 10)
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
	for _, animals := range data.animals {
		data.field.SpawnAnimals(animals)
		maxCount[animals.GetType()] = animals.GetCount()
		startCount[animals.GetType()] = animals.GetCount()
	}
	finishMessage := "Epoch Finished!!! Max count: "

	for i := 0; i < data.time2Live; i++ {
		start := time.Now()
		workOfLoop(data)
		since := time.Since(start)
		if data.field.GetSize() > 35 && since > 0 {
			fmt.Println(since)
		}
		curTotalCount := 0
		for _, animals := range data.animals {
			curCount := animals.GetCount()
			curTotalCount += curCount
			maxCount[animals.GetType()] = int(math.Max(float64(maxCount[animals.GetType()]), float64(curCount)))
		}
		if curTotalCount == 0 {
			finishMessage = "All died!!! Max count: "
			break
		}
	}

	for _, animals := range data.animals {
		totalCount[animals.GetType()] = animals.GetInternalIndex()
	}

	fmt.Println(
		finishMessage,
		maxCount,
		" Start count:",
		startCount,
		" Total count: ",
		totalCount,
	)
}

type loopData struct {
	animals   []Animals
	field     Field
	time2Live int
}
