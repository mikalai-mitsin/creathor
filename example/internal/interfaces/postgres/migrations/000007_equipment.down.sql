DROP TABLE public.equipment;

DELETE
FROM public.permissions
WHERE id IN (
    'equipment_list',
    'equipment_detail',
    'equipment_create',
    'equipment_update',
    'equipment_delete'
);
