-- Add owner_id to resources table to track resource ownership

-- Add owner_id column
ALTER TABLE resources ADD COLUMN IF NOT EXISTS owner_id INT;

-- Add foreign key constraint
ALTER TABLE resources ADD CONSTRAINT fk_resources_owner
    FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE SET NULL;

-- Create index for faster owner queries
CREATE INDEX IF NOT EXISTS idx_resources_owner_id ON resources(owner_id);

-- Add comment
COMMENT ON COLUMN resources.owner_id IS 'Owner of the resource (user who created it)';

-- Update existing resources to assign them to first admin user or first user
-- This ensures existing data remains valid
DO $$
DECLARE
    first_user_id INT;
BEGIN
    -- Try to get first admin user, otherwise get first user
    SELECT id INTO first_user_id FROM users WHERE role = 'admin' LIMIT 1;

    IF first_user_id IS NULL THEN
        SELECT id INTO first_user_id FROM users LIMIT 1;
    END IF;

    -- Update resources that don't have an owner
    IF first_user_id IS NOT NULL THEN
        UPDATE resources SET owner_id = first_user_id WHERE owner_id IS NULL;
    END IF;
END $$;
