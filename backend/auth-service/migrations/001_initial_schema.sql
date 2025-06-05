-- Drop existing tables and indexes
DROP TABLE IF EXISTS refresh_tokens CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'user',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create refresh_tokens table
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token ON refresh_tokens(token);

-- Создаем функцию для синхронизации пользователей
CREATE OR REPLACE FUNCTION sync_users_to_forum()
RETURNS TRIGGER AS $$
BEGIN
    -- Подключаемся к базе данных forum
    PERFORM dblink_connect('forum_db', 'dbname=forum user=postgres password=postgres');
    
    -- Синхронизируем пользователя
    IF TG_OP = 'INSERT' THEN
        PERFORM dblink_exec('forum_db', format(
            'INSERT INTO users (id, username, email, password_hash, role, created_at, updated_at)
             VALUES (%L, %L, %L, %L, %L, %L, %L)
             ON CONFLICT (id) DO UPDATE SET
             username = EXCLUDED.username,
             email = EXCLUDED.email,
             password_hash = EXCLUDED.password_hash,
             role = EXCLUDED.role,
             updated_at = EXCLUDED.updated_at',
            NEW.id, NEW.username, NEW.email, NEW.password_hash, NEW.role, NEW.created_at, NEW.updated_at
        ));
    ELSIF TG_OP = 'UPDATE' THEN
        PERFORM dblink_exec('forum_db', format(
            'UPDATE users SET
             username = %L,
             email = %L,
             password_hash = %L,
             role = %L,
             updated_at = %L
             WHERE id = %L',
            NEW.username, NEW.email, NEW.password_hash, NEW.role, NEW.updated_at, NEW.id
        ));
    ELSIF TG_OP = 'DELETE' THEN
        PERFORM dblink_exec('forum_db', format(
            'DELETE FROM users WHERE id = %L',
            OLD.id
        ));
    END IF;
    
    -- Отключаемся от базы данных forum
    PERFORM dblink_disconnect('forum_db');
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Создаем триггер для синхронизации
DROP TRIGGER IF EXISTS sync_users_trigger ON users;
CREATE TRIGGER sync_users_trigger
AFTER INSERT OR UPDATE OR DELETE ON users
FOR EACH ROW
EXECUTE FUNCTION sync_users_to_forum();

-- Синхронизируем существующих пользователей
DO $$
DECLARE
    user_record RECORD;
BEGIN
    FOR user_record IN SELECT * FROM users LOOP
        PERFORM dblink_connect('forum_db', 'dbname=forum user=postgres password=postgres');
        PERFORM dblink_exec('forum_db', format(
            'INSERT INTO users (id, username, email, password_hash, role, created_at, updated_at)
             VALUES (%L, %L, %L, %L, %L, %L, %L)
             ON CONFLICT (id) DO UPDATE SET
             username = EXCLUDED.username,
             email = EXCLUDED.email,
             password_hash = EXCLUDED.password_hash,
             role = EXCLUDED.role,
             updated_at = EXCLUDED.updated_at',
            user_record.id, user_record.username, user_record.email, user_record.password_hash,
            user_record.role, user_record.created_at, user_record.updated_at
        ));
        PERFORM dblink_disconnect('forum_db');
    END LOOP;
END $$; 