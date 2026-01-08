package packagemanager

import (
	"reflect"
	"testing"
)

// To get this json, copy and paste the value of psScript from
// update_windows.go into a PowerShell window and run it.
var multipleAvailableWindowsUpdatesJson = `[
  {
    "Title": "Update für Windows Defender Antivirus-Antischadsoftwareplattform – KB4052623 (Version 4.18.2001.10)",
    "Description": "Dieses Paket aktualisiert die Komponenten der Windows Defender Antivirus-Antischadsoftwareplattform auf dem Computer des Nutzers.",
    "KBArticleIDs": [
      "4052623"
    ],
    "Categories": [
      "Microsoft Defender Antivirus",
      "Updates"
    ],
    "SupportUrl": "https://go.microsoft.com/fwlink/?linkid=862339",
    "RebootRequired": false,
    "EulaAccepted": true,
    "UpdateID": "c01629fc-64ea-45f3-b7cb-cabc7d566933",
    "RevisionNumber": 200
  },
  {
    "Title": "Update für Windows Security platform - KB5007651 (Version 10.0.29429.1000)",
    "Description": "Dieses Paket aktualisiert Windows Security platform Komponenten auf dem Benutzercomputer.",
    "KBArticleIDs": [
      "5007651"
    ],
    "Categories": [
      "Definition Updates",
      "Windows Security platform"
    ],
    "SupportUrl": "https://go.microsoft.com/fwlink/?LinkId=52661",
    "RebootRequired": false,
    "EulaAccepted": true,
    "UpdateID": "a32ca1d0-ddd4-486b-b708-d941db4f1061",
    "RevisionNumber": 205
  },
  {
    "Title": "Windows-Tool zum Entfernen bösartiger Software x64 - v5.138 (KB890830)",
    "Description": "Dieses Tool wird nach dem Herunterladen einmal ausgeführt, um den Computer auf Infektionen durch bestimmte, besonders schädliche Software (einschließlich Blaster, Sasser und Mydoom) zu überprüfen. Das Tool unterstützt Sie zudem beim Entfernen von entdeckten Infektionen. Wenn eine Infektion gefunden wurde, wird beim nächsten Starten des Computers ein Statusbericht angezeigt. Jeden Monat ist eine neue Version des Tools verfügbar. Wenn Sie das Tool manuell ausführen möchten, können Sie es im Microsoft Download Center herunterladen oder eine Onlineversion von microsoft.com ausführen. Dieses Tool kann kein Antivirenprodukt ersetzen. Sie sollten daher ein Antivirenprodukt verwenden, um zum Schutz Ihres Computers beizutragen.",
    "KBArticleIDs": [
      "890830"
    ],
    "Categories": [
      "EU Browser Choice Update-For Europe Only",
      "Update Rollups",
      "Windows 10",
      "Windows 10 LTSB",
      "Windows 10, version 1903 and later",
      "Windows 11"
    ],
    "SupportUrl": "http://support.microsoft.com",
    "RebootRequired": false,
    "EulaAccepted": true,
    "UpdateID": "9ba69000-b60b-4aae-96bb-a5506eb6ccfe",
    "RevisionNumber": 200
  },
  {
    "Title": "Update für Microsoft Defender Antivirus Antischadsoftwareplattform – KB4052623 (Version 4.18.25110.6) – Aktueller Kanal (Allgemein)",
    "Description": "Dieses Paket aktualisiert Microsoft Defender Antivirus Komponenten der Antischadsoftwareplattform auf dem Benutzercomputer.",
    "KBArticleIDs": [
      "4052623"
    ],
    "Categories": [
      "Definition Updates",
      "Microsoft Defender Antivirus"
    ],
    "SupportUrl": "https://go.microsoft.com/fwlink/?linkid=862339",
    "RebootRequired": false,
    "EulaAccepted": true,
    "UpdateID": "33ac787c-4192-4508-a736-21fd0fd16a18",
    "RevisionNumber": 200
  },
  {
    "Title": "Security Intelligence-Update für Microsoft Defender Antivirus – KB2267602 (Version 1.443.563.0) – Aktueller Kanal (Allgemein)",
    "Description": "Installieren Sie dieses Update, um die Dateien zu überarbeiten, die zum Erkennen von Viren, Spyware und anderer potenziell unerwünschter Software verwendet werden. Nachdem Sie dieses Element installiert haben, kann es nicht mehr entfernt werden.",
    "KBArticleIDs": [
      "2267602"
    ],
    "Categories": [
      "Definition Updates",
      "Microsoft Defender Antivirus"
    ],
    "SupportUrl": "https://go.microsoft.com/fwlink/?LinkId=52661",
    "RebootRequired": false,
    "EulaAccepted": true,
    "UpdateID": "4ba1ad8c-4958-467c-8f3e-7943c637fe8b",
    "RevisionNumber": 200
  },
  {
    "Title": "2025-10 Kumulatives Update für .NET Framework 3.5 und 4.8.1 für Windows 11, version 25H2 für x64 (KB5066128)",
    "Description": "In einem Microsoft-Softwareprodukt wurde ein Sicherheitsproblem festgestellt, das Auswirkungen auf Ihr System haben könnte. Durch die Installation dieses Updates von Microsoft können Sie zum Schutz Ihres Systems beitragen. Eine vollständige Liste der Problembehebungen in diesem Update finden Sie in dem entsprechenden Microsoft Knowledge Base-Artikel. Nach der Installation dieses Updates müssen Sie das System gegebenenfalls neu starten.",
    "KBArticleIDs": [
      "5066128"
    ],
    "Categories": [
      "Security Updates"
    ],
    "SupportUrl": "https://support.microsoft.com/help/5066128",
    "RebootRequired": false,
    "EulaAccepted": true,
    "UpdateID": "4f7e95d4-b8f0-404c-b682-ea2eabff590f",
    "RevisionNumber": 1
  },
  {
    "Title": "2025-12 Sicherheitsupdate (KB5072033) (26200.7462)",
    "Description": "Installieren Sie dieses Update, um Probleme in Windows zu beheben. Eine vollständige Liste der im Update enthaltenen Probleme finden Sie im zugehörigen Microsoft Knowledge Base-Artikel. Nach der Installation dieses Elements müssen Sie den Computer möglicherweise neu starten.",
    "KBArticleIDs": [
      "5072033"
    ],
    "Categories": [
      "Security Updates"
    ],
    "SupportUrl": "https://support.microsoft.com/help/5072033",
    "RebootRequired": false,
    "EulaAccepted": true,
    "UpdateID": "9769e8ee-5e2e-4f54-92ee-d616ce08ca7f",
    "RevisionNumber": 1
  }
]`

// Just Windows things -.-
var singelAvailableWindowsUpdatesJson = `{
  "Title": "JBL Driver Update (2.0.0.6)",
  "Description": "JBL HIDClass  driver update released in  December 2025",
  "KBArticleIDs": [],
  "Categories": [
    "Drivers"
  ],
  "SupportUrl": "http://support.microsoft.com/select/?target=hub",
  "RebootRequired": false,
  "EulaAccepted": true,
  "UpdateID": "fd9a5356-bf43-44a2-9561-064d84fde280",
  "RevisionNumber": 1
}`

func TestParsePowerShellUpdateSessionOutput_MultipleUpdates(t *testing.T) {
	manager := WindowsUpdatesManager{}
	updates, err := manager.parsePowerShellUpdateSessionOutput(multipleAvailableWindowsUpdatesJson)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(updates) != 7 {
		t.Errorf("expected 7 updates, got %d", len(updates))
	}
	// Check first update fields
	if updates[0].Title != "Update für Windows Defender Antivirus-Antischadsoftwareplattform – KB4052623 (Version 4.18.2001.10)" {
		t.Errorf("unexpected Title: %s", updates[0].Title)
	}
	if !reflect.DeepEqual(updates[0].KBArticleIDs, []string{"4052623"}) {
		t.Errorf("unexpected KBArticleIDs: %v", updates[0].KBArticleIDs)
	}
	if updates[0].IsInstalled != false {
		t.Errorf("expected IsInstalled to be false")
	}
	if updates[0].IsSecurityUpdate != true {
		t.Errorf("expected IsSecurityUpdate to be true")
	}
}

func TestParsePowerShellUpdateSessionOutput_SingleUpdate(t *testing.T) {
	manager := WindowsUpdatesManager{}
	updates, err := manager.parsePowerShellUpdateSessionOutput(singelAvailableWindowsUpdatesJson)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(updates) != 1 {
		t.Errorf("expected 1 update, got %d", len(updates))
	}
	if updates[0].Title != "JBL Driver Update (2.0.0.6)" {
		t.Errorf("unexpected Title: %s", updates[0].Title)
	}
	if len(updates[0].KBArticleIDs) != 0 {
		t.Errorf("expected KBArticleIDs to be empty")
	}
	if updates[0].IsInstalled != false {
		t.Errorf("expected IsInstalled to be false")
	}
}

func TestParsePowerShellUpdateSessionOutput_EmptyOutput(t *testing.T) {
	manager := WindowsUpdatesManager{}
	updates, err := manager.parsePowerShellUpdateSessionOutput("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(updates) != 0 {
		t.Errorf("expected 0 updates, got %d", len(updates))
	}
}

func TestParsePowerShellUpdateSessionOutput_InvalidJson(t *testing.T) {
	manager := WindowsUpdatesManager{}
	_, err := manager.parsePowerShellUpdateSessionOutput("{invalid json}")
	if err == nil {
		t.Errorf("expected error for invalid json, got nil")
	}
}
