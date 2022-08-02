<!-- LTeX: enabled=false -->

# Entwicklung eines Datenintegrationsprozesses von der SAP Blockchain in Signavio Process Intelligence

## Gliederung

- Einleitung
- Theoretische Grundlagen
	- Process Mining und Signavio
		- BPMN
		- Datenquellen und ETL
		- Investigations
	- Blockchain
		- HANA Blockchain
			- (grober Aufbau des Ledgers, Felder und Blockchain-Bezug)
	- ccWF
		- (HS2)
		- Views
	- SQL & HANA Views
- Konzeptionierung
	- Datenvorbereitung
		- Großer View
		- ETL
		- Mischung &rArr; Event-View + ETL
	- Datenbereitstellung
		- API
		- HANA Connector
		- File-Import
		- (Data Intelligence)
- Implementierung der Datenintegration
- Auswertung und Future Work

## Notizen

- Analyse und Vergleich verschiedener Varianten
	- Datenvorbereitung
	- Großer Event-View auf der HANA
	- Direktes Einbinden der ccWF Views (in die Signavio ETL)
	- Signavio Data Ingestion API
- Auswahl und Implementierung einer Variante

- Auf Footprint der Lösungen eingehen (in kritischer Reflexion)
- e.g. brauche ich noch ne Data Intelligence Lizenz, etc?

- Doku schreiben im Wiki oder auf GitHub


---

- Process Mining: RWTH
	- Prozesse innerhalb eines Unternehmens

+ geteilte Eventlogs

- warum Blockchain für geteilte Eventlogs
- Auditierbar
- Firmen müssen sich nicht vertrauen
- immutable

- Motivation warum Blockchain

- Datenstrukur

Ausblick:
- KPIs über geteilten Workflow