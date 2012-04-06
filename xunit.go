package xunit

import (
	"io"
	"os"
	"time"
)

type Test struct {
	Name string
	Time time.Duration
	Error string
}

func parseOutput(reader io.Reader) ([]Test, error) {
	tests := []Test{}

	return tests, nil
}

func main() {
	parseOutput(os.Stdin)
}
