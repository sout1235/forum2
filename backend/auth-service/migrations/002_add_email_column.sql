-- Добавляем поле email в таблицу users как nullable
ALTER TABLE users ADD COLUMN IF NOT EXISTS email VARCHAR(255);

-- Обновляем существующие записи, устанавливая email на основе username
UPDATE users SET email = username || '@example.com' WHERE email IS NULL;

-- Теперь делаем поле NOT NULL
ALTER TABLE users ALTER COLUMN email SET NOT NULL;

-- Добавляем уникальный индекс
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Синхронизируем изменения с базой forum
DO $$
BEGIN
    -- Подключаемся к базе данных forum
    PERFORM dblink_connect('forum_db', 'dbname=forum user=postgres password=postgres');
    
    -- Добавляем поле email в таблицу users в базе forum как nullable
    PERFORM dblink_exec('forum_db', 'ALTER TABLE users ADD COLUMN IF NOT EXISTS email VARCHAR(255)');
    
    -- Обновляем существующие записи в базе forum
    PERFORM dblink_exec('forum_db', 'UPDATE users SET email = username || ''@example.com'' WHERE email IS NULL');
    
    -- Делаем поле NOT NULL в базе forum
    PERFORM dblink_exec('forum_db', 'ALTER TABLE users ALTER COLUMN email SET NOT NULL');
    
    -- Создаем уникальный индекс в базе forum
    PERFORM dblink_exec('forum_db', 'CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users(email)');
    
    -- Отключаемся от базы данных forum
    PERFORM dblink_disconnect('forum_db');
END $$; 