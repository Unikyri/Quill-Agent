ALTER TABLE universes ADD COLUMN genre VARCHAR(50);
UPDATE universes SET genre = genre_tags[1] WHERE cardinality(genre_tags) > 0;

ALTER TABLE universes ADD COLUMN format VARCHAR(50);
UPDATE universes AS u
SET format = COALESCE((
    SELECT mode() WITHIN GROUP (ORDER BY w.type)
    FROM works AS w
    WHERE w.universe_id = u.id
), 'novel');
ALTER TABLE universes ALTER COLUMN format SET NOT NULL;

ALTER TABLE entities DROP CONSTRAINT entities_type_check;
ALTER TABLE works DROP CONSTRAINT works_type_check;
ALTER TABLE universes DROP COLUMN genre_tags;
