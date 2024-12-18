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

CREATE TABLE products (
                          id SERIAL PRIMARY KEY,
                          name VARCHAR(100) NOT NULL,
                          description TEXT,
                          price DECIMAL(10, 2) NOT NULL,
                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                          updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO products (name, description, price) VALUES
    ('The Last of Us Part II', 'Эмоциональное продолжение культового приключения в постапокалиптическом мире.', 59.99),
    ('Spider-Man: Miles Morales', 'Новый супергеройский опыт с учетом высоких технологий и городской жизни.', 49.99),
    ('Demons Souls', 'Переосмысленный классический RPG с захватывающим миром и сложными противниками.', 69.99),
    ('Ghost of Tsushima', 'Эпическое приключение самурая в открытом мире Японии.', 59.99),
    ('Final Fantasy VII Remake', 'Современная переработка классической JRPG с новой графикой и сюжетом.', 79.99),
    ('Ratchet & Clank: Rift Apart', 'Красочная платформенная игра с увлекательным сюжетом и динамичным геймплеем.', 69.99),
    ('Horizon Forbidden West', 'Захватывающее продолжение приключений Элой в мире, полном опасностей и загадок.', 59.99),
    ('Resident Evil Village', 'Страшный и захватывающий хоррор с элементами выживания и расследования.', 59.99),
    ('Gran Turismo 7', 'Совершенная автосимуляция с реалистичной графикой и широким выбором автомобилей.', 69.99),
    ('Returnal', 'Новообразный шутер с элементами roguelike и увлекательным сюжетом.', 59.99),
    ('Cyberpunk 2077', 'Открытый мир с высокими технологиями и многогранной историей, погружающий игроков в будущее.', 59.99),
    ('Heavy Rain', 'Интерактивный детективный триллер с несколькими концовками, основанный на принятии решений.', 39.99),
    ('Detroit: Become Human', 'Динамичное приключение, исследующее вопросы человеческости через призму андроидов и выборов.', 49.99);