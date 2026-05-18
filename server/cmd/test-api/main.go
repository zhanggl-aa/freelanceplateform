package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Testing project API...")

	// First test local API endpoints
	fmt.Println("\n1. Testing project listing at http://localhost:8080/api/v1/projects...")
	resp, err := http.Get("http://localhost:8080/api/v1/projects")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Status: %d\n", resp.StatusCode)
		resp.Body.Close()
	}

	// Test developer endpoint
	fmt.Println("\n2. Testing developer endpoint at http://localhost:8080/api/v1/developers/b0000000-0000-4000-8000-000000000012...")
	resp, err = http.Get("http://localhost:8080/api/v1/developers/b0000000-0000-4000-8000-000000000012")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Status: %d\n", resp.StatusCode)
		resp.Body.Close()
	}

	// Test with user ID too
	fmt.Println("\n3. Testing developer endpoint with user ID at http://localhost:8080/api/v1/developers/a0000000-0000-4000-8000-000000000012...")
	resp, err = http.Get("http://localhost:8080/api/v1/developers/a0000000-0000-4000-8000-000000000012")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Status: %d\n", resp.StatusCode)
		resp.Body.Close()
	}

	fmt.Println("\nDone! Remember to restart the server if you haven't already!")
}
