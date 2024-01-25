package intergation

type GenericHandler[T any] func(msgType byte, T any) error

type RouteMsg struct{}

func RouteMsgHandler(msgType byte, msg *RouteMsg) error {
	return nil
}

func NewGenericMap() {
	router := make(map[byte]any)
	router[1] = RouteMsgHandler
}
