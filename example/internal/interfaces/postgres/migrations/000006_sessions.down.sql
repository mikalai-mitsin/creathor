DROP TABLE public.sessions;

DELETE
FROM public.permissions
WHERE id IN (
    'session_list',
    'session_detail',
    'session_create',
    'session_update',
    'session_delete'
);
