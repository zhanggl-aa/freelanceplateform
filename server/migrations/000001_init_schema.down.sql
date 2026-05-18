-- Drop all tables in reverse dependency order
DROP TABLE IF EXISTS admin_roles;
DROP TABLE IF EXISTS disputes;
DROP TABLE IF EXISTS bookmarks;
DROP TABLE IF EXISTS file_attachments;
DROP TABLE IF EXISTS notification_settings;
DROP TABLE IF EXISTS notifications;
DROP TABLE IF EXISTS reviews;
DROP TABLE IF EXISTS chat_messages;
DROP TABLE IF EXISTS conversation_participants;
DROP TABLE IF EXISTS chat_conversations;
DROP TABLE IF EXISTS wallet_transactions;
DROP TABLE IF EXISTS platform_wallets;
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS contracts;
DROP TABLE IF EXISTS bids;
DROP TABLE IF EXISTS project_milestones;
DROP TABLE IF EXISTS projects;
DROP TABLE IF EXISTS project_categories;
DROP TABLE IF EXISTS client_profiles;
DROP TABLE IF EXISTS developer_portfolio;
DROP TABLE IF EXISTS developer_skills;
DROP TABLE IF EXISTS developer_profiles;
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS oauth_accounts;
DROP TABLE IF EXISTS users;

-- Drop trigger function
DROP FUNCTION IF EXISTS update_updated_at_column();
