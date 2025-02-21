-- add email column to company table
ALTER TABLE company ADD COLUMN email TEXT;
ALTER TABLE company ADD CONSTRAINT email_unique UNIQUE(email);

-- add api_key column to company table
ALTER TABLE company ADD COLUMN api_key TEXT;
ALTER TABLE company ADD CONSTRAINT api_key UNIQUE(api_key);
