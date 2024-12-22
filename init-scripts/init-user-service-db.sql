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

INSERT INTO users (id, username, email, password) VALUES
                                                      ('e97f96c0-3136-4705-a5f5-6c28a5eae9b5', 'Иван Иванов', 'ivan.ivanov@example.com', crypt('password1', gen_salt('bf'))),
                                                      ('bfbf025d-05bc-4e3c-b7b4-0c1eb7b648a2', 'Мария Петрова', 'maria.petrova@example.com', crypt('password2', gen_salt('bf'))),
                                                      ('a3885686-73b0-4c44-bb30-86d2f476c160', 'Алексей Смирнов', 'alexey.smirnov@example.com', crypt('password3', gen_salt('bf'))),
                                                      ('9620d84e-bf49-4a24-8c7f-6c52f1c095ab', 'Ольга Кузнецова', 'olga.kuznetsova@example.com', crypt('password4', gen_salt('bf'))),
                                                      ('9a81be64-e1e8-4a7a-aa0e-4d024b0cbded', 'Дмитрий Васильев', 'dmitry.vasiliev@example.com', crypt('password5', gen_salt('bf'))),
                                                      ('3ea5bc56-614e-474f-8a8a-2d77d5de14f8', 'Анна Сергеева', 'anna.sergeeva@example.com', crypt('password6', gen_salt('bf'))),
                                                      ('5d1a5e23-7821-49f6-8c04-ac174cf8cf12', 'Сергей Орлов', 'sergey.orlov@example.com', crypt('password7', gen_salt('bf'))),
                                                      ('15aa0517-f8fd-4db9-86af-339b3838706f', 'Елена Фёдорова', 'elena.fyodorova@example.com', crypt('password8', gen_salt('bf'))),
                                                      ('b7c0992d-0b4e-4b85-9d04-15727c5c0f5d', 'Виктория Лебедева', 'victoria.lebedyeva@example.com', crypt('password9', gen_salt('bf'))),
                                                      ('b382bd29-4f2b-49ef-8189-e2da94a1325e', 'Артем Громов', 'artem.gromov@example.com', crypt('password10', gen_salt('bf'))),
                                                      ('26b1b2f6-db49-4170-87ee-62d7341ed6e5', 'Наталья Сидорова', 'natalya.sidorova@example.com', crypt('password11', gen_salt('bf'))),
                                                      ('fd693d5f-0c08-4fad-a90c-646027b496a6', 'Роман Николаев', 'roman.nikolaev@example.com', crypt('password12', gen_salt('bf'))),
                                                      ('47e31b0b-0c5c-4324-b29b-8e45b23ed46f', 'Татьяна Михайлова', 'tatiana.mikhaylova@example.com', crypt('password13', gen_salt('bf'))),
                                                      ('33395518-53f6-4eb8-a32d-717f5b827e04', 'Станислав Крылов', 'stanislav.krylov@example.com', crypt('password14', gen_salt('bf'))),
                                                      ('caf46f2b-17dc-46d9-858b-723f3424899d', 'Ксения Алексеева', 'kseniya.alekseeva@example.com', crypt('password15', gen_salt('bf'))),
                                                      ('eae81bd9-ffca-4c0f-b4a4-8abf81102ebb', 'Павел Гречишкин', 'pavel.grecheshkin@example.com', crypt('password16', gen_salt('bf'))),
                                                      ('da66f449-d978-4f26-85c4-c3577c6e3c9b', 'Ирина Дмитриева', 'irina.dmitrieva@example.com', crypt('password17', gen_salt('bf')));