CREATE VIEW "USER".HS2_YM_DEMO_split (
	"c_caseid",
	"CASE_GROUP",
	"c_time",
	"c_eventname",
	"SENDER",
	"SENDER_TYPE",
	"RECIEVER",
	"RECEIVER_TYPE",
	"PRICE",
	"PRODUCT"
) AS

SELECT
	SUBSTR_AFTER(tasks.task_business_key, '-') AS c_caseid,
	SUBSTR_BEFORE(SUBSTR_AFTER(tasks.task_business_key, '-'), '-' ) AS case_group,
	tasks.tx_ts AS c_time,
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
	recievers.member_name AS reciever,
	recievers.member_role AS reciever_type,
	CAST(0 AS INTEGER) AS price,
	'' AS product

FROM "USER".TASKS AS tasks
LEFT JOIN "USER".MEMBERS AS senders ON tasks.publickey = senders.publickey
LEFT JOIN "USER".CONFIRMATIONREQUESTS AS confirmationrequest ON tasks.id = confirmationrequest.confirmationrequest_task_id
LEFT JOIN "USER".MEMBERS AS recievers ON recievers.publickey = confirmationrequest.publickey

WHERE (
	(
		tasks.task_tasktemplate_id LIKE 'S15K62OMJIMOFF8O' OR
		tasks.task_tasktemplate_id LIKE 'P2OKGZBL1GJ3TU9J' OR
		tasks.task_tasktemplate_id LIKE 'MCWWVIX5Z2D0WW1D' OR
		tasks.task_tasktemplate_id LIKE 'LDH299VRNU4A3CIT'
	) AND
	tasks.task_business_key LIKE '6b60d509-%-%'
)



UNION ALL



SELECT
	SUBSTR_AFTER(tasks.task_business_key, '-') AS c_caseid,
	SUBSTR_BEFORE(SUBSTR_AFTER(tasks.task_business_key, '-'), '-' ) AS case_group,
	confirmations.tx_ts AS c_time,
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

FROM "USER".CONFIRMATIONS AS confirmations
LEFT JOIN "USER".TASKS AS tasks ON tasks.id = confirmations.confirmation_task_id
LEFT JOIN "USER".MEMBERS AS senders ON senders.publickey = confirmations.publickey
LEFT JOIN "USER".CONTEXT AS c1 ON tasks.ID = c1.ID AND c1.KEY = 'Price'
LEFT JOIN "USER".CONTEXT AS c2 ON tasks.ID = c2.ID AND c2.KEY = 'Product'
LEFT JOIN "USER".MEMBERS AS recievers ON recievers.publickey = tasks.publickey

WHERE (
	(
		tasks.task_tasktemplate_id LIKE 'S15K62OMJIMOFF8O' OR
		tasks.task_tasktemplate_id LIKE 'P2OKGZBL1GJ3TU9J' OR
		tasks.task_tasktemplate_id LIKE 'MCWWVIX5Z2D0WW1D' OR
		tasks.task_tasktemplate_id LIKE 'LDH299VRNU4A3CIT'
	) AND
	tasks.task_business_key LIKE '6b60d509-%-%'
)



UNION ALL



SELECT
	SUBSTR_AFTER(tasks.task_business_key, '-') AS c_caseid,
	SUBSTR_BEFORE(SUBSTR_AFTER(tasks.task_business_key, '-'), '-' ) AS case_group,
	ADD_SECONDS(confirmations.tx_ts, 1) AS c_time,
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

FROM "USER".CONFIRMATIONS AS confirmations
LEFT JOIN "USER".TASKS AS tasks ON tasks.id = confirmations.confirmation_task_id
LEFT JOIN "USER".MEMBERS AS senders ON senders.publickey = confirmations.publickey
LEFT JOIN "USER".MEMBERS AS recievers ON recievers.publickey = tasks.publickey

WHERE (
	(
		tasks.task_tasktemplate_id LIKE 'S15K62OMJIMOFF8O' OR
		tasks.task_tasktemplate_id LIKE 'P2OKGZBL1GJ3TU9J' OR
		tasks.task_tasktemplate_id LIKE 'MCWWVIX5Z2D0WW1D' OR
		tasks.task_tasktemplate_id LIKE 'LDH299VRNU4A3CIT'
	) AND
	tasks.task_business_key LIKE '6b60d509-%-%'
)

ORDER BY
	c_caseid ASC,
	c_time ASC,
	c_eventname DESC;