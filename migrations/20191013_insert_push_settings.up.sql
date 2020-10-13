INSERT INTO public.setting_groups (id, created_at, updated_at, deleted_at, name) VALUES (3, null, null, null, 'PUSH');
INSERT INTO public.settings (id, created_at, updated_at, deleted_at, setting_group_id, key, value, label, description, value_type) VALUES (9, null, null, null, 3, 'privateKey', null, 'Private key', null, 'text');
INSERT INTO public.settings (id, created_at, updated_at, deleted_at, setting_group_id, key, value, label, description, value_type) VALUES (10, null, null, null, 3, 'publicKey', null, 'Public key', null, 'text');
