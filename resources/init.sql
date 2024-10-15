INSERT INTO casbin_rule (ptype, v0, v1, v2)
VALUES ('p', 'guest', '/swagger/*', 'GET'),
       ('p', 'ipc_device', '/ipc/ipc_device/upload_image', 'POST');