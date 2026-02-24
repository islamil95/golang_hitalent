-- +goose Up
CREATE TABLE IF NOT EXISTS departments (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    parent_id INTEGER REFERENCES departments(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_departments_name_parent ON departments (parent_id, name) WHERE parent_id IS NOT NULL;
CREATE UNIQUE INDEX idx_departments_name_root ON departments (name) WHERE parent_id IS NULL;
CREATE INDEX idx_departments_parent_id ON departments(parent_id);

CREATE TABLE IF NOT EXISTS employees (
    id SERIAL PRIMARY KEY,
    department_id INTEGER NOT NULL REFERENCES departments(id) ON DELETE CASCADE,
    full_name VARCHAR(200) NOT NULL,
    position VARCHAR(200) NOT NULL,
    hired_at DATE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_employees_department_id ON employees(department_id);

-- +goose Down
DROP TABLE IF EXISTS employees;
DROP TABLE IF EXISTS departments;
