# tSF_Installer v2.2.2 (Java) / v1.0.0 (Go)
Tactical Shift Framework installer

### Как пользоваться (CLI версия)
- [Скачать .EXE](https://raw.githubusercontent.com/10Dozen/tSF_Installer/refs/heads/master/tSF_Installer_v1.0.0.zip) 
- Распаковать и положить в любое удобное место
- Запустить `cmd` / `PowerShell` и перейти в директорию куда положили .exe файл (альтернативно - запустить `run_cmd.bat`)
- В окне консоли выполнить:
```bash
# Посмотреть опции
tSFInstaller.exe install -h

# Установить в директорию path\to\mission\folder
tSFInstaller.exe install --dir="path\to\mission\folder"
```

### Как пользоваться (UI версия, устарело, но работает)
- Установить Java 8+ ([link](https://www.java.com/download/ie_manual.jsp))
- [Скачать JAR](https://github.com/10Dozen/tSF_Installer/raw/master/tSF_Installer_v2.2.2.jar)
- Положить .jar файл в любое удобное место
- Запустить .jar файл, выбрать целевую директорию (пустую или с вашей миссией) и кликнуть Install


#### Changelog
##### GO v1.0.0
- Добавилось сравнение существующих файлов в указанной директории с файлами скачанными из репозиториев. Для ряда файлов будут созданы копии и HTML-файл показывающий различия между текущими файлами миссии и файлами из свежего фреймворка. Данные проверки можно отключить флагом `--nobackup`.
