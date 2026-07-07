CREATE SCHEMA IF NOT EXISTS user_schema;

CREATE TABLE user_schema.users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    username VARCHAR(128) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- CREATE TRIGGER set_users_updated_at
-- BEFORE UPDATE ON user_schema.users
-- FOR EACH ROW
-- EXECUTE FUNCTION public.update_updated_at_column();

ALTER TABLE todo_schema.todos
ADD CONSTRAINT fk_todos_user_id
FOREIGN KEY (user_id)
REFERENCES user_schema.users(id)
ON DELETE CASCADE;
