package file

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"syscall"
	"testing"
)

func TestGetLinesFromNonExistentFile(t *testing.T) {
	data, err := GetLinesFromFilename("non_existent_filename")

	if data != nil {
		t.Errorf("Expected data to be nil, got %v", data)
	}

	if err == nil {
		t.Errorf("Expected error not to be nil, got %v", err)
	}
}

func TestGetLinesFromEmptyFile(t *testing.T) {
	file, _ := ioutil.TempFile("", "TestGetLinesFromFilename")
	defer syscall.Unlink(file.Name())

	data, err := GetLinesFromFilename(file.Name())
	dataType := reflect.TypeOf(data).String()

	if dataType != "[]string" {
		t.Errorf("Expected data to be []string, got %s", dataType)
	}

	if len(data) != 0 {
		t.Errorf("Expected data len to be 0, got %d", len(data))
	}

	if err != nil {
		t.Errorf("Expected error to be nil, got %v", err)
	}
}

func TestGetLinesFromFile(t *testing.T) {
	file, _ := ioutil.TempFile("", "TestGetLinesFromFilename")
	defer syscall.Unlink(file.Name())

	bytes := []byte("one\ntwo\nthree")
	ioutil.WriteFile(file.Name(), bytes, 0644)

	data, err := GetLinesFromFilename(file.Name())

	if fmt.Sprintf("%v", data) != "[one two three]" {
		t.Errorf("got %v", data)
	}

	if err != nil {
		t.Errorf("Expected error to be nil, got %v", err)
	}
}
