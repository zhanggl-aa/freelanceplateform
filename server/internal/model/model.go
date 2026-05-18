package model

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWT Claims
type JWTClaims struct {
	UserID   string `json:"user_id"`
	UserType string `json:"user_type"`
	jwt.RegisteredClaims
}

// ─── Request DTOs ───

type RegisterRequest struct {
	Email    *string `json:"email" binding:"omitempty,email"`
	Phone    *string `json:"phone"`
	Password string  `json:"password" binding:"required,min=6"`
	Nickname string  `json:"nickname" binding:"required,min=2,max=50"`
	UserType string  `json:"user_type" binding:"required,oneof=client developer both"`
}

type LoginRequest struct {
	Email    *string `json:"email"`
	Phone    *string `json:"phone"`
	Password string  `json:"password" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type PaginationQuery struct {
	Page     int `form:"page,default=1" binding:"min=1"`
	PageSize int `form:"page_size,default=20" binding:"min=1,max=100"`
}

func (p PaginationQuery) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// ─── User ───

type User struct {
	ID             string     `json:"id"`
	Email          *string    `json:"email,omitempty"`
	Phone          *string    `json:"phone,omitempty"`
	PasswordHash   string     `json:"-"`
	WechatOpenID   *string    `json:"-"`
	WechatUnionID  *string    `json:"-"`
	AvatarURL      *string    `json:"avatar_url,omitempty"`
	Nickname       string     `json:"nickname"`
	UserType       string     `json:"user_type"`
	Status         string     `json:"status"`
	EmailVerified  bool       `json:"email_verified"`
	PhoneVerified  bool       `json:"phone_verified"`
	LastLoginAt    *time.Time `json:"last_login_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// ─── OAuth Account ───

type OAuthAccount struct {
	ID             string     `json:"id"`
	UserID         string     `json:"user_id"`
	Provider       string     `json:"provider"`
	ProviderID     string     `json:"provider_id"`
	ProviderData   *string    `json:"provider_data,omitempty"`
	AccessToken    *string    `json:"-"`
	RefreshToken   *string    `json:"-"`
	TokenExpiresAt *time.Time `json:"token_expires_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// ─── Refresh Token ───

type RefreshToken struct {
	ID         string     `json:"id"`
	UserID     string     `json:"user_id"`
	Token      string     `json:"-"`
	DeviceInfo *string    `json:"device_info,omitempty"`
	IPAddress  *string    `json:"-"`
	ExpiresAt  time.Time  `json:"expires_at"`
	Revoked    bool       `json:"revoked"`
	CreatedAt  time.Time  `json:"created_at"`
}

// ─── Developer Profile ───

type DeveloperProfile struct {
	ID               string     `json:"id"`
	UserID           string     `json:"user_id"`
	RealName         *string    `json:"real_name,omitempty"`
	Title            *string    `json:"title,omitempty"`
	Bio              *string    `json:"bio,omitempty"`
	HourlyRate       *float64   `json:"hourly_rate,omitempty"`
	Availability     string     `json:"availability"`
	ExperienceYears  int        `json:"experience_years"`
	Location         *string    `json:"location,omitempty"`
	GithubURL        *string    `json:"github_url,omitempty"`
	LinkedinURL      *string    `json:"linkedin_url,omitempty"`
	WebsiteURL       *string    `json:"website_url,omitempty"`
	Verified         bool       `json:"verified"`
	VerifiedAt       *time.Time `json:"verified_at,omitempty"`
	RatingAvg        float64    `json:"rating_avg"`
	RatingCount      int        `json:"rating_count"`
	TotalEarnings    float64    `json:"total_earnings"`
	CompletedProjects int       `json:"completed_projects"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// ─── Developer Skill ───

type DeveloperSkill struct {
	ID              string  `json:"id"`
	DeveloperID     string  `json:"developer_id"`
	SkillName       string  `json:"skill_name"`
	Proficiency     string  `json:"proficiency"`
	YearsExperience *int    `json:"years_experience,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

// ─── Developer Portfolio ───

type DeveloperPortfolio struct {
	ID          string     `json:"id"`
	DeveloperID string     `json:"developer_id"`
	Title       string     `json:"title"`
	Description *string    `json:"description,omitempty"`
	ProjectURL  *string    `json:"project_url,omitempty"`
	ImageURLs   []string   `json:"image_urls"`
	TechStack   []string   `json:"tech_stack"`
	StartDate   *time.Time `json:"start_date,omitempty"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	SortOrder   int        `json:"sort_order"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// ─── Client Profile ───

type ClientProfile struct {
	ID             string     `json:"id"`
	UserID         string     `json:"user_id"`
	CompanyName    *string    `json:"company_name,omitempty"`
	CompanyLogoURL *string    `json:"company_logo_url,omitempty"`
	CompanyWebsite *string    `json:"company_website,omitempty"`
	Industry       *string    `json:"industry,omitempty"`
	CompanySize    *string    `json:"company_size,omitempty"`
	Verified       bool       `json:"verified"`
	VerifiedAt     *time.Time `json:"verified_at,omitempty"`
	TotalSpent     float64    `json:"total_spent"`
	PostedProjects int        `json:"posted_projects"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// ─── Project Category ───

type ProjectCategory struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Description *string `json:"description,omitempty"`
	IconURL     *string `json:"icon_url,omitempty"`
	ParentID    *string `json:"parent_id,omitempty"`
	SortOrder   int     `json:"sort_order"`
	IsActive    bool    `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	Children    []ProjectCategory `json:"children,omitempty"`
}

// ─── Project ───

type Project struct {
	ID                  string     `json:"id"`
	ClientID            string     `json:"client_id"`
	CategoryID          string     `json:"category_id"`
	Title               string     `json:"title"`
	Description         string     `json:"description"`
	BudgetMin           *float64   `json:"budget_min,omitempty"`
	BudgetMax           *float64   `json:"budget_max,omitempty"`
	BudgetType          string     `json:"budget_type"`
	Deadline            *time.Time `json:"deadline,omitempty"`
	TechStack           []string   `json:"tech_stack"`
	Status              string     `json:"status"`
	ViewCount           int        `json:"view_count"`
	BookmarkCount       int        `json:"bookmark_count"`
	BidCount            int        `json:"bid_count"`
	BidDeadline         *time.Time `json:"bid_deadline,omitempty"`
	AssignedDeveloperID *string    `json:"assigned_developer_id,omitempty"`
	Featured            bool       `json:"featured"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	ClientName          *string    `json:"client_name,omitempty"`
	ClientAvatar        *string    `json:"client_avatar,omitempty"`
	CategoryName        *string    `json:"category_name,omitempty"`
}

// ─── Project Milestone ───

type ProjectMilestone struct {
	ID             string     `json:"id"`
	ProjectID      string     `json:"project_id"`
	Title          string     `json:"title"`
	Description    *string    `json:"description,omitempty"`
	Amount         float64    `json:"amount"`
	Deadline       *time.Time `json:"deadline,omitempty"`
	Status         string     `json:"status"`
	SortOrder      int        `json:"sort_order"`
	DeliverableURLs []string  `json:"deliverable_urls"`
	ClientFeedback *string    `json:"client_feedback,omitempty"`
	SubmittedAt    *time.Time `json:"submitted_at,omitempty"`
	ApprovedAt     *time.Time `json:"approved_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// ─── Bid ───

type Bid struct {
	ID             string     `json:"id"`
	ProjectID      string     `json:"project_id"`
	DeveloperID    string     `json:"developer_id"`
	CoverLetter    string     `json:"cover_letter"`
	EstimatedDays  int        `json:"estimated_days"`
	ProposedBudget float64    `json:"proposed_budget"`
	BudgetType     string     `json:"budget_type"`
	MilestonePlan  *string    `json:"milestone_plan,omitempty"`
	Status         string     `json:"status"`
	ClientMessage  *string    `json:"client_message,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeveloperName  *string    `json:"developer_name,omitempty"`
	DeveloperAvatar *string   `json:"developer_avatar,omitempty"`
}

// ─── Contract ───

type Contract struct {
	ID             string     `json:"id"`
	ProjectID      string     `json:"project_id"`
	ClientID       string     `json:"client_id"`
	DeveloperID    string     `json:"developer_id"`
	BidID          string     `json:"bid_id"`
	TotalAmount    float64    `json:"total_amount"`
	PlatformFeeRate float64   `json:"platform_fee_rate"`
	PlatformFee    float64    `json:"platform_fee"`
	DeveloperPayout float64   `json:"developer_payout"`
	PaidAmount     float64    `json:"paid_amount"`
	ReleasedAmount float64    `json:"released_amount"`
	Status         string     `json:"status"`
	StartDate      time.Time  `json:"start_date"`
	EndDate        *time.Time `json:"end_date,omitempty"`
	Terms          *string    `json:"terms,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// ─── Payment ───

type Payment struct {
	ID            string     `json:"id"`
	ContractID    string     `json:"contract_id"`
	MilestoneID   *string    `json:"milestone_id,omitempty"`
	PayerID       string     `json:"payer_id"`
	PayeeID       string     `json:"payee_id"`
	Amount        float64    `json:"amount"`
	PlatformFee   float64    `json:"platform_fee"`
	NetAmount     float64    `json:"net_amount"`
	PaymentMethod string     `json:"payment_method"`
	ExternalTxID  *string    `json:"external_tx_id,omitempty"`
	Status        string     `json:"status"`
	EscrowAt      *time.Time `json:"escrow_at,omitempty"`
	ReleasedAt    *time.Time `json:"released_at,omitempty"`
	RefundedAt    *time.Time `json:"refunded_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// ─── Platform Wallet ───

type PlatformWallet struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	Balance        float64   `json:"balance"`
	FrozenAmount   float64   `json:"frozen_amount"`
	TotalDeposited float64   `json:"total_deposited"`
	TotalWithdrawn float64   `json:"total_withdrawn"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ─── Wallet Transaction ───

type WalletTransaction struct {
	ID           string    `json:"id"`
	WalletID     string    `json:"wallet_id"`
	PaymentID    *string   `json:"payment_id,omitempty"`
	Type         string    `json:"type"`
	Amount       float64   `json:"amount"`
	BalanceAfter float64   `json:"balance_after"`
	Description  *string   `json:"description,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

// ─── Chat Conversation ───

type ChatConversation struct {
	ID            string     `json:"id"`
	Type          string     `json:"type"`
	ProjectID     *string    `json:"project_id,omitempty"`
	LastMessageAt time.Time  `json:"last_message_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	Participants  []ConversationParticipant `json:"participants,omitempty"`
	LastMessage   *ChatMessage `json:"last_message,omitempty"`
}

// ─── Conversation Participant ───

type ConversationParticipant struct {
	ID             string    `json:"id"`
	ConversationID string    `json:"conversation_id"`
	UserID         string    `json:"user_id"`
	LastReadAt     time.Time `json:"last_read_at"`
	IsMuted        bool      `json:"is_muted"`
	Nickname       *string   `json:"nickname,omitempty"`
	AvatarURL      *string   `json:"avatar_url,omitempty"`
}

// ─── Chat Message ───

type ChatMessage struct {
	ID             string     `json:"id"`
	ConversationID string     `json:"conversation_id"`
	SenderID       string     `json:"sender_id"`
	Content        *string    `json:"content,omitempty"`
	MessageType    string     `json:"message_type"`
	FileURL        *string    `json:"file_url,omitempty"`
	FileName       *string    `json:"file_name,omitempty"`
	FileSize       *int64     `json:"file_size,omitempty"`
	IsRead         bool       `json:"is_read"`
	CreatedAt      time.Time  `json:"created_at"`
	SenderName     *string    `json:"sender_name,omitempty"`
	SenderAvatar   *string    `json:"sender_avatar,omitempty"`
}

// ─── Review ───

type Review struct {
	ID                 string    `json:"id"`
	ProjectID          string    `json:"project_id"`
	ContractID         string    `json:"contract_id"`
	ReviewerID         string    `json:"reviewer_id"`
	RevieweeID         string    `json:"reviewee_id"`
	QualityRating      int       `json:"quality_rating"`
	CommunicationRating int      `json:"communication_rating"`
	TimelinessRating   int       `json:"timeliness_rating"`
	OverallRating      float64   `json:"overall_rating"`
	Comment            *string   `json:"comment,omitempty"`
	IsPublic           bool      `json:"is_public"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	ReviewerName       *string   `json:"reviewer_name,omitempty"`
	ReviewerAvatar     *string   `json:"reviewer_avatar,omitempty"`
}

// ─── Notification ───

type Notification struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	Type      string     `json:"type"`
	Title     string     `json:"title"`
	Content   *string    `json:"content,omitempty"`
	Data      *string    `json:"data,omitempty"`
	IsRead    bool       `json:"is_read"`
	ReadAt    *time.Time `json:"read_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

// ─── Notification Settings ───

type NotificationSettings struct {
	ID                    string `json:"id"`
	UserID                string `json:"user_id"`
	EmailEnabled          bool   `json:"email_enabled"`
	SMSEnabled            bool   `json:"sms_enabled"`
	PushEnabled           bool   `json:"push_enabled"`
	InAppEnabled          bool   `json:"in_app_enabled"`
	BidNotifications      bool   `json:"bid_notifications"`
	MessageNotifications  bool   `json:"message_notifications"`
	PaymentNotifications  bool   `json:"payment_notifications"`
	ProjectNotifications  bool   `json:"project_notifications"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// ─── File Attachment ───

type FileAttachment struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	FileName   string    `json:"file_name"`
	FilePath   string    `json:"file_path"`
	FileSize   int64     `json:"file_size"`
	FileType   string    `json:"file_type"`
	MimeType   string    `json:"mime_type"`
	StorageType string  `json:"storage_type"`
EntityType *string  `json:"entity_type,omitempty"`
	EntityID   *string  `json:"entity_id,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

// ─── Bookmark ───

type Bookmark struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	ProjectID string    `json:"project_id"`
	CreatedAt time.Time `json:"created_at"`
}

// ─── Dispute ───

type Dispute struct {
	ID             string     `json:"id"`
	ContractID     string     `json:"contract_id"`
	MilestoneID    *string    `json:"milestone_id,omitempty"`
	ReporterID     string     `json:"reporter_id"`
	ReportedID     string     `json:"reported_id"`
	Reason         string     `json:"reason"`
	EvidenceURLs   []string   `json:"evidence_urls"`
	Status         string     `json:"status"`
	Resolution     *string    `json:"resolution,omitempty"`
	ResolutionType *string    `json:"resolution_type,omitempty"`
	ResolvedBy     *string    `json:"resolved_by,omitempty"`
	ResolvedAt     *time.Time `json:"resolved_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// ─── Admin Role ───

type AdminRole struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Role        string    `json:"role"`
	Permissions *string   `json:"permissions,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
