package formatter

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"text/tabwriter"
)

// TableFormatter outputs data as aligned text tables using tabwriter.
type TableFormatter struct{}

func (f *TableFormatter) Format(data any, writer io.Writer) error {
	if data == nil {
		fmt.Fprintln(writer, "No data")
		return nil
	}

	v := reflect.ValueOf(data)

	// Dereference pointer
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			fmt.Fprintln(writer, "No data")
			return nil
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		return formatSlice(v, writer)
	case reflect.Struct:
		return formatStruct(v, writer)
	case reflect.Map:
		return formatMap(v, writer)
	default:
		_, err := fmt.Fprintf(writer, "%v\n", v.Interface())
		return err
	}
}

func formatSlice(v reflect.Value, writer io.Writer) error {
	if v.Len() == 0 {
		fmt.Fprintln(writer, "No data")
		return nil
	}

	// Check if elements are structs (or pointer to struct)
	elemType := v.Type().Elem()
	if elemType.Kind() == reflect.Ptr {
		elemType = elemType.Elem()
	}

	if elemType.Kind() == reflect.Struct {
		return formatStructSlice(v, elemType, writer)
	}

	// Fallback: one value per line
	tw := tabwriter.NewWriter(writer, 0, 4, 2, ' ', 0)
	for i := 0; i < v.Len(); i++ {
		fmt.Fprintln(tw, stringify(v.Index(i).Interface()))
	}
	return tw.Flush()
}

func formatStructSlice(v reflect.Value, elemType reflect.Type, writer io.Writer) error {
	tw := tabwriter.NewWriter(writer, 0, 4, 2, ' ', 0)

	// Print header from exported field names or json tags
	headers := structHeaders(elemType)
	for i, h := range headers {
		if i > 0 {
			fmt.Fprint(tw, "\t")
		}
		fmt.Fprint(tw, h)
	}
	fmt.Fprintln(tw)

	// Print rows
	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i)
		if elem.Kind() == reflect.Ptr {
			elem = elem.Elem()
		}
		for j := 0; j < elem.NumField(); j++ {
			if !elem.Type().Field(j).IsExported() {
				continue
			}
			if j > 0 {
				fmt.Fprint(tw, "\t")
			}
			fmt.Fprint(tw, stringify(elem.Field(j).Interface()))
		}
		fmt.Fprintln(tw)
	}

	return tw.Flush()
}

func formatStruct(v reflect.Value, writer io.Writer) error {
	tw := tabwriter.NewWriter(writer, 0, 4, 2, ' ', 0)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		name := jsonFieldName(field)
		fmt.Fprintf(tw, "%s\t%s\n", name, stringify(v.Field(i).Interface()))
	}

	return tw.Flush()
}

func formatMap(v reflect.Value, writer io.Writer) error {
	tw := tabwriter.NewWriter(writer, 0, 4, 2, ' ', 0)

	for _, key := range v.MapKeys() {
		fmt.Fprintf(tw, "%s\t%s\n", stringify(key.Interface()), stringify(v.MapIndex(key).Interface()))
	}

	return tw.Flush()
}

func structHeaders(t reflect.Type) []string {
	var headers []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		headers = append(headers, jsonFieldName(field))
	}
	return headers
}

func jsonFieldName(field reflect.StructField) string {
	tag := field.Tag.Get("json")
	if tag != "" && tag != "-" {
		// Take the name part before any comma
		for i := 0; i < len(tag); i++ {
			if tag[i] == ',' {
				return tag[:i]
			}
		}
		return tag
	}
	return field.Name
}

func stringify(v any) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case json.Number:
		return string(val)
	default:
		return fmt.Sprintf("%v", val)
	}
}

// formatQueryResultsTable formats query results as an aligned table.
func formatQueryResultsTable(columns []string, rows [][]any, writer io.Writer) error {
	tw := tabwriter.NewWriter(writer, 0, 4, 2, ' ', 0)

	// Header
	for i, col := range columns {
		if i > 0 {
			fmt.Fprint(tw, "\t")
		}
		fmt.Fprint(tw, col)
	}
	fmt.Fprintln(tw)

	// Rows
	for _, row := range rows {
		for i, val := range row {
			if i > 0 {
				fmt.Fprint(tw, "\t")
			}
			fmt.Fprint(tw, stringify(val))
		}
		fmt.Fprintln(tw)
	}

	return tw.Flush()
}
