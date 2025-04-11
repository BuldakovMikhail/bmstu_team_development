-- Проверка пользователей
SELECT *
FROM users;

-- Проверка задач
SELECT t.*, u.user_name
FROM task t
         JOIN users u ON t.user_id = u.id_user;

-- Проверка категорий
SELECT c.*, u.user_name
FROM category c
         JOIN users u ON c.user_id = u.id_user;

-- Проверка связей задач и категорий
SELECT tc.task_id,
       t.title AS task_title,
       tc.category_id,
       c.name  AS category_name,
       u.user_name
FROM task_category tc
         JOIN task t ON tc.task_id = t.id_task
         JOIN category c ON tc.category_id = c.id_category
         JOIN users u ON t.user_id = u.id_user;