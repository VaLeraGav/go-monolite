CREATE TABLE product_images (
  id BIGSERIAL PRIMARY KEY,
  product_uuid UUID REFERENCES products(uuid) ON DELETE CASCADE,
  url TEXT NOT NULL,
  position INT DEFAULT 0,
  is_main BOOLEAN DEFAULT FALSE -- основная картинка
);