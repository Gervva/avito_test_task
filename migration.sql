CREATE TABLE IF NOT EXISTS Banner (
  id UUID PRIMARY KEY,
  group_id BIGSERIAL NOT NULL,
  feature_id BIGINT NOT NULL, 
  version BIGINT NOT NULL DEFAULT 1,
  content jsonb NOT NULL,
  is_active BOOLEAN NOT NULL, 
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

CREATE TABLE IF NOT EXISTS Tag (
    tag_id BIGINT NOT NULL,
    banner_id UUID NOT NULL REFERENCES Banner (id) ON DELETE CASCADE
);

CREATE INDEX banner_feature ON Banner(feature_id);
CREATE INDEX banner_tag ON Tag(id);
