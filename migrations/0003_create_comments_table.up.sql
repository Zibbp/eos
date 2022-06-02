CREATE TABLE IF NOT EXISTS comments (
    ID text PRIMARY KEY,
    Text text NOT NULL,
    Timestamp bigint NOT NULL,
    Like_Count integer NOT NULL,
    Is_Favorited boolean,
    -- Some names are over 250 chars long??
    Author text NOT NULL,
    Author_ID varchar(250) NOT NULL,
    Author_Thumbnail text,
    Author_Is_Uploader boolean,
    Parent varchar(250) NOT NULL,
    Video_ID varchar(250) NOT NULL,
    FOREIGN KEY (Video_ID) REFERENCES videos(ID) ON DELETE CASCADE
)