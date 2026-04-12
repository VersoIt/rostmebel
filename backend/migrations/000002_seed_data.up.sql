-- Seed Categories
INSERT INTO project_categories (name, slug, icon, sort_order) VALUES
('Кухни', 'kitchens', 'LucideUtensils', 1),
('Шкафы-купе', 'wardrobes', 'LucideBox', 2);

-- Seed Projects
INSERT INTO projects (project_category_id, name, slug, description, price, price_old, images, specs, ai_tags, status) VALUES
(1, 'Кухня "Антрацит Лофт"', 'kitchen-anthracite-loft', 'Современная кухня с фасадами из итальянского пластика Fenix. Столешница из искусственного камня, фурнитура Blum с доводчиками.', 185000, 210000, 
 '[{"url": "https://img.freepik.com/premium-photo/modern-kitchen-furnished-with-white-wooden-cabinets-platform-light_126745-4248.jpg?semt=ais_hybrid", "is_main": true}]',
 '{"Материал": "Fenix NTM", "Фурнитура": "Blum", "Стиль": "Лофт"}', 
 'лофт, антрацит, темная кухня, премиум', 'published'),

(1, 'Скандинавская светлая кухня', 'scandi-white-kitchen', 'Уютная и функциональная кухня в скандинавском стиле. Фасады — матовая эмаль, интегрированные ручки.', 145000, NULL, 
 '[{"url": "https://img.freepik.com/premium-photo/modern-kitchen-furnished-with-white-wooden-cabinets-platform-light_126745-4248.jpg?semt=ais_hybrid", "is_main": true}]',
 '{"Материал": "МДФ Эмаль", "Фурнитура": "Hettich", "Стиль": "Скандинавский"}', 
 'сканди, белая кухня, светлый интерьер', 'published'),

(2, 'Шкаф-купе "Зеркальный минимализм"', 'wardrobe-mirror-minimal', 'Встроенный шкаф-купе во всю стену. Зеркальные полотна визуально расширяют пространство спальни.', 95000, 115000, 
 '[{"url": "https://yandex-images.clstorage.net/5yEC0x183/3aaf52V3Qv/Pxf9RBfb9Vgh2o4Lg0qc5Wk98tC1Ukw5Yv_KhEScyfDxVdhjusm74sX8AcYd1W2vGS5qPOhaqIXPXBTRFV43OH8G_lwMi7oLNo-ZTp0SwPRyJbzRo15TSFPb_-l8tu-U22z_efTW-KSsmxSOa7kmg1tQFGKOKMKrli1hpHBqSR4sBqYIfsPRYTHUnir0RKHBJTTD755ndUVG_dDCQpzcpdNf3kr2vf6Qo5dMk2iUOYhCVzVOmSj1sps7S6FNb2YiKCCTO3nu0mgE3aU29GiN5QwD-dSbXWF1X8D58Wi47rHITtB0wavPoJyCZqdYtxmYTFAuYoEd1aqxCmGTYnhxDigKjTJr-c1pGdeuO9lykc4WJen2uX1NWFr08v5yj_eo33H6U_ua_IedgmGzaJ0mi2NwFEO5MdeMpyZlknlyTj8IKoUGTO7mcSLrlS7JZZPvIwHx0ItdRG963Pj6T4runtFf00z6l_i6iqJdgEywL75geTN4lDvorbghXY1LZk4SAzW1AWvj8moi3pse50Wa8TAN5PS9QV1QVd_c_ke686zoU-Zx3pDes6OWV5dKlhmdWGQ-dZsq946WHn2xV0ZxHCc1qCVP4O5JDMySDstJovI1AeH4iU9ie2fi2N9ksPC4-kLHYcS95qKBjH25ZpQRq29xOUmCM8qytSZft2ZtRw4WMr0Bdsr1dRf6tgv0coLDHwzl_JR6Qlpt1vj_Sb3pp8ht3G7-q9Cpqathr3aVKJ5sUgtuoD7CrIwzb65sSko4BBaCLXfKwHc63Yot4UmV5AML9t2hQk1NX9HV42Gc75n7d_tCzZ_PmYahX5thnR2aTlEiWq0g946GEGiVeUhzLz0qtB921stpNNyQFOZSs9cJK8nUh0RVbUHq9-hlo9a02E_wWe2-wYuah3qfQI4YlGR8MVubMeiCvSBRp1hPVxYHN4w-aOPCUT7iqD7oa7LXDzP0-ZlAWWV20u3ZbIg", "is_main": true}]',
 '{"Профиль": "Aristo", "Наполнение": "Egger", "Тип": "Встроенный"}', 
 'шкаф, зеркало, спальня, минимализм', 'published');
