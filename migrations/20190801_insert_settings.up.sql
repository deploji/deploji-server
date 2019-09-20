INSERT INTO public.setting_groups (id, created_at, updated_at, deleted_at, name) VALUES (1, null, null, null, 'LDAP');
INSERT INTO public.settings (id, created_at, updated_at, deleted_at, setting_group_id, key, value, label, description, value_type) VALUES (1, null, null, null, 1, 'enabled', null, 'LDAP enabled', null, 'bool');
INSERT INTO public.settings (id, created_at, updated_at, deleted_at, setting_group_id, key, value, label, description, value_type) VALUES (2, null, null, null, 1, 'host', null, 'Host', null, 'text');
INSERT INTO public.settings (id, created_at, updated_at, deleted_at, setting_group_id, key, value, label, description, value_type) VALUES (3, null, null, null, 1, 'domain', null, 'Domain', null, 'text');
