<template>
  <div class="tab-bar">
    <div class="tab-list">
      <div
        v-for="tab in tabsStore.tabs"
        :key="tab.path"
        :class="['tab-item', { active: tabsStore.activeTab === tab.path }]"
        @click="handleTabClick(tab.path)"
        @contextmenu.prevent="handleContextMenu($event, tab)"
      >
        <span class="tab-title">{{ tab.title }}</span>
        <n-icon
          v-if="tab.closable"
          class="tab-close"
          :component="CloseOutline"
          @click.stop="handleTabClose(tab.path)"
        />
      </div>
    </div>

    <!-- 右键菜单 -->
    <n-dropdown
      :show="showContextMenu"
      :x="contextMenuX"
      :y="contextMenuY"
      :options="contextMenuOptions"
      placement="bottom-start"
      @select="handleContextMenuSelect"
      @clickoutside="showContextMenu = false"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import { useTabsStore, type TabItem } from '@/store/tabs';
import { CloseOutline } from '@vicons/ionicons5';

const { t } = useI18n();
const router = useRouter();
const tabsStore = useTabsStore();

const showContextMenu = ref(false);
const contextMenuX = ref(0);
const contextMenuY = ref(0);
const currentContextTab = ref<TabItem | null>(null);

const contextMenuOptions = computed(() => {
  const options = [];
  if (currentContextTab.value?.closable) {
    options.push({ label: t('tabs.close'), key: 'close' });
  }
  options.push(
    { label: t('tabs.closeOther'), key: 'closeOther' },
    { label: t('tabs.closeRight'), key: 'closeRight' },
    { label: t('tabs.closeAll'), key: 'closeAll' }
  );
  return options;
});

const handleTabClick = (path: string) => {
  tabsStore.setActiveTab(path);
  router.push(path);
};

const handleTabClose = (path: string) => {
  const targetPath = tabsStore.closeTab(path);
  if (targetPath && tabsStore.activeTab !== path) {
    // 如果关闭的不是当前标签，不需要跳转
  } else if (targetPath) {
    router.push(targetPath);
  }
};

const handleContextMenu = (e: MouseEvent, tab: TabItem) => {
  currentContextTab.value = tab;
  contextMenuX.value = e.clientX;
  contextMenuY.value = e.clientY;
  showContextMenu.value = true;
};

const handleContextMenuSelect = (key: string) => {
  showContextMenu.value = false;
  if (!currentContextTab.value) return;

  const path = currentContextTab.value.path;

  switch (key) {
    case 'close':
      handleTabClose(path);
      break;
    case 'closeOther':
      tabsStore.closeOtherTabs(path);
      router.push(path);
      break;
    case 'closeRight':
      tabsStore.closeRightTabs(path);
      break;
    case 'closeAll':
      const targetPath = tabsStore.closeAllTabs();
      router.push(targetPath);
      break;
  }
};
</script>

<style scoped>
.tab-bar {
  background: white;
  padding: 8px 16px 0;
  border-bottom: 1px solid #e8e8e8;
}

.tab-list {
  display: flex;
  gap: 4px;
  overflow-x: auto;
}

.tab-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  background: #f5f5f5;
  border-radius: 4px 4px 0 0;
  cursor: pointer;
  white-space: nowrap;
  font-size: 14px;
  color: #666;
  transition: all 0.2s;
  border: 1px solid #e8e8e8;
  border-bottom: none;
  user-select: none;
}

.tab-item:hover {
  background: #e6f7ff;
  color: #1890ff;
}

.tab-item.active {
  background: #1890ff;
  color: white;
  border-color: #1890ff;
}

.tab-title {
  max-width: 120px;
  overflow: hidden;
  text-overflow: ellipsis;
}

.tab-close {
  font-size: 12px;
  border-radius: 50%;
  padding: 2px;
  transition: all 0.2s;
}

.tab-close:hover {
  background: rgba(0, 0, 0, 0.1);
}

.tab-item.active .tab-close:hover {
  background: rgba(255, 255, 255, 0.2);
}
</style>
