<template>
  <div class="dashboard">
    <!-- 统计卡片 -->
    <n-grid :cols="4" :x-gap="16" :y-gap="16">
      <n-gi>
        <n-card class="stat-card stat-card-primary">
          <n-statistic :label="t('dashboard.totalUsers')" :value="stats.userCount">
            <template #prefix>
              <n-icon :component="PeopleOutline" />
            </template>
          </n-statistic>
        </n-card>
      </n-gi>
      <n-gi>
        <n-card class="stat-card stat-card-success">
          <n-statistic :label="t('dashboard.totalRoles')" :value="stats.roleCount">
            <template #prefix>
              <n-icon :component="ShieldCheckmarkOutline" />
            </template>
          </n-statistic>
        </n-card>
      </n-gi>
      <n-gi>
        <n-card class="stat-card stat-card-warning">
          <n-statistic :label="t('dashboard.totalMenus')" :value="stats.menuCount">
            <template #prefix>
              <n-icon :component="MenuOutline" />
            </template>
          </n-statistic>
        </n-card>
      </n-gi>
      <n-gi>
        <n-card class="stat-card stat-card-info">
          <n-statistic :label="t('dashboard.todayVisits')" :value="stats.todayVisit">
            <template #prefix>
              <n-icon :component="EyeOutline" />
            </template>
          </n-statistic>
        </n-card>
      </n-gi>
    </n-grid>

    <!-- 图表区域 -->
    <n-grid :cols="2" :x-gap="16" :y-gap="16" style="margin-top: 16px;">
      <n-gi>
        <n-card :title="t('dashboard.visitTrend')">
          <v-chart :option="visitTrendOption" autoresize style="height: 300px;" />
        </n-card>
      </n-gi>
      <n-gi>
        <n-card :title="t('dashboard.userGrowth')">
          <v-chart :option="userGrowthOption" autoresize style="height: 300px;" />
        </n-card>
      </n-gi>
    </n-grid>

    <n-grid :cols="2" :x-gap="16" :y-gap="16" style="margin-top: 16px;">
      <n-gi>
        <n-card :title="t('dashboard.moduleStats')">
          <v-chart :option="moduleStatsOption" autoresize style="height: 300px;" />
        </n-card>
      </n-gi>
      <n-gi>
        <n-card :title="t('dashboard.recentLogs')">
          <n-list bordered>
            <n-list-item v-for="(log, index) in charts.latestLogs" :key="index">
              <n-thing>
                <template #header>
                  <n-text strong>{{ log.username }}</n-text>
                  <n-tag size="small" :type="getTagType(log.action)" style="margin-left: 8px;">
                    {{ log.action }}
                  </n-tag>
                </template>
                <template #description>
                  {{ log.module }} · {{ log.createdAt }}
                </template>
              </n-thing>
            </n-list-item>
            <n-empty v-if="charts.latestLogs.length === 0" description="暂无日志" />
          </n-list>
        </n-card>
      </n-gi>
    </n-grid>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { useUserStore } from '@/store/user';
import { getDashboardStats, getDashboardCharts } from '@/api/dashboard';
import { PeopleOutline, ShieldCheckmarkOutline, MenuOutline, EyeOutline } from '@vicons/ionicons5';
import VChart from 'vue-echarts';
import { use } from 'echarts/core';
import { CanvasRenderer } from 'echarts/renderers';
import { LineChart, PieChart, BarChart } from 'echarts/charts';
import { GridComponent, TooltipComponent, LegendComponent, TitleComponent } from 'echarts/components';

use([CanvasRenderer, LineChart, PieChart, BarChart, GridComponent, TooltipComponent, LegendComponent, TitleComponent]);

const { t } = useI18n();
const userStore = useUserStore();

const stats = ref({
  userCount: 0,
  roleCount: 0,
  menuCount: 0,
  todayVisit: 0
});

const charts = ref({
  visitTrend: { categories: [] as string[], series: [] as number[] },
  userGrowth: { categories: [] as string[], series: [] as number[] },
  moduleStats: [] as { name: string; value: number }[],
  latestLogs: [] as { username: string; action: string; module: string; createdAt: string }[]
});

// 访问趋势图表配置
const visitTrendOption = computed(() => ({
  tooltip: { trigger: 'axis' },
  grid: { left: '5%', right: '4%', bottom: '5%' },
  xAxis: {
    type: 'category',
    boundaryGap: false,
    data: charts.value.visitTrend.categories
  },
  yAxis: { type: 'value' },
  series: [{
    name: t('dashboard.visits'),
    type: 'line',
    smooth: true,
    areaStyle: {
      color: {
        type: 'linear',
        x: 0, y: 0, x2: 0, y2: 1,
        colorStops: [
          { offset: 0, color: 'rgba(24, 144, 255, 0.3)' },
          { offset: 1, color: 'rgba(24, 144, 255, 0.05)' }
        ]
      }
    },
    lineStyle: { color: '#1890ff', width: 2 },
    itemStyle: { color: '#1890ff' },
    data: charts.value.visitTrend.series
  }]
}));

// 用户增长图表配置
const userGrowthOption = computed(() => ({
  tooltip: { trigger: 'axis' },
  grid: { left: '5%', right: '4%', bottom: '5%' },
  xAxis: {
    type: 'category',
    data: charts.value.userGrowth.categories
  },
  yAxis: { type: 'value' },
  series: [{
    name: t('dashboard.userCount'),
    type: 'bar',
    barWidth: '50%',
    itemStyle: {
      color: {
        type: 'linear',
        x: 0, y: 0, x2: 0, y2: 1,
        colorStops: [
          { offset: 0, color: '#52c41a' },
          { offset: 1, color: '#95de64' }
        ]
      },
      borderRadius: [4, 4, 0, 0]
    },
    data: charts.value.userGrowth.series
  }]
}));

// 模块访问统计图表配置
const moduleStatsOption = computed(() => ({
  tooltip: { trigger: 'item', formatter: '{b}: {c} ({d}%)' },
  legend: { orient: 'vertical', left: 'left' },
  series: [{
    name: t('dashboard.moduleAccess'),
    type: 'pie',
    radius: ['40%', '70%'],
    avoidLabelOverlap: false,
    itemStyle: { borderRadius: 10, borderColor: '#fff', borderWidth: 2 },
    label: { show: false, position: 'center' },
    emphasis: {
      label: { show: true, fontSize: 16, fontWeight: 'bold' }
    },
    labelLine: { show: false },
    data: charts.value.moduleStats.map((item, index) => ({
      ...item,
      itemStyle: { color: ['#1890ff', '#52c41a', '#faad14', '#f5222d', '#722ed1', '#13c2c2'][index % 6] }
    }))
  }]
}));

const getTagType = (action: string) => {
  switch (action) {
    case '新增': return 'success';
    case '编辑': return 'info';
    case '删除': return 'error';
    case '查询': return 'default';
    default: return 'default';
  }
};

const fetchStats = async () => {
  try {
    const data: any = await getDashboardStats();
    stats.value = data;
  } catch (error) {
    console.error('Failed to fetch dashboard stats:', error);
  }
};

const fetchCharts = async () => {
  try {
    const data: any = await getDashboardCharts();
    charts.value = data;
  } catch (error) {
    console.error('Failed to fetch dashboard charts:', error);
  }
};

onMounted(() => {
  fetchStats();
  fetchCharts();
});
</script>

<style scoped>
.dashboard {
  min-height: calc(100vh - 64px - 48px);
}

.stat-card {
  text-align: center;
  transition: transform 0.3s, box-shadow 0.3s;
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.stat-card :deep(.n-statistic-value) {
  font-size: 32px;
  font-weight: 600;
}

.stat-card-primary :deep(.n-statistic-value) {
  color: #1890ff;
}

.stat-card-success :deep(.n-statistic-value) {
  color: #52c41a;
}

.stat-card-warning :deep(.n-statistic-value) {
  color: #faad14;
}

.stat-card-info :deep(.n-statistic-value) {
  color: #722ed1;
}

.stat-card :deep(.n-icon) {
  font-size: 24px;
}

.stat-card-primary :deep(.n-icon) {
  color: #1890ff;
}

.stat-card-success :deep(.n-icon) {
  color: #52c41a;
}

.stat-card-warning :deep(.n-icon) {
  color: #faad14;
}

.stat-card-info :deep(.n-icon) {
  color: #722ed1;
}
</style>
