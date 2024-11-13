package graph

import "graph-api/internal/database"

// Resolver serves as dependency injection for your app.
type Resolver struct {
	User *database.User
}
