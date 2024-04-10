CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    task_text VARCHAR(255),
    task_date TIMESTAMP,
    status VARCHAR(50)
);