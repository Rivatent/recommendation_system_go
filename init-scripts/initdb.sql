CREATE TABLE IF NOT EXISTS users (
                                     id SERIAL PRIMARY KEY,
                                     name TEXT NOT NULL,
                                     email TEXT UNIQUE NOT NULL
);

INSERT INTO users (name, email) VALUES
                                    ('Иван Иванов', 'ivan.ivanov@example.com'),
                                    ('Мария Петрова', 'maria.petrova@example.com'),
                                    ('Алексей Смирнов', 'alexey.smirnov@example.com'),
                                    ('Ольга Кузнецова', 'olga.kuznetsova@example.com'),
                                    ('Дмитрий Васильев', 'dmitry.vasiliev@example.com');