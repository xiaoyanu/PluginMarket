<script setup lang="ts">
defineProps({
  title: {
    type: String,
    required: true
  },
  desc: {
    type: String,
    required: true
  },
  icon: {
    type: String,
    required: true
  },
  frameIcon: {
    type: String,
    required: true
  },
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
  },
  responsive: {
    type: Boolean,
    default: false
  },
})

</script>

<template>
  <div :class="['card', { 'card--responsive': responsive }]">
    <div class="card-left">
      <img draggable="false" class="main-img" :src="icon" alt="">
      <img draggable="false" class="frame-img" :src="frameIcon" alt="">
    </div>
    <div class="card-right">
      <div class="card-title">{{ title }}</div>
      <div class="card-desc">{{ desc }}</div>
      <div>
        <plugin-card-tag :type="type" :stars="stars" :download="download" :views="views"/>
      </div>
    </div>

  </div>
</template>

<style scoped lang="scss">
.card {
  display: flex;
  align-items: center;
  width: 340px;
  height: 120px;
  gap: 12px;
  padding: 16px;
  background: #FFF;
  border: 1px solid rgba(0, 0, 0, 0.05);
  border-radius: 12px;
  cursor: pointer;
  user-select: none;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.03);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);

  &:hover {
    transform: translateY(-3px);
    box-shadow: 0 8px 20px rgba(0, 0, 0, 0.06);
    border-color: rgba(0, 0, 0, 0.08);
  }

  .card-left {
    position: relative;
    width: 60px;
    height: 60px;
    flex-shrink: 0;

    .main-img {
      width: 60px;
      height: 60px;
      border-radius: 12px;
      border: 1px solid #FFF;
      box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);

    }

    .frame-img {
      width: 22px;
      height: 22px;
      border-radius: 50%;
      position: absolute;
      bottom: -4px;
      right: -4px;
      border: 2px solid #FFF;
      background-color: #FFF;
      box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    }
  }

  .card-right {
    flex: 1;
    min-width: 0;

    .card-title {
      font-size: 14px;
      font-weight: 600;
      margin-bottom: 4px;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
      color: #1E293B;
    }

    .card-desc {
      font-size: 12px;
      color: #64748B;
      margin-bottom: 8px;
      overflow: hidden;
      height: 33px;
      text-overflow: ellipsis;
    }
  }
}

@media (max-width: 767px) {
  .card--responsive {
    width: 100%;
    min-width: 0;
    height: auto;
    min-height: 112px;
    padding: 14px;
    gap: 12px;

    .card-left {
      width: 54px;
      height: 54px;

      .main-img {
        width: 54px;
        height: 54px;
      }
    }

    .card-right .card-desc {
      display: -webkit-box;
      height: auto;
      min-height: 32px;
      -webkit-box-orient: vertical;
      -webkit-line-clamp: 2;
    }
  }
}
</style>