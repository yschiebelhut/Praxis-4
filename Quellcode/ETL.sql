SELECT
	REGEXP_REPLACE(tasks.task_business_key, '^[a-z0-9]*-') AS c_caseid,
  REGEXP_REPLACE(REGEXP_REPLACE(tasks.task_business_key, '^[a-z0-9]*-'), '-.*$') AS case_group,
	FROM_UNIXTIME(tasks.tx_ts / 1000) AS c_time,
	(
    CASE
     	WHEN tasks.task_tasktemplate_id LIKE 'S15K62OMJIMOFF8O'
				THEN 'Request for Quotation Sent'
			WHEN tasks.task_tasktemplate_id LIKE 'P2OKGZBL1GJ3TU9J'
				THEN 'Purchase Order Raised'
			WHEN tasks.task_tasktemplate_id LIKE 'MCWWVIX5Z2D0WW1D'
				THEN 'Asset Handover Confirmation Requested'
			WHEN tasks.task_tasktemplate_id LIKE 'LDH299VRNU4A3CIT'
				THEN 'Final Inspection Requested'
		END
  ) AS c_eventname,
  senders.member_name AS sender,
  senders.member_role AS sender_type,
  receiver.member_name AS reciever,
  receiver.member_role AS reciever_type,
  CAST(0 AS INTEGER) AS price,
  '' as product
    
FROM tasks 
LEFT JOIN tasktemplates ON tasks.task_tasktemplate_id = tasktemplates.id
LEFT JOIN members AS senders ON tasks.publickey = senders.publickey 
LEFT JOIN confirmationrequests ON tasks.id = confirmationrequests.confirmationrequest_task_id 
LEFT JOIN members AS receiver ON receiver.publickey = confirmationrequests.publickey

WHERE (
	(
    tasks.task_tasktemplate_id LIKE 'S15K62OMJIMOFF8O' OR
		tasks.task_tasktemplate_id LIKE 'P2OKGZBL1GJ3TU9J' OR
		tasks.task_tasktemplate_id LIKE 'MCWWVIX5Z2D0WW1D' OR
		tasks.task_tasktemplate_id LIKE 'LDH299VRNU4A3CIT'
  ) AND	
  tasks.TASK_BUSINESS_KEY LIKE '6b60d508-%-%'
)



UNION ALL



SELECT
	REGEXP_REPLACE(tasks.task_business_key, '^[a-z0-9]*-') AS c_caseid,
  REGEXP_REPLACE(REGEXP_REPLACE(tasks.task_business_key, '^[a-z0-9]*-'), '-.*$') AS CASE_GROUP,
	FROM_UNIXTIME(confirmations.tx_ts / 1000) AS c_time,
	(
    CASE
      WHEN tasks.task_tasktemplate_id LIKE 'S15K62OMJIMOFF8O'
				THEN 'Request for Quotation Received'
			WHEN tasks.task_tasktemplate_id LIKE 'P2OKGZBL1GJ3TU9J'
				THEN 'Purchase Order Received'
			WHEN tasks.task_tasktemplate_id LIKE 'MCWWVIX5Z2D0WW1D'
				THEN 'Asset Handover Confirmation Received'
			WHEN tasks.task_tasktemplate_id LIKE 'LDH299VRNU4A3CIT'
				THEN 'Final Inspection Received'
		END
  ) AS c_eventname,
  senders.member_name AS sender,
  senders.member_role AS sender_type,
  recievers.member_name AS reciever,
  recievers.member_role AS reciever_type,
  CAST(c1.value AS INTEGER) AS price,
  c2.value AS product
    
FROM confirmations 
LEFT JOIN tasks ON tasks.id = confirmations.confirmation_task_id 
LEFT JOIN tasktemplates	ON tasktemplates.id = tasks.task_tasktemplate_id
LEFT JOIN members AS senders ON senders.publickey = confirmations.publickey 
LEFT JOIN members AS recievers ON recievers.publickey = tasks.publickey
LEFT JOIN context AS c1 ON tasks.ID = c1.ID AND c1.KEY = 'Price' 
LEFT JOIN context AS c2 ON tasks.ID = c2.ID AND c2.KEY = 'Product' 

WHERE (
	(
    tasks.task_tasktemplate_id LIKE 'S15K62OMJIMOFF8O' OR
		tasks.task_tasktemplate_id LIKE 'P2OKGZBL1GJ3TU9J' OR
		tasks.task_tasktemplate_id LIKE 'MCWWVIX5Z2D0WW1D' OR
		tasks.task_tasktemplate_id LIKE 'LDH299VRNU4A3CIT'
  ) AND	
  tasks.TASK_BUSINESS_KEY LIKE '6b60d508-%-%'
)



UNION ALL



SELECT
	REGEXP_REPLACE(tasks.task_business_key, '^[a-z0-9]*-') AS c_caseid,
  REGEXP_REPLACE(REGEXP_REPLACE(tasks.task_business_key, '^[a-z0-9]*-'), '-.*$') AS CASE_GROUP,
	FROM_UNIXTIME((confirmations.tx_ts + 1000) / 1000) AS c_time,
	(
    CASE
      WHEN tasks.task_tasktemplate_id LIKE 'S15K62OMJIMOFF8O'
				THEN (
          CASE
            WHEN confirmations.confirmation_decisions IN('["approve"]', '["Approve"]')
							THEN 'Request for Quotation Approved'
						WHEN confirmations.confirmation_decisions IN('["reject"]', '["Reject"]')
							THEN 'Request for Quotation Rejected'
					END
        )
			WHEN tasks.task_tasktemplate_id LIKE 'P2OKGZBL1GJ3TU9J'
				THEN (
          CASE
            WHEN confirmations.confirmation_decisions IN('["approve"]', '["Approve"]')
							THEN 'Purchase Order Approved'
						WHEN confirmations.confirmation_decisions IN('["reject"]', '["Reject"]')
							THEN 'Purchase Order Rejected'
					END
        )
			WHEN tasks.task_tasktemplate_id LIKE 'MCWWVIX5Z2D0WW1D'
				THEN (
          CASE
            WHEN confirmations.confirmation_decisions IN('["approve"]', '["Approve"]')
							THEN 'Asset Handover Approved'
						WHEN confirmations.confirmation_decisions IN('["reject"]', '["Reject"]')
							THEN 'Asset Handover Rejected'
					END
        )
			WHEN tasks.task_tasktemplate_id LIKE 'LDH299VRNU4A3CIT'
				THEN (
          CASE
            WHEN confirmations.confirmation_decisions IN('["approve"]', '["Approve"]')
							THEN 'Final Inspection Passed'
						WHEN confirmations.confirmation_decisions IN('["reject"]', '["Reject"]')
							THEN 'Final Inspection Failed'
					END
        )
		END
  ) AS c_eventname,
  senders.member_name AS sender,
  senders.member_role AS sender_type,
  recievers.member_name AS reciever,
  recievers.member_role AS reciever_type,
  CAST(0 AS INTEGER) AS price,
  '' AS product
    
FROM confirmations
LEFT JOIN tasks ON tasks.id = confirmations.confirmation_task_id 
LEFT JOIN tasktemplates ON tasktemplates.id = tasks.task_tasktemplate_id
LEFT JOIN members AS senders ON senders.publickey = confirmations.publickey 
LEFT JOIN members AS recievers ON recievers.publickey = tasks.publickey

WHERE (
	(
    tasks.task_tasktemplate_id LIKE 'S15K62OMJIMOFF8O' OR
		tasks.task_tasktemplate_id LIKE 'P2OKGZBL1GJ3TU9J' OR
		tasks.task_tasktemplate_id LIKE 'MCWWVIX5Z2D0WW1D' OR
		tasks.task_tasktemplate_id LIKE 'LDH299VRNU4A3CIT'
  ) AND	
  tasks.TASK_BUSINESS_KEY LIKE '6b60d508-%-%'
)

ORDER BY 
	c_caseid ASC, 
	c_time ASC, 
	c_eventname DESC