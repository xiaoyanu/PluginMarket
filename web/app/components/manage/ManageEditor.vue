<script setup lang="ts">
import {config, MdEditor, type ToolbarNames, type UploadImgEvent} from 'md-editor-v3';
import 'md-editor-v3/lib/style.css';

const text = defineModel<string>({default: ''});
const assetUrl = useAssetUrl();
const toolbars: ToolbarNames[] = [
  'bold', 'underline', 'italic', 'strikeThrough', '-',
  'title', 'sub', 'sup', 'quote', 'unorderedList', 'orderedList', 'task', '-',
  'codeRow', 'code', 'link', 'image', 'table', 'mermaid', 'katex', '-',
  'revoke', 'next', '=',
  'prettier', 'pageFullscreen', 'preview', 'previewOnly', 'catalog'
];
config({
  editorExtensions: {
    highlight: {
      js: "https://s4.zstatic.net/ajax/libs/highlight.js/11.11.1/highlight.min.js",
      css: {
        'atom': {
          light: "https://s4.zstatic.net/ajax/libs/highlight.js/11.11.1/styles/atom-one-light.min.css",
          dark: 'https://s4.zstatic.net/ajax/libs/highlight.js/11.11.1/styles/atom-one-dark.min.css'
        }
      }
    },
    cropper: {
      js: "https://s4.zstatic.net/npm/cropperjs@1.6.2/dist/cropper.min.js",
      css: "https://s4.zstatic.net/npm/cropperjs@1.6.2/dist/cropper.min.css"
    },
    echarts: {
      js: "https://s4.zstatic.net/ajax/libs/echarts/6.0.0/echarts.min.js"
    },
    katex: {
      js: "https://s4.zstatic.net/ajax/libs/KaTeX/0.16.4/katex.min.js",
      css: "https://s4.zstatic.net/ajax/libs/KaTeX/0.16.4/katex.min.css"
    },
    mermaid: {
      js: "https://s4.zstatic.net/ajax/libs/mermaid/11.12.0/mermaid.min.js"
    },
    prettier: {
      parserMarkdownJs: "https://npm.onmicrosoft.cn/prettier@3.8.1/plugins/markdown.js",
      standaloneJs: "https://npm.onmicrosoft.cn/prettier@3.8.1/standalone.js"
    },
    screenfull: {
      js: "https://s4.zstatic.net/ajax/libs/screenfull.js/5.2.0/screenfull.js"
    }
  }
});

const onUploadImg: UploadImgEvent = async (files, callback) => {
  try {
    const urls = await Promise.all(files.map(async (file) => {
      const body = new FormData();
      body.append('file', file);
      const response: any = await useApiFetch('/upload/image', {
        method: 'POST',
        body
      });
      const url = response.data?.url;
      if (!url) throw new Error('图片上传成功，但未返回图片地址');
      return assetUrl(url);
    }));
    callback(urls);
  } catch (error) {
    showApiError(error, '正文图片上传失败');
  }
}
</script>

<template>
  <MdEditor v-model="text" class="pm-editor" :toolbars="toolbars" placeholder="插件详细介绍"
            @onUploadImg="onUploadImg"/>
</template>

<style scoped lang="scss">
:deep(.md-editor-footer) {
  align-items: center;
}

.pm-editor {
  width: 100%;
  min-width: 0;
  border-radius: 4px !important;
}

@media (max-width: 767px) {
  :deep(.md-editor-toolbar-wrapper) {
    overflow-x: auto;
    overscroll-behavior-x: contain;
  }

  :deep(.md-editor-toolbar) {
    width: max-content;
  }
}
</style>
