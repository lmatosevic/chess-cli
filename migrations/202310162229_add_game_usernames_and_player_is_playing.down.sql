ALTER TABLE game
    DROP COLUMN "whitePlayerUsername",
    DROP COLUMN "blackPlayerUsername";

ALTER TABLE player
    DROP COLUMN "isPlaying";
