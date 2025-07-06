CREATE TABLE categories (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID NOT NULL UNIQUE,
    slug VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL, 
    parent_uuid UUID REFERENCES categories(uuid) ON DELETE SET NULL, 
    active CHAR(1) NOT NULL DEFAULT 'Y' CHECK (active IN ('Y', 'N')),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX categories_slug_idx ON categories (slug);
CREATE UNIQUE INDEX category_uuid_idx ON categories (uuid);
