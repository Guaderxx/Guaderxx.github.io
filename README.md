# Guaderxx's Blog

A Hexo blog published to <https://guaderxx.github.io>.

## Stack

| | |
|---|---|
| Hexo | 8.1.2 |
| Theme | Butterfly 5.7.0 (installed from npm) |
| Generated output | `public/` |
| Published to | `gh-pages` branch of `Guaderxx.github.io` |

## How deployment works

`main` holds the source. `hexo deploy` generates the site and force-pushes the
output to the `gh-pages` branch of
`git@github.com:Guaderxx/Guaderxx.github.io.git`, which GitHub Pages serves.
`.deploy_git/` is the deployer's scratch clone — it is git-ignored and safe to
delete at any time.

## Commands

```bash
npm install      # install dependencies
npm run server   # local preview at http://localhost:4000
npm run build    # generate into public/
npm run clean    # remove public/ and db.json
npm run deploy   # generate and push to gh-pages
```

Typical publish flow:

```bash
npx hexo clean && npx hexo deploy
```

## Notes

- **Theme** is an npm dependency (`hexo-theme-butterfly`), so `npm install`
  fully reproduces the blog. Site-specific settings live in
  `_config.butterfly.yml`, which Hexo deep-merges over the theme's own
  `_config.yml` defaults.
- **`scripts/gist.js`** restores the `{% gist %}` tag that Hexo 8 removed from
  core. Three published posts use it. Don't delete this file unless those
  `{% gist %}` usages are replaced with something else.
- **`_config.butterfly.yml`** was migrated from the Butterfly 4 format to
  Butterfly 5, which renamed many settings and does not read the old names.
  The header comment in that file lists every rename.
  Removed by Butterfly 5 and therefore no longer available: `pangu`
  (automatic space between Chinese and Latin characters), `aplayerInject`,
  `daovoice`, `messenger`, AddThis.
- Drafts under `source/_drafts/` are rendered even though they are not
  published, so a broken tag in a draft will fail the build.
