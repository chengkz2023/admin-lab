package fileflow

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type fileSnapshot struct {
	firstSeen time.Time
	emitted   bool
}

// FileWatcher 轮询目录，并投递通过 Filter 与 ReadyStrategy 的文件。
type FileWatcher struct {
	dir           string
	interval      time.Duration
	filter        func(name string) bool
	readyStrategy ReadyStrategy
	logger        *slog.Logger
	seen          map[string]fileSnapshot
	mu            sync.Mutex
	out           chan FileEvent
}

func newFileWatcher(cfg Config) *FileWatcher {
	return &FileWatcher{
		dir:           cfg.WatchDir,
		interval:      cfg.Interval,
		filter:        cfg.Filter,
		readyStrategy: cfg.ReadyStrategy,
		logger:        cfg.Logger,
		seen:          make(map[string]fileSnapshot),
		out:           make(chan FileEvent, cfg.EventBuffer),
	}
}

// Watch 启动轮询，并在 ctx 取消后关闭事件通道。
func (w *FileWatcher) Watch(ctx context.Context) <-chan FileEvent {
	go func() {
		t := time.NewTicker(w.interval)
		defer t.Stop()
		defer close(w.out)
		w.scan(ctx)
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				w.scan(ctx)
			}
		}
	}()
	return w.out
}

func (w *FileWatcher) scan(ctx context.Context) {
	entries, err := os.ReadDir(w.dir)
	if err != nil {
		if w.logger != nil {
			w.logger.Warn("fileflow: read watch dir failed", "dir", w.dir, "err", err)
		}
		return
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	now := time.Now()
	current := make(map[string]struct{}, len(entries))

	for _, entry := range entries {
		if ctx.Err() != nil {
			return
		}
		if entry.IsDir() {
			continue
		}
		if w.filter != nil && !w.filter(entry.Name()) {
			continue
		}

		path := filepath.Join(w.dir, entry.Name())
		info, err := entry.Info()
		if err != nil {
			continue
		}
		current[path] = struct{}{}

		snap, known := w.seen[path]
		if !known {
			snap.firstSeen = now
			w.seen[path] = snap
		}
		if snap.emitted {
			continue
		}

		ready, err := w.readyStrategy.Ready(ctx, path, info)
		if err != nil {
			if w.logger != nil {
				w.logger.Warn("fileflow: ready check failed", "file", entry.Name(), "err", err)
			}
			continue
		}
		if !ready {
			continue
		}

		event := FileEvent{
			Path:      path,
			Name:      entry.Name(),
			Size:      info.Size(),
			ModTime:   info.ModTime(),
			CreatedAt: snap.firstSeen,
		}
		select {
		case w.out <- event:
			snap.emitted = true
			w.seen[path] = snap
		default:
			// 保持 emitted=false，让下一轮扫描重试，避免丢事件。
		}
	}

	for path := range w.seen {
		if _, ok := current[path]; !ok {
			delete(w.seen, path)
		}
	}
}
