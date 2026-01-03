package main

import (
	"encoding/json"
	"os"
	"time"
)

// #region agent log
func writeLog(location, message string, data map[string]interface{}, hypothesisId string) {
	logData := map[string]interface{}{
		"sessionId":    "debug-session",
		"runId":        "run1",
		"hypothesisId": hypothesisId,
		"location":     location,
		"message":      message,
		"data":         data,
		"timestamp":    time.Now().UnixMilli(),
	}
	if f, err := os.OpenFile("c:\\code\\clippytui\\.cursor\\debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
		json.NewEncoder(f).Encode(logData)
		f.Close()
	}
}

// #endregion

