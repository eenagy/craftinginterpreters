package main

type Environment struct {
	values    map[string]interface{}
	enclosing *Environment
}

func NewEnvironment(env *Environment) *Environment {
	return &Environment{
		values:    make(map[string]interface{}),
		enclosing: env,
	}
}

func (env *Environment) Copy() *Environment {
	// Create a new environment with the same enclosing environment
	newEnv := &Environment{
		values:    make(map[string]interface{}),
		enclosing: env.enclosing,
	}

	// Copy the values from the current environment to the new one
	for key, value := range env.values {
		newEnv.values[key] = value
	}

	return newEnv
}
func (e *Environment) Define(name string, value interface{}) {
	e.values[name] = value
}

func (e Environment) Get(name Token) (interface{}, error) {
	value, exists := e.values[name.Lexeme]
	if !exists {
		if e.enclosing != nil {
			return e.enclosing.Get(name)
		}
		return nil, RuntimeError{Operator: name, Message: "Undefined variable '" + name.Lexeme + "'."}
	} else {
		return value, nil
	}
}

func (e *Environment) Assign(name Token, value interface{}) error {
	_, exists := e.values[name.Lexeme]
	if !exists {
		if e.enclosing != nil {
			return e.enclosing.Assign(name, value)
		}
		return RuntimeError{Operator: name, Message: "Undefined variable '" + name.Lexeme + "'."}
	} else {
		e.values[name.Lexeme] = value
		return nil
	}
}
