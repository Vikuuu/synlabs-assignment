-- +goose Up 
CREATE TYPE user_type AS ENUM('applicant', 'admin');

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    address TEXT NOT NULL,
    user_type user_type NOT NULL DEFAULT 'applicant',
    password_hash TEXT NOT NULL,
    profile_headline TEXT NOT NULL,
    profile_id INT
);

CREATE TABLE profile (
    applicant INT PRIMARY KEY REFERENCES users(id),
    resume_file_address TEXT NOT NULL,
    skills TEXT,
    education TEXT,
    name TEXT,
    email TEXT,
    phone TEXT
);

ALTER TABLE users
ADD CONSTRAINT fk_profile FOREIGN KEY (profile_id)
REFERENCES profile(applicant);

CREATE TABLE job (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    posted_on TIMESTAMP NOT NULL,
    total_applications INT DEFAULT 0,
    company_name TEXT NOT NULL,
    posted_by INT NOT NULL,
    CONSTRAINT fk_posted_by FOREIGN KEY (posted_by)
        REFERENCES users(id)
);

-- +goose Down
DROP TABLE job;
DROP TABLE profile;
DROP TABLE users;
DROP TYPE type;
