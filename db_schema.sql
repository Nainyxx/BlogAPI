-- 1. Таблица пользователей
CREATE TABLE users (
                       id UUID PRIMARY KEY,
                       name VARCHAR(100) NOT NULL CHECK (length(name) >= 1),
                       surname VARCHAR(100) NOT NULL CHECK (length(surname) >= 1),
                       login VARCHAR(50) NOT NULL UNIQUE CHECK (length(login) >= 3),
                       email VARCHAR(255) NOT NULL UNIQUE,
                       password_hash VARCHAR(255) NOT NULL CHECK (length(password_hash) > 0),
                       created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_login ON users(login);

-- 2. Таблица постов
CREATE TABLE posts (
                       id UUID PRIMARY KEY,
                       title VARCHAR(255) NOT NULL CHECK (length(title) >= 1),
                       body TEXT NOT NULL CHECK (length(body) >= 1),
                       author_id UUID NOT NULL,
                       created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

                       CONSTRAINT fk_author FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_posts_author_id ON posts(author_id);
CREATE INDEX idx_posts_created_at ON posts(created_at DESC);

-- 3. Таблица картинок к постам
CREATE TABLE images (
                        id UUID PRIMARY KEY,
                        image_url VARCHAR(2048) NOT NULL CHECK (length(image_url) >= 10),
                        post_id UUID NOT NULL,
                        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

                        CONSTRAINT fk_post FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
);

CREATE INDEX idx_images_post_id ON images(post_id);

-- 4. Таблица комментариев
CREATE TABLE comments (
                          id UUID PRIMARY KEY,
                          body TEXT NOT NULL CHECK (length(body) >= 1),
                          author_id UUID NOT NULL,
                          post_id UUID NOT NULL,
                          created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                          updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

                          CONSTRAINT fk_author FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE,
                          CONSTRAINT fk_post FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
);

CREATE INDEX idx_comments_post_id ON comments(post_id);
CREATE INDEX idx_comments_author_id ON comments(author_id);

-- 5. Таблица лайков
CREATE TABLE likes (
                       id UUID PRIMARY KEY,
                       user_id UUID NOT NULL,
                       post_id UUID,
                       comment_id UUID,
                       created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

                       CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
                       CONSTRAINT fk_post FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
                       CONSTRAINT fk_comment FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE,
                       CONSTRAINT check_one_type CHECK (
                           (post_id IS NOT NULL AND comment_id IS NULL) OR
                           (post_id IS NULL AND comment_id IS NOT NULL)
                           ),
                       UNIQUE(user_id, post_id, comment_id)
);

CREATE INDEX idx_likes_user_id ON likes(user_id);
CREATE INDEX idx_likes_post_id ON likes(post_id);
CREATE INDEX idx_likes_comment_id ON likes(comment_id);