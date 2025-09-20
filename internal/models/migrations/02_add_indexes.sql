-- Add index on name for faster searches
CREATE INDEX idx_classifiers_name ON classifiers(name);

-- Add index on created_at for faster sorting in list operation
CREATE INDEX idx_classifiers_created_at ON classifiers(created_at DESC);

-- Update statistics
ANALYZE TABLE classifiers;