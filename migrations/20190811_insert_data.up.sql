INSERT INTO public.projects (id, created_at, updated_at, deleted_at, name, repo_url, repo_branch, repo_user, ssh_key_id)
VALUES (1, '2019-08-11 20:56:42.911997', '2019-08-11 20:56:42.911997', null, 'Demo project',
        'git@github.com:deploji/ansible-playbooks.git', 'master', 'git', 0);
INSERT INTO public.inventories (id, created_at, updated_at, deleted_at, name, project_id, source_file)
VALUES (1, '2019-08-11 20:57:08.641661', '2019-08-11 20:57:08.641661', null, 'staging', 1, 'staging.ini');
INSERT INTO public.repositories (id, created_at, updated_at, deleted_at, name, type, url, username, password)
VALUES (1, '2019-08-11 20:58:10.044090', '2019-08-11 20:58:10.044090', null, 'Docker Hub', 'docker-v1',
        'https://registry.hub.docker.com', '', '');
INSERT INTO public.applications (id, created_at, updated_at, deleted_at, name, ansible_name, project_id, repository_id,
                                 repository_artifact, ansible_playbook)
VALUES (2, '2019-08-11 20:59:33.174062', '2019-08-11 21:10:50.605307', null, 'Deploji server', 'deploji-server',
        1, 1, 'deploji/deploji-server', 'deploy.yml');
INSERT INTO public.applications (id, created_at, updated_at, deleted_at, name, ansible_name, project_id, repository_id,
                                 repository_artifact, ansible_playbook)
VALUES (3, '2019-08-11 21:00:01.407079', '2019-08-11 21:10:55.794526', null, 'Deploji worker', 'deploji-worker',
        1, 1, 'deploji/deploji-worker', 'deploy.yml');
INSERT INTO public.applications (id, created_at, updated_at, deleted_at, name, ansible_name, project_id, repository_id,
                                 repository_artifact, ansible_playbook)
VALUES (1, '2019-08-11 20:58:52.646516', '2019-08-11 21:11:09.667299', null, 'Deploji frontend', 'deploji-front',
        1, 1, 'deploji/deploji', 'deploy.yml');
