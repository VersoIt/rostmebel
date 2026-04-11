-- Seed Categories
INSERT INTO project_categories (name, slug, icon, sort_order) VALUES
('Кухни', 'kitchens', 'LucideUtensils', 1),
('Шкафы-купе', 'wardrobes', 'LucideBox', 2);

-- Seed Projects
INSERT INTO projects (project_category_id, name, slug, description, price, price_old, images, specs, ai_tags, status) VALUES
(1, 'Кухня "Антрацит Лофт"', 'kitchen-anthracite-loft', 'Современная кухня с фасадами из итальянского пластика Fenix. Столешница из искусственного камня, фурнитура Blum с доводчиками.', 185000, 210000, 
 '[{"url": "https://images.unsplash.com/photo-1556911223-05345a39365e", "is_main": true}]', 
 '{"Материал": "Fenix NTM", "Фурнитура": "Blum", "Стиль": "Лофт"}', 
 'лофт, антрацит, темная кухня, премиум', 'published'),

(1, 'Скандинавская светлая кухня', 'scandi-white-kitchen', 'Уютная и функциональная кухня в скандинавском стиле. Фасады — матовая эмаль, интегрированные ручки.', 145000, NULL, 
 '[{"url": "https://images.unsplash.com/photo-1556909114-f6e7ad7d3136", "is_main": true}]', 
 '{"Материал": "МДФ Эмаль", "Фурнитура": "Hettich", "Стиль": "Скандинавский"}', 
 'сканди, белая кухня, светлый интерьер', 'published'),

(2, 'Шкаф-купе "Зеркальный минимализм"', 'wardrobe-mirror-minimal', 'Встроенный шкаф-купе во всю стену. Зеркальные полотна визуально расширяют пространство спальни.', 95000, 115000, 
 '[{"url": "https://images.unsplash.com/photo-1595428774223-ef52624120d2", "is_main": true}]', 
 '{"Профиль": "Aristo", "Наполнение": "Egger", "Тип": "Встроенный"}', 
 'шкаф, зеркало, спальня, минимализм', 'published');
