# Анализ инструментов и технологий для проекта s4s

## 1. Дизайн-система: Untitled UI
### Сравнение с альтернативами
Критерии: совместимость с React, наличие Figma-версии, кастомизация, производительность, сообщество.

| Дизайн-система | Плюсы | Минусы | Цена/Лицензия | Применимость для s4s |
|----------------|-------|-------|---------------|----------------------|
| **Untitled UI** | - Крупнейшая коллекция React-компонентов на Tailwind.<br>- Полная интеграция с Figma (Auto Layout 5.0).<br>- 1000+ компонентов, иконок, шаблонов.<br>- Лёгкий (PRO LITE версия на 55% легче).<br>- Бесплатные обновления lifetime. | - Зависимость от Tailwind (overhead, если не используем).<br>- Не полностью OSS.<br>- Дорогие лицензии для больших команд. | - PRO SOLO: ~$99/год (индивидуал).<br>- PRO TEAM: ~$299/год (до 5 users).<br>- PRO ENTERPRISE: custom (>12 users).<br>- Лицензия: Unlimited personal/commercial projects, но нельзя resale или создавать competing kits. Sharing только с team-лицензией. | Идеально для минималистичного UI sales-tool, быстрый старт с готовыми блоками (дашборд, editor). Рекомендую PRO TEAM для команды 2-5 devs. |
| **Shadcn UI** | - Полностью OSS (MIT).<br>- Лёгкий, копи-паст компоненты на Tailwind.<br>- Высокая кастомизация, нет vendor lock-in. | - Нет полноценного Figma-kita.<br>- Меньше готовых шаблонов. | Free (OSS). | Альтернатива для tight бюджета, но Untitled быстрее для MVP. |
| **Chakra UI** | - Доступный (a11y-focused).<br>- Themeable, responsive.<br>- Большое сообщество. | - Тяжёлый bundle size.<br>- Меньше фокуса на Figma. | Free (MIT). | Overkill для minimalism s4s. |
| **Material UI** | - Зрелый, 1000+ компонентов.<br>- Google-backed, сильная документация.<br>- Figma-kit доступен. | - Google-style, не минималистичный.<br>- Большой bundle size. | Free (MIT), premium themes ~$59+. | Не наш стиль, material design не нужен. |
| **Tailwind UI** | - Pure Tailwind, no JS overhead.<br>- Готовые шаблоны страниц. | - Не полноценные React-компоненты (HTML+CSS).<br>- Нет Figma. | $299 one-time (unlimited). | Базовый, Untitled даёт больше React-ready компонентов. |

**Вывод**: Untitled UI оптимален для React + Figma, ускоряет прототипирование. Лицензия PRO TEAM (~$299/год) подходит для команды. Нет ограничений на commercial use для нашего случая.

## 2. Workflow Editor: xyflow (React Flow)
### Сравнение с альтернативами
Критерии: React-support, кастомизация, производительность, сообщество.

| Библиотека | Плюсы | Минусы | Лицензия | Применимость для s4s |
|------------|-------|-------|----------|-----------------------|
| **xyflow (React Flow)** | - MIT OSS, полная кастомизация nodes/connections.<br>- Drag-and-drop, zoom, mini-map out-of-box.<br>- Используют Stripe, Typeform — scalable.<br>- Хорошая документация, примеры для workflows. | - Low-level (больше кода для full editor).<br>- Нет built-in persistence. | MIT (free commercial use, no restrictions). | Идеально для React-based node editor, подходит для простых sales-flows. OSS, нет лимитов на users/projects. |
| **Rete.js** | - Modular, plugin-based.<br>- Хороша для complex graphs.<br>- Vue/React/Svelte support. | - Крутая learning curve.<br>- Меньше сообщество. | MIT (free). | Альтернатива для plugin-heavy систем, но xyflow проще для MVP. |
| **Drawflow** | - Лёгкий, vanilla JS.<br>- Simple API. | - Нет native React.<br>- Limited features. | MIT. | Базовый, xyflow богаче для interactive UI. |
| **GoJS** | - Профессиональный, много diagram types.<br>- High perf. | - Платный (~$895+).<br>- Не OSS. | Commercial license (restrictions). | Дорогой, overkill; xyflow free. |
| **LiteGraph.js** | - Lightweight, canvas-based.<br>- Good для real-time. | - Не React-specific.<br>- Старый maintenance. | MIT. | Простой, xyflow modernнее. |

**Вывод**: xyflow — лучший выбор для React-based workflow editor. Полностью OSS (MIT), нет ограничений на commercial use/users. Подходит для node-based UI, минималистичных sales-flows.

## 3. Анализ n8n (бенчмарк)
- **Лицензия**: Fair-code (Sustainable Use License). Разрешено internal business use, forking, self-host. Commercial use ok, но нельзя resale как SaaS без enterprise license. Для s4s — безопасно, т.к. мы строим собственный продукт, вдохновлённый n8n.
- **Открытые компоненты**:
  - Workflow engine: Полное ядро выполнения (JSON parsing, execution).
  - Nodes: 400+ интеграций, custom JS/Python, триггеры/действия.
  - Editor: Drag-and-drop UI (Vue-based).
  - Cloud-only: SSO, multi-user spaces, advanced analytics — реализуем сами.
- **Ресёрч**:
  - Код: Node.js backend, Vue frontend, BullMQ для queues, PostgreSQL для storage.
  - Для s4s: Переводим на GoLang для performance (concurrency для orchestration). UI на React (xyflow + Untitled UI).
  - LLaMA: Помогла разобрать n8n orchestration (trigger -> node execution -> logging). Полезно для архитектуры и логики.
  - Конкуренты (Zapier, Make, Pipedream): Фокус на templates (Zapier: 5k+ integrations). Для s4s — 10-20 sales-specific templates (leads, follow-ups) для MVP.

## 4. Прогресс и выводы
- **Untitled UI**: Рекомендую PRO TEAM лицензию ($299/год) для команды до 5 человек. Подходит для быстрого старта дашборда/editor’а.
- **xyflow**: Выбрано за React-support, OSS, кастомизацию. Идеально для workflow editor.
- **n8n**: Полезен для изучения workflow engine. Переносим логику на GoLang, добавляем sales-specific шаблоны.
- **Конкуренты**: Упор на templates для non-tech users (маркетологи, сейлы). Для MVP s4s — фокус на 10-20 интеграций (Google Sheets, Email, Slack, CRM).
- **Следующие шаги**:
  - Залить таблицы в Google Sheets в папку "Состав продукта".
  - Настроить CI/CD (GitHub Actions) и local dev env (Docker, PostgreSQL, Redis).
  - Начать прототипирование editor’а (xyflow) и дашборда (Untitled UI).
  - Провести встречу для согласования tech stack (Go vs PHP) и приоритетов Sprint 1.