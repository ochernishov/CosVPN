# CosVPN for macOS (and iOS)

VPN-клиент для macOS, основанный на WireGuard протоколе. Ребрендинг оригинального wireguard-apple.

## Требования

- **macOS 12.0+** (deployment target)
- **Xcode 15+** с Command Line Tools
- **Go 1.19+** (для компиляции wireguard-go bridge)
- **SwiftLint** (для проверки стиля кода)

## Установка зависимостей

```bash
brew install swiftlint go
```

Или установка Go вручную:
```bash
curl -L -o /tmp/go.tar.gz "https://go.dev/dl/go1.24.1.darwin-arm64.tar.gz"
mkdir -p ~/go-sdk && tar -C ~/go-sdk -xzf /tmp/go.tar.gz
export PATH="$HOME/go-sdk/go/bin:$PATH"
```

## Настройка Developer.xcconfig

Скопируйте шаблон и заполните данные:

```bash
cp Sources/WireGuardApp/Config/Developer.xcconfig.template Sources/WireGuardApp/Config/Developer.xcconfig
```

Заполните:
```
DEVELOPMENT_TEAM = <ваш Apple Team ID>
APP_ID_IOS = com.wireguard.ios
APP_ID_MACOS = com.cosinn.vpn.macos
```

**Важно:** Для работы Network Extension необходим Apple Developer account с Network Extension capability.

## Сборка

### Через Xcode

```bash
open WireGuard.xcodeproj
```

Выберите scheme **WireGuardmacOS** и нажмите **Product -> Build** (Cmd+B).

### Через командную строку

```bash
export PATH="$HOME/go-sdk/go/bin:$PATH"  # если Go установлен вручную

# Debug сборка (без подписи, для разработки)
xcodebuild -project WireGuard.xcodeproj \
  -scheme WireGuardmacOS \
  -configuration Debug build \
  CODE_SIGN_IDENTITY="" \
  CODE_SIGNING_REQUIRED=NO \
  CODE_SIGNING_ALLOWED=NO \
  ONLY_ACTIVE_ARCH=YES

# Release сборка (с подписью, DEVELOPMENT_TEAM должен быть заполнен)
xcodebuild -project WireGuard.xcodeproj \
  -scheme WireGuardmacOS \
  -configuration Release build
```

## Ребрендинг

Проект ребрендирован с WireGuard на CosVPN:

| Параметр | Было | Стало |
|----------|------|-------|
| Product Name | WireGuard | CosVPN |
| Bundle ID | com.wireguard.macos | com.cosinn.vpn.macos |
| Network Extension | com.wireguard.macos.network-extension | com.cosinn.vpn.macos.network-extension |
| Login Item Helper | com.wireguard.macos.login-item-helper | com.cosinn.vpn.macos.login-item-helper |
| Copyright | WireGuard LLC | CosinnDev |

### Изменённые файлы

- `Sources/WireGuardApp/Config/Developer.xcconfig` -- bundle ID
- `WireGuard.xcodeproj/project.pbxproj` -- PRODUCT_NAME
- `Sources/WireGuardApp/UI/macOS/Info.plist` -- copyright, NSPrincipalClass, custom keys
- `Sources/WireGuardApp/UI/macOS/LoginItemHelper/Info.plist` -- copyright, custom keys
- `Sources/WireGuardApp/UI/macOS/LoginItemHelper/main.m` -- plist key references
- `Sources/Shared/FileManager+Extension.swift` -- plist key reference
- `Sources/WireGuardApp/Base.lproj/Localizable.strings` -- UI strings (English)
- `Sources/WireGuardApp/ru.lproj/Localizable.strings` -- UI strings (Russian)
- `Sources/WireGuardApp/UI/macOS/Assets.xcassets/AppIcon.appiconset/` -- app icons
- `Sources/WireGuardApp/UI/macOS/Assets.xcassets/StatusBarIcon*.imageset/` -- status bar icons
- `Sources/WireGuardKitC/WireGuardKitC.h` -- compatibility fix for macOS 26 SDK

## VPN-сервер для тестирования

```
Server: 72.56.26.254:51820
Config: test-configs/cosvpn-macos-test1.conf
```

## Структура таргетов

- **WireGuardmacOS** -- основное macOS приложение (CosVPN.app)
- **WireGuardNetworkExtensionmacOS** -- Network Extension (CosVPNNetworkExtension.appex)
- **WireGuardGoBridgemacOS** -- Go bridge для wireguard-go
- **WireGuardmacOSLoginItemHelper** -- Login Item Helper (CosVPNLoginItemHelper.app)

## Лицензия

MIT License. Оригинальный проект: [wireguard-apple](https://git.zx2c4.com/wireguard-apple).
