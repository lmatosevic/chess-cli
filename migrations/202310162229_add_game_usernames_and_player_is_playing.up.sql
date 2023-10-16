ALTER TABLE game
    ADD COLUMN "whitePlayerUsername" character varying NULL,
    ADD COLUMN "blackPlayerUsername" character varying NULL;

UPDATE game
SET "whitePlayerUsername" = (SELECT player."username" FROM player WHERE player.id = game."whitePlayerId")
WHERE game."whitePlayerId" > 0;

UPDATE game
SET "blackPlayerUsername" = (SELECT player."username" FROM player WHERE player.id = game."blackPlayerId")
WHERE game."blackPlayerId" > 0;


ALTER TABLE player
    ADD COLUMN "isPlaying" boolean NOT NULL DEFAULT false;

UPDATE player
SET "isPlaying" = true
WHERE (SELECT COUNT(*)
       FROM game
       WHERE (player.id = game."whitePlayerId" OR
              player.id = game."blackPlayerId")
         AND game."inProgress") > 0;