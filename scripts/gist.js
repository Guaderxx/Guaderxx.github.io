'use strict';

/**
 * Gist tag — restores the `{% gist %}` tag that Hexo removed in v8.
 *
 * Hexo 6 shipped this tag in core (lib/plugins/tag/gist.js); Hexo 8 dropped it,
 * which broke every post using it. This is the original implementation, kept
 * byte-for-byte compatible so existing posts render exactly as before.
 *
 * Syntax:
 *   {% gist gist_id [filename] %}
 */

hexo.extend.tag.register('gist', function(args) {
  const id = args.shift();
  const file = args.length ? `?file=${args[0]}` : '';

  return `<script src="//gist.github.com/${id}.js${file}"></script>`;
});
