CREATE OR REPLACE FUNCTION update_project_search_vector()
RETURNS TRIGGER AS $$
DECLARE
    category_name TEXT := '';
    specs_text    TEXT := '';
BEGIN
    IF NEW.project_category_id IS NOT NULL THEN
        SELECT name
          INTO category_name
          FROM project_categories
         WHERE id = NEW.project_category_id;
    END IF;

    SELECT COALESCE(string_agg(spec.key || ' ' || spec.value, ' '), '')
      INTO specs_text
      FROM jsonb_each_text(COALESCE(NEW.specs, '{}'::jsonb)) AS spec(key, value);

    NEW.search_vector := to_tsvector('russian',
        COALESCE(NEW.name, '') || ' ' ||
        COALESCE(category_name, '') || ' ' ||
        COALESCE(NEW.description, '') || ' ' ||
        COALESCE(specs_text, '') || ' ' ||
        COALESCE(NEW.ai_tags, '')
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION refresh_projects_search_vector_for_category()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE projects
       SET name = name
     WHERE project_category_id = NEW.id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_project_category_search_refresh ON project_categories;
CREATE TRIGGER trg_project_category_search_refresh
AFTER UPDATE OF name, slug ON project_categories
FOR EACH ROW
WHEN (OLD.name IS DISTINCT FROM NEW.name OR OLD.slug IS DISTINCT FROM NEW.slug)
EXECUTE FUNCTION refresh_projects_search_vector_for_category();

UPDATE projects
   SET name = name;
