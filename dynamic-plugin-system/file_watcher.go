package main

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

// FileWatcher 文件监控器
type FileWatcher struct {
	watchDir string
	callback func(string)
	files    map[string]time.Time
	stop     chan bool
}

// NewFileWatcher 创建新的文件监控器
func NewFileWatcher(watchDir string, callback func(string)) *FileWatcher {
	return &FileWatcher{
		watchDir: watchDir,
		callback: callback,
		files:    make(map[string]time.Time),
		stop:     make(chan bool),
	}
}

// Start 开始监控
func (fw *FileWatcher) Start() {
	// 初始扫描
	fw.scanFiles()
	
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			fw.checkForChanges()
		case <-fw.stop:
			return
		}
	}
}

// Stop 停止监控
func (fw *FileWatcher) Stop() {
	close(fw.stop)
}

// scanFiles 扫描目录中的所有文件
func (fw *FileWatcher) scanFiles() {
	filepath.Walk(fw.watchDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		
		if !info.IsDir() && filepath.Ext(path) == ".go" {
			fw.files[path] = info.ModTime()
		}
		
		return nil
	})
}

// checkForChanges 检查文件变化
func (fw *FileWatcher) checkForChanges() {
	filepath.Walk(fw.watchDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		
		if !info.IsDir() && filepath.Ext(path) == ".go" {
			modTime := info.ModTime()
			
			if lastModTime, exists := fw.files[path]; exists {
				// 文件已存在，检查是否被修改
				if modTime.After(lastModTime) {
					fw.files[path] = modTime
					if fw.callback != nil {
						fw.callback(path)
					}
				}
			} else {
				// 新文件
				fw.files[path] = modTime
				if fw.callback != nil {
					fw.callback(path)
				}
			}
		}
		
		return nil
	})
	
	// 检查删除的文件
	for path := range fw.files {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			delete(fw.files, path)
		}
	}
}