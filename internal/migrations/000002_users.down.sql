ALTER TABLE todo_schema.todos DROP CONSTRAINT IF EXISTS fk_todos_user_id;

DROP TRIGGER IF EXISTS set_users_updated_at ON user_schema.users;
DROP TABLE IF EXISTS user_schema.refresh_tokens;
DROP TABLE IF EXISTS user_schema.users;
DROP SCHEMA IF EXISTS user_schema;
