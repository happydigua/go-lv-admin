<template>
  <n-card title="菜单管理">
    <template #header-extra>
      <n-button type="primary" @click="handleAdd(null)">
        <template #icon><n-icon :component="AddCircleOutline" /></template>
        新增菜单
      </n-button>
    </template>
    <n-data-table
      :columns="columns"
      :data="tableData"
      :loading="loading"
      :row-key="(row: any) => row.ID"
      default-expand-all
      :bordered="false"
    />
  </n-card>
  
  <!-- 编辑/新增弹窗 -->
  <n-modal v-model:show="showModal" preset="dialog" :title="modalTitle" style="width: 600px;">
    <n-form
      ref="formRef"
      :model="formData"
      :rules="formRules"
      label-placement="left"
      label-width="80"
    >
      <n-form-item label="上级菜单" path="parent_id">
        <n-tree-select
          v-model:value="formData.parent_id"
          :options="parentOptions"
          placeholder="无（顶级菜单）"
          clearable
          default-expand-all
          key-field="ID"
          label-field="title"
        />
      </n-form-item>
      <n-form-item label="菜单类型" path="type">
        <n-radio-group v-model:value="formData.type">
          <n-radio :value="1">目录</n-radio>
          <n-radio :value="2">菜单</n-radio>
          <n-radio :value="3">按钮</n-radio>
        </n-radio-group>
      </n-form-item>
      <n-form-item label="菜单名称" path="title">
        <n-input v-model:value="formData.title" placeholder="请输入菜单名称" />
      </n-form-item>
      <n-form-item v-if="formData.type !== 3" label="图标" path="icon">
        <n-input v-model:value="formData.icon" placeholder="如: HomeOutline" />
      </n-form-item>
      <n-form-item v-if="formData.type !== 3" label="路由路径" path="path">
        <n-input v-model:value="formData.path" placeholder="如: /system/user" />
      </n-form-item>
      <n-form-item v-if="formData.type === 2" label="组件路径" path="component">
        <n-input v-model:value="formData.component" placeholder="如: views/system/user/index" />
      </n-form-item>
      <n-form-item v-if="formData.type === 3" label="权限标识" path="permission">
        <n-input v-model:value="formData.permission" placeholder="如: system:user:add" />
      </n-form-item>
      <n-form-item label="排序" path="sort">
        <n-input-number v-model:value="formData.sort" :min="0" style="width: 100%;" />
      </n-form-item>
    </n-form>
    <template #action>
      <n-button @click="showModal = false">取消</n-button>
      <n-button type="primary" :loading="submitLoading" @click="handleSubmit">确定</n-button>
    </template>
  </n-modal>
</template>

<script setup lang="ts">
import { h, ref, computed, onMounted, markRaw } from 'vue';
import { NButton, NSpace, NTag, NIcon, useMessage, useDialog } from 'naive-ui';
import { AddCircleOutline, HomeOutline, SettingsOutline, PersonOutline, PeopleOutline, MenuOutline } from '@vicons/ionicons5';
import { getMenuList, createMenu, updateMenu, deleteMenu } from '@/api/system/menu';

const message = useMessage();
const dialog = useDialog();

const loading = ref(false);
const submitLoading = ref(false);
const showModal = ref(false);
const isEdit = ref(false);
const formRef = ref();
const modalTitle = ref('新增菜单');
const tableData = ref<any[]>([]);

const formData = ref({
  ID: 0,
  parent_id: null as number | null,
  title: '',
  icon: '',
  path: '',
  component: '',
  permission: '',
  sort: 0,
  type: 2
});

const formRules = {
  title: { required: true, message: '请输入菜单名称', trigger: 'blur' },
  path: { required: true, message: '请输入路由路径', trigger: 'blur' }
};

const iconMap: Record<string, any> = {
  HomeOutline: markRaw(HomeOutline),
  SettingsOutline: markRaw(SettingsOutline),
  PersonOutline: markRaw(PersonOutline),
  PeopleOutline: markRaw(PeopleOutline),
  MenuOutline: markRaw(MenuOutline),
  AddCircleOutline: markRaw(AddCircleOutline)
};

const columns = [
  { title: '菜单名称', key: 'title' },
  { 
    title: '图标', 
    key: 'icon', 
    render: (row: any) => {
      const icon = iconMap[row.icon];
      return icon ? h(NIcon, { size: 18 }, { default: () => h(icon) }) : '-';
    }
  },
  { title: '路径', key: 'path' },
  { title: '组件', key: 'component' },
  { title: '排序', key: 'sort', width: 80 },
  {
    title: '类型',
    key: 'type',
    render: (row: any) => {
      const types: Record<number, { type: string; label: string }> = {
        1: { type: 'info', label: '目录' },
        2: { type: 'success', label: '菜单' },
        3: { type: 'warning', label: '按钮' }
      };
      const t = types[row.type] || { type: 'default', label: '未知' };
      return h(NTag, { type: t.type as any, size: 'small' }, { default: () => t.label });
    }
  },
  {
    title: '操作',
    key: 'actions',
    width: 200,
    render: (row: any) => h(NSpace, null, {
      default: () => [
        h(NButton, { size: 'small', tertiary: true, type: 'primary', onClick: () => handleAdd(row) }, { default: () => '新增' }),
        h(NButton, { size: 'small', tertiary: true, type: 'info', onClick: () => handleEdit(row) }, { default: () => '编辑' }),
        h(NButton, { size: 'small', tertiary: true, type: 'error', onClick: () => handleDelete(row) }, { default: () => '删除' })
      ]
    })
  }
];

// 构建上级菜单选项
const parentOptions = computed(() => {
  const buildOptions = (items: any[]): any[] => {
    return items.filter(item => item.type !== 3).map(item => ({
      ...item,
      children: item.Children ? buildOptions(item.Children) : undefined
    }));
  };
  return buildOptions(tableData.value);
});

const fetchData = async () => {
  loading.value = true;
  try {
    const res: any = await getMenuList();
    tableData.value = res || [];
  } catch (error) {
    console.error('Failed to fetch menus:', error);
  } finally {
    loading.value = false;
  }
};

const handleAdd = (parent: any) => {
  isEdit.value = false;
  modalTitle.value = parent ? `新增子菜单 - ${parent.title}` : '新增菜单';
  formData.value = {
    ID: 0,
    parent_id: parent?.ID || null,
    title: '',
    icon: '',
    path: '',
    component: '',
    permission: '',
    sort: 0,
    type: 2
  };
  showModal.value = true;
};

const handleEdit = (row: any) => {
  isEdit.value = true;
  modalTitle.value = '编辑菜单';
  formData.value = { 
    ID: row.ID,
    parent_id: row.parent_id,
    title: row.title,
    icon: row.icon,
    path: row.path,
    component: row.component,
    permission: row.permission,
    sort: row.sort,
    type: row.type
  };
  showModal.value = true;
};

const handleDelete = (row: any) => {
  if (row.Children && row.Children.length > 0) {
    message.error('请先删除子菜单');
    return;
  }
  dialog.error({
    title: '删除确认',
    content: `确定要删除菜单 "${row.title}" 吗？`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await deleteMenu(row.ID);
        message.success('删除成功');
        fetchData();
      } catch (error) {
        message.error('删除失败');
      }
    }
  });
};

const handleSubmit = () => {
  formRef.value?.validate(async (errors: any) => {
    if (!errors) {
      submitLoading.value = true;
      try {
        if (isEdit.value) {
          await updateMenu(formData.value.ID, formData.value);
          message.success('更新成功');
        } else {
          await createMenu(formData.value);
          message.success('创建成功');
        }
        showModal.value = false;
        fetchData();
      } catch (error) {
        message.error(isEdit.value ? '更新失败' : '创建失败');
      } finally {
        submitLoading.value = false;
      }
    }
  });
};

onMounted(() => {
  fetchData();
});
</script>
