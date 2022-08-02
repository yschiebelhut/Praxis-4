SELECT DISTINCT
	REGEXP_REPLACE(tasks.task_business_key, '^[a-z0-9]*-') AS c_caseid
FROM tasks
WHERE
	tasks.task_business_key LIKE '6b60d509-%-%'