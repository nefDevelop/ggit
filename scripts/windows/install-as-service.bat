@ECHO off

:: This script relies on nssm.exe to work.
:: Please, download it and make it available on the system path,
:: or copy it to the ggit path.
:: https://nssm.cc/download
:: This script itself should run in the ggit path, too.
:: In case of startup failure, please read carefully the log file.
:: Make sure Ggit work running manually with "gogs web" before running
:: this script.
:: And, please, read carefully the installation docs first:
:: https://gogs.io/fine-tuning/run-as-service#windows
:: To uninstall the service, run "nssm remove gogs" and restart Windows.

:: Set the folder where you extracted Gogs. Omit the last slash.
SET gogspath=C:\gogs

nssm install ggit "%gogspath%\gogs.exe"
nssm set ggit AppParameters "web"
nssm set ggit Description "A painless self-hosted Git service."
nssm set ggit DisplayName "Gogs"
nssm set ggit Start SERVICE_DELAYED_AUTO_START
nssm set ggit AppStdout "%gogspath%\gogs.log"
nssm start gogs
pause
