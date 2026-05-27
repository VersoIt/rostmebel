DROP TRIGGER IF EXISTS trg_project_category_search_refresh ON project_categories;
DROP TRIGGER IF EXISTS trg_project_search_vector ON projects;

DROP FUNCTION IF EXISTS refresh_projects_search_vector_for_category();
DROP FUNCTION IF EXISTS update_project_search_vector();

DROP TABLE IF EXISTS reviews;
DROP TABLE IF EXISTS ip_blocks;
DROP TABLE IF EXISTS admins;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS projects;
DROP TABLE IF EXISTS project_categories;
