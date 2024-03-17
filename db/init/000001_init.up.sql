CREATE TABLE Users
(
    id SERIAL PRIMARY KEY,
    username VARCHAR NOT NULL UNIQUE,
    password_hash VARCHAR NOT NULL, 
    is_admin BOOLEAN NOT NULL
);

CREATE TABLE Actors 
(
    id SERIAL PRIMARY KEY,
    first_name VARCHAR NOT NULL,
    last_name VARCHAR NOT NULL,
    gender VARCHAR NOT NULL,
    date_of_birth DATE NOT NULL
);

CREATE TABLE Movies 
(
    id SERIAL PRIMARY KEY,
    title VARCHAR(150) NOT NULL,
    description VARCHAR(1000) NOT NULL,
    release_date DATE NOT NULL,
    rating INT NOT NULL CHECK (rating >= 0 AND rating <= 10)
);

CREATE TABLE MoviesActors 
(
    actor_id INTEGER,
    movie_id INTEGER,
    FOREIGN KEY (actor_id) REFERENCES Actors(id),
    FOREIGN KEY (movie_id) REFERENCES Movies(id),
    PRIMARY KEY (actor_id, movie_id)
);
