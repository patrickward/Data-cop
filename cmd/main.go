package main

import "fmt"

func main() {
	fmt.Println("Demonstrating nil map behavior in Go")
	fmt.Println("====================================")

	// Create a nil map (not initialized)
	var nilMap map[string][]string

	// Is map actually nil?
	fmt.Println("\n0. Is the map nil?")
	fmt.Printf("nilMap == nil: %v\n", nil == nilMap)

	// Demonstrate lookup operations
	fmt.Println("\n1. Map Lookup Operations:")
	value, exists := nilMap["any_key"]
	fmt.Printf("value, exists := nilMap[\"any_key\"]\n")
	fmt.Printf("- Value: %v (zero value of []string)\n", value)
	fmt.Printf("- Exists: %v\n", exists)

	// Demonstrate length operations
	fmt.Println("\n2. Length Operations:")
	length := len(nilMap)
	fmt.Printf("len(nilMap) = %d\n", length)

	// Demonstrate range operations
	fmt.Println("\n3. Range Operations:")
	fmt.Println("for k, v := range nilMap {")
	itemCount := 0
	for k, v := range nilMap {
		itemCount++
		fmt.Printf("  - Key: %v, Value: %v\n", k, v)
	}
	fmt.Println("}")
	fmt.Printf("Items iterated: %d\n", itemCount)

	// Demonstrate safe comparisons
	fmt.Println("\n4. Safe Comparisons:")
	fmt.Printf("nilMap == nil: %v\n", nilMap == nil)

	// Demonstrate safe field access patterns
	fmt.Println("\n5. Common Safe Patterns:")
	slice, ok := nilMap["test"]
	if ok && len(slice) > 0 {
		fmt.Println("  - This won't print because lookup returns false")
	}

	// Show what would panic (commented out)
	fmt.Println("\n6. Operations that would panic (commented out):")
	fmt.Println("// nilMap[\"key\"] = []string{\"value\"}  // PANIC: assignment to nil map")

	// Demonstrate proper initialization
	fmt.Println("\n7. Proper initialization fixes panic:")
	nilMap = make(map[string][]string)
	nilMap["key"] = []string{"value"} // Now safe
	fmt.Printf("After initialization: %v\n", nilMap)
}
