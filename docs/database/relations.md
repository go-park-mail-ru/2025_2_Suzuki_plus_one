# Нормализация и функциональные зависимости

В этом документе перечислены ключевые отношения схемы PopFilms (см. `docs/database/base.md`) и заданы их функциональные зависимости (ФЗ). После этого показано, что каждое отношение соответствует 1NF, 2NF, 3NF и (по возможности) BCNF.

Обозначения: фигурные скобки { } — набор атрибутов; запись A -> B означает, что A функционально определяет B.

Важно: для компактности приведены основные (не тривиальные) отношения: `user`, `playlist`, `playlist_media`, `media`, `media_episode`, `media_genre`, `asset`, `asset_video`, `asset_image`, `media_video`, `media_video_asset`, `asset_audio`, `media_audio`, `asset_subtitle`, `media_subtitle`, `actor`, `actor_role`, `user_watch_history`, `user_like_media`, `user_comment_media`, `user_rating_media`, `saved_media`.

Каждое отношение сопровождается списком атрибутов, ФЗ и кратким объяснением нормализации.

## Relation: user

Атрибуты:
{user_id, username, email, password_hash, created_at, updated_at}

Функциональные зависимости:
{user_id} -> username, email, password_hash, created_at, updated_at
{username} -> user_id, email, password_hash, created_at, updated_at
{email} -> user_id, username, password_hash, created_at, updated_at

Пояснение по нормальным формам:
- 1NF: все атрибуты атомарны (текстовые, таймстемпы). Нет повторяющихся групп.
- 2NF: ключ отношения — {user_id} (одноатрибутный), поэтому нет частичных зависимостей.
- 3NF: все неключевые атрибуты функционально зависят от ключа и не зависят транзитивно друг от друга (username/email — также кандидаты, но они уникальны и определяют остальные атрибуты; это допустимо). Нет нетривиальных зависимостей неключевых атрибутов друг от друга.
- BCNF: каждая детерминирующая сторона любой нетривиальной ФЗ — суперключ (username и email являются кандидатом/суперключом из-за уникальности), значит отношение удовлетворяет BCNF.

## Relation: playlist

Атрибуты:
{playlist_id, user_id, name, description, visibility, created_at, updated_at}

ФЗ:
{playlist_id} -> user_id, name, description, visibility, created_at, updated_at

Пояснение:
- 1NF: все атомарно.
- 2NF: ключ — {playlist_id}, нет частичных зависимостей.
- 3NF: все неключевые атрибуты напрямую зависят от ключа; Нет транзитивных зависимостей.
- BCNF: единственная детерминирующая сторона — {playlist_id} (суперключ), поэтому BCNF выполнен.

## Relation: playlist_media

Атрибуты:
{playlist_id, media_id, created_at, updated_at}

ФЗ:
{playlist_id, media_id} -> created_at, updated_at

Пояснение:
- 1NF: атомарно.
- 2NF: составной ключ {playlist_id, media_id} — нет частичных зависимостей, т.к. created_at/updated_at зависят от комбинации (вставки в плейлист).
- 3NF: нет транзитивных зависимостей.
- BCNF: детерминанты — только ключ; BCNF выполнен.

## Relation: media

Атрибуты:
{media_id, media_type, title, description, release_date, created_at, updated_at}

ФЗ:
{media_id} -> media_type, title, description, release_date, created_at, updated_at

Пояснение:
- 1NF: атомарно.
- 2NF: ключ — {media_id}, нет частичных зависимостей.
- 3NF: неключевые атрибуты зависят только от ключа; возможные кандидаты (например, title) не детерминируют остальные атрибуты в общем случае.
- BCNF: детерминанты — суперключи (media_id), BCNF выполнен.

## Relation: media_genre

Атрибуты:
{media_genre_id, media_id, name, created_at, updated_at}

ФЗ:
{media_genre_id} -> media_id, name, created_at, updated_at
{media_id, name} -> media_genre_id, created_at, updated_at

Пояснение:
- 1NF: атомарно.
- 2NF: ключ — {media_genre_id} (одноатрибутный), поэтому 2NF выполнен.
- 3NF: name — это жанр; предполагается, что одна запись media_genre описывает принадлежность медиа к жанру. Нет транзитивных зависимостей.
- BCNF: любые детерминанты — суперключи (media_genre_id или составной {media_id,name}), BCNF выполнен.

## Relation: media_episode

Атрибуты:
{episode_id, series_id, season_number, episode_number, created_at, updated_at}

ФЗ:
{episode_id} -> series_id, season_number, episode_number, created_at, updated_at
{series_id, season_number, episode_number} -> episode_id, created_at, updated_at

Пояснение:
- 1NF: атомарно.
- 2NF: ключ — {episode_id}, частичных зависимостей нет.
- 3NF: все неключевые атрибуты зависят непосредственно от ключа; вторичная зависимость {series_id, season_number, episode_number} -> episode_id отражает естественную уникальность эпизода в рамках сериала — это кандидатный ключ, если модельирует уникальность.
- BCNF: возможная проблема — если мы считаем, что {series_id, season_number, episode_number} детерминирует episode_id, то детерминант является кандидатом и, следовательно, BCNF выполняется. Мы добавили проверки (CHECK/EXISTS) чтобы гарантировать соответствие типов media.

## Relation: asset

Атрибуты:
{asset_id, s3_key, mime_type, size_mb, created_at, updated_at}

ФЗ:
{asset_id} -> s3_key, mime_type, size_mb, created_at, updated_at
{s3_key} -> asset_id, mime_type, size_mb, created_at, updated_at

Пояснение:
- 1NF: атомарно.
- 2NF: ключ — {asset_id}, нет частичных зависимостей.
- 3NF: mime_type и size_mb напрямую зависят от asset_id; s3_key уникален и потому является кандидатом — BCNF выполнен.

## Relation: asset_image

Атрибуты:
{asset_image_id, asset_id, resolution_width, resolution_height, created_at, updated_at}

ФЗ:
{asset_image_id} -> asset_id, resolution_width, resolution_height, created_at, updated_at
{asset_id} -> asset_image_id, resolution_width, resolution_height, created_at, updated_at

Пояснение:
- 1NF: атомарно.
- 2NF: ключ — {asset_image_id}, нет частичных зависимостей.
- 3NF: атрибуты зависят от ключа; связь 1:1 с `asset` обеспечена уникальностью asset_id.
- BCNF: детерминанты (asset_image_id и asset_id) являются суперключами, BCNF выполнен.

## Relation: asset_video

Атрибуты:
{asset_video_id, asset_id, quality, resolution_width, resolution_height, created_at, updated_at}

ФЗ:
{asset_video_id} -> asset_id, quality, resolution_width, resolution_height, created_at, updated_at
{asset_id} -> asset_video_id, quality, resolution_width, resolution_height, created_at, updated_at

Пояснение:
- Аналогично `asset_image`: атомарность, отсутствие частичных зависимостей, BCNF выполнен благодаря 1:1 с `asset`.

## Relation: media_video

Атрибуты:
{media_video_id, media_id, video_type, created_at, updated_at}

ФЗ:
{media_video_id} -> media_id, video_type, created_at, updated_at
{media_id, video_type} -> media_video_id, created_at, updated_at  -- если в модели допускается только один media_video данного типа для media

Пояснение:
- 1NF/2NF/3NF: media_video_id — ключ; нет частичных или транзитивных зависимостей.
- BCNF: детерминанты — суперключи (media_video_id или, при ограничении, {media_id,video_type}).

## Relation: media_video_asset

Атрибуты:
{media_video_id, asset_video_id, created_at, updated_at}

ФЗ:
{media_video_id, asset_video_id} -> created_at, updated_at

Пояснение: типичное соединительное отношение (many-to-many); ключ — составной; удовлетворяет 1NF/2NF/3NF/BCNF.

## Relation: asset_audio

Атрибуты:
{asset_audio_id, asset_id, language, created_at, updated_at}

ФЗ:
{asset_audio_id} -> asset_id, language, created_at, updated_at
{asset_id} -> asset_audio_id, language, created_at, updated_at

Пояснение: 1:1 с asset, нормальные формы соблюдены.

## Relation: media_audio

Атрибуты:
{media_video_id, asset_audio_id, created_at, updated_at}

ФЗ:
{media_video_id, asset_audio_id} -> created_at, updated_at

Пояснение: соединительная таблица, 1NF/2NF/3NF/BCNF выполнены.

## Relation: asset_subtitle

Атрибуты:
{asset_subtitle_id, asset_id, language, format, created_at, updated_at}

ФЗ:
{asset_subtitle_id} -> asset_id, language, format, created_at, updated_at
{asset_id} -> asset_subtitle_id, language, format, created_at, updated_at

Пояснение: 1:1 с asset, нормальные формы соблюдены.

## Relation: media_subtitle

Атрибуты:
{media_video_id, asset_subtitle_id, created_at, updated_at}

ФЗ:
{media_video_id, asset_subtitle_id} -> created_at, updated_at

Пояснение: соединительное отношение, все нормальные формы выполнены.

## Relation: actor

Атрибуты:
{actor_id, name, birth_date, bio, created_at, updated_at}

ФЗ:
{actor_id} -> name, birth_date, bio, created_at, updated_at
{name, birth_date} -> actor_id, bio, created_at, updated_at  -- при предположении уникальности имени + даты рождения

Пояснение: 1NF/2NF/3NF/BCNF при условии выбранных кандидатов-ключей.

## Relation: actor_role

Атрибуты:
{actor_role_id, actor_id, media_id, role_name, created_at, updated_at}

ФЗ:
{actor_role_id} -> actor_id, media_id, role_name, created_at, updated_at
{actor_id, media_id} -> actor_role_id, role_name, created_at, updated_at

Пояснение: многие-ко-многим (актер участвует в нескольких медиа). Ключ — actor_role_id или составной {actor_id,media_id} — нормальные формы соблюдены.

## Relation: user_watch_history

Атрибуты:
{watch_history_id, user_id, media_id, watched_at, progress_seconds, created_at, updated_at}

ФЗ:
{watch_history_id} -> user_id, media_id, watched_at, progress_seconds, created_at, updated_at
{user_id, media_id, watched_at} -> watch_history_id, progress_seconds, created_at, updated_at  -- если требуется уникальность записи по времени

Пояснение: ключ — watch_history_id; отсутствие частичных/транзитивных зависимостей; 1NF..BCNF выполнены.

## Relation: user_like_media

Атрибуты:
{user_id, media_id, created_at, updated_at}

ФЗ:
{user_id, media_id} -> created_at, updated_at

Пояснение: составной ключ; соединительная таблица; 1NF/2NF/3NF/BCNF выполнены.

## Relation: user_comment_media

Атрибуты:
{user_comment_id, user_id, media_id, content, created_at, updated_at}

ФЗ:
{user_comment_id} -> user_id, media_id, content, created_at, updated_at

Пояснение: ключ — user_comment_id; нормальные формы выполнены.

## Relation: user_rating_media

Атрибуты:
{user_id, media_id, rating, created_at, updated_at}

ФЗ:
{user_id, media_id} -> rating, created_at, updated_at

Пояснение: составной ключ; рейтинг непосредственно зависит от ключа; 1NF..BCNF выполнены.

## Relation: saved_media

Атрибуты:
{saved_id, user_id, media_id, created_at, updated_at}

ФЗ:
{saved_id} -> user_id, media_id, created_at, updated_at
{user_id, media_id} -> saved_id, created_at, updated_at

Пояснение: ключ — saved_id или составной {user_id,media_id}; нормальные формы выполнены.

---

# Общее доказательство нормализации (1NF, 2NF, 3NF, BCNF)

1. Первая нормальная форма (1NF):
   - Все отношения проектированы с атомарными атрибутами (нет множественных значений в одном атрибуте). Списки жанров, медиа в плейлисте и т.п. представлены отдельными записями в соответствующих соединительных таблицах (`media_genre`, `playlist_media`). Таким образом 1NF выполняется во всей схеме.

2. Вторая нормальная форма (2NF):
   - Все отношения с составными ключами (например, `playlist_media`, `media_video_asset`, `media_audio` и т.д.) не содержат частичных зависимостей: все неключевые атрибуты зависят от всей комбинации ключа. Отношения с единичными ключами (например, `media`, `asset`, `user`) автоматически удовлетворяют 2NF.

3. Третья нормальная форма (3NF):
   - Мы исключили транзитивные зависимости, переносив повторяющиеся и зависящие данные в отдельные таблицы: файловые метаданные хранятся в `asset`, типизированная мета-информация — в `asset_video`, `asset_image`, `asset_audio`, `asset_subtitle`. Информация об актерах, ролях, комментариях, лайках и плейлистах — в собственных таблицах. Следовательно, неключевые атрибуты зависят только от ключа отношения.

4. BCNF (Нормальная форма Бойса-Кодда):
   - Во всех отношениях детерминанты (левая часть невырожденных ФЗ) либо являются ключами, либо мы явно сделали кандидатные ключи (например, сочетание {media_id, name} в `media_genre`). В случаях, где семантика требует специальных проверок (например, гарантировать, что `media_episode.series_id` указывает на запись `media` с `media_type='series'`), в миграции добавлены контрольные проверки и/или это оставлено на логику приложения/триггеры. В целом схема спроектирована так, чтобы детерминанты были суперключами, следовательно BCNF выполняется.

---

Если нужны строгие доказательства для каждой таблицы (например, формальные множества всех возможных ФЗ, включая тривиальные и производные), могу развернуть документ и дать полные вычисления замыкания множеств атрибутов и привести кандидаты-ключи для каждой таблицы. Также могу добавить SQL-триггеры или ограничители для обеспечения семантических ограничений уровня БД (например, гарантировать media_type для series/episode), если вы хотите усилить целостность на уровне СУБД.

Если всё устроит — помечу задачу документирования как завершённую.