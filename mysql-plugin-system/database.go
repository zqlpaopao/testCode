package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// DatabaseManager 数据库管理器
type DatabaseManager struct {
	db *sql.DB
}

// PluginRecord 插件记录结构
type PluginRecord struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Version     string    `json:"version"`
	Description string    `json:"description"`
	SourceCode  string    `json:"source_code"`
	BinaryData  []byte    `json:"binary_data"`
	Hash        string    `json:"hash"`
	Status      string    `json:"status"` // active, inactive, error
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewDatabaseManager 创建新的数据库管理器
func NewDatabaseManager(dsn string) (*DatabaseManager, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %v", err)
	}

	// 测试连接
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("数据库连接测试失败: %v", err)
	}

	dm := &DatabaseManager{db: db}
	
	// 初始化数据库表
	if err := dm.initTables(); err != nil {
		return nil, fmt.Errorf("初始化数据库表失败: %v", err)
	}

	return dm, nil
}

// initTables 初始化数据库表
func (dm *DatabaseManager) initTables() error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS plugins (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL UNIQUE,
		version VARCHAR(50) NOT NULL DEFAULT '1.0.0',
		description TEXT,
		source_code LONGTEXT NOT NULL,
		binary_data LONGBLOB,
		hash VARCHAR(64) NOT NULL,
		status ENUM('active', 'inactive', 'error') DEFAULT 'active',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		INDEX idx_name (name),
		INDEX idx_status (status),
		INDEX idx_hash (hash)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`

	if _, err := dm.db.Exec(createTableSQL); err != nil {
		return fmt.Errorf("创建插件表失败: %v", err)
	}

	// 创建插件变更日志表
	createLogTableSQL := `
	CREATE TABLE IF NOT EXISTS plugin_changes (
		id INT AUTO_INCREMENT PRIMARY KEY,
		plugin_name VARCHAR(255) NOT NULL,
		action ENUM('insert', 'update', 'delete') NOT NULL,
		old_hash VARCHAR(64),
		new_hash VARCHAR(64),
		changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		INDEX idx_plugin_name (plugin_name),
		INDEX idx_changed_at (changed_at)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`

	if _, err := dm.db.Exec(createLogTableSQL); err != nil {
		return fmt.Errorf("创建插件变更日志表失败: %v", err)
	}

	return nil
}

// SavePlugin 保存插件到数据库
func (dm *DatabaseManager) SavePlugin(plugin *PluginRecord) error {
	// 检查插件是否已存在
	var existingID int
	err := dm.db.QueryRow("SELECT id FROM plugins WHERE name = ?", plugin.Name).Scan(&existingID)
	
	if err == sql.ErrNoRows {
		// 插入新插件
		insertSQL := `
		INSERT INTO plugins (name, version, description, source_code, binary_data, hash, status)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		`
		result, err := dm.db.Exec(insertSQL, 
			plugin.Name, plugin.Version, plugin.Description, 
			plugin.SourceCode, plugin.BinaryData, plugin.Hash, plugin.Status)
		
		if err != nil {
			return fmt.Errorf("插入插件失败: %v", err)
		}

		id, _ := result.LastInsertId()
		plugin.ID = int(id)

		// 记录变更日志
		dm.logPluginChange(plugin.Name, "insert", "", plugin.Hash)
		
		log.Printf("✅ 插件 '%s' 已保存到数据库 (ID: %d)", plugin.Name, plugin.ID)
		
	} else if err != nil {
		return fmt.Errorf("查询插件失败: %v", err)
	} else {
		// 更新现有插件
		var oldHash string
		dm.db.QueryRow("SELECT hash FROM plugins WHERE id = ?", existingID).Scan(&oldHash)
		
		updateSQL := `
		UPDATE plugins 
		SET version = ?, description = ?, source_code = ?, binary_data = ?, hash = ?, status = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
		`
		_, err := dm.db.Exec(updateSQL,
			plugin.Version, plugin.Description, plugin.SourceCode, 
			plugin.BinaryData, plugin.Hash, plugin.Status, existingID)
		
		if err != nil {
			return fmt.Errorf("更新插件失败: %v", err)
		}

		plugin.ID = existingID

		// 记录变更日志
		dm.logPluginChange(plugin.Name, "update", oldHash, plugin.Hash)
		
		log.Printf("🔄 插件 '%s' 已更新到数据库 (ID: %d)", plugin.Name, plugin.ID)
	}

	return nil
}

// GetPlugin 从数据库获取插件
func (dm *DatabaseManager) GetPlugin(name string) (*PluginRecord, error) {
	plugin := &PluginRecord{}
	
	query := `
	SELECT id, name, version, description, source_code, binary_data, hash, status, created_at, updated_at
	FROM plugins 
	WHERE name = ? AND status = 'active'
	`
	
	err := dm.db.QueryRow(query, name).Scan(
		&plugin.ID, &plugin.Name, &plugin.Version, &plugin.Description,
		&plugin.SourceCode, &plugin.BinaryData, &plugin.Hash, &plugin.Status,
		&plugin.CreatedAt, &plugin.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("插件 '%s' 不存在", name)
	} else if err != nil {
		return nil, fmt.Errorf("查询插件失败: %v", err)
	}
	
	return plugin, nil
}

// GetAllPlugins 获取所有活跃的插件
func (dm *DatabaseManager) GetAllPlugins() ([]*PluginRecord, error) {
	query := `
	SELECT id, name, version, description, source_code, binary_data, hash, status, created_at, updated_at
	FROM plugins 
	WHERE status = 'active'
	ORDER BY name
	`
	
	rows, err := dm.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("查询所有插件失败: %v", err)
	}
	defer rows.Close()
	
	var plugins []*PluginRecord
	
	for rows.Next() {
		plugin := &PluginRecord{}
		err := rows.Scan(
			&plugin.ID, &plugin.Name, &plugin.Version, &plugin.Description,
			&plugin.SourceCode, &plugin.BinaryData, &plugin.Hash, &plugin.Status,
			&plugin.CreatedAt, &plugin.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("扫描插件记录失败: %v", err)
		}
		plugins = append(plugins, plugin)
	}
	
	return plugins, nil
}

// DeletePlugin 删除插件
func (dm *DatabaseManager) DeletePlugin(name string) error {
	// 获取旧哈希值
	var oldHash string
	dm.db.QueryRow("SELECT hash FROM plugins WHERE name = ?", name).Scan(&oldHash)
	
	_, err := dm.db.Exec("UPDATE plugins SET status = 'inactive' WHERE name = ?", name)
	if err != nil {
		return fmt.Errorf("删除插件失败: %v", err)
	}

	// 记录变更日志
	dm.logPluginChange(name, "delete", oldHash, "")
	
	log.Printf("🗑️ 插件 '%s' 已标记为非活跃状态", name)
	return nil
}

// GetPluginChanges 获取插件变更记录
func (dm *DatabaseManager) GetPluginChanges(since time.Time) ([]PluginChange, error) {
	query := `
	SELECT plugin_name, action, old_hash, new_hash, changed_at
	FROM plugin_changes 
	WHERE changed_at > ?
	ORDER BY changed_at DESC
	`
	
	rows, err := dm.db.Query(query, since)
	if err != nil {
		return nil, fmt.Errorf("查询插件变更失败: %v", err)
	}
	defer rows.Close()
	
	var changes []PluginChange
	
	for rows.Next() {
		change := PluginChange{}
		err := rows.Scan(
			&change.PluginName, &change.Action, &change.OldHash, &change.NewHash, &change.ChangedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("扫描变更记录失败: %v", err)
		}
		changes = append(changes, change)
	}
	
	return changes, nil
}

// logPluginChange 记录插件变更日志
func (dm *DatabaseManager) logPluginChange(pluginName, action, oldHash, newHash string) {
	insertSQL := `
	INSERT INTO plugin_changes (plugin_name, action, old_hash, new_hash)
	VALUES (?, ?, ?, ?)
	`
	
	_, err := dm.db.Exec(insertSQL, pluginName, action, 
		sql.NullString{String: oldHash, Valid: oldHash != ""}, 
		sql.NullString{String: newHash, Valid: newHash != ""})
	
	if err != nil {
		log.Printf("❌ 记录插件变更日志失败: %v", err)
	}
}

// GetPluginHash 获取插件哈希值
func (dm *DatabaseManager) GetPluginHash(name string) (string, error) {
	var hash string
	err := dm.db.QueryRow("SELECT hash FROM plugins WHERE name = ? AND status = 'active'", name).Scan(&hash)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("插件 '%s' 不存在", name)
	} else if err != nil {
		return "", fmt.Errorf("查询插件哈希失败: %v", err)
	}
	return hash, nil
}

// Close 关闭数据库连接
func (dm *DatabaseManager) Close() error {
	if dm.db != nil {
		return dm.db.Close()
	}
	return nil
}

// PluginChange 插件变更记录
type PluginChange struct {
	PluginName string    `json:"plugin_name"`
	Action     string    `json:"action"`
	OldHash    string    `json:"old_hash"`
	NewHash    string    `json:"new_hash"`
	ChangedAt  time.Time `json:"changed_at"`
}

// GetPluginStats 获取插件统计信息
func (dm *DatabaseManager) GetPluginStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})
	
	// 总插件数
	var totalCount int
	err := dm.db.QueryRow("SELECT COUNT(*) FROM plugins WHERE status = 'active'").Scan(&totalCount)
	if err != nil {
		return nil, err
	}
	stats["total_plugins"] = totalCount
	
	// 按状态统计
	statusQuery := `
	SELECT status, COUNT(*) 
	FROM plugins 
	GROUP BY status
	`
	rows, err := dm.db.Query(statusQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	statusStats := make(map[string]int)
	for rows.Next() {
		var status string
		var count int
		rows.Scan(&status, &count)
		statusStats[status] = count
	}
	stats["by_status"] = statusStats
	
	// 最近变更数量
	var recentChanges int
	since := time.Now().Add(-24 * time.Hour)
	err = dm.db.QueryRow("SELECT COUNT(*) FROM plugin_changes WHERE changed_at > ?", since).Scan(&recentChanges)
	if err != nil {
		return nil, err
	}
	stats["recent_changes_24h"] = recentChanges
	
	return stats, nil
}