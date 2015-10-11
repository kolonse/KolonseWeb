package KolonseWeb

import (
	. "github.com/kolonse/KolonseWeb/Type"
)

type MiddleWares struct {
	Do       DoStep
	nextStep bool
}

func (middleWares *MiddleWares) Next() {
	middleWares.nextStep = true
}

func (middleWares *MiddleWares) IsGoNext() bool {
	return middleWares.nextStep
}

func NewMiddleWares() *MiddleWares {
	mdiddleWares := &MiddleWares{}
	mdiddleWares.nextStep = false
	mdiddleWares.Do = DefaultDoStep
	return mdiddleWares
}
