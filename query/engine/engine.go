package engine

// Connect provides a basic connection config struct for creating the necessary
// mongodb connection object.
type Connect struct {
	Host     string
	AuthDb   string
	Username string
	Pass     string
	Db       string
}

// Engine provides a basic mongodb query controller, ensuring concurrent request
// to the given database instance.
type Engine struct {
}
