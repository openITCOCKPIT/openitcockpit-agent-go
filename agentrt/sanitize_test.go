package agentrt

import (
	"encoding/json"
	"math"
	"testing"
)

func TestSanitizeFloats_NaN(t *testing.T) {
	input := map[string]interface{}{
		"valid":   42.5,
		"nan_val": math.NaN(),
		"inf_val": math.Inf(1),
		"neg_inf": math.Inf(-1),
	}

	result := sanitizeFloats(input)

	// Serialize and re-parse to verify JSON round-trip
	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("json.Marshal failed after sanitize: %v", err)
	}

	var parsed map[string]float64
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if parsed["valid"] != 42.5 {
		t.Errorf("expected 42.5, got %v", parsed["valid"])
	}
	if parsed["nan_val"] != 0 {
		t.Errorf("expected 0 for NaN, got %v", parsed["nan_val"])
	}
	if parsed["inf_val"] != 0 {
		t.Errorf("expected 0 for Inf, got %v", parsed["inf_val"])
	}
	if parsed["neg_inf"] != 0 {
		t.Errorf("expected 0 for -Inf, got %v", parsed["neg_inf"])
	}
}

func TestSanitizeFloats_NestedMap(t *testing.T) {
	input := map[string]interface{}{
		"cpu": map[string]interface{}{
			"total": math.NaN(),
			"cores": []interface{}{1.5, math.NaN(), 3.0, math.Inf(1)},
		},
		"name": "test",
	}

	result := sanitizeFloats(input)

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var parsed map[string]json.RawMessage
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	var cpu map[string]interface{}
	if err := json.Unmarshal(parsed["cpu"], &cpu); err != nil {
		t.Fatalf("json.Unmarshal cpu failed: %v", err)
	}

	if cpu["total"].(float64) != 0 {
		t.Errorf("expected 0 for nested NaN, got %v", cpu["total"])
	}
}

func TestSanitizeFloats_JsonMarshalSucceeds(t *testing.T) {
	input := map[string]interface{}{
		"value":  math.NaN(),
		"nested": map[string]interface{}{"x": math.Inf(-1)},
		"list":   []interface{}{math.NaN(), 1.0},
	}

	// Without sanitize, this would fail
	_, err := json.Marshal(input)
	if err == nil {
		t.Fatal("expected json.Marshal to fail on NaN without sanitization")
	}

	// With sanitize, it should succeed
	sanitized := sanitizeFloats(input)
	data, err := json.Marshal(sanitized)
	if err != nil {
		t.Fatalf("json.Marshal failed after sanitization: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}
}

func TestSanitizeFloats_Nil(t *testing.T) {
	if sanitizeFloats(nil) != nil {
		t.Error("expected nil for nil input")
	}
}

func TestSanitizeFloats_NoFloats(t *testing.T) {
	input := map[string]interface{}{
		"name":  "test",
		"count": 42,
		"tags":  []interface{}{"a", "b"},
	}

	result := sanitizeFloats(input)

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if parsed["name"].(string) != "test" {
		t.Errorf("expected test, got %v", parsed["name"])
	}
}

func TestSanitizeFloats_PreservesInt(t *testing.T) {
	input := map[string]interface{}{
		"pid":     int64(728),
		"cpu":     float64(42.5),
		"nan":     math.NaN(),
		"name":    "test",
		"running": true,
	}

	result := sanitizeFloats(input)

	// Must not panic on json.Marshal
	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	// JSON numbers are float64 after unmarshal
	if parsed["pid"].(float64) != 728 {
		t.Errorf("expected 728, got %v", parsed["pid"])
	}
	if parsed["cpu"].(float64) != 42.5 {
		t.Errorf("expected 42.5, got %v", parsed["cpu"])
	}
	if parsed["nan"].(float64) != 0 {
		t.Errorf("expected 0 for NaN, got %v", parsed["nan"])
	}
}

func TestSanitizeFloats_StructWithNaN(t *testing.T) {
	// Simulates a check result struct with NaN values (like resultDiskIo)
	type checkResult struct {
		Device      string  `json:"device"`
		LoadPercent float64 `json:"load_percent"`
		ReadWait    float64 `json:"read_wait"`
		WriteWait   float64 `json:"write_wait"`
	}

	input := map[string]interface{}{
		"disk_io": map[string]*checkResult{
			"C:": {
				Device:      "C:",
				LoadPercent: math.NaN(),
				ReadWait:    1.5,
				WriteWait:   math.Inf(1),
			},
		},
		"cpu": map[string]interface{}{
			"total": 42.5,
		},
	}

	// Without sanitize, this fails
	_, err := json.Marshal(input)
	if err == nil {
		t.Fatal("expected json.Marshal to fail on struct with NaN")
	}

	// With sanitize, should succeed and repair NaN values
	result := sanitizeFloats(input)
	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("json.Marshal failed after sanitization: %v", err)
	}

	var parsed map[string]json.RawMessage
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	// disk_io should be present and repaired
	var diskIo map[string]checkResult
	if err := json.Unmarshal(parsed["disk_io"], &diskIo); err != nil {
		t.Fatalf("json.Unmarshal disk_io failed: %v", err)
	}

	if diskIo["C:"].LoadPercent != 0 {
		t.Errorf("expected 0 for NaN LoadPercent, got %v", diskIo["C:"].LoadPercent)
	}
	if diskIo["C:"].ReadWait != 1.5 {
		t.Errorf("expected 1.5 for ReadWait, got %v", diskIo["C:"].ReadWait)
	}
	if diskIo["C:"].WriteWait != 0 {
		t.Errorf("expected 0 for Inf WriteWait, got %v", diskIo["C:"].WriteWait)
	}

	// cpu should be unaffected
	if _, exists := parsed["cpu"]; !exists {
		t.Error("cpu key should still exist")
	}
}

func TestSanitizeFloats_DeeplyNested(t *testing.T) {
	input := map[string]interface{}{
		"level1": map[string]interface{}{
			"level2": map[string]interface{}{
				"level3": []interface{}{
					map[string]interface{}{
						"value": math.NaN(),
						"ok":    42.0,
					},
				},
			},
		},
	}

	result := sanitizeFloats(input)
	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("json.Marshal failed on deeply nested: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}
}

func TestSanitizeFloats_PassthroughNonMap(t *testing.T) {
	// Non-map input should pass through unchanged
	result := sanitizeFloats("hello")
	if result.(string) != "hello" {
		t.Errorf("expected passthrough, got %v", result)
	}

	result = sanitizeFloats(42)
	if result.(int) != 42 {
		t.Errorf("expected passthrough, got %v", result)
	}
}
