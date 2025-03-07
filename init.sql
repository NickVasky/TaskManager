BEGIN;
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        username VARCHAR(100) UNIQUE NOT NULL,
        password VARCHAR(100) NOT NULL,
        first_name VARCHAR(100),
        second_name VARCHAR(100),
        created_at TIMESTAMP WITH TIME ZONE
    );

    CREATE TABLE IF NOT EXISTS sessions (
        id SERIAL PRIMARY KEY,
        user_id INT REFERENCES users (id),
        session_token VARCHAR(64),
        csrf_token VARCHAR(64),
        expires_at TIMESTAMP WITH TIME ZONE
    );

    CREATE TABLE IF NOT EXISTS tasks (
        id SERIAL PRIMARY KEY,
        user_id INT REFERENCES users (id),
        title VARCHAR(256),
        goal TEXT,
        measure TEXT,
        relevance TEXT,
        is_done BOOLEAN,
        deadline TIMESTAMP WITH TIME ZONE,
        created_at TIMESTAMP WITH TIME ZONE,
        finished_at TIMESTAMP WITH TIME ZONE
    );
COMMIT;