INSERT INTO roles(id, name, description)
VALUES
(1, 'user', 'default and general role'),
(2, 'admin', 'managing role');

INSERT INTO permissions(id, name, description)
VALUES
(1, 'create', 'resources creation permission'),
(2, 'readOwn', 'own resources read/get permission'),
(3, 'readAll', 'all resources read/get permission'),
(4, 'updateOwn', 'own resources update permission'),
(5, 'updateAll', 'all resources update permission'),
(6, 'deleteOwn', 'own resources deletion permission'),
(7, 'deleteAll', 'own resources deletion permission');

INSERT INTO resources(id, name, description)
VALUES
(1, 'users', 'registered users data'),
(2, 'videos', 'stored video data'),
(3, 'rentals', 'users rental records'),
(4, 'payments', 'users payment records');


INSERT INTO rbac(role_id, permission_id, resource_id)
VALUES
-- user permissions
(1, 2, 1), -- user can read own user profile
(1, 4, 1), -- user can update own user profile
(1, 6, 1), -- user can delete own user profile

(1, 3, 2), -- user can read all videos

(1, 1, 3), -- user can create rental
(1, 2, 3), -- user can read own rental

(1, 1, 4), -- user can create payment
(1, 2, 4), -- user can read own payment

-- admin permissions
(2, 2, 1), -- admin can read own profile
(2, 3, 1), -- admin can read all users
(2, 4, 1), -- admin can update own profile
(2, 5, 1), -- admin can update all users
(2, 7, 1), -- admin can delete all users

(2, 1, 2), -- admin can create videos
(2, 3, 2), -- admin can read all videos
(2, 5, 2), -- admin can update all videos
(2, 7, 2), -- admin can delete all videos

(2, 3, 3), -- admin can read all rentals
(2, 5, 3), -- admin can update all rentals
(2, 7, 3), -- admin can delete all rentals

(2, 3, 4), -- admin can read all payments
(2, 7, 4); -- admin can delete all payments