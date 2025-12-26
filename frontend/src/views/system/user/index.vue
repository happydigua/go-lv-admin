<template>
  <n-card title="用户管理">
    <template #header-extra>
      <n-button type="primary" @click="handleAdd">
        <template #icon><n-icon :component="PersonAddOutline" /></template>
        新增用户
      </n-button>
    </template>
    
    <!-- 搜索区域 -->
    <n-space style="margin-bottom: 16px;">
      <n-input v-model:value="searchForm.username" placeholder="用户名" clearable style="width: 200px;" />
      <n-input v-model:value="searchForm.phone" placeholder="手机号" clearable style="width: 200px;" />
      <n-select v-model:value="searchForm.status" placeholder="状态" clearable style="width: 120px;" :options="statusOptions" />
      <n-button type="primary" @click="fetchData">搜索</n-button>
      <n-button @click="handleReset">重置</n-button>
    </n-space>
    
    <n-data-table
      :columns="columns"
      :data="tableData"
      :loading="loading"
      :pagination="pagination"
      :bordered="false"
      @update:page="handlePageChange"
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
      <n-form-item label="用户名" path="username">
        <n-input v-model:value="formData.username" :disabled="isEdit" placeholder="请输入用户名" />
      </n-form-item>
      <n-form-item v-if="!isEdit" label="密码" path="password">
        <n-input v-model:value="formData.password" type="password" placeholder="请输入密码" />
      </n-form-item>
      <n-form-item label="昵称" path="nickname">
        <n-input v-model:value="formData.nickname" placeholder="请输入昵称" />
      </n-form-item>
      <n-form-item label="邮箱" path="email">
        <n-input v-model:value="formData.email" placeholder="请输入邮箱" />
      </n-form-item>
      <n-form-item label="手机号" path="phone">
        <n-input v-model:value="formData.phone" placeholder="请输入手机号" />
      </n-form-item>
      <n-form-item label="角色" path="role_id">
        <n-select v-model:value="formData.role_id" :options="roleOptions" placeholder="请选择角色" />
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
</template>

<script setup lang="ts">
import { h, ref, onMounted, reactive } from 'vue';
import { NButton, NSpace, NTag, useMessage, useDialog } from 'naive-ui';
import { PersonAddOutline } from '@vicons/ionicons5';
import { getUserList, createUser, updateUser, deleteUser, resetPassword, getRoleOptions } from '@/api/system/user';

const message = useMessage();
const dialog = useDialog();

const loading = ref(false);
const submitLoading = ref(false);
const searchForm = ref({ username: '', phone: '', status: null as number | null });
const showModal = ref(false);
const isEdit = ref(false);
const formRef = ref();
const tableData = ref<any[]>([]);
const roleOptions = ref<any[]>([]);

const pagination = reactive({
  page: 1,
  pageSize: 10,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 50],
  onChange: (page: number) => {
    pagination.page = page;
    fetchData();
  },
  onUpdatePageSize: (pageSize: number) => {
    pagination.pageSize = pageSize;
    pagination.page = 1;
    fetchData();
  }
});

const formData = ref({
  ID: 0,
  username: '',
  password: '',
  nickname: '',
  email: '',
  phone: '',
  role_id: null as number | null,
  status: 1
});

const formRules = {
  username: { required: true, message: '请输入用户名', trigger: 'blur' },
  password: { required: true, message: '请输入密码', trigger: 'blur' },
  nickname: { required: true, message: '请输入昵称', trigger: 'blur' },
  role_id: { required: true, message: '请选择角色', trigger: 'change', type: 'number' }
};

const statusOptions = [
  { label: '正常', value: 1 },
  { label: '禁用', value: 0 }
];

const modalTitle = ref('新增用户');

const columns = [
  { title: 'ID', key: 'ID', width: 80 },
  { title: '用户名', key: 'username' },
  { title: '昵称', key: 'nickname' },
  { title: '邮箱', key: 'email' },
  { title: '手机号', key: 'phone' },
  {
    title: '角色',
    key: 'Role',
    render: (row: any) => row.Role?.name || '-'
  },
  {
    title: '状态',
    key: 'status',
    render: (row: any) => h(NTag, { type: row.status === 1 ? 'success' : 'error' }, { default: () => row.status === 1 ? '正常' : '禁用' })
  },
  {
    title: '操作',
    key: 'actions',
    width: 200,
    render: (row: any) => h(NSpace, null, {
      default: () => [
        h(NButton, { size: 'small', tertiary: true, type: 'info', onClick: () => handleEdit(row) }, { default: () => '编辑' }),
        h(NButton, { size: 'small', tertiary: true, type: 'warning', onClick: () => handleResetPwd(row) }, { default: () => '重置密码' }),
        h(NButton, { size: 'small', tertiary: true, type: 'error', disabled: row.ID === 1, onClick: () => handleDelete(row) }, { default: () => '删除' })
      ]
    })
  }
];

const fetchData = async () => {
  loading.value = true;
  try {
    const res: any = await getUserList({
      page: pagination.page,
      pageSize: pagination.pageSize,
      username: searchForm.value.username,
      phone: searchForm.value.phone,
      status: searchForm.value.status !== null ? searchForm.value.status : undefined
    });
    tableData.value = res.list || [];
    pagination.itemCount = res.total || 0;
  } catch (error) {
    console.error('Failed to fetch users:', error);
  } finally {
    loading.value = false;
  }
};

const fetchRoleOptions = async () => {
  try {
    const res: any = await getRoleOptions();
    roleOptions.value = res || [];
  } catch (error) {
    console.error('Failed to fetch role options:', error);
  }
};

const handleReset = () => {
  searchForm.value = { username: '', phone: '', status: null };
  pagination.page = 1;
  fetchData();
};

const handlePageChange = (page: number) => {
  pagination.page = page;
  fetchData();
};

const handleAdd = () => {
  isEdit.value = false;
  modalTitle.value = '新增用户';
  formData.value = { ID: 0, username: '', password: '', nickname: '', email: '', phone: '', role_id: null, status: 1 };
  showModal.value = true;
};

const handleEdit = (row: any) => {
  isEdit.value = true;
  modalTitle.value = '编辑用户';
  formData.value = { 
    ID: row.ID,
    username: row.username,
    password: '',
    nickname: row.nickname,
    email: row.email,
    phone: row.phone,
    role_id: row.role_id,
    status: row.status
  };
  showModal.value = true;
};

const handleResetPwd = (row: any) => {
  dialog.warning({
    title: '重置密码',
    content: `确定要重置用户 "${row.username}" 的密码吗？密码将被重置为 123456`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await resetPassword(row.ID);
        message.success('密码已重置为: 123456');
      } catch (error) {
        message.error('重置密码失败');
      }
    }
  });
};

const handleDelete = (row: any) => {
  dialog.error({
    title: '删除确认',
    content: `确定要删除用户 "${row.username}" 吗？此操作不可恢复！`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await deleteUser(row.ID);
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
          await updateUser(formData.value.ID, formData.value);
          message.success('更新成功');
        } else {
          await createUser(formData.value);
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
  fetchRoleOptions();
});
</script>
