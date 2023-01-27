DROP TABLE public.days;

DELETE
FROM public.permissions
WHERE id IN (
    'day_list',
    'day_detail',
    'day_create',
    'day_update',
    'day_delete'
);
