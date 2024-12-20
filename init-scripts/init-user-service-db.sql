CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
   CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users (
                       id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                       username VARCHAR(50) UNIQUE NOT NULL,
                       email VARCHAR(100) UNIQUE NOT NULL,
                       password VARCHAR(255) NOT NULL,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO users (username, email, password) VALUES
                                                  ('Иван Иванов', 'ivan.ivanov@example.com', crypt('password1', gen_salt('bf'))),
                                                  ('Мария Петрова', 'maria.petrova@example.com', crypt('password2', gen_salt('bf'))),
                                                  ('Алексей Смирнов', 'alexey.smirnov@example.com', crypt('password3', gen_salt('bf'))),
                                                  ('Ольга Кузнецова', 'olga.kuznetsova@example.com', crypt('password4', gen_salt('bf'))),
                                                  ('Дмитрий Васильев', 'dmitry.vasiliev@example.com', crypt('password5', gen_salt('bf'))),
                                                  ('Анна Сергеева', 'anna.sergeeva@example.com', crypt('password6', gen_salt('bf'))),
                                                  ('Сергей Орлов', 'sergey.orlov@example.com', crypt('password7', gen_salt('bf'))),
                                                  ('Елена Фёдорова', 'elena.fyodorova@example.com', crypt('password8', gen_salt('bf'))),
                                                  ('Виктория Лебедева', 'victoria.lebedyeva@example.com', crypt('password9', gen_salt('bf'))),
                                                  ('Артем Громов', 'artem.gromov@example.com', crypt('password10', gen_salt('bf'))),
                                                  ('Наталья Сидорова', 'natalya.sidorova@example.com', crypt('password11', gen_salt('bf'))),
                                                  ('Роман Николаев', 'roman.nikolaev@example.com', crypt('password12', gen_salt('bf'))),
                                                  ('Татьяна Михайлова', 'tatiana.mikhaylova@example.com', crypt('password13', gen_salt('bf'))),
                                                  ('Станислав Крылов', 'stanislav.krylov@example.com', crypt('password14', gen_salt('bf'))),
                                                  ('Ксения Алексеева', 'kseniya.alekseeva@example.com', crypt('password15', gen_salt('bf'))),
                                                  ('Павел Гречишкин', 'pavel.grecheshkin@example.com', crypt('password16', gen_salt('bf'))),
                                                  ('Ирина Дмитриева', 'irina.dmitrieva@example.com', crypt('password17', gen_salt('bf')));

