package engine

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/smtp"
	"strings"
	"time"
)

// Node represents a workflow node
type Node struct {
	ID       string                 `json:"id"`
	Type     string                 `json:"type"`
	Data     map[string]interface{} `json:"data"`
	Position Position               `json:"position"`
}

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// NodeExecutor interface for all node types
type NodeExecutor interface {
	Execute(ctx context.Context, node *Node, input map[string]interface{}) (map[string]interface{}, error)
}

// HTTPRequestExecutor executes HTTP requests
type HTTPRequestExecutor struct{}

func (h *HTTPRequestExecutor) Execute(ctx context.Context, node *Node, input map[string]interface{}) (map[string]interface{}, error) {
	config, ok := node.Data["config"].(map[string]interface{})
	if !ok {
		return nil, errors.New("invalid http request configuration")
	}

	method, _ := config["method"].(string)
	url, _ := config["url"].(string)
	headers, _ := config["headers"].(map[string]interface{})
	body, _ := config["body"]

	if method == "" {
		method = "GET"
	}
	if url == "" {
		return nil, errors.New("url is required")
	}

	url = replaceVariables(url, input)

	var bodyReader io.Reader
	if body != nil {
		bodyJSON, _ := json.Marshal(body)
		bodyReader = strings.NewReader(string(bodyJSON))
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for key, value := range headers {
		if strValue, ok := value.(string); ok {
			req.Header.Set(key, strValue)
		}
	}

	if req.Header.Get("Content-Type") == "" && bodyReader != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	var responseData map[string]interface{}
	if err := json.Unmarshal(respBody, &responseData); err != nil {
		responseData = map[string]interface{}{
			"body": string(respBody),
		}
	}

	return map[string]interface{}{
		"http_response": responseData,
		"status_code":   resp.StatusCode,
	}, nil
}

// EmailExecutor sends emails
type EmailExecutor struct{}

func (e *EmailExecutor) Execute(ctx context.Context, node *Node, input map[string]interface{}) (map[string]interface{}, error) {
	config, ok := node.Data["config"].(map[string]interface{})
	if !ok {
		return nil, errors.New("invalid email configuration")
	}

	to, _ := config["to"].(string)
	subject, _ := config["subject"].(string)
	body, _ := config["body"].(string)
	from, _ := config["from"].(string)
	smtpHost, _ := config["smtp_host"].(string)
	smtpPort, _ := config["smtp_port"].(string)
	smtpUser, _ := config["smtp_user"].(string)
	smtpPass, _ := config["smtp_pass"].(string)

	if to == "" || subject == "" || body == "" {
		return nil, errors.New("to, subject, and body are required")
	}

	to = replaceVariables(to, input)
	subject = replaceVariables(subject, input)
	body = replaceVariables(body, input)

	if smtpHost == "" {
		smtpHost = "smtp.gmail.com"
	}
	if smtpPort == "" {
		smtpPort = "587"
	}
	if from == "" {
		from = smtpUser
	}

	message := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s\r\n", to, subject, body))

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	err := smtp.SendMail(addr, auth, from, []string{to}, message)
	if err != nil {
		return nil, fmt.Errorf("failed to send email: %w", err)
	}

	return map[string]interface{}{
		"email_sent": true,
		"to":         to,
		"subject":    subject,
	}, nil
}

// WebhookExecutor handles webhook triggers
type WebhookExecutor struct{}

func (w *WebhookExecutor) Execute(ctx context.Context, node *Node, input map[string]interface{}) (map[string]interface{}, error) {
	return input, nil
}

// DelayExecutor waits for a specified duration
type DelayExecutor struct{}

func (d *DelayExecutor) Execute(ctx context.Context, node *Node, input map[string]interface{}) (map[string]interface{}, error) {
	config, ok := node.Data["config"].(map[string]interface{})
	if !ok {
		return nil, errors.New("invalid delay configuration")
	}

	delaySeconds, ok := config["seconds"].(float64)
	if !ok {
		delaySeconds = 1
	}

	duration := time.Duration(delaySeconds) * time.Second

	select {
	case <-time.After(duration):
		return input, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// IfExecutor handles conditional logic
type IfExecutor struct{}

func (i *IfExecutor) Execute(ctx context.Context, node *Node, input map[string]interface{}) (map[string]interface{}, error) {
	config, ok := node.Data["config"].(map[string]interface{})
	if !ok {
		return nil, errors.New("invalid if configuration")
	}

	condition, _ := config["condition"].(string)
	result := evaluateCondition(condition, input)

	return map[string]interface{}{
		"condition_result": result,
	}, nil
}

// Helper functions
func replaceVariables(text string, data map[string]interface{}) string {
	result := text
	for key, value := range data {
		placeholder := fmt.Sprintf("{{%s}}", key)
		if strValue, ok := value.(string); ok {
			result = strings.ReplaceAll(result, placeholder, strValue)
		} else {
			result = strings.ReplaceAll(result, placeholder, fmt.Sprintf("%v", value))
		}
	}
	return result
}

func evaluateCondition(condition string, data map[string]interface{}) bool {
	if strings.Contains(condition, "==") {
		parts := strings.Split(condition, "==")
		if len(parts) == 2 {
			left := strings.TrimSpace(parts[0])
			right := strings.TrimSpace(parts[1])

			leftValue := data[left]
			return fmt.Sprintf("%v", leftValue) == right
		}
	}

	return false
}
