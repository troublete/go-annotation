package internal

import "github.com/troublete/go-chariot/inspect"

type Processor func(map[string]string, inspect.Function) error

type ProcessorRegister map[string]Processor

func NewProcessorRegister() ProcessorRegister {
	return map[string]Processor{}
}

func (pc ProcessorRegister) Register(name string, proc Processor) ProcessorRegister {
	pc[name] = proc
	return pc
}
