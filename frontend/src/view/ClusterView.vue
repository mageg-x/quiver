<template>
    <div class="cluster-management w-full p-6">
        <!-- 标题和操作按钮 -->
        <div class="flex justify-between items-center mb-6">
            <div>
                <h2 class="text-2xl font-bold text-gray-800">集群管理</h2>
                <p class="text-gray-600">管理{{ selectedApp ? selectedApp.name : '应用' }}的集群配置</p>
            </div>
            <button @click="showAddClusterDialog = true"
                class="bg-primary-600 hover:bg-primary-700 text-white px-4 py-2 rounded-lg flex items-center transition-colors">
                <i class="fas fa-plus mr-2"></i> 添加集群
            </button>
        </div>

        <!-- 应用选择器 -->
        <div class="bg-white rounded-xl shadow-sm p-4 mb-6 border border-gray-200">
            <div class="flex items-center space-x-4">
                <div class="w-full">
                    <label class="block text-sm font-medium text-gray-700 mb-2">选择应用</label>
                    <div class="relative">
                        <select v-model="selectedApp"
                            class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent">
                            <option v-for="app in apps" :key="app.id" :value="app">{{ app.name }}</option>
                        </select>
                        <i class="fas fa-chevron-down absolute right-3 top-3 text-gray-400 pointer-events-none"></i>
                    </div>
                </div>

                <div class="w-full">
                    <label class="block text-sm font-medium text-gray-700 mb-2">环境</label>
                    <div class="relative">
                        <select v-model="selectedEnv"
                            class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent">
                            <option v-for="env in environments" :key="env.value" :value="env.value">{{ env.label }}
                            </option>
                        </select>
                        <i class="fas fa-chevron-down absolute right-3 top-3 text-gray-400 pointer-events-none"></i>
                    </div>
                </div>
            </div>
        </div>

        <!-- 集群列表 -->
        <div class="bg-white rounded-xl shadow-sm overflow-hidden border border-gray-200">
            <div class="overflow-x-auto">
                <table class="min-w-full divide-y divide-gray-200">
                    <thead class="bg-gray-50">
                        <tr>
                            <th scope="col"
                                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                集群名称</th>
                            <th scope="col"
                                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                描述</th>
                            <th scope="col"
                                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                环境</th>
                            <th scope="col"
                                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                节点数量</th>
                            <th scope="col"
                                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                创建时间</th>
                            <th scope="col"
                                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                状态</th>
                            <th scope="col"
                                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                操作</th>
                        </tr>
                    </thead>
                    <tbody class="bg-white divide-y divide-gray-200">
                        <tr v-for="cluster in filteredClusters" :key="cluster.id"
                            class="hover:bg-gray-50 transition-colors">
                            <td class="px-6 py-4 whitespace-nowrap">
                                <router-link :to="`/clusters/${cluster.id}/namespaces`"
                                    class="text-primary-600 font-medium hover:underline">
                                    {{ cluster.name }}
                                </router-link>
                            </td>
                            <td class="px-6 py-4">
                                <div class="text-sm text-gray-900 max-w-xs truncate">{{ cluster.description }}</div>
                            </td>
                            <td class="px-6 py-4 whitespace-nowrap">
                                <span :class="envClass(cluster.env)" class="px-2 py-1 text-xs rounded-full">
                                    {{ envLabel(cluster.env) }}
                                </span>
                            </td>
                            <td class="px-6 py-4 whitespace-nowrap">
                                <div class="flex items-center">
                                    <span class="font-medium mr-1">{{ cluster.nodeCount }}</span>
                                    <span class="text-gray-500">个</span>
                                </div>
                            </td>
                            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                {{ formatDate(cluster.createdAt) }}
                            </td>
                            <td class="px-6 py-4 whitespace-nowrap">
                                <span :class="statusClass(cluster.status)" class="px-2 py-1 text-xs rounded-full">
                                    {{ statusLabels[cluster.status] }}
                                </span>
                            </td>
                            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                <div class="flex space-x-3">
                                    <button @click="editCluster(cluster)"
                                        class="action-btn text-blue-600 hover:text-blue-800">
                                        <i class="fas fa-edit"></i>
                                    </button>
                                    <button @click="confirmDelete(cluster)"
                                        class="action-btn text-red-600 hover:text-red-800">
                                        <i class="fas fa-trash"></i>
                                    </button>
                                </div>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>

            <!-- 空状态 -->
            <div v-if="filteredClusters.length === 0" class="text-center py-12">
                <div class="mx-auto w-24 h-24 rounded-full bg-gray-100 flex items-center justify-center mb-4">
                    <i class="fas fa-server text-gray-400 text-3xl"></i>
                </div>
                <h3 class="text-lg font-medium text-gray-900 mb-1">暂无集群数据</h3>
                <p class="text-gray-500 max-w-md mx-auto">
                    当前应用下还没有创建集群，点击"添加集群"按钮创建您的第一个集群
                </p>
                <button @click="showAddClusterDialog = true"
                    class="mt-4 bg-primary-600 hover:bg-primary-700 text-white px-4 py-2 rounded-lg inline-flex items-center transition-colors">
                    <i class="fas fa-plus mr-2"></i> 添加集群
                </button>
            </div>

            <!-- 分页控件 -->
            <div v-if="filteredClusters.length > 0"
                class="bg-white px-4 py-3 flex items-center justify-between border-t border-gray-200 sm:px-6">
                <div class="hidden sm:flex-1 sm:flex sm:items-center sm:justify-between">
                    <div>
                        <p class="text-sm text-gray-700">
                            显示第 <span class="font-medium">1</span> 到 <span class="font-medium">{{
                                filteredClusters.length }}</span> 条，共 <span class="font-medium">{{
                                filteredClusters.length }}</span> 条记录
                        </p>
                    </div>
                    <div>
                        <nav class="relative z-0 inline-flex rounded-md shadow-sm -space-x-px" aria-label="Pagination">
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
                    <i class="fas fa-server text-blue-600 text-2xl"></i>
                </div>
                <div>
                    <div class="text-gray-500 text-sm">集群总数</div>
                    <div class="text-2xl font-bold text-gray-800">{{ filteredClusters.length }}</div>
                </div>
            </div>

            <div class="stat-card bg-white rounded-xl shadow-sm p-6 flex items-center border border-gray-200">
                <div class="bg-green-100 p-4 rounded-lg mr-4">
                    <i class="fas fa-check-circle text-green-600 text-2xl"></i>
                </div>
                <div>
                    <div class="text-gray-500 text-sm">运行中集群</div>
                    <div class="text-2xl font-bold text-gray-800">{{ runningClustersCount }}</div>
                </div>
            </div>

            <div class="stat-card bg-white rounded-xl shadow-sm p-6 flex items-center border border-gray-200">
                <div class="bg-amber-100 p-4 rounded-lg mr-4">
                    <i class="fas fa-microchip text-amber-600 text-2xl"></i>
                </div>
                <div>
                    <div class="text-gray-500 text-sm">总节点数</div>
                    <div class="text-2xl font-bold text-gray-800">{{ totalNodesCount }}</div>
                </div>
            </div>
        </div>

        <!-- 添加/编辑集群对话框 -->
        <transition name="fade">
            <div v-if="showAddClusterDialog" class="fixed inset-0 flex items-center justify-center z-50 p-4">
                <div class="dialog-overlay fixed inset-0 bg-black bg-opacity-50" @click="closeDialog"></div>
                <div class="bg-white rounded-xl shadow-xl w-full max-w-md relative z-10 transform transition-all">
                    <div class="p-6">
                        <div class="flex justify-between items-center mb-4">
                            <h3 class="text-lg font-bold text-gray-800">{{ editingCluster ? '编辑集群' : '添加新集群' }}</h3>
                            <button @click="closeDialog" class="text-gray-400 hover:text-gray-500">
                                <i class="fas fa-times"></i>
                            </button>
                        </div>

                        <div class="space-y-4">
                            <div>
                                <label class="block text-sm font-medium text-gray-700 mb-1">集群名称</label>
                                <input v-model="currentCluster.name" type="text" placeholder="如：北京生产集群"
                                    class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent">
                            </div>

                            <div>
                                <label class="block text-sm font-medium text-gray-700 mb-1">描述</label>
                                <textarea v-model="currentCluster.description" rows="3" placeholder="集群用途描述..."
                                    class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"></textarea>
                            </div>

                            <div class="grid grid-cols-2 gap-4">
                                <div>
                                    <label class="block text-sm font-medium text-gray-700 mb-1">所属应用</label>
                                    <select v-model="currentCluster.appId"
                                        class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent">
                                        <option v-for="app in apps" :key="app.id" :value="app.id">{{ app.name }}
                                        </option>
                                    </select>
                                </div>

                                <div>
                                    <label class="block text-sm font-medium text-gray-700 mb-1">环境</label>
                                    <select v-model="currentCluster.env"
                                        class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent">
                                        <option v-for="env in environments" :key="env.value" :value="env.value">{{
                                            env.label }}</option>
                                    </select>
                                </div>
                            </div>

                            <div class="grid grid-cols-2 gap-4">
                                <div>
                                    <label class="block text-sm font-medium text-gray-700 mb-1">节点数量</label>
                                    <input v-model="currentCluster.nodeCount" type="number" min="1"
                                        class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent">
                                </div>

                                <div>
                                    <label class="block text-sm font-medium text-gray-700 mb-1">状态</label>
                                    <div class="grid grid-cols-2 gap-2">
                                        <button @click="currentCluster.status = 'running'"
                                            :class="{ 'bg-green-100 border-green-500 text-green-700': currentCluster.status === 'running', 'bg-gray-100 border-gray-300 text-gray-700': currentCluster.status !== 'running' }"
                                            class="py-2 px-3 border rounded-lg text-sm transition-colors">
                                            运行中
                                        </button>
                                        <button @click="currentCluster.status = 'stopped'"
                                            :class="{ 'bg-red-100 border-red-500 text-red-700': currentCluster.status === 'stopped', 'bg-gray-100 border-gray-300 text-gray-700': currentCluster.status !== 'stopped' }"
                                            class="py-2 px-3 border rounded-lg text-sm transition-colors">
                                            已停止
                                        </button>
                                    </div>
                                </div>
                            </div>
                        </div>

                        <div class="mt-6 flex justify-end space-x-3">
                            <button @click="closeDialog"
                                class="px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors">取消</button>
                            <button @click="saveCluster"
                                class="px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 transition-colors">
                                {{ editingCluster ? '保存更改' : '添加集群' }}
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        </transition>

        <!-- 删除确认对话框 -->
        <transition name="fade">
            <div v-if="showDeleteDialog" class="fixed inset-0 flex items-center justify-center z-50 p-4">
                <div class="dialog-overlay fixed inset-0 bg-black bg-opacity-50" @click="showDeleteDialog = false">
                </div>
                <div class="bg-white rounded-xl shadow-xl w-full max-w-md relative z-10 transform transition-all">
                    <div class="p-6">
                        <div class="flex justify-center mb-4">
                            <div class="w-16 h-16 rounded-full bg-red-100 flex items-center justify-center">
                                <i class="fas fa-exclamation-triangle text-red-600 text-2xl"></i>
                            </div>
                        </div>

                        <h3 class="text-lg font-bold text-center text-gray-800 mb-2">确认删除集群</h3>
                        <p class="text-gray-600 text-center mb-6">确定要删除集群 <span class="font-semibold">{{
                                currentCluster.name }}</span> 吗？此操作不可恢复。</p>

                        <div class="flex justify-center space-x-3">
                            <button @click="showDeleteDialog = false"
                                class="px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors">取消</button>
                            <button @click="deleteCluster"
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
import { ref, computed, onMounted } from 'vue';

// 应用数据
const apps = ref([
    { id: 1, name: '订单服务' },
    { id: 2, name: '支付服务' },
    { id: 3, name: '用户服务' },
    { id: 4, name: '库存服务' },
]);

// 环境选项
const environments = ref([
    { label: '开发环境', value: 'dev' },
    { label: '测试环境', value: 'test' },
    { label: '预发环境', value: 'stage' },
    { label: '生产环境', value: 'prod' },
]);

// 集群数据
const clusters = ref([
    {
        id: 1,
        appId: 1,
        name: '北京生产集群',
        description: '处理华北地区订单的主要生产集群',
        env: 'prod',
        nodeCount: 12,
        status: 'running',
        createdAt: '2023-01-15'
    },
    {
        id: 2,
        appId: 1,
        name: '上海测试集群',
        description: '用于新功能测试的集群',
        env: 'test',
        nodeCount: 4,
        status: 'running',
        createdAt: '2023-02-20'
    },
    {
        id: 3,
        appId: 2,
        name: '支付生产集群',
        description: '处理所有支付交易的生产集群',
        env: 'prod',
        nodeCount: 8,
        status: 'running',
        createdAt: '2023-03-10'
    },
    {
        id: 4,
        appId: 2,
        name: '支付开发集群',
        description: '支付服务的开发测试环境',
        env: 'dev',
        nodeCount: 3,
        status: 'stopped',
        createdAt: '2023-04-05'
    },
    {
        id: 5,
        appId: 3,
        name: '用户服务集群',
        description: '用户认证和管理的主集群',
        env: 'prod',
        nodeCount: 6,
        status: 'running',
        createdAt: '2023-05-18'
    },
]);

// 状态标签
const statusLabels = ref({
    running: '运行中',
    stopped: '已停止',
    pending: '部署中',
    error: '异常'
});

// 环境标签
const envLabel = (env) => {
    const envObj = environments.value.find(e => e.value === env);
    return envObj ? envObj.label : env;
};

// 状态类
const statusClass = (status) => {
    return {
        'bg-green-100 text-green-800': status === 'running',
        'bg-red-100 text-red-800': status === 'stopped',
        'bg-yellow-100 text-yellow-800': status === 'pending',
        'bg-gray-100 text-gray-800': status === 'error',
    };
};

// 环境类
const envClass = (env) => {
    return {
        'bg-blue-100 text-blue-800': env === 'dev',
        'bg-green-100 text-green-800': env === 'test',
        'bg-amber-100 text-amber-800': env === 'stage',
        'bg-purple-100 text-purple-800': env === 'prod',
    };
};

// 选中的应用
const selectedApp = ref(apps.value[0]);
// 选中的环境
const selectedEnv = ref('');

// 过滤后的集群列表
const filteredClusters = computed(() => {
    let result = clusters.value.filter(cluster => cluster.appId === selectedApp.value.id);

    if (selectedEnv.value) {
        result = result.filter(cluster => cluster.env === selectedEnv.value);
    }

    return result;
});

// 运行中集群数量
const runningClustersCount = computed(() => {
    return filteredClusters.value.filter(cluster => cluster.status === 'running').length;
});

// 总节点数
const totalNodesCount = computed(() => {
    return filteredClusters.value.reduce((sum, cluster) => sum + cluster.nodeCount, 0);
});

// 对话框状态
const showAddClusterDialog = ref(false);
const showDeleteDialog = ref(false);
const editingCluster = ref(false);

// 当前操作的集群
const currentCluster = ref({
    id: null,
    appId: selectedApp.value.id,
    name: '',
    description: '',
    env: 'prod',
    nodeCount: 1,
    status: 'running'
});

// 提示信息
const showToast = ref(false);
const toastMessage = ref('');

// 日期格式化
const formatDate = (dateString) => {
    const options = { year: 'numeric', month: 'long', day: 'numeric' };
    return new Date(dateString).toLocaleDateString('zh-CN', options);
};

// 编辑集群
const editCluster = (cluster) => {
    currentCluster.value = { ...cluster };
    editingCluster.value = true;
    showAddClusterDialog.value = true;
};

// 确认删除
const confirmDelete = (cluster) => {
    currentCluster.value = { ...cluster };
    showDeleteDialog.value = true;
};

// 删除集群
const deleteCluster = () => {
    clusters.value = clusters.value.filter(c => c.id !== currentCluster.value.id);
    showDeleteDialog.value = false;
    showToastMessage('集群已成功删除');
};

// 保存集群
const saveCluster = () => {
    if (editingCluster.value) {
        // 更新现有集群
        const index = clusters.value.findIndex(c => c.id === currentCluster.value.id);
        if (index !== -1) {
            clusters.value[index] = { ...currentCluster.value };
        }
        showToastMessage('集群信息已更新');
    } else {
        // 添加新集群
        const newId = Math.max(...clusters.value.map(c => c.id), 0) + 1;
        clusters.value.push({
            id: newId,
            ...currentCluster.value,
            createdAt: new Date().toISOString().split('T')[0]
        });
        showToastMessage('新集群已成功添加');
    }
    closeDialog();
};

// 关闭对话框
const closeDialog = () => {
    showAddClusterDialog.value = false;
    showDeleteDialog.value = false;
    editingCluster.value = false;
    currentCluster.value = {
        id: null,
        appId: selectedApp.value.id,
        name: '',
        description: '',
        env: 'prod',
        nodeCount: 1,
        status: 'running'
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

// 初始化
onMounted(() => {
    // 模拟从API加载数据
    setTimeout(() => {
        // 这里可以添加实际的数据加载逻辑
    }, 500);
});
</script>


<script>
export default { name: 'ClusterView' }
</script>

<style scoped>
.cluster-management {
    font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}

.stat-card {
    transition: all 0.3s ease;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05);
}

.stat-card:hover {
    transform: translateY(-3px);
}

.action-btn {
    transition: all 0.2s ease;
}

.action-btn:hover {
    transform: scale(1.1);
}

.fade-enter-active,
.fade-leave-active {
    transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
    opacity: 0;
}

.dialog-overlay {
    background: rgba(0, 0, 0, 0.5);
}
</style>