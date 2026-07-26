import { config } from 'md-editor-v3'

export const configureMarkdownEditor = () => config({
  editorExtensions: {
    highlight: {
      js: 'https://s4.zstatic.net/ajax/libs/highlight.js/11.11.1/highlight.min.js',
      css: {
        atom: {
          light: 'https://s4.zstatic.net/ajax/libs/highlight.js/11.11.1/styles/atom-one-light.min.css',
          dark: 'https://s4.zstatic.net/ajax/libs/highlight.js/11.11.1/styles/atom-one-dark.min.css',
        },
      },
    },
    cropper: {
      js: 'https://s4.zstatic.net/npm/cropperjs@1.6.2/dist/cropper.min.js',
      css: 'https://s4.zstatic.net/npm/cropperjs@1.6.2/dist/cropper.min.css',
    },
    echarts: {
      js: 'https://s4.zstatic.net/ajax/libs/echarts/6.0.0/echarts.min.js',
    },
    katex: {
      js: 'https://s4.zstatic.net/ajax/libs/KaTeX/0.16.4/katex.min.js',
      css: 'https://s4.zstatic.net/ajax/libs/KaTeX/0.16.4/katex.min.css',
    },
    mermaid: {
      js: 'https://s4.zstatic.net/ajax/libs/mermaid/11.12.0/mermaid.min.js',
    },
    prettier: {
      parserMarkdownJs: 'https://npm.onmicrosoft.cn/prettier@3.8.1/plugins/markdown.js',
      standaloneJs: 'https://npm.onmicrosoft.cn/prettier@3.8.1/standalone.js',
    },
    screenfull: {
      js: 'https://s4.zstatic.net/ajax/libs/screenfull.js/5.2.0/screenfull.js',
    },
  },
})
