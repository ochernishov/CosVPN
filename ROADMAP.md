# CosVPN -- Дорожная карта ребрендинга и сборки клиентов

> Дата создания: 2026-03-20
> Сервер: 72.56.26.254:443 (Ubuntu 22.04, CosVPN)
> Тестовый клиент: test-client1 (10.0.0.2)
> Монорепо: https://github.com/ochernishov/CosVPN

---

## Этап 1: CosVPN-Android -- Ребрендинг и сборка

> Приоритет: ВЫСШИЙ
> Папка: CosVPN-Android/
> Стек: Kotlin/Java, Gradle, Android SDK 36, minSdk 24
> Текущий package: com.wireguard.android
> Целевой package: com.cosinn.vpn.android (или аналогичный)

### Блок 1.1 — Подготовка окружения сборки

- Установить Android SDK, NDK, CMake (требуется для нативной библиотеки libwg-go.so)
- Проверить наличие JDK 17+ (требование gradle)
- Проверить сборку ОРИГИНАЛЬНОГО проекта без изменений (`./gradlew assembleDebug`)
- Убедиться, что нативные библиотеки (tunnel/tools/CMakeLists.txt) собираются корректно
- Зафиксировать baseline: проект собирается, APK устанавливается на устройство/эмулятор

### Блок 1.2 — Ребрендинг: идентификаторы и метаданные

- Изменить `wireguardPackageName` в `gradle.properties` (com.wireguard.android -> com.cosinn.vpn.android)
- Изменить `rootProject.name` в `settings.gradle.kts` (wireguard-android -> cosvpn-android)
- Обновить `wireguardVersionName` и `wireguardVersionCode` в `gradle.properties`
- Обновить POM-метаданные в `tunnel/build.gradle.kts`: name, description, url, scm, developers
- Обновить copyright-заголовки в файлах (WireGuard LLC -> CosinnDev)
- Обновить `app_name` в `ui/src/main/res/values/strings.xml` (WireGuard -> CosVPN)
- Обновить русскую локализацию `ui/src/main/res/values-ru/strings.xml`
- Переместить Java/Kotlin-пакеты из com.wireguard.android в com.cosinn.vpn.android
- Обновить все import-ссылки во всех .kt/.java файлах
- Обновить intent-actions в AndroidManifest.xml (com.wireguard.android.action.* -> com.cosinn.vpn.android.action.*)

### Блок 1.3 — Ребрендинг: иконки и визуал

- Создать иконку CosVPN (ic_launcher) во всех разрешениях: mdpi, hdpi, xhdpi, xxhdpi, xxxhdpi
- Создать круглую иконку (ic_launcher_round) во всех разрешениях
- Создать баннер для Android TV (mipmap-xhdpi/banner.png)
- Обновить ic_launcher_foreground.xml (adaptive icon foreground)
- Обновить ic_tile.xml (Quick Settings tile icon)
- Обновить цветовую схему в styles.xml / themes при необходимости
- Заменить текст "WireGuard" в любых drawable-ресурсах

### Блок 1.4 — Подключение к серверу и тестирование

- Создать нового клиента на сервере через add-client.sh для Android-устройства
- Собрать debug APK (`./gradlew assembleDebug`)
- Установить на эмулятор Android (API 24+) и проверить запуск
- Импортировать конфигурацию (ручной ввод или QR-код)
- Проверить подключение к серверу 72.56.26.254:51820
- Проверить туннелирование трафика (ping, curl через VPN)
- Проверить Quick Settings tile
- Проверить auto-connect при загрузке устройства
- Установить на реальное Android-устройство и повторить тесты

### Блок 1.5 — Финализация Android

- Собрать release APK с подписью (`./gradlew assembleRelease`)
- Настроить keystore для подписи (создать или использовать существующий)
- Проверить ProGuard/R8 минификацию (proguard-android-optimize.txt)
- Протестировать release-сборку на устройстве
- Задокументировать процесс сборки в README.md внутри CosVPN-Android/

---

## Этап 2: CosVPN-Apple (macOS) -- Ребрендинг и сборка

> Приоритет: ВЫСОКИЙ
> Папка: CosVPN-Apple/
> Стек: Swift, Xcode, Network Extension framework
> Текущий bundle ID: задан через APP_ID_MACOS в Xcode build settings
> Xcode project: WireGuard.xcodeproj
> Таргеты: WireGuardmacOS, WireGuardNetworkExtension (macOS), WireGuardmacOSLoginItemHelper

### Блок 2.1 — Подготовка окружения macOS ✅

- Убедиться в наличии Xcode 15+ с Command Line Tools
- Открыть WireGuard.xcodeproj и проверить, что проект компилируется без изменений
- Настроить Apple Developer Team ID в Xcode (DEVELOPMENT_TEAM)
- Проверить, что Network Extension entitlement доступен в профиле разработчика
- Зафиксировать baseline: проект собирается для macOS таргета

### Блок 2.2 — Ребрендинг: идентификаторы macOS ✅

- Изменить PRODUCT_NAME для macOS таргета (WireGuard -> CosVPN) в pbxproj
- Изменить APP_ID_MACOS build setting (com.wireguard.macos -> com.cosinn.vpn.macos)
- Обновить PRODUCT_BUNDLE_IDENTIFIER для всех macOS таргетов:
  - Приложение: $(APP_ID_MACOS)
  - Network Extension: $(APP_ID_MACOS).network-extension
  - Login Item Helper: $(APP_ID_MACOS).login-item-helper
- Обновить app group: $(DEVELOPMENT_TEAM).group.$(APP_ID_MACOS)
- Обновить copyright в Info.plist (WireGuard LLC -> CosinnDev)
- Обновить NSPrincipalClass в Info.plist если содержит WireGuard в пути
- Обновить Localizable.strings для macOS (если есть упоминания WireGuard в UI-строках)
- Переименовать WireGuard.xcodeproj -> CosVPN.xcodeproj (опционально, может сломать ссылки)

### Блок 2.3 — Ребрендинг: иконки macOS ✅

- Создать AppIcon для macOS (Assets.xcassets/AppIcon.appiconset) -- все размеры (16-1024px)
- Обновить StatusBar иконки (StatusBarIcon, StatusBarIconDimmed, StatusBarIconDot1-3)
- Обновить StatusCircleYellow если используется в UI
- Создать иконку в стиле CosVPN (согласовать с Android-иконкой)

### Блок 2.4 — Подключение к серверу и тестирование macOS ✅

- Создать нового клиента на сервере для macOS через add-client.sh
- Собрать приложение через Xcode (Product -> Build)
- Запустить на локальном маке
- Импортировать конфигурацию (.conf файл или ручной ввод)
- Проверить подключение к серверу 72.56.26.254:51820
- Проверить туннелирование трафика
- Проверить status bar menu (StatusItemController)
- Проверить auto-connect и login item helper
- Проверить логи через встроенный log viewer

### Блок 2.5 — Финализация macOS ✅

- Настроить code signing для distribution
- Собрать .app bundle для распространения
- Проверить notarization (если нужна установка вне App Store)
- Задокументировать процесс сборки

---

## Этап 3: CosVPN-Apple (iOS) -- Ребрендинг и сборка

> Приоритет: ВЫСОКИЙ
> Папка: CosVPN-Apple/ (тот же проект, iOS-таргет)
> Стек: Swift, Xcode, NetworkExtension framework
> Текущий bundle ID: задан через APP_ID_IOS в Xcode build settings
> Таргеты: WireGuard (iOS), WireGuardNetworkExtension (iOS)

### Блок 3.1 — Подготовка окружения iOS

- Убедиться, что iOS Simulator установлен в Xcode
- Проверить, что iOS-таргет компилируется без изменений
- Настроить provisioning profile для iOS (если тестирование на реальном устройстве)
- Зафиксировать baseline: проект собирается для iOS таргета

### Блок 3.2 — Ребрендинг: идентификаторы iOS

- Изменить APP_ID_IOS build setting (com.wireguard.ios -> com.cosinn.vpn.ios)
- Обновить PRODUCT_BUNDLE_IDENTIFIER для iOS таргетов:
  - Приложение: $(APP_ID_IOS)
  - Network Extension: $(APP_ID_IOS).network-extension
- Обновить app group: group.$(APP_ID_IOS)
- Обновить PRODUCT_NAME для iOS таргета (WireGuard -> CosVPN)
- Обновить CFBundleDisplayName и CFBundleName в Info.plist
- Обновить UTTypeIdentifier в Info.plist (com.wireguard.config.quick -> com.cosinn.vpn.config.quick)
- Обновить CFBundleTypeIconFiles (заменить wireguard_doc_logo_*.png)
- Обновить UTTypeDescription и CFBundleTypeName
- Обновить Localizable.strings для iOS

### Блок 3.3 — Ребрендинг: иконки iOS

- Создать AppIcon для iOS (Assets.xcassets/AppIcon.appiconset) -- все размеры
- Создать document icons (wireguard_doc_logo -> cosvpn_doc_logo) 22x29, 44x58, 64x64, 320x320
- Обновить LaunchScreen storyboard если содержит логотип WireGuard
- Обновить wireguard.imageset в iOS Assets

### Блок 3.4 — Подключение к серверу и тестирование iOS

- Создать нового клиента на сервере для iOS через add-client.sh
- Собрать и запустить на iOS Simulator
- Импортировать конфигурацию (ручной ввод, файл .conf, или QR)
- Проверить подключение к серверу 72.56.26.254:51820
- Проверить туннелирование трафика
- Проверить on-demand правила (WiFi/Cellular)
- Проверить Face ID / biometric protection
- Установить на реальное iOS-устройство (через Xcode) и повторить тесты

### Блок 3.5 — Финализация iOS

- Настроить code signing для distribution (Ad Hoc или App Store)
- Собрать .ipa для распространения
- Протестировать на разных устройствах (iPhone, iPad)
- Задокументировать процесс сборки

---

## Этап 4: Общие задачи и инфраструктура

> Приоритет: СРЕДНИЙ
> Выполняется параллельно или после этапов 1-3

### Блок 4.1 — Серверная инфраструктура

- Автоматизировать add-client.sh: параметризация имени клиента, IP, DNS
- Создать скрипт генерации QR-кодов для мобильных клиентов
- Настроить мониторинг WireGuard-сервера (wg show, подключённые пиры)
- Документировать серверную конфигурацию

### Блок 4.2 — CI/CD и автоматизация сборок

- Настроить GitHub Actions для сборки Android APK
- Настроить GitHub Actions для сборки macOS .app (на macOS runner)
- Автоматическое версионирование (tag-based)
- Артефакты сборок в GitHub Releases

### Блок 4.3 — Документация и README

- Обновить корневой README.md монорепо с описанием всех компонентов
- Создать CONTRIBUTING.md с инструкциями по сборке каждой платформы
- Создать docs/ с архитектурной документацией
- Описать процесс добавления нового клиента (end-to-end)

### Блок 4.4 — Дизайн и брендинг

- Утвердить цветовую палитру CosVPN
- Создать логотип CosVPN в векторном формате (SVG)
- Экспортировать иконки для всех платформ из единого источника
- Создать splash screen / launch screen дизайн

---

## ✅ Этап 5: CosVPN Protocol — Собственный протокол с обфускацией

> Приоритет: КРИТИЧЕСКИЙ
> Папка: CosVPN-Go/
> Стек: Go 1.23+, crypto/tls, ChaCha20-Poly1305
> Статус: ЗАВЕРШЁН

### ✅ Блок 5.1 — Ребрендинг протокола

- Все идентификаторы WireGuard → CosVPN (96 файлов)
- Протокольные константы: CosVPNIdentifier, CosVPNLabelMAC1, CosVPNLabelCookie
- Module path: github.com/ochernishov/cosvpn
- Copyright: CosinnDev

### ✅ Блок 5.2 — Пакет обфускации (obfs/)

- XOR header obfuscation (первые 16 байт)
- Random padding (0-64 байт)
- Junk packet injection (30% вероятность)
- 5 тестов — все PASS

### ✅ Блок 5.3 — Интеграция в send/receive

- ObfsConfig в Device struct
- Обфускация перед UDP-отправкой (send.go)
- Деобфускация при получении (receive.go)
- UAPI: obfuscation_key, obfuscation_mode

### ✅ Блок 5.4 — TLS 1.3 обёртка

- TLSListener + TLSClient с самоподписанным ECDSA P-256
- Framing: [2 bytes BE length][payload]
- Минимальная версия TLS 1.3

### ✅ Блок 5.5 — Деплой на сервер

- Бинарник cosvpn-go задеплоен на 72.56.26.254
- ObfuscationKey сгенерирован и сохранён
- Порт: 443 (маскировка под HTTPS)
