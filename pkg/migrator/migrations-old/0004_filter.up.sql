CREATE TABLE filter (
    id BIGSERIAL PRIMARY KEY,
    property_uuid UUID NOT NULL,
    category_uuid UUID NOT NULL,
    sort INTEGER NOT NULL DEFAULT 0,
    active CHAR(1) NOT NULL DEFAULT 'Y' CHECK (active IN ('Y', 'N')),
    unit VARCHAR(255),
    min_value DOUBLE PRECISION,
    max_value DOUBLE PRECISION,
    string_value VARCHAR(255),  -- для хранения значений типа "Строка"
    CONSTRAINT fk_property_uuid FOREIGN KEY (property_uuid) REFERENCES property(uuid),
    CONSTRAINT fk_category_uuid FOREIGN KEY (category_uuid) REFERENCES categories(uuid),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);