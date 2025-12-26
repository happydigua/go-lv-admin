<template>
  <n-layout has-sider style="height: 100vh;">
    <n-layout-sider
      bordered
      collapse-mode="width"
      :collapsed-width="64"
      :width="240"
      :collapsed="collapsed"
      show-trigger
      @collapse="handleCollapse"
      @expand="handleExpand"
      :native-scrollbar="false"
      inverted
      class="sider"
    >
      <div class="logo">
        <img v-if="settingStore.siteLogo" :src="settingStore.siteLogo" alt="logo" class="logo-img" />
        <span v-if="!collapsed" class="logo-text">{{ settingStore.siteName }}</span>
        <span v-else-if="!settingStore.siteLogo">LV</span>
      </div>
      <n-menu
        :collapsed="collapsed"
        :collapsed-width="64"
        :collapsed-icon-size="22"
        :options="menuOptions"
        :value="activeKey"
        @update:value="handleMenuUpdate"
        inverted
      />
    </n-layout-sider>
    <n-layout>
      <n-layout-header bordered class="header">
        <div class="header-left">
          <n-breadcrumb>
            <n-breadcrumb-item>首页</n-breadcrumb-item>
            <n-breadcrumb-item v-if="currentTitle">{{ currentTitle }}</n-breadcrumb-item>
          </n-breadcrumb>
        </div>
        <div class="header-right">
          <LocaleSwitcher />
          <n-dropdown :options="userDropdownOptions" @select="handleUserDropdown">
            <n-space align="center" style="cursor: pointer;">
              <n-avatar
                round
                size="small"
                :src="userStore.userInfo?.avatar"
                fallback-src="https://via.placeholder.com/40"
              />
              <span>{{ userStore.userInfo?.nickname || 'Admin' }}</span>
            </n-space>
          </n-dropdown>
        </div>
      </n-layout-header>
      <!-- 多标签栏 -->
      <TabBar />
      <n-layout-content class="content" :native-scrollbar="false">
        <router-view v-slot="{ Component }">
          <keep-alive>
            <component :is="Component" />
          </keep-alive>
        </router-view>
      </n-layout-content>
    </n-layout>
  </n-layout>
</template>

<script setup lang="ts">
import { ref, computed, h, onMounted, watch, markRaw } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useUserStore } from '@/store/user';
import { useTabsStore } from '@/store/tabs';
import { NIcon } from 'naive-ui';
import TabBar from '@/components/TabBar.vue';
import LocaleSwitcher from '@/components/LocaleSwitcher.vue';
import {
  HomeOutline,
  PersonOutline,
  PeopleOutline,
  MenuOutline,
  SettingsOutline,
  LogOutOutline,
  PersonCircleOutline,
  DocumentTextOutline,
  ConstructOutline,
  CodeOutline,
  FolderOutline
} from '@vicons/ionicons5';

import { useSettingStore } from '@/store/setting';

const router = useRouter();
const route = useRoute();
const userStore = useUserStore();
const tabsStore = useTabsStore();
const settingStore = useSettingStore();

// 从 localStorage 读取折叠状态
const collapsed = ref(localStorage.getItem('sidebarCollapsed') === 'true');

// 监听折叠变化并保存
const handleCollapse = () => {
  collapsed.value = true;
  localStorage.setItem('sidebarCollapsed', 'true');
};

const handleExpand = () => {
  collapsed.value = false;
  localStorage.setItem('sidebarCollapsed', 'false');
};

const activeKey = computed(() => route.path);

const currentTitle = computed(() => {
  const matched = route.matched.find(r => r.meta?.title);
  return matched?.meta?.title || '';
});

// 图标映射
const iconMap: Record<string, any> = {
  HomeOutline: markRaw(HomeOutline),
  SettingsOutline: markRaw(SettingsOutline),
  PersonOutline: markRaw(PersonOutline),
  PeopleOutline: markRaw(PeopleOutline),
  MenuOutline: markRaw(MenuOutline),
  PersonCircleOutline: markRaw(PersonCircleOutline),
  DocumentTextOutline: markRaw(DocumentTextOutline),
  ConstructOutline: markRaw(ConstructOutline),
  CodeOutline: markRaw(CodeOutline),
  FolderOutline: markRaw(FolderOutline)
};

// 渲染图标辅助函数
function renderIcon(iconName: string) {
  const icon = iconMap[iconName];
  return icon ? () => h(NIcon, null, { default: () => h(icon) }) : undefined;
}

// 将数据库菜单转换为 Naive UI 菜单格式
function transformMenus(menus: any[]): any[] {
  return menus.map(menu => {
    const item: any = {
      label: menu.title,
      // 目录类型(type=1)不设置可导航的 key，只用于展开子菜单
      key: menu.type === 1 ? `dir-${menu.ID}` : menu.path,
      icon: menu.icon ? renderIcon(menu.icon) : undefined
    };
    if (menu.children && menu.children.length > 0) {
      item.children = transformMenus(menu.children);
    }
    return item;
  });
}

// 动态菜单配置
const menuOptions = computed(() => {
  return transformMenus(userStore.menus);
});

// 用户下拉菜单
const userDropdownOptions = [
  {
    label: '个人中心',
    key: 'profile'
  },
  {
    type: 'divider',
    key: 'd1'
  },
  {
    label: '退出登录',
    key: 'logout',
    icon: () => h(NIcon, null, { default: () => h(LogOutOutline) })
  }
];

const handleMenuUpdate = (key: string) => {
  // 目录类型的 key 以 dir- 开头，不进行路由跳转
  if (!key.startsWith('dir-')) {
    router.push(key);
  }
};

const handleUserDropdown = (key: string) => {
  if (key === 'logout') {
    userStore.logout();
    router.push('/login');
  } else if (key === 'profile') {
    router.push('/profile');
  }
};

// 监听路由变化，自动添加标签
watch(
  () => route.path,
  (path) => {
    if (path === '/login') return;
    
    const title = route.meta?.title as string || '未命名';
    const name = route.name as string || '';
    
    tabsStore.addTab({
      name,
      path,
      title,
      closable: path !== '/dashboard'
    });
  },
  { immediate: true }
);

// 页面加载时获取菜单
onMounted(async () => {
  if (userStore.menus.length === 0) {
    await userStore.fetchMenus();
  }
  if (userStore.permissions.length === 0) {
    await userStore.fetchPermissions();
  }
});
</script>

<style scoped>
.sider {
  background: #001529;
}

.logo {
  height: 64px;
  display: flex;
  justify-content: center;
  align-items: center;
  font-size: 18px;
  font-weight: bold;
  color: #fff;
  background-color: #002140; /* 略深一点的背景 */
  overflow: hidden;
  white-space: nowrap;
}

.logo-img {
  width: 32px;
  height: 32px;
  margin-right: 8px;
}

.logo-text {
  font-size: 16px;
}

.header {
  height: 64px;
  padding: 0 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: white;
}

.header-left {
  display: flex;
  align-items: center;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.content {
  padding: 24px;
  background: #f0f2f5;
}

:deep(.n-menu) {
  background: #001529;
}

:deep(.n-menu-item-content) {
  color: rgba(255, 255, 255, 0.65) !important;
}

:deep(.n-menu-item-content:hover) {
  color: white !important;
}

:deep(.n-menu-item-content--selected) {
  color: white !important;
  background: #1890ff !important;
}
</style>
