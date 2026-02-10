package packagemanager

import (
	"reflect"
	"testing"
)

var softwareUpdateOutputSample = `Software Update Tool

Finding available software
Software Update found the following new or updated software:
* Label: Command Line Tools for Xcode 26.2-26.2
	Title: Command Line Tools for Xcode 26.2, Version: 26.2, Size: 858715KiB, Recommended: YES,
* Label: macOS Tahoe 26.2-25C56
	Title: macOS Tahoe 26.2, Version: 26.2, Size: 3693075KiB, Recommended: YES, Action: restart,

`

// Currently unused, but kept for future tests
// pkgutil --pkgs
var pkgutilOutputSample = `pkgutil --pkgs                                                                                                                                                               ✔  3s   07:58:16  
com.apple.pkg.XProtectPlistConfigData_10_15.16U4408
com.apple.pkg.MAContent10_AssetPack_0637_AppleLoopsDrummerKyle
com.apple.pkg.MAContent10_AssetPack_0593_DrummerSoCalGBLogic
com.apple.pkg.CLTools_Executables
com.apple.files.data-template
com.apple.pkg.CLTools_SDK_macOS_LMOS
com.apple.pkg.GatekeeperCompatibilityData.16U1906
com.apple.pkg.MobileAssets
com.apple.pkg.MAContent10_AssetPack_0317_AppleLoopsModernRnB1
com.apple.pkg.MAContent10_AssetPack_0537_DrummerShaker
com.apple.pkg.MAContent10_AssetPack_0482_EXS_OrchWoodwindAltoSax
com.apple.pkg.XProtectPayloads_10_15.16U4404
com.apple.pkg.CLTools_SDK_macOS_NMOS
com.apple.pkg.MAContent10_AssetPack_0048_AlchemyPadsDigitalHolyGhost
com.apple.pkg.MAContent10_AssetPack_0539_DrummerTambourine
com.apple.pkg.MAContent10_AssetPack_0323_AppleLoopsVintageBreaks
com.apple.pkg.MAContent10_AssetPack_0487_EXS_OrchWoodwindFluteSolo
com.apple.pkg.MAContent10_AssetPack_0557_IRsSharedAUX
com.apple.pkg.MAContent10_AssetPack_0484_EXS_OrchWoodwindClarinetSolo
com.apple.pkg.MRTConfigData_10_15.16U4211
com.apple.pkg.MAContent10_AssetPack_0310_UB_DrumMachineDesignerGB
com.apple.pkg.MobileDevice
com.apple.pkg.MAContent10_AssetPack_0316_AppleLoopsDubstep1
com.apple.pkg.MAContent10_AssetPack_0491_EXS_OrchBrass
com.apple.pkg.MAContent10_AssetPack_0560_LTPBasicPiano1
com.apple.pkg.MAContent10_AssetPack_0315_AppleLoopsElectroHouse1
com.apple.pkg.MAContent10_AssetPack_0322_AppleLoopsDiscoFunk1
com.apple.pkg.MAContent10_AssetPack_0597_LTPChordTrainer
com.apple.pkg.MAContent10_AssetPack_0540_PlugInSettingsGB
com.apple.pkg.MAContent10_AssetPack_0375_EXS_GuitarsVintageStrat
com.apple.pkg.MAContent10_AssetPack_0806_PlugInSettingsGBLogic
com.apple.pkg.MAContent10_AssetPack_0615_GBLogicAlchemyEssentials
com.apple.pkg.MAContent10_AssetPack_0509_EXS_StringsEnsemble
com.apple.pkg.MAContent10_AssetPack_0324_AppleLoopsBluesGarage
com.apple.pkg.MAContent10_AssetPack_0314_AppleLoopsHipHop1
com.apple.pkg.RosettaUpdateAuto
com.apple.pkg.MAContent10_AssetPack_0325_AppleLoopsGarageBand1
com.apple.pkg.MAContent10_AssetPack_0371_EXS_GuitarsAcoustic
com.apple.pkg.MAContent10_AssetPack_0312_UB_UltrabeatKitsGBLogic
com.apple.pkg.MAContent10_AssetPack_0646_AppleLoopsDrummerElectronic
com.apple.pkg.MAContent10_AssetPack_0598_LTPBasicGuitar1
com.apple.pkg.MAContent10_AssetPack_0354_EXS_PianoSteinway
com.apple.pkg.CoreTypes.1900A26
com.apple.pkg.MAContent10_AssetPack_0321_AppleLoopsIndieDisco
com.apple.pkg.MAContent10_AssetPack_0538_DrummerSticks
com.apple.pkg.MAContent10_AssetPack_0358_EXS_BassElectricFingerStyle
com.apple.pkg.MobileDeviceDevelopment
com.apple.pkg.CLTools_SwiftBackDeploy
com.apple.pkg.MAContent10_AssetPack_0554_AppleLoopsDiscoFunk2
com.apple.pkg.CLTools_macOS_SDK
com.apple.pkg.XcodeSystemResources
com.apple.pkg.MAContent10_AssetPack_0536_DrummerClapsCowbell
com.apple.pkg.XProtectPlistConfigData_10_15.16U4407
com.apple.pkg.MAContent10_AssetPack_0357_EXS_BassAcousticUprightJazz
com.apple.pkg.MAContent10_AssetPack_0320_AppleLoopsChillwave1
com.apple.pkg.XProtectPlistConfigData_10_15.16U4406
org.openvpn.client.pkg
com.fortinet.forticlient.vpnservice
com.apple.cdm.pkg.Pages_MASReceipt
com.microsoft.edgemac
com.apple.pkg.Pages14
com.microsoft.package.Microsoft_AU_Bootstrapper.app
com.adobe.acrobat.AcroRdrSCADCUpd2500120529_MUI
com.citrix.ICAClient
org.golang.go
com.microsoft.dlp.ux
com.tinyspeck.slackmacgap
com.adobe.acrobat.AcroRdrSCADCUpd2500120476_MUI
com.microsoft.wdav.shim
com.apple.cdm.pkg.Keynote_MASReceipt
com.it-novum.openitcockpit.agent
org.virtualbox.pkg.virtualbox
com.adobe.acrobat.DC.sca.config.application.1.pkg.MUI
com.microsoft.dlp.daemon
org.virtualbox.pkg.virtualboxcli
com.apple.pkg.iMovie_AppStore
com.microsoft.package.Microsoft_Excel.app
com.citrix.devicetrust.client.ica
com.fortinet.forticlient.preinstall
com.citrix.devicetrust.client
com.adobe.acrobat.DC.scamini.app.pkg.MUI
com.adobe.acrobat.AcroRdrSCADCUpd2500120577_MUI
com.adobe.acrobat.AcroRdrSCADCUpd2500120841_MUI
com.adobe.acrobat.AcroRdrSCADCUpd2500120756_MUI
com.microsoft.CompanyPortalMac
com.adobe.acrobat.AcroRdrSCADCUpd2500120566_MUI
com.microsoft.package.Microsoft_Word.app
com.microsoft.OneDrive
com.apple.pkg.GarageBand_AppStore
com.adobe.acrobat.DC.scamini.appsupport.pkg.MUI
com.microsoft.package.Frameworks
com.apple.cdm.pkg.Numbers_MASReceipt
org.openvpn.helper_framework.pkg
com.microsoft.teams
com.microsoft.package.Microsoft_OneNote.app
com.microsoft.wdav
com.microsoft.package.Microsoft_Outlook.app
com.microsoft.dlp.agent
com.adobe.acrobat.AcroRdrSCADCUpd2500120693_MUI
com.apple.cdm.pkg.GarageBand_MASReceipt
com.microsoft.package.Proofing_Tools
com.microsoft.package.DFonts
com.apple.pkg.Xcode
com.citrix.enterprisebrowserinstaller
com.fortinet.forticlient.Uninstall
org.openvpn.client_framework.pkg
com.apple.cdm.pkg.iMovie_MASReceipt
com.apple.pkg.Keynote14
com.microsoft.intuneMDMAgent
com.microsoft.package.Microsoft_AutoUpdate.app
com.adobe.armdc.app.pkg
com.fortinet.forticlient.postinstall
com.fortinet.forticlient.fssoagent
com.adobe.acrobat.AcroRdrSCADCUpd2500120997_MUI
org.openvpn.client_launch.pkg
com.microsoft.MSTeamsAudioDevice
com.microsoft.pkg.licensing
com.adobe.acrobat.AcroRdrSCADCUpd2500120937_MUI
org.openvpn.helper_launch.pkg
com.microsoft.SkypeForBusiness
com.adobe.acrobat.AcroRdrSCADCUpd2500120630_MUI
org.ocsinventory-ng.agent.macosx
com.fortinet.forticlient.commservice
com.displaylink.displaylinkloginscreenext
com.fortinet.forticlient.FortiClientarm64
com.adobe.acrobat.AcroRdrSCADCUpd2500120643_MUI
com.apple.pkg.Numbers14
org.openvpn.client_uninstall.pkg
com.displaylink.displaylinkmanagerapp
com.microsoft.teams2
com.RootRiseTechnologies.Read-CHM
com.microsoft.package.Microsoft_PowerPoint.appv

`

// system_profiler SPApplicationsDataType -json
var systemProfilerOutputSample = `{
  "SPApplicationsDataType": [
    {
      "_name": "App Store",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/App Store.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "3.0"
    },
    {
      "_name": "Apps",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/Apps.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "Automator",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/Automator.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "2.10"
    },
    {
      "_name": "Bücher",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/Books.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "8.1"
    },
    {
      "_name": "Rechner",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/Calculator.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "12.0"
    },
    {
      "_name": "Kalender",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/Calendar.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "16.0"
    },
    {
      "_name": "Schach",
      "arch_kind": "arch_arm_i64",
      "info": "3.18, Copyright 2003–2024 Apple Inc.",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/Chess.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "3.18"
    },
    {
      "_name": "Uhr",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/Clock.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.1"
    },
    {
      "_name": "Kontakte",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/Contacts.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "14.0"
    },
    {
      "_name": "Passwörter",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/Passwords.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "2.1"
    },
    {
      "_name": "Telefon",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/Phone.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "Photo Booth",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/Photo Booth.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "13.1"
    },
    {
      "_name": "Fotos",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/Photos.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "11.0"
    },
    {
      "_name": "Podcasts",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/Podcasts.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.1.0"
    },
    {
      "_name": "Vorschau",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/Preview.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "11.0"
    },
    {
      "_name": "QuickTime Player",
      "arch_kind": "arch_arm_i64",
      "info": "10.5, Copyright © 2009-2025 Apple Inc. All Rights Reserved.",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/QuickTime Player.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "10.5"
    },
    {
      "_name": "Erinnerungen",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/Reminders.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "7.0"
    },
    {
      "_name": "Kurzbefehle",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/Shortcuts.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "7.0"
    },
    {
      "_name": "Siri",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/Siri.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "Notizzettel",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/Stickies.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "10.3"
    },
    {
      "_name": "Aktien",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/Stocks.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "8.1"
    },
    {
      "_name": "Systemeinstellungen",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/System Settings.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "15.0"
    },
    {
      "_name": "TV",
      "arch_kind": "arch_arm_i64",
      "info": "TV 1.6.1.44, © 2019–2025 Apple Inc. All rights reserved.",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/TV.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.6.1"
    },
    {
      "_name": "TextEdit",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/TextEdit.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.20"
    },
    {
      "_name": "Time Machine",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/Time Machine.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.3"
    },
    {
      "_name": "Tipps",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/Tips.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "26.1"
    },
    {
      "_name": "Aktivitätsanzeige",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/Utilities/Activity Monitor.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "10.14"
    },
    {
      "_name": "AirPort-Dienstprogramm",
      "arch_kind": "arch_arm_i64",
      "info": "6.3.9, Copyright 2001 -2025 Apple Inc.",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/Utilities/AirPort Utility.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "6.3.9"
    },
    {
      "_name": "Audio-MIDI-Setup",
      "arch_kind": "arch_arm_i64",
      "info": "3.7, Copyright 2002–2025 Apple Inc.",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/Utilities/Audio MIDI Setup.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "3.7"
    },
    {
      "_name": "Bluetooth-Datenaustausch",
      "arch_kind": "arch_arm_i64",
      "info": "7.0.0, Copyright © 2002-2018 Apple Inc. All rights reserved.",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/Utilities/Bluetooth File Exchange.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "9.0"
    },
    {
      "_name": "Boot Camp-Assistent",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/Utilities/Boot Camp Assistant.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "6.1.0"
    },
    {
      "_name": "ColorSync-Dienstprogramm",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/Utilities/ColorSync Utility.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "12.2.0"
    },
    {
      "_name": "Konsole",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/Utilities/Console.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.1"
    },
    {
      "_name": "Digital Color Meter",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/Utilities/Digital Color Meter.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "6.10"
    },
    {
      "_name": "Festplattendienstprogramm",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/Utilities/Disk Utility.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "22.7"
    },
    {
      "_name": "Grapher",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/Utilities/Grapher.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "2.8"
    },
    {
      "_name": "Lupe",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/Utilities/Magnifier.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "Migrationsassistent",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/Utilities/Migration Assistant.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "26.1"
    },
    {
      "_name": "Druckzentrale",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/Utilities/Print Center.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "Bildschirmfreigabe",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/Utilities/Screen Sharing.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "6.1"
    },
    {
      "_name": "Bildschirmfoto",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/Utilities/Screenshot.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "Skripteditor",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/Utilities/Script Editor.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "2.11"
    },
    {
      "_name": "Systeminformationen",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/Utilities/System Information.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "11.0"
    },
    {
      "_name": "Terminal",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/Utilities/Terminal.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "2.15"
    },
    {
      "_name": "VoiceOver-Dienstprogramm",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/Utilities/VoiceOver Utility.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "10"
    },
    {
      "_name": "Sprachmemos",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/VoiceMemos.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "3.2"
    },
    {
      "_name": "Wetter",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/Weather.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "6.0"
    },
    {
      "_name": "iPhone-Synchronisierung",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/iPhone Mirroring.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.5"
    },
    {
      "_name": "Über diesen Mac",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/Applications/About This Mac.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "Archivierungsprogramm",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/Applications/Archive Utility.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "10.15"
    },
    {
      "_name": "DVD-Player",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/Applications/DVD Player.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "6.0"
    },
    {
      "_name": "Schreibtischansicht",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/Applications/Desk View.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "2.0"
    },
    {
      "_name": "Verzeichnisdienste",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/Applications/Directory Utility.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "7.0"
    },
    {
      "_name": "Dienstprogramm für Erweiterungssteckplätze",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/Applications/Expansion Slot Utility.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "2.0"
    },
    {
      "_name": "Feedback-Assistent",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/Applications/Feedback Assistant.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "5.1"
    },
    {
      "_name": "SystemUIServer",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/SystemUIServer.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.7"
    },
    {
      "_name": "TextInputMenuAgent",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/TextInputMenuAgent.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "TextInputSwitcher",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/TextInputSwitcher.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.1"
    },
    {
      "_name": "ThermalTrap",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/ThermalTrap.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "TMHelperAgent",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/TimeMachine/TMHelperAgent.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "13"
    },
    {
      "_name": "Tipps",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/TipsSpotlightHandler.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "UIKitSystem",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/UIKitSystem.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "UniversalAccessControl",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/UniversalAccessControl.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "7.0"
    },
    {
      "_name": "Universelle Steuerung",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/UniversalControl.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "UnmountAssistantAgent",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/UnmountAssistantAgent.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "5.0"
    },
    {
      "_name": "UserNotificationCenter",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/UserNotificationCenter.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "82"
    },
    {
      "_name": "VoiceOver",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/VoiceOver.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "10"
    },
    {
      "_name": "Hintergrundbild",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/WallpaperAgent.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "Zifferblatt-Hilfe",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/WatchFaceAlert.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "WiFiAgent",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/WiFiAgent.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "WidgetKit Simulator",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/WidgetKit Simulator.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "WindowManager",
      "arch_kind": "arch_arm_i64",
      "info": "WindowManager",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/WindowManager.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "341.1.4"
    },
    {
      "_name": "WindowManagerShowDesktopEducation",
      "arch_kind": "arch_arm_i64",
      "info": "WindowManagerShowDesktopEducation",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/WindowManagerShowDesktopEducation.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "WorkoutAlert-Mac",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/WorkoutAlert-Mac.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "iCloud+",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/iCloud+.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "iCloud",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/iCloud.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "liquiddetectiond",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/liquiddetectiond.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ]
    },
    {
      "_name": "loginwindow",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/loginwindow.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "9.0"
    },
    {
      "_name": "rcd",
      "arch_kind": "arch_arm_i64",
      "info": "362",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/rcd.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "362"
    },
    {
      "_name": "screencaptureui",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/screencaptureui.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "Ordneraktionen konfigurieren",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/Applications/Folder Actions Setup.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.2"
    },
    {
      "_name": "Schlüsselbundverwaltung",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/Applications/Keychain Access.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "11.0"
    },
    {
      "_name": "Ticket-Viewer",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/Applications/Ticket Viewer.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "4.1"
    },
    {
      "_name": "Diagnose für drahtlose Umgebungen",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/Applications/Wireless Diagnostics.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "11.0"
    },
    {
      "_name": "iOS-App-Installationsprogramm",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/Applications/iOS App Installer.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "Finder",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/Finder.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "26.1"
    },
    {
      "_name": "AirDrop",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/Finder.app/Contents/Applications/AirDrop.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "26.1"
    },
    {
      "_name": "Computer",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/Finder.app/Contents/Applications/Computer.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "26.1"
    },
    {
      "_name": "Netzwerk",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/Finder.app/Contents/Applications/Network.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "26.1"
    },
    {
      "_name": "Zuletzt benutzt",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/Finder.app/Contents/Applications/Recents.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "26.1"
    },
    {
      "_name": "iCloud Drive",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/Finder.app/Contents/Applications/iCloud Drive.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "26.1"
    },
    {
      "_name": "Systemassistent",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/Setup Assistant.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "10.10"
    },
    {
      "_name": "ShortcutDroplet",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/ShortcutDroplet.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "Shortcuts Events",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/Shortcuts Events.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "ShortcutsActions",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/ShortcutsActions.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ]
    },
    {
      "_name": "Siri",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/Siri.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "3505.11.1"
    },
    {
      "_name": "Softwareupdate",
      "arch_kind": "arch_arm_i64",
      "info": "Software Update version 4.0, Copyright © 2000-2009, Apple Inc. All rights reserved.",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/Software Update.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "6"
    },
    {
      "_name": "SpacesTouchBarAgent",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/SpacesTouchBarAgent.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "Spotlight",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/Spotlight.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "StageManagerOnboarding",
      "arch_kind": "arch_arm_i64",
      "info": "StageManagerOnboarding",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/StageManagerOnboarding.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "System Events",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/System Events.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.3.6"
    },
    {
      "_name": "SystemIntents",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/SystemIntents.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "Automator Application Stub",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "unknown",
      "path": "/System/Library/CoreServices/Automator Application Stub.app",
      "version": "1.3"
    },
    {
      "_name": "Automator-Installationsprogramm",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/Automator Installer.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "2.10"
    },
    {
      "_name": "Batterien",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/Batteries.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "Bluetooth-Assistent",
      "arch_kind": "arch_arm_i64",
      "info": "9.0 (1)",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/BluetoothSetupAssistant.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "9.0"
    },
    {
      "_name": "BluetoothUIServer",
      "arch_kind": "arch_arm_i64",
      "info": "kBluetoothCFBundleGetInfoString",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/BluetoothUIServer.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "9.0"
    },
    {
      "_name": "BluetoothUIService",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/BluetoothUIService.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "CalendarFileHandler",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/CalendarFileHandler.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "8.0"
    },
    {
      "_name": "Captive Network Assistant",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/Captive Network Assistant.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "5.0"
    },
    {
      "_name": "Zertifikatsassistent",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/Certificate Assistant.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "5.0"
    },
    {
      "_name": "Kontrollzentrum",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/ControlCenter.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "ControlStrip",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/ControlStrip.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "CoreLocationAgent",
      "arch_kind": "arch_arm_i64",
      "info": "Copyright © 2013 Apple Inc.",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/CoreLocationAgent.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "3060.0.18"
    },
    {
      "_name": "CoreServicesUIAgent",
      "arch_kind": "arch_arm_i64",
      "info": "Copyright © 2009 Apple Inc.",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/CoreServicesUIAgent.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "369"
    },
    {
      "_name": "Abdeckungsdetails",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/Coverage Details.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "Database Events",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/Database Events.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0.6"
    },
    {
      "_name": "Diagnostics Reporter",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/Diagnostics Reporter.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "DiscHelper",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/DiscHelper.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "DiskImageMounter",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/DiskImageMounter.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "Dock",
      "arch_kind": "arch_arm_i64",
      "info": "Dock 1.8",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/Dock.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.8"
    },
    {
      "_name": "Dwell Control",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/Dwell Control.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "Erweiterte Protokollierung",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/Enhanced Logging.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "Löschassistent",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/Erase Assistant.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "EscrowSecurityAlert",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/EscrowSecurityAlert.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "Family",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/Family.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "FamilyExtensionHost",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/FamilyExtensionHost.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "FileProvider-Feedback",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/FileProvider-Feedback.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "FolderActionsDispatcher",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/FolderActionsDispatcher.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "Game Center",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/Game Center.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "GameOverlayUI",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/GameOverlayUI.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "Spiele",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/GameTrampoline.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "820.1.14"
    },
    {
      "_name": "IOUIAgent",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/IOUIAgent.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "Image Events",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/Image Events.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.1.6"
    },
    {
      "_name": "Install Command Line Developer Tools",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/Install Command Line Developer Tools.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "2416"
    },
    {
      "_name": "Installation wird durchgeführt",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/Install in Progress.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "3.0"
    },
    {
      "_name": "Installer Progress",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/Installer Progress.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "Installationsprogramm",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/Installer.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "6.2.0"
    },
    {
      "_name": "JavaLauncher",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/JavaLauncher.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "326"
    },
    {
      "_name": "KeyboardAccessAgent",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/KeyboardAccessAgent.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "10"
    },
    {
      "_name": "KeyboardSetupAssistant",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/KeyboardSetupAssistant.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "Keychain Circle Notification",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/Keychain Circle Notification.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "Sprachauswahl",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/Language Chooser.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "MDMMigrationTrampoline",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/MDMMigrationTrampoline.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "MTLReplayer",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/MTLReplayer.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "311.2"
    },
    {
      "_name": "ManagedClient",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/ManagedClient.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "18.0"
    },
    {
      "_name": "Dienstprogramm für Speichersteckplätze",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/Memory Slot Utility.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.5.3"
    },
    {
      "_name": "Musikerkennung",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/MusicRecognitionMac.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "NetAuthAgent",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/NetAuthAgent.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "6.2"
    },
    {
      "_name": "Mitteilungszentrale",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/NotificationCenter.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "NowPlayingTouchUI",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/NowPlayingTouchUI.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "OBEXAgent",
      "arch_kind": "arch_arm_i64",
      "info": "7.0.0, Copyright © 2002-2018 Apple Inc. All rights reserved.",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/OBEXAgent.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "9.0"
    },
    {
      "_name": "ODSAgent",
      "arch_kind": "arch_arm_i64",
      "info": "1.9 (190.2), Copyright © 2007-2009 Apple Inc. All Rights Reserved.",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/ODSAgent.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.9"
    },
    {
      "_name": "OSDUIHelper",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/OSDUIHelper.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "PIPAgent",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/PIPAgent.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "2.0"
    },
    {
      "_name": "Gekoppelte Geräte",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/Paired Devices.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "7.1"
    },
    {
      "_name": "Pass Viewer",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/Pass Viewer.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "PassViewer",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/PassViewer.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ]
    },
    {
      "_name": "PeopleMessageService",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/PeopleMessageService.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "Kontakte",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/PeopleViewService.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "PosterBoard",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/PosterBoard.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "PowerChime",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/PowerChime.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "PreviewShell",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/PreviewShell.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "16.0"
    },
    {
      "_name": "Pro Display-Kalibrierung",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/Pro Display Calibrator.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "211.25"
    },
    {
      "_name": "Problem Reporter",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/Problem Reporter.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "10.13"
    },
    {
      "_name": "Profilinstallation",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/ProfileHelper.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "RapportUIAgent",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/RapportUIAgent.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "7.1"
    },
    {
      "_name": "RegisterPluginIMApp",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/RegisterPluginIMApp.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "26.200"
    },
    {
      "_name": "ARDAgent",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/RemoteManagement/ARDAgent.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "3.9.8"
    },
    {
      "_name": "Remote Desktop-Nachricht",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/RemoteManagement/Remote Desktop Message.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "3.9.8"
    },
    {
      "_name": "SSMenuAgent",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/RemoteManagement/SSMenuAgent.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "3.9.8"
    },
    {
      "_name": "Rosetta 2 Updater",
      "arch_kind": "arch_arm",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/Rosetta 2 Updater.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "Bildschirmzeit",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/Screen Time.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "3.0"
    },
    {
      "_name": "ScreenSaverEngine",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/ScreenSaverEngine.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "5.0"
    },
    {
      "_name": "Skript-Menü",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/Script Menu.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1"
    },
    {
      "_name": "ScriptMonitor",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/ScriptMonitor.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0.1"
    },
    {
      "_name": "CharacterPalette",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/Input Methods/CharacterPalette.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "2.0.1"
    },
    {
      "_name": "Diktat",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/Input Methods/DictationIM.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "6.2.41.5"
    },
    {
      "_name": "EmojiFunctionRowIM",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/Input Methods/EmojiFunctionRowIM.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "JapaneseIM-KanaTyping",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/Input Methods/JapaneseIM-KanaTyping.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "6.3"
    },
    {
      "_name": "JapaneseIM-RomajiTyping",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/Input Methods/JapaneseIM-RomajiTyping.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "6.3"
    },
    {
      "_name": "KoreanIM",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/Input Methods/KoreanIM.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "PluginIM",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/Input Methods/PluginIM.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "26.200"
    },
    {
      "_name": "PressAndHold",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/Input Methods/PressAndHold.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "SCIM",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/Input Methods/SCIM.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "104"
    },
    {
      "_name": "TCIM",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/Input Methods/TCIM.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "104"
    },
    {
      "_name": "TYIM",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/Input Methods/TYIM.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "104"
    },
    {
      "_name": "TamilIM",
      "arch_kind": "arch_arm_i64",
      "info": "Tamil Input Method 1.5",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/Input Methods/TamilIM.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.6"
    },
    {
      "_name": "TrackpadIM",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/Input Methods/TrackpadIM.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "TransliterationIM",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/Input Methods/TransliterationIM.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "VietnameseIM",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/Input Methods/VietnameseIM.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "AOSAlertManager",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/PrivateFrameworks/AOSKit.framework/Versions/A/Helpers/AOSAlertManager.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.07"
    },
    {
      "_name": "AOSHeartbeat",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/PrivateFrameworks/AOSKit.framework/Versions/A/Helpers/AOSHeartbeat.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.07"
    },
    {
      "_name": "AOSPushRelay",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/PrivateFrameworks/AOSKit.framework/Versions/A/Helpers/AOSPushRelay.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.07"
    },
    {
      "_name": "ClassroomStudentMenuExtra",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/Classroom/ClassroomStudentMenuExtra.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "Kalibrierungsassistent",
      "arch_kind": "arch_arm_i64",
      "info": "4.19, Copyright 2014-2024 Apple Computer, Inc.",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/ColorSync/Calibrators/Display Calibrator.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "4.19"
    },
    {
      "_name": "SpeechRecognitionServer",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/PrivateFrameworks/SpeechObjects.framework/Versions/A/SpeechRecognitionServer.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "9.0.78"
    },
    {
      "_name": "DiskImages UI Agent",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/PrivateFrameworks/DiskImages.framework/Versions/A/Resources/DiskImages UI Agent.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "680"
    },
    {
      "_name": "AXVisualSupportAgent",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/PrivateFrameworks/UniversalAccess.framework/Versions/A/Resources/AXVisualSupportAgent.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "Bedienungshilfen-Einführung",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/PrivateFrameworks/UniversalAccess.framework/Versions/A/Resources/Accessibility Tutorial.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "DFRHUD",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/PrivateFrameworks/UniversalAccess.framework/Versions/A/Resources/DFRHUD.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "universalAccessAuthWarn",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/PrivateFrameworks/UniversalAccess.framework/Versions/A/Resources/universalAccessAuthWarn.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "imagent",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/PrivateFrameworks/IMCore.framework/imagent.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "10.0"
    },
    {
      "_name": "IMAutomaticHistoryDeletionAgent",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/PrivateFrameworks/IMDPersistence.framework/IMAutomaticHistoryDeletionAgent.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "10.0"
    },
    {
      "_name": "AutoFillPanelService",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/PrivateFrameworks/AutoFillUI.framework/Contents/AutoFillPanelService.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "AutomationModeUI",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/PrivateFrameworks/AutomationMode.framework/AutomationModeUI.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "FindMyMacMessenger",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/PrivateFrameworks/FindMyMac.framework/Versions/A/Resources/FindMyMacMessenger.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "4.1"
    },
    {
      "_name": "STMUIHelper",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/PrivateFrameworks/StorageManagement.framework/Versions/A/Resources/STMUIHelper.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "privatecloudcomputed",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/PrivateFrameworks/PrivateCloudCompute.framework/privatecloudcomputed.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "Webseite erstellen",
      "arch_kind": "arch_arm_i64",
      "info": "10.1, © Copyright 2003-2014 Apple  Inc. All rights reserved.",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/Image Capture/Automatic Tasks/Build Web Page.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "10.1"
    },
    {
      "_name": "MakePDF",
      "arch_kind": "arch_arm_i64",
      "info": "10.1, © Copyright 2003-2015 Apple Inc. All rights reserved.",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/Image Capture/Automatic Tasks/MakePDF.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "10.1"
    },
    {
      "_name": "AirScanScanner",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/Image Capture/Devices/AirScanScanner.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "18"
    },
    {
      "_name": "50onPaletteServer",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/Input Methods/50onPaletteServer.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.1.0"
    },
    {
      "_name": "AinuIM",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/Input Methods/AinuIM.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "Assistive Control",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/Input Methods/Assistive Control.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "2.0"
    },
    {
      "_name": "BackgroundTaskManagementAgent",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/PrivateFrameworks/BackgroundTaskManagement.framework/Support/BackgroundTaskManagementAgent.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "nbagent",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/PrivateFrameworks/Noticeboard.framework/Versions/A/Resources/nbagent.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "iCloud Drive",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/PrivateFrameworks/iCloudDriveCore.framework/Versions/A/Resources/iCloud Drive.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "iCloudUserNotificationsd",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/PrivateFrameworks/AOSAccounts.framework/Versions/A/Resources/iCloudUserNotificationsd.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "IMTransferAgent",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/PrivateFrameworks/IMTransferServices.framework/IMTransferAgent.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "10.0"
    },
    {
      "_name": "Live-Untertitel",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/PrivateFrameworks/AccessibilitySharedSupport.framework/Versions/A/Resources/Live Captions.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "LiveSpeech",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/PrivateFrameworks/AccessibilitySharedSupport.framework/Versions/A/Resources/LiveSpeech.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "identityservicesd",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/PrivateFrameworks/IDS.framework/identityservicesd.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "10.0"
    },
    {
      "_name": "IDSRemoteURLConnectionAgent",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/PrivateFrameworks/IDSFoundation.framework/IDSRemoteURLConnectionAgent.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "10.0"
    },
    {
      "_name": "eaptlstrust",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/PrivateFrameworks/EAP8021X.framework/Support/eaptlstrust.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "13.0"
    },
    {
      "_name": "storeuid",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/PrivateFrameworks/CommerceKit.framework/Versions/A/Resources/storeuid.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "ScreenReaderUIServer",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/PrivateFrameworks/ScreenReader.framework/Versions/A/Resources/ScreenReaderUIServer.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "10"
    },
    {
      "_name": "VoiceOver-Kurzübersicht",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/PrivateFrameworks/ScreenReader.framework/Versions/A/Resources/VoiceOver Quickstart.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "10"
    },
    {
      "_name": "sociallayerd",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/PrivateFrameworks/SocialLayer.framework/sociallayerd.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "SoftwareUpdateNotificationManager",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/PrivateFrameworks/SoftwareUpdate.framework/Versions/A/Resources/SoftwareUpdateNotificationManager.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "Sprach-Downloader",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/PrivateFrameworks/SpeechObjects.framework/Versions/A/SpeechDataInstallerd.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "9.0.78"
    },
    {
      "_name": "UASharedPasteboardProgressUI",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/PrivateFrameworks/UserActivity.framework/Agents/UASharedPasteboardProgressUI.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "54.1"
    },
    {
      "_name": "ParentalControls",
      "arch_kind": "arch_arm_i64",
      "info": "2.0, Copyright Apple Inc. 2007-2019",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/PrivateFrameworks/FamilyControls.framework/Versions/A/Resources/ParentalControls.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "4.1"
    },
    {
      "_name": "FeedbackRemoteView",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/PrivateFrameworks/FeedbackService.framework/Versions/A/Support/FeedbackRemoteView.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "FollowUpUI",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/PrivateFrameworks/CoreFollowUp.framework/Versions/A/Resources/FollowUpUI.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "AppSSOAgent",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/PrivateFrameworks/AppSSO.framework/Support/AppSSOAgent.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "KerberosMenuExtra",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/PrivateFrameworks/AppSSOKerberos.framework/Support/KerberosMenuExtra.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "AskPermissionUI",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/PrivateFrameworks/AskPermission.framework/Versions/A/Resources/AskPermissionUI.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "Konfliktlöser",
      "arch_kind": "arch_arm_i64",
      "info": "1.0, Copyright Apple Computer Inc. 2004",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/PrivateFrameworks/SyncServicesUI.framework/Versions/A/Resources/Conflict Resolver.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "8.1"
    },
    {
      "_name": "syncuid",
      "arch_kind": "arch_arm_i64",
      "info": "4.0, Copyright Apple Computer Inc. 2004",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/PrivateFrameworks/SyncServicesUI.framework/Versions/A/Resources/syncuid.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "8.1"
    },
    {
      "_name": "Calibration Assistant",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/PrivateFrameworks/AmbientDisplay.framework/Versions/A/Resources/Calibration Assistant.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "AMSEngagementViewService",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/PrivateFrameworks/AppleMediaServicesUI.framework/Versions/A/Resources/AMSEngagementViewService.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "AquaAppearanceHelper",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/PrivateFrameworks/SkyLight.framework/Versions/A/Resources/AquaAppearanceHelper.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "CIMFindInputCodeTool",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/PrivateFrameworks/CoreChineseEngine.framework/Versions/A/SharedSupport/CIMFindInputCodeTool.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "104"
    },
    {
      "_name": "AccessibilityVisualsAgent",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/PrivateFrameworks/AccessibilitySupport.framework/Versions/A/Resources/AccessibilityVisualsAgent.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "Systemsprachausgabe",
      "arch_kind": "arch_arm_i64",
      "info": "9.2.22",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/Frameworks/ApplicationServices.framework/Versions/A/Frameworks/SpeechSynthesis.framework/Versions/A/System Speech.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "9.2.22"
    },
    {
      "_name": "qlmanage",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/Frameworks/QuickLook.framework/Versions/A/Resources/qlmanage.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "quicklookd",
      "arch_kind": "arch_arm_i64",
      "info": "5.0, Copyright Apple Inc. 2007-2013",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/Frameworks/QuickLook.framework/Versions/A/Resources/quicklookd.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "5.0"
    },
    {
      "_name": "SyncServer",
      "arch_kind": "arch_arm_i64",
      "info": "© 2002-2003 Apple",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/Frameworks/SyncServices.framework/Versions/A/Resources/SyncServer.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "8.1"
    },
    {
      "_name": "FontRegistryUIAgent",
      "arch_kind": "arch_arm_i64",
      "info": "Copyright © 2008-2013 Apple Inc.",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/Frameworks/ApplicationServices.framework/Versions/A/Frameworks/ATS.framework/Versions/A/Support/FontRegistryUIAgent.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "81.0"
    },
    {
      "_name": "SmartCard-Verknüpfung",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/Frameworks/CryptoTokenKit.framework/ctkbind.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "Quick Look Simulator",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/Frameworks/QuickLookUI.framework/Versions/A/Resources/Quick Look Simulator.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "QuickLookUIHelper",
      "arch_kind": "arch_arm_i64",
      "info": "5.0, Copyright Apple Inc. 2007-2013",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/Frameworks/QuickLookUI.framework/Versions/A/Resources/QuickLookUIHelper.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "5.0"
    },
    {
      "_name": "Wish",
      "arch_kind": "arch_arm_i64",
      "info": "Wish Shell 8.5.9,\nCopyright © 1989-2025 Tcl Core Team,\nCopyright © 2002-2025 Daniel A. Steffen,\nCopyright © 2001-2009 Apple Inc.,\nCopyright © 2001-2002 Jim Ingham & Ian Reid",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/Frameworks/Tk.framework/Versions/8.5/Resources/Wish.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "8.5.9"
    },
    {
      "_name": "LinkedNotesUIService",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/Frameworks/PaperKit.framework/Contents/LinkedNotesUIService.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "CinematicFramingOnboardingUI",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/Frameworks/CoreMediaIO.framework/Versions/A/Resources/CinematicFramingOnboardingUI.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "ContinuityCaptureOnboardingUI",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/Frameworks/CoreMediaIO.framework/Versions/A/Resources/ContinuityCaptureOnboardingUI.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "AddressBookSync",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/Frameworks/AddressBook.framework/Helpers/AddressBookSync.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "14.0"
    },
    {
      "_name": "ABAssistantService",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/Frameworks/AddressBook.framework/Versions/A/Helpers/ABAssistantService.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "14.0"
    },
    {
      "_name": "AddressBookManager",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/Frameworks/AddressBook.framework/Versions/A/Helpers/AddressBookManager.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "14.0"
    },
    {
      "_name": "AddressBookSourceSync",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/Frameworks/AddressBook.framework/Versions/A/Helpers/AddressBookSourceSync.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "14.0"
    },
    {
      "_name": "AppleSpell",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/Services/AppleSpell.service",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "2.4"
    },
    {
      "_name": "ChinesischerTextKonvertierungsService",
      "arch_kind": "arch_arm_i64",
      "info": "Chinese Text Converter 1.1",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/Services/ChineseTextConverterService.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "2.1"
    },
    {
      "_name": "OpenSpell",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/Services/OpenSpell.service",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "SpeechService",
      "arch_kind": "arch_arm_i64",
      "info": "9.2.22",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/Services/SpeechService.service",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "9.2.22"
    },
    {
      "_name": "Spotlight",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/Services/Spotlight.service",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "3.0"
    },
    {
      "_name": "Zusammenfassung",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/Services/Summary Service.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "2.0"
    },
    {
      "_name": "AOSUIPrefPaneLauncher",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/AOSUIPrefPaneLauncher.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "AVB-Konfiguration",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/AVB Configuration.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1400.10"
    },
    {
      "_name": "Reader für Bedienungshilfen",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/Accessibility Reader.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "AccessibilityUIServer",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/AccessibilityUIServer.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "AddPrinter",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/AddPrinter.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "607"
    },
    {
      "_name": "AddressBookUrlForwarder",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/AddressBookUrlForwarder.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "14.0"
    },
    {
      "_name": "AirPlayUIAgent",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/AirPlayUIAgent.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "2.0"
    },
    {
      "_name": "AirPort-Basisstation-Agent",
      "arch_kind": "arch_arm_i64",
      "info": "2.2.1 (221.12), Copyright © 2006 -2024 Apple Inc. All Rights Reserved.",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/AirPort Base Station Agent.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "2.2.1"
    },
    {
      "_name": "Apple Diagnose",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/Apple Diagnostics.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "AppleScript-Dienstprogramm",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Library/CoreServices/AppleScript Utility.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.1.2"
    },
    {
      "_name": "Lexikon",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/Dictionary.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "2.3.0"
    },
    {
      "_name": "FaceTime",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/FaceTime.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "36"
    },
    {
      "_name": "Wo ist?",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/FindMy.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "4.0"
    },
    {
      "_name": "Schriftsammlung",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/Font Book.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "11.0"
    },
    {
      "_name": "Freeform",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/Freeform.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "4.1"
    },
    {
      "_name": "Spiele",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/Games.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "Home",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/Home.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "10.0"
    },
    {
      "_name": "Digitale Bilder",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/Image Capture.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "8.0"
    },
    {
      "_name": "Image Playground",
      "arch_kind": "arch_arm",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/Image Playground.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "Journal",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/Journal.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "2.0"
    },
    {
      "_name": "Mail",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/Mail.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "16.0"
    },
    {
      "_name": "Karten",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/Maps.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "3.0"
    },
    {
      "_name": "Nachrichten",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/Messages.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "26.0"
    },
    {
      "_name": "Mission Control",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/Mission Control.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.2"
    },
    {
      "_name": "Musik",
      "arch_kind": "arch_arm_i64",
      "info": "Music 1.6.1.44, © 2019–2025 Apple Inc. All rights reserved.",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/Music.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.6.1"
    },
    {
      "_name": "Notizen",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/System/Applications/Notes.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "4.13"
    },
    {
      "_name": "iTerm",
      "arch_kind": "arch_arm_i64",
      "info": "3.6.6",
      "lastModified": "2025-11-18T07:14:07Z",
      "obtained_from": "identified_developer",
      "path": "/Applications/iTerm.app",
      "signed_by": [
        "Developer ID Application: GEORGE NACHMAN (H7V7XYVQ7D)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "3.6.6"
    },
    {
      "_name": "Citrix Workspace",
      "arch_kind": "arch_arm_i64",
      "info": "25.08.10",
      "lastModified": "2025-11-07T12:10:22Z",
      "obtained_from": "identified_developer",
      "path": "/Applications/Citrix Workspace.app",
      "signed_by": [
        "Developer ID Application: Citrix Systems, Inc. (S272Y5R93J)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "25.08.10"
    },
    {
      "_name": "Deinstallieren von Citrix Workspace",
      "arch_kind": "arch_arm_i64",
      "info": "25.08.10",
      "lastModified": "2025-11-07T12:10:22Z",
      "obtained_from": "identified_developer",
      "path": "/Library/Application Support/Citrix Receiver/Uninstall Citrix Workspace.app",
      "signed_by": [
        "Developer ID Application: Citrix Systems, Inc. (S272Y5R93J)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "25.08.10"
    },
    {
      "_name": "Citrix Enterprise Browser",
      "arch_kind": "arch_arm_i64",
      "info": "Citrix Enterprise Browser 139.1.1.27, Copyright © 2021-2025. Citrix Systems, Inc. All rights reserved.",
      "lastModified": "2025-11-07T12:10:22Z",
      "obtained_from": "identified_developer",
      "path": "/Library/Application Support/Citrix Receiver/Citrix Enterprise Browser.app",
      "signed_by": [
        "Developer ID Application: Citrix Systems, Inc. (S272Y5R93J)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "139.1.1.27"
    },
    {
      "_name": "CitrixEndpointAnalysis",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-11-07T12:10:19Z",
      "obtained_from": "identified_developer",
      "path": "/Users/ibering/Library/Application Support/Citrix/EPAPlugin/CitrixEndpointAnalysis.app",
      "signed_by": [
        "Developer ID Application: Citrix Systems, Inc. (S272Y5R93J)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "25.6.10"
    },
    {
      "_name": "deviceTRUST Location",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-11-07T12:10:21Z",
      "obtained_from": "identified_developer",
      "path": "/Library/Application Support/deviceTRUST/ICA Client/Bin/deviceTRUST Location.app",
      "signed_by": [
        "Developer ID Application: Citrix Systems, Inc. (S272Y5R93J)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "25.3"
    },
    {
      "_name": "FortiClient",
      "arch_kind": "arch_arm",
      "lastModified": "2024-12-06T23:39:32Z",
      "obtained_from": "identified_developer",
      "path": "/Applications/FortiClient.app",
      "signed_by": [
        "Developer ID Application: Fortinet, Inc (AH4XFXJ7DK)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "7.4.2.1717"
    },
    {
      "_name": "FortiClientUninstaller",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2024-12-06T23:39:27Z",
      "obtained_from": "identified_developer",
      "path": "/Applications/FortiClientUninstaller.app",
      "signed_by": [
        "Developer ID Application: Fortinet, Inc (AH4XFXJ7DK)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.2"
    },
    {
      "_name": "Google Chrome",
      "arch_kind": "arch_arm_i64",
      "info": "Google Chrome 143.0.7499.170, Copyright 2025 Google LLC. All rights reserved.",
      "lastModified": "2025-03-08T01:26:02Z",
      "obtained_from": "identified_developer",
      "path": "/Applications/Google Chrome.app",
      "signed_by": [
        "Developer ID Application: Google LLC (EQHXZ8M8AV)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "143.0.7499.170"
    },
    {
      "_name": "Firefox",
      "arch_kind": "arch_arm_i64",
      "info": "Firefox 146.0.1",
      "lastModified": "2026-01-05T06:19:03Z",
      "obtained_from": "identified_developer",
      "path": "/Applications/Firefox.app",
      "signed_by": [
        "Developer ID Application: Mozilla Corporation (43AQ936H96)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "146.0.1"
    },
    {
      "_name": "Screenshot Capture",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2024-08-04T10:31:41Z",
      "obtained_from": "unknown",
      "path": "/Users/ibering/Applications/Screenshot Capture.app",
      "version": "1.0"
    },
    {
      "_name": "DockDoor",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-02-21T01:35:50Z",
      "obtained_from": "identified_developer",
      "path": "/Applications/DockDoor.app",
      "signed_by": [
        "Developer ID Application: Ethan Bills (2Q775S63Q3)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.6.2"
    },
    {
      "_name": "DisplayLink End-User Cleaner",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2024-01-18T11:57:44Z",
      "obtained_from": "identified_developer",
      "path": "/Users/ibering/Downloads/DisplayLink End-User Cleaner.app",
      "signed_by": [
        "Developer ID Application: DisplayLink Corp (73YQY62QM3)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0.0"
    },
    {
      "_name": "The Unarchiver",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-03-18T08:40:07Z",
      "obtained_from": "identified_developer",
      "path": "/Applications/The Unarchiver.app",
      "signed_by": [
        "Developer ID Application: MacPaw Way Ltd (S8EX82NJP6)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "4.3.9"
    },
    {
      "_name": "Sublime Text",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-05-23T07:55:45Z",
      "obtained_from": "identified_developer",
      "path": "/Applications/Sublime Text.app",
      "signed_by": [
        "Developer ID Application: Sublime HQ Pty Ltd (Z6D26JE4Y4)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "Build 4200"
    },
    {
      "_name": "JetBrains Toolbox",
      "arch_kind": "arch_arm",
      "info": "",
      "lastModified": "2025-12-18T10:36:16Z",
      "obtained_from": "identified_developer",
      "path": "/Applications/JetBrains Toolbox.app",
      "signed_by": [
        "Developer ID Application: JetBrains s.r.o. (2ZEFAR8TH3)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "3.2.0.65851"
    },
    {
      "_name": "GoogleUpdater",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-12-03T06:17:35Z",
      "obtained_from": "identified_developer",
      "path": "/Users/ibering/Library/Application Support/Google/GoogleUpdater/144.0.7547.0/GoogleUpdater.app",
      "signed_by": [
        "Developer ID Application: Google LLC (EQHXZ8M8AV)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "144.0.7547.0"
    },
    {
      "_name": "AppleMobileDeviceHelper",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/Library/Apple/System/Library/PrivateFrameworks/MobileDevice.framework/Versions/A/AppleMobileDeviceHelper.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "5.0"
    },
    {
      "_name": "MobileDeviceUpdater",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/Library/Apple/System/Library/PrivateFrameworks/MobileDevice.framework/Versions/A/Resources/MobileDeviceUpdater.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0"
    },
    {
      "_name": "AppleMobileSync",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/Library/Apple/System/Library/PrivateFrameworks/MobileDevice.framework/Versions/A/AppleMobileSync.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "5.0"
    },
    {
      "_name": "AirScanLegacyDiscovery",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "apple",
      "path": "/Library/Image Capture/Support/LegacyDeviceDiscoveryHelpers/AirScanLegacyDiscovery.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "607"
    },
    {
      "_name": "Recursive File Processing Droplet",
      "arch_kind": "arch_other",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "unknown",
      "path": "/Library/Application Support/Script Editor/Templates/Droplets/Recursive File Processing Droplet.app",
      "version": "1.0"
    },
    {
      "_name": "Droplet with Settable Properties",
      "arch_kind": "arch_other",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "unknown",
      "path": "/Library/Application Support/Script Editor/Templates/Droplets/Droplet with Settable Properties.app",
      "version": "1.0"
    },
    {
      "_name": "Recursive Image File Processing Droplet",
      "arch_kind": "arch_other",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "unknown",
      "path": "/Library/Application Support/Script Editor/Templates/Droplets/Recursive Image File Processing Droplet.app",
      "version": "1.0"
    },
    {
      "_name": "Cocoa-AppleScript Applet",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-29T01:21:04Z",
      "obtained_from": "unknown",
      "path": "/Library/Application Support/Script Editor/Templates/Cocoa-AppleScript Applet.app",
      "version": "1.0"
    },
    {
      "_name": "GarageBand",
      "arch_kind": "arch_arm_i64",
      "info": "GarageBand 10.4.13, Copyright © 2004–2025 Apple Inc. All Rights Reserved",
      "lastModified": "2025-12-17T07:49:16Z",
      "obtained_from": "mac_app_store",
      "path": "/Applications/GarageBand.app",
      "signed_by": [
        "Apple Mac OS Application Signing",
        "Apple Worldwide Developer Relations Certification Authority",
        "Apple Root CA"
      ],
      "version": "10.4.13"
    },
    {
      "_name": "Accessibility Inspector",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-12-15T07:41:44Z",
      "obtained_from": "apple",
      "path": "/Applications/Xcode.app/Contents/Applications/Accessibility Inspector.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "5.0"
    },
    {
      "_name": "DisplayLink Manager",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-12-15T13:47:00Z",
      "obtained_from": "identified_developer",
      "path": "/Applications/DisplayLink Manager.app",
      "signed_by": [
        "Developer ID Application: DisplayLink Corp (73YQY62QM3)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "14.2.0"
    },
    {
      "_name": "Microsoft AutoUpdate",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-12-15T07:40:03Z",
      "obtained_from": "identified_developer",
      "path": "/Library/Application Support/Microsoft/MAU2.0/Microsoft AutoUpdate.app",
      "signed_by": [
        "Developer ID Application: Microsoft Corporation (UBF8T346G9)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "4.81.2"
    },
    {
      "_name": "EdgeUpdater",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-07-11T13:57:15Z",
      "obtained_from": "identified_developer",
      "path": "/Library/Application Support/Microsoft/EdgeUpdater/137.0.3249.0/EdgeUpdater.app",
      "signed_by": [
        "Developer ID Application: Microsoft Corporation (UBF8T346G9)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "137.0.3249.0"
    },
    {
      "_name": "Microsoft Defender",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-12-19T07:03:42Z",
      "obtained_from": "identified_developer",
      "path": "/Applications/Microsoft Defender.app",
      "signed_by": [
        "Developer ID Application: Microsoft Corporation (UBF8T346G9)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "101.25102.0019"
    },
    {
      "_name": "Microsoft Intune Agent",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-12-10T06:35:35Z",
      "obtained_from": "identified_developer",
      "path": "/Library/Intune/Microsoft Intune Agent.app",
      "signed_by": [
        "Developer ID Application: Microsoft Corporation (UBF8T346G9)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "2512.003"
    },
    {
      "_name": "OpenVPN Connect",
      "arch_kind": "arch_i64",
      "lastModified": "2026-01-08T10:55:49Z",
      "obtained_from": "identified_developer",
      "path": "/Applications/OpenVPN Connect/OpenVPN Connect.app",
      "signed_by": [
        "Developer ID Application: OPENVPN TECHNOLOGIES, INC. (ACV7L3WCD8)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "3.2.7"
    },
    {
      "_name": "Uninstall OpenVPN Connect",
      "arch_kind": "arch_i64",
      "info": "OpenVPN Uninstaller",
      "lastModified": "2026-01-08T10:55:51Z",
      "obtained_from": "identified_developer",
      "path": "/Applications/OpenVPN Connect/Uninstall OpenVPN Connect.app",
      "signed_by": [
        "Developer ID Application: OPENVPN TECHNOLOGIES, INC. (ACV7L3WCD8)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "0.1"
    },
    {
      "_name": "BetterDisplay",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-12-09T06:44:01Z",
      "obtained_from": "identified_developer",
      "path": "/Applications/BetterDisplay.app",
      "signed_by": [
        "Developer ID Application: Istvan Toth (299YSU96J7)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "4.1.1"
    },
    {
      "_name": "EPSON Scanner",
      "arch_kind": "arch_i64",
      "info": "5.7.24, Copyright(C) Seiko Epson Corporation 2002-2015 All rights reserved.",
      "lastModified": "2025-12-01T20:31:10Z",
      "obtained_from": "identified_developer",
      "path": "/Library/Image Capture/Devices/EPSON Scanner.app",
      "signed_by": [
        "Developer ID Application: EPSON (TXAEAV5RN4)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "5.7.24"
    },
    {
      "_name": "Discord",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-12-15T09:19:14Z",
      "obtained_from": "identified_developer",
      "path": "/Applications/Discord.app",
      "signed_by": [
        "Developer ID Application: Discord, Inc. (53Q6R32WPB)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "0.0.371"
    },
    {
      "_name": "OneDrive",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-12-16T06:30:07Z",
      "obtained_from": "identified_developer",
      "path": "/Applications/OneDrive.app",
      "signed_by": [
        "Developer ID Application: Microsoft Corporation (UBF8T346G9)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "25.222.1112"
    },
    {
      "_name": "MRT",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-12-02T10:28:21Z",
      "obtained_from": "apple",
      "path": "/Library/Apple/System/Library/CoreServices/MRT.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.93"
    },
    {
      "_name": "XProtect",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-12-02T10:28:27Z",
      "obtained_from": "apple",
      "path": "/Library/Apple/System/Library/CoreServices/XProtect.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "156"
    },
    {
      "_name": "Unternehmensportal",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-11-26T06:00:13Z",
      "obtained_from": "identified_developer",
      "path": "/Applications/Company Portal.app",
      "signed_by": [
        "Developer ID Application: Microsoft Corporation (UBF8T346G9)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "5.2510.1"
    },
    {
      "_name": "PhpStorm",
      "arch_kind": "arch_ios",
      "info": "PhpStorm 2025.2.5, build PS-252.28238.9. Copyright JetBrains s.r.o., (c) 2000-2025",
      "lastModified": "2025-12-02T09:40:36Z",
      "obtained_from": "identified_developer",
      "path": "/Users/ibering/Applications/PhpStorm.app",
      "signed_by": [
        "Developer ID Application: JetBrains s.r.o. (2ZEFAR8TH3)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "2025.2.5"
    },
    {
      "_name": "Sequel Ace",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2024-11-22T07:11:10Z",
      "obtained_from": "mac_app_store",
      "path": "/Applications/Sequel Ace.app",
      "signed_by": [
        "Apple Mac OS Application Signing",
        "Apple Worldwide Developer Relations Certification Authority",
        "Apple Root CA"
      ],
      "version": "4.1.5"
    },
    {
      "_name": "LimeChat",
      "arch_kind": "arch_i64",
      "info": "LimeChat for Mac, Copyright 2007-2020 Satoshi Nakagawa",
      "lastModified": "2021-04-28T06:46:41Z",
      "obtained_from": "mac_app_store",
      "path": "/Applications/LimeChat.app",
      "signed_by": [
        "Apple Mac OS Application Signing",
        "Apple Worldwide Developer Relations Certification Authority",
        "Apple Root CA"
      ],
      "version": "2.47"
    },
    {
      "_name": "VisualXML",
      "arch_kind": "arch_i64",
      "lastModified": "2021-04-28T06:46:34Z",
      "obtained_from": "mac_app_store",
      "path": "/Applications/VisualXML.app",
      "signed_by": [
        "Apple Mac OS Application Signing",
        "Apple Worldwide Developer Relations Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.4.1"
    },
    {
      "_name": "draw.io",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2024-12-03T11:25:10Z",
      "obtained_from": "identified_developer",
      "path": "/Applications/draw.io.app",
      "signed_by": [
        "Developer ID Application: JGraph Ltd (UZEUFB4N53)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "25.0.2"
    },
    {
      "_name": "Imagine",
      "arch_kind": "arch_i64",
      "lastModified": "2022-12-06T08:14:32Z",
      "obtained_from": "unknown",
      "path": "/Applications/Imagine.app",
      "version": "0.7.3"
    },
    {
      "_name": "Microsoft Remote Desktop",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2024-08-29T06:36:45Z",
      "obtained_from": "mac_app_store",
      "path": "/Applications/Microsoft Remote Desktop.app",
      "signed_by": [
        "Apple Mac OS Application Signing",
        "Apple Worldwide Developer Relations Certification Authority",
        "Apple Root CA"
      ],
      "version": "10.9.10"
    },
    {
      "_name": "GoToMeeting",
      "arch_kind": "arch_i64",
      "info": "GoToMeeting v10.20.0.19992, Copyright © 2024 LogMeIn, Inc.",
      "lastModified": "2024-04-10T12:35:46Z",
      "obtained_from": "identified_developer",
      "path": "/Applications/GoToMeeting.app",
      "signed_by": [
        "Developer ID Application: LogMeIn, Inc. (GFNFVT632V)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "10.20.0.19992"
    },
    {
      "_name": "Affinity Publisher 2",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2023-11-03T09:55:53Z",
      "obtained_from": "identified_developer",
      "path": "/Applications/Affinity Publisher 2.app",
      "signed_by": [
        "Developer ID Application: Serif (Europe) Ltd. (6LVTQB9699)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "2.2.1"
    },
    {
      "_name": "Skype for Business",
      "arch_kind": "arch_i64",
      "lastModified": "2025-08-19T08:49:47Z",
      "obtained_from": "identified_developer",
      "path": "/Applications/Skype for Business.app",
      "signed_by": [
        "Developer ID Application: Skype Communications S.a.r.l (AL798K98FX)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "16.31.11"
    },
    {
      "_name": "Spotify",
      "arch_kind": "arch_arm",
      "lastModified": "2025-08-20T15:03:00Z",
      "obtained_from": "identified_developer",
      "path": "/Applications/Spotify.app",
      "signed_by": [
        "Developer ID Application: Spotify (2FNC3A47ZF)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.2.70.409"
    },
    {
      "_name": "Docker",
      "arch_kind": "arch_arm",
      "lastModified": "2025-07-15T09:51:25Z",
      "obtained_from": "identified_developer",
      "path": "/Applications/Docker.app",
      "signed_by": [
        "Developer ID Application: Docker Inc (9BNSXJN65R)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "4.43.2"
    },
    {
      "_name": "com.microsoft.dlp.daemon",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-12-19T07:03:42Z",
      "obtained_from": "identified_developer",
      "path": "/Library/Application Support/Microsoft/DLP/com.microsoft.dlp.daemon.app",
      "signed_by": [
        "Developer ID Application: Microsoft Corporation (UBF8T346G9)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.25102.103"
    },
    {
      "_name": "com.microsoft.dlp.ux",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-12-19T07:03:42Z",
      "obtained_from": "identified_developer",
      "path": "/Library/Application Support/Microsoft/DLP/com.microsoft.dlp.ux.app",
      "signed_by": [
        "Developer ID Application: Microsoft Corporation (UBF8T346G9)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.25102.103"
    },
    {
      "_name": "com.microsoft.dlp.agent",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-12-19T07:03:42Z",
      "obtained_from": "identified_developer",
      "path": "/Library/Application Support/Microsoft/DLP/com.microsoft.dlp.agent.app",
      "signed_by": [
        "Developer ID Application: Microsoft Corporation (UBF8T346G9)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.25102.103"
    },
    {
      "_name": "Electron",
      "arch_kind": "arch_i64",
      "lastModified": "2024-12-04T08:49:19Z",
      "obtained_from": "unknown",
      "path": "/Users/ibering/openITCOCKPIT-desktop/node_modules/electron/dist/Electron.app",
      "version": "33.2.1"
    },
    {
      "_name": "Install Spotify",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2023-01-31T14:02:42Z",
      "obtained_from": "identified_developer",
      "path": "/Users/ibering/Downloads/Install Spotify.app",
      "signed_by": [
        "Developer ID Application: Spotify (2FNC3A47ZF)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.2.4.912.g949d5fd0"
    },
    {
      "_name": "Affinity Photo 2",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-11-03T09:24:46Z",
      "obtained_from": "identified_developer",
      "path": "/Applications/Affinity Photo 2.app",
      "signed_by": [
        "Developer ID Application: Serif (Europe) Ltd. (6LVTQB9699)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "2.6.5"
    },
    {
      "_name": "YouTube",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-11-13T12:01:37Z",
      "obtained_from": "unknown",
      "path": "/Users/ibering/Applications/Chrome Apps.localized/YouTube.app",
      "version": ""
    },
    {
      "_name": "jquery",
      "arch_kind": "arch_other",
      "lastModified": "2023-10-18T09:57:44Z",
      "obtained_from": "unknown",
      "path": "/Users/ibering/GIT/ngUpgrade/node_modules/angular-ui-router/bower_components/DefinitelyTyped/jquery.placeholder"
    },
    {
      "_name": "Double-Click To Start Support Session",
      "arch_kind": "arch_arm_i64",
      "info": "23.2.2.1c736bc987b79d38772272b47827c465b5b1b888 Copyright © 2002-2023 BeyondTrust Corporation. Redistribution Prohibited. All Rights Reserved.",
      "lastModified": "2025-07-11T12:14:02Z",
      "obtained_from": "identified_developer",
      "path": "/Applications/Double-Click To Start Support Session.app",
      "signed_by": [
        "Developer ID Application: Bomgar (B65TM49E24)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "23.2.2.1c736bc987b79d38772272b47827c465b5b1b888"
    },
    {
      "_name": "FileZilla",
      "arch_kind": "arch_i64",
      "info": "FileZilla Client 3.67.1, Copyright (C) 2004-2024  Tim Kosse, Website: https://filezilla-project.org",
      "lastModified": "2024-07-10T14:23:16Z",
      "obtained_from": "identified_developer",
      "path": "/Applications/FileZilla.app",
      "signed_by": [
        "Developer ID Application: Tim Kosse (5VPGKXL75N)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "3.67.1"
    },
    {
      "_name": "OBS",
      "arch_kind": "arch_i64",
      "lastModified": "2021-10-04T17:54:59Z",
      "obtained_from": "identified_developer",
      "path": "/Applications/OBS.app",
      "signed_by": [
        "Developer ID Application: Wizards of OBS LLC (2MMRE5MTB8)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "27.1.3"
    },
    {
      "_name": "Inkscape",
      "arch_kind": "arch_i64",
      "lastModified": "2021-01-15T17:29:23Z",
      "obtained_from": "identified_developer",
      "path": "/Applications/Inkscape.app",
      "signed_by": [
        "Developer ID Application: Rene de Hesselle (SW3D6BB6A6)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0.2 (e86c8708)"
    },
    {
      "_name": "VisualDesigner",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2023-03-13T07:41:58Z",
      "obtained_from": "mac_app_store",
      "path": "/Applications/VisualDesigner.app",
      "signed_by": [
        "Apple Mac OS Application Signing",
        "Apple Worldwide Developer Relations Certification Authority",
        "Apple Root CA"
      ],
      "version": "5.11"
    },
    {
      "_name": "TeamViewer",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2021-08-12T14:05:13Z",
      "obtained_from": "identified_developer",
      "path": "/Applications/TeamViewer.app",
      "signed_by": [
        "Developer ID Application: TeamViewer GmbH (H7UGFBUGV6)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "15.21.4"
    },
    {
      "_name": "zoom.us",
      "arch_kind": "arch_i64",
      "lastModified": "2024-08-15T07:59:18Z",
      "obtained_from": "identified_developer",
      "path": "/Applications/zoom.us.app",
      "signed_by": [
        "Developer ID Application: Zoom Video Communications, Inc. (BJ4HAAB9B3)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "6.1.6 (37851)"
    },
    {
      "_name": "Visual Studio Code",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-10-14T22:51:18Z",
      "obtained_from": "identified_developer",
      "path": "/Applications/Visual Studio Code.app",
      "signed_by": [
        "Developer ID Application: Microsoft Corporation (UBF8T346G9)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.105.1"
    },
    {
      "_name": "coconutBattery",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2024-04-22T17:34:38Z",
      "obtained_from": "identified_developer",
      "path": "/Applications/coconutBattery.app",
      "signed_by": [
        "Developer ID Application: Christoph Sinai (R5SC3K86L5)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "3.9.17"
    },
    {
      "_name": "VirtualBox",
      "arch_kind": "arch_ios",
      "info": "Oracle VirtualBox Manager 7.1.6, © 2007-2025 Oracle and/or its affiliates",
      "lastModified": "2025-03-17T14:12:50Z",
      "obtained_from": "identified_developer",
      "path": "/Applications/VirtualBox.app",
      "signed_by": [
        "Developer ID Application: Oracle America, Inc. (VB5E2TV963)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "7.1.6"
    },
    {
      "_name": "Uninstall_Citrix_Endpoint_Analysis",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-11-07T12:10:19Z",
      "obtained_from": "identified_developer",
      "path": "/Users/ibering/Library/Application Support/Citrix/EPAPlugin/Uninstall_Citrix_Endpoint_Analysis.app",
      "signed_by": [
        "Developer ID Application: Citrix Systems, Inc. (S272Y5R93J)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "23.11.16"
    },
    {
      "_name": "Canon IJScanner6",
      "arch_kind": "arch_i64",
      "info": "Canon IJScanner6 version 4.0.0, Copyright CANON INC. 2009-2014",
      "lastModified": "2014-06-26T03:00:00Z",
      "obtained_from": "identified_developer",
      "path": "/Library/Image Capture/Devices/Canon IJScanner6.app",
      "signed_by": [
        "Developer ID Application: Canon Inc. (XE2XNRRXZ5)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "4.0.0"
    },
    {
      "_name": "https+++ngxtension.netlify",
      "arch_kind": "arch_other",
      "lastModified": "2025-09-24T05:23:59Z",
      "obtained_from": "unknown",
      "path": "/Users/ibering/Library/Application Support/Firefox/Profiles/hlq66utr.default-release/storage/default/https+++ngxtension.netlify.app"
    },
    {
      "_name": "https+++nice-angular-authguard.netlify",
      "arch_kind": "arch_other",
      "lastModified": "2025-09-24T05:23:58Z",
      "obtained_from": "unknown",
      "path": "/Users/ibering/Library/Application Support/Firefox/Profiles/hlq66utr.default-release/storage/default/https+++nice-angular-authguard.netlify.app"
    },
    {
      "_name": "iMovie",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-11-12T07:08:57Z",
      "obtained_from": "mac_app_store",
      "path": "/Applications/iMovie.app",
      "signed_by": [
        "Apple Mac OS Application Signing",
        "Apple Worldwide Developer Relations Certification Authority",
        "Apple Root CA"
      ],
      "version": "10.4.3"
    },
    {
      "_name": "https+++4d9qp0.csb",
      "arch_kind": "arch_other",
      "lastModified": "2025-09-24T05:24:10Z",
      "obtained_from": "unknown",
      "path": "/Users/ibering/Library/Application Support/Firefox/Profiles/hlq66utr.default-release/storage/default/https+++4d9qp0.csb.app"
    },
    {
      "_name": "NativeMessagingHost",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-05-28T06:40:03Z",
      "obtained_from": "identified_developer",
      "path": "/Library/Application Support/Adobe/WebExtnUtils/NativeMessagingHost.app",
      "signed_by": [
        "Developer ID Application: Adobe Inc. (JQ525L2MZD)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "6.0"
    },
    {
      "_name": "GoToOpener",
      "arch_kind": "arch_i64",
      "lastModified": "2021-04-21T15:02:59Z",
      "obtained_from": "identified_developer",
      "path": "/Users/ibering/Library/Application Support/GoToOpener/GoToOpener.app",
      "signed_by": [
        "Developer ID Application: LogMeIn, Inc. (GFNFVT632V)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.0.544"
    },
    {
      "_name": "IDLE 3",
      "arch_kind": "arch_other",
      "info": "3.11.4, © 2001-2023 Python Software Foundation",
      "lastModified": "2023-07-13T13:33:25Z",
      "obtained_from": "unknown",
      "path": "/usr/local/Cellar/python@3.11/3.11.4_1/IDLE 3.app",
      "version": "3.11.4"
    },
    {
      "_name": "Python Launcher 3",
      "arch_kind": "arch_i64",
      "info": "3.11.4, © 2001-2023 Python Software Foundation",
      "lastModified": "2023-07-13T13:33:24Z",
      "obtained_from": "unknown",
      "path": "/usr/local/Cellar/python@3.11/3.11.4_1/Python Launcher 3.app",
      "version": "3.11.4"
    },
    {
      "_name": "Python",
      "arch_kind": "arch_i64",
      "info": "3.12.5, (c) 2001-2023 Python Software Foundation.",
      "lastModified": "2024-08-06T19:08:49Z",
      "obtained_from": "unknown",
      "path": "/usr/local/Cellar/python@3.12/3.12.5/Frameworks/Python.framework/Versions/3.12/Resources/Python.app",
      "version": "3.12.5"
    },
    {
      "_name": "group.is.workflow.my",
      "arch_kind": "arch_other",
      "lastModified": "2024-08-12T13:38:23Z",
      "obtained_from": "unknown",
      "path": "/Users/ibering/Library/Application Scripts/group.is.workflow.my.app"
    },
    {
      "_name": "Pages",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-11-12T07:04:51Z",
      "obtained_from": "mac_app_store",
      "path": "/Applications/Pages.app",
      "signed_by": [
        "Apple Mac OS Application Signing",
        "Apple Worldwide Developer Relations Certification Authority",
        "Apple Root CA"
      ],
      "version": "14.4"
    },
    {
      "_name": "Numbers",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-11-12T07:05:39Z",
      "obtained_from": "mac_app_store",
      "path": "/Applications/Numbers.app",
      "signed_by": [
        "Apple Mac OS Application Signing",
        "Apple Worldwide Developer Relations Certification Authority",
        "Apple Root CA"
      ],
      "version": "14.4"
    },
    {
      "_name": "Keynote",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-11-12T07:07:31Z",
      "obtained_from": "mac_app_store",
      "path": "/Applications/Keynote.app",
      "signed_by": [
        "Apple Mac OS Application Signing",
        "Apple Worldwide Developer Relations Certification Authority",
        "Apple Root CA"
      ],
      "version": "14.4"
    },
    {
      "_name": "IDLE 3",
      "arch_kind": "arch_other",
      "info": "3.12.5, © 2001-2023 Python Software Foundation",
      "lastModified": "2024-08-06T19:08:49Z",
      "obtained_from": "unknown",
      "path": "/usr/local/Cellar/python@3.12/3.12.5/IDLE 3.app",
      "version": "3.12.5"
    },
    {
      "_name": "Python Launcher 3",
      "arch_kind": "arch_i64",
      "info": "3.12.5, © 2001-2023 Python Software Foundation",
      "lastModified": "2024-08-06T19:08:49Z",
      "obtained_from": "unknown",
      "path": "/usr/local/Cellar/python@3.12/3.12.5/Python Launcher 3.app",
      "version": "3.12.5"
    },
    {
      "_name": "NativeMessagingHost",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-04-22T18:02:25Z",
      "obtained_from": "identified_developer",
      "path": "/Library/Application Support/Adobe/WebExtnUtils/DC/NativeMessagingHost.app",
      "signed_by": [
        "Developer ID Application: Adobe Inc. (JQ525L2MZD)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "5.0"
    },
    {
      "_name": "Acrobat Update Helper",
      "arch_kind": "arch_arm_i64",
      "info": "1 . 2 . 6, ©2009-2015 Adobe Systems Incorporated. All rights reserved.",
      "lastModified": "2024-12-04T07:24:40Z",
      "obtained_from": "identified_developer",
      "path": "/Library/Application Support/Adobe/ARMDC/Application/Acrobat Update Helper.app",
      "signed_by": [
        "Developer ID Application: Adobe Inc. (JQ525L2MZD)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "1 . 2 . 6"
    },
    {
      "_name": "Adobe Acrobat Updater",
      "arch_kind": "arch_other",
      "info": "1 . 2 . 6, ©2009-2015 Adobe Systems Incorporated. All rights reserved.",
      "lastModified": "2024-12-04T07:31:59Z",
      "obtained_from": "unknown",
      "path": "/Library/Application Support/Adobe/ARMDC/Application/Adobe Acrobat Updater.app",
      "version": "1 . 2 . 6"
    },
    {
      "_name": "https+++xgkft.csb",
      "arch_kind": "arch_other",
      "lastModified": "2021-09-02T12:26:27Z",
      "obtained_from": "unknown",
      "path": "/Users/ibering/Library/Application Support/Firefox/Profiles/hlq66utr.default-release/storage/archives/0/2023-08-15/default/https+++xgkft.csb.app"
    },
    {
      "_name": "https+++zzun",
      "arch_kind": "arch_other",
      "lastModified": "2022-05-17T13:13:20Z",
      "obtained_from": "unknown",
      "path": "/Users/ibering/Library/Application Support/Firefox/Profiles/hlq66utr.default-release/storage/archives/0/2023-08-15/default/https+++zzun.app"
    },
    {
      "_name": "Brother Status Monitor",
      "arch_kind": "arch_i64",
      "lastModified": "2021-06-21T12:46:12Z",
      "obtained_from": "identified_developer",
      "path": "/Library/Printers/Brother/Utilities/BrStatusMonitor.app",
      "signed_by": [
        "Developer ID Application: Brother Industries, LTD.",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "3.23.0"
    },
    {
      "_name": "Brother Scanner",
      "arch_kind": "arch_i64",
      "info": "2.9.0, © 2007-2016 Brother Industries, Ltd. All Rights Reserved.",
      "lastModified": "2021-06-21T12:46:12Z",
      "obtained_from": "identified_developer",
      "path": "/Library/Image Capture/Devices/Brother Scanner.app",
      "signed_by": [
        "Developer ID Application: Brother Industries, LTD.",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "2.9.0"
    },
    {
      "_name": "Canon IJScanner2",
      "arch_kind": "arch_i64",
      "info": "Canon IJScanner2 version 4.0.0, Copyright CANON INC. 2009-2014",
      "lastModified": "2014-06-26T03:00:00Z",
      "obtained_from": "identified_developer",
      "path": "/Library/Image Capture/Devices/Canon IJScanner2.app",
      "signed_by": [
        "Developer ID Application: Canon Inc. (XE2XNRRXZ5)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "4.0.0"
    },
    {
      "_name": "EPFaxAutoSetupTool",
      "arch_kind": "arch_i64",
      "info": "Copyright(C) Seiko Epson Corporation 2009-2015. All rights reserved.",
      "lastModified": "2015-08-19T15:00:00Z",
      "obtained_from": "identified_developer",
      "path": "/Library/Printers/EPSON/Fax/AutoSetupTool/EPFaxAutoSetupTool.app",
      "signed_by": [
        "Developer ID Application: EPSON (TXAEAV5RN4)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.71"
    },
    {
      "_name": "epsonfax",
      "arch_kind": "arch_i64",
      "info": "Copyright(C) Seiko Epson Corporation 2009-2015. All rights reserved.",
      "lastModified": "2015-08-19T15:00:00Z",
      "obtained_from": "identified_developer",
      "path": "/Library/Printers/EPSON/Fax/FaxIOSupport/epsonfax.app",
      "signed_by": [
        "Developer ID Application: EPSON (TXAEAV5RN4)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.71"
    },
    {
      "_name": "commandFilter",
      "arch_kind": "arch_i64",
      "info": "Copyright(C) Seiko Epson Corporation 2012-2015. All rights reserved.",
      "lastModified": "2015-08-19T15:00:00Z",
      "obtained_from": "identified_developer",
      "path": "/Library/Printers/EPSON/Fax/Filter/commandFilter.app",
      "signed_by": [
        "Developer ID Application: EPSON (TXAEAV5RN4)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.71"
    },
    {
      "_name": "rastertoepfax",
      "arch_kind": "arch_i64",
      "info": "Copyright(C) Seiko Epson Corporation 2009-2015. All rights reserved.",
      "lastModified": "2015-08-19T15:00:00Z",
      "obtained_from": "identified_developer",
      "path": "/Library/Printers/EPSON/Fax/Filter/rastertoepfax.app",
      "signed_by": [
        "Developer ID Application: EPSON (TXAEAV5RN4)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.71"
    },
    {
      "_name": "FAX Utility",
      "arch_kind": "arch_i64",
      "info": "1.73, Copyright(C) Seiko Epson Corporation 2009-2015. All rights reserved.",
      "lastModified": "2015-08-19T15:00:00Z",
      "obtained_from": "identified_developer",
      "path": "/Library/Printers/EPSON/Fax/Utility/FAX Utility.app",
      "signed_by": [
        "Developer ID Application: EPSON (TXAEAV5RN4)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.73"
    },
    {
      "_name": "Monitor zum Faxempfang",
      "arch_kind": "arch_i64",
      "info": "1.71, Copyright(C) Seiko Epson Corporation 2009-2015. All rights reserved.",
      "lastModified": "2015-08-19T15:00:00Z",
      "obtained_from": "identified_developer",
      "path": "/Library/Printers/EPSON/Fax/Utility/Fax Receive Monitor.app",
      "signed_by": [
        "Developer ID Application: EPSON (TXAEAV5RN4)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.71"
    },
    {
      "_name": "IDLE 3",
      "arch_kind": "arch_other",
      "info": "3.10.14, © 2001-2023 Python Software Foundation",
      "lastModified": "2024-03-19T21:46:16Z",
      "obtained_from": "unknown",
      "path": "/usr/local/Cellar/python@3.10/3.10.14_1/IDLE 3.app",
      "version": "3.10.14"
    },
    {
      "_name": "Python Launcher 3",
      "arch_kind": "arch_i64",
      "info": "3.10.14, © 2001-2023 Python Software Foundation",
      "lastModified": "2024-03-19T21:46:16Z",
      "obtained_from": "unknown",
      "path": "/usr/local/Cellar/python@3.10/3.10.14_1/Python Launcher 3.app",
      "version": "3.10.14"
    },
    {
      "_name": "Python",
      "arch_kind": "arch_i64",
      "info": "3.11.4, (c) 2001-2023 Python Software Foundation.",
      "lastModified": "2023-07-13T13:33:18Z",
      "obtained_from": "unknown",
      "path": "/usr/local/Cellar/python@3.11/3.11.4_1/Frameworks/Python.framework/Versions/3.11/Resources/Python.app",
      "version": "3.11.4"
    },
    {
      "_name": "Canon IJScanner4",
      "arch_kind": "arch_i64",
      "info": "Canon IJScanner4 version 4.0.0, Copyright CANON INC. 2009-2014",
      "lastModified": "2014-06-26T03:00:00Z",
      "obtained_from": "identified_developer",
      "path": "/Library/Image Capture/Devices/Canon IJScanner4.app",
      "signed_by": [
        "Developer ID Application: Canon Inc. (XE2XNRRXZ5)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "4.0.0"
    },
    {
      "_name": "NewApp",
      "arch_kind": "arch_other",
      "lastModified": "2023-05-30T11:57:21Z",
      "obtained_from": "unknown",
      "path": "/usr/local/Homebrew/Library/Homebrew/test/support/fixtures/cask/NewApp.app"
    },
    {
      "_name": "https+++appcode",
      "arch_kind": "arch_other",
      "lastModified": "2022-06-24T12:22:34Z",
      "obtained_from": "unknown",
      "path": "/Users/ibering/Library/Application Support/Firefox/Profiles/hlq66utr.default-release/storage/archives/0/2023-08-15/default/https+++appcode.app"
    },
    {
      "_name": "JetBrains Toolbox",
      "arch_kind": "arch_i64",
      "info": "",
      "lastModified": "2023-11-10T13:22:43Z",
      "obtained_from": "identified_developer",
      "path": "/Users/ibering/Library/Application Support/JetBrains/Toolbox/apps/Toolbox/self/2.1.0.18144/JetBrains Toolbox.app",
      "signed_by": [
        "Developer ID Application: JetBrains s.r.o. (2ZEFAR8TH3)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "2.1.0.18144"
    },
    {
      "_name": "Python",
      "arch_kind": "arch_i64",
      "info": "3.10.14, (c) 2001-2023 Python Software Foundation.",
      "lastModified": "2024-03-19T21:46:16Z",
      "obtained_from": "unknown",
      "path": "/usr/local/Cellar/python@3.10/3.10.14_1/Frameworks/Python.framework/Versions/3.10/Resources/Python.app",
      "version": "3.10.14"
    },
    {
      "_name": "Google Drive",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-09-15T06:11:50Z",
      "obtained_from": "unknown",
      "path": "/Users/ibering/Applications/Chrome Apps.localized/Google Drive.app",
      "version": ""
    },
    {
      "_name": "Tabellen",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-09-15T06:11:51Z",
      "obtained_from": "unknown",
      "path": "/Users/ibering/Applications/Chrome Apps.localized/Tabellen.app",
      "version": ""
    },
    {
      "_name": "Gmail",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-09-15T06:11:51Z",
      "obtained_from": "unknown",
      "path": "/Users/ibering/Applications/Chrome Apps.localized/Gmail.app",
      "version": ""
    },
    {
      "_name": "Präsentationen",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-09-15T06:11:53Z",
      "obtained_from": "unknown",
      "path": "/Users/ibering/Applications/Chrome Apps.localized/Präsentationen.app",
      "version": ""
    },
    {
      "_name": "Dokumente",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-09-15T06:11:54Z",
      "obtained_from": "unknown",
      "path": "/Users/ibering/Applications/Chrome Apps.localized/Dokumente.app",
      "version": ""
    },
    {
      "_name": "Mockoon",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-06-30T14:43:28Z",
      "obtained_from": "identified_developer",
      "path": "/Applications/Mockoon.app",
      "signed_by": [
        "Developer ID Application: 1kB SARL-S (8443RQQKK6)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "9.3.0"
    },
    {
      "_name": "Read CHM",
      "arch_kind": "arch_i64",
      "lastModified": "2025-09-24T06:11:00Z",
      "obtained_from": "mac_app_store",
      "path": "/Applications/Read CHM.app",
      "signed_by": [
        "Apple Mac OS Application Signing",
        "Apple Worldwide Developer Relations Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.6"
    },
    {
      "_name": "Spotify",
      "arch_kind": "arch_arm",
      "lastModified": "2025-10-23T08:08:13Z",
      "obtained_from": "identified_developer",
      "path": "/Users/ibering/Library/Application Support/Spotify/PersistentCache/Update/temp/Spotify.app",
      "signed_by": [
        "Developer ID Application: Spotify (2FNC3A47ZF)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.2.74.477"
    },
    {
      "_name": "Adobe Acrobat",
      "arch_kind": "arch_arm_i64",
      "info": "Adobe Acrobat X 25.001.20997, ©1984 -2025 Adobe Systems Incorporated. All rights reserved.",
      "lastModified": "2025-12-15T07:21:27Z",
      "obtained_from": "identified_developer",
      "path": "/Applications/Adobe Acrobat DC/Adobe Acrobat.app",
      "signed_by": [
        "Developer ID Application: Adobe Inc. (JQ525L2MZD)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "25.001.20997"
    },
    {
      "_name": "Xcode",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-12-15T07:41:44Z",
      "obtained_from": "mac_app_store",
      "path": "/Applications/Xcode.app",
      "signed_by": [
        "Apple Mac OS Application Signing",
        "Apple Worldwide Developer Relations Certification Authority",
        "Apple Root CA"
      ],
      "version": "26.2"
    },
    {
      "_name": "Simulator",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-12-15T07:41:44Z",
      "obtained_from": "apple",
      "path": "/Applications/Xcode.app/Contents/Developer/Applications/Simulator.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "26.0"
    },
    {
      "_name": "Create ML",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-12-15T07:41:44Z",
      "obtained_from": "apple",
      "path": "/Applications/Xcode.app/Contents/Applications/Create ML.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "6.2"
    },
    {
      "_name": "Instruments",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-12-15T07:41:44Z",
      "obtained_from": "apple",
      "path": "/Applications/Xcode.app/Contents/Applications/Instruments.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "26.2"
    },
    {
      "_name": "FileMerge",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-12-15T07:41:44Z",
      "obtained_from": "apple",
      "path": "/Applications/Xcode.app/Contents/Applications/FileMerge.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "2.11"
    },
    {
      "_name": "Icon Composer",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-12-15T07:41:44Z",
      "obtained_from": "apple",
      "path": "/Applications/Xcode.app/Contents/Applications/Icon Composer.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "1.2"
    },
    {
      "_name": "Reality Composer Pro",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-12-15T07:41:44Z",
      "obtained_from": "apple",
      "path": "/Applications/Xcode.app/Contents/Applications/Reality Composer Pro.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "2.0"
    },
    {
      "_name": "Microsoft Edge",
      "arch_kind": "arch_arm_i64",
      "info": "Microsoft Edge 143.0.3650.96, © 2025 Microsoft Corporation. All rights reserved.",
      "lastModified": "2025-12-18T09:45:43Z",
      "obtained_from": "identified_developer",
      "path": "/Applications/Microsoft Edge.app",
      "signed_by": [
        "Developer ID Application: Microsoft Corporation (UBF8T346G9)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "143.0.3650.96"
    },
    {
      "_name": "openITCOCKPIT-Desktop",
      "arch_kind": "arch_ios",
      "lastModified": "2025-11-04T10:16:25Z",
      "obtained_from": "identified_developer",
      "path": "/Applications/openITCOCKPIT-Desktop.app",
      "signed_by": [
        "Developer ID Application: Allgeier IT Services GmbH (QL3HZGCQ4U)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "5.2.0"
    },
    {
      "_name": "Affinity Designer 2",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-11-05T07:20:11Z",
      "obtained_from": "identified_developer",
      "path": "/Applications/Affinity Designer 2.app",
      "signed_by": [
        "Developer ID Application: Serif (Europe) Ltd. (6LVTQB9699)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "2.6.5"
    },
    {
      "_name": "Python",
      "arch_kind": "arch_arm_i64",
      "info": "3.9.6, (c) 2001-2020 Python Software Foundation.",
      "lastModified": "2025-11-12T07:33:33Z",
      "obtained_from": "apple",
      "path": "/Library/Developer/CommandLineTools/Library/Frameworks/Python3.framework/Versions/3.9/Resources/Python.app",
      "signed_by": [
        "Software Signing",
        "Apple Code Signing Certification Authority",
        "Apple Root CA"
      ],
      "version": "3.9.6"
    },
    {
      "_name": "Microsoft Teams",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2025-11-21T07:17:35Z",
      "obtained_from": "identified_developer",
      "path": "/Applications/Microsoft Teams.app",
      "signed_by": [
        "Developer ID Application: Microsoft Corporation (UBF8T346G9)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "25306.805.4102.7211"
    },
    {
      "_name": "Citrix Workspace Updater",
      "arch_kind": "arch_arm_i64",
      "info": "25.08.10",
      "lastModified": "2025-11-07T12:10:22Z",
      "obtained_from": "identified_developer",
      "path": "/Library/Application Support/Citrix Workspace Updater/Citrix Workspace Updater.app",
      "signed_by": [
        "Developer ID Application: Citrix Systems, Inc. (S272Y5R93J)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "25.08.10"
    },
    {
      "_name": "https+++itf-full-stack-essentials.netlify",
      "arch_kind": "arch_other",
      "lastModified": "2026-01-09T06:17:10Z",
      "obtained_from": "unknown",
      "path": "/Users/ibering/Library/Application Support/Firefox/Profiles/hlq66utr.default-release/storage/default/https+++itf-full-stack-essentials.netlify.app"
    },
    {
      "_name": "Microsoft OneNote",
      "arch_kind": "arch_arm_i64",
      "info": "SZLONGVERSION",
      "lastModified": "2026-01-07T07:42:10Z",
      "obtained_from": "identified_developer",
      "path": "/Applications/Microsoft OneNote.app",
      "signed_by": [
        "Developer ID Application: Microsoft Corporation (UBF8T346G9)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "16.104.1"
    },
    {
      "_name": "Microsoft PowerPoint",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2026-01-07T07:42:50Z",
      "obtained_from": "identified_developer",
      "path": "/Applications/Microsoft PowerPoint.app",
      "signed_by": [
        "Developer ID Application: Microsoft Corporation (UBF8T346G9)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "16.104.1"
    },
    {
      "_name": "Microsoft Excel",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2026-01-07T07:43:43Z",
      "obtained_from": "identified_developer",
      "path": "/Applications/Microsoft Excel.app",
      "signed_by": [
        "Developer ID Application: Microsoft Corporation (UBF8T346G9)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "16.104.1"
    },
    {
      "_name": "Microsoft Outlook",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2026-01-07T07:44:47Z",
      "obtained_from": "identified_developer",
      "path": "/Applications/Microsoft Outlook.app",
      "signed_by": [
        "Developer ID Application: Microsoft Corporation (UBF8T346G9)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "16.104.1"
    },
    {
      "_name": "Microsoft Word",
      "arch_kind": "arch_arm_i64",
      "lastModified": "2026-01-07T07:44:15Z",
      "obtained_from": "identified_developer",
      "path": "/Applications/Microsoft Word.app",
      "signed_by": [
        "Developer ID Application: Microsoft Corporation (UBF8T346G9)",
        "Developer ID Certification Authority",
        "Apple Root CA"
      ],
      "version": "16.104.1"
    }
  ]
}`

func TestParseMacOSSoftwareUpdateOutput(t *testing.T) {
	updates, err := parseMacOSSoftwareUpdateOutput(softwareUpdateOutputSample)
	if err != nil {
		t.Fatalf("Error parsing software update output: %v", err)
	}

	if len(updates) != 2 {
		t.Fatalf("Expected 2 updates, got %d", len(updates))
	}

	expectedFirst := MacosUpdate{
		Name:        "Command Line Tools for Xcode 26.2-26.2",
		Description: "Command Line Tools for Xcode 26.2",
		Version:     "26.2",
	}

	if updates[0] != expectedFirst {
		t.Errorf("First update does not match expected.\nGot: %+v\nExpected: %+v", updates[0], expectedFirst)
	}

	expectedSecond := MacosUpdate{
		Name:        "macOS Tahoe 26.2-25C56",
		Description: "macOS Tahoe 26.2",
		Version:     "26.2",
	}

	if updates[1] != expectedSecond {
		t.Errorf("Second update does not match expected.\nGot: %+v\nExpected: %+v", updates[1], expectedSecond)
	}
}

func TestParseMacOSInstalledAppsOutput_ValidInput(t *testing.T) {
	input := `{
		"SPApplicationsDataType": [
			{
				"_name": "Safari",
				"version": "16.4",
				"info": "Apple web browser"
			},
			{
				"_name": "Xcode",
				"version": "14.3",
				"info": "Apple IDE"
			}
		]
	}`

	expected := []Package{
		{Name: "Safari", Version: "16.4", Description: "Apple web browser"},
		{Name: "Xcode", Version: "14.3", Description: "Apple IDE"},
	}

	result, err := parseMacOSInstalledAppsOutput(input)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestParseMacOSInstalledAppsOutput_RealWorld(t *testing.T) {

	expected := []Package{
		{Name: "openITCOCKPIT-Desktop", Version: "5.2.0", Description: ""},
		{Name: "Microsoft Edge", Version: "143.0.3650.96", Description: "Microsoft Edge 143.0.3650.96, © 2025 Microsoft Corporation. All rights reserved."},
		{Name: "Xcode", Version: "26.2", Description: ""},
		{Name: "Python", Version: "3.9.6", Description: "3.9.6, (c) 2001-2020 Python Software Foundation."},
	}

	result, err := parseMacOSInstalledAppsOutput(systemProfilerOutputSample)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check that the given examples are within the result
	foundCount := 0
	for _, expPkg := range expected {
		for _, resPkg := range result {
			if expPkg == resPkg {
				foundCount++
				break
			}
		}
	}

	if foundCount != len(expected) {
		t.Errorf("Expected to find all %d expected packages, but found %d", len(expected), foundCount)
	}
}

func TestParseMacOSInstalledAppsOutput_EmptyList(t *testing.T) {
	input := `{"SPApplicationsDataType": []}`

	result, err := parseMacOSInstalledAppsOutput(input)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(result) != 0 {
		t.Errorf("Expected empty list, got %v", result)
	}
}

func TestParseMacOSInstalledAppsOutput_MissingFields(t *testing.T) {
	input := `{
		"SPApplicationsDataType": [
			{
				"_name": "Terminal"
			}
		]
	}`

	expected := []Package{
		{Name: "Terminal", Version: "", Description: ""},
	}

	result, err := parseMacOSInstalledAppsOutput(input)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestParseMacOSInstalledAppsOutput_InvalidJSON(t *testing.T) {
	input := `{"SPApplicationsDataType": [ { "_name": "Safari", "version": "16.4", "info": "Apple web browser" }`

	_, err := parseMacOSInstalledAppsOutput(input)
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

func TestParseMacOSInstalledAppsOutput_ExtraFieldsIgnored(t *testing.T) {
	input := `{
		"SPApplicationsDataType": [
			{
				"_name": "Pages",
				"version": "12.1",
				"info": "Word processor",
				"extra_field": "should be ignored"
			}
		]
	}`

	expected := []Package{
		{Name: "Pages", Version: "12.1", Description: "Word processor"},
	}

	result, err := parseMacOSInstalledAppsOutput(input)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}
