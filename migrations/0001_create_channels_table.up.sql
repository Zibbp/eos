CREATE TABLE IF NOT EXISTS channels (
    ID varchar(250) PRIMARY KEY,
    Name text NOT NULL UNIQUE,
    Channel_Image_Path text NOT NULL,
    Created_At TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    Updated_At TIMESTAMPTZ NOT NULL DEFAULT NOW()
)