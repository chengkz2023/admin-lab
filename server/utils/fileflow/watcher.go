package fileflow

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type fileSnapshot struct {
	size    int64
	modTime time.Time
	// stable=true 代表该文件已投递过事件，避免重复触发。
	stable bool
	// firstSeen 用于记录首次发现时间，写入到 FileEvent.CreatedAt。
	firstSeen time.Time
}

// FileWatcher 通过轮询目录发现“稳定文件”并产生事件。
type FileWatcher struct {
	dir      string
	interval time.Duration
	filter   func(name string) bool
	seen     map[string]fileSnapshot
	mu       sync.Mutex
	out      chan FileEvent
}

// newFileWatcher 创建轮询 watcher。
func newFileWatcher(cfg Config) *FileWatcher {
	return &FileWatcher{
		dir:      cfg.WatchDir,
		interval: cfg.Interval,
		filter:   cfg.Filter,
		seen:     make(map[string]fileSnapshot),
		out:      make(chan FileEvent, cfg.EventBuffer),
	}
}

// Watch 启动轮询协程并返回事件通道，ctx 取消后自动关闭通道。
func (w *FileWatcher) Watch(ctx context.Context) <-chan FileEvent {
	go func() {
		t := time.NewTicker(w.interval)
		defer t.Stop()
		defer close(w.out)
		w.scan()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				w.scan()
			}
		}
	}()
	return w.out
}

// scan 单次扫描目录：
// 1) 首次发现先记录快照；
// 2) 连续两次 size+mtime 不变才判定为稳定；
// 3) 事件通道满时不标记 stable，下轮继续尝试投递，避免静默丢事件。
func (w *FileWatcher) scan() {
	entries, err := os.ReadDir(w.dir)
	if err != nil {
		return
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	now := time.Now()
	current := make(map[string]struct{}, len(entries))

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if w.filter != nil && !w.filter(e.Name()) {
			continue
		}

		path := filepath.Join(w.dir, e.Name())
		info, err := e.Info()
		if err != nil {
			continue
		}
		current[path] = struct{}{}

		prev, known := w.seen[path]
		snap := fileSnapshot{
			size:      info.Size(),
			modTime:   info.ModTime(),
			firstSeen: now,
		}

		if known {
			snap.firstSeen = prev.firstSeen
		}

		if !known {
			w.seen[path] = snap
			continue
		}

		if prev.stable {
			continue
		}

		if prev.size == snap.size && prev.modTime.Equal(snap.modTime) {
			event := FileEvent{
				Path:      path,
				Name:      e.Name(),
				Size:      info.Size(),
				ModTime:   info.ModTime(),
				CreatedAt: snap.firstSeen,
			}
			select {
			case w.out <- event:
				snap.stable = true
			default:
				// Keep stable=false to retry next scan instead of silently dropping.
			}
		}

		w.seen[path] = snap
	}

	for p := range w.seen {
		if _, ok := current[p]; !ok {
			delete(w.seen, p)
		}
	}
}
