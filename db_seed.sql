DROP TABLE IF EXISTS comment;
DROP TABLE IF EXISTS post;
DROP TABLE IF EXISTS users;

CREATE TABLE users
(
    id SERIAL PRIMARY KEY,
    name varchar(255)
);

CREATE TABLE post 
(
    id SERIAL PRIMARY KEY,
    author_users_id INTEGER NOT NULL,
    text text,
    created_time TIMESTAMP NOT NULL DEFAULT now(),
    FOREIGN KEY (author_users_id) REFERENCES users(id)
);

CREATE TABLE comment
(
    id SERIAL PRIMARY KEY,
    post_id INTEGER NOT NULL,
    author_users_id INTEGER NOT NULL,
    created_time TIMESTAMP NOT NULL DEFAULT now(),
    text text,
    FOREIGN KEY (post_id) REFERENCES post(id),
    FOREIGN KEY (author_users_id) REFERENCES users(id)
);

insert into users(id, name)
values (1, 'Sam'),
       (2, 'Pat'),
       (3, 'Riley'),
       (4, 'Jamie'),
       (10, 'Me');

insert into post(id, author_users_id, text)
values (10, 1, 'New pair of shoes!'),
       (20, 1, 'I am at the beach!'),
       (30, 2, 'Me and my boo'),
       (40, 3, 'Just emptied my retirement into Dodgecoin');

insert into comment(id, post_id, author_users_id, text)
values (110, 10, 4, 'Love them!'),
       (120, 10, 3, 'I have the same pair!'),
       (130, 10, 1, 'Thank you for the positive vibes!'),

       (140, 20, 4, 'I was there last week.'),
       (150, 20, 1, 'Wow where did you stay?'),
       (160, 20, 4, 'I stayed in an AirBnB off Ocean Avenue.'),

       (170, 30, 3, 'So beautiful!'),

       (180, 40, 1, 'To the moon!'),
       (190, 40, 2, 'Diamond hands go brrrr'),
       (200, 40, 3, 'How do I buy?'),
       (210, 40, 1, 'Go to cryptoscam dot com and give them all your money'),
       (220, 40, 4, 'I mean come on scam is literally in the name, please don''t use them!'),
       (230, 40, 1, 'I think you need to do your own research, I had a lot of good luck with them, they are legit.'),
       (240, 40, 3, 'Wow I just made a ton of money!!');