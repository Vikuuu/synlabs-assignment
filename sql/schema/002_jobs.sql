-- +goose Up 
CREATE TABLE apply_jobs (
    applicant_id INT REFERENCES users(id),
    job_id INT REFERENCES job(id)
);

-- +goose Down
DROP TABLE apply_jobs;
