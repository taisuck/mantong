// Package logging 은 로그 레벨 파싱과 KST 날짜별 로그 파일 출력을 제공한다.
package logging

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var kst = time.FixedZone("KST", 9*60*60)

// ParseLevel 은 "debug", "info", "warn", "error" 문자열을 slog.Level 로 변환한다.
func ParseLevel(s string) (slog.Level, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "debug":
		return slog.LevelDebug, nil
	case "", "info":
		return slog.LevelInfo, nil
	case "warn", "warning":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return slog.LevelInfo, fmt.Errorf("알 수 없는 로그 레벨: %q", s)
	}
}

// Setup 은 stdout 과 dir 아래 KST 날짜별 로그 파일(mantong-YYYY-MM-DD.log)에 함께 기록하는
// 로거를 기본 로거로 설정한다. 반환된 io.Closer 로 로그 파일을 닫는다.
func Setup(dir, level string) (io.Closer, error) {
	lv, err := ParseLevel(level)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("로그 디렉터리 생성 실패: %w", err)
	}

	fw := &dailyFileWriter{dir: dir, prefix: "mantong"}
	handler := slog.NewTextHandler(io.MultiWriter(os.Stdout, fw), &slog.HandlerOptions{
		Level: lv,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey && a.Value.Kind() == slog.KindTime {
				a.Value = slog.StringValue(a.Value.Time().In(kst).Format("2006-01-02T15:04:05.000+09:00"))
			}
			return a
		},
	})
	slog.SetDefault(slog.New(handler))
	return fw, nil
}

// dailyFileWriter 는 KST 날짜가 바뀌면 새 로그 파일로 전환한다.
type dailyFileWriter struct {
	mu     sync.Mutex
	dir    string
	prefix string
	date   string
	file   *os.File
}

func (w *dailyFileWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	date := time.Now().In(kst).Format("2006-01-02")
	if w.file == nil || date != w.date {
		if w.file != nil {
			w.file.Close()
		}
		name := filepath.Join(w.dir, fmt.Sprintf("%s-%s.log", w.prefix, date))
		f, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
		if err != nil {
			w.file = nil
			return 0, err
		}
		w.file = f
		w.date = date
	}
	return w.file.Write(p)
}

func (w *dailyFileWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.file == nil {
		return nil
	}
	err := w.file.Close()
	w.file = nil
	return err
}
