package auth

import (
	"fmt"
	"sync"
)

// User represents a user in the system
type User struct {
	Username string
	Password string
}

// AuthManager handles authentication operations
type AuthManager struct {
	users          map[string]User
	activeSessions map[string]bool
	mu             sync.RWMutex
}

// NewAuthManager creates a new authentication manager
func NewAuthManager() *AuthManager {
	return &AuthManager{
		users:          make(map[string]User),
		activeSessions: make(map[string]bool),
	}
}

// Register creates a new user account
func (am *AuthManager) Register(username, password string) error {
	am.mu.Lock()
	defer am.mu.Unlock()

	if _, exists := am.users[username]; exists {
		return fmt.Errorf("user already exists")
	}

	am.users[username] = User{
		Username: username,
		Password: password,
	}

	fmt.Printf("User %s registered successfully\n", username)
	return nil
}

// Login authenticates a user
func (am *AuthManager) Login(username, password string) error {
	am.mu.Lock()
	defer am.mu.Unlock()

	user, exists := am.users[username]
	if !exists {
		return fmt.Errorf("user not found")
	}

	if user.Password != password {
		return fmt.Errorf("invalid password")
	}

	am.activeSessions[username] = true
	fmt.Printf("User %s logged in successfully\n", username)
	return nil
}

// Logout ends a user session
func (am *AuthManager) Logout(username string) error {
	am.mu.Lock()
	defer am.mu.Unlock()

	if !am.activeSessions[username] {
		return fmt.Errorf("user not logged in")
	}

	delete(am.activeSessions, username)
	fmt.Printf("User %s logged out successfully\n", username)
	return nil
}

// IsLoggedIn checks if a user is currently logged in
func (am *AuthManager) IsLoggedIn(username string) bool {
	am.mu.RLock()
	defer am.mu.RUnlock()
	return am.activeSessions[username]
}

func main() {
	auth := NewAuthManager()

	// Register some users
	auth.Register("user1", "password1")
	auth.Register("user2", "password2")

	// Login attempts
	fmt.Println("\nLogin attempts:")
	err := auth.Login("user1", "password1")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}

	err = auth.Login("user2", "wrongpassword") // This will fail
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}

	// Check login status
	fmt.Println("\nLogin status:")
	fmt.Printf("user1 logged in: %v\n", auth.IsLoggedIn("user1"))
	fmt.Printf("user2 logged in: %v\n", auth.IsLoggedIn("user2"))

	// Logout
	fmt.Println("\nLogout:")
	err = auth.Logout("user1")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}

	// Check login status again
	fmt.Println("\nLogin status after logout:")
	fmt.Printf("user1 logged in: %v\n", auth.IsLoggedIn("user1"))
}
