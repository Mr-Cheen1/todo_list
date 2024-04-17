CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    task_text VARCHAR(255),
    createdDate TIMESTAMP,
    expectedDate TIMESTAMP,
    status INTEGER
);