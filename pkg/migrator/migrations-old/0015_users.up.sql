CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  email VARCHAR(255) UNIQUE,
  phone VARCHAR(50) UNIQUE,
  name VARCHAR(100),
  last_name VARCHAR(100),
  second_name VARCHAR(100),
  city_id INT REFERENCES cities(id),
  user_type VARCHAR(10) NOT NULL CHECK (user_type IN ('individual', 'legal')),
  inn VARCHAR(20) REFERENCES companies(inn),
  active CHAR(1) NOT NULL DEFAULT 'Y' CHECK (active IN ('Y', 'N')),
  password_hash VARCHAR(255) NOT NULL,
  checkword VARCHAR(255),
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

  token_version VARCHAR(10), -- позволяют контролировать или отзывать access, при необходимости меняем версию и пользовательский  токен становится не валидный 

  -- Гарантирует, что хотя бы email или phone должен быть заполнен
  CONSTRAINT chk_email_or_phone CHECK (
    email IS NOT NULL
    OR phone IS NOT NULL
  )
);