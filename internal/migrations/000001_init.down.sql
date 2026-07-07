DROP TRIGGER IF EXISTS set_todos_updated_at ON todo_schema.todos;
DROP TABLE IF EXISTS todo_schema.todos;
DROP SCHEMA IF EXISTS todo_schema;
DROP FUNCTION IF EXISTS public.update_updated_at_column();
