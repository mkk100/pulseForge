CREATE TABLE users(
    userId int PRIMARY KEY,  
    userName string,
);

CREATE TABLE posts(
    postId int PRIMARY KEY,
    postTitle varchar(255),
    postDescription varchar(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    userId int NOT NULL references user(userId) ON DELETE CASCADE
);

-- on delete cascade will remove the dependent child when the parent is deleted
-- clustered index means sorted in the order of index
-- unclustered index is a structure that points to rows in the table.

CREATE INDEX index_posts_created_at ON posts(created_at DESC);
CREATE INDEX created_at_and_author ON posts(userID, created_at DESC)