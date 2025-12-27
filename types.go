package sendpigeon

// EmailStatus represents the status of an email.
type EmailStatus string

const (
	EmailStatusScheduled  EmailStatus = "scheduled"
	EmailStatusCancelled  EmailStatus = "cancelled"
	EmailStatusPending    EmailStatus = "pending"
	EmailStatusSent       EmailStatus = "sent"
	EmailStatusDelivered  EmailStatus = "delivered"
	EmailStatusBounced    EmailStatus = "bounced"
	EmailStatusComplained EmailStatus = "complained"
	EmailStatusFailed     EmailStatus = "failed"
)

// DomainStatus represents the status of a domain.
type DomainStatus string

const (
	DomainStatusPending          DomainStatus = "pending"
	DomainStatusVerified         DomainStatus = "verified"
	DomainStatusTemporaryFailure DomainStatus = "temporary_failure"
	DomainStatusFailed           DomainStatus = "failed"
)

// APIKeyMode represents the mode of an API key.
type APIKeyMode string

const (
	APIKeyModeLive APIKeyMode = "live"
	APIKeyModeTest APIKeyMode = "test"
)

// APIKeyPermission represents the permission level of an API key.
type APIKeyPermission string

const (
	APIKeyPermissionFullAccess APIKeyPermission = "full_access"
	APIKeyPermissionSending    APIKeyPermission = "sending"
)

// Attachment represents an email attachment.
type Attachment struct {
	Filename    string `json:"filename"`
	Content     string `json:"content,omitempty"`
	Path        string `json:"path,omitempty"`
	ContentType string `json:"contentType,omitempty"`
}

// AttachmentMeta represents attachment metadata returned from API.
type AttachmentMeta struct {
	Filename    string `json:"filename"`
	Size        int    `json:"size"`
	ContentType string `json:"content_type"`
}

// SendEmailRequest represents a request to send an email.
type SendEmailRequest struct {
	To             []string          `json:"to"`
	From           string            `json:"from,omitempty"`
	Subject        string            `json:"subject,omitempty"`
	HTML           string            `json:"html,omitempty"`
	Text           string            `json:"text,omitempty"`
	CC             []string          `json:"cc,omitempty"`
	BCC            []string          `json:"bcc,omitempty"`
	ReplyTo        string            `json:"replyTo,omitempty"`
	TemplateID     string            `json:"templateId,omitempty"`
	Variables      map[string]string `json:"variables,omitempty"`
	Attachments    []Attachment      `json:"attachments,omitempty"`
	Tags           []string          `json:"tags,omitempty"`
	Metadata       map[string]string `json:"metadata,omitempty"`
	Headers        map[string]string `json:"headers,omitempty"`
	ScheduledAt    string            `json:"scheduled_at,omitempty"`
	IdempotencyKey string            `json:"-"` // Sent as header
}

// SendEmailResponse represents the response from sending an email.
type SendEmailResponse struct {
	ID          string      `json:"id"`
	Status      EmailStatus `json:"status"`
	ScheduledAt string      `json:"scheduled_at,omitempty"`
	Suppressed  []string    `json:"suppressed,omitempty"`
}

// BatchEmailResult represents the result for a single email in a batch.
type BatchEmailResult struct {
	Index      int                    `json:"index"`
	Status     string                 `json:"status"`
	ID         string                 `json:"id,omitempty"`
	Suppressed []string               `json:"suppressed,omitempty"`
	Error      map[string]interface{} `json:"error,omitempty"`
}

// SendBatchResponse represents the response from sending batch emails.
type SendBatchResponse struct {
	Data    []BatchEmailResult     `json:"data"`
	Summary map[string]interface{} `json:"summary"`
}

// EmailDetail represents detailed email information.
type EmailDetail struct {
	ID            string           `json:"id"`
	FromAddress   string           `json:"from_address"`
	ToAddress     string           `json:"to_address"`
	Subject       string           `json:"subject"`
	Status        EmailStatus      `json:"status"`
	CreatedAt     string           `json:"created_at"`
	CCAddress     string           `json:"cc_address,omitempty"`
	BCCAddress    string           `json:"bcc_address,omitempty"`
	Tags          []string         `json:"tags,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
	SentAt        string           `json:"sent_at,omitempty"`
	DeliveredAt   string           `json:"delivered_at,omitempty"`
	BouncedAt     string           `json:"bounced_at,omitempty"`
	ComplainedAt  string           `json:"complained_at,omitempty"`
	BounceType    string           `json:"bounce_type,omitempty"`
	ComplaintType string           `json:"complaint_type,omitempty"`
	Attachments   []AttachmentMeta `json:"attachments,omitempty"`
	HasBody       bool             `json:"has_body"`
}

// Template represents an email template.
type Template struct {
	ID        string                 `json:"id"`
	Name      string                 `json:"name"`
	Subject   string                 `json:"subject"`
	Variables []string               `json:"variables"`
	CreatedAt string                 `json:"created_at"`
	UpdatedAt string                 `json:"updated_at"`
	HTML      string                 `json:"html,omitempty"`
	Text      string                 `json:"text,omitempty"`
	Domain    map[string]interface{} `json:"domain,omitempty"`
}

// CreateTemplateRequest represents a request to create a template.
type CreateTemplateRequest struct {
	Name     string `json:"name"`
	Subject  string `json:"subject"`
	HTML     string `json:"html,omitempty"`
	Text     string `json:"text,omitempty"`
	DomainID string `json:"domainId,omitempty"`
}

// UpdateTemplateRequest represents a request to update a template.
type UpdateTemplateRequest struct {
	Name    string `json:"name,omitempty"`
	Subject string `json:"subject,omitempty"`
	HTML    string `json:"html,omitempty"`
	Text    string `json:"text,omitempty"`
}

// DNSRecord represents a DNS record for domain verification.
type DNSRecord struct {
	Type     string `json:"type"`
	Name     string `json:"name"`
	Value    string `json:"value"`
	Priority int    `json:"priority,omitempty"`
}

// Domain represents domain information.
type Domain struct {
	ID            string       `json:"id"`
	Name          string       `json:"name"`
	Status        DomainStatus `json:"status"`
	CreatedAt     string       `json:"created_at"`
	VerifiedAt    string       `json:"verified_at,omitempty"`
	LastCheckedAt string       `json:"last_checked_at,omitempty"`
	FailingSince  string       `json:"failing_since,omitempty"`
}

// DomainWithDNSRecords represents domain with DNS records for setup.
type DomainWithDNSRecords struct {
	Domain
	DNSRecords []DNSRecord `json:"dns_records"`
}

// DomainVerificationResult represents result of domain verification.
type DomainVerificationResult struct {
	Verified   bool         `json:"verified"`
	Status     DomainStatus `json:"status"`
	DNSRecords []DNSRecord  `json:"dns_records"`
}

// APIKey represents API key information (without secret).
type APIKey struct {
	ID         string                 `json:"id"`
	Name       string                 `json:"name"`
	KeyPrefix  string                 `json:"key_prefix"`
	Mode       APIKeyMode             `json:"mode"`
	Permission APIKeyPermission       `json:"permission"`
	CreatedAt  string                 `json:"created_at"`
	LastUsedAt string                 `json:"last_used_at,omitempty"`
	ExpiresAt  string                 `json:"expires_at,omitempty"`
	Domain     map[string]interface{} `json:"domain,omitempty"`
}

// APIKeyWithSecret represents API key with secret (only on creation).
type APIKeyWithSecret struct {
	APIKey
	Key string `json:"key"`
}

// CreateAPIKeyRequest represents a request to create an API key.
type CreateAPIKeyRequest struct {
	Name       string           `json:"name"`
	Mode       APIKeyMode       `json:"mode,omitempty"`
	Permission APIKeyPermission `json:"permission,omitempty"`
	DomainID   string           `json:"domainId,omitempty"`
	ExpiresAt  string           `json:"expiresAt,omitempty"`
}

// ListOptions represents options for list endpoints.
type ListOptions struct {
	Limit  int    `json:"limit,omitempty"`
	Offset int    `json:"offset,omitempty"`
	Cursor string `json:"cursor,omitempty"`
}

// ListEmailsOptions represents options for listing emails.
type ListEmailsOptions struct {
	ListOptions
	Status EmailStatus `json:"status,omitempty"`
	Tag    string      `json:"tag,omitempty"`
}
