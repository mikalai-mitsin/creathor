DROP TABLE public.plans;

DELETE
FROM public.permissions
WHERE id IN (
    'plan_list',
    'plan_detail',
    'plan_create',
    'plan_update',
    'plan_delete'
);
