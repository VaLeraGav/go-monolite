CREATE TABLE filter_values (
    id BIGSERIAL PRIMARY KEY,
    filter_id BIGINT NOT NULL,
    property_values_id BIGINT NOT NULL,
    CONSTRAINT fk_filter FOREIGN KEY (filter_id) REFERENCES filter(id),
    CONSTRAINT fk_property_values FOREIGN KEY (property_values_id) REFERENCES property_values(id)
);

CREATE INDEX idx_filter_values_filter_id ON filter_values USING HASH (filter_id);
CREATE INDEX idx_filter_values_property_values_id ON filter_values USING HASH (property_values_id);
