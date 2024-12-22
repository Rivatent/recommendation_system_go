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
                                                             ('e97f96c0-3136-4705-a5f5-6c28a5eae9b5', 'dfd719d3-b51f-4f60-b2e4-58c61c6efc10', 4.50),  -- Иван Иванов
                                                             ('e97f96c0-3136-4705-a5f5-6c28a5eae9b5', '18abf4be-7a3b-4b62-a30e-3347ca8ca613', 3.80),  -- Иван Иванов
                                                             ('bfbf025d-05bc-4e3c-b7b4-0c1eb7b648a2', 'c0a3eefb-a429-4789-84b3-b1d2d10b7da6', 5.00),  -- Мария Петрова
                                                             ('bfbf025d-05bc-4e3c-b7b4-0c1eb7b648a2', '90b0a9eb-afe7-44c9-9e01-34a1a0cb6f29', 2.30),  -- Мария Петрова
                                                             ('a3885686-73b0-4c44-bb30-86d2f476c160', 'a4ad81cb-af80-4754-b987-09e4444d7cc0', 4.20),  -- Алексей Смирнов
                                                             ('9620d84e-bf49-4a24-8c7f-6c52f1c095ab', 'dfd719d3-b51f-4f60-b2e4-58c61c6efc10', 3.60),  -- Ольга Кузнецова
                                                             ('9a81be64-e1e8-4a7a-aa0e-4d024b0cbded', '18abf4be-7a3b-4b62-a30e-3347ca8ca613', 4.40),  -- Дмитрий Васильев
                                                             ('3ea5bc56-614e-474f-8a8a-2d77d5de14f8', 'c0a3eefb-a429-4789-84b3-b1d2d10b7da6', 4.75),  -- Анна Сергеева
                                                             ('5d1a5e23-7821-49f6-8c04-ac174cf8cf12', '90b0a9eb-afe7-44c9-9e01-34a1a0cb6f29', 2.00),  -- Сергей Орлов
                                                             ('9a81be64-e1e8-4a7a-aa0e-4d024b0cbded', 'a4ad81cb-af80-4754-b987-09e4444d7cc0', 5.00);  -- Дмитрий Васильев