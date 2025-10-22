package service

import "time"

// timeToStringPtr time.Time を *string に変換
func timeToStringPtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format("2006-01-02")
	return &s
}

// timeToString time.Time を string に変換
func timeToString(t time.Time) string {
	return t.Format("2006-01-02")
}
