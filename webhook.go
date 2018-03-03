package paypalsdk

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// CreateWebhook creates a new webhook in Paypal
//
// Allows for the customisation of the payment experience
//
// Endpoint: POST /v1/notifications/webhooks
func (c *Client) CreateWebhook(w Webhook) (*Webhook, error) {
	url := fmt.Sprintf("%s/v1/notifications/webhooks", c.APIBase)
	req, err := c.NewRequest("POST", url, w)
	if err != nil {
		return &Webhook{}, err
	}

	response := &Webhook{}

	if err = c.SendWithAuth(req, response); err != nil {
		return response, err
	}

	return response, nil
}

// GetWebhook gets an exists payment experience from Paypal
//
// Endpoint: GET /v1/notifications/webhooks/<webhook-id>
func (c *Client) GetWebhook(webhookID string) (*Webhook, error) {
	var w Webhook

	url := fmt.Sprintf("%s/v1/notifications/webhooks/%s", c.APIBase, webhookID)
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return &w, err
	}

	if err = c.SendWithAuth(req, &w); err != nil {
		return &w, err
	}

	if w.ID == "" {
		return &w, fmt.Errorf("paypalsdk: unable to get web webhook with ID = %s", webhookID)
	}

	return &w, nil
}

// GetWebhooks retrieves webhooks from Paypal
//
// Endpoint: GET /v1/notifications/webhooks
func (c *Client) GetWebhooks() ([]Webhook, error) {
	var ws struct {
		Webhooks []Webhook `json:"webhooks"`
	}

	url := fmt.Sprintf("%s/v1/notifications/webhooks", c.APIBase)
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return ws.Webhooks, err
	}

	if err = c.SendWithAuth(req, &ws); err != nil {
		return ws.Webhooks, err
	}

	return ws.Webhooks, nil
}

// SetWebhook sets a webhook in Paypal with given id
//
// Endpoint: PUT /v1/notifications/webhooks
func (c *Client) SetWebhook(w Webhook) error {
	if w.ID == "" {
		return fmt.Errorf("paypalsdk: no ID specified for Webhook")
	}

	url := fmt.Sprintf("%s/v1/notifications/webhooks/%s", c.APIBase, w.ID)

	p := []WebhookPatch{
		WebhookPatch{
			Operation: "replace",
			Path:      "/url",
			Value:     w.URL,
		},
		WebhookPatch{
			Operation: "replace",
			Path:      "/event_types",
			Value:     w.EventTypes,
		},
	}

	req, err := c.NewRequest("POST", url, p)

	if err != nil {
		return err
	}

	if err = c.SendWithAuth(req, nil); err != nil {
		return err
	}

	return nil
}

// DeleteWebhook deletes a webhook from Paypal with given id
//
// Endpoint: DELETE /v1/notifications/webhooks
func (c *Client) DeleteWebhook(webhookID string) error {
	url := fmt.Sprintf("%s/v1/notifications/webhooks/%s", c.APIBase, webhookID)

	req, err := c.NewRequest("DELETE", url, nil)

	if err != nil {
		return err
	}

	if err = c.SendWithAuth(req, nil); err != nil {
		return err
	}

	return nil
}

// GetWebhookEventTypes lists events to which an app can subscribe.
//
// Endpoint: GET /v1/notifications/webhooks-event-types
func (c *Client) GetWebhookEventTypes() (*EventTypeList, error) {
	var e EventTypeList

	url := fmt.Sprintf("%s/v1/notifications/webhooks-event-types", c.APIBase)
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return &e, err
	}

	if err = c.SendWithAuth(req, &e); err != nil {
		return &e, err
	}

	return &e, nil
}

// GetWebhookEvent get event notification details
//
// Endpoint: GET /v1/notifications/webhook-events/<event-id>
func (c *Client) GetWebhookEvent(eventID string) (*Event, error) {
	var e Event

	url := fmt.Sprintf("%s/v1/notifications/webhook-events/%s", c.APIBase, eventID)
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return &e, err
	}

	if err = c.SendWithAuth(req, &e); err != nil {
		return &e, err
	}

	if e.ID == "" {
		return &e, fmt.Errorf("paypalsdk: unable to get web event with ID = %s", eventID)
	}

	return &e, nil
}

// ResendWebhookEvent resends an event notification, by event ID.
//
// Endpoint: POST /v1/notifications/webhook-events/<event-id>/resend
func (c *Client) ResendWebhookEvent(eventID string, webhookIDs []string) (*Event, error) {
	var e Event

	url := fmt.Sprintf("%s/v1/notifications/webhook-events/%s/resend", c.APIBase, eventID)

	w := WebhookIDList{
		WebhookIDs: webhookIDs,
	}

	req, err := c.NewRequest("POST", url, w)

	if err != nil {
		return nil, err
	}

	if err = c.SendWithAuth(req, &e); err != nil {
		return nil, err
	}

	return &e, nil
}

// GetWebhookEvents get event notification details
//
// Endpoint: GET /v1/notifications/webhook-events
func (c *Client) GetWebhookEvents(f GetWebhookEventsFilter) (*[]Event, error) {
	var e []Event

	var qs url.Values
	if f.PageSize != 0 {
		qs.Set("page_size", strconv.FormatInt(int64(f.PageSize), 10))
	}
	if !f.StartTime.IsZero() {
		qs.Set("start_time", f.StartTime.UTC().Format(time.RFC3339))
	}
	if !f.EndTime.IsZero() {
		qs.Set("end_time", f.EndTime.UTC().Format(time.RFC3339))
	}
	if f.TransactionID != "" {
		qs.Set("transaction_id", f.TransactionID)
	}
	if f.EventType != "" {
		qs.Set("event_type", f.EventType)
	}

	url := fmt.Sprintf("%s/v1/notifications/webhook-events?%s", c.APIBase, qs.Encode())
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return &e, err
	}

	if err = c.SendWithAuth(req, &e); err != nil {
		return &e, err
	}

	return &e, nil
}

// SimulateWebhookEvent simulates a webhook event
//
// Endpoint: POST /v1/notifications/simulate-event
func (c *Client) SimulateWebhookEvent(r SimulateEventReq) (*Event, error) {
	var e Event

	url := fmt.Sprintf("%s/v1/notifications/simulate-event", c.APIBase)

	buf, _ := json.Marshal(r)
	log.Println(string(buf))

	req, err := c.NewRequest("POST", url, r)

	if err != nil {
		return nil, err
	}

	if err = c.SendWithAuth(req, &e); err != nil {
		return nil, err
	}

	return &e, nil
}

// VerifyWebhookSignature verify a webhook signature.
//
// Endpoint: POST /v1/notifications/verify-webhook-signature
func (c *Client) VerifyWebhookSignature(r WebhookRequest) (*VerificationStatus, error) {
	var e VerificationStatus

	url := fmt.Sprintf("%s/v1/notifications/verify-webhook-signature", c.APIBase)

	req, err := c.NewRequest("POST", url, r)

	if err != nil {
		return nil, err
	}

	if err = c.SendWithAuth(req, &e); err != nil {
		return nil, err
	}

	return &e, nil
}
