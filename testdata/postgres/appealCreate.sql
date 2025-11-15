DROP TABLE IF EXISTS user_appeal;
DROP TABLE IF EXISTS user_appeal_message;
DROP TABLE IF EXISTS user_appeal_response;

CREATE TABLE user_appeal (
    user_appeal_id BIGINT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    tag TEXT NOT NULL CHECK (tag IN ('bug', 'feature', 'other')),
    name TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('open', 'in_progress', 'resolved')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES USER(user_id) ON DELETE CASCADE
);

CREATE TABLE USER_APPEAL_MESSAGE (
     user_appeal_message_id BIGINT PRIMARY KEY,
     user_appeal_id BIGINT NOT NULL,
     message_content TEXT NOT NULL,
     created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
     updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
     FOREIGN KEY (user_appeal_id) REFERENCES USER_APPEAL(user_appeal_id) ON DELETE CASCADE
);

CREATE TABLE USER_APPEAL_RESPONSE (
      user_appeal_response_id BIGINT PRIMARY KEY,
      user_appeal_id BIGINT NOT NULL,
      response_message TEXT NOT NULL,
      created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
      updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
      FOREIGN KEY (user_appeal_id) REFERENCES USER_APPEAL(user_appeal_id) ON DELETE CASCADE
);

CREATE INDEX idx_user_appeal_user_id ON USER_APPEAL(user_id);
CREATE INDEX idx_user_appeal_status ON USER_APPEAL(status);
CREATE INDEX idx_user_appeal_created_at ON USER_APPEAL(created_at);
CREATE INDEX idx_user_appeal_message_user_appeal_id ON USER_APPEAL_MESSAGE(user_appeal_id);
CREATE INDEX idx_user_appeal_response_user_appeal_id ON USER_APPEAL_RESPONSE(user_appeal_id);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_user_appeal_updated_at
    BEFORE UPDATE ON USER_APPEAL
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_appeal_message_updated_at
    BEFORE UPDATE ON USER_APPEAL_MESSAGE
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_appeal_response_updated_at
    BEFORE UPDATE ON USER_APPEAL_RESPONSE
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
