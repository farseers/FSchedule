SELECT client_name,
       execute_status,
       COUNT(*) AS count
from fschedule.fschedule_task
WHERE create_at >= (NOW() - INTERVAL 30 MINUTE) and client_name !=''
GROUP BY client_name,execute_status;