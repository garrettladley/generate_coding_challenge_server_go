CREATE DOMAIN nuid_domain AS varchar(9)
    CHECK (value ~ '^[0-9]{9}$');

CREATE DOMAIN applicant_name_domain AS varchar(256)
    CHECK (value !~ '[/()"<>\\{}]');

CREATE TABLE IF NOT EXISTS applicants (
    nuid nuid_domain PRIMARY KEY,
    applicant_name applicant_name_domain NOT NULL,
    registration_time timestamp with time zone NOT NULL,
    token uuid UNIQUE NOT NULL,
    challenge text[] NOT NULL,
    solution text[] NOT NULL
);

CREATE TABLE IF NOT EXISTS submissions (
    submission_id serial PRIMARY KEY,
    nuid nuid_domain NOT NULL REFERENCES applicants (nuid),
    correct boolean NOT NULL,
    submission_time timestamp with time zone NOT NULL
);