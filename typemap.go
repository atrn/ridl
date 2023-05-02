package main

import (
	"encoding/json"
	"io"
	"os"
)

type TypeMap struct {
	GoType    string `json:"go-type"`
	CppType   string `json:"cpp-type"`
	PassByRef bool   `json:"pass-by-ref"`
}

var (
	typeMap map[string]TypeMap

	defaultTypeMap = []TypeMap{
		{"byte", "std::byte", false},
		{"error", "std::runtime_error", true},
		{"string", "std::string", true},
		{"float32", "float", false},
		{"float64", "double", false},
		{"rune", "uint32_t", false},
		{"bool", "bool", false},
		{"float", "double", false},
		{"int", "int", false},
		{"uint", "unsigned int", false},
		{"int8", "int8_t", false},
		{"uint8", "uint8_t", false},
		{"int16", "int16_t", false},
		{"uint16", "uint16_t", false},
		{"int32", "int32_t", false},
		{"uint32", "uint32_t", false},
		{"int64", "int64_t", false},
		{"uint64", "uint64_t", false},
		{"uintptr", "ptrdiff_t", false},
		{"complex32", "std::complex<float>", false},
		{"complex64", "std::complex<double>", false},
	}
)

func initTypeMap() {
	typeMap = make(map[string]TypeMap)
	for _, t := range defaultTypeMap {
		typeMap[t.GoType] = t
	}
}

func readTypeMap(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	var mappings []TypeMap
	d := json.NewDecoder(file)
	err = d.Decode(&mappings)
	if err == nil {
		for _, t := range mappings {
			typeMap[t.GoType] = t
		}
	}
	return nil
}

func mapGoToCpp(goType string) (string, bool) {
	if t, found := typeMap[goType]; found {
		return t.CppType, t.PassByRef
	}
	return goType, false
}

func writeTypeMap(w io.Writer) {
	e := json.NewEncoder(w)
	e.Encode(typeMap)
}
