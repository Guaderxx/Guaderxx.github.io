---
title: mastodon
date: 2026-02-18 13:43:55
tags:
  - mastodon
categories:
  - mastodon
keywords:
  - mastodon
copyright:
copyright_author_href:
copyright_info:
---

Follow by [install mastodon][install_mastodon]

1. Create a user `mastodon`
  - `sudo adduser --disabled-password mastodon`
2. Install NodeJS
  - Use Nvm, install 20LTS
  - `corepack enable` , install `yarn`   
3. Install PgSQL 
  - Follow [install pgsql][install_pgsql]
  - Optional: change configuration by [pgtune][pgtune]
4. Install Mastodon
5. Acquire an SSL certificate
  - `certbot certonly --nginx -d guaderxx.cc`
6. OVER [mastodon][mastodon]


[install_mastodon]: https://docs.joinmastodon.org/admin/install/
[install_pgsql]: https://www.postgresql.org/download/linux/debian/
[pgtune]: https://pgtune.leopard.in.ua/
[mastodon]: https://mastodon.guaderxx.cc/
