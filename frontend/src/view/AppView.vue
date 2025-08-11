<template>
    <div class="app-management w-full">
        <div class="p-6">
            <!-- 标题和操作按钮 -->
            <div class="flex justify-between items-center mb-6">
                <div>
                    <h2 class="text-2xl font-bold text-gray-800">应用管理</h2>
                    <p class="text-gray-600">管理您的所有应用配置</p>
                </div>
                <button @click="showAddDialog = true"
                    class="bg-primary-600 hover:bg-primary-700 text-white px-4 py-2 rounded-lg flex items-center transition-colors">
                    <i class="fas fa-plus mr-2"></i> 添加应用
                </button>
            </div>

            <!-- 搜索和筛选区域 -->
            <div class="bg-white rounded-xl shadow-sm p-4 mb-6 border border-gray-200">
                <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
                    <div class="relative">
                        <input v-model="searchQuery" type="text" placeholder="搜索应用名称..."
                            class="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent">
                        <i class="fas fa-search absolute left-3 top-3 text-gray-400"></i>
                    </div>

                    <div>
                        <select v-model="selectedOwner"
                            class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent">
                            <option value="">所有负责人</option>
                            <option v-for="owner in owners" :value="owner">{{ owner }}</option>
                        </select>
                    </div>

                    <div>
                        <button @click="resetFilters"
                            class="w-full bg-gray-100 hover:bg-gray-200 text-gray-800 px-4 py-2 rounded-lg transition-colors">
                            重置筛选
                        </button>
                    </div>
                </div>
            </div>

            <!-- 应用列表表格 -->
            <div class="bg-white rounded-xl shadow-sm overflow-hidden border border-gray-200">
                <div class="overflow-x-auto">
                    <table class="min-w-full divide-y divide-gray-200">
                        <thead class="bg-gray-50">
                            <tr>
                                <th scope="col"
                                    class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    应用名称</th>
                                <th scope="col"
                                    class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    描述</th>
                                <th scope="col"
                                    class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    负责人</th>
                                <th scope="col"
                                    class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    创建时间</th>
                                <th scope="col"
                                    class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    操作</th>
                            </tr>
                        </thead>
                        <tbody class="bg-white divide-y divide-gray-200">
                            <tr v-for="app in filteredApps" :key="app.id"
                                class="app-table-row hover:bg-gray-50 transition-colors">
                                <td class="px-6 py-4 whitespace-nowrap">
                                    <a :href="'/apps/' + app.id" class="text-primary-600 font-medium hover:underline">{{
                                        app.name }}</a>
                                </td>
                                <td class="px-6 py-4">
                                    <div class="text-sm text-gray-900 max-w-xs truncate">{{ app.description }}</div>
                                </td>
                                <td class="px-6 py-4 whitespace-nowrap">
                                    <div class="flex items-center">
                                        <div
                                            class="w-8 h-8 rounded-full bg-blue-100 flex items-center justify-center mr-2">
                                            <i class="fas fa-user text-blue-600 text-sm"></i>
                                        </div>
                                        <span>{{ app.owner }}</span>
                                    </div>
                                </td>
                                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                    {{ formatDate(app.createdAt) }}
                                </td>
                                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                    <div class="flex space-x-3">
                                        <button @click="editApp(app)"
                                            class="action-btn text-blue-600 hover:text-blue-800">
                                            <i class="fas fa-edit"></i>
                                        </button>
                                        <button @click="confirmDelete(app)"
                                            class="action-btn text-red-600 hover:text-red-800">
                                            <i class="fas fa-trash"></i>
                                        </button>
                                    </div>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>

                <!-- 分页控件 -->
                <div class="bg-white px-4 py-3 flex items-center justify-between border-t border-gray-200 sm:px-6">
                    <div class="hidden sm:flex-1 sm:flex sm:items-center sm:justify-between">
                        <div>
                            <p class="text-sm text-gray-700">
                                显示第 <span class="font-medium">1</span> 到 <span class="font-medium">8</span> 条，共 <span
                                    class="font-medium">{{ apps.length }}</span> 条记录
                            </p>
                        </div>
                        <div>
                            <nav class="relative z-0 inline-flex rounded-md shadow-sm -space-x-px"
                                aria-label="Pagination">
                                <button
                                    class="relative inline-flex items-center px-2 py-2 rounded-l-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50">
                                    <span class="sr-only">上一页</span>
                                    <i class="fas fa-chevron-left"></i>
                                </button>
                                <button aria-current="page"
                                    class="z-10 bg-primary-50 border-primary-500 text-primary-600 relative inline-flex items-center px-4 py-2 border text-sm font-medium">
                                    1
                                </button>
                                <button
                                    class="relative inline-flex items-center px-4 py-2 border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50">
                                    2
                                </button>
                                <button
                                    class="relative inline-flex items-center px-4 py-2 border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50">
                                    3
                                </button>
                                <button
                                    class="relative inline-flex items-center px-2 py-2 rounded-r-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50">
                                    <span class="sr-only">下一页</span>
                                    <i class="fas fa-chevron-right"></i>
                                </button>
                            </nav>
                        </div>
                    </div>
                </div>
            </div>

            <!-- 统计卡片 -->
            <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mt-8">
                <div class="stat-card bg-white rounded-xl shadow-sm p-6 flex items-center border border-gray-200">
                    <div class="bg-blue-100 p-4 rounded-lg mr-4">
                        <i class="fas fa-th-large text-blue-600 text-2xl"></i>
                    </div>
                    <div>
                        <div class="text-gray-500 text-sm">应用总数</div>
                        <div class="text-2xl font-bold text-gray-800">{{ apps.length }}</div>
                    </div>
                </div>

                <div class="stat-card bg-white rounded-xl shadow-sm p-6 flex items-center border border-gray-200">
                    <div class="bg-green-100 p-4 rounded-lg mr-4">
                        <i class="fas fa-check-circle text-green-600 text-2xl"></i>
                    </div>
                    <div>
                        <div class="text-gray-500 text-sm">活跃应用</div>
                        <div class="text-2xl font-bold text-gray-800">{{ activeAppsCount }}</div>
                    </div>
                </div>

                <div class="stat-card bg-white rounded-xl shadow-sm p-6 flex items-center border border-gray-200">
                    <div class="bg-amber-100 p-4 rounded-lg mr-4">
                        <i class="fas fa-users text-amber-600 text-2xl"></i>
                    </div>
                    <div>
                        <div class="text-gray-500 text-sm">负责人数量</div>
                        <div class="text-2xl font-bold text-gray-800">{{ owners.length }}</div>
                    </div>
                </div>
            </div>
        </div>

        <!-- 添加/编辑应用对话框 -->
        <transition name="fade">
            <div v-if="showAddDialog" class="fixed inset-0 flex items-center justify-center z-50 p-4">
                <div class="dialog-overlay fixed inset-0" @click="closeDialog"></div>
                <div class="bg-white rounded-xl shadow-xl w-full max-w-md relative z-10 transform transition-all">
                    <div class="p-6">
                        <div class="flex justify-between items-center mb-4">
                            <h3 class="text-lg font-bold text-gray-800">{{ editingApp ? '编辑应用' : '添加新应用' }}</h3>
                            <button @click="closeDialog" class="text-gray-400 hover:text-gray-500">
                                <i class="fas fa-times"></i>
                            </button>
                        </div>

                        <div class="space-y-4">
                            <div>
                                <label class="block text-sm font-medium text-gray-700 mb-1">应用名称</label>
                                <input v-model="currentApp.name" type="text"
                                    class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent">
                            </div>

                            <div>
                                <label class="block text-sm font-medium text-gray-700 mb-1">描述</label>
                                <textarea v-model="currentApp.description" rows="3"
                                    class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"></textarea>
                            </div>

                            <div>
                                <label class="block text-sm font-medium text-gray-700 mb-1">负责人</label>
                                <select v-model="currentApp.owner"
                                    class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent">
                                    <option v-for="owner in owners" :value="owner">{{ owner }}</option>
                                </select>
                            </div>

                        </div>

                        <div class="mt-6 flex justify-end space-x-3">
                            <button @click="closeDialog"
                                class="px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors">取消</button>
                            <button @click="saveApp"
                                class="px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 transition-colors">{{
                                    editingApp ? '保存更改' : '添加应用' }}</button>
                        </div>
                    </div>
                </div>
            </div>
        </transition>

        <!-- 删除确认对话框 -->
        <transition name="fade">
            <div v-if="showDeleteDialog" class="fixed inset-0 flex items-center justify-center z-50 p-4">
                <div class="dialog-overlay fixed inset-0" @click="showDeleteDialog = false"></div>
                <div class="bg-white rounded-xl shadow-xl w-full max-w-md relative z-10 transform transition-all">
                    <div class="p-6">
                        <div class="flex justify-center mb-4">
                            <div class="w-16 h-16 rounded-full bg-red-100 flex items-center justify-center">
                                <i class="fas fa-exclamation-triangle text-red-600 text-2xl"></i>
                            </div>
                        </div>

                        <h3 class="text-lg font-bold text-center text-gray-800 mb-2">确认删除应用</h3>
                        <p class="text-gray-600 text-center mb-6">确定要删除应用 <span class="font-semibold">{{ currentApp.name
                        }}</span> 吗？此操作不可恢复。</p>

                        <div class="flex justify-center space-x-3">
                            <button @click="showDeleteDialog = false"
                                class="px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors">取消</button>
                            <button @click="deleteApp"
                                class="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors">确认删除</button>
                        </div>
                    </div>
                </div>
            </div>
        </transition>

        <!-- 操作成功提示 -->
        <transition name="fade">
            <div v-if="showToast" class="fixed top-4 right-4 z-50">
                <div
                    class="bg-green-100 border border-green-400 text-green-700 px-4 py-3 rounded-lg shadow-lg flex items-center">
                    <i class="fas fa-check-circle mr-2"></i>
                    <span>{{ toastMessage }}</span>
                    <button @click="showToast = false" class="ml-4 text-green-700 hover:text-green-900">
                        <i class="fas fa-times"></i>
                    </button>
                </div>
            </div>
        </transition>
    </div>
</template>

<script setup>
import { ref, computed } from 'vue'

// 应用数据
const apps = ref([
    { id: 1, name: '订单服务', description: '处理电商平台订单流程的核心服务', owner: '张工', createdAt: '2023-01-15' },
    { id: 2, name: '支付服务', description: '负责处理所有支付交易', owner: '李工', createdAt: '2023-02-10' },
    { id: 3, name: '用户服务', description: '管理用户账户和认证信息', owner: '王工', createdAt: '2023-03-22' },
    { id: 4, name: '库存服务', description: '实时跟踪商品库存状态', owner: '赵工', createdAt: '2023-04-05' },
    { id: 5, name: '推荐服务', description: '提供个性化商品推荐', owner: '陈工', createdAt: '2023-05-18' },
    { id: 6, name: '搜索服务', description: '提供商品搜索功能', owner: '刘工', createdAt: '2023-06-30' },
    { id: 7, name: '物流服务', description: '管理订单配送和物流信息', owner: '孙工', createdAt: '2023-07-12' },
    { id: 8, name: '消息服务', description: '处理系统内消息通知', owner: '钱工', createdAt: '2023-08-25' },
]);




// 负责人列表
const owners = computed(() => {
    const uniqueOwners = new Set(apps.value.map(app => app.owner));
    return Array.from(uniqueOwners);
});

// 活跃应用数量
const activeAppsCount = computed(() => {
    return 3
});

// 搜索和筛选
const searchQuery = ref('');
const selectedOwner = ref('');

// 过滤后的应用列表
const filteredApps = computed(() => {
    return apps.value.filter(app => {
        const matchesSearch = app.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
            app.description.toLowerCase().includes(searchQuery.value.toLowerCase());

        const matchesOwner = selectedOwner.value ? app.owner === selectedOwner.value : true;

        return matchesSearch && matchesOwner;
    });
});

// 对话框状态
const showAddDialog = ref(false);
const showDeleteDialog = ref(false);
const editingApp = ref(false);

// 当前操作的应用
const currentApp = ref({
    id: null,
    name: '',
    description: '',
    owner: '',
});

// 提示信息
const showToast = ref(false);
const toastMessage = ref('');

// 日期格式化
const formatDate = (dateString) => {
    const options = { year: 'numeric', month: 'long', day: 'numeric' };
    return new Date(dateString).toLocaleDateString('zh-CN', options);
};

// 重置筛选
const resetFilters = () => {
    searchQuery.value = '';
    selectedOwner.value = '';
};

// 编辑应用
const editApp = (app) => {
    currentApp.value = { ...app };
    editingApp.value = true;
    showAddDialog.value = true;
};

// 确认删除
const confirmDelete = (app) => {
    currentApp.value = { ...app };
    showDeleteDialog.value = true;
};

// 删除应用
const deleteApp = () => {
    apps.value = apps.value.filter(a => a.id !== currentApp.value.id);
    showDeleteDialog.value = false;
    showToastMessage('应用已成功删除');
};

// 保存应用
const saveApp = () => {
    if (editingApp.value) {
        // 更新现有应用
        const index = apps.value.findIndex(a => a.id === currentApp.value.id);
        if (index !== -1) {
            apps.value[index] = { ...currentApp.value };
        }
        showToastMessage('应用信息已更新');
    } else {
        // 添加新应用
        const newId = Math.max(...apps.value.map(a => a.id), 0) + 1;
        apps.value.push({
            id: newId,
            ...currentApp.value,
            createdAt: new Date().toISOString().split('T')[0]
        });
        showToastMessage('新应用已成功添加');
    }
    closeDialog();
};

// 关闭对话框
const closeDialog = () => {
    showAddDialog.value = false;
    showDeleteDialog.value = false;
    editingApp.value = false;
    currentApp.value = {
        id: null,
        name: '',
        description: '',
        owner: '',
    };
};

// 显示提示信息
const showToastMessage = (message) => {
    toastMessage.value = message;
    showToast.value = true;
    setTimeout(() => {
        showToast.value = false;
    }, 3000);
};
</script>


<script>
export default { name: 'AppView' }
</script>