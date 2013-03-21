package processor

import (
	"log"
)

type CalmingFunction func()

func Calm() {
	if r := recover(); r != nil {
		return
	}
}

func CalmAndLog() {
	if r := recover(); r != nil {
		log.Println(r)
	}
}

func CalmAndLogFunc(name string) CalmingFunction {
	return func() {
		if r := recover(); r != nil {
			log.Printf("%s() : Error %s\n", name, r)
		}
	}
}