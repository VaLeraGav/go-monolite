CREATE TABLE product_prices(
    id SERIAL PRIMARY KEY,
    product_uuid UUID NOT NULL, -- нет связи REFERENCES products(uuid), так как обменом может прийти цена раньше чем товар
    type_price_uuid UUID NOT NULL REFERENCES type_price(uuid) ON DELETE CASCADE,
    active CHAR(1) NOT NULL DEFAULT 'Y' CHECK (active IN ('Y', 'N')),
    price NUMERIC(12, 2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
