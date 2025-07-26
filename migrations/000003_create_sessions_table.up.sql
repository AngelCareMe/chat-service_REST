-- Create sessions table
CREATE TABLE IF NOT EXISTS sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(512) NOT NULL UNIQUE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Add comments
COMMENT ON TABLE sessions IS 'User sessions table for authentication';
COMMENT ON COLUMN sessions.id IS 'Unique identifier for the session';
COMMENT ON COLUMN sessions.user_id IS 'Reference to the user who owns the session';
COMMENT ON COLUMN sessions.token IS 'JWT token for the session';
COMMENT ON COLUMN sessions.expires_at IS 'Timestamp when session expires';
COMMENT ON COLUMN sessions.created_at IS 'Timestamp when session was created';