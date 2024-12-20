CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE recommendations (
                                 id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),        -- Уникальный идентификатор записи
                                 user_id UUID,-- REFERENCES users(id) ON DELETE CASCADE,   -- Идентификатор пользователя
                                 product_id UUID,-- REFERENCES products(id) ON DELETE CASCADE, -- Идентификатор продукта
                                 score DECIMAL(3, 2) NOT NULL DEFAULT 0.00,              -- Оценка рекомендации
                                 created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,         -- Дата и время создания записи
                                 updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,         -- Дата и время последнего обновления
                                 UNIQUE (user_id, product_id)                            -- Уникальная пара (пользователь-продукт)
);

INSERT INTO recommendations (user_id, product_id, score) VALUES
                                                             ('c73c1a27-48e1-4c2c-b48d-e02957cca1b0', 'c73c1a27-48e1-4c2c-b48d-e02957cca1b1', 4.5);