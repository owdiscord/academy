-- migrate:up
ALTER TABLE staff 
  ADD COLUMN has_avatar BOOLEAN DEFAULT false,
  ADD COLUMN avatar_hash VARCHAR(128) NULL;

-- migrate:down
ALTER TABLE staff
  DROP COLUMN has_avatar,
  DROP COLUMN avatar_hash;
