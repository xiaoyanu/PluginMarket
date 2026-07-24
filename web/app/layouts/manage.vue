<script setup lang="ts">
import {
  PhAppWindow,
  PhGear,
  PhHouse,
  PhIdentificationCard,
  PhMedal,
  PhMegaphone,
  PhPuzzlePiece,
  PhShieldCheck,
  PhStar,
  PhTag,
  PhUsers
} from "@phosphor-icons/vue";

const route = useRoute();
const userStore = useUserStore()

const menuItems = computed(() => {
  const isAdmin = userStore.userInfo?.power === 1
  return [
    { name: '主页', path: '/manage', icon: PhHouse },
    { name: '我的插件', path: '/manage/plugin', icon: PhPuzzlePiece },
    { name: '我的收藏', path: '/manage/star', icon: PhStar },
    { name: '我的信息', path: '/manage/info', icon: PhIdentificationCard },
    { name: '审核插件', path: '/manage/check', icon: PhShieldCheck, adminOnly: true },
    { name: '用户管理', path: '/manage/user', icon: PhUsers, adminOnly: true },
    { name: '通知管理', path: '/manage/notification', icon: PhMegaphone, adminOnly: true },
    { name: '框架管理', path: '/manage/frame', icon: PhAppWindow, adminOnly: true },
    { name: '称号管理', path: '/manage/title', icon: PhMedal, adminOnly: true },
    { name: '标签管理', path: '/manage/tag', icon: PhTag, adminOnly: true },
    { name: '全局设置', path: '/manage/setting', icon: PhGear, adminOnly: true },
  ].filter(item => !item.adminOnly || isAdmin)
})
</script>

<template>
  <NuxtLayout name="default">
    <div class="manage-layout flex flex-col md:flex-row items-stretch md:items-start gap-3 md:gap-5 max-w-350 m-auto px-3 md:px-5 mb-6 md:mb-10">
      <aside class="manage-nav w-full md:w-60 bg-white rounded-[16px] md:rounded-[20px] p-3 md:p-5 shadow-pmbox md:sticky md:top-25 shrink-0">
        <div class="item-box">
          <div class="item-title">管理</div>
          <div class="item-content">
            <ul class="framework-list">
              <li
                  v-for="item in menuItems"
                  :key="item.path"
                  :class="{ 'active': route.path === item.path }"
                  @click="navigateTo(item.path)"
              >
                <component :is="item.icon" :size="18"/>
                {{ item.name }}
              </li>
            </ul>
          </div>
        </div>
      </aside>
      <main class="manage-content flex-1 min-w-0 w-full select-none">
        <slot/>
      </main>
    </div>
  </NuxtLayout>
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
    margin-bottom: 12px;
    display: flex;
    align-items: center;
    color: #1E293B;

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

  .item-content {
    color: #475569;

    .framework-list {
      list-style: none;
      padding: 0;
      margin: 0;

      li {
        padding: 8px 12px;
        border-radius: 7px;
        margin-bottom: 4px;
        cursor: pointer;
        font-size: 14px;
        display: flex;
        align-items: center;
        gap: 8px;
        color: #555;
        transition: all .3s;

        &:hover {
          background-color: #f5f6f8;

        }

        &.active {
          background-color: #5692F5;
          color: #FFF;
          box-shadow: 0 2px 8px rgba(66, 133, 244, 0.35);
        }
      }
    }
  }
}

@media (max-width: 767px) {
  .manage-nav {
    overflow-x: auto;
    scrollbar-width: none;

    &::-webkit-scrollbar {
      display: none;
    }

    .item-box {
      margin: 0;
    }

    .item-title {
      display: none;
    }

    .framework-list {
      display: flex;
      width: max-content;
      gap: 8px;

      li {
        min-height: 40px;
        margin: 0;
        padding: 10px 12px;
        white-space: nowrap;
      }
    }
  }
}
</style>
