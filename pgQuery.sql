CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE NULL,
    updated_at TIMESTAMP WITH TIME ZONE NULL,
    deleted_at TIMESTAMP WITH TIME ZONE NULL
);

CREATE INDEX idx_roles_deleted_at ON roles (deleted_at);

INSERT INTO roles (id,name) VALUES (1,'ADMIN');
INSERT INTO roles (id,name) VALUES (2,'VP');
INSERT INTO roles (id,name) VALUES (3,'HR');
INSERT INTO roles (id,name) VALUES (4,'EMPLOYEE');


SELECT id, name, created_at, updated_at, deleted_at
	FROM roles;