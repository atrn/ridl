package main

import (
	"encoding/json"
	"os"
	"testing"
)

func TestTypeMap1(t *testing.T) {
	const filename = "test-map.json"

	defer os.Remove(filename)

	f, err := os.Create(filename)
	if err != nil {
		t.Fatal(err)
	}

	var mapping []TypeMap
	mapping = append(mapping, TypeMap{"Timepoint", "std::chrono::steady_clock::timepoint", false})
	mapping = append(mapping, TypeMap{"StringMap", "std::map<std::string, std::string>", true})

	e := json.NewEncoder(f)
	err = e.Encode(mapping)
	if err != nil {
		f.Close()
		os.Remove(f.Name())
		t.Fatal(err)
	}
	err = f.Close()
	if err != nil {
		t.Fatal(err)
	}

	initTypeMap()

	cpp, ref := mapGoToCpp("int")
	if cpp != "int" {
		t.Fatalf("Go %q mapped to C++ %q", "int", cpp)
	}
	if ref {
		t.Fatalf("Go \"int\" mapped to pass-by-ref type %q", cpp)
	}

	err = readTypeMap(filename)
	if err != nil {
		t.Fatal(err)
	}

	cpp, ref = mapGoToCpp("Timepoint")
	if cpp != "std::chrono::steady_clock::timepoint" {
		t.Fatalf("Go %q mapped to C++ %q", "Timepoint", cpp)
	}

}
