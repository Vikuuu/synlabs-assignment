-- name: CreateUser :one
INSERT INTO users (name, email, address, user_type, password_hash, profile_headline)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, name, email, user_type;

-- name: GetUser :one
SELECT id, password_hash, user_type FROM users
WHERE email = $1;

-- name: GetUserFromID :one
SELECT user_type FROM users
WHERE id = $1;

-- name: CreateApplicantProfile :one
INSERT INTO profile (applicant)
VALUES ($1)
RETURNING applicant;

-- name: AddProfileIDInUser :exec
UPDATE users
SET profile_id = $1;

-- name: UpdateProfile :one
UPDATE profile
SET name = $1, email = $2, phone=$3, skills = $4, education = $5
WHERE applicant = $6
RETURNING name, email, phone, skills, education;

-- name: CreateJob :one
INSERT INTO job (title, description, posted_on, company_name, posted_by)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, title, description, posted_on, company_name, posted_by;

-- name: GetJob :one
SELECT title, description, posted_on, company_name, posted_by
FROM job
WHERE id = $1;
