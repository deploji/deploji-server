INSERT INTO public.setting_groups (id, created_at, updated_at, deleted_at, name) VALUES (2, null, null, null, 'SMTP');
INSERT INTO public.settings (id, created_at, updated_at, deleted_at, setting_group_id, key, value, label, description, value_type) VALUES (4, null, null, null, 2, 'host', null, 'Host', null, 'text');
INSERT INTO public.settings (id, created_at, updated_at, deleted_at, setting_group_id, key, value, label, description, value_type) VALUES (5, null, null, null, 2, 'port', null, 'Port', null, 'text');
INSERT INTO public.settings (id, created_at, updated_at, deleted_at, setting_group_id, key, value, label, description, value_type) VALUES (6, null, null, null, 2, 'username', null, 'Username', null, 'text');
INSERT INTO public.settings (id, created_at, updated_at, deleted_at, setting_group_id, key, value, label, description, value_type) VALUES (7, null, null, null, 2, 'password', null, 'Password', null, 'text');
INSERT INTO public.settings (id, created_at, updated_at, deleted_at, setting_group_id, key, value, label, description, value_type) VALUES (8, null, null, null, 2, 'sender', null, 'Sender email', null, 'text');
