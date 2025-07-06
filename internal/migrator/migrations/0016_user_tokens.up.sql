CREATE TABLE user_tokens (
  id BIGSERIAL PRIMARY KEY,
  refresh_token TEXT NOT NULL,
  user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
  device_id VARCHAR(255) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(user_id, device_id) -- один активный refresh на устройство (можно убрать, если нужно хранить историю)
);