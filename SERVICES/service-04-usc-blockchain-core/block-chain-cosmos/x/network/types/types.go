package types

import (
	"fmt"
	"time"
)

const (
	// ModuleName defines the module name
	ModuleName = "network"

	// RouterKey defines the message route for the network module
	RouterKey = ModuleName

	// QuerierRoute defines the querier route for the network module
	QuerierRoute = ModuleName
)

// Event types
const (
	EventTypeNetworkCreated        = "network_created"
	EventTypeNetworkUpdated        = "network_updated"
	EventTypeNetworkDeleted        = "network_deleted"
	EventTypeNodeJoined            = "node_joined"
	EventTypeNodeLeft              = "node_left"
	EventTypeConnectionEstablished = "connection_established"
	EventTypeConnectionLost        = "connection_lost"
	EventTypeNetworkSync           = "network_sync"
	EventTypeNetworkHealth         = "network_health"

	AttributeKeyNetworkID    = "network_id"
	AttributeKeyNodeID       = "node_id"
	AttributeKeyConnectionID = "connection_id"
	AttributeKeyModule       = ModuleName
)

// NetworkStatus represents the status of a network
type NetworkStatus string

const (
	NetworkStatusActive      NetworkStatus = "active"
	NetworkStatusInactive    NetworkStatus = "inactive"
	NetworkStatusMaintenance NetworkStatus = "maintenance"
	NetworkStatusError       NetworkStatus = "error"
)

// NodeStatus represents the status of a node
type NodeStatus string

const (
	NodeStatusOnline     NodeStatus = "online"
	NodeStatusOffline    NodeStatus = "offline"
	NodeStatusConnecting NodeStatus = "connecting"
	NodeStatusError      NodeStatus = "error"
)

// ConnectionStatus represents the status of a connection
type ConnectionStatus string

const (
	ConnectionStatusActive   ConnectionStatus = "active"
	ConnectionStatusInactive ConnectionStatus = "inactive"
	ConnectionStatusPending  ConnectionStatus = "pending"
	ConnectionStatusError    ConnectionStatus = "error"
)

// Network represents a network configuration
type Network struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Status      NetworkStatus     `json:"status"`
	Nodes       []Node            `json:"nodes"`
	Connections []Connection      `json:"connections"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Metadata    map[string]string `json:"metadata"`
}

// Node represents a network node
type Node struct {
	ID        string            `json:"id"`
	NetworkID string            `json:"network_id"`
	Name      string            `json:"name"`
	Address   string            `json:"address"`
	Port      int               `json:"port"`
	Status    NodeStatus        `json:"status"`
	LastSeen  time.Time         `json:"last_seen"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	Metadata  map[string]string `json:"metadata"`
}

// Connection represents a network connection
type Connection struct {
	ID         string            `json:"id"`
	NetworkID  string            `json:"network_id"`
	FromNodeID string            `json:"from_node_id"`
	ToNodeID   string            `json:"to_node_id"`
	Status     ConnectionStatus  `json:"status"`
	Latency    int64             `json:"latency"`   // milliseconds
	Bandwidth  int64             `json:"bandwidth"` // bytes per second
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
	Metadata   map[string]string `json:"metadata"`
}

// NetworkSync represents network synchronization status
type NetworkSync struct {
	ID            string            `json:"id"`
	NetworkID     string            `json:"network_id"`
	NodeID        string            `json:"node_id"`
	Status        string            `json:"status"`   // e.g., "syncing", "synced", "error"
	Progress      int64             `json:"progress"` // 0 to 100 (percentage)
	StartHeight   int64             `json:"start_height"`
	EndHeight     int64             `json:"end_height"`
	CurrentHeight int64             `json:"current_height"`
	StartedAt     time.Time         `json:"started_at"`
	CompletedAt   *time.Time        `json:"completed_at,omitempty"`
	Metadata      map[string]string `json:"metadata"`
}

// NetworkHealth represents network health metrics
type NetworkHealth struct {
	ID                string            `json:"id"`
	NetworkID         string            `json:"network_id"`
	Timestamp         time.Time         `json:"timestamp"`
	ActiveNodes       int               `json:"active_nodes"`
	AvgLatency        int64             `json:"avg_latency"`
	TotalConnections  int               `json:"total_connections"`
	ActiveConnections int               `json:"active_connections"`
	Throughput        int64             `json:"throughput"`   // bytes per second
	HealthScore       int64             `json:"health_score"` // 0 to 100 (percentage)
	Metadata          map[string]string `json:"metadata"`
}

// GenesisState defines the network module's genesis state
type GenesisState struct {
	Networks      []Network       `json:"networks"`
	Nodes         []Node          `json:"nodes"`
	Connections   []Connection    `json:"connections"`
	Syncs         []NetworkSync   `json:"syncs"`
	HealthMetrics []NetworkHealth `json:"health_metrics"`
	Params        Params          `json:"params"`
}

// Params defines the network module parameters
type Params struct {
	MaxNodesPerNetwork  int64 `json:"max_nodes_per_network"`
	MaxConnections      int64 `json:"max_connections"`
	SyncInterval        int64 `json:"sync_interval"`         // seconds
	HealthCheckInterval int64 `json:"health_check_interval"` // seconds
}

// DefaultParams returns default network module parameters
func DefaultParams() Params {
	return Params{
		MaxNodesPerNetwork:  1000,
		MaxConnections:      10000,
		SyncInterval:        60,
		HealthCheckInterval: 30,
	}
}

// Validate validates the network module parameters
func (p Params) Validate() error {
	if p.MaxNodesPerNetwork <= 0 {
		return fmt.Errorf("max_nodes_per_network must be positive")
	}
	if p.MaxConnections <= 0 {
		return fmt.Errorf("max_connections must be positive")
	}
	if p.SyncInterval <= 0 {
		return fmt.Errorf("sync_interval must be positive")
	}
	if p.HealthCheckInterval <= 0 {
		return fmt.Errorf("health_check_interval must be positive")
	}
	return nil
}

// Validate validates a network
func (n Network) Validate() error {
	if n.ID == "" {
		return fmt.Errorf("network ID cannot be empty")
	}
	if n.Name == "" {
		return fmt.Errorf("network name cannot be empty")
	}
	return nil
}

// Validate validates a node
func (n Node) Validate() error {
	if n.ID == "" {
		return fmt.Errorf("node ID cannot be empty")
	}
	if n.NetworkID == "" {
		return fmt.Errorf("network ID cannot be empty")
	}
	if n.Address == "" {
		return fmt.Errorf("node address cannot be empty")
	}
	if n.Port <= 0 || n.Port > 65535 {
		return fmt.Errorf("node port must be between 1 and 65535")
	}
	return nil
}

// Validate validates a connection
func (c Connection) Validate() error {
	if c.ID == "" {
		return fmt.Errorf("connection ID cannot be empty")
	}
	if c.NetworkID == "" {
		return fmt.Errorf("network ID cannot be empty")
	}
	if c.FromNodeID == "" {
		return fmt.Errorf("from node ID cannot be empty")
	}
	if c.ToNodeID == "" {
		return fmt.Errorf("to node ID cannot be empty")
	}
	return nil
}

// Validate validates network sync
func (ns NetworkSync) Validate() error {
	if ns.ID == "" {
		return fmt.Errorf("sync ID cannot be empty")
	}
	if ns.NetworkID == "" {
		return fmt.Errorf("network ID cannot be empty")
	}
	if ns.NodeID == "" {
		return fmt.Errorf("node ID cannot be empty")
	}
	if ns.Progress < 0 || ns.Progress > 100 {
		return fmt.Errorf("progress must be between 0 and 100")
	}
	return nil
}

// Validate validates network health
func (nh NetworkHealth) Validate() error {
	if nh.ID == "" {
		return fmt.Errorf("health ID cannot be empty")
	}
	if nh.NetworkID == "" {
		return fmt.Errorf("network ID cannot be empty")
	}
	if nh.HealthScore < 0 || nh.HealthScore > 100 {
		return fmt.Errorf("health score must be between 0 and 100")
	}
	return nil
}
