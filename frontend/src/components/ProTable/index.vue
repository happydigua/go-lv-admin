<template>
  <div class="pro-table">
    <!-- 搜索栏 -->
    <n-card v-if="hasSearch" class="search-card" :bordered="false">
      <n-form
        inline
        :model="searchParams"
        label-placement="left"
        @keyup.enter="handleSearch"
      >
        <template v-for="(config, field) in searchSchema" :key="field">
          <n-form-item :label="config.label">
            <!-- 输入框 -->
            <n-input
              v-if="config.type === 'input'"
              v-model:value="searchParams[field]"
              :placeholder="config.placeholder || '请输入' + config.label"
              clearable
            />
            <!-- 下拉框 -->
            <n-select
              v-else-if="config.type === 'select'"
              v-model:value="searchParams[field]"
              :options="config.options"
              :placeholder="config.placeholder || '请选择' + config.label"
              clearable
              style="width: 160px"
            />
            <!-- 日期选择 -->
            <n-date-picker
              v-else-if="config.type === 'date-picker'"
              v-model:value="searchParams[field]"
              type="date"
              clearable
            />
          </n-form-item>
        </template>
        <n-form-item>
          <n-space>
            <n-button type="primary" @click="handleSearch">
              <template #icon><n-icon :component="SearchOutline" /></template>
              搜索
            </n-button>
            <n-button @click="handleReset">
              <template #icon><n-icon :component="RefreshOutline" /></template>
              重置
            </n-button>
          </n-space>
        </n-form-item>
      </n-form>
    </n-card>

    <!-- 工具栏和表格 -->
    <n-card class="table-card" :bordered="false">
      <!-- 头部工具栏 -->
      <div class="table-header">
        <div class="table-title">
          <slot name="headerTitle">{{ title }}</slot>
        </div>
        <div class="table-tools">
          <n-space>
            <slot name="toolbar"></slot>
            <n-divider vertical />
            <n-tooltip trigger="hover">
              <template #trigger>
                <n-button circle quaternary @click="refresh">
                  <template #icon><n-icon :component="ReloadOutline" /></template>
                </n-button>
              </template>
              刷新
            </n-tooltip>
            <n-tooltip trigger="hover">
              <template #trigger>
                <n-popselect v-model:value="tableSize" :options="sizeOptions" trigger="click">
                  <n-button circle quaternary>
                    <template #icon><n-icon :component="ResizeOutline" /></template>
                  </n-button>
                </n-popselect>
              </template>
              密度
            </n-tooltip>
          </n-space>
        </div>
      </div>

      <!-- 数据表格 -->
      <n-data-table
        remote
        :loading="loading"
        :columns="columns"
        :data="data"
        :pagination="pagination"
        :row-key="actualRowKey"
        :size="tableSize"
        @update:page="handlePageChange"
        @update:page-size="handlePageSizeChange"
        @update:sorter="handleSorterChange"
      />
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue';
import { SearchOutline, RefreshOutline, ReloadOutline, ResizeOutline } from '@vicons/ionicons5';
import type { DataTableColumns, PaginationProps } from 'naive-ui';
import type { ProTableProps } from './types';

const props = defineProps<ProTableProps>();
const emit = defineEmits(['update:checked-row-keys']);

// 状态
const loading = ref(false);
const data = ref<any[]>([]);
const tableSize = ref<'small' | 'medium' | 'large'>('medium');
const sortState = ref<{ columnKey: string | number; order: 'ascend' | 'descend' | false } | null>(null);

const sizeOptions = [
  { label: '紧凑', value: 'small' },
  { label: '默认', value: 'medium' },
  { label: '宽松', value: 'large' }
];

// 搜索参数
const searchParams = reactive<any>({});
const hasSearch = computed(() => props.searchSchema && Object.keys(props.searchSchema).length > 0);

// 分页配置
const pagination = reactive<PaginationProps>({
  page: 1,
  pageSize: 10,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 50, 100],
  prefix: ({ itemCount }) => `共 ${itemCount} 条`
});

// 初始化搜索参数默认值
if (props.searchSchema) {
  Object.entries(props.searchSchema).forEach(([key, config]) => {
    if (config.defaultValue !== undefined) {
      searchParams[key] = config.defaultValue;
    }
  });
}

// 处理 rowKey
const actualRowKey = computed(() => {
  if (typeof props.rowKey === 'string') {
    return (row: any) => row[props.rowKey as string];
  }
  return props.rowKey;
});

// 获取数据
const fetchData = async () => {
  loading.value = true;
  try {
    const params = {
      page: pagination.page,
      pageSize: pagination.pageSize,
      ...searchParams,
      ...(sortState.value?.order ? {
        sortField: sortState.value.columnKey,
        sortOrder: sortState.value.order
      } : {})
    };
    
    const res = await props.request(params);
    data.value = res.list;
    pagination.itemCount = res.total;
  } catch (error) {
    console.error('ProTable fetch data failed:', error);
  } finally {
    loading.value = false;
  }
};

// 事件处理
const handlePageChange = (page: number) => {
  pagination.page = page;
  fetchData();
};

const handlePageSizeChange = (pageSize: number) => {
  pagination.pageSize = pageSize;
  pagination.page = 1;
  fetchData();
};

const handleSorterChange = (sorter: any) => {
  sortState.value = sorter;
  fetchData();
};

const handleSearch = () => {
  pagination.page = 1;
  fetchData();
};

const handleReset = () => {
  // 重置搜索参数
  Object.keys(searchParams).forEach(key => {
    searchParams[key] = props.searchSchema?.[key]?.defaultValue || null;
  });
  handleSearch();
};

const refresh = () => {
  fetchData();
};

// 暴露给父组件的方法
defineExpose({
  refresh,
  handleSearch
});

onMounted(() => {
  fetchData();
});
</script>

<style scoped>
.pro-table {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.search-card :deep(.n-card__content) {
  padding-bottom: 0;
}

.table-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.table-title {
  font-size: 16px;
  font-weight: 500;
}
</style>
