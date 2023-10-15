CREATE TABLE "player"
(
    "id"           SERIAL                 NOT NULL,
    "username"     character varying(250) NOT NULL,
    "passwordHash" character varying      NOT NULL,
    "wins"         integer                NOT NULL DEFAULT 0,
    "losses"       integer                NOT NULL DEFAULT 0,
    "draws"        integer                NOT NULL DEFAULT 0,
    "rate"         float                  NOT NULL DEFAULT 0,
    "elo"          integer                NOT NULL DEFAULT 1000,
    "lastPlayedAt" TIMESTAMP              NULL,
    "createdAt"    TIMESTAMP              NOT NULL DEFAULT (now() at time zone 'utc'),
    "updatedAt"    TIMESTAMP              NOT NULL DEFAULT (now() at time zone 'utc'),
    CONSTRAINT "PK_player_id" PRIMARY KEY ("id"),
    CONSTRAINT "UQ_player_username" UNIQUE ("username")
);

CREATE TABLE "game"
(
    "id"                  SERIAL                NOT NULL,
    "name"                character varying     NOT NULL,
    "passwordHash"        character varying     NULL,
    "turnDurationSeconds" integer               NULL,
    "whitePlayerId"       integer               NULL,
    "blackPlayerId"       integer               NULL,
    "creatorId"           integer               NOT NULL,
    "winnerId"            integer               NULL,
    "tiles"               character varying(64) NOT NULL,
    "inProgress"          boolean               NOT NULL DEFAULT false,
    "lastMovePlayedAt"    TIMESTAMP             NULL,
    "startedAt"           TIMESTAMP             NULL,
    "endedAt"             TIMESTAMP             NULL,
    "createdAt"           TIMESTAMP             NOT NULL DEFAULT (now() at time zone 'utc'),
    "updatedAt"           TIMESTAMP             NOT NULL DEFAULT (now() at time zone 'utc'),
    CONSTRAINT "PK_game_id" PRIMARY KEY ("id"),
    CONSTRAINT "FK_game_creator_id" FOREIGN KEY ("creatorId") REFERENCES "player" ("id") ON DELETE SET NULL ON UPDATE NO ACTION,
    CONSTRAINT "FK_game_winner_id" FOREIGN KEY ("winnerId") REFERENCES "player" ("id") ON DELETE SET NULL ON UPDATE NO ACTION,
    CONSTRAINT "FK_game_white_player_id" FOREIGN KEY ("whitePlayerId") REFERENCES "player" ("id") ON DELETE SET NULL ON UPDATE NO ACTION,
    CONSTRAINT "FK_game_black_player_id" FOREIGN KEY ("blackPlayerId") REFERENCES "player" ("id") ON DELETE SET NULL ON UPDATE NO ACTION
);

CREATE TABLE "game_move"
(
    "id"        SERIAL               NOT NULL,
    "gameId"    integer              NOT NULL,
    "playerId"  integer              NULL,
    "move"      character varying(8) NOT NULL,
    "createdAt" TIMESTAMP            NOT NULL DEFAULT (now() at time zone 'utc'),
    "updatedAt" TIMESTAMP            NOT NULL DEFAULT (now() at time zone 'utc'),
    CONSTRAINT "PK_game_move_id" PRIMARY KEY ("id"),
    CONSTRAINT "FK_game_move_game_id" FOREIGN KEY ("gameId") REFERENCES "game" ("id") ON DELETE CASCADE ON UPDATE NO ACTION,
    CONSTRAINT "FK_game_move_player_id" FOREIGN KEY ("playerId") REFERENCES "player" ("id") ON DELETE SET NULL ON UPDATE NO ACTION
);

CREATE TABLE "access_token"
(
    "id"        SERIAL                NOT NULL,
    "playerId"  integer               NULL,
    "token"     character varying(36) NOT NULL,
    "createdAt" TIMESTAMP             NOT NULL DEFAULT (now() at time zone 'utc'),
    "updatedAt" TIMESTAMP             NOT NULL DEFAULT (now() at time zone 'utc'),
    CONSTRAINT "PK_access_token_id" PRIMARY KEY ("id"),

    CONSTRAINT "UQ_access_token_token" UNIQUE ("token"),
    CONSTRAINT "FK_access_token_player_id" FOREIGN KEY ("playerId") REFERENCES "player" ("id") ON DELETE CASCADE ON UPDATE NO ACTION
);