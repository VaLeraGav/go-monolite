CREATE TABLE products (
  id BIGSERIAL PRIMARY KEY,
  uuid UUID NOT NULL UNIQUE,
  name VARCHAR(455) NOT NULL,
  unit VARCHAR(32),
  code INT NOT NULL UNIQUE,
  article VARCHAR(32),
  slug VARCHAR(255) NOT NULL UNIQUE,
  active CHAR(1) NOT NULL DEFAULT 'Y' CHECK (active IN ('Y', 'N')),
  step INT CHECK (step > 0),
  brand_uuid VARCHAR(50) DEFAULT NULL REFERENCES property_values(key),
  property JSON,
  weight DOUBLE PRECISION,
  width DOUBLE PRECISION,
  length DOUBLE PRECISION,
  height DOUBLE PRECISION,
  volume DOUBLE PRECISION,
  category_uuid UUID NOT NULL REFERENCES categories(uuid),
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX products_slug_idx ON products (slug);
CREATE UNIQUE INDEX products_uuid_idx ON products (uuid);