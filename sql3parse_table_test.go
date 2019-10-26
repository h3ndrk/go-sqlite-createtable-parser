package parse

import (
	"fmt"
	"testing"
)

func TestFromString(t *testing.T) {
	table, err := FromString("CREATE TABLE a (b INTEGER)")
	fmt.Printf("t: %+v, err: %+v", table, err)
}
