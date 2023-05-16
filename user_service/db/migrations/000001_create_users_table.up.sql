CREATE TABLE  users (
                        id uuid PRIMARY KEY,
                        name text NOT NULL,
                        description text,
                        login text NOT NULL,
                        password text NOT NULL
);