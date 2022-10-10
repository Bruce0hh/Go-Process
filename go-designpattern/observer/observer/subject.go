package observer

type Subject interface {
	Register(observer Observer)
	deregister(observer Observer)
	notifyAll()
}
