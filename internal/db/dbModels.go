package db

import "time"

// Migration represents a migration entry in the database
type Migrations struct {
	ID        int64  `gorm:"primaryKey"`
	Timestamp int64  `gorm:"not null"`
	Name      string `gorm:"unique;not null"`
}

type AuditFields struct {
	Status    int32 `gorm:"not null"`
	CreatedBy int64 `gorm:"not null"`
	CreatedAt int64 `gorm:"not null"` // store Unix seconds
	UpdatedBy int64 `gorm:"not null"`
	UpdatedAt int64 `gorm:"not null"` // store Unix seconds
}

// Subscription
type Subscription struct {
	ID            int64  `gorm:"primaryKey"`
	Name          string `gorm:"size:255;not null"`
	Type          string `gorm:"size:50;not null"`
	UserLimit     int
	PipelineLimit int
	IsCustomPlan  bool `gorm:"default:false"`
	AuditFields
	Tenants []Tenant `gorm:"foreignKey:SubscriptionID"`
}

// Tenant
type Tenant struct {
	ID             int64  `gorm:"primaryKey"`
	Name           string `gorm:"size:255;not null"`
	SubscriptionID int64
	Subscription   Subscription `gorm:"foreignKey:SubscriptionID"`
	AuditFields
	Users      []User           `gorm:"foreignKey:TenantID"`
	Pipelines  []Pipeline       `gorm:"foreignKey:TenantID"`
	Connectors []Connector      `gorm:"foreignKey:TenantID"`
	Transforms []Transformation `gorm:"foreignKey:TenantID"`
	AuditLogs  []AuditLog       `gorm:"foreignKey:TenantID"`
}

// User
type User struct {
	ID           int64  `gorm:"primaryKey"`
	Username     string `gorm:"size:255;not null"`
	Email        string `gorm:"size:255;not null;unique"`
	PasswordHash string `gorm:"type:text;not null"`
	Role         string `gorm:"size:50;not null"`
	TenantID     int64
	Tenant       Tenant `gorm:"foreignKey:TenantID"`
	AuditFields
	AuditLogs []AuditLog `gorm:"foreignKey:UserID"`
}

// Connector
type Connector struct {
	ID               int64 `gorm:"primaryKey"`
	TenantID         int64
	Tenant           Tenant `gorm:"foreignKey:TenantID"`
	Name             string `gorm:"size:255;not null"`
	ConnectionTypeID int64
	ConnectionType   ConnectorType `gorm:"foreignKey:ConnectionTypeID"`
	AuditFields
	Properties []ConnectorProperty `gorm:"foreignKey:ConnectorID"`
}

// ConnectorType
type ConnectorType struct {
	ID          int64  `gorm:"primaryKey"`
	Name        string `gorm:"size:255;not null"`
	Type        int
	Properties  string `gorm:"type:jsonb"`
	DocumentURL string `gorm:"size:500"`
	AuditFields
	Connectors []Connector `gorm:"foreignKey:ConnectionTypeID"`
}

// ConnectorProperty
type ConnectorProperty struct {
	ID            int64 `gorm:"primaryKey"`
	ConnectorID   int64
	Connector     Connector `gorm:"foreignKey:ConnectorID"`
	PropertyKey   string    `gorm:"size:255;not null"`
	PropertyValue string    `gorm:"type:text"`
	IsSensitive   bool      `gorm:"default:false"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}

// Pipeline
type Pipeline struct {
	ID          int64 `gorm:"primaryKey"`
	TenantID    int64
	Tenant      Tenant `gorm:"foreignKey:TenantID"`
	Name        string `gorm:"size:255;not null"`
	Description string `gorm:"type:text"`
	Schedule    string `gorm:"size:255"`
	AuditFields
	Connectors      []PipelineConnector `gorm:"foreignKey:PipelineID"`
	Transformations []Transformation    `gorm:"foreignKey:PipelineID"`
	Runs            []PipelineRun       `gorm:"foreignKey:PipelineID"`
}

// PipelineConnector
type PipelineConnector struct {
	ID          int64 `gorm:"primaryKey"`
	PipelineID  int64
	Pipeline    Pipeline `gorm:"foreignKey:PipelineID"`
	ConnectorID int64
	Connector   Connector `gorm:"foreignKey:ConnectorID"`
	Role        string    `gorm:"size:50;not null"`
	AuditFields
}

// Transformation
type Transformation struct {
	ID          int64 `gorm:"primaryKey"`
	TenantID    int64
	Tenant      Tenant `gorm:"foreignKey:TenantID"`
	PipelineID  int64
	Pipeline    Pipeline `gorm:"foreignKey:PipelineID"`
	Name        string   `gorm:"size:255;not null"`
	Description string   `gorm:"type:text"`
	Language    string   `gorm:"size:50;not null"`
	Code        string   `gorm:"type:text"`
	StepOrder   int
	AuditFields
}

// PipelineRun
type PipelineRun struct {
	ID               int64 `gorm:"primaryKey"`
	PipelineID       int64
	Pipeline         Pipeline `gorm:"foreignKey:PipelineID"`
	Status           string   `gorm:"size:50;not null"`
	StartTime        int64
	EndTime          int64
	RecordsProcessed int64
	ErrorMessage     string `gorm:"type:text"`
	CreatedAt        int64  `gorm:"autoCreateTime"`
	Jobs             []Job  `gorm:"foreignKey:PipelineRunID"`
}

// Job
type Job struct {
	ID            int64 `gorm:"primaryKey"`
	PipelineRunID int64
	PipelineRun   PipelineRun `gorm:"foreignKey:PipelineRunID"`
	StepName      string      `gorm:"size:255"`
	StepType      string      `gorm:"size:255"`
	Status        string      `gorm:"size:50;not null"`
	StartTime     time.Time
	EndTime       time.Time
	Log           string    `gorm:"type:text"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
}

// AuditLog
type AuditLog struct {
	ID        int64 `gorm:"primaryKey"`
	TenantID  int64
	Tenant    Tenant `gorm:"foreignKey:TenantID"`
	UserID    int64
	User      User      `gorm:"foreignKey:UserID"`
	Action    string    `gorm:"size:255;not null"`
	Details   string    `gorm:"type:jsonb"`
	IPAddress string    `gorm:"size:100"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
