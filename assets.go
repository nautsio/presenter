package main

import "fmt"

// Asset loads and returns the requested asset.
func Asset(name string) ([]byte, error) {
	return nil, fmt.Errorf("Asset %s not found", name)
}
