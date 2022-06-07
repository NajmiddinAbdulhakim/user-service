CREATE TABLE IF NOT EXISTS users (
    id uuid NOT NULL UNIQUE PRIMARY KEY,
    first_name varchar(50) NOT NULL,
    last_name varchar(50) NOT NULL,
    user_name varchar(50) NOT NULL  UNIQUE,
    email varchar(50) NOT NULL UNIQUE,
    phone_number TEXT [],
    bio TEXT,
    status varchar(20),
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS adresses (
    user_id uuid NOT NULL,
    country varchar(50) NOT NULL,
    city varchar(50) NOT NULL,
    district varchar(50) NOT NULL,
    postalCode integer NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);