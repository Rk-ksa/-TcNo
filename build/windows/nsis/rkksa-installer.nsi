Unicode true
!include "MUI2.nsh"

!define PRODUCT_NAME "TcNo Account Switcher - Rk-ksa Edition"
!define PRODUCT_EXE "TcNo-Acc-Switcher-Rk-ksa.exe"
!define PRODUCT_VERSION "4.0.8"
!define UNINSTALL_KEY "Software\Microsoft\Windows\CurrentVersion\Uninstall\TcNo Account Switcher - Rk-ksa Edition"

Name "${PRODUCT_NAME}"
Caption "${PRODUCT_NAME} Setup"
OutFile ".tmp\release-assets\TcNo-Acc-Switcher-Rk-ksa-Setup.exe"
InstallDir "$LOCALAPPDATA\Programs\TcNo Account Switcher - Rk-ksa Edition"
InstallDirRegKey HKCU "Software\Rk-ksa\TcNo Account Switcher" "InstallDir"
RequestExecutionLevel user
ShowInstDetails show
ShowUninstDetails show

!define MUI_ABORTWARNING
!define MUI_ICON "build\windows\icon.ico"
!insertmacro MUI_PAGE_WELCOME
!insertmacro MUI_PAGE_LICENSE "LICENSE"
!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_INSTFILES
!define MUI_FINISHPAGE_RUN "$INSTDIR\${PRODUCT_EXE}"
!define MUI_FINISHPAGE_RUN_TEXT "Launch TcNo Account Switcher - Rk-ksa Edition"
!insertmacro MUI_PAGE_FINISH
!insertmacro MUI_UNPAGE_CONFIRM
!insertmacro MUI_UNPAGE_INSTFILES
!insertmacro MUI_LANGUAGE "English"

Section "Install"
  SetShellVarContext current
  SetOutPath "$INSTDIR"
  File /oname=${PRODUCT_EXE} "bin\TcNo-Acc-Switcher-Polished.exe"
  File "OPEN_SOURCE_LICENSES.txt"
  WriteUninstaller "$INSTDIR\Uninstall.exe"

  CreateDirectory "$SMPROGRAMS\TcNo Account Switcher - Rk-ksa Edition"
  CreateShortcut "$SMPROGRAMS\TcNo Account Switcher - Rk-ksa Edition\TcNo Account Switcher.lnk" "$INSTDIR\${PRODUCT_EXE}"
  CreateShortcut "$SMPROGRAMS\TcNo Account Switcher - Rk-ksa Edition\Uninstall.lnk" "$INSTDIR\Uninstall.exe"

  WriteRegStr HKCU "Software\Rk-ksa\TcNo Account Switcher" "InstallDir" "$INSTDIR"
  WriteRegStr HKCU "${UNINSTALL_KEY}" "DisplayName" "${PRODUCT_NAME}"
  WriteRegStr HKCU "${UNINSTALL_KEY}" "DisplayVersion" "${PRODUCT_VERSION}"
  WriteRegStr HKCU "${UNINSTALL_KEY}" "UninstallString" '"$INSTDIR\Uninstall.exe"'
  WriteRegStr HKCU "${UNINSTALL_KEY}" "InstallLocation" "$INSTDIR"
  WriteRegDWORD HKCU "${UNINSTALL_KEY}" "NoModify" 1
  WriteRegDWORD HKCU "${UNINSTALL_KEY}" "NoRepair" 1
SectionEnd

Section "Uninstall"
  Delete "$INSTDIR\${PRODUCT_EXE}"
  Delete "$INSTDIR\OPEN_SOURCE_LICENSES.txt"
  Delete "$INSTDIR\Uninstall.exe"
  RMDir "$INSTDIR"
  Delete "$SMPROGRAMS\TcNo Account Switcher - Rk-ksa Edition\TcNo Account Switcher.lnk"
  Delete "$SMPROGRAMS\TcNo Account Switcher - Rk-ksa Edition\Uninstall.lnk"
  RMDir "$SMPROGRAMS\TcNo Account Switcher - Rk-ksa Edition"
  DeleteRegKey HKCU "${UNINSTALL_KEY}"
  DeleteRegKey HKCU "Software\Rk-ksa\TcNo Account Switcher"
SectionEnd
