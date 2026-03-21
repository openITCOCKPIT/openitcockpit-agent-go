package agentrt

import (
	"encoding/json"
	"math"
	"reflect"
	"strings"

	log "github.com/sirupsen/logrus"
)

// sanitizeFloats prevents json.Marshal from failing on NaN/Inf float values.
//
// Strategy:
// 1. Try json.Marshal on the entire result — if it works, no NaN present (fast path).
// 2. If it fails, marshal each check result individually.
// 3. For checks that fail, repair NaN/Inf values in-place using reflect,
//    then marshal again.
//
// This approach avoids creating struct copies (which caused type mismatch panics)
// and isolates NaN errors to individual checks instead of losing all data.
//
// See: https://github.com/openITCOCKPIT/openitcockpit-agent-go/issues/88
func sanitizeFloats(v interface{}) interface{} {
	if v == nil {
		return v
	}

	resultMap, ok := v.(map[string]interface{})
	if !ok {
		return v
	}

	// Fast path: try marshal everything at once — no NaN means no work needed.
	// Return json.RawMessage so the caller's json.Marshal passes it through
	// without serializing a second time.
	if data, err := json.Marshal(resultMap); err == nil {
		return json.RawMessage(data)
	}

	// Slow path: at least one check has NaN/Inf — handle each key individually
	sanitized := make(map[string]interface{}, len(resultMap))
	for key, value := range resultMap {
		// Direct float values in the map (not wrapped in a struct)
		if f, ok := value.(float64); ok {
			if math.IsNaN(f) || math.IsInf(f, 0) {
				value = float64(0)
			}
			sanitized[key] = value
			continue
		}
		if f, ok := value.(float32); ok {
			if math.IsNaN(float64(f)) || math.IsInf(float64(f), 0) {
				value = float32(0)
			}
			sanitized[key] = value
			continue
		}

		// Try to sanitize nested maps/slices first (covers most check results)
		value = sanitizeValue(value)

		data, err := json.Marshal(value)
		if err != nil && strings.Contains(err.Error(), "unsupported value") {
			// NaN/Inf is inside a struct — repair struct fields in-place via reflect
			repairNaNInPlace(reflect.ValueOf(value))

			data, err = json.Marshal(value)
			if err != nil {
				log.Errorf("sanitizeFloats: check '%s' could not be repaired: %s", key, err)
				sanitized[key] = map[string]string{"error": "check result contained invalid float values"}
				continue
			}
			log.Debugf("sanitizeFloats: repaired NaN/Inf values in check '%s'", key)
		} else if err != nil {
			sanitized[key] = map[string]string{"error": err.Error()}
			continue
		}
		sanitized[key] = json.RawMessage(data)
	}
	return sanitized
}

// sanitizeValue recursively cleans NaN/Inf in dynamic types (maps, slices, floats).
// For structs, it returns the value unchanged — struct repair is done via repairNaNInPlace.
func sanitizeValue(v interface{}) interface{} {
	switch val := v.(type) {
	case float64:
		if math.IsNaN(val) || math.IsInf(val, 0) {
			return float64(0)
		}
		return val
	case float32:
		if math.IsNaN(float64(val)) || math.IsInf(float64(val), 0) {
			return float32(0)
		}
		return val
	case map[string]interface{}:
		for k, v := range val {
			val[k] = sanitizeValue(v)
		}
		return val
	case []interface{}:
		for i, v := range val {
			val[i] = sanitizeValue(v)
		}
		return val
	default:
		return v
	}
}

// repairNaNInPlace recursively walks the value via reflect and sets
// any NaN or Inf float64/float32 fields to 0. Modifies values in-place
// without creating copies, which avoids type mismatch issues.
func repairNaNInPlace(v reflect.Value) {
	switch v.Kind() {
	case reflect.Float32, reflect.Float64:
		if v.CanSet() {
			f := v.Float()
			if math.IsNaN(f) || math.IsInf(f, 0) {
				v.SetFloat(0)
			}
		}

	case reflect.Ptr:
		if !v.IsNil() {
			repairNaNInPlace(v.Elem())
		}

	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			if field.CanSet() {
				repairNaNInPlace(field)
			}
		}

	case reflect.Map:
		for _, key := range v.MapKeys() {
			elem := v.MapIndex(key)
			// Map values aren't addressable — for float values we must replace
			if elem.Kind() == reflect.Float64 {
				f := elem.Float()
				if math.IsNaN(f) || math.IsInf(f, 0) {
					v.SetMapIndex(key, reflect.ValueOf(float64(0)))
				}
			} else if elem.Kind() == reflect.Float32 {
				f := elem.Float()
				if math.IsNaN(f) || math.IsInf(f, 0) {
					v.SetMapIndex(key, reflect.ValueOf(float32(0)))
				}
			} else if elem.Kind() == reflect.Interface && !elem.IsNil() {
				// Unwrap interface and check if it's a float
				inner := elem.Elem()
				if inner.Kind() == reflect.Float64 {
					f := inner.Float()
					if math.IsNaN(f) || math.IsInf(f, 0) {
						v.SetMapIndex(key, reflect.ValueOf(float64(0)))
					}
				} else if inner.Kind() == reflect.Float32 {
				f := inner.Float()
				if math.IsNaN(f) || math.IsInf(f, 0) {
					v.SetMapIndex(key, reflect.ValueOf(float32(0)))
				}
			} else {
					repairNaNInPlace(inner)
				}
			} else {
				repairNaNInPlace(elem)
			}
		}

	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			repairNaNInPlace(v.Index(i))
		}

	case reflect.Interface:
		if !v.IsNil() {
			repairNaNInPlace(v.Elem())
		}
	}
}
