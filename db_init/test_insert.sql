-- Тестовые данные
-- Удаляем существующие данные (на случай перезапуска)
DELETE
FROM task_category;
DELETE
FROM task;
DELETE
FROM category;
DELETE
FROM users;

-- Вставляем пользователей
INSERT INTO users (user_name, password_hash)
VALUES ('alice', 'hash1'),
       ('bob', 'hash2') RETURNING id_user, user_name;

-- Запоминаем ID пользователей
WITH user_ids AS (SELECT id_user
                  FROM users
                  WHERE user_name IN ('alice', 'bob'))
-- Вставляем задачи
INSERT
INTO task (user_id, title, description, is_done) VALUES
                                                            ((SELECT id_user FROM users WHERE user_name = 'alice'), 'Купить продукты', 'Молоко, хлеб, яйца', false),
                                                            ((SELECT id_user FROM users WHERE user_name = 'alice'), 'Сделать ДЗ', 'Математика и физика', true),
                                                            ((SELECT id_user FROM users WHERE user_name = 'bob'), 'Починить кран', NULL, false)
RETURNING id_task, title;

-- Вставляем категории
INSERT INTO category (user_id, name)
VALUES ((SELECT id_user FROM users WHERE user_name = 'alice'), 'Домашние дела'),
       ((SELECT id_user FROM users WHERE user_name = 'alice'), 'Учеба'),
       ((SELECT id_user FROM users WHERE user_name = 'bob'), 'Ремонт') RETURNING id_category, name;

-- Связываем задачи с категориями
INSERT INTO task_category (task_id, category_id)
VALUES ((SELECT id_task FROM task WHERE title = 'Купить продукты'),
        (SELECT id_category FROM category WHERE name = 'Домашние дела')),

       ((SELECT id_task FROM task WHERE title = 'Сделать ДЗ'),
        (SELECT id_category FROM category WHERE name = 'Учеба')),

       ((SELECT id_task FROM task WHERE title = 'Починить кран'),
        (SELECT id_category FROM category WHERE name = 'Ремонт'));