-- Create messages table
CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Add comments
COMMENT ON TABLE messages IS 'Messages table for chat application';
COMMENT ON COLUMN messages.id IS 'Unique identifier for the message';
COMMENT ON COLUMN messages.user_id IS 'Reference to the user who sent the message';
COMMENT ON COLUMN messages.content IS 'Content of the message';
COMMENT ON COLUMN messages.created_at IS 'Timestamp when message was created';
COMMENT ON COLUMN messages.updated_at IS 'Timestamp when message was last updated';

-- Add triggers for updated_at
CREATE TRIGGER update_messages_updated_at 
    BEFORE UPDATE ON messages 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();