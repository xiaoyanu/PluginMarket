<script setup lang="ts">
import {MdEditor, type ToolbarNames, type UploadImgEvent} from 'md-editor-v3';
import 'md-editor-v3/lib/style.css';
import { configureMarkdownEditor } from '~/utils/md-editor-config';

configureMarkdownEditor();

const text = defineModel<string>({default: ''});
const assetUrl = useAssetUrl();
const toolbars: ToolbarNames[] = [
  'bold', 'underline', 'italic', 'strikeThrough', '-',
  'title', 'sub', 'sup', 'quote', 'unorderedList', 'orderedList', 'task', '-',
  'codeRow', 'code', 'link', 'image', 'table', 'mermaid', 'katex', '-',
  'revoke', 'next', '=',
  'prettier', 'pageFullscreen', 'preview', 'previewOnly', 'catalog'
];
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
