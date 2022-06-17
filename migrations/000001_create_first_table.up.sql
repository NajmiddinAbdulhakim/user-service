CREATE TABLE IF NOT EXISTS students (
    id uuid PRIMARY KEY,
    first_name varchar(200),
    last_name varchar(200),
    age int NOT NULL
)