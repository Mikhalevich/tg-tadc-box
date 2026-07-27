BEGIN;

-- common
INSERT INTO reward(type, uri, created_at)
VALUES
    ('common', 'assets/common/bubble/1.jpg', NOW()),

    ('common', 'assets/common/caine/1.jpg', NOW()),
    ('common', 'assets/common/caine/2.jpg', NOW()),
    ('common', 'assets/common/caine/3.jpg', NOW()),
    ('common', 'assets/common/caine/4.jpg', NOW()),

    ('common', 'assets/common/gangle/1.jpg', NOW()),
    ('common', 'assets/common/gangle/2.jpg', NOW()),
    ('common', 'assets/common/gangle/3.jpg', NOW()),

    ('common', 'assets/common/jax/1.jpg', NOW()),
    ('common', 'assets/common/jax/2.jpg', NOW()),
    ('common', 'assets/common/jax/3.jpg', NOW()),
    ('common', 'assets/common/jax/4.jpg', NOW()),
    ('common', 'assets/common/jax/5.jpg', NOW()),

    ('common', 'assets/common/kinger/1.jpg', NOW()),

    ('common', 'assets/common/pomni/1.jpg', NOW()),
    ('common', 'assets/common/pomni/2.jpg', NOW()),
    ('common', 'assets/common/pomni/3.jpg', NOW()),
    ('common', 'assets/common/pomni/4.jpg', NOW()),

    ('common', 'assets/common/ragatha/1.jpg', NOW()),
    ('common', 'assets/common/ragatha/2.jpg', NOW()),

    ('common', 'assets/common/zooble/1.jpg', NOW())

    ON CONFLICT (uri) DO NOTHING;

-- rare
INSERT INTO reward(type, uri, created_at)
VALUES
    ('rare', 'assets/rare/caine/1.jpg', NOW()),
    ('rare', 'assets/rare/caine/2.jpg', NOW()),

    ('rare', 'assets/rare/gangle/1.jpg', NOW()),
    ('rare', 'assets/rare/gangle/2.jpg', NOW()),

    ('rare', 'assets/rare/jax/1.jpg', NOW()),
    ('rare', 'assets/rare/jax/2.jpg', NOW()),

    ('rare', 'assets/rare/kinger/1.jpg', NOW()),
    ('rare', 'assets/rare/kinger/2.jpg', NOW()),

    ('rare', 'assets/rare/npc/1.jpg', NOW()),

    ('rare', 'assets/rare/pomni/1.jpg', NOW()),
    ('rare', 'assets/rare/pomni/2.jpg', NOW()),
    ('rare', 'assets/rare/pomni/3.jpg', NOW()),
    ('rare', 'assets/rare/pomni/4.jpg', NOW()),

    ('rare', 'assets/rare/ragatha/1.jpg', NOW()),
    ('rare', 'assets/rare/ragatha/2.jpg', NOW()),

    ('rare', 'assets/rare/zooble/1.jpg', NOW()),
    ('rare', 'assets/rare/zooble/2.jpg', NOW())

    ON CONFLICT (uri) DO NOTHING;

-- epic
INSERT INTO reward(type, uri, created_at)
VALUES
    ('epic', 'assets/epic/caine/1.jpg', NOW()),
    ('epic', 'assets/epic/caine/2.jpg', NOW()),
    ('epic', 'assets/epic/caine/3.jpg', NOW()),
    ('epic', 'assets/epic/caine/4.jpg', NOW()),

    ('epic', 'assets/epic/jax/1.jpg', NOW()),
    ('epic', 'assets/epic/jax/2.jpg', NOW()),
    ('epic', 'assets/epic/jax/3.jpg', NOW()),

    ('epic', 'assets/epic/kinger/1.jpg', NOW()),

    ('epic', 'assets/epic/pomni/1.jpg', NOW()),
    ('epic', 'assets/epic/pomni/2.jpg', NOW()),
    ('epic', 'assets/epic/pomni/3.jpg', NOW()),
    ('epic', 'assets/epic/pomni/4.jpg', NOW()),
    ('epic', 'assets/epic/pomni/5.jpg', NOW()),
    ('epic', 'assets/epic/pomni/6.jpg', NOW()),

    ('epic', 'assets/epic/ragatha/1.jpg', NOW()),
    ('epic', 'assets/epic/ragatha/2.jpg', NOW()),

    ('epic', 'assets/epic/ribbit/1.jpg', NOW()),

    ('epic', 'assets/epic/zooble/1.jpg', NOW())

    ON CONFLICT (uri) DO NOTHING;

-- legendary
INSERT INTO reward(type, uri, created_at)
VALUES
    ('legendary', 'assets/legendary/caine/1.jpg', NOW()),

    ('legendary', 'assets/legendary/gangle/1.jpg', NOW()),

    ('legendary', 'assets/legendary/jax/1.jpg', NOW()),
    ('legendary', 'assets/legendary/jax/2.jpg', NOW()),

    ('legendary', 'assets/legendary/pomni/1.jpg', NOW())

    ON CONFLICT (uri) DO NOTHING;

COMMIT;