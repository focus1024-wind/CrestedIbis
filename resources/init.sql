INSERT INTO casbin_rule (ptype, v0, v1, v2)
VALUES ('p', 'guest', '/swagger/*', 'GET'),
       ('p', 'guest', '/system/user/login', 'POST'),
       ('p', 'guest', '/system/user/register', 'POST'),
       ('p', 'ipc_device', '/ipc/ipc_device/upload_image', 'POST');

INSERT INTO sys_users (sys_users.user_id, username, password)
VALUES (1, 'admin', '$2a$10$/wq8mbPLfTdXQV/wvr77YuEim/uxcuZmqXM9TEiBlFD6dtt3aQ/Ha');

INSERT INTO role_groups (role_id, role_name)
VALUES (1, 'admin');

INSERT INTO user_role_groups(sys_user_user_id, role_group_role_id)
VALUES (1, 1);