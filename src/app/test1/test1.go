package test1

type Mux interface {
	HandleFunc(s int)
	ServeHTTP(s string)
}
