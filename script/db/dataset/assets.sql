BEGIN;

-- common
INSERT INTO reward(type, uri, created_at)
VALUES
    ('common', 'assets/common/bubble/1.jpg', NOW()),

    ('common', 'assets/common/caine/1.jpg', NOW()),
    ('common', 'assets/common/caine/2.jpg', NOW()),
    ('common', 'assets/common/caine/3.jpg', NOW()),
    ('common', 'assets/common/caine/4.jpg', NOW()),
    ('common', 'assets/common/caine/5.jpg', NOW()),
    ('common', 'assets/common/caine/6.jpg', NOW()),
    ('common', 'assets/common/caine/7.jpg', NOW()),
    ('common', 'assets/common/caine/8.jpg', NOW()),

    ('common', 'assets/common/gangle/1.jpg', NOW()),
    ('common', 'assets/common/gangle/2.jpg', NOW()),
    ('common', 'assets/common/gangle/3.jpg', NOW()),
    ('common', 'assets/common/gangle/4.jpg', NOW()),
    ('common', 'assets/common/gangle/5.jpg', NOW()),
    ('common', 'assets/common/gangle/6.jpg', NOW()),
    ('common', 'assets/common/gangle/7.jpg', NOW()),
    ('common', 'assets/common/gangle/8.jpg', NOW()),

    ('common', 'assets/common/jax/1.jpg', NOW()),
    ('common', 'assets/common/jax/2.jpg', NOW()),
    ('common', 'assets/common/jax/3.jpg', NOW()),
    ('common', 'assets/common/jax/4.jpg', NOW()),
    ('common', 'assets/common/jax/5.jpg', NOW()),
    ('common', 'assets/common/jax/6.jpg', NOW()),
    ('common', 'assets/common/jax/7.jpg', NOW()),
    ('common', 'assets/common/jax/8.jpg', NOW()),
    ('common', 'assets/common/jax/9.jpg', NOW()),
    ('common', 'assets/common/jax/10.jpg', NOW()),

    ('common', 'assets/common/kinger/1.jpg', NOW()),
    ('common', 'assets/common/kinger/2.jpg', NOW()),
    ('common', 'assets/common/kinger/3.jpg', NOW()),
    ('common', 'assets/common/kinger/4.jpg', NOW()),

    ('common', 'assets/common/pomni/1.jpg', NOW()),
    ('common', 'assets/common/pomni/2.jpg', NOW()),
    ('common', 'assets/common/pomni/3.jpg', NOW()),
    ('common', 'assets/common/pomni/4.jpg', NOW()),
    ('common', 'assets/common/pomni/5.jpg', NOW()),
    ('common', 'assets/common/pomni/6.jpg', NOW()),
    ('common', 'assets/common/pomni/7.jpg', NOW()),
    ('common', 'assets/common/pomni/8.jpg', NOW()),
    ('common', 'assets/common/pomni/9.jpg', NOW()),
    ('common', 'assets/common/pomni/10.jpg', NOW()),
    ('common', 'assets/common/pomni/11.jpg', NOW()),

    ('common', 'assets/common/ragatha/1.jpg', NOW()),
    ('common', 'assets/common/ragatha/2.jpg', NOW()),
    ('common', 'assets/common/ragatha/3.jpg', NOW()),
    ('common', 'assets/common/ragatha/4.jpg', NOW()),
    ('common', 'assets/common/ragatha/5.jpg', NOW()),
    ('common', 'assets/common/ragatha/6.jpg', NOW()),
    ('common', 'assets/common/ragatha/7.jpg', NOW()),
    ('common', 'assets/common/ragatha/8.jpg', NOW()),

    ('common', 'assets/common/zooble/1.jpg', NOW()),
    ('common', 'assets/common/zooble/2.jpg', NOW()),
    ('common', 'assets/common/zooble/3.jpg', NOW())

    ON CONFLICT (uri) DO NOTHING;

-- rare
INSERT INTO reward(type, uri, created_at)
VALUES
    ('rare', 'assets/rare/caine/1.jpg', NOW()),
    ('rare', 'assets/rare/caine/2.jpg', NOW()),
    ('rare', 'assets/rare/caine/3.jpg', NOW()),
    ('rare', 'assets/rare/caine/4.jpg', NOW()),
    ('rare', 'assets/rare/caine/5.jpg', NOW()),

    ('rare', 'assets/rare/gangle/1.jpg', NOW()),
    ('rare', 'assets/rare/gangle/2.jpg', NOW()),
    ('rare', 'assets/rare/gangle/3.jpg', NOW()),
    ('rare', 'assets/rare/gangle/4.jpg', NOW()),

    ('rare', 'assets/rare/jax/1.jpg', NOW()),
    ('rare', 'assets/rare/jax/2.jpg', NOW()),
    ('rare', 'assets/rare/jax/3.jpg', NOW()),
    ('rare', 'assets/rare/jax/4.jpg', NOW()),

    ('rare', 'assets/rare/kinger/1.jpg', NOW()),
    ('rare', 'assets/rare/kinger/2.jpg', NOW()),

    ('rare', 'assets/rare/npc/1.jpg', NOW()),

    ('rare', 'assets/rare/pomni/1.jpg', NOW()),
    ('rare', 'assets/rare/pomni/2.jpg', NOW()),
    ('rare', 'assets/rare/pomni/3.jpg', NOW()),
    ('rare', 'assets/rare/pomni/4.jpg', NOW()),
    ('rare', 'assets/rare/pomni/5.jpg', NOW()),
    ('rare', 'assets/rare/pomni/6.jpg', NOW()),
    ('rare', 'assets/rare/pomni/7.jpg', NOW()),

    ('rare', 'assets/rare/ragatha/1.jpg', NOW()),
    ('rare', 'assets/rare/ragatha/2.jpg', NOW()),
    ('rare', 'assets/rare/ragatha/3.jpg', NOW()),
    ('rare', 'assets/rare/ragatha/4.jpg', NOW()),

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
    ('epic', 'assets/epic/caine/5.jpg', NOW()),

    ('epic', 'assets/epic/gangle/1.jpg', NOW()),

    ('epic', 'assets/epic/jax/1.jpg', NOW()),
    ('epic', 'assets/epic/jax/2.jpg', NOW()),
    ('epic', 'assets/epic/jax/3.jpg', NOW()),

    ('epic', 'assets/epic/kinger/1.jpg', NOW()),

    ('epic', 'assets/epic/npc/1.jpg', NOW()),
    ('epic', 'assets/epic/npc/2.jpg', NOW()),
    ('epic', 'assets/epic/npc/3.jpg', NOW()),

    ('epic', 'assets/epic/pomni/1.jpg', NOW()),
    ('epic', 'assets/epic/pomni/2.jpg', NOW()),
    ('epic', 'assets/epic/pomni/3.jpg', NOW()),
    ('epic', 'assets/epic/pomni/4.jpg', NOW()),
    ('epic', 'assets/epic/pomni/5.jpg', NOW()),
    ('epic', 'assets/epic/pomni/6.jpg', NOW()),

    ('epic', 'assets/epic/ragatha/1.jpg', NOW()),
    ('epic', 'assets/epic/ragatha/2.jpg', NOW()),

    ('epic', 'assets/epic/ribbit/1.jpg', NOW()),

    ('epic', 'assets/epic/zooble/1.jpg', NOW()),
    ('epic', 'assets/epic/zooble/2.jpg', NOW())

    ON CONFLICT (uri) DO NOTHING;

-- legendary
INSERT INTO reward(type, uri, created_at)
VALUES
    ('legendary', 'assets/legendary/caine/1.jpg', NOW()),

    ('legendary', 'assets/legendary/gangle/1.jpg', NOW()),

    ('legendary', 'assets/legendary/jax/1.jpg', NOW()),
    ('legendary', 'assets/legendary/jax/2.jpg', NOW()),

    ('legendary', 'assets/legendary/kinger/1.jpg', NOW()),

    ('legendary', 'assets/legendary/pomni/1.jpg', NOW())

    ON CONFLICT (uri) DO NOTHING;

COMMIT;