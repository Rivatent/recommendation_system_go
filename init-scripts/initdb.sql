CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       username VARCHAR(50) UNIQUE NOT NULL,
                       email VARCHAR(100) UNIQUE NOT NULL,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO users (username, email) VALUES
                                        ('Иван Иванов', 'ivan.ivanov@example.com'),
                                        ('Мария Петрова', 'maria.petrova@example.com'),
                                        ('Алексей Смирнов', 'alexey.smirnov@example.com'),
                                        ('Ольга Кузнецова', 'olga.kuznetsova@example.com'),
                                        ('Дмитрий Васильев', 'dmitry.vasiliev@example.com'),
                                        ('Анна Сергеева', 'anna.sergeeva@example.com'),
                                        ('Сергей Орлов', 'sergey.orlov@example.com'),
                                        ('Елена Фёдорова', 'elena.fyodorova@example.com'),
                                        ('Виктория Лебедева', 'victoria.lebedyeva@example.com'),
                                        ('Артем Громов', 'artem.gromov@example.com'),
                                        ('Наталья Сидорова', 'natalya.sidorova@example.com'),
                                        ('Роман Николаев', 'roman.nikolaev@example.com'),
                                        ('Татьяна Михайлова', 'tatiana.mikhaylova@example.com'),
                                        ('Станислав Крылов', 'stanislav.krylov@example.com'),
                                        ('Ксения Алексеева', 'kseniya.alekseeva@example.com'),
                                        ('Павел Гречишкин', 'pavel.grecheshkin@example.com'),
                                        ('Ирина Дмитриева', 'irina.dmitrieva@example.com');