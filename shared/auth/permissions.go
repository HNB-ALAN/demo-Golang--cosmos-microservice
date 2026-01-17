// Package auth provides authentication and authorization utilities for USC platform services.
package auth

import (
	"errors"
	"strings"
)

// Permission represents a permission
type Permission struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Resource    string   `json:"resource"`
	Action      string   `json:"action"`
	Conditions  []string `json:"conditions"`
}

// Role represents a role with permissions
type Role struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Permissions []Permission `json:"permissions"`
	Inherits    []string     `json:"inherits"`
}

// PermissionManager manages permissions and roles
type PermissionManager struct {
	roles       map[string]*Role
	permissions map[string]*Permission
}

// NewPermissionManager creates a new permission manager
func NewPermissionManager() *PermissionManager {
	pm := &PermissionManager{
		roles:       make(map[string]*Role),
		permissions: make(map[string]*Permission),
	}

	// Initialize with default permissions and roles
	pm.initializeDefaults()

	return pm
}

// AddPermission adds a new permission
func (pm *PermissionManager) AddPermission(permission *Permission) error {
	if permission.Name == "" {
		return errors.New("permission name cannot be empty")
	}

	if _, exists := pm.permissions[permission.Name]; exists {
		return errors.New("permission already exists")
	}

	pm.permissions[permission.Name] = permission
	return nil
}

// GetPermission retrieves a permission by name
func (pm *PermissionManager) GetPermission(name string) (*Permission, error) {
	permission, exists := pm.permissions[name]
	if !exists {
		return nil, errors.New("permission not found")
	}
	return permission, nil
}

// AddRole adds a new role
func (pm *PermissionManager) AddRole(role *Role) error {
	if role.Name == "" {
		return errors.New("role name cannot be empty")
	}

	if _, exists := pm.roles[role.Name]; exists {
		return errors.New("role already exists")
	}

	// Validate permissions
	for _, permName := range role.Permissions {
		if _, exists := pm.permissions[permName.Name]; !exists {
			return errors.New("permission not found: " + permName.Name)
		}
	}

	pm.roles[role.Name] = role
	return nil
}

// GetRole retrieves a role by name
func (pm *PermissionManager) GetRole(name string) (*Role, error) {
	role, exists := pm.roles[name]
	if !exists {
		return nil, errors.New("role not found")
	}
	return role, nil
}

// GetRolePermissions returns all permissions for a role (including inherited)
func (pm *PermissionManager) GetRolePermissions(roleName string) ([]Permission, error) {
	role, err := pm.GetRole(roleName)
	if err != nil {
		return nil, err
	}

	permissions := make([]Permission, 0)
	permissionMap := make(map[string]bool)

	// Add direct permissions
	for _, perm := range role.Permissions {
		if !permissionMap[perm.Name] {
			permissions = append(permissions, perm)
			permissionMap[perm.Name] = true
		}
	}

	// Add inherited permissions
	for _, inheritedRoleName := range role.Inherits {
		inheritedPermissions, err := pm.GetRolePermissions(inheritedRoleName)
		if err != nil {
			continue // Skip invalid inherited roles
		}

		for _, perm := range inheritedPermissions {
			if !permissionMap[perm.Name] {
				permissions = append(permissions, perm)
				permissionMap[perm.Name] = true
			}
		}
	}

	return permissions, nil
}

// HasPermission checks if a role has a specific permission
func (pm *PermissionManager) HasPermission(roleName, permissionName string) bool {
	permissions, err := pm.GetRolePermissions(roleName)
	if err != nil {
		return false
	}

	for _, perm := range permissions {
		if perm.Name == permissionName {
			return true
		}
	}

	return false
}

// CheckPermission checks if user has permission for specific resource and action
func (pm *PermissionManager) CheckPermission(userRole, resource, action string) bool {
	permissions, err := pm.GetRolePermissions(userRole)
	if err != nil {
		return false
	}

	for _, perm := range permissions {
		if perm.Resource == resource && perm.Action == action {
			return true
		}
	}

	return false
}

// CheckPermissionWithConditions checks permission with additional conditions
func (pm *PermissionManager) CheckPermissionWithConditions(userRole, resource, action string, conditions map[string]string) bool {
	permissions, err := pm.GetRolePermissions(userRole)
	if err != nil {
		return false
	}

	for _, perm := range permissions {
		if perm.Resource == resource && perm.Action == action {
			// Check conditions if any
			if len(perm.Conditions) == 0 {
				return true
			}

			// All conditions must be met
			for _, condition := range perm.Conditions {
				if !pm.evaluateCondition(condition, conditions) {
					return false
				}
			}
			return true
		}
	}

	return false
}

// evaluateCondition evaluates a permission condition
func (pm *PermissionManager) evaluateCondition(condition string, context map[string]string) bool {
	// Simple condition evaluation
	// Format: "key:value" or "key!=value"
	if strings.Contains(condition, "!=") {
		parts := strings.Split(condition, "!=")
		if len(parts) != 2 {
			return false
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		return context[key] != value
	}

	if strings.Contains(condition, ":") {
		parts := strings.Split(condition, ":")
		if len(parts) != 2 {
			return false
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		return context[key] == value
	}

	return false
}

// ListRoles returns all available roles
func (pm *PermissionManager) ListRoles() []string {
	roles := make([]string, 0, len(pm.roles))
	for roleName := range pm.roles {
		roles = append(roles, roleName)
	}
	return roles
}

// ListPermissions returns all available permissions
func (pm *PermissionManager) ListPermissions() []string {
	permissions := make([]string, 0, len(pm.permissions))
	for permName := range pm.permissions {
		permissions = append(permissions, permName)
	}
	return permissions
}

// initializeDefaults initializes default permissions and roles
func (pm *PermissionManager) initializeDefaults() {
	// Default permissions
	defaultPermissions := []Permission{
		// User management
		{Name: "user:read", Description: "Read user information", Resource: "user", Action: "read"},
		{Name: "user:write", Description: "Create/update user information", Resource: "user", Action: "write"},
		{Name: "user:delete", Description: "Delete user", Resource: "user", Action: "delete"},
		{Name: "user:list", Description: "List users", Resource: "user", Action: "list"},

		// Content management
		{Name: "content:read", Description: "Read content", Resource: "content", Action: "read"},
		{Name: "content:write", Description: "Create/update content", Resource: "content", Action: "write"},
		{Name: "content:delete", Description: "Delete content", Resource: "content", Action: "delete"},
		{Name: "content:list", Description: "List content", Resource: "content", Action: "list"},

		// Analytics
		{Name: "analytics:read", Description: "Read analytics data", Resource: "analytics", Action: "read"},
		{Name: "analytics:write", Description: "Write analytics data", Resource: "analytics", Action: "write"},

		// System administration
		{Name: "system:read", Description: "Read system information", Resource: "system", Action: "read"},
		{Name: "system:write", Description: "Modify system settings", Resource: "system", Action: "write"},
		{Name: "system:admin", Description: "Full system administration", Resource: "system", Action: "admin"},

		// API access
		{Name: "api:read", Description: "Read API data", Resource: "api", Action: "read"},
		{Name: "api:write", Description: "Write API data", Resource: "api", Action: "write"},
		{Name: "api:admin", Description: "API administration", Resource: "api", Action: "admin"},
	}

	// Add default permissions
	for _, perm := range defaultPermissions {
		pm.permissions[perm.Name] = &perm
	}

	// Default roles
	defaultRoles := []Role{
		{
			Name:        "admin",
			Description: "System administrator with full access",
			Permissions: []Permission{
				{Name: "user:read"}, {Name: "user:write"}, {Name: "user:delete"}, {Name: "user:list"},
				{Name: "content:read"}, {Name: "content:write"}, {Name: "content:delete"}, {Name: "content:list"},
				{Name: "analytics:read"}, {Name: "analytics:write"},
				{Name: "system:read"}, {Name: "system:write"}, {Name: "system:admin"},
				{Name: "api:read"}, {Name: "api:write"}, {Name: "api:admin"},
			},
		},
		{
			Name:        "moderator",
			Description: "Content moderator with content management access",
			Permissions: []Permission{
				{Name: "user:read"}, {Name: "user:list"},
				{Name: "content:read"}, {Name: "content:write"}, {Name: "content:delete"}, {Name: "content:list"},
				{Name: "analytics:read"},
				{Name: "api:read"},
			},
		},
		{
			Name:        "user",
			Description: "Regular user with basic access",
			Permissions: []Permission{
				{Name: "user:read"},
				{Name: "content:read"}, {Name: "content:write"},
				{Name: "api:read"},
			},
		},
		{
			Name:        "guest",
			Description: "Guest user with read-only access",
			Permissions: []Permission{
				{Name: "content:read"},
				{Name: "api:read"},
			},
		},
	}

	// Add default roles
	for _, role := range defaultRoles {
		pm.roles[role.Name] = &role
	}
}

// ValidateUserPermissions validates user permissions against required permissions
func (pm *PermissionManager) ValidateUserPermissions(userRole string, requiredPermissions []string) (bool, []string) {
	userPermissions, err := pm.GetRolePermissions(userRole)
	if err != nil {
		return false, []string{"Invalid user role"}
	}

	userPermMap := make(map[string]bool)
	for _, perm := range userPermissions {
		userPermMap[perm.Name] = true
	}

	missingPermissions := make([]string, 0)
	for _, required := range requiredPermissions {
		if !userPermMap[required] {
			missingPermissions = append(missingPermissions, required)
		}
	}

	return len(missingPermissions) == 0, missingPermissions
}

// GetResourcePermissions returns all permissions for a specific resource
func (pm *PermissionManager) GetResourcePermissions(resource string) []Permission {
	permissions := make([]Permission, 0)
	for _, perm := range pm.permissions {
		if perm.Resource == resource {
			permissions = append(permissions, *perm)
		}
	}
	return permissions
}

// GetActionPermissions returns all permissions for a specific action
func (pm *PermissionManager) GetActionPermissions(action string) []Permission {
	permissions := make([]Permission, 0)
	for _, perm := range pm.permissions {
		if perm.Action == action {
			permissions = append(permissions, *perm)
		}
	}
	return permissions
}
