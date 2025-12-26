<template>
  <n-card title="角色管理">
    <template #header-extra>
      <n-button type="primary" @click="handleAdd">
        <template #icon><n-icon :component="AddCircleOutline" /></template>
        新增角色
      </n-button>
    </template>
    
    <n-data-table
      :columns="columns"
      :data="tableData"
      :loading="loading"
      :bordered="false"
    />
  </n-card>
  
  <!-- 编辑/新增弹窗 -->
  <n-modal v-model:show="showModal" preset="dialog" :title="modalTitle" style="width: 500px;">
    <n-form
      ref="formRef"
      :model="formData"
      :rules="formRules"
      label-placement="left"
      label-width="80"
    >
      <n-form-item label="角色名称" path="name">
        <n-input v-model:value="formData.name" placeholder="请输入角色名称" />
      </n-form-item>
      <n-form-item label="角色标识" path="keyword">
        <n-input v-model:value="formData.keyword" :disabled="isEdit" placeholder="如: admin, user" />
      </n-form-item>
      <n-form-item label="描述" path="desc">
        <n-input v-model:value="formData.desc" type="textarea" placeholder="请输入描述" />
      </n-form-item>
      <n-form-item label="排序" path="sort">
        <n-input-number v-model:value="formData.sort" :min="0" style="width: 100%;" />
      </n-form-item>
      <n-form-item label="状态" path="status">
        <n-switch v-model:value="formData.status" :checked-value="1" :unchecked-value="0">
          <template #checked>正常</template>
          <template #unchecked>禁用</template>
        </n-switch>
      </n-form-item>
    </n-form>
    <template #action>
      <n-button @click="showModal = false">取消</n-button>
      <n-button type="primary" :loading="submitLoading" @click="handleSubmit">确定</n-button>
    </template>
  </n-modal>

  <!-- 权限分配弹窗 -->
  <n-modal v-model:show="showPermModal" preset="dialog" title="权限分配" style="width: 500px;">
    <n-spin :show="permLoading">
      <div style="max-height: 400px; overflow-y: auto;">
        <n-tree
          :data="menuTreeData"
          :checked-keys="checkedMenus"
          :default-expand-all="true"
          checkable
          cascade
          selectable
          key-field="key"
          label-field="label"
          children-field="children"
          @update:checked-keys="handleCheckMenu"
        />
      </div>
    </n-spin>
    <template #action>
      <n-button @click="showPermModal = false">取消</n-button>
      <n-button type="primary" :loading="permSubmitLoading" @click="handleSubmitPerm">保存</n-button>
    </template>
  </n-modal>
</template>

<script setup lang="ts">
import { h, ref, onMounted } from 'vue';
import { NButton, NSpace, NTag, useMessage, useDialog } from 'naive-ui';
import { AddCircleOutline } from '@vicons/ionicons5';
import { getRoleList, createRole, updateRole, deleteRole } from '@/api/system/role';
import { getMenuList } from '@/api/system/menu';
import { getRoleMenus, setRoleMenus } from '@/api/permission';

const message = useMessage();
const dialog = useDialog();

const loading = ref(false);
const submitLoading = ref(false);
const showModal = ref(false);
const isEdit = ref(false);
const formRef = ref();
const tableData = ref<any[]>([]);

// 权限分配相关
const showPermModal = ref(false);
const permLoading = ref(false);
const permSubmitLoading = ref(false);
const currentRoleId = ref(0);
const checkedMenus = ref<number[]>([]);
const menuTreeData = ref<any[]>([]);

const formData = ref({
  ID: 0,
  name: '',
  keyword: '',
  desc: '',
  sort: 0,
  status: 1
});

const formRules = {
  name: { required: true, message: '请输入角色名称', trigger: 'blur' },
  keyword: { required: true, message: '请输入角色标识', trigger: 'blur' }
};

const modalTitle = ref('新增角色');

// 转换菜单数据为树形结构（供 n-tree 使用）
const convertToTreeData = (menus: any[]): any[] => {
  return menus.map(menu => ({
    key: menu.ID,
    label: menu.title,
    children: menu.children && menu.children.length > 0 ? convertToTreeData(menu.children) : undefined
  }));
};

const columns = [
  { title: 'ID', key: 'ID', width: 80 },
  { title: '角色名称', key: 'name' },
  { title: '角色标识', key: 'keyword' },
  { title: '描述', key: 'desc' },
  { title: '排序', key: 'sort', width: 80 },
  {
    title: '状态',
    key: 'status',
    render: (row: any) => h(NTag, { type: row.status === 1 ? 'success' : 'error' }, { default: () => row.status === 1 ? '正常' : '禁用' })
  },
  {
    title: '操作',
    key: 'actions',
    width: 260,
    render: (row: any) => h(NSpace, null, {
      default: () => [
        h(NButton, { size: 'small', tertiary: true, type: 'warning', onClick: () => handlePerm(row) }, { default: () => '权限' }),
        h(NButton, { size: 'small', tertiary: true, type: 'info', onClick: () => handleEdit(row) }, { default: () => '编辑' }),
        h(NButton, { size: 'small', tertiary: true, type: 'error', disabled: row.ID === 1, onClick: () => handleDelete(row) }, { default: () => '删除' })
      ]
    })
  }
];

const fetchData = async () => {
  loading.value = true;
  try {
    const res: any = await getRoleList();
    tableData.value = res || [];
  } catch (error) {
    console.error('Failed to fetch roles:', error);
  } finally {
    loading.value = false;
  }
};

const handleAdd = () => {
  isEdit.value = false;
  modalTitle.value = '新增角色';
  formData.value = { ID: 0, name: '', keyword: '', desc: '', sort: 0, status: 1 };
  showModal.value = true;
};

const handleEdit = (row: any) => {
  isEdit.value = true;
  modalTitle.value = '编辑角色';
  formData.value = { ...row };
  showModal.value = true;
};

const handleDelete = (row: any) => {
  dialog.error({
    title: '删除确认',
    content: `确定要删除角色 "${row.name}" 吗？`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await deleteRole(row.ID);
        message.success('删除成功');
        fetchData();
      } catch (error) {
        message.error('删除失败');
      }
    }
  });
};

const handlePerm = async (row: any) => {
  currentRoleId.value = row.ID;
  showPermModal.value = true;
  permLoading.value = true;

  try {
    // 获取所有菜单（树形）
    const menuRes: any = await getMenuList();
    menuTreeData.value = convertToTreeData(menuRes || []);

    // 获取角色已分配的菜单
    const roleMenuRes: any = await getRoleMenus(row.ID);
    checkedMenus.value = roleMenuRes || [];
  } catch (error) {
    console.error('Failed to fetch permission data:', error);
    message.error('获取权限数据失败');
  } finally {
    permLoading.value = false;
  }
};

const handleCheckMenu = (keys: number[]) => {
  checkedMenus.value = keys;
};

const handleSubmitPerm = async () => {
  permSubmitLoading.value = true;
  try {
    await setRoleMenus(currentRoleId.value, checkedMenus.value);
    message.success('权限分配成功');
    showPermModal.value = false;
  } catch (error) {
    message.error('权限分配失败');
  } finally {
    permSubmitLoading.value = false;
  }
};

const handleSubmit = () => {
  formRef.value?.validate(async (errors: any) => {
    if (!errors) {
      submitLoading.value = true;
      try {
        if (isEdit.value) {
          await updateRole(formData.value.ID, formData.value);
          message.success('更新成功');
        } else {
          await createRole(formData.value);
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
