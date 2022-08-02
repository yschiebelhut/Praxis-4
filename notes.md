fix sql ac



## Schwierigkeiten:
- Signavio in aktiver Entwicklung, Features fallen unerwartet aus
- Signavio ist nicht für diese Art von Datenauswertung vorgesehen
- Gesamtes Projekt ist ein Prototyp, dadurch sind Dokumentationen lückenhaft und einige technische Voraussetzungen sind nicht erfüllt
- Signavio Dokumentation bietet Optimierungspotential
- Mangelhafte Fehlermeldungen
- Signavio SAP HANA Connector ist in Betaphase, dadurch sind einige Konfigurationsmöglichkeiten nicht gegeben und es müssen Umwege dafür gefunden werden
- mangelhafte Implementierung in Signavio, sorgt für zusätzlich notwendige Schritte



## Ledger:
\begin{itemize}
    \item Testprojekt, Prototyp, nicht kommerziell verfügbar
    \item kombiniert Vorteile von Blockchain Technologie und der SAP HANA Cloud
    \item dadurch werden viele Nachteile der Blockchain Technologie annulliert
    \item reduziert Deployment Aufwand und Kosten im Vergleich zu klassischen Blockchain Lösungen \cite{ledger.jam}
    \item in vollständiges SAP Ökosystem integriert
    
    \item technische Zugriffsrechte
    \item bestimmte Spalten, die für eine Verkettung der Transaktionen genutzt wird
    \item diese Verkettung ist von Blockchain Technologie inspiriert
    \item \enquote{append-only}
    
    \item Spalten in zwei \enquote{Gruppen}
    \item ID, Document, Timestamp, Public Key, Hash, Signature werden auf Client Seite erstellt und als Teil einer \enquote{create new transaction} Request an den Server geschickt
    \item Werte der verbleibenden Spalten werden vom Service selbst serverseitig generiert
    \item Service Signature bestätigt, dass der Service die Integrität der eingehenden Transaktion überprüft hat
    

    \item kryptografische Signatur des Clients verhindert Manipulation der Transaktionsdaten
    \item nur der Besitzer des Private Keys könnte die Transaktion manipulieren, allerdings kann er ohne den Key des Shared Ledger Services nicht die Serivce Signatur erneuern
    \item dadurch kann weder der Ersteller noch der Ledger Bereitsteller die Transaktionsdaten ändern
    \item Wer? (Public Key), Was? (Document), Wann? (Timestamp) sind unveränderbar und unwiderruflich im Ledger verzeichnet
    \item Verifikaiton ist durch alle Teilnehmer mit Leseberechtigung auf den Ledger möglich
    
    \item der Service Provider kann zwar nicht die Daten beeinflussen, jedoch könnte er theoretisch den Service Timestamp manipulieren und damit die Reihenfolge der Transaktionen (oder auch die Transaction ID)
    \item dies würde allerdings direkt den Service Hash verändern
    \item somit reicht es aus, den Service Hash abzusichern
    
    \item deshalb -- Service Hashes der Transaktionen in Blocks abspeichern
    \item Blöcke werden dadurch verkettet, dass ein Block (außer dem ersten) den Service Hash seines Vorgängerblocks enthält und vom Service signiert wird
    \item wird eine Transaktion gefälscht, wird automatisch der Block Hash ungültig
    \item wird der Block Hash neu bestimmt, so ist die Referenz des nachfolgenden Blocks auf den aktuellen ungültig
    \item deshalb müssten, bei Manipulation einer Transaktion, der aktuelle Block und alle nachfolgenden Blöcke neu gehasht werden
    \item in Folge dessen würde es den Teilnehmern reichen, sich den Hash des aktuellsten Blocks zu merken, um zu einem späteren Zeitpunkt Modifikationen beweisen zu können
    
    \item Zeitstempel einer Transaktion müssen zwischen dem des vorangehenden Blocks und dem des aktuellen Blocks liegen
    \item deshalb wird der Schreibzugriff auf die Tabelle während der Erstellung eines Blocks gesperrt
    
    \item für jeden Ledger werden Blöcke alle 5 Sekunden generiert, falls neue Transaktionen erstellt wurden
    \item Blöcke sind auch Transaktionen
\end{itemize}




## CCWF
\begin{itemize}
    \item SAP Workflow mit unternehmensübergreifenden Features erweitern
    \item Nutzung von Blockchain Technologien
    
    \item Zusammenarbeit mit Partnern muss verbessert werden, um Produkte und dienstleistungen günstig und reaktionsschnell entwickeln zu können
    \item Geschäftsprozesse über Unternehmensgrenzen hinaus automatisiert in vollständige Wertschöpfungsketten organisieren
    \item Möglichkeit, anhand von cross-company Tasks unternehmensspezifische Workflows zur Verarbeitung dieser zu triggern
    \item no-code Werkzeuge zur Verwaltung unternehmensübergreifender Szenarien
    \item alle relevanten Schritte der Zusammenarbeit auf Shared Ledger gespeichert
    \item out-of-the-box Erweiterung für SAP Workflow Management
    \item manuelle Integration in Tools Dritter möglich
    
    \item unterstützt digitale Prozess Autonominisierung über diverse Services
    \item Workflows digitalisieren
    \item Prozessentscheidungen verwalten (Angebotsannahme, etc)
    \item Überblick über alle Teile des Prozesses bekommen
    \item Entwickeln von Erweiterungen für Kunden und Partner möglich
    
    \item 4 Kernkompetenzen
    \item i direkte Integration in bestehende SAP und Drittanbieter Applikationen möglich, damit Kunden diese nicht ersetzen müssen
    \item ii out-of-the-box Einbindung in SAP Workflow Management
    \item iii fortschrittliche firmenübergreifende Automation -- jede beteiligte Partei kann einen Workflow definieren, welcher bei eingehender Task eines externen Partners automatisch angestoßen wird
    \item iv konsistenter und unveränderbarer Audit-Log gewährt für alle Teilnehmer einen Überblick über den aktuellen Prozessstand; reduziert Risiko Betrug, Missverständnisse oder Uneinigkeiten, welche sonst nur mittels aufwändiger und kostspieliger Kontrolle aufgeklärt werden könnten
    
    \item Firmen können digitale Vertrauensbasis nutzen, um eigene Risiken zugunsten von gesteigerter Effizienz und Nachhaltigkeit zu minimieren
    \item Manuelle Prozesse zu digitalisieren senkt unter anderem Betriebskosten, reduziert das Risiko für Verzögerungen des Projekts oder verspätete Zahlungen und verbessert die \ce{CO2}-Bilanz
    
    \item an einer Integration ins Gesundheitswesen, den öffentlichen- und Banksektor, Lieferketten und Transportnetz wird gearbeitet
    
    \cite{Westphal2021}
\end{itemize}

Ein weiterer Pluspunkt ist der gewählte No-Code-Ansatz.
Dieser bedeutet, dass das Erstellen von Prozessen, Aufgaben und Ähnlichem einfach grafisch und ohne Programmierkenntnisse geschehen kann, was die Nutzerfreundlichkeit und Supportkosten enorm verringern kann.
