CREATE TABLE IF NOT EXISTS chapters (
    ID text PRIMARY KEY,
    Start_Time numeric,
    Title text,
    End_Time numeric,
    Video_ID varchar(250) NOT NULL,
    FOREIGN KEY (Video_ID) REFERENCES videos(ID) ON DELETE CASCADE
)