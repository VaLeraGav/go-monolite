CREATE TABLE auth_codes (
  id SERIAL PRIMARY KEY,
  email VARCHAR(255),                       
  phone VARCHAR(50),                      
  code VARCHAR(10) NOT NULL,               
  expires_at TIMESTAMP NOT NULL,            -- срок действия
  used BOOLEAN DEFAULT FALSE,               -- был ли использован
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)