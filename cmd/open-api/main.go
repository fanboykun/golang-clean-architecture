package main

import (
	"encoding/json"
	"golang-clean-architecture/internal/delivery/http/route"
	"log"
	"os"
	"path/filepath"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// Create a minimal Fiber app (won't be started)
	app := fiber.New()

	// Create Huma API with the same config as the main app
	humaConfig := huma.DefaultConfig("Backend API", "1.0.0")
	humaConfig.Servers = []*huma.Server{
		{URL: "http://localhost:8080", Description: "Development server"},
		{URL: "https://api.example.com", Description: "Production server"},
	}

	// Add security scheme for bearer token
	humaConfig.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"bearer": {
			Type:         "http",
			Scheme:       "bearer",
			BearerFormat: "JWT",
			Description:  "Bearer token authentication",
		},
	}

	api := humafiber.New(app, humaConfig)

	// Register all routes using the shared route configuration
	// We only need the API for OpenAPI generation, so we pass nil for controllers
	routeConfig := route.RouteConfig{
		App:               app,
		Api:               api,
		UserController:    nil,
		ContactController: nil,
		AddressController: nil,
		AuthMiddleware:    nil,
	}

	// Register only the Huma operations (skip Fiber routes since we're not running the server)
	registerOpenAPIRoutes(&routeConfig)

	// Get the OpenAPI spec
	spec := api.OpenAPI()

	// Marshal to JSON
	data, err := json.MarshalIndent(spec, "", "  ")
	if err != nil {
		log.Fatal("Failed to marshal OpenAPI spec:", err)
	}

	// Determine output file
	outputFile := "api/openapi.json"
	if len(os.Args) > 1 {
		outputFile = os.Args[1]
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(outputFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatal("Failed to create output directory:", err)
	}

	// Write to file
	if err := os.WriteFile(outputFile, data, 0644); err != nil {
		log.Fatal("Failed to write OpenAPI spec:", err)
	}

	log.Printf("✅ OpenAPI spec generated successfully: %s\n", outputFile)
	log.Printf("📊 Total operations: %d\n", countOperations(spec))
	log.Printf("🏷️  Tags: %v\n", getTags(spec))
}

func countOperations(spec *huma.OpenAPI) int {
	count := 0
	for _, pathItem := range spec.Paths {
		if pathItem.Get != nil {
			count++
		}
		if pathItem.Post != nil {
			count++
		}
		if pathItem.Put != nil {
			count++
		}
		if pathItem.Patch != nil {
			count++
		}
		if pathItem.Delete != nil {
			count++
		}
	}
	return count
}

func getTags(spec *huma.OpenAPI) []string {
	tagMap := make(map[string]bool)
	for _, pathItem := range spec.Paths {
		operations := []*huma.Operation{
			pathItem.Get, pathItem.Post, pathItem.Put,
			pathItem.Patch, pathItem.Delete,
		}
		for _, op := range operations {
			if op != nil {
				for _, tag := range op.Tags {
					tagMap[tag] = true
				}
			}
		}
	}

	tags := make([]string, 0, len(tagMap))
	for tag := range tagMap {
		tags = append(tags, tag)
	}
	return tags
}

func registerOpenAPIRoutes(c *route.RouteConfig) {
	// This is a simplified version that only registers Huma operations
	// without the Fiber handlers (since we're not running the server)
	route.RegisterHumaOperations(c)
}
