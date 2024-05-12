-- Создание таблицы tasks.
CREATE TABLE tasks (
    -- Первичный ключ id с автоинкрементом.
    id SERIAL PRIMARY KEY,
    
    -- Текст задачи (максимум 255 символов).
    task_text VARCHAR(255),
    
    -- Дата создания задачи.
    createdDate DATE,
    
    -- Ожидаемая дата выполнения задачи.
    expectedDate DATE,
    
    -- Статус задачи (целое число).
    status INTEGER
);