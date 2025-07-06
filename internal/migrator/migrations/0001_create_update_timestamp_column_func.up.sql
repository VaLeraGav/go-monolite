-- Функция для автоматического обновления `updated_at`
CREATE
    OR REPLACE FUNCTION update_timestamp_column()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at
        = NOW();
    RETURN NEW;
END;
$$
    LANGUAGE plpgsql;

