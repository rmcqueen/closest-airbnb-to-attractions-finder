package connections

import (
	"reflect"
	"testing"
)

func TestConnect_existingPsqlConnectionReturned(t *testing.T) {
	connection := Init()
	secondConnection := Init()

	if !reflect.DeepEqual(connection, secondConnection) {
		t.Errorf("Expected connection to be re-used, was not.")
	}
}
