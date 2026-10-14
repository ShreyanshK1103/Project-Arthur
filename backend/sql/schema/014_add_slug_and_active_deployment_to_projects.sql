-- +goose Up
ALTER TABLE projects ADD COLUMN slug TEXT;
ALTER TABLE projects ADD COLUMN active_deployment_id UUID
    REFERENCES deployments(id) ON DELETE SET NULL;

-- Backfill slugs for existing projects: sanitized name + 4 chars of the id hash
UPDATE projects
SET slug = COALESCE(
        NULLIF(trim(both '-' from lower(regexp_replace(name, '[^a-zA-Z0-9]+', '-', 'g'))), ''),
        'project'
    ) || '-' || substr(md5(id::text), 1, 4);

-- Backfill active deployment: latest successful one per project
UPDATE projects p
SET active_deployment_id = (
    SELECT d.id
    FROM deployments d
    WHERE d.project_id = p.id AND d.status = 'success'
    ORDER BY d.created_at DESC
    LIMIT 1
);

ALTER TABLE projects ALTER COLUMN slug SET NOT NULL;
CREATE UNIQUE INDEX projects_slug_key ON projects (slug);

-- +goose Down
DROP INDEX projects_slug_key;
ALTER TABLE projects DROP COLUMN active_deployment_id;
ALTER TABLE projects DROP COLUMN slug;