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

// TemplateStatus represents the status of a template.
type TemplateStatus string

const (
	TemplateStatusDraft     TemplateStatus = "draft"
	TemplateStatusPublished TemplateStatus = "published"
)

// TemplateVariableType represents the type of a template variable.
type TemplateVariableType string

const (
	TemplateVariableTypeString  TemplateVariableType = "string"
	TemplateVariableTypeNumber  TemplateVariableType = "number"
	TemplateVariableTypeBoolean TemplateVariableType = "boolean"
)

// TrackingOptions represents per-email tracking options.
type TrackingOptions struct {
	Opens  *bool `json:"opens,omitempty"`
	Clicks *bool `json:"clicks,omitempty"`
}

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
	Tracking       *TrackingOptions  `json:"tracking,omitempty"`
	IdempotencyKey string            `json:"-"` // Sent as header
}

// SendEmailResponse represents the response from sending an email.
type SendEmailResponse struct {
	ID          string      `json:"id"`
	Status      EmailStatus `json:"status"`
	ScheduledAt string      `json:"scheduled_at,omitempty"`
	Suppressed  []string    `json:"suppressed,omitempty"`
	Warnings    []string    `json:"warnings,omitempty"`
}

// BatchEmailResult represents the result for a single email in a batch.
type BatchEmailResult struct {
	Index      int                    `json:"index"`
	Status     string                 `json:"status"`
	ID         string                 `json:"id,omitempty"`
	Suppressed []string               `json:"suppressed,omitempty"`
	Warnings   []string               `json:"warnings,omitempty"`
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

// TemplateVariable represents a typed variable in a template.
type TemplateVariable struct {
	Key           string               `json:"key"`
	Type          TemplateVariableType `json:"type"`
	FallbackValue string               `json:"fallbackValue,omitempty"`
}

// Template represents an email template.
type Template struct {
	ID         string                 `json:"id"`
	TemplateID string                 `json:"templateId"`
	Name       string                 `json:"name,omitempty"`
	Subject    string                 `json:"subject"`
	Variables  []TemplateVariable     `json:"variables"`
	Status     TemplateStatus         `json:"status"`
	CreatedAt  string                 `json:"createdAt"`
	UpdatedAt  string                 `json:"updatedAt"`
	HTML       string                 `json:"html,omitempty"`
	Text       string                 `json:"text,omitempty"`
	Domain     map[string]interface{} `json:"domain,omitempty"`
}

// CreateTemplateRequest represents a request to create a template.
type CreateTemplateRequest struct {
	TemplateID string             `json:"templateId"`
	Name       string             `json:"name,omitempty"`
	Subject    string             `json:"subject"`
	HTML       string             `json:"html,omitempty"`
	Text       string             `json:"text,omitempty"`
	Variables  []TemplateVariable `json:"variables,omitempty"`
	DomainID   string             `json:"domainId,omitempty"`
}

// UpdateTemplateRequest represents a request to update a template.
type UpdateTemplateRequest struct {
	Name      string             `json:"name,omitempty"`
	Subject   string             `json:"subject,omitempty"`
	HTML      string             `json:"html,omitempty"`
	Text      string             `json:"text,omitempty"`
	Variables []TemplateVariable `json:"variables,omitempty"`
}

// TestTemplateRequest represents a request to test a template.
type TestTemplateRequest struct {
	To        string            `json:"to"`
	Variables map[string]string `json:"variables,omitempty"`
}

// TestTemplateResponse represents the response from testing a template.
type TestTemplateResponse struct {
	Message string `json:"message"`
	EmailID string `json:"emailId"`
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

// SuppressionReason represents the reason for suppression.
type SuppressionReason string

const (
	SuppressionReasonHardBounce SuppressionReason = "hard_bounce"
	SuppressionReasonComplaint  SuppressionReason = "complaint"
)

// Suppression represents a suppressed email address.
type Suppression struct {
	ID            string            `json:"id"`
	Email         string            `json:"email"`
	Reason        SuppressionReason `json:"reason"`
	SourceEmailID string            `json:"sourceEmailId,omitempty"`
	CreatedAt     string            `json:"createdAt"`
}

// SuppressionListResponse represents the response from listing suppressions.
type SuppressionListResponse struct {
	Data  []Suppression `json:"data"`
	Total int           `json:"total"`
}

// TrackingDefaults represents organization tracking defaults.
type TrackingDefaults struct {
	TrackingEnabled     bool `json:"trackingEnabled"`
	PrivacyMode         bool `json:"privacyMode"`
	WebhookOnEveryOpen  bool `json:"webhookOnEveryOpen"`
	WebhookOnEveryClick bool `json:"webhookOnEveryClick"`
}

// UpdateTrackingDefaultsRequest represents a request to update tracking defaults.
type UpdateTrackingDefaultsRequest struct {
	TrackingEnabled     *bool `json:"trackingEnabled,omitempty"`
	PrivacyMode         *bool `json:"privacyMode,omitempty"`
	WebhookOnEveryOpen  *bool `json:"webhookOnEveryOpen,omitempty"`
	WebhookOnEveryClick *bool `json:"webhookOnEveryClick,omitempty"`
}

// ContactStatus represents the status of a contact.
type ContactStatus string

const (
	ContactStatusActive       ContactStatus = "ACTIVE"
	ContactStatusUnsubscribed ContactStatus = "UNSUBSCRIBED"
	ContactStatusBounced      ContactStatus = "BOUNCED"
	ContactStatusComplained   ContactStatus = "COMPLAINED"
)

// Contact represents a contact in the audience.
type Contact struct {
	ID             string            `json:"id"`
	Email          string            `json:"email"`
	Fields         map[string]string `json:"fields,omitempty"`
	Tags           []string          `json:"tags,omitempty"`
	Status         ContactStatus     `json:"status"`
	UnsubscribedAt string            `json:"unsubscribedAt,omitempty"`
	BouncedAt      string            `json:"bouncedAt,omitempty"`
	ComplainedAt   string            `json:"complainedAt,omitempty"`
	CreatedAt      string            `json:"createdAt"`
	UpdatedAt      string            `json:"updatedAt"`
}

// CreateContactRequest represents a request to create a contact.
type CreateContactRequest struct {
	Email  string            `json:"email"`
	Fields map[string]string `json:"fields,omitempty"`
	Tags   []string          `json:"tags,omitempty"`
}

// UpdateContactRequest represents a request to update a contact.
type UpdateContactRequest struct {
	Email  string            `json:"email,omitempty"`
	Fields map[string]string `json:"fields,omitempty"`
	Tags   []string          `json:"tags,omitempty"`
}

// BatchContactInput represents a single contact in a batch operation.
type BatchContactInput struct {
	Email  string            `json:"email"`
	Fields map[string]string `json:"fields,omitempty"`
	Tags   []string          `json:"tags,omitempty"`
}

// BatchContactResult represents the result for a single contact in batch.
type BatchContactResult struct {
	Index   int      `json:"index"`
	Status  string   `json:"status"`
	ID      string   `json:"id,omitempty"`
	Email   string   `json:"email,omitempty"`
	Error   string   `json:"error,omitempty"`
	Message string   `json:"message,omitempty"`
}

// BatchContactResponse represents the response from batch contact creation.
type BatchContactResponse struct {
	Data    []BatchContactResult `json:"data"`
	Summary struct {
		Total   int `json:"total"`
		Created int `json:"created"`
		Updated int `json:"updated"`
		Failed  int `json:"failed"`
	} `json:"summary"`
}

// AudienceStats represents contact statistics.
type AudienceStats struct {
	Total        int `json:"total"`
	Active       int `json:"active"`
	Unsubscribed int `json:"unsubscribed"`
	Bounced      int `json:"bounced"`
	Complained   int `json:"complained"`
}

// ListContactsOptions represents options for listing contacts.
type ListContactsOptions struct {
	Limit  int      `json:"limit,omitempty"`
	Offset int      `json:"offset,omitempty"`
	Cursor string   `json:"cursor,omitempty"`
	Tags   []string `json:"tags,omitempty"`
	Status string   `json:"status,omitempty"`
	Search string   `json:"search,omitempty"`
}

// BroadcastStatus represents the status of a broadcast.
type BroadcastStatus string

const (
	BroadcastStatusDraft     BroadcastStatus = "DRAFT"
	BroadcastStatusScheduled BroadcastStatus = "SCHEDULED"
	BroadcastStatusSending   BroadcastStatus = "SENDING"
	BroadcastStatusSent      BroadcastStatus = "SENT"
	BroadcastStatusCancelled BroadcastStatus = "CANCELLED"
	BroadcastStatusFailed    BroadcastStatus = "FAILED"
)

// BroadcastStats represents broadcast delivery statistics.
type BroadcastStats struct {
	Total      int `json:"total"`
	Sent       int `json:"sent"`
	Delivered  int `json:"delivered"`
	Opened     int `json:"opened"`
	Clicked    int `json:"clicked"`
	Bounced    int `json:"bounced"`
	Complained int `json:"complained"`
	Failed     int `json:"failed"`
}

// Broadcast represents a broadcast campaign.
type Broadcast struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Subject     string          `json:"subject"`
	PreviewText string          `json:"previewText,omitempty"`
	FromEmail   string          `json:"fromEmail,omitempty"`
	FromName    string          `json:"fromName,omitempty"`
	ReplyTo     string          `json:"replyTo,omitempty"`
	HTMLContent string          `json:"htmlContent,omitempty"`
	TextContent string          `json:"textContent,omitempty"`
	Tags        []string        `json:"tags,omitempty"`
	Status      BroadcastStatus `json:"status"`
	Stats       *BroadcastStats `json:"stats,omitempty"`
	ScheduledAt string          `json:"scheduledAt,omitempty"`
	SentAt      string          `json:"sentAt,omitempty"`
	CreatedAt   string          `json:"createdAt"`
	UpdatedAt   string          `json:"updatedAt"`
}

// CreateBroadcastRequest represents a request to create a broadcast.
type CreateBroadcastRequest struct {
	Name        string   `json:"name"`
	Subject     string   `json:"subject"`
	PreviewText string   `json:"previewText,omitempty"`
	FromEmail   string   `json:"fromEmail,omitempty"`
	FromName    string   `json:"fromName,omitempty"`
	ReplyTo     string   `json:"replyTo,omitempty"`
	HTMLContent string   `json:"htmlContent,omitempty"`
	TextContent string   `json:"textContent,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	TemplateID  string   `json:"templateId,omitempty"`
}

// UpdateBroadcastRequest represents a request to update a broadcast.
type UpdateBroadcastRequest struct {
	Name        string   `json:"name,omitempty"`
	Subject     string   `json:"subject,omitempty"`
	PreviewText string   `json:"previewText,omitempty"`
	FromEmail   string   `json:"fromEmail,omitempty"`
	FromName    string   `json:"fromName,omitempty"`
	ReplyTo     string   `json:"replyTo,omitempty"`
	HTMLContent string   `json:"htmlContent,omitempty"`
	TextContent string   `json:"textContent,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

// BroadcastTargeting represents tag-based targeting for broadcasts.
type BroadcastTargeting struct {
	// Only send to contacts with ANY of these tags. Empty = all active contacts.
	IncludeTags []string `json:"includeTags,omitempty"`
	// Exclude contacts with ANY of these tags.
	ExcludeTags []string `json:"excludeTags,omitempty"`
}

// SendBroadcastRequest represents a request to send a broadcast.
type SendBroadcastRequest struct {
	BroadcastTargeting
}

// ScheduleBroadcastRequest represents a request to schedule a broadcast.
type ScheduleBroadcastRequest struct {
	ScheduledAt string `json:"scheduledAt"`
	BroadcastTargeting
}

// TestBroadcastRequest represents a request to send a test broadcast.
type TestBroadcastRequest struct {
	To []string `json:"to"`
}

// TestBroadcastResponse represents the response from sending a test.
type TestBroadcastResponse struct {
	Message  string   `json:"message"`
	EmailIDs []string `json:"emailIds"`
}

// BroadcastRecipientStatus represents the status of a broadcast recipient.
type BroadcastRecipientStatus string

const (
	BroadcastRecipientStatusPending    BroadcastRecipientStatus = "PENDING"
	BroadcastRecipientStatusSent       BroadcastRecipientStatus = "SENT"
	BroadcastRecipientStatusDelivered  BroadcastRecipientStatus = "DELIVERED"
	BroadcastRecipientStatusOpened     BroadcastRecipientStatus = "OPENED"
	BroadcastRecipientStatusClicked    BroadcastRecipientStatus = "CLICKED"
	BroadcastRecipientStatusBounced    BroadcastRecipientStatus = "BOUNCED"
	BroadcastRecipientStatusComplained BroadcastRecipientStatus = "COMPLAINED"
	BroadcastRecipientStatusFailed     BroadcastRecipientStatus = "FAILED"
)

// BroadcastRecipient represents a recipient of a broadcast.
type BroadcastRecipient struct {
	ID          string                   `json:"id"`
	ContactID   string                   `json:"contactId"`
	Email       string                   `json:"email"`
	Status      BroadcastRecipientStatus `json:"status"`
	SentAt      string                   `json:"sentAt,omitempty"`
	DeliveredAt string                   `json:"deliveredAt,omitempty"`
	OpenedAt    string                   `json:"openedAt,omitempty"`
	ClickedAt   string                   `json:"clickedAt,omitempty"`
	BouncedAt   string                   `json:"bouncedAt,omitempty"`
	FailedAt    string                   `json:"failedAt,omitempty"`
}

// ListBroadcastsOptions represents options for listing broadcasts.
type ListBroadcastsOptions struct {
	Limit  int    `json:"limit,omitempty"`
	Offset int    `json:"offset,omitempty"`
	Cursor string `json:"cursor,omitempty"`
	Status string `json:"status,omitempty"`
}

// ListRecipientsOptions represents options for listing broadcast recipients.
type ListRecipientsOptions struct {
	Limit  int    `json:"limit,omitempty"`
	Offset int    `json:"offset,omitempty"`
	Cursor string `json:"cursor,omitempty"`
	Status string `json:"status,omitempty"`
}

// OpensOverTime represents opens data for a time period.
type OpensOverTime struct {
	Date   string `json:"date"`
	Opens  int    `json:"opens"`
	Unique int    `json:"unique"`
}

// LinkPerformance represents click data for a link.
type LinkPerformance struct {
	URL    string `json:"url"`
	Clicks int    `json:"clicks"`
	Unique int    `json:"unique"`
}

// BroadcastAnalytics represents detailed broadcast analytics.
type BroadcastAnalytics struct {
	OpensOverTime   []OpensOverTime   `json:"opensOverTime"`
	LinkPerformance []LinkPerformance `json:"linkPerformance"`
}

// ─── Events ───

type Event struct {
	ID         string                 `json:"id"`
	ContactID  string                 `json:"contactId"`
	Name       string                 `json:"name"`
	Properties map[string]interface{} `json:"properties,omitempty"`
	CreatedAt  string                 `json:"createdAt"`
}

type FireEventRequest struct {
	Email             string                 `json:"email"`
	EventName         string                 `json:"eventName"`
	Properties        map[string]interface{} `json:"properties,omitempty"`
	ContactProperties map[string]interface{} `json:"contactProperties,omitempty"`
	Tags              []string               `json:"tags,omitempty"`
	IdempotencyKey    string                 `json:"idempotencyKey,omitempty"`
}

type FireEventResponse struct {
	EventID            string                   `json:"eventId"`
	ContactID          string                   `json:"contactId"`
	ContactCreated     bool                     `json:"contactCreated"`
	SequencesTriggered []TriggeredSequence      `json:"sequencesTriggered"`
}

type TriggeredSequence struct {
	SequenceID   string `json:"sequenceId"`
	EnrollmentID string `json:"enrollmentId"`
}

type ListEventsOptions struct {
	Name      string
	ContactID string
	Since     string
	Until     string
	Limit     int
	Offset    int
}

// ─── Sequences ───

type SequenceStatus = string

const (
	SequenceStatusDraft    SequenceStatus = "DRAFT"
	SequenceStatusActive   SequenceStatus = "ACTIVE"
	SequenceStatusPaused   SequenceStatus = "PAUSED"
	SequenceStatusArchived SequenceStatus = "ARCHIVED"
)

type SequenceTriggerType = string

const (
	TriggerTypeEvent          SequenceTriggerType = "EVENT"
	TriggerTypeContactCreated SequenceTriggerType = "CONTACT_CREATED"
	TriggerTypeTagAdded       SequenceTriggerType = "TAG_ADDED"
	TriggerTypeAPICall        SequenceTriggerType = "API_CALL"
	TriggerTypeDateProperty   SequenceTriggerType = "DATE_PROPERTY"
)

type SequenceStepType = string

const (
	StepTypeSendEmail     SequenceStepType = "SEND_EMAIL"
	StepTypeWait          SequenceStepType = "WAIT"
	StepTypeBranch        SequenceStepType = "BRANCH"
	StepTypeUpdateContact SequenceStepType = "UPDATE_CONTACT"
	StepTypeWebhook       SequenceStepType = "WEBHOOK"
)

type EnrollmentStatus = string

const (
	EnrollmentStatusActive    EnrollmentStatus = "ACTIVE"
	EnrollmentStatusPaused    EnrollmentStatus = "PAUSED"
	EnrollmentStatusCompleted EnrollmentStatus = "COMPLETED"
	EnrollmentStatusExited    EnrollmentStatus = "EXITED"
	EnrollmentStatusFailed    EnrollmentStatus = "FAILED"
)

type SequenceStep struct {
	ID                string                 `json:"id"`
	SequenceID        string                 `json:"sequenceId"`
	SequenceVersion   int                    `json:"sequenceVersion"`
	Type              SequenceStepType       `json:"type"`
	Config            map[string]interface{} `json:"config"`
	Position          int                    `json:"position"`
	NextStepID        *string                `json:"nextStepId"`
	BranchTrueStepID  *string                `json:"branchTrueStepId"`
	BranchFalseStepID *string                `json:"branchFalseStepId"`
	EnteredCount      int                    `json:"enteredCount"`
	CompletedCount    int                    `json:"completedCount"`
	CreatedAt         string                 `json:"createdAt"`
}

type Sequence struct {
	ID             string                 `json:"id"`
	OrganizationID string                 `json:"organizationId"`
	Name           string                 `json:"name"`
	Description    *string                `json:"description"`
	Status         SequenceStatus         `json:"status"`
	TriggerType    SequenceTriggerType    `json:"triggerType"`
	TriggerConfig  map[string]interface{} `json:"triggerConfig"`
	TriggerFilters map[string]interface{} `json:"triggerFilters"`
	ReentryAllowed bool                   `json:"reentryAllowed"`
	ReentryDelay   *int                   `json:"reentryDelay"`
	EntryLimit     *int                   `json:"entryLimit"`
	Version        int                    `json:"version"`
	TotalEnrolled  int                    `json:"totalEnrolled"`
	ActiveEnrolled int                    `json:"activeEnrolled"`
	CompletedCount int                    `json:"completedCount"`
	ExitedCount    int                    `json:"exitedCount"`
	CreatedAt      string                 `json:"createdAt"`
	UpdatedAt      string                 `json:"updatedAt"`
}

type SequenceWithSteps struct {
	Sequence
	Steps []SequenceStep `json:"steps"`
}

type CreateSequenceRequest struct {
	Name           string                 `json:"name"`
	Description    string                 `json:"description,omitempty"`
	TriggerType    SequenceTriggerType    `json:"triggerType"`
	TriggerConfig  map[string]interface{} `json:"triggerConfig,omitempty"`
	TriggerFilters map[string]interface{} `json:"triggerFilters,omitempty"`
	ReentryAllowed bool                   `json:"reentryAllowed,omitempty"`
	ReentryDelay   *int                   `json:"reentryDelay,omitempty"`
	EntryLimit     *int                   `json:"entryLimit,omitempty"`
}

type UpdateSequenceRequest struct {
	Name           *string                `json:"name,omitempty"`
	Description    *string                `json:"description,omitempty"`
	TriggerConfig  map[string]interface{} `json:"triggerConfig,omitempty"`
	TriggerFilters map[string]interface{} `json:"triggerFilters,omitempty"`
	ReentryAllowed *bool                  `json:"reentryAllowed,omitempty"`
	ReentryDelay   *int                   `json:"reentryDelay,omitempty"`
	EntryLimit     *int                   `json:"entryLimit,omitempty"`
}

type AddStepRequest struct {
	Type              SequenceStepType       `json:"type"`
	Config            map[string]interface{} `json:"config"`
	Position          *int                   `json:"position,omitempty"`
	NextStepID        string                 `json:"nextStepId,omitempty"`
	BranchTrueStepID  string                 `json:"branchTrueStepId,omitempty"`
	BranchFalseStepID string                 `json:"branchFalseStepId,omitempty"`
}

type EnrollContact struct {
	Email     string `json:"email,omitempty"`
	ContactID string `json:"contactId,omitempty"`
}

type EnrollRequest struct {
	Contacts    []EnrollContact        `json:"contacts"`
	TriggerData map[string]interface{} `json:"triggerData,omitempty"`
}

type EnrollResultItem struct {
	EnrollmentID *string `json:"enrollmentId"`
	ContactID    string  `json:"contactId"`
	Status       string  `json:"status"`
	Skipped      bool    `json:"skipped"`
	Reason       string  `json:"reason,omitempty"`
}

type EnrollResult struct {
	Enrollments []EnrollResultItem `json:"enrollments"`
}

type SequenceEnrollment struct {
	ID              string                 `json:"id"`
	SequenceID      string                 `json:"sequenceId"`
	ContactID       string                 `json:"contactId"`
	ContactEmail    string                 `json:"contactEmail,omitempty"`
	SequenceVersion int                    `json:"sequenceVersion"`
	Status          EnrollmentStatus       `json:"status"`
	CurrentStepID   *string                `json:"currentStepId"`
	ExitReason      *string                `json:"exitReason"`
	EnteredAt       string                 `json:"enteredAt"`
	CompletedAt     *string                `json:"completedAt"`
	ExitedAt        *string                `json:"exitedAt"`
	NextActionAt    *string                `json:"nextActionAt"`
	TriggerData     map[string]interface{} `json:"triggerData"`
}

type ListSequencesOptions struct {
	Status SequenceStatus
	Limit  int
	Offset int
}

type ListEnrollmentsOptions struct {
	Status EnrollmentStatus
	Limit  int
	Offset int
}

type EmailStepStats struct {
	Sent      int     `json:"sent"`
	Opened    int     `json:"opened"`
	Clicked   int     `json:"clicked"`
	Bounced   int     `json:"bounced"`
	OpenRate  float64 `json:"openRate"`
	ClickRate float64 `json:"clickRate"`
}

type StepAnalytics struct {
	StepID         string          `json:"stepId"`
	Type           string          `json:"type"`
	Position       int             `json:"position"`
	EnteredCount   int             `json:"enteredCount"`
	CompletedCount int             `json:"completedCount"`
	DropOffRate    float64         `json:"dropOffRate"`
	EmailStats     *EmailStepStats `json:"emailStats"`
}

type SequenceAnalytics struct {
	TotalEnrolled  int             `json:"totalEnrolled"`
	ActiveEnrolled int             `json:"activeEnrolled"`
	CompletedCount int             `json:"completedCount"`
	ExitedCount    int             `json:"exitedCount"`
	CompletionRate float64         `json:"completionRate"`
	Steps          []StepAnalytics `json:"steps"`
}
