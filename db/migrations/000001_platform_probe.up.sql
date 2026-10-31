BEGIN;
-- Infrastructure probe only: no company, identity, or asset schema in Phase 0.
CREATE TABLE platform_probe (
    singleton boolean PRIMARY KEY DEFAULT true CHECK (singleton),
    installed_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);
INSERT INTO platform_probe (singleton) VALUES (true);
COMMIT;
