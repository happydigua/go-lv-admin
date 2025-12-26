import { defineStore } from 'pinia';
import { ref } from 'vue';

export interface TabItem {
    name: string;
    path: string;
    title: string;
    closable: boolean;
}

export const useTabsStore = defineStore('tabs', () => {
    const tabs = ref<TabItem[]>([
        { name: 'Dashboard', path: '/dashboard', title: '仪表盘', closable: false }
    ]);
    const activeTab = ref('/dashboard');

    // 添加标签
    const addTab = (tab: TabItem) => {
        const exists = tabs.value.find(t => t.path === tab.path);
        if (!exists) {
            tabs.value.push(tab);
        }
        activeTab.value = tab.path;
    };

    // 关闭标签
    const closeTab = (path: string) => {
        const index = tabs.value.findIndex(t => t.path === path);
        if (index === -1) return;

        // 不能关闭仪表盘
        if (!tabs.value[index].closable) return;

        // 如果关闭的是当前标签，需要切换到其他标签
        if (activeTab.value === path) {
            // 优先切换到右边，没有就切换到左边
            const nextTab = tabs.value[index + 1] || tabs.value[index - 1];
            if (nextTab) {
                activeTab.value = nextTab.path;
            }
        }

        tabs.value.splice(index, 1);
        return activeTab.value;
    };

    // 关闭其他标签
    const closeOtherTabs = (path: string) => {
        tabs.value = tabs.value.filter(t => !t.closable || t.path === path);
        activeTab.value = path;
    };

    // 关闭右侧标签
    const closeRightTabs = (path: string) => {
        const index = tabs.value.findIndex(t => t.path === path);
        if (index !== -1) {
            tabs.value = tabs.value.filter((t, i) => i <= index || !t.closable);
        }
    };

    // 关闭所有标签（除了首页）
    const closeAllTabs = () => {
        tabs.value = tabs.value.filter(t => !t.closable);
        activeTab.value = '/dashboard';
        return activeTab.value;
    };

    // 设置当前标签
    const setActiveTab = (path: string) => {
        activeTab.value = path;
    };

    return {
        tabs,
        activeTab,
        addTab,
        closeTab,
        closeOtherTabs,
        closeRightTabs,
        closeAllTabs,
        setActiveTab
    };
});
