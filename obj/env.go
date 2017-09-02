package obj

func NewEnclosedEnv(outer *Env) *Env {
	env := NewEnv()
	env.outer = outer
	return env
}

func NewEnv() *Env {
	s := make(map[string]Obj)
	return &Env{store: s, outer: nil}
}

type Env struct {
	store map[string]Obj
	outer *Env
}

func (e *Env) Get(name string) (Obj, bool) {
	o, ok := e.store[name]
	if !ok && e.outer != nil {
		o, ok = e.outer.Get(name)
	}
	return o, ok
}

func (e *Env) Set(name string, val Obj) Obj {
	e.store[name] = val
	return val
}
