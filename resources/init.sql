INSERT INTO casbin_rule (ptype, v0, v1, v2)
VALUES ('p', 'guest', '/swagger/*', 'GET'),
       ('p', 'guest', '/system/user/login', 'POST'),
       ('p', 'guest', '/system/user/register', 'POST'),
       ('p', 'ipc_device', '/ipc/device/devices', 'GET'),
       ('p', 'ipc_device', '/ipc/device/channels', 'GET'),
       ('p', 'ipc_device', '/ipc/device/upload_image', 'POST'),
       ('p', 'guest', '/record/*', 'GET');

INSERT INTO sys_users (sys_users.id, username, password)
VALUES (1, 'admin', '$2a$10$/wq8mbPLfTdXQV/wvr77YuEim/uxcuZmqXM9TEiBlFD6dtt3aQ/Ha');

INSERT INTO role_groups (id, role_name)
VALUES (1, 'admin'),
       (2, 'guest');

INSERT INTO user_role_groups(sys_user_id, role_group_id)
VALUES (1, 1);