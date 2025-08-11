<template>
    <div class="namespace-view p-6">
        <!-- 操作工具栏 -->
        <div class="bg-white rounded-xl shadow-sm p-5 mb-6">
            <div class="flex flex-wrap items-center justify-between">
                <div>
                    <h2 class="text-2xl font-bold text-gray-800 flex items-center">
                        <i class="fas fa-folder text-primary-600 mr-3"></i>
                        命名空间管理
                        <span class="ml-3 text-base font-normal bg-primary-100 text-primary-800 px-3 py-1 rounded-full">
                            当前命名空间: {{ namespaceName }}
                        </span>
                    </h2>
                    <div class="mt-2 text-sm text-gray-600 flex items-center">
                        <span class="font-medium">最新版本:</span>
                        <span class="ml-2 bg-gray-100 px-2 py-1 rounded font-mono">release-{{ currentReleaseId }}</span>
                        <span class="ml-4 font-medium">最后发布时间:</span>
                        <span class="ml-2">{{ lastPublishTime }}</span>
                    </div>
                </div>

                <div class="flex flex-wrap gap-3 mt-4 sm:mt-0">
                    <button class="btn bg-primary-500 hover:bg-primary-600 text-white">
                        <i class="fas fa-plus mr-2"></i>新增配置项
                    </button>
                    <button class="btn bg-green-500 hover:bg-green-600 text-white">
                        <i class="fas fa-paper-plane mr-2"></i>发布版本
                    </button>
                    <button class="btn bg-blue-500 hover:bg-blue-600 text-white">
                        <i class="fas fa-history mr-2"></i>发布历史
                    </button>
                    <button class="btn bg-orange-500 hover:bg-orange-600 text-white">
                        <i class="fas fa-undo mr-2"></i>回滚
                    </button>
                    <button class="btn bg-purple-500 hover:bg-purple-600 text-white">
                        <i class="fas fa-adjust mr-2"></i>灰度发布
                    </button>
                </div>
            </div>
        </div>

        <!-- 搜索和过滤区域 -->
        <div class="bg-white rounded-xl shadow-sm p-5 mb-6">
            <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
                <div class="md:col-span-2">
                    <div class="relative">
                        <input type="text" placeholder="搜索配置项名称或值..."
                            class="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-400 focus:border-primary-500">
                        <i class="fas fa-search absolute left-3 top-3 text-gray-400"></i>
                    </div>
                </div>

                <div>
                    <select
                        class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-400 focus:border-primary-500">
                        <option value="">所有类型</option>
                        <option v-for="type in types" :key="type" :value="type">{{ type }}</option>
                    </select>
                </div>

                <div>
                    <select
                        class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-400 focus:border-primary-500">
                        <option value="">所有状态</option>
                        <option v-for="status in statuses" :key="status" :value="status">{{ status }}</option>
                    </select>
                </div>
            </div>

            <div class="flex flex-wrap gap-3 mt-4">
                <button class="flex items-center text-sm text-gray-600 hover:text-primary-600">
                    <i class="fas fa-filter mr-1"></i> 高级筛选
                </button>
                <button class="flex items-center text-sm text-gray-600 hover:text-primary-600">
                    <i class="fas fa-sync-alt mr-1"></i> 重置筛选
                </button>
                <button class="flex items-center text-sm text-gray-600 hover:text-primary-600 ml-auto">
                    <i class="fas fa-columns mr-1"></i> 自定义列
                </button>
            </div>
        </div>

        <!-- 配置项表格 -->
        <div class="bg-white rounded-xl shadow-sm overflow-hidden">
            <div class="overflow-x-auto">
                <table class="min-w-full">
                    <thead class="bg-gray-50">
                        <tr>
                            <th class="py-3 px-6 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                发布状态
                            </th>
                            <th class="py-3 px-6 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Key
                            </th>
                            <th class="py-3 px-6 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                类型
                            </th>
                            <th class="py-3 px-6 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Value
                            </th>
                            <th class="py-3 px-6 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                最后修改人
                            </th>
                            <th class="py-3 px-6 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                最后修改时间
                            </th>
                            <th class="py-3 px-6 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                                操作
                            </th>
                        </tr>
                    </thead>
                    <tbody class="bg-white divide-y divide-gray-200">
                        <tr v-for="(item, index) in items" :key="index" class="hover:bg-gray-50 transition-colors">
                            <td class="py-4 px-6 whitespace-nowrap">
                                <span :class="{
                                    'px-3 py-1 rounded-full text-xs font-medium': true,
                                    'bg-green-100 text-green-800': item.status === '已发布',
                                    'bg-yellow-100 text-yellow-800': item.status === '待发布',
                                    'bg-gray-100 text-gray-800': item.status === '已下线'
                                }">
                                    {{ item.status }}
                                </span>
                            </td>
                            <td class="py-4 px-6 whitespace-nowrap">
                                <div class="flex items-center">
                                    <span class="font-medium text-gray-900">{{ item.key }}</span>
                                    <i v-if="item.isEncrypted" class="fas fa-lock text-yellow-500 ml-2"></i>
                                </div>
                            </td>
                            <td class="py-4 px-6 whitespace-nowrap text-sm text-gray-500">
                                <span class="px-2 py-1 bg-blue-100 text-blue-800 rounded text-xs font-medium">
                                    {{ item.type }}
                                </span>
                            </td>
                            <td class="py-4 px-6 max-w-xs truncate">
                                <div class="flex items-center">
                                    <span class="text-gray-600 truncate">{{ item.value }}</span>
                                    <button class="ml-2 text-gray-400 hover:text-primary-500"
                                        @click="copyValue(item.value)">
                                        <i class="fas fa-copy"></i>
                                    </button>
                                </div>
                            </td>
                            <td class="py-4 px-6 whitespace-nowrap text-sm text-gray-500">
                                <div class="flex items-center">
                                    <div
                                        class="w-6 h-6 rounded-full bg-blue-500 flex items-center justify-center text-white text-xs mr-2">
                                        {{ item.modifier.charAt(0) }}
                                    </div>
                                    {{ item.modifier }}
                                </div>
                            </td>
                            <td class="py-4 px-6 whitespace-nowrap text-sm text-gray-500">
                                {{ item.modifiedTime }}
                            </td>
                            <td class="py-4 px-6 whitespace-nowrap text-right text-sm font-medium">
                                <div class="flex justify-end space-x-3">
                                    <button class="text-blue-500 hover:text-blue-700" title="编辑">
                                        <i class="fas fa-edit"></i>
                                    </button>
                                    <button class="text-green-500 hover:text-green-700" title="克隆">
                                        <i class="fas fa-clone"></i>
                                    </button>
                                    <button class="text-gray-500 hover:text-red-500" title="删除">
                                        <i class="fas fa-trash-alt"></i>
                                    </button>
                                </div>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>

            <!-- 表格底部 -->
            <div class="px-6 py-4 bg-gray-50 border-t border-gray-200 flex items-center justify-between">
                <div class="text-sm text-gray-700">
                    显示 <span class="font-medium">{{ startItem }}</span> 到
                    <span class="font-medium">{{ endItem }}</span> 项，共
                    <span class="font-medium">{{ totalItems }}</span> 项
                </div>

                <div class="flex items-center space-x-4">
                    <div class="flex items-center text-sm text-gray-700">
                        每页显示
                        <select
                            class="mx-2 px-2 py-1 border border-gray-300 rounded focus:ring-primary-500 focus:border-primary-500">
                            <option>10</option>
                            <option>20</option>
                            <option>50</option>
                        </select>
                        项
                    </div>

                    <div class="flex items-center space-x-1">
                        <button class="p-2 rounded hover:bg-gray-200 disabled:opacity-50" :disabled="currentPage === 1">
                            <i class="fas fa-chevron-left"></i>
                        </button>
                        <button v-for="page in visiblePages" :key="page"
                            class="w-8 h-8 rounded flex items-center justify-center" :class="{
                                'bg-primary-500 text-white': currentPage === page,
                                'text-gray-700 hover:bg-gray-200': currentPage !== page
                            }">
                            {{ page }}
                        </button>
                        <button class="p-2 rounded hover:bg-gray-200 disabled:opacity-50"
                            :disabled="currentPage === totalPages">
                            <i class="fas fa-chevron-right"></i>
                        </button>
                    </div>
                </div>
            </div>
        </div>

        <!-- 空状态提示 -->
        <div v-if="items.length === 0" class="bg-white rounded-xl shadow-sm p-12 text-center mt-6">
            <div class="mx-auto w-24 h-24 rounded-full bg-gray-100 flex items-center justify-center mb-6">
                <i class="fas fa-file-alt text-gray-400 text-4xl"></i>
            </div>
            <h3 class="text-xl font-medium text-gray-700 mb-2">暂无配置项</h3>
            <p class="text-gray-500 mb-6">当前命名空间下还没有任何配置项，请点击下方按钮添加</p>
            <button class="btn bg-primary-500 hover:bg-primary-600 text-white">
                <i class="fas fa-plus mr-2"></i>添加配置项
            </button>
        </div>

        <!-- 使用说明卡片 -->
        <div class="mt-8 grid grid-cols-1 md:grid-cols-3 gap-6">
            <div class="bg-blue-50 border border-blue-100 rounded-xl p-5">
                <div class="flex items-center mb-3">
                    <div class="w-10 h-10 rounded-full bg-blue-500 flex items-center justify-center text-white mr-3">
                        <i class="fas fa-book"></i>
                    </div>
                    <h4 class="font-medium text-gray-800">命名空间指南</h4>
                </div>
                <p class="text-sm text-gray-600 mb-4">命名空间用于对配置项进行逻辑分组管理，每个命名空间相当于一个独立的配置文件。</p>
                <button class="text-sm text-blue-600 hover:text-blue-800 font-medium">
                    查看文档 <i class="fas fa-arrow-right ml-1 text-xs"></i>
                </button>
            </div>

            <div class="bg-green-50 border border-green-100 rounded-xl p-5">
                <div class="flex items-center mb-3">
                    <div class="w-10 h-10 rounded-full bg-green-500 flex items-center justify-center text-white mr-3">
                        <i class="fas fa-shield-alt"></i>
                    </div>
                    <h4 class="font-medium text-gray-800">安全配置建议</h4>
                </div>
                <p class="text-sm text-gray-600 mb-4">敏感配置项应启用加密存储，并严格控制访问权限，确保配置安全。</p>
                <button class="text-sm text-green-600 hover:text-green-800 font-medium">
                    安全设置 <i class="fas fa-arrow-right ml-1 text-xs"></i>
                </button>
            </div>

            <div class="bg-purple-50 border border-purple-100 rounded-xl p-5">
                <div class="flex items-center mb-3">
                    <div class="w-10 h-10 rounded-full bg-purple-500 flex items-center justify-center text-white mr-3">
                        <i class="fas fa-history"></i>
                    </div>
                    <h4 class="font-medium text-gray-800">版本与发布</h4>
                </div>
                <p class="text-sm text-gray-600 mb-4">每次发布都会生成新版本，支持版本回滚和灰度发布，确保配置变更安全可控。</p>
                <button class="text-sm text-purple-600 hover:text-purple-800 font-medium">
                    发布管理 <i class="fas fa-arrow-right ml-1 text-xs"></i>
                </button>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, computed } from 'vue';

// 模拟数据
const namespaceName = ref('payment-service');
const currentReleaseId = ref('20230811001');
const lastPublishTime = ref('2023-08-11 14:30:22');
const types = ref(['字符串', '数字', '布尔值', 'JSON', '加密文本']);
const statuses = ref(['已发布', '待发布', '已下线']);

// 分页相关
const currentPage = ref(1);
const itemsPerPage = ref(10);
const totalItems = ref(45);

const startItem = computed(() => (currentPage.value - 1) * itemsPerPage.value + 1);
const endItem = computed(() => Math.min(currentPage.value * itemsPerPage.value, totalItems.value));
const totalPages = computed(() => Math.ceil(totalItems.value / itemsPerPage.value));

// 生成可见页码
const visiblePages = computed(() => {
    const pages = [];
    const maxVisible = 5;
    let start = Math.max(1, currentPage.value - Math.floor(maxVisible / 2));
    let end = Math.min(totalPages.value, start + maxVisible - 1);

    if (end - start < maxVisible - 1) {
        start = Math.max(1, end - maxVisible + 1);
    }

    for (let i = start; i <= end; i++) {
        pages.push(i);
    }

    return pages;
});

// 配置项数据
const items = ref([
    {
        key: 'api.timeout',
        value: '5000',
        type: '数字',
        status: '已发布',
        modifier: '张明',
        modifiedTime: '2023-08-10 09:24:17',
        isEncrypted: false
    },
    {
        key: 'logging.level',
        value: 'INFO',
        type: '字符串',
        status: '已发布',
        modifier: '李华',
        modifiedTime: '2023-08-09 14:32:45',
        isEncrypted: false
    },
    {
        key: 'payment.gateway.url',
        value: 'https://pay.example.com/v2',
        type: '字符串',
        status: '待发布',
        modifier: '王芳',
        modifiedTime: '2023-08-11 10:15:33',
        isEncrypted: false
    },
    {
        key: 'retry.maxAttempts',
        value: '3',
        type: '数字',
        status: '已发布',
        modifier: '张明',
        modifiedTime: '2023-08-08 16:45:21',
        isEncrypted: false
    },
    {
        key: 'feature.newCheckout',
        value: 'true',
        type: '布尔值',
        status: '已发布',
        modifier: '赵强',
        modifiedTime: '2023-08-07 11:30:18',
        isEncrypted: false
    },
    {
        key: 'db.connection',
        value: '{"host":"db.example.com","port":3306,"user":"app_user","password":"******"}',
        type: 'JSON',
        status: '已发布',
        modifier: '李华',
        modifiedTime: '2023-08-06 13:22:57',
        isEncrypted: true
    },
    {
        key: 'cache.ttl',
        value: '3600',
        type: '数字',
        status: '已下线',
        modifier: '王芳',
        modifiedTime: '2023-08-05 09:11:34',
        isEncrypted: false
    },
    {
        key: 'auth.jwtSecret',
        value: 'encrypted:******',
        type: '加密文本',
        status: '已发布',
        modifier: '张明',
        modifiedTime: '2023-08-04 15:43:22',
        isEncrypted: true
    },
    {
        key: 'analytics.enabled',
        value: 'false',
        type: '布尔值',
        status: '已发布',
        modifier: '赵强',
        modifiedTime: '2023-08-03 17:55:41',
        isEncrypted: false
    },
    {
        key: 'max.concurrent.requests',
        value: '100',
        type: '数字',
        status: '待发布',
        modifier: '李华',
        modifiedTime: '2023-08-11 08:20:15',
        isEncrypted: false
    }
]);

// 复制值到剪贴板
const copyValue = (value) => {
    navigator.clipboard.writeText(value);
    // 这里可以添加复制成功的提示
    console.log('已复制值: ', value);
};
</script>

<style scoped>
@reference "tailwindcss";
.namespace-view {
    font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}

.btn {
    @apply px-4 py-2 rounded-lg transition-colors duration-200 flex items-center font-medium text-sm;
}

table {
    border-collapse: separate;
    border-spacing: 0;
}

th {
    @apply text-left font-semibold;
}

tr:not(:last-child) td {
    border-bottom: 1px solid #f3f4f6;
}

tbody tr:hover {
    @apply bg-gray-50;
}
</style>