DROP TABLE IF EXISTS user_appeal;
DROP TABLE IF EXISTS user_appeal_message;

CREATE TABLE IF NOT EXISTS "user" (
    user_id BIGINT PRIMARY KEY,
    username TEXT NOT NULL
);

CREATE TABLE user_appeal (
     user_appeal_id BIGINT PRIMARY KEY,
     user_id BIGINT NOT NULL,
     tag TEXT NOT NULL CHECK (tag IN ('bug', 'feature', 'other')),
     name TEXT NOT NULL,
     status TEXT NOT NULL CHECK (status IN ('open', 'in_progress', 'resolved')),
     created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
     updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
     FOREIGN KEY (user_id) REFERENCES "user"(user_id) ON DELETE CASCADE
);

CREATE TABLE user_appeal_message (
     user_appeal_message_id BIGINT PRIMARY KEY,
     user_appeal_id BIGINT NOT NULL,
     is_response BOOLEAN NOT NULL DEFAULT FALSE,
     message_content TEXT NOT NULL,
     created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
     updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
     FOREIGN KEY (user_appeal_id) REFERENCES user_appeal(user_appeal_id) ON DELETE CASCADE
);

CREATE INDEX idx_user_appeal_user_id ON user_appeal(user_id);
CREATE INDEX idx_user_appeal_status ON user_appeal(status);
CREATE INDEX idx_user_appeal_created_at ON user_appeal(created_at);
CREATE INDEX idx_user_appeal_message_user_appeal_id ON USER_APPEAL_MESSAGE(user_appeal_id);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_user_appeal_updated_at
    BEFORE UPDATE ON user_appeal
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_appeal_message_updated_at
    BEFORE UPDATE ON user_appeal_message
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
