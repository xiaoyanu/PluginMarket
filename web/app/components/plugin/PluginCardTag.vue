<script setup lang="ts">
import {PhDownloadSimple, PhEye, PhStar} from "@phosphor-icons/vue";
defineProps({
  type: {
    type: Number,
    required: true
  },
  views: {
    type: Number,
    required: true
  },
  download: {
    type: Number,
    required: true
  },
  stars: {
    type: Number,
    required: true
  }
})
</script>

<template>
  <div class="card-footer">
    <span v-if="type==0" class="tag tag-free">免费</span>
    <span v-else class="tag tag-buy">收费</span>
    <span class="tag"><PhEye weight="duotone"/>&nbsp;{{ views }}</span>
    <span class="tag"><PhDownloadSimple weight="duotone"/>&nbsp;{{ download }}</span>
    <span class="tag"><PhStar weight="duotone"/>&nbsp;{{ stars }}</span>
  </div>
</template>

<style scoped lang="scss">
.card-footer {
  color: #909399;
  display: flex;
  align-items: center;
  gap: 5px;


  .tag {
    display: flex;
    align-items: center;
    padding: 2px 6px;
    border-radius: 6px;
    font-size: 11px;
    background-color: #F4F4F5;
    color: #909399;
  }

  .tag-buy {
    background: #FFEDD9;
    color: #FF8400;
  }

  .tag-free {
    background: #E2F6DA;
    color: #39C408;
  }

}
</style>