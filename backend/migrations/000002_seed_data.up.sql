-- Seed Categories
INSERT INTO categories (name, slug, icon, sort_order) VALUES
('Диваны', 'sofas', '🛋️', 1),
('Кровати', 'beds', '🛏️', 2),
('Столы', 'tables', '🍽️', 3);

-- Seed Products
INSERT INTO products (category_id, name, slug, description, price, price_old, images, specs, ai_tags, status) VALUES
(1, 'Диван "Скандинавия"', 'sofa-scandinavia', 'Минималистичный диван в сером цвете с деревянными ножками.', 45000, 52000, 
 '[{"url": "https://images.unsplash.com/photo-1555041469-a586c61ea9bc", "is_main": true}]', 
 '{"Материал": "Ткань, Дуб", "Размер": "210х90х85"}', 
 'скандинавский, серый, гостиная, минимализм', 'published'),

(2, 'Кровать "Берген"', 'bed-bergen', 'Просторная двуспальная кровать с мягким изголовьем.', 68000, NULL, 
 '[{"url": "https://images.unsplash.com/photo-1505693419148-ad3035ce6121", "is_main": true}]', 
 '{"Материал": "Велюр", "Размер": "180х200"}', 
 'спальня, уют, современный, светлый', 'published'),

(3, 'Стол "Лофт"', 'table-loft', 'Массивный обеденный стол из цельного дуба на стальном каркасе.', 32000, 38000, 
 '[{"url": "https://images.unsplash.com/photo-1530018607912-eff2df114f11", "is_main": true}]', 
 '{"Материал": "Дуб, Сталь", "Размер": "160х80"}', 
 'кухня, лофт, дерево, обеденный', 'published');
