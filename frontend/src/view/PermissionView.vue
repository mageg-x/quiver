<template>
    <div class="permission-management w-full p-6">
        <!-- 标题和操作按钮 -->
        <div class="flex justify-between items-center mb-6">
            <div>
                <h2 class="text-2xl font-bold text-gray-800">权限管理</h2>
                <p class="text-gray-600">管理用户对资源的操作权限</p>
            </div>
            <button @click="showAddPermissionDialog = true"
                class="bg-primary-600 hover:bg-primary-700 text-white px-4 py-2 rounded-lg flex items-center transition-colors">
                <i class="fas fa-plus mr-2"></i> 添加权限
            </button>
        </div>

        <!-- 搜索和筛选区域 -->
        <div class="bg-white rounded-xl shadow-sm p-4 mb-6 border border-gray-200">
            <div class="grid grid-cols-1 md:grid-cols-5 gap-4">
                <div class="relative">
                    <input v-model="searchQuery" type="text" placeholder="搜索用户或资源..."
                        class="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent">
                    <i class="fas fa-search absolute left-3 top-3 text-gray-400"></i>
                </div>

                <div>
                    <select v-model="selectedResourceType"
                        class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent">
                        <option value="">所有资源类型</option>
                        <option value="app">应用</option>
                        <option value="cluster">集群</option>
                        <option value="namespace">命名空间</option>
                        <option value="global">全局</option>
                    </select>
                </div>

                <div>
                    <select v-model="selectedUser"
                        class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent">
                        <option value="">所有用户</option>
                        <option v-for="user in users" :value="user.id">{{ user.name }}</option>
                    </select>
                </div>

                <div>
                    <select v-model="selectedAction"
                        class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent">
                        <option value="">所有操作</option>
                        <option value="read">读取</option>
                        <option value="create">创建</option>
                        <option value="update">修改</option>
                        <option value="delete">删除</option>
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

        <!-- 权限列表 -->
        <div class="bg-white rounded-xl shadow-sm overflow-hidden border border-gray-200">
            <div class="overflow-x-auto">
                <table class="min-w-full divide-y divide-gray-200">
                    <thead class="bg-gray-50">
                        <tr>
                            <th scope="col"
                                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                用户</th>
                            <th scope="col"
                                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                资源类型</th>
                            <th scope="col"
                                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                资源名称</th>
                            <th scope="col"
                                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                操作权限</th>
                            <th scope="col"
                                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                创建时间</th>
                            <th scope="col"
                                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                操作</th>
                        </tr>
                    </thead>
                    <tbody class="bg-white divide-y divide-gray-200">
                        <tr v-for="permission in filteredPermissions" :key="permission.id"
                            class="hover:bg-gray-50 transition-colors">
                            <td class="px-6 py-4 whitespace-nowrap">
                                <div class="flex items-center">
                                    <div class="flex-shrink-0 h-10 w-10">
                                        <div
                                            class="h-10 w-10 rounded-full bg-blue-100 flex items-center justify-center">
                                            <span class="text-blue-600 font-medium">{{ getUserAvatar(permission.userId)
                                                }}</span>
                                        </div>
                                    </div>
                                    <div class="ml-4">
                                        <div class="text-sm font-medium text-gray-900">{{ getUserName(permission.userId)
                                            }}</div>
                                    </div>
                                </div>
                            </td>
                            <td class="px-6 py-4 whitespace-nowrap">
                                <span :class="resourceTypeClass(permission.resourceType)"
                                    class="px-2 py-1 text-xs rounded-full">
                                    {{ resourceTypeLabel(permission.resourceType) }}
                                </span>
                            </td>
                            <td class="px-6 py-4">
                                <div class="text-sm text-gray-900">
                                    {{ permission.resourceName || '所有资源' }}
                                </div>
                            </td>
                            <td class="px-6 py-4 whitespace-nowrap">
                                <div class="flex flex-wrap gap-1">
                                    <span v-for="action in permission.actions" :key="action"
                                        :class="actionClass(action)" class="px-2 py-1 text-xs rounded-full">
                                        {{ actionLabel(action) }}
                                    </span>
                                </div>
                            </td>
                            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                {{ formatDate(permission.createdAt) }}
                            </td>
                            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                <div class="flex space-x-2">
                                    <button @click="editPermission(permission)"
                                        class="action-btn text-blue-600 hover:text-blue-800">
                                        <i class="fas fa-edit"></i>
                                    </button>
                                    <button @click="confirmDelete(permission)"
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
            <div v-if="filteredPermissions.length === 0" class="text-center py-12">
                <div class="mx-auto w-24 h-24 rounded-full bg-gray-100 flex items-center justify-center mb-4">
                    <i class="fas fa-shield-alt text-gray-400 text-3xl"></i>
                </div>
                <h3 class="text-lg font-medium text-gray-900 mb-1">暂无权限规则</h3>
                <p class="text-gray-500 max-w-md mx-auto">
                    当前系统中还没有定义权限规则，点击"添加权限"按钮创建第一条规则
                </p>
                <button @click="showAddPermissionDialog = true"
                    class="mt-4 bg-primary-600 hover:bg-primary-700 text-white px-4 py-2 rounded-lg inline-flex items-center transition-colors">
                    <i class="fas fa-plus mr-2"></i> 添加权限
                </button>
            </div>

            <!-- 分页控件 -->
            <div v-if="filteredPermissions.length > 0"
                class="bg-white px-4 py-3 flex items-center justify-between border-t border-gray-200 sm:px-6">
                <div class="hidden sm:flex-1 sm:flex sm:items-center sm:justify-between">
                    <div>
                        <p class="text-sm text-gray-700">
                            显示第 <span class="font-medium">1</span> 到 <span class="font-medium">{{
                                filteredPermissions.length }}</span> 条，共 <span class="font-medium">{{
                                filteredPermissions.length }}</span> 条记录
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
        <div class="grid grid-cols-1 md:grid-cols-4 gap-6 mt-8">
            <div class="stat-card bg-white rounded-xl shadow-sm p-6 flex items-center border border-gray-200">
                <div class="bg-blue-100 p-4 rounded-lg mr-4">
                    <i class="fas fa-user-shield text-blue-600 text-2xl"></i>
                </div>
                <div>
                    <div class="text-gray-500 text-sm">权限规则总数</div>
                    <div class="text-2xl font-bold text-gray-800">{{ permissions.length }}</div>
                </div>
            </div>

            <div class="stat-card bg-white rounded-xl shadow-sm p-6 flex items-center border border-gray-200">
                <div class="bg-green-100 p-4 rounded-lg mr-4">
                    <i class="fas fa-users text-green-600 text-2xl"></i>
                </div>
                <div>
                    <div class="text-gray-500 text-sm">已授权用户</div>
                    <div class="text-2xl font-bold text-gray-800">{{ uniqueUsersCount }}</div>
                </div>
            </div>

            <div class="stat-card bg-white rounded-xl shadow-sm p-6 flex items-center border border-gray-200">
                <div class="bg-amber-100 p-4 rounded-lg mr-4">
                    <i class="fas fa-cubes text-amber-600 text-2xl"></i>
                </div>
                <div>
                    <div class="text-gray-500 text-sm">受控资源</div>
                    <div class="text-2xl font-bold text-gray-800">{{ uniqueResourcesCount }}</div>
                </div>
            </div>

            <div class="stat-card bg-white rounded-xl shadow-sm p-6 flex items-center border border-gray-200">
                <div class="bg-purple-100 p-4 rounded-lg mr-4">
                    <i class="fas fa-key text-purple-600 text-2xl"></i>
                </div>
                <div>
                    <div class="text-gray-500 text-sm">最高频操作</div>
                    <div class="text-2xl font-bold text-gray-800">{{ mostCommonAction }}</div>
                </div>
            </div>
        </div>

        <!-- 添加/编辑权限对话框 -->
        <transition name="fade">
            <div v-if="showAddPermissionDialog" class="fixed inset-0 flex items-center justify-center z-50 p-4">
                <div class="dialog-overlay fixed inset-0 bg-black bg-opacity-50" @click="closeDialog"></div>
                <div class="bg-white rounded-xl shadow-xl w-full max-w-2xl relative z-10 transform transition-all">
                    <div class="p-6">
                        <div class="flex justify-between items-center mb-4">
                            <h3 class="text-lg font-bold text-gray-800">{{ editingPermission ? '编辑权限规则' : '添加新权限规则' }}
                            </h3>
                            <button @click="closeDialog" class="text-gray-400 hover:text-gray-500">
                                <i class="fas fa-times"></i>
                            </button>
                        </div>

                        <div class="space-y-6">
                            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                                <div>
                                    <label class="block text-sm font-medium text-gray-700 mb-1">选择用户</label>
                                    <div
                                        class="space-y-2 max-h-60 overflow-y-auto p-2 border border-gray-300 rounded-lg">
                                        <div v-for="user in users" :key="user.id" @click="selectUser(user.id)"
                                            :class="{ 'bg-blue-50 border-blue-200': currentPermission.userId === user.id, 'bg-white border-gray-200': currentPermission.userId !== user.id }"
                                            class="p-3 border rounded-lg cursor-pointer transition-colors">
                                            <div class="flex items-center">
                                                <div class="flex-shrink-0 h-8 w-8">
                                                    <div
                                                        class="h-8 w-8 rounded-full bg-blue-100 flex items-center justify-center">
                                                        <span class="text-blue-600 text-sm font-medium">{{ user.avatar
                                                            }}</span>
                                                    </div>
                                                </div>
                                                <div class="ml-3">
                                                    <div class="text-sm font-medium text-gray-900">{{ user.name }}</div>
                                                    <div class="text-xs text-gray-500">{{ user.email }}</div>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                </div>

                                <div>
                                    <label class="block text-sm font-medium text-gray-700 mb-1">资源类型</label>
                                    <select v-model="currentPermission.resourceType" @change="updateResourceOptions"
                                        class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent">
                                        <option value="app">应用</option>
                                        <option value="cluster">集群</option>
                                        <option value="namespace">命名空间</option>
                                        <option value="global">全局</option>
                                    </select>

                                    <div v-if="currentPermission.resourceType !== 'global'" class="mt-3">
                                        <label class="block text-sm font-medium text-gray-700 mb-1">选择资源</label>
                                        <select v-model="currentPermission.resourceId"
                                            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent">
                                            <option value="">所有资源</option>
                                            <option v-for="resource in resourceOptions" :key="resource.id"
                                                :value="resource.id">
                                                {{ resource.name }}
                                            </option>
                                        </select>
                                    </div>
                                </div>
                            </div>

                            <div>
                                <label class="block text-sm font-medium text-gray-700 mb-1">操作权限</label>
                                <div class="grid grid-cols-2 md:grid-cols-4 gap-3">
                                    <div v-for="action in actionOptions" :key="action.value"
                                        @click="toggleAction(action.value)"
                                        :class="{ 'bg-blue-100 border-blue-500 text-blue-700': currentPermission.actions.includes(action.value), 'bg-gray-100 border-gray-300 text-gray-700': !currentPermission.actions.includes(action.value) }"
                                        class="p-3 border rounded-lg cursor-pointer flex items-center transition-colors">
                                        <div class="w-6 h-6 flex items-center justify-center mr-2">
                                            <i :class="action.icon" class="text-lg"></i>
                                        </div>
                                        <span>{{ action.label }}</span>
                                    </div>
                                </div>
                            </div>

                            <div>
                                <label class="block text-sm font-medium text-gray-700 mb-1">权限描述</label>
                                <textarea v-model="currentPermission.description" rows="2" placeholder="可选：添加权限规则描述"
                                    class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"></textarea>
                            </div>
                        </div>

                        <div class="mt-6 flex justify-end space-x-3">
                            <button @click="closeDialog"
                                class="px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors">取消</button>
                            <button @click="savePermission"
                                class="px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 transition-colors">
                                {{ editingPermission ? '保存更改' : '添加权限' }}
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

                        <h3 class="text-lg font-bold text-center text-gray-800 mb-2">确认删除权限规则</h3>
                        <p class="text-gray-600 text-center mb-6">确定要删除此权限规则吗？此操作不可恢复。</p>

                        <div class="flex justify-center space-x-3">
                            <button @click="showDeleteDialog = false"
                                class="px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors">取消</button>
                            <button @click="deletePermission"
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

// 用户数据
const users = ref([
    { id: 1, name: '张明', email: 'zhangming@example.com', avatar: '张' },
    { id: 2, name: '李红', email: 'lihong@example.com', avatar: '李' },
    { id: 3, name: '王伟', email: 'wangwei@example.com', avatar: '王' },
    { id: 4, name: '赵静', email: 'zhaojing@example.com', avatar: '赵' },
]);

// 资源数据
const apps = ref([
    { id: 1, name: '订单服务' },
    { id: 2, name: '支付服务' },
    { id: 3, name: '用户服务' },
]);

const clusters = ref([
    { id: 1, name: '北京生产集群', appId: 1 },
    { id: 2, name: '上海测试集群', appId: 1 },
    { id: 3, name: '支付生产集群', appId: 2 },
]);

const namespaces = ref([
    { id: 1, name: '订单命名空间', clusterId: 1 },
    { id: 2, name: '支付命名空间', clusterId: 3 },
    { id: 3, name: '用户命名空间', clusterId: 3 },
]);

// 权限数据
const permissions = ref([
    {
        id: 1,
        userId: 1,
        resourceType: 'app',
        resourceId: 1,
        resourceName: '订单服务',
        actions: ['read', 'update'],
        description: '允许对订单服务应用的读取和修改权限',
        createdAt: '2023-05-10'
    },
    {
        id: 2,
        userId: 2,
        resourceType: 'cluster',
        resourceId: 3,
        resourceName: '支付生产集群',
        actions: ['read', 'create', 'update', 'delete'],
        description: '支付集群的完全控制权限',
        createdAt: '2023-06-15'
    },
    {
        id: 3,
        userId: 3,
        resourceType: 'global',
        resourceId: null,
        resourceName: null,
        actions: ['read'],
        description: '全局只读权限',
        createdAt: '2023-07-20'
    },
    {
        id: 4,
        userId: 1,
        resourceType: 'namespace',
        resourceId: 2,
        resourceName: '支付命名空间',
        actions: ['read', 'update'],
        description: '支付命名空间的读写权限',
        createdAt: '2023-08-05'
    },
]);

// 操作选项
const actionOptions = ref([
    { value: 'read', label: '读取', icon: 'fas fa-eye' },
    { value: 'create', label: '创建', icon: 'fas fa-plus-circle' },
    { value: 'update', label: '修改', icon: 'fas fa-edit' },
    { value: 'delete', label: '删除', icon: 'fas fa-trash' },
]);

// 资源类型标签
const resourceTypeLabel = (type) => {
    const labels = {
        app: '应用',
        cluster: '集群',
        namespace: '命名空间',
        global: '全局'
    };
    return labels[type] || type;
};

// 资源类型样式类
const resourceTypeClass = (type) => {
    return {
        'bg-blue-100 text-blue-800': type === 'app',
        'bg-green-100 text-green-800': type === 'cluster',
        'bg-purple-100 text-purple-800': type === 'namespace',
        'bg-gray-100 text-gray-800': type === 'global',
    };
};

// 操作标签
const actionLabel = (action) => {
    const labels = {
        read: '读取',
        create: '创建',
        update: '修改',
        delete: '删除'
    };
    return labels[action] || action;
};

// 操作样式类
const actionClass = (action) => {
    return {
        'bg-blue-100 text-blue-800': action === 'read',
        'bg-green-100 text-green-800': action === 'create',
        'bg-yellow-100 text-yellow-800': action === 'update',
        'bg-red-100 text-red-800': action === 'delete',
    };
};

// 搜索和筛选
const searchQuery = ref('');
const selectedResourceType = ref('');
const selectedUser = ref('');
const selectedAction = ref('');

// 过滤后的权限列表
const filteredPermissions = computed(() => {
    return permissions.value.filter(permission => {
        const user = users.value.find(u => u.id === permission.userId);
        const matchesSearch = user?.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
            (permission.resourceName && permission.resourceName.toLowerCase().includes(searchQuery.value.toLowerCase())) ||
            permission.description?.toLowerCase().includes(searchQuery.value.toLowerCase());

        const matchesResourceType = selectedResourceType.value ?
            permission.resourceType === selectedResourceType.value : true;

        const matchesUser = selectedUser.value ?
            permission.userId === parseInt(selectedUser.value) : true;

        const matchesAction = selectedAction.value ?
            permission.actions.includes(selectedAction.value) : true;

        return matchesSearch && matchesResourceType && matchesUser && matchesAction;
    });
});

// 统计信息
const uniqueUsersCount = computed(() => {
    const userIds = new Set(permissions.value.map(p => p.userId));
    return userIds.size;
});

const uniqueResourcesCount = computed(() => {
    const resourceIds = new Set(
        permissions.value
            .filter(p => p.resourceId)
            .map(p => `${p.resourceType}-${p.resourceId}`)
    );
    return resourceIds.size;
});

const mostCommonAction = computed(() => {
    const actionCounts = {
        read: 0,
        create: 0,
        update: 0,
        delete: 0
    };

    permissions.value.forEach(p => {
        p.actions.forEach(action => {
            if (actionCounts[action] !== undefined) {
                actionCounts[action]++;
            }
        });
    });

    let mostCommon = '读取';
    let maxCount = 0;

    for (const [action, count] of Object.entries(actionCounts)) {
        if (count > maxCount) {
            maxCount = count;
            mostCommon = actionLabel(action);
        }
    }

    return mostCommon;
});

// 对话框状态
const showAddPermissionDialog = ref(false);
const showDeleteDialog = ref(false);
const editingPermission = ref(false);

// 当前操作权限
const currentPermission = ref({
    id: null,
    userId: null,
    resourceType: 'app',
    resourceId: null,
    actions: ['read'],
    description: ''
});

// 资源选项
const resourceOptions = ref([]);

// 提示信息
const showToast = ref(false);
const toastMessage = ref('');

// 日期格式化
const formatDate = (dateString) => {
    const options = { year: 'numeric', month: 'long', day: 'numeric' };
    return new Date(dateString).toLocaleDateString('zh-CN', options);
};

// 获取用户名
const getUserName = (userId) => {
    const user = users.value.find(u => u.id === userId);
    return user ? user.name : '未知用户';
};

// 获取用户头像
const getUserAvatar = (userId) => {
    const user = users.value.find(u => u.id === userId);
    return user ? user.avatar : '?';
};

// 重置筛选
const resetFilters = () => {
    searchQuery.value = '';
    selectedResourceType.value = '';
    selectedUser.value = '';
    selectedAction.value = '';
};

// 更新资源选项
const updateResourceOptions = () => {
    switch (currentPermission.value.resourceType) {
        case 'app':
            resourceOptions.value = apps.value;
            break;
        case 'cluster':
            resourceOptions.value = clusters.value;
            break;
        case 'namespace':
            resourceOptions.value = namespaces.value;
            break;
        default:
            resourceOptions.value = [];
    }
};

// 选择用户
const selectUser = (userId) => {
    currentPermission.value.userId = userId;
};

// 切换操作
const toggleAction = (action) => {
    const index = currentPermission.value.actions.indexOf(action);
    if (index > -1) {
        currentPermission.value.actions.splice(index, 1);
    } else {
        currentPermission.value.actions.push(action);
    }
};

// 编辑权限
const editPermission = (permission) => {
    currentPermission.value = { ...permission };
    updateResourceOptions();
    editingPermission.value = true;
    showAddPermissionDialog.value = true;
};

// 确认删除
const confirmDelete = (permission) => {
    currentPermission.value = { ...permission };
    showDeleteDialog.value = true;
};

// 删除权限
const deletePermission = () => {
    permissions.value = permissions.value.filter(p => p.id !== currentPermission.value.id);
    showDeleteDialog.value = false;
    showToastMessage('权限规则已删除');
};

// 保存权限
const savePermission = () => {
    if (!currentPermission.value.userId) {
        showToastMessage('请选择用户', 'error');
        return;
    }

    if (currentPermission.value.resourceType !== 'global' && !currentPermission.value.resourceId) {
        showToastMessage('请选择资源', 'error');
        return;
    }

    if (currentPermission.value.actions.length === 0) {
        showToastMessage('请至少选择一个操作权限', 'error');
        return;
    }

    // 获取资源名称
    if (currentPermission.value.resourceId) {
        let resource = null;
        switch (currentPermission.value.resourceType) {
            case 'app':
                resource = apps.value.find(a => a.id === currentPermission.value.resourceId);
                break;
            case 'cluster':
                resource = clusters.value.find(c => c.id === currentPermission.value.resourceId);
                break;
            case 'namespace':
                resource = namespaces.value.find(n => n.id === currentPermission.value.resourceId);
                break;
        }
        currentPermission.value.resourceName = resource ? resource.name : '';
    } else {
        currentPermission.value.resourceName = null;
    }

    if (editingPermission.value) {
        // 更新现有权限
        const index = permissions.value.findIndex(p => p.id === currentPermission.value.id);
        if (index !== -1) {
            permissions.value[index] = { ...currentPermission.value };
        }
        showToastMessage('权限规则已更新');
    } else {
        // 添加新权限
        const newId = Math.max(...permissions.value.map(p => p.id), 0) + 1;
        permissions.value.push({
            id: newId,
            ...currentPermission.value,
            createdAt: new Date().toISOString().split('T')[0]
        });
        showToastMessage('新权限规则已添加');
    }
    closeDialog();
};

// 关闭对话框
const closeDialog = () => {
    showAddPermissionDialog.value = false;
    showDeleteDialog.value = false;
    editingPermission.value = false;
    currentPermission.value = {
        id: null,
        userId: null,
        resourceType: 'app',
        resourceId: null,
        actions: ['read'],
        description: ''
    };
};

// 显示提示信息
const showToastMessage = (message, type = 'success') => {
    toastMessage.value = message;
    showToast.value = true;
    setTimeout(() => {
        showToast.value = false;
    }, 3000);
};

// 初始化
onMounted(() => {
    updateResourceOptions();
});
</script>


<script>
export default { name: 'PermissionView' }
</script>

<style scoped>
.permission-management {
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

.permission-action {
    transition: all 0.2s ease;
}

.permission-action:hover {
    transform: translateY(-2px);
}
</style>