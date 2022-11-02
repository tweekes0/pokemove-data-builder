CREATE TABLE pokemon (
    id SERIAL NOT NULL,
    poke_id INTEGER UNIQUE CONSTRAINT "unique poke_id" PRIMARY KEY,
    origin_gen INTEGER NOT NULL,
    name TEXT NOT NULL,
    sprite TEXT,
    species TEXT NOT NULL
);

CREATE TABLE pokemon_moves (
    id SERIAL PRIMARY KEY, 
    name TEXT NOT NULL,
    move_id INTEGER NOT NULL,
    accuracy INTEGER NOT NULL,
    power INTEGER NOT NULL,
    power_points INTEGER NOT NULL,
    generation INTEGER NOT NULL,
    type TEXT NOT NULL,
    damage_type TEXT NOT NULL,
    description TEXT NOT NULL 
);

CREATE TABLE pokemon_abilities (
    id SERIAL PRIMARY KEY,
    ability_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    generation INTEGER NOT NULL
);

CREATE TABLE pokemon_move_rels(
    id SERIAL PRIMARY KEY,
    poke_id INTEGER NOT NULL,
    move_id INTEGER NOT NULL,
    generation INTEGER NOT NULL, 
    level_learned INTEGER NOT NULL,
    learn_method TEXT NOT NULL,
    game_name TEXT NOT NULL
);

CREATE TABLE pokemon_ability_rels (
    id SERIAL PRIMARY KEY,
    poke_id INTEGER NOT NULL,
    ability_id INTEGER NOT NULL,
    slot INTEGER NOT NULL,
    hidden BOOLEAN NOT NULL
);