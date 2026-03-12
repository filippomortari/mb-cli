package client

import (
	"encoding/json"
	"strconv"
	"strings"
)

func buildCardQueryParameters(card *Card, params map[string]string) []QueryParameter {
	if len(params) == 0 {
		return nil
	}

	resolved := make([]QueryParameter, 0, len(params))
	tags := map[string]TemplateTag(nil)
	if card != nil && card.DatasetQuery != nil && card.DatasetQuery.Native != nil {
		tags = card.DatasetQuery.Native.TemplateTags
	}

	for key, rawValue := range params {
		param := QueryParameter{
			ID:    key,
			Value: coerceQueryParameterValue(rawValue),
		}
		if _, tag, ok := resolveTemplateTag(tags, key); ok && tag.ID != "" {
			param.ID = tag.ID
		}
		resolved = append(resolved, param)
	}

	return resolved
}

func buildDashboardQueryParameters(dashboard *Dashboard, params map[string]string) []QueryParameter {
	if len(params) == 0 {
		return nil
	}

	resolved := make([]QueryParameter, 0, len(params))
	for key, rawValue := range params {
		param := QueryParameter{
			ID:    key,
			Value: coerceQueryParameterValue(rawValue),
		}
		if dashboardParam, ok := resolveDashboardParameterValue(dashboard, key); ok {
			param.ID = dashboardParam.ID
			param.Type = dashboardParam.Type
		}
		resolved = append(resolved, param)
	}

	return resolved
}

func resolveTemplateTag(tags map[string]TemplateTag, input string) (string, TemplateTag, bool) {
	if len(tags) == 0 {
		return "", TemplateTag{}, false
	}
	if tag, ok := tags[input]; ok {
		return input, tag, true
	}

	for key, tag := range tags {
		if tag.ID == input || strings.EqualFold(tag.Name, input) || strings.EqualFold(tag.DisplayName, input) {
			return key, tag, true
		}
	}

	return "", TemplateTag{}, false
}

func resolveDashboardParameterValue(dashboard *Dashboard, input string) (*DashParameter, bool) {
	if dashboard == nil {
		return nil, false
	}
	for i := range dashboard.Parameters {
		parameter := &dashboard.Parameters[i]
		if parameter.ID == input || parameter.Slug == input || strings.EqualFold(parameter.Name, input) {
			return parameter, true
		}
	}

	return nil, false
}

func coerceQueryParameterValue(raw string) any {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return ""
	}

	if strings.HasPrefix(trimmed, "[") || strings.HasPrefix(trimmed, "{") || strings.HasPrefix(trimmed, `"`) {
		var decoded any
		if err := json.Unmarshal([]byte(trimmed), &decoded); err == nil {
			return decoded
		}
	}

	switch strings.ToLower(trimmed) {
	case "true":
		return true
	case "false":
		return false
	case "null":
		return nil
	}

	if i, err := strconv.Atoi(trimmed); err == nil {
		return i
	}
	if f, err := strconv.ParseFloat(trimmed, 64); err == nil {
		return f
	}

	return raw
}
