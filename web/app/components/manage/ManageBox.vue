<script setup lang="ts">
defineProps({
  title: {
    type: String,
    required: true
  },
  value: {
    type: Number,
    default: -1
  }
})
</script>

<template>
  <section class="manage-box bg-white rounded-[16px] md:rounded-[20px] shadow-pmbox p-3 md:p-5">
    <div class="item-box">
      <div class="manage-box__header flex justify-between items-center mb-4 md:mb-5 flex-wrap gap-2">
        <div class="item-title flex items-center">{{ title }}
          <el-tag class="ml-1.25" v-if="value>-1" size="small">{{ value }}</el-tag>
        </div>
        <div class="manage-box__actions min-h-8 flex items-center">
          <slot name="header"/>
        </div>
      </div>
      <div class="manage-box__content text-[#475569] flex gap-4 md:gap-5 flex-col min-w-0">
        <slot/>
      </div>
    </div>
  </section>
</template>

<style scoped lang="scss">
.item-box {
  margin-bottom: 24px;
  user-select: none;

  &:last-child {
    margin-bottom: 0;
  }

  .item-title {
    font-size: 16px;
    font-weight: bold;
    color: #1E293B;
    flex-shrink: 0;

    &:before {
      content: '';
      display: inline-block;
      width: 4px;
      height: 20px;
      background-color: #00BAAD;
      margin-right: 8px;
      border-radius: 2px;
    }
  }
}

:deep(.el-tag) {
  border-radius: 10px;
}

:deep(.searchForm) {
  .el-select {
    --el-select-width: 100px !important;
  }

  .el-input__wrapper, .el-select__wrapper {
    border-radius: 7px !important;
  }
}

@media (max-width: 767px) {
  .manage-box {
    width: 100%;
    min-width: 0;
  }

  .manage-box__actions {
    max-width: 100%;
  }

  .manage-box__content {
    width: 100%;
    min-width: 0;
  }

  :deep(.searchForm) {
    display: flex;
    flex-wrap: wrap;
    gap: 12px;

    .el-form-item {
      width: 100%;
      margin: 0;
    }

    .el-form-item__content,
    .el-input,
    .el-select {
      width: 100%;
      min-width: 0;
      --el-select-width: 100% !important;
    }
  }

  :deep(.el-table) {
    width: 100%;
    min-width: 0;
  }

  :deep(.pageBox-diy) {
    max-width: 100%;
    justify-content: center;
    flex-wrap: wrap;

    .el-pager li {
      min-width: 30px;
    }
  }
}
</style>
