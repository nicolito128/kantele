package gateway

type eventFunc func(any)

type eventHandler map[string][]eventFunc

func (eh eventHandler) Call(eventName string, data any) {
	if eh == nil {
		panic("invalid event handler")
	}

	evList := eh[eventName]
	for _, handle := range evList {
		handle(data)
	}
}

func (eh eventHandler) Append(eventName string, h eventFunc) {
	if eh == nil {
		panic("invalid event handler")
	}

	evList, ok := eh[eventName]
	if !ok {
		evList = make([]eventFunc, 0)
	}

	eh[eventName] = append(evList, h)
}
