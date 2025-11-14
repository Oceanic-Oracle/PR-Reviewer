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

CREATE OR REPLACE FUNCTION prevent_reviewer_change_if_merged()
RETURNS TRIGGER AS $$
DECLARE
    pr_status TEXT;
    pr_id TEXT;
BEGIN
    IF TG_OP = 'DELETE' THEN
        pr_id := OLD.pull_requests_id;
    ELSE
        pr_id := NEW.pull_requests_id;
    END IF;

    SELECT status INTO pr_status
    FROM pull_requests
    WHERE id = pr_id;

    IF pr_status = 'MERGED' THEN
        RAISE EXCEPTION 'Cannot modify reviewers for a merged pull request';
    END IF;

    IF TG_OP = 'DELETE' THEN
        RETURN OLD;
    ELSE
        RETURN NEW;
    END IF;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_prevent_reviewer_change_if_merged
BEFORE INSERT OR UPDATE OR DELETE ON users_pull_requests
FOR EACH ROW
EXECUTE FUNCTION prevent_reviewer_change_if_merged();

--------------------------------------------------------------------
--------------------------------------------------------------------

CREATE OR REPLACE FUNCTION switch_status()
RETURNS TRIGGER AS $$
BEGIN

IF OLD.status != NEW.status THEN
    IF NEW.status = 'OPEN' THEN
        NEW.merged_at = NULL;
    ELSIF NEW.status = 'MERGED' THEN
        NEW.merged_at = NOW();
    END IF;
END IF;

RETURN NEW;

END
BEFORE UPDATE ON pull_requests
FOR EACH ROW EXECUTE FUNCTION switch_status();