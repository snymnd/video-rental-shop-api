INSERT INTO videos(title, overview, format, rent_price, production_company, cover_path, total_stock, available_stock, genre_ids) 
VALUES 
('The Matrix', 'A computer hacker learns about reality and his role in the fight against machines.', 'bluray', 4.99, 'Warner Bros.', '/covers/matrix.jpg', 10, 8, ARRAY[1, 8]),
('Inception', 'A thief who steals corporate secrets through dream-sharing technology is given a final job.', 'dvd', 3.99, 'Warner Bros.', '/covers/inception.jpg', 15, 12, ARRAY[1, 8, 5]),
('Pulp Fiction', 'The lives of two mob hitmen, a boxer, a gangster and his wife intertwine.', 'bluray', 4.50, 'Miramax', '/covers/pulpfiction.jpg', 8, 5, ARRAY[1, 5]),
('The Shawshank Redemption', 'Two imprisoned men bond over a number of years.', 'dvd', 3.50, 'Castle Rock Entertainment', '/covers/shawshank.jpg', 12, 10, ARRAY[5]),
('Avengers: Endgame', 'The Avengers take a final stand against Thanos.', 'dvd', 5.99, 'Marvel Studios', '/covers/endgame.jpg', 20, 0, ARRAY[1, 8]),
('Parasite', 'A poor family schemes to become employed by a wealthy family.', 'bluray', 4.99, 'CJ Entertainment', '/covers/parasite.jpg', 7, 0, ARRAY[5, 3]),
('The Lion King', 'A young lion prince flees his kingdom after the murder of his father.', 'dvd', 3.50, 'Walt Disney Pictures', '/covers/lionking.jpg', 18, 15, ARRAY[7]),
('Joker', 'A mentally troubled comedian embarks on a downward spiral.', 'vhs', 5.50, 'Warner Bros.', '/covers/joker.jpg', 10, 9, ARRAY[5]),
('The Godfather', 'The aging patriarch of an organized crime dynasty transfers control to his son.', 'bluray', 4.99, 'Paramount Pictures', '/covers/godfather.jpg', 5, 3, ARRAY[1, 5]),
('Titanic', 'A seventeen-year-old aristocrat falls in love with a kind but poor artist aboard the luxurious, ill-fated R.M.S. Titanic.', 'dvd', 3.99, 'Paramount Pictures', '/covers/titanic.jpg', 15, 12, ARRAY[2, 5]);
