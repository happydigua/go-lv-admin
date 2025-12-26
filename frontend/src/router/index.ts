import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router';

const routes: Array<RouteRecordRaw> = [
    {
        path: '/login',
        name: 'Login',
        component: () => import('@/views/login/index.vue'),
        meta: { title: 'Login' }
    },
    {
        path: '/',
        component: () => import('@/layouts/AdminLayout.vue'),
        redirect: '/dashboard',
        meta: { requiresAuth: true },
        children: [
            {
                path: 'dashboard',
                name: 'Dashboard',
                component: () => import('@/views/dashboard/index.vue'),
                meta: { title: '仪表盘', requiresAuth: true }
            },
            {
                path: 'profile',
                name: 'Profile',
                component: () => import('@/views/profile/index.vue'),
                meta: { title: '个人中心', requiresAuth: true }
            },
            {
                path: 'system/user',
                name: 'SystemUser',
                component: () => import('@/views/system/user/index.vue'),
                meta: { title: '用户管理', requiresAuth: true }
            },
            {
                path: 'system/role',
                name: 'SystemRole',
                component: () => import('@/views/system/role/index.vue'),
                meta: { title: '角色管理', requiresAuth: true }
            },
            {
                path: 'system/menu',
                name: 'SystemMenu',
                component: () => import('@/views/system/menu/index.vue'),
                meta: { title: '菜单管理', requiresAuth: true }
            },
            {
                path: 'system/log',
                name: 'SystemLog',
                component: () => import('@/views/system/log/index.vue'),
                meta: { title: '操作日志', requiresAuth: true }
            },
            {
                path: 'system/file',
                name: 'SystemFile',
                component: () => import('@/views/system/file/index.vue'),
                meta: { title: '文件管理', requiresAuth: true }
            },
            {
                path: 'system/setting',
                name: 'SystemSetting',
                component: () => import('@/views/system/setting/index.vue'),
                meta: { title: '系统设置', requiresAuth: true }
            },
            {
                path: 'tool/generator',
                name: 'ToolGenerator',
                component: () => import('@/views/tool/generator/index.vue'),
                meta: { title: '代码生成', requiresAuth: true }
            },
            {
                path: 'demo/table',
                name: 'DemoTable',
                component: () => import('@/views/demo/table/index.vue'),
                meta: { title: '复杂表格', requiresAuth: true }
            }
        ]
    }
];


const router = createRouter({
    history: createWebHistory(),
    routes
});

router.beforeEach((to, _from, next) => {
    const token = localStorage.getItem('token');
    if (to.meta.requiresAuth && !token) {
        next('/login');
    } else {
        next();
    }
});

export default router;
