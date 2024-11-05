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
SET name = $1, email = $2, phone=$3, skills = $4, education = $5, resume_file_address = $6
WHERE applicant = $7
RETURNING name, email, phone, skills, education;

-- name: CreateJob :one
INSERT INTO job (title, description, posted_on, company_name, posted_by)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, title, description, posted_on, company_name, posted_by;

-- name: GetJob :one
SELECT title, description, posted_on, company_name, posted_by, total_applications
FROM job
WHERE id = $1;

-- name: GetApplicants :many
SELECT name, email, address, profile_headline
FROM users
WHERE user_type = 'applicant';

-- name: GetApplicant :one
SELECT u.name, u.email, u.address, u.profile_headline, p.resume_file_address, p.skills, p.education, p.phone
FROM users u
JOIN profile p ON u.id = p.applicant
WHERE u.id = $1;

-- name: GetJobsApplicant :many
SELECT title, description, posted_on, total_applications, company_name, posted_by
FROM job;

-- name: ApplyJob :exec
INSERT INTO apply_jobs (applicant_id, job_id)
VALUES ($1, $2);

-- name: UpdateTotalApplications :exec
UPDATE job
SET total_applications = total_applications + 1
WHERE id = $1;
