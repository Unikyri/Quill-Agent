-- Close the universe genre and work format vocabularies, and standardize entity types.

ALTER TABLE universes ADD COLUMN genre_tags TEXT[] NOT NULL DEFAULT '{}';

UPDATE universes
SET genre_tags = CASE genre
    WHEN 'sci-fi' THEN ARRAY['science-fiction']
    WHEN 'fantasy' THEN ARRAY['fantasy']
    WHEN 'mystery' THEN ARRAY['mystery']
    WHEN 'romance' THEN ARRAY['romance']
    WHEN 'horror' THEN ARRAY['horror']
    WHEN 'thriller' THEN ARRAY['thriller']
    WHEN 'historical' THEN ARRAY['historical']
    WHEN 'adventure' THEN ARRAY['adventure']
    WHEN 'non-fiction' THEN ARRAY['literary']
    WHEN 'comedy' THEN ARRAY['literary']
    WHEN 'drama' THEN ARRAY['literary']
    WHEN NULL THEN '{}'
    ELSE '{}'
END;

UPDATE works AS w
SET type = u.format
FROM universes AS u
WHERE w.universe_id = u.id
  AND u.format IN ('novel', 'novella', 'short-story');

UPDATE works
SET type = 'novel'
WHERE type IS NULL OR type NOT IN ('novel', 'novella', 'short-story');

ALTER TABLE works
    ADD CONSTRAINT works_type_check
    CHECK (type IN ('novel', 'novella', 'short-story'));

UPDATE entities SET type = 'world_rule' WHERE type IN ('worldrule', 'rule');
UPDATE entities SET type = 'place' WHERE type = 'location';
UPDATE entities SET type = 'object'
WHERE type NOT IN ('character', 'place', 'object', 'faction', 'event', 'world_rule', 'plot_arc');

ALTER TABLE entities
    ADD CONSTRAINT entities_type_check
    CHECK (type IN ('character', 'place', 'object', 'faction', 'event', 'world_rule', 'plot_arc'));

ALTER TABLE universes DROP COLUMN genre;
ALTER TABLE universes DROP COLUMN format;
