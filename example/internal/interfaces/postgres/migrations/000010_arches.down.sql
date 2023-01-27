DROP TABLE public.arches;

DELETE
FROM public.permissions
WHERE id IN (
    'arch_list',
    'arch_detail',
    'arch_create',
    'arch_update',
    'arch_delete'
);
