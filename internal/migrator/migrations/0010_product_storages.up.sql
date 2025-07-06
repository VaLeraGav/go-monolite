CREATE TABLE product_storages (
    id SERIAL PRIMARY KEY,
    product_uuid UUID NOT NULL, -- нет связи REFERENCES products(uuid), так как обменом может прийти цена раньше чем товар
    storage_uuid UUID NOT NULL REFERENCES storage(uuid) ON DELETE CASCADE,
    active CHAR(1) NOT NULL DEFAULT 'Y' CHECK (active IN ('Y', 'N')),
    quantity INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
