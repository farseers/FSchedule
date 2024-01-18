select
    count(1) as 'Count'
from fschedule.fschedule_task_log as log
    inner join fschedule.fschedule_task as task on log.task_id=task.id
where 1=1 %s