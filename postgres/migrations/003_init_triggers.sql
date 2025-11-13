CREATE OR REPLACE FUNCTION max_two_reviewers()
RETURNS TRIGGER AS $$

DECLARE
    rev_count INT;
BEGIN

SELECT COUNT(*) INTO rev_count
    FROM users_pull_requests
    WHERE pull_requests_id = pull_requests_id;

IF rev_count >= 2 THEN
    RAISE EXCEPTION 'Maximum 2 reviewers allowed per PR';
END IF;

RETURN NEW;

END
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_max_two_reviewers
BEFORE INSERT ON users_pull_requests
FOR EACH ROW EXECUTE FUNCTION max_two_reviewers();

--------------------------------------------------------------------
--------------------------------------------------------------------

CREATE OR REPLACE FUNCTION merged_requests()
RETURNS TRIGGER AS $$
BEGIN

IF OLD.status = 'MERGED' THEN
    IF NEW.status = 'OPEN' THEN
        RETURN NEW;
    ELSE
        RAISE EXCEPTION 'Cannot update merged pull request except unmerging to OPEN';
    END IF;
END IF;

END
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_merged_requests
BEFORE UPDATE ON pull_requests
FOR EACH ROW EXECUTE FUNCTION merged_requests();