import axios, { type AxiosInstance, type AxiosResponse } from 'axios';
import { createDiscreteApi } from 'naive-ui';

const { message } = createDiscreteApi(['message']);

const service: AxiosInstance = axios.create({
    baseURL: import.meta.env.VITE_BASE_URL || '/api',
    timeout: 10000,
});

// Request Interceptor
service.interceptors.request.use(
    (config: any) => {
        const token = localStorage.getItem('token');
        if (token) {
            config.headers = {
                ...config.headers,
                'Authorization': `Bearer ${token}`
            }
        }
        return config;
    },
    (error: any) => {
        return Promise.reject(error);
    }
);

// Response Interceptor
service.interceptors.response.use(
    (response: AxiosResponse) => {
        const res = response.data;
        // Assuming backend returns { code: 0, data: ..., msg: ... }
        if (res.code !== 0) {
            message.error(res.msg || 'Error');
            // 如果是 401，跳转到登录页
            if (res.code === 401) {
                localStorage.removeItem('token');
                window.location.href = '/login';
            }
            return Promise.reject(new Error(res.msg || 'Error'));
        } else {
            return res.data;
        }
    },
    (error: any) => {
        console.log('err' + error);
        // 处理 HTTP 401 错误
        if (error.response?.status === 401) {
            message.error('登录已过期，请重新登录');
            localStorage.removeItem('token');
            window.location.href = '/login';
        } else {
            message.error(error.message);
        }
        return Promise.reject(error);
    }
);

export default service;
