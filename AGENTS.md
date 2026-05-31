# AGENTS.md

Pravidla pro AI coding agenty pracující v tomto projektu.

Tento soubor je určený pro Codex, Antigravity, Gemini nebo jiného AI coding agenta. Při práci v projektu se řiď těmito pravidly.

## Komunikace

- Se uživatelem vždy komunikuj česky.
- Kód, názvy souborů, technické identifikátory, příkazy, chybové hlášky, názvy knihoven a názvy API ponechávej v původním jazyce.
- Odpovídej stručně, prakticky a věcně.
- Nepředpokládej kontext z předchozího chatu, pokud není zapsaný v souborech projektu.

## Orientace v projektu

Na začátku práce zkontroluj aktuální stav:

pwd
git branch --show-current
git status --short

Přečti si relevantní kontextové soubory, zejména:

- PROJECT_CONTEXT.md
- TODO.md
- DECISIONS.md
- START_NEW_AI_CHAT.md
- README.md

## Práce se změnami

- Nerevertuj nesouvisející necommitované změny.
- Nepřepisuj práci uživatele.
- Před úpravami si prohlédni relevantní soubory.
- Dělej malé, cílené změny.
- Nepřepisuj celé soubory, pokud stačí menší úprava.
- Nedělej velké refaktory bez výslovného důvodu.
- Neměň architekturu projektu, pokud to aktuální úkol nevyžaduje.
- Zachovej stávající styl kódu, strukturu projektu a použité technologie.
- Neinstaluj nové závislosti bez jasného důvodu.

## Existující nebo rozpracovaný projekt

Pokud jde o starší nebo rozpracovaný projekt:

- Nejprve pochop aktuální stav.
- Respektuj existující strukturu.
- Neprováděj automatické „vylepšování“ mimo aktuální úkol.
- Pokud najdeš technický dluh, stručně ho zmiň nebo zapiš do TODO.md, pokud je to vhodné.
- Upřednostni bezpečné pokračování před ideálním přepisem.

## Ověření

Po změnách spusť odpovídající ověření, pokud je v projektu zřejmé.

Použij například podle typu projektu:

Go:

go test ./...
go build ./...

Node / frontend:

npm test
npm run build

PHP / Laravel:

php artisan test
composer test

Pokud není jasné, jak projekt ověřit, podívej se do README.md, PROJECT_CONTEXT.md nebo konfiguračních souborů.

Pokud ověření nelze spustit, napiš proč.

## Kontextové soubory

Pokud se během práce změnil stav projektu, aktualizuj podle potřeby:

- PROJECT_CONTEXT.md
- TODO.md
- DECISIONS.md
- START_NEW_AI_CHAT.md

Zapisuj jen relevantní informace. Nezapisuj domněnky jako fakta.

## Závěr práce

Na konci stručně uveď:

- co bylo změněno
- které soubory byly upraveny
- jak bylo ověřeno
- co případně zůstává nedokončené

Potom ukaž:

git status --short

A navrhni commit message.

## WhisprGo specifika

- Bezne nastaveni patri do `config.yaml` a maji jit menit i pres `whisprgo config ...`.
- API klice nepatri do `config.yaml`; pouzivej `whisprgo auth ...` a `.env` soubory.
