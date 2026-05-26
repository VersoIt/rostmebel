DROP TRIGGER IF EXISTS trg_project_category_search_refresh ON project_categories;
DROP FUNCTION IF EXISTS refresh_projects_search_vector_for_category();

CREATE OR REPLACE FUNCTION update_project_search_vector()
RETURNS TRIGGER AS $$
BEGIN
    NEW.search_vector := to_tsvector('russian',
        COALESCE(NEW.name, '') || ' ' ||
        COALESCE(NEW.description, '') || ' ' ||
        COALESCE(NEW.ai_tags, '')
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

UPDATE projects
   SET name = name;
