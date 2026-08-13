INSERT INTO project_categories (name, slug, icon, sort_order)
VALUES
    ('Кухни', 'kitchens', 'LucideUtensils', 1),
    ('Шкафы-купе', 'wardrobes', 'LucideBox', 2),
    ('Гардеробные', 'dressing-rooms', 'LucideShirt', 3),
    ('Прихожие', 'hallways', 'LucideDoorOpen', 4),
    ('Столы', 'tables', 'LucideTable2', 5),
    ('Детские', 'children-rooms', 'LucideBedSingle', 6),
    ('ТВ-тумбы и гостиные', 'living-rooms', 'LucideTv', 7),
    ('Коммерческая мебель', 'commercial-furniture', 'LucideBriefcaseBusiness', 8)
ON CONFLICT (slug) DO UPDATE
SET
    name = EXCLUDED.name,
    icon = EXCLUDED.icon,
    sort_order = EXCLUDED.sort_order,
    updated_at = NOW();
