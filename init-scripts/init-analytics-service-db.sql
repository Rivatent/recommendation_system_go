CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE analytics (
                          id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                          total_users INT DEFAULT 0,
                          total_sales INT DEFAULT 0,
                          sales_progression_rate DECIMAL(10,2) DEFAULT 0.00,
                          users_progression_rate DECIMAL(10,2) DEFAULT 0.00,
                          product_average_rating DECIMAL(10,2) DEFAULT 0.00,
                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                          updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO analytics (total_users, total_sales, sales_progression_rate, users_progression_rate, product_average_rating)
VALUES
    (5, 3, 0.60, 0.50, 4.00),
    (10, 7, 0.70, 1.00, 4.00),
    (17, 13, 0.76, 1.70, 4.00);