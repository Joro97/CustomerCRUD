CREATE TABLE IF NOT EXISTS customers (
                                         id UUID PRIMARY KEY,
                                         first_name TEXT NOT NULL,
                                         middle_name TEXT,
                                         last_name TEXT NOT NULL,
                                         email TEXT NOT NULL UNIQUE,
                                         phone_number TEXT
);
