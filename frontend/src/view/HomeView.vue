<template>
    <div class="layout-wrapper flex flex-col min-h-screen">
        <!-- 顶部导航栏 -->
        <div
            class="layout-topbar flex justify-between items-center bg-gradient-to-r from-primary-700 to-primary-800 text-white p-4 shadow-topbar">
            <div class="flex items-center">
                <button @click="toggleSidebar"
                    class="p-2 rounded-lg hover:bg-primary-600 transition-colors focus:outline-none focus:ring-2 focus:ring-primary-400 mr-2">
                    <i class="fas fa-bars text-xl aspect-square"></i>
                </button>
                <h1 class="text-xl font-bold">配置中心管理系统</h1>
            </div>

            <div class="flex items-center space-x-5">
                <div class="relative">
                    <select v-model="currentEnv"
                        class="bg-primary-600 border border-primary-500 rounded-lg py-2 pl-3 pr-8 text-white focus:outline-none focus:ring-2 focus:ring-primary-400 appearance-none">
                        <option v-for="env in environments" :key="env.value" :value="env.value">
                            {{ env.label }}
                        </option>
                    </select>
                    <i class="fas fa-chevron-down absolute right-3 top-3 pointer-events-none"></i>
                </div>

                <div
                    class="flex items-center space-x-2 cursor-pointer hover:bg-primary-600/20 p-2 rounded-lg transition-colors">
                    <div class="w-8 h-8 rounded-full bg-primary-500 flex items-center justify-center">
                        <i class="fas fa-user"></i>
                    </div>
                    <span class="font-medium">管理员</span>
                </div>

                <button @click="handleLogout"
                    class="flex items-center space-x-1 hover:bg-primary-600/20 p-2 rounded-lg transition-colors">
                    <i class="fas fa-sign-out-alt"></i>
                    <span>退出</span>
                </button>
            </div>
        </div>

        <div class="flex flex-1">
            <!-- 侧边栏 -->
            <div :class="{ 'w-20': !sidebarVisible, 'w-54': sidebarVisible, 'sidebar-collapsed': !sidebarVisible }"
                class="layout-sidebar flex flex-col bg-white h-full shadow-sidebar transition-all duration-300 ease-in-out">
                <div class="p-4 flex-1 overflow-y-auto">
                    <div class="space-y-1 mb-6">
                        <router-link to="/" class="menu-item flex items-center px-4 py-3 rounded-lg"
                            active-class="" exact-active-class="router-link-active">
                            <i class="fas fa-home text-lg"></i>
                            <span class="menu-text ml-3 transition-all duration-300">首页</span>
                        </router-link>
                        <router-link to="/apps" class="menu-item flex items-center px-4 py-3 rounded-lg"
                            active-class="router-link-active">
                            <i class="fas fa-th-large text-lg"></i>
                            <span class="menu-text ml-3 transition-all duration-300">应用管理</span>
                        </router-link>

                        <router-link to="/clusters" class="menu-item flex items-center px-4 py-3 rounded-lg"
                            active-class="router-link-active">
                            <i class="fas fa-server text-lg"></i>
                            <span class="menu-text ml-3 transition-all duration-300">集群管理</span>
                        </router-link>

                        <router-link to="/namespaces" class="menu-item flex items-center px-4 py-3 rounded-lg"
                            active-class="router-link-active">
                            <i class="fas fa-folder text-lg"></i>
                            <span class="menu-text ml-3 transition-all duration-300">命名空间</span>
                        </router-link>

                        <router-link to="/permissions" class="menu-item flex items-center px-4 py-3 rounded-lg"
                            active-class="router-link-active">
                            <i class="fas fa-shield-alt text-lg"></i>
                            <span class="menu-text ml-3 transition-all duration-300">权限管理</span>
                        </router-link>

                        <router-link to="/users" class="menu-item flex items-center px-4 py-3 rounded-lg"
                            active-class="router-link-active">
                            <i class="fas fa-users text-lg"></i>
                            <span class="menu-text ml-3 transition-all duration-300">用户管理</span>
                        </router-link>

                        <router-link to="/tokens" class="menu-item flex items-center px-4 py-3 rounded-lg"
                            active-class="router-link-active">
                            <i class="fas fa-key text-lg"></i>
                            <span class="menu-text ml-3 transition-all duration-300">令牌管理</span>
                        </router-link>
                    </div>
                    <hr class="h-px my-8 bg-gray-200 border-0 dark:bg-gray-700">
                    <div class="mt-8  pt-4">
                        <div class="text-xs text-gray-500 px-4 mb-2" :class="{ 'hidden': !sidebarVisible }">快捷导航</div>
                        <div class="space-y-1">
                            <a href="#" class="menu-item flex items-center px-4 py-3 rounded-lg">
                                <i class="fas fa-history text-lg"></i>
                                <span class="menu-text ml-3 transition-all duration-300">最近访问</span>
                            </a>
                            <a href="#" class="menu-item flex items-center px-4 py-3 rounded-lg">
                                <i class="fas fa-star text-lg"></i>
                                <span class="menu-text ml-3 transition-all duration-300">收藏应用</span>
                            </a>
                        </div>
                    </div>
                </div>
                <hr class="h-px my-8 bg-gray-200 border-0 dark:bg-gray-700 mx-4">
                <div class="p-4">
                    <div class="text-center">
                        <div class="text-xs text-gray-500 mb-1">系统版本 v2.5.1</div>
                        <div class="text-xs text-gray-500">© 2023 配置中心</div>
                    </div>
                </div>
            </div>

            <!-- 主内容区 -->
            <div class="layout-main flex-1 overflow-auto bg-gray-50">
                <div v-if="hasPermission" class="h-full w-full">
                    <router-view />
                </div>
                <DenieView v-else class="h-full w-full" />
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import { useToast } from 'primevue/usetoast';
import DenieView from '@/view/DenieView.vue'

const router = useRouter();
const toast = useToast();

// 响应式数据
const sidebarVisible = ref(true);
const currentEnv = ref('dev');
const hasPermission = computed(() => {
    return true;
});

// 环境选项
const environments = [
    { label: '开发环境', value: 'dev' },
    { label: '生产环境', value: 'pro' }
];



// 切换侧边栏
const toggleSidebar = () => {
    sidebarVisible.value = !sidebarVisible.value;
};

// 退出登录
const handleLogout = () => {
    toast.add({
        severity: 'success',
        summary: '退出成功',
        detail: '您已安全退出系统',
        life: 3000
    });

    // 模拟退出逻辑
    setTimeout(() => {
        router.push('/login');
    }, 1000);
};
</script>

<script>
export default { name: 'HomeView' }
</script>

<style scoped>
.layout-wrapper {
    font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    background-color: #f8fafc;
}

.layout-topbar {
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.layout-sidebar {
    box-shadow: 4px 0 15px rgba(0, 0, 0, 0.05);
}

.menu-item {
    transition: all 0.3s ease;
}

.menu-item:hover {
    background: linear-gradient(90deg, rgba(14, 165, 233, 0.1) 0%, rgba(14, 165, 233, 0) 100%);
    border-left: 3px solid #0ea5e9;
}

.router-link-active {
    background: linear-gradient(90deg, rgba(14, 165, 233, 0.15) 0%, rgba(14, 165, 233, 0) 100%);
    border-left: 3px solid #0ea5e9;
    color: #0ea5e9;
    font-weight: 500;
}

.stat-card,
.app-card {
    transition: all 0.3s ease;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05);
}

.stat-card:hover {
    transform: translateY(-3px);
}

.app-card:hover {
    transform: translateY(-3px);
    box-shadow: 0 10px 25px rgba(0, 0, 0, 0.1);
}

/* 侧边栏折叠时样式 */
:deep(.sidebar-collapsed) .menu-text {
    opacity: 0;
    width: 0;
    position: absolute;
    display: none;
}

:deep(.sidebar-collapsed) .menu-item {
    justify-content: center;
    padding-left: 0.75rem;
    padding-right: 0.75rem;
}

/* 过渡动画 */
.fade-enter-active,
.fade-leave-active {
    transition: opacity 0.3s;
}

.fade-enter,
.fade-leave-to {
    opacity: 0;
}
</style>