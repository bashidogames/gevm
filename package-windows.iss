#define MyAppPublisher "Bashido Games"
#define MyAppURL "https://github.com/bashidogames/gdvm"
#define MyAppExeName "gdvm.exe"

[Setup]
AppId={#MyAppID}
AppName={#MyAppName}
AppVersion={#MyAppVersion}
AppVerName={#MyAppName} {#MyAppVersion}
AppPublisher={#MyAppPublisher}
AppPublisherURL={#MyAppURL}
AppSupportURL={#MyAppURL}
AppUpdatesURL={#MyAppURL}
OutputBaseFilename=gdvm-windows-{#MyAppArch}-installer
DefaultDirName={autopf}\gdvm\bin\{#MyAppArch}
DisableDirPage=yes
DisableProgramGroupPage=yes
PrivilegesRequired=lowest
OutputDir=.
Compression=lzma
SolidCompression=yes
WizardStyle=modern
DisableWelcomePage=no
ChangesEnvironment=true
SetupLogging=yes

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"

[Files]
Source: "{#MyAppExeName}"; DestDir: "{app}"; Flags: ignoreversion

[Tasks]
Name: addPath; Description: "Add to PATH"

[Code]
const UninstallKey = 'Software\Microsoft\Windows\CurrentVersion\Uninstall\{#MyAppID}_is1';
const EnvironmentKey = 'Environment';
const LogUninstallFilename = 'gdvm-{#MyAppArch}-uninstall.log';
const LogInstallFilename = 'gdvm-{#MyAppArch}-install.log';
const PathValue = 'Path';
const LogParam = '/log';

procedure AddLogParam(ValueName: string);
var
    NewValue, OldValue: string;
begin
    if not RegQueryStringValue(HKEY_CURRENT_USER, UninstallKey, ValueName, OldValue) then begin
        Log(Format('ERROR: Registry key/value not found: [HKEY_CURRENT_USER\%s] - [%s]', [UninstallKey, ValueName]));
        exit;
    end;

    NewValue := OldValue + ' ' + LogParam;

    if not RegWriteStringValue(HKEY_CURRENT_USER, UninstallKey, ValueName, NewValue) then
        Log(Format('ERROR: Registry not written: [HKEY_CURRENT_USER\%s] - [%s] = [%s] => [%s]', [UninstallKey, ValueName, OldValue, NewValue]))
    else
        Log(Format('Registry written: [HKEY_CURRENT_USER\%s] - [%s] = [%s] => [%s]', [UninstallKey, ValueName, OldValue, NewValue]));
end;

procedure RemovePath(Path: string);
var
    CountOffset, IndexOffset, Index: Integer;
    NewPaths, OldPaths: string;
begin
    if not RegQueryStringValue(HKEY_CURRENT_USER, EnvironmentKey, PathValue, OldPaths) then begin
        Log(Format('[%s] not found in PATH: [%s]', [Path, OldPaths]));
        exit;
    end;

    Index := Pos(';' + Uppercase(Path) + '\;',  ';' + Uppercase(OldPaths) + ';');
    CountOffset := 2;

    if Index = 0 then begin
        Index := Pos(';' + Uppercase(Path) + ';',  ';' + Uppercase(OldPaths) + ';');
        CountOffset := 1;
    end;

    if Index = 0 then begin
        Log(Format('[%s] not found in PATH: [%s]', [Path, OldPaths]));
        exit;
    end;

    NewPaths := OldPaths

    if Index <> 1 then
        IndexOffset := 1
    else
        IndexOffset := 0;

    Delete(NewPaths, Index - IndexOffset, Length(Path) + CountOffset);

    if not RegWriteStringValue(HKEY_CURRENT_USER, EnvironmentKey, PathValue, NewPaths) then
        Log(Format('ERROR: [%s] not removed from PATH: [%s] => [%s]', [Path, OldPaths, NewPaths]))
    else
        Log(Format('[%s] removed from PATH: [%s] => [%s]', [Path, OldPaths, NewPaths]));
end;

procedure AddPath(Path: string);
var
    NewPaths, OldPaths: string;
begin
    if not RegQueryStringValue(HKEY_CURRENT_USER, EnvironmentKey, PathValue, OldPaths) then
        OldPaths := '';

    if Pos(';' + Uppercase(Path) + '\;',  ';' + Uppercase(OldPaths) + ';') > 0 then begin
        Log(Format('[%s] already exists in PATH: [%s]', [Path, OldPaths]));
        exit;
    end;

    if Pos(';' + Uppercase(Path) + ';',  ';' + Uppercase(OldPaths) + ';') > 0 then begin
        Log(Format('[%s] already exists in PATH: [%s]', [Path, OldPaths]));
        exit;
    end;

    if OldPaths[length(OldPaths)] <> ';' then
        NewPaths := OldPaths + ';' + Path + ';'
    else
        NewPaths := OldPaths + Path + ';';

    if not RegWriteStringValue(HKEY_CURRENT_USER, EnvironmentKey, PathValue, NewPaths) then
        Log(Format('ERROR: [%s] not added to PATH: [%s] => [%s]', [Path, OldPaths, NewPaths]))
    else
        Log(Format('[%s] added to PATH: [%s] => [%s]', [Path, OldPaths, NewPaths]));
end;

procedure CopyLog(DestDir, Filename: string);
var
  SourcePath, DestPath: string;
begin
  SourcePath := ExpandConstant('{log}');
  DestPath := DestDir + '\' + Filename;
  FileCopy(SourcePath, DestPath, false);
end;

procedure CurUninstallStepChanged(CurUninstallStep: TUninstallStep);
begin
    if CurUninstallStep = usPostUninstall then
        RemovePath(ExpandConstant('{app}'));
end;

procedure CurStepChanged(CurStep: TSetupStep);
begin
    if CurStep = ssPostInstall then begin
        if WizardIsTaskSelected('addPath') then
            AddPath(ExpandConstant('{app}'));
        
        AddLogParam('QuietUninstallString');
        AddLogParam('UninstallString');
    end;
end;

procedure DeinitializeUninstall();
begin
    CopyLog(ExpandConstant(GetEnv('USERPROFILE')), LogUninstallFilename);
end;

procedure DeinitializeSetup();
begin
    CopyLog(ExpandConstant('{src}'), LogInstallFilename);
end;
