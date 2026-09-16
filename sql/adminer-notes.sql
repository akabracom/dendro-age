-- Проверки для защиты лабы 2 в Adminer (SQL command).

SELECT id, title, status, event_year, ring_count, image_url, video_url, creator_id
FROM dendro_scale_patterns;

SELECT * FROM dendro_researchers;
SELECT * FROM dendro_pattern_likes;

-- Логическое удаление через статус (как кнопка на плитке):
-- UPDATE dendro_scale_patterns SET status = 'deleted' WHERE id = 1;

-- Смена полей по теме и проверка, что лайки из м-м видны на плитках:
-- UPDATE dendro_scale_patterns SET event_year = 1976, ring_count = 4 WHERE id = 1;
