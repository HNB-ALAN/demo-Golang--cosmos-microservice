package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/usc-platform/shared/config"
	"github.com/usc-platform/shared/database"
	"github.com/usc-platform/shared/logging"
)

// Migration represents a database migration
type Migration struct {
	Version     int
	Name        string
	UpSQL       string
	DownSQL     string
	Description string
	CreatedAt   time.Time
}

// MigrationStatus represents the status of a migration
type MigrationStatus struct {
	Version   int       `json:"version"`
	Name      string    `json:"name"`
	Applied   bool      `json:"applied"`
	AppliedAt time.Time `json:"applied_at,omitempty"`
}

// MigrationManager manages database migrations for USC Blockchain Core Service
type MigrationManager struct {
	dbManager      *database.DatabaseManager
	config         *config.Config
	logger         logging.Logger
	migrationsPath string
}

// NewMigrationManager creates a new migration manager
func NewMigrationManager(dbManager *database.DatabaseManager, cfg *config.Config, logger logging.Logger) (*MigrationManager, error) {
	return &MigrationManager{
		dbManager:      dbManager,
		config:         cfg,
		logger:         logger,
		migrationsPath: filepath.Join("migrations", "postgresql"),
	}, nil
}

// CreateMigrationsTable creates the migrations tracking table
func (mm *MigrationManager) CreateMigrationsTable(ctx context.Context) error {
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			description TEXT
		);
	`

	db := mm.dbManager.PostgreSQL()
	if db == nil {
		return fmt.Errorf("PostgreSQL connection not available")
	}

	_, err := db.ExecContext(ctx, createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	mm.logger.Info("Migrations table created successfully")
	return nil
}

// GetAppliedMigrations returns a list of applied migrations
func (mm *MigrationManager) GetAppliedMigrations(ctx context.Context) ([]MigrationStatus, error) {
	db := mm.dbManager.PostgreSQL()
	if db == nil {
		return nil, fmt.Errorf("PostgreSQL connection not available")
	}

	query := `
		SELECT version, name, applied_at, description 
		FROM schema_migrations 
		ORDER BY version ASC
	`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query applied migrations: %w", err)
	}
	defer rows.Close()

	var appliedMigrations []MigrationStatus
	for rows.Next() {
		var status MigrationStatus
		var description sql.NullString

		err := rows.Scan(&status.Version, &status.Name, &status.AppliedAt, &description)
		if err != nil {
			return nil, fmt.Errorf("failed to scan migration row: %w", err)
		}

		status.Applied = true
		appliedMigrations = append(appliedMigrations, status)
	}

	return appliedMigrations, nil
}

// GetPendingMigrations returns a list of pending migrations
func (mm *MigrationManager) GetPendingMigrations(ctx context.Context) ([]Migration, error) {
	mm.logger.Info("Getting pending migrations", logging.String("path", mm.migrationsPath))

	// Read all migration files
	allMigrations, err := mm.readMigrationFiles()
	if err != nil {
		return nil, fmt.Errorf("failed to read migration files: %w", err)
	}

	// Get applied migrations
	appliedMigrations, err := mm.GetAppliedMigrations(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Create lookup map for applied migrations
	appliedVersions := make(map[int]bool)
	for _, applied := range appliedMigrations {
		appliedVersions[applied.Version] = true
	}

	// Filter out applied migrations
	var pendingMigrations []Migration
	for _, migration := range allMigrations {
		if !appliedVersions[migration.Version] {
			pendingMigrations = append(pendingMigrations, migration)
		}
	}

	// Sort by version
	sort.Slice(pendingMigrations, func(i, j int) bool {
		return pendingMigrations[i].Version < pendingMigrations[j].Version
	})

	mm.logger.Info("Found pending migrations", logging.Int("count", len(pendingMigrations)))
	return pendingMigrations, nil
}

// MigrateUp applies all pending migrations
func (mm *MigrationManager) MigrateUp(ctx context.Context) error {
	// Ensure migrations table exists
	if err := mm.CreateMigrationsTable(ctx); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get pending migrations
	pendingMigrations, err := mm.GetPendingMigrations(ctx)
	if err != nil {
		return fmt.Errorf("failed to get pending migrations: %w", err)
	}

	if len(pendingMigrations) == 0 {
		mm.logger.Info("No pending migrations to apply")
		return nil
	}

	mm.logger.Info("Applying pending migrations", logging.Int("count", len(pendingMigrations)))

	// Apply migrations
	for _, migration := range pendingMigrations {
		mm.logger.Info("Applying migration",
			logging.Int("version", migration.Version),
			logging.String("name", migration.Name))

		if err := mm.ApplyMigration(ctx, migration); err != nil {
			return fmt.Errorf("failed to apply migration %d (%s): %w",
				migration.Version, migration.Name, err)
		}

		mm.logger.Info("Migration applied successfully",
			logging.Int("version", migration.Version),
			logging.String("name", migration.Name))
	}

	mm.logger.Info("All pending migrations applied successfully",
		logging.Int("count", len(pendingMigrations)))

	return nil
}

// MigrateDown rolls back the last applied migration
func (mm *MigrationManager) MigrateDown(ctx context.Context) error {
	// Get applied migrations
	appliedMigrations, err := mm.GetAppliedMigrations(ctx)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	if len(appliedMigrations) == 0 {
		mm.logger.Info("No migrations to rollback")
		return nil
	}

	// Find the highest version number (last applied migration)
	highestMigration := appliedMigrations[len(appliedMigrations)-1]

	// Get all migration files to find the down SQL
	allMigrations, err := mm.readMigrationFiles()
	if err != nil {
		return fmt.Errorf("failed to read migration files: %w", err)
	}

	// Find the migration with down SQL
	var rollbackMigration Migration
	for _, migration := range allMigrations {
		if migration.Version == highestMigration.Version {
			rollbackMigration = migration
			break
		}
	}

	if rollbackMigration.Version == 0 {
		return fmt.Errorf("migration %d not found in files", highestMigration.Version)
	}

	// Rollback the migration
	if err := mm.RollbackMigration(ctx, rollbackMigration); err != nil {
		return fmt.Errorf("failed to rollback migration %d: %w", highestMigration.Version, err)
	}

	mm.logger.Info("Migration rolled back successfully",
		logging.Int("version", highestMigration.Version),
		logging.String("name", highestMigration.Name))
	return nil
}

// GetMigrationStatus returns the current migration status
func (mm *MigrationManager) GetMigrationStatus(ctx context.Context) ([]MigrationStatus, error) {
	// Get applied migrations
	appliedMigrations, err := mm.GetAppliedMigrations(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Get all migration files
	allMigrations, err := mm.readMigrationFiles()
	if err != nil {
		return nil, fmt.Errorf("failed to read migration files: %w", err)
	}

	// Create lookup map for applied migrations
	appliedMap := make(map[int]MigrationStatus)
	for _, applied := range appliedMigrations {
		appliedMap[applied.Version] = applied
	}

	// Create status list
	var statusList []MigrationStatus
	for _, migration := range allMigrations {
		status := MigrationStatus{
			Version: migration.Version,
			Name:    migration.Name,
			Applied: false,
		}

		if appliedStatus, exists := appliedMap[migration.Version]; exists {
			status.Applied = true
			status.AppliedAt = appliedStatus.AppliedAt
		}

		statusList = append(statusList, status)
	}

	// Sort by version
	sort.Slice(statusList, func(i, j int) bool {
		return statusList[i].Version < statusList[j].Version
	})

	return statusList, nil
}

// ValidateMigrations validates that all migrations are properly formatted
func (mm *MigrationManager) ValidateMigrations(ctx context.Context) error {
	// Get all migration files
	allMigrations, err := mm.readMigrationFiles()
	if err != nil {
		return fmt.Errorf("failed to read migration files: %w", err)
	}

	if len(allMigrations) == 0 {
		mm.logger.Info("No migration files found to validate")
		return nil
	}

	// Validate each migration
	for _, migration := range allMigrations {
		if err := mm.validateMigration(migration); err != nil {
			return err
		}
	}

	// Validate version sequence
	if err := mm.validateVersionSequence(allMigrations); err != nil {
		return err
	}

	mm.logger.Info("All migrations validated successfully",
		logging.Int("count", len(allMigrations)))

	return nil
}

// ApplyMigration applies a single migration
func (mm *MigrationManager) ApplyMigration(ctx context.Context, migration Migration) error {
	db := mm.dbManager.PostgreSQL()
	if db == nil {
		return fmt.Errorf("PostgreSQL connection not available")
	}

	// Begin transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Execute migration SQL
	_, err = tx.ExecContext(ctx, migration.UpSQL)
	if err != nil {
		return fmt.Errorf("failed to execute migration %d: %w", migration.Version, err)
	}

	// Record migration in tracking table
	insertSQL := `
		INSERT INTO schema_migrations (version, name, description) 
		VALUES ($1, $2, $3)
	`
	_, err = tx.ExecContext(ctx, insertSQL, migration.Version, migration.Name, migration.Description)
	if err != nil {
		return fmt.Errorf("failed to record migration %d: %w", migration.Version, err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit migration: %w", err)
	}

	mm.logger.Info("Migration applied successfully",
		logging.Int("version", migration.Version),
		logging.String("name", migration.Name))

	return nil
}

// RollbackMigration rolls back a single migration
func (mm *MigrationManager) RollbackMigration(ctx context.Context, migration Migration) error {
	db := mm.dbManager.PostgreSQL()
	if db == nil {
		return fmt.Errorf("PostgreSQL connection not available")
	}

	// Begin transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Execute rollback SQL
	_, err = tx.ExecContext(ctx, migration.DownSQL)
	if err != nil {
		return fmt.Errorf("failed to rollback migration %d: %w", migration.Version, err)
	}

	// Remove migration from tracking table
	deleteSQL := `DELETE FROM schema_migrations WHERE version = $1`
	_, err = tx.ExecContext(ctx, deleteSQL, migration.Version)
	if err != nil {
		return fmt.Errorf("failed to remove migration record %d: %w", migration.Version, err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit rollback: %w", err)
	}

	mm.logger.Info("Migration rolled back successfully",
		logging.Int("version", migration.Version),
		logging.String("name", migration.Name))

	return nil
}

// readMigrationFiles reads all migration files from the filesystem
func (mm *MigrationManager) readMigrationFiles() ([]Migration, error) {
	// Check if migrations directory exists
	if _, err := os.Stat(mm.migrationsPath); os.IsNotExist(err) {
		mm.logger.Warn("Migrations directory does not exist", logging.String("path", mm.migrationsPath))
		return []Migration{}, nil
	}

	// Read all files in the migrations directory
	files, err := os.ReadDir(mm.migrationsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Parse migration files
	migrationMap := make(map[int]*Migration)
	upFileRegex := regexp.MustCompile(`^(\d+)_(.+)\.up\.sql$`)
	downFileRegex := regexp.MustCompile(`^(\d+)_(.+)\.down\.sql$`)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filename := file.Name()

		// Parse up files
		if matches := upFileRegex.FindStringSubmatch(filename); matches != nil {
			version, err := strconv.Atoi(matches[1])
			if err != nil {
				mm.logger.Warn("Invalid migration version", logging.String("file", filename))
				continue
			}

			name := matches[2]
			upSQL, err := mm.readMigrationFile(filepath.Join(mm.migrationsPath, filename))
			if err != nil {
				mm.logger.Warn("Failed to read up migration file",
					logging.String("file", filename),
					logging.Error(err))
				continue
			}

			// Create or update migration
			if migration, exists := migrationMap[version]; exists {
				migration.UpSQL = upSQL
			} else {
				info, err := file.Info()
				if err != nil {
					mm.logger.Warn("Failed to get file info", logging.String("file", filename), logging.Error(err))
					continue
				}
				migrationMap[version] = &Migration{
					Version:     version,
					Name:        name,
					UpSQL:       upSQL,
					Description: fmt.Sprintf("Migration %d: %s", version, name),
					CreatedAt:   info.ModTime(),
				}
			}
		}

		// Parse down files
		if matches := downFileRegex.FindStringSubmatch(filename); matches != nil {
			version, err := strconv.Atoi(matches[1])
			if err != nil {
				mm.logger.Warn("Invalid migration version", logging.String("file", filename))
				continue
			}

			downSQL, err := mm.readMigrationFile(filepath.Join(mm.migrationsPath, filename))
			if err != nil {
				mm.logger.Warn("Failed to read down migration file",
					logging.String("file", filename),
					logging.Error(err))
				continue
			}

			// Create or update migration
			if migration, exists := migrationMap[version]; exists {
				migration.DownSQL = downSQL
			} else {
				info, err := file.Info()
				if err != nil {
					mm.logger.Warn("Failed to get file info", logging.String("file", filename), logging.Error(err))
					continue
				}
				migrationMap[version] = &Migration{
					Version:     version,
					Name:        matches[2],
					DownSQL:     downSQL,
					Description: fmt.Sprintf("Migration %d: %s", version, matches[2]),
					CreatedAt:   info.ModTime(),
				}
			}
		}
	}

	// Convert map to slice and sort by version
	var migrations []Migration
	for _, migration := range migrationMap {
		migrations = append(migrations, *migration)
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	mm.logger.Info("Loaded migration files",
		logging.Int("count", len(migrations)),
		logging.String("path", mm.migrationsPath))

	return migrations, nil
}

// readMigrationFile reads the content of a migration file
func (mm *MigrationManager) readMigrationFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	return strings.TrimSpace(string(content)), nil
}

// validateMigration validates a single migration
func (mm *MigrationManager) validateMigration(migration Migration) error {
	// Validate migration structure
	if migration.Version <= 0 {
		return fmt.Errorf("invalid migration version: %d", migration.Version)
	}

	if strings.TrimSpace(migration.UpSQL) == "" {
		return fmt.Errorf("empty UpSQL for migration %d (%s)", migration.Version, migration.Name)
	}

	if strings.TrimSpace(migration.DownSQL) == "" {
		return fmt.Errorf("empty DownSQL for migration %d (%s)", migration.Version, migration.Name)
	}

	// Validate SQL syntax (basic check)
	upSQLUpper := strings.ToUpper(migration.UpSQL)
	if !strings.Contains(upSQLUpper, "CREATE") &&
		!strings.Contains(upSQLUpper, "ALTER") &&
		!strings.Contains(upSQLUpper, "INSERT") {
		mm.logger.Warn("Migration UpSQL may not contain DDL/DML statements",
			logging.Int("version", migration.Version),
			logging.String("name", migration.Name))
	}

	downSQLUpper := strings.ToUpper(migration.DownSQL)
	if !strings.Contains(downSQLUpper, "DROP") &&
		!strings.Contains(downSQLUpper, "ALTER") &&
		!strings.Contains(downSQLUpper, "DELETE") {
		mm.logger.Warn("Migration DownSQL may not contain rollback statements",
			logging.Int("version", migration.Version),
			logging.String("name", migration.Name))
	}

	return nil
}

// validateVersionSequence validates that migration versions are sequential
func (mm *MigrationManager) validateVersionSequence(migrations []Migration) error {
	for i := 1; i < len(migrations); i++ {
		if migrations[i].Version != migrations[i-1].Version+1 {
			mm.logger.Warn("Non-sequential migration versions detected",
				logging.Int("expected", migrations[i-1].Version+1),
				logging.Int("found", migrations[i].Version))
		}
	}
	return nil
}
