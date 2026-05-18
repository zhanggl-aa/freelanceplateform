-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Auto-update updated_at trigger
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- ═══════════════════════════════════════════
-- USERS
-- ═══════════════════════════════════════════
CREATE TABLE users (
    id              uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    email           varchar(255) UNIQUE,
    phone           varchar(20)  UNIQUE,
    password_hash   varchar(255) NOT NULL,
    wechat_openid   varchar(128) UNIQUE,
    wechat_unionid  varchar(128),
    avatar_url      varchar(512),
    nickname        varchar(100) NOT NULL,
    user_type       varchar(20)  NOT NULL CHECK (user_type IN ('client', 'developer', 'both')),
    status          varchar(20)  NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'suspended', 'deleted')),
    email_verified  boolean      NOT NULL DEFAULT false,
    phone_verified  boolean      NOT NULL DEFAULT false,
    last_login_at   timestamptz,
    created_at      timestamptz  NOT NULL DEFAULT now(),
    updated_at      timestamptz  NOT NULL DEFAULT now(),
    CONSTRAINT users_email_or_phone CHECK (email IS NOT NULL OR phone IS NOT NULL)
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_phone ON users(phone);
CREATE INDEX idx_users_wechat_openid ON users(wechat_openid);
CREATE INDEX idx_users_status ON users(status);

CREATE TRIGGER trg_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ═══════════════════════════════════════════
-- OAUTH ACCOUNTS
-- ═══════════════════════════════════════════
CREATE TABLE oauth_accounts (
    id              uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         uuid        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider        varchar(30) NOT NULL CHECK (provider IN ('github', 'google', 'wechat')),
    provider_id     varchar(255) NOT NULL,
    provider_data   jsonb,
    access_token    varchar(512),
    refresh_token   varchar(512),
    token_expires_at timestamptz,
    created_at      timestamptz  NOT NULL DEFAULT now(),
    updated_at      timestamptz  NOT NULL DEFAULT now(),
    CONSTRAINT uq_oauth_provider_id UNIQUE (provider, provider_id)
);

CREATE TRIGGER trg_oauth_accounts_updated_at BEFORE UPDATE ON oauth_accounts
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ═══════════════════════════════════════════
-- REFRESH TOKENS
-- ═══════════════════════════════════════════
CREATE TABLE refresh_tokens (
    id          uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     uuid        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token       varchar(512) NOT NULL UNIQUE,
    device_info varchar(255),
    ip_address  inet,
    expires_at  timestamptz  NOT NULL,
    revoked     boolean      NOT NULL DEFAULT false,
    created_at  timestamptz  NOT NULL DEFAULT now()
);

CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_token ON refresh_tokens(token);

-- ═══════════════════════════════════════════
-- DEVELOPER PROFILES
-- ═══════════════════════════════════════════
CREATE TABLE developer_profiles (
    id                 uuid         PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id            uuid         NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    real_name          varchar(100),
    title              varchar(100),
    bio                text,
    hourly_rate        decimal(10,2),
    availability       varchar(20)  NOT NULL DEFAULT 'available' CHECK (availability IN ('available', 'busy', 'unavailable')),
    experience_years   smallint     NOT NULL DEFAULT 0,
    location           varchar(100),
    github_url         varchar(512),
    linkedin_url       varchar(512),
    website_url        varchar(512),
    verified           boolean      NOT NULL DEFAULT false,
    verified_at        timestamptz,
    rating_avg         decimal(3,2) NOT NULL DEFAULT 0.00,
    rating_count       int          NOT NULL DEFAULT 0,
    total_earnings     decimal(12,2) NOT NULL DEFAULT 0.00,
    completed_projects int          NOT NULL DEFAULT 0,
    created_at         timestamptz  NOT NULL DEFAULT now(),
    updated_at         timestamptz  NOT NULL DEFAULT now()
);

CREATE INDEX idx_developer_profiles_availability ON developer_profiles(availability);
CREATE INDEX idx_developer_profiles_hourly_rate ON developer_profiles(hourly_rate);
CREATE INDEX idx_developer_profiles_rating_avg ON developer_profiles(rating_avg);

CREATE TRIGGER trg_developer_profiles_updated_at BEFORE UPDATE ON developer_profiles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ═══════════════════════════════════════════
-- DEVELOPER SKILLS
-- ═══════════════════════════════════════════
CREATE TABLE developer_skills (
    id               uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    developer_id     uuid        NOT NULL REFERENCES developer_profiles(id) ON DELETE CASCADE,
    skill_name       varchar(80) NOT NULL,
    proficiency      varchar(20) NOT NULL CHECK (proficiency IN ('beginner', 'intermediate', 'advanced', 'expert')),
    years_experience smallint,
    created_at       timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT uq_developer_skill UNIQUE (developer_id, skill_name)
);

CREATE INDEX idx_developer_skills_skill_name ON developer_skills(skill_name);

-- ═══════════════════════════════════════════
-- DEVELOPER PORTFOLIO
-- ═══════════════════════════════════════════
CREATE TABLE developer_portfolio (
    id           uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    developer_id uuid        NOT NULL REFERENCES developer_profiles(id) ON DELETE CASCADE,
    title        varchar(200) NOT NULL,
    description  text,
    project_url  varchar(512),
    image_urls   jsonb       NOT NULL DEFAULT '[]',
    tech_stack   jsonb       NOT NULL DEFAULT '[]',
    start_date   date,
    end_date     date,
    sort_order   smallint    NOT NULL DEFAULT 0,
    created_at   timestamptz NOT NULL DEFAULT now(),
    updated_at   timestamptz NOT NULL DEFAULT now()
);

CREATE TRIGGER trg_developer_portfolio_updated_at BEFORE UPDATE ON developer_portfolio
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ═══════════════════════════════════════════
-- CLIENT PROFILES
-- ═══════════════════════════════════════════
CREATE TABLE client_profiles (
    id               uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id          uuid        NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    company_name     varchar(200),
    company_logo_url varchar(512),
    company_website  varchar(512),
    industry         varchar(100),
    company_size     varchar(20) CHECK (company_size IN ('1-10', '11-50', '51-200', '201-500', '500+')),
    business_license varchar(255),
    verified         boolean     NOT NULL DEFAULT false,
    verified_at      timestamptz,
    total_spent      decimal(12,2) NOT NULL DEFAULT 0.00,
    posted_projects  int         NOT NULL DEFAULT 0,
    created_at       timestamptz NOT NULL DEFAULT now(),
    updated_at       timestamptz NOT NULL DEFAULT now()
);

CREATE TRIGGER trg_client_profiles_updated_at BEFORE UPDATE ON client_profiles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ═══════════════════════════════════════════
-- PROJECT CATEGORIES
-- ═══════════════════════════════════════════
CREATE TABLE project_categories (
    id          uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    name        varchar(100) NOT NULL,
    slug        varchar(100) NOT NULL UNIQUE,
    description text,
    icon_url    varchar(512),
    parent_id   uuid REFERENCES project_categories(id) ON DELETE SET NULL,
    sort_order  smallint    NOT NULL DEFAULT 0,
    is_active   boolean     NOT NULL DEFAULT true,
    created_at  timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_project_categories_parent_id ON project_categories(parent_id);
CREATE INDEX idx_project_categories_slug ON project_categories(slug);

-- ═══════════════════════════════════════════
-- PROJECTS
-- ═══════════════════════════════════════════
CREATE TABLE projects (
    id                    uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id             uuid        NOT NULL REFERENCES users(id),
    category_id           uuid        NOT NULL REFERENCES project_categories(id),
    title                 varchar(200) NOT NULL,
    description           text        NOT NULL,
    budget_min            decimal(12,2),
    budget_max            decimal(12,2),
    budget_type           varchar(20) NOT NULL CHECK (budget_type IN ('fixed', 'hourly')),
    deadline              date,
    tech_stack            jsonb       NOT NULL DEFAULT '[]',
    status                varchar(20) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'published', 'bidding', 'in_progress', 'delivered', 'completed', 'closed', 'cancelled')),
    view_count            int         NOT NULL DEFAULT 0,
    bookmark_count        int         NOT NULL DEFAULT 0,
    bid_count             int         NOT NULL DEFAULT 0,
    bid_deadline          timestamptz,
    assigned_developer_id uuid REFERENCES users(id),
    featured              boolean     NOT NULL DEFAULT false,
    created_at            timestamptz NOT NULL DEFAULT now(),
    updated_at            timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_projects_client_id ON projects(client_id);
CREATE INDEX idx_projects_category_id ON projects(category_id);
CREATE INDEX idx_projects_status ON projects(status);
CREATE INDEX idx_projects_created_at ON projects(created_at);
CREATE INDEX idx_projects_budget_min ON projects(budget_min);

CREATE TRIGGER trg_projects_updated_at BEFORE UPDATE ON projects
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ═══════════════════════════════════════════
-- PROJECT MILESTONES
-- ═══════════════════════════════════════════
CREATE TABLE project_milestones (
    id               uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id       uuid        NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    title            varchar(200) NOT NULL,
    description      text,
    amount           decimal(12,2) NOT NULL CHECK (amount > 0),
    deadline         date,
    status           varchar(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'submitted', 'approved', 'rejected', 'disputed')),
    sort_order       smallint    NOT NULL DEFAULT 0,
    deliverable_urls jsonb       NOT NULL DEFAULT '[]',
    client_feedback  text,
    submitted_at     timestamptz,
    approved_at      timestamptz,
    created_at       timestamptz NOT NULL DEFAULT now(),
    updated_at       timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_project_milestones_project_id ON project_milestones(project_id);
CREATE INDEX idx_project_milestones_status ON project_milestones(status);

CREATE TRIGGER trg_project_milestones_updated_at BEFORE UPDATE ON project_milestones
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ═══════════════════════════════════════════
-- BIDS
-- ═══════════════════════════════════════════
CREATE TABLE bids (
    id              uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id      uuid        NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    developer_id    uuid        NOT NULL REFERENCES users(id),
    cover_letter    text        NOT NULL,
    estimated_days  smallint    NOT NULL CHECK (estimated_days > 0),
    proposed_budget decimal(12,2) NOT NULL CHECK (proposed_budget > 0),
    budget_type     varchar(20) NOT NULL CHECK (budget_type IN ('fixed', 'hourly')),
    milestone_plan  jsonb,
    status          varchar(20) NOT NULL DEFAULT 'submitted' CHECK (status IN ('submitted', 'shortlisted', 'interview', 'accepted', 'rejected', 'withdrawn', 'counter_offered')),
    client_message  text,
    created_at      timestamptz NOT NULL DEFAULT now(),
    updated_at      timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT uq_bid_project_developer UNIQUE (project_id, developer_id)
);

CREATE INDEX idx_bids_project_id ON bids(project_id);
CREATE INDEX idx_bids_developer_id ON bids(developer_id);
CREATE INDEX idx_bids_status ON bids(status);

CREATE TRIGGER trg_bids_updated_at BEFORE UPDATE ON bids
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ═══════════════════════════════════════════
-- CONTRACTS
-- ═══════════════════════════════════════════
CREATE TABLE contracts (
    id                uuid         PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id        uuid         NOT NULL UNIQUE REFERENCES projects(id),
    client_id         uuid         NOT NULL REFERENCES users(id),
    developer_id      uuid         NOT NULL REFERENCES users(id),
    bid_id            uuid         NOT NULL REFERENCES bids(id),
    total_amount      decimal(12,2) NOT NULL CHECK (total_amount > 0),
    platform_fee_rate decimal(5,4)  NOT NULL DEFAULT 0.1000,
    platform_fee      decimal(12,2) NOT NULL,
    developer_payout  decimal(12,2) NOT NULL,
    paid_amount       decimal(12,2) NOT NULL DEFAULT 0.00,
    released_amount   decimal(12,2) NOT NULL DEFAULT 0.00,
    status            varchar(20)   NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'completed', 'cancelled', 'disputed')),
    start_date        date          NOT NULL,
    end_date          date,
    terms             text,
    created_at        timestamptz   NOT NULL DEFAULT now(),
    updated_at        timestamptz   NOT NULL DEFAULT now()
);

CREATE INDEX idx_contracts_client_id ON contracts(client_id);
CREATE INDEX idx_contracts_developer_id ON contracts(developer_id);
CREATE INDEX idx_contracts_status ON contracts(status);

CREATE TRIGGER trg_contracts_updated_at BEFORE UPDATE ON contracts
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ═══════════════════════════════════════════
-- PAYMENTS
-- ═══════════════════════════════════════════
CREATE TABLE payments (
    id             uuid         PRIMARY KEY DEFAULT gen_random_uuid(),
    contract_id    uuid         NOT NULL REFERENCES contracts(id),
    milestone_id   uuid         REFERENCES project_milestones(id),
    payer_id       uuid         NOT NULL REFERENCES users(id),
    payee_id       uuid         NOT NULL REFERENCES users(id),
    amount         decimal(12,2) NOT NULL CHECK (amount > 0),
    platform_fee   decimal(12,2) NOT NULL DEFAULT 0.00,
    net_amount     decimal(12,2) NOT NULL,
    payment_method varchar(30)  NOT NULL CHECK (payment_method IN ('wechat_pay', 'alipay', 'bank_transfer', 'balance')),
    external_tx_id varchar(255),
    status         varchar(20)  NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'escrow', 'released', 'refunding', 'refunded', 'failed')),
    escrow_at      timestamptz,
    released_at    timestamptz,
    refunded_at    timestamptz,
    created_at     timestamptz  NOT NULL DEFAULT now(),
    updated_at     timestamptz  NOT NULL DEFAULT now()
);

CREATE INDEX idx_payments_contract_id ON payments(contract_id);
CREATE INDEX idx_payments_payer_id ON payments(payer_id);
CREATE INDEX idx_payments_payee_id ON payments(payee_id);
CREATE INDEX idx_payments_status ON payments(status);

CREATE TRIGGER trg_payments_updated_at BEFORE UPDATE ON payments
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ═══════════════════════════════════════════
-- PLATFORM WALLETS
-- ═══════════════════════════════════════════
CREATE TABLE platform_wallets (
    id              uuid         PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         uuid         NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    balance         decimal(12,2) NOT NULL DEFAULT 0.00 CHECK (balance >= 0),
    frozen_amount   decimal(12,2) NOT NULL DEFAULT 0.00,
    total_deposited decimal(12,2) NOT NULL DEFAULT 0.00,
    total_withdrawn decimal(12,2) NOT NULL DEFAULT 0.00,
    created_at      timestamptz  NOT NULL DEFAULT now(),
    updated_at      timestamptz  NOT NULL DEFAULT now()
);

CREATE TRIGGER trg_platform_wallets_updated_at BEFORE UPDATE ON platform_wallets
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ═══════════════════════════════════════════
-- WALLET TRANSACTIONS
-- ═══════════════════════════════════════════
CREATE TABLE wallet_transactions (
    id            uuid         PRIMARY KEY DEFAULT gen_random_uuid(),
    wallet_id     uuid         NOT NULL REFERENCES platform_wallets(id),
    payment_id    uuid         REFERENCES payments(id),
    type          varchar(30)  NOT NULL CHECK (type IN ('deposit', 'withdraw', 'escrow_hold', 'escrow_release', 'escrow_refund', 'commission')),
    amount        decimal(12,2) NOT NULL,
    balance_after decimal(12,2) NOT NULL,
    description   varchar(255),
    created_at    timestamptz  NOT NULL DEFAULT now()
);

CREATE INDEX idx_wallet_transactions_wallet_id ON wallet_transactions(wallet_id);
CREATE INDEX idx_wallet_transactions_created_at ON wallet_transactions(created_at);

-- ═══════════════════════════════════════════
-- CHAT CONVERSATIONS
-- ═══════════════════════════════════════════
CREATE TABLE chat_conversations (
    id             uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    type           varchar(20) NOT NULL CHECK (type IN ('direct', 'project')),
    project_id     uuid REFERENCES projects(id) ON DELETE SET NULL,
    last_message_at timestamptz NOT NULL DEFAULT now(),
    created_at     timestamptz NOT NULL DEFAULT now(),
    updated_at     timestamptz NOT NULL DEFAULT now()
);

CREATE TRIGGER trg_chat_conversations_updated_at BEFORE UPDATE ON chat_conversations
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ═══════════════════════════════════════════
-- CONVERSATION PARTICIPANTS
-- ═══════════════════════════════════════════
CREATE TABLE conversation_participants (
    id              uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    conversation_id uuid        NOT NULL REFERENCES chat_conversations(id) ON DELETE CASCADE,
    user_id         uuid        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    last_read_at    timestamptz NOT NULL DEFAULT '1970-01-01',
    is_muted        boolean     NOT NULL DEFAULT false,
    created_at      timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT uq_conversation_participant UNIQUE (conversation_id, user_id)
);

CREATE INDEX idx_conversation_participants_user_id ON conversation_participants(user_id);

-- ═══════════════════════════════════════════
-- CHAT MESSAGES
-- ═══════════════════════════════════════════
CREATE TABLE chat_messages (
    id              uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    conversation_id uuid        NOT NULL REFERENCES chat_conversations(id) ON DELETE CASCADE,
    sender_id       uuid        NOT NULL REFERENCES users(id),
    content         text,
    message_type    varchar(20) NOT NULL CHECK (message_type IN ('text', 'image', 'file', 'system')),
    file_url        varchar(512),
    file_name       varchar(255),
    file_size       bigint,
    is_read         boolean     NOT NULL DEFAULT false,
    created_at      timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_chat_messages_conversation_id_created_at ON chat_messages(conversation_id, created_at);
CREATE INDEX idx_chat_messages_sender_id ON chat_messages(sender_id);

-- ═══════════════════════════════════════════
-- REVIEWS
-- ═══════════════════════════════════════════
CREATE TABLE reviews (
    id                  uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id          uuid        NOT NULL REFERENCES projects(id),
    contract_id         uuid        NOT NULL REFERENCES contracts(id),
    reviewer_id         uuid        NOT NULL REFERENCES users(id),
    reviewee_id         uuid        NOT NULL REFERENCES users(id),
    quality_rating      smallint    NOT NULL CHECK (quality_rating BETWEEN 1 AND 5),
    communication_rating smallint   NOT NULL CHECK (communication_rating BETWEEN 1 AND 5),
    timeliness_rating   smallint    NOT NULL CHECK (timeliness_rating BETWEEN 1 AND 5),
    overall_rating      decimal(3,2) NOT NULL,
    comment             text,
    is_public           boolean     NOT NULL DEFAULT true,
    created_at          timestamptz NOT NULL DEFAULT now(),
    updated_at          timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT uq_review_per_project UNIQUE (project_id, reviewer_id),
    CONSTRAINT chk_reviewer_not_reviewee CHECK (reviewer_id != reviewee_id)
);

CREATE INDEX idx_reviews_reviewee_id ON reviews(reviewee_id);
CREATE INDEX idx_reviews_reviewer_id ON reviews(reviewer_id);

CREATE TRIGGER trg_reviews_updated_at BEFORE UPDATE ON reviews
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ═══════════════════════════════════════════
-- NOTIFICATIONS
-- ═══════════════════════════════════════════
CREATE TABLE notifications (
    id         uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    uuid        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type       varchar(50) NOT NULL,
    title      varchar(200) NOT NULL,
    content    text,
    data       jsonb,
    is_read    boolean     NOT NULL DEFAULT false,
    read_at    timestamptz,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_notifications_user_id_is_read ON notifications(user_id, is_read);
CREATE INDEX idx_notifications_created_at ON notifications(created_at);

-- ═══════════════════════════════════════════
-- NOTIFICATION SETTINGS
-- ═══════════════════════════════════════════
CREATE TABLE notification_settings (
    id                     uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id                uuid        NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    email_enabled          boolean     NOT NULL DEFAULT true,
    sms_enabled            boolean     NOT NULL DEFAULT false,
    push_enabled           boolean     NOT NULL DEFAULT true,
    in_app_enabled         boolean     NOT NULL DEFAULT true,
    bid_notifications      boolean     NOT NULL DEFAULT true,
    message_notifications  boolean     NOT NULL DEFAULT true,
    payment_notifications  boolean     NOT NULL DEFAULT true,
    project_notifications  boolean     NOT NULL DEFAULT true,
    created_at             timestamptz NOT NULL DEFAULT now(),
    updated_at             timestamptz NOT NULL DEFAULT now()
);

CREATE TRIGGER trg_notification_settings_updated_at BEFORE UPDATE ON notification_settings
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ═══════════════════════════════════════════
-- FILE ATTACHMENTS
-- ═══════════════════════════════════════════
CREATE TABLE file_attachments (
    id           uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      uuid        NOT NULL REFERENCES users(id),
    file_name    varchar(255) NOT NULL,
    file_path    varchar(512) NOT NULL,
    file_size    bigint      NOT NULL,
    file_type    varchar(50) NOT NULL,
    mime_type    varchar(100) NOT NULL,
    storage_type varchar(20) NOT NULL DEFAULT 'local' CHECK (storage_type IN ('local', 'oss')),
    entity_type  varchar(50),
    entity_id    uuid,
    created_at   timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_file_attachments_entity ON file_attachments(entity_type, entity_id);
CREATE INDEX idx_file_attachments_user_id ON file_attachments(user_id);

-- ═══════════════════════════════════════════
-- BOOKMARKS
-- ═══════════════════════════════════════════
CREATE TABLE bookmarks (
    id         uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    uuid        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    project_id uuid        NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    created_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT uq_bookmark UNIQUE (user_id, project_id)
);

-- ═══════════════════════════════════════════
-- DISPUTES
-- ═══════════════════════════════════════════
CREATE TABLE disputes (
    id              uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    contract_id     uuid        NOT NULL REFERENCES contracts(id),
    milestone_id    uuid        REFERENCES project_milestones(id),
    reporter_id     uuid        NOT NULL REFERENCES users(id),
    reported_id     uuid        NOT NULL REFERENCES users(id),
    reason          text        NOT NULL,
    evidence_urls   jsonb       NOT NULL DEFAULT '[]',
    status          varchar(20) NOT NULL DEFAULT 'open' CHECK (status IN ('open', 'under_review', 'resolved', 'closed')),
    resolution      text,
    resolution_type varchar(30) CHECK (resolution_type IN ('refund_full', 'refund_partial', 'release_full', 'release_partial', 'other')),
    resolved_by     uuid REFERENCES users(id),
    resolved_at     timestamptz,
    created_at      timestamptz NOT NULL DEFAULT now(),
    updated_at      timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_disputes_contract_id ON disputes(contract_id);
CREATE INDEX idx_disputes_status ON disputes(status);

CREATE TRIGGER trg_disputes_updated_at BEFORE UPDATE ON disputes
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ═══════════════════════════════════════════
-- ADMIN ROLES
-- ═══════════════════════════════════════════
CREATE TABLE admin_roles (
    id          uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     uuid        NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    role        varchar(30) NOT NULL CHECK (role IN ('super_admin', 'admin', 'moderator', 'finance')),
    permissions jsonb       NOT NULL DEFAULT '{}',
    created_at  timestamptz NOT NULL DEFAULT now(),
    updated_at  timestamptz NOT NULL DEFAULT now()
);

CREATE TRIGGER trg_admin_roles_updated_at BEFORE UPDATE ON admin_roles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ═══════════════════════════════════════════
-- SEED DATA: Default project categories
-- ═══════════════════════════════════════════
INSERT INTO project_categories (name, slug, sort_order) VALUES
('网站开发', 'web-development', 1),
('移动应用', 'mobile-app', 2),
('小程序开发', 'mini-program', 3),
('前端开发', 'frontend', 4),
('后端开发', 'backend', 5),
('AI/机器学习', 'ai-ml', 6),
('UI/UX设计', 'ui-ux-design', 7),
('数据分析', 'data-analysis', 8),
('DevOps', 'devops', 9),
('游戏开发', 'game-development', 10),
('区块链', 'blockchain', 11),
('其他', 'other', 99);
