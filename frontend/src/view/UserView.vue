<template>
    <div class="user-management w-full p-6">
        <!-- 标题和操作按钮 -->
        <div class="flex justify-between items-center mb-6">
            <div>
                <h2 class="text-2xl font-bold text-gray-800">用户管理</h2>
                <p class="text-gray-600">管理系统用户账户和权限</p>
            </div>
            <button @click="showAddUserDialog = true"
                class="bg-primary-600 hover:bg-primary-700 text-white px-4 py-2 rounded-lg flex items-center transition-colors">
                <i class="fas fa-plus mr-2"></i> 添加用户
            </button>
        </div>

        <!-- 搜索和筛选区域 -->
        <div class="bg-white rounded-xl shadow-sm p-4 mb-6 border border-gray-200">
            <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
                <div class="relative">
                    <input v-model="searchQuery" type="text" placeholder="搜索用户名或邮箱..."
                        class="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent">
                    <i class="fas fa-search absolute left-3 top-3 text-gray-400"></i>
                </div>

                <div>
                    <select v-model="selectedRole"
                        class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent">
                        <option value="">所有角色</option>
                        <option v-for="role in roles" :value="role.value">{{ role.label }}</option>
                    </select>
                </div>

                <div>
                    <select v-model="selectedStatus"
                        class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent">
                        <option value="">所有状态</option>
                        <option value="active">正常</option>
                        <option value="disabled">已禁用</option>
                        <option value="pending">待激活</option>
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

        <!-- 用户列表 -->
        <div class="bg-white rounded-xl shadow-sm overflow-hidden border border-gray-200">
            <div class="overflow-x-auto">
                <table class="min-w-full divide-y divide-gray-200">
                    <thead class="bg-gray-50">
                        <tr>
                            <th scope="col"
                                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                用户信息</th>
                            <th scope="col"
                                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                角色</th>
                            <th scope="col"
                                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                状态</th>
                            <th scope="col"
                                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                最后登录</th>
                            <th scope="col"
                                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                创建时间</th>
                            <th scope="col"
                                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                操作</th>
                        </tr>
                    </thead>
                    <tbody class="bg-white divide-y divide-gray-200">
                        <tr v-for="user in filteredUsers" :key="user.id" class="hover:bg-gray-50 transition-colors">
                            <td class="px-6 py-4">
                                <div class="flex items-center">
                                    <div class="flex-shrink-0 h-10 w-10">
                                        <div
                                            class="h-10 w-10 rounded-full bg-blue-100 flex items-center justify-center">
                                            <span class="text-blue-600 font-medium">{{ user.avatar }}</span>
                                        </div>
                                    </div>
                                    <div class="ml-4">
                                        <div class="text-sm font-medium text-gray-900">{{ user.name }}</div>
                                        <div class="text-sm text-gray-500">{{ user.email }}</div>
                                    </div>
                                </div>
                            </td>
                            <td class="px-6 py-4 whitespace-nowrap">
                                <span :class="roleClass(user.role)" class="px-2 py-1 text-xs rounded-full">
                                    {{ roleLabel(user.role) }}
                                </span>
                            </td>
                            <td class="px-6 py-4 whitespace-nowrap">
                                <span :class="statusClass(user.status)" class="px-2 py-1 text-xs rounded-full">
                                    {{ 正常}}
                                </span>
                            </td>
                            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                {{ formatDateTime(user.lastLogin) }}
                            </td>
                            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                {{ formatDate(user.createdAt) }}
                            </td>
                            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                <div class="flex space-x-2">
                                    <button @click="editUser(user)"
                                        class="action-btn text-blue-600 hover:text-blue-800">
                                        <i class="fas fa-edit"></i>
                                    </button>
                                    <button @click="toggleUserStatus(user)"
                                        :class="{ 'text-green-600 hover:text-green-800': user.status === 'disabled', 'text-orange-600 hover:text-orange-800': user.status === 'active' }"
                                        class="action-btn">
                                        <i v-if="user.status === 'disabled'" class="fas fa-user-check"></i>
                                        <i v-else class="fas fa-user-slash"></i>
                                    </button>
                                    <button @click="resetPassword(user)"
                                        class="action-btn text-purple-600 hover:text-purple-800">
                                        <i class="fas fa-key"></i>
                                    </button>
                                    <button @click="confirmDelete(user)"
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
            <div v-if="filteredUsers.length === 0" class="text-center py-12">
                <div class="mx-auto w-24 h-24 rounded-full bg-gray-100 flex items-center justify-center mb-4">
                    <i class="fas fa-users text-gray-400 text-3xl"></i>
                </div>
                <h3 class="text-lg font-medium text-gray-900 mb-1">暂无用户数据</h3>
                <p class="text-gray-500 max-w-md mx-auto">
                    当前系统中还没有用户，点击"添加用户"按钮创建第一个用户
                </p>
                <button @click="showAddUserDialog = true"
                    class="mt-4 bg-primary-600 hover:bg-primary-700 text-white px-4 py-2 rounded-lg inline-flex items-center transition-colors">
                    <i class="fas fa-plus mr-2"></i> 添加用户
                </button>
            </div>

            <!-- 分页控件 -->
            <div v-if="filteredUsers.length > 0"
                class="bg-white px-4 py-3 flex items-center justify-between border-t border-gray-200 sm:px-6">
                <div class="hidden sm:flex-1 sm:flex sm:items-center sm:justify-between">
                    <div>
                        <p class="text-sm text-gray-700">
                            显示第 <span class="font-medium">1</span> 到 <span class="font-medium">{{ filteredUsers.length
                                }}</span> 条，共 <span class="font-medium">{{ filteredUsers.length }}</span> 条记录
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
                    <i class="fas fa-users text-blue-600 text-2xl"></i>
                </div>
                <div>
                    <div class="text-gray-500 text-sm">用户总数</div>
                    <div class="text-2xl font-bold text-gray-800">{{ filteredUsers.length }}</div>
                </div>
            </div>

            <div class="stat-card bg-white rounded-xl shadow-sm p-6 flex items-center border border-gray-200">
                <div class="bg-green-100 p-4 rounded-lg mr-4">
                    <i class="fas fa-user-check text-green-600 text-2xl"></i>
                </div>
                <div>
                    <div class="text-gray-500 text-sm">活跃用户</div>
                    <div class="text-2xl font-bold text-gray-800">{{ activeUsersCount }}</div>
                </div>
            </div>

            <div class="stat-card bg-white rounded-xl shadow-sm p-6 flex items-center border border-gray-200">
                <div class="bg-amber-100 p-4 rounded-lg mr-4">
                    <i class="fas fa-user-shield text-amber-600 text-2xl"></i>
                </div>
                <div>
                    <div class="text-gray-500 text-sm">管理员</div>
                    <div class="text-2xl font-bold text-gray-800">{{ adminUsersCount }}</div>
                </div>
            </div>

            <div class="stat-card bg-white rounded-xl shadow-sm p-6 flex items-center border border-gray-200">
                <div class="bg-purple-100 p-4 rounded-lg mr-4">
                    <i class="fas fa-clock text-purple-600 text-2xl"></i>
                </div>
                <div>
                    <div class="text-gray-500 text-sm">最近7天登录</div>
                    <div class="text-2xl font-bold text-gray-800">{{ recentLoginCount }}</div>
                </div>
            </div>
        </div>

        <!-- 添加/编辑用户对话框 -->
        <transition name="fade">
            <div v-if="showAddUserDialog" class="fixed inset-0 flex items-center justify-center z-50 p-4">
                <div class="dialog-overlay fixed inset-0 bg-black bg-opacity-50" @click="closeDialog"></div>
                <div class="bg-white rounded-xl shadow-xl w-full max-w-md relative z-10 transform transition-all">
                    <div class="p-6">
                        <div class="flex justify-between items-center mb-4">
                            <h3 class="text-lg font-bold text-gray-800">{{ editingUser ? '编辑用户' : '添加新用户' }}</h3>
                            <button @click="closeDialog" class="text-gray-400 hover:text-gray-500">
                                <i class="fas fa-times"></i>
                            </button>
                        </div>

                        <div class="space-y-4">
                            <div class="grid grid-cols-2 gap-4">
                                <div>
                                    <label class="block text-sm font-medium text-gray-700 mb-1">姓名</label>
                                    <input v-model="currentUser.name" type="text" placeholder="请输入姓名"
                                        class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent">
                                </div>

                                <div>
                                    <label class="block text-sm font-medium text-gray-700 mb-1">用户名</label>
                                    <input v-model="currentUser.username" type="text" placeholder="请输入用户名"
                                        class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent">
                                </div>
                            </div>

                            <div>
                                <label class="block text-sm font-medium text-gray-700 mb-1">邮箱</label>
                                <input v-model="currentUser.email" type="email" placeholder="user@example.com"
                                    class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent">
                            </div>

                            <div class="grid grid-cols-2 gap-4">
                                <div>
                                    <label class="block text-sm font-medium text-gray-700 mb-1">角色</label>
                                    <select v-model="currentUser.role"
                                        class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent">
                                        <option v-for="role in roles" :key="role.value" :value="role.value">{{
                                            role.label }}</option>
                                    </select>
                                </div>

                                <div>
                                    <label class="block text-sm font-medium text-gray-700 mb-1">状态</label>
                                    <select v-model="currentUser.status"
                                        class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent">
                                        <option value="active">正常</option>
                                        <option value="disabled">已禁用</option>
                                        <option value="pending">待激活</option>
                                    </select>
                                </div>
                            </div>

                            <div v-if="!editingUser" class="grid grid-cols-2 gap-4">
                                <div>
                                    <label class="block text-sm font-medium text-gray-700 mb-1">密码</label>
                                    <input v-model="currentUser.password" type="password" placeholder="设置用户密码"
                                        class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent">
                                </div>

                                <div>
                                    <label class="block text-sm font-medium text-gray-700 mb-1">确认密码</label>
                                    <input v-model="currentUser.confirmPassword" type="password" placeholder="再次输入密码"
                                        class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent">
                                </div>
                            </div>

                            <div>
                                <label class="block text-sm font-medium text-gray-700 mb-1">电话号码</label>
                                <input v-model="currentUser.phone" type="tel" placeholder="请输入电话号码"
                                    class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent">
                            </div>
                        </div>

                        <div class="mt-6 flex justify-end space-x-3">
                            <button @click="closeDialog"
                                class="px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors">取消</button>
                            <button @click="saveUser"
                                class="px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 transition-colors">
                                {{ editingUser ? '保存更改' : '添加用户' }}
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        </transition>

        <!-- 重置密码对话框 -->
        <transition name="fade">
            <div v-if="showResetPasswordDialog" class="fixed inset-0 flex items-center justify-center z-50 p-4">
                <div class="dialog-overlay fixed inset-0 bg-black bg-opacity-50"
                    @click="showResetPasswordDialog = false"></div>
                <div class="bg-white rounded-xl shadow-xl w-full max-w-md relative z-10 transform transition-all">
                    <div class="p-6">
                        <div class="flex justify-center mb-4">
                            <div class="w-16 h-16 rounded-full bg-blue-100 flex items-center justify-center">
                                <i class="fas fa-key text-blue-600 text-2xl"></i>
                            </div>
                        </div>

                        <h3 class="text-lg font-bold text-center text-gray-800 mb-2">重置用户密码</h3>
                        <p class="text-gray-600 text-center mb-4">确定要为用户 <span class="font-semibold">{{ currentUser.name
                                }}</span> 重置密码吗？</p>

                        <div class="space-y-4">
                            <div>
                                <label class="block text-sm font-medium text-gray-700 mb-1">新密码</label>
                                <input v-model="newPassword" type="password" placeholder="输入新密码"
                                    class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent">
                            </div>

                            <div>
                                <label class="block text-sm font-medium text-gray-700 mb-1">确认新密码</label>
                                <input v-model="confirmNewPassword" type="password" placeholder="再次输入新密码"
                                    class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent">
                            </div>
                        </div>

                        <div class="mt-6 flex justify-center space-x-3">
                            <button @click="showResetPasswordDialog = false"
                                class="px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors">取消</button>
                            <button @click="confirmResetPassword"
                                class="px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 transition-colors">重置密码</button>
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

                        <h3 class="text-lg font-bold text-center text-gray-800 mb-2">确认删除用户</h3>
                        <p class="text-gray-600 text-center mb-6">确定要删除用户 <span class="font-semibold">{{
                                currentUser.name }}</span> 吗？此操作不可恢复。</p>

                        <div class="flex justify-center space-x-3">
                            <button @click="showDeleteDialog = false"
                                class="px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors">取消</button>
                            <button @click="deleteUser"
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

// 角色选项
const roles = ref([
    { label: '系统管理员', value: 'admin' },
    { label: '普通用户', value: 'user' },
    { label: '配置管理员', value: 'config-admin' },
    { label: '观察者', value: 'viewer' },
]);

// 状态标签
const statusLabels = ref({
    active: '正常',
    disabled: '已禁用',
    pending: '待激活'
});

// 角色标签
const roleLabel = (role) => {
    const roleObj = roles.value.find(r => r.value === role);
    return roleObj ? roleObj.label : role;
};

// 状态类
const statusClass = (status) => {
    return {
        'bg-green-100 text-green-800': status === 'active',
        'bg-red-100 text-red-800': status === 'disabled',
        'bg-yellow-100 text-yellow-800': status === 'pending',
    };
};

// 角色类
const roleClass = (role) => {
    return {
        'bg-purple-100 text-purple-800': role === 'admin',
        'bg-blue-100 text-blue-800': role === 'config-admin',
        'bg-green-100 text-green-800': role === 'user',
        'bg-gray-100 text-gray-800': role === 'viewer',
    };
};

// 用户数据
const users = ref([
    {
        id: 1,
        name: '张明',
        username: 'zhangming',
        email: 'zhangming@example.com',
        role: 'admin',
        status: 'active',
        phone: '13800138000',
        lastLogin: '2023-08-15T14:30:00',
        createdAt: '2023-01-10',
        avatar: '张'
    },
    {
        id: 2,
        name: '李红',
        username: 'lihong',
        email: 'lihong@example.com',
        role: 'config-admin',
        status: 'active',
        phone: '13900139000',
        lastLogin: '2023-08-20T09:15:00',
        createdAt: '2023-02-15',
        avatar: '李'
    },
    {
        id: 3,
        name: '王伟',
        username: 'wangwei',
        email: 'wangwei@example.com',
        role: 'user',
        status: 'active',
        phone: '13700137000',
        lastLogin: '2023-08-18T16:45:00',
        createdAt: '2023-03-22',
        avatar: '王'
    },
    {
        id: 4,
        name: '赵静',
        username: 'zhaojing',
        email: 'zhaojing@example.com',
        role: 'viewer',
        status: 'disabled',
        phone: '13600136000',
        lastLogin: '2023-07-05T11:20:00',
        createdAt: '2023-04-05',
        avatar: '赵'
    },
    {
        id: 5,
        name: '陈华',
        username: 'chenhua',
        email: 'chenhua@example.com',
        role: 'user',
        status: 'pending',
        phone: '13500135000',
        lastLogin: null,
        createdAt: '2023-05-18',
        avatar: '陈'
    },
]);

// 搜索和筛选
const searchQuery = ref('');
const selectedRole = ref('');
const selectedStatus = ref('');

// 过滤后的用户列表
const filteredUsers = computed(() => {
    return users.value.filter(user => {
        const matchesSearch = user.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
            user.email.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
            user.username.toLowerCase().includes(searchQuery.value.toLowerCase());
        const matchesRole = selectedRole.value ? user.role === selectedRole.value : true;
        const matchesStatus = selectedStatus.value ? user.status === selectedStatus.value : true;

        return matchesSearch && matchesRole && matchesStatus;
    });
});

// 活跃用户数量
const activeUsersCount = computed(() => {
    return filteredUsers.value.filter(user => user.status === 'active').length;
});

// 管理员用户数量
const adminUsersCount = computed(() => {
    return filteredUsers.value.filter(user => user.role === 'admin').length;
});

// 最近登录用户数量（最近7天）
const recentLoginCount = computed(() => {
    const sevenDaysAgo = new Date();
    sevenDaysAgo.setDate(sevenDaysAgo.getDate() - 7);

    return filteredUsers.value.filter(user => {
        if (!user.lastLogin) return false;
        const loginDate = new Date(user.lastLogin);
        return loginDate > sevenDaysAgo;
    }).length;
});

// 对话框状态
const showAddUserDialog = ref(false);
const showResetPasswordDialog = ref(false);
const showDeleteDialog = ref(false);
const editingUser = ref(false);

// 当前操作用户
const currentUser = ref({
    id: null,
    name: '',
    username: '',
    email: '',
    role: 'user',
    status: 'active',
    phone: '',
    password: '',
    confirmPassword: ''
});

// 密码重置相关
const newPassword = ref('');
const confirmNewPassword = ref('');

// 提示信息
const showToast = ref(false);
const toastMessage = ref('');

// 日期格式化
const formatDate = (dateString) => {
    const options = { year: 'numeric', month: 'long', day: 'numeric' };
    return new Date(dateString).toLocaleDateString('zh-CN', options);
};

// 日期时间格式化
const formatDateTime = (dateTimeString) => {
    if (!dateTimeString) return '从未登录';
    const options = { year: 'numeric', month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' };
    return new Date(dateTimeString).toLocaleDateString('zh-CN', options);
};

// 重置筛选
const resetFilters = () => {
    searchQuery.value = '';
    selectedRole.value = '';
    selectedStatus.value = '';
};

// 编辑用户
const editUser = (user) => {
    currentUser.value = { ...user, password: '', confirmPassword: '' };
    editingUser.value = true;
    showAddUserDialog.value = true;
};

// 切换用户状态
const toggleUserStatus = (user) => {
    const newStatus = user.status === 'active' ? 'disabled' : 'active';
    const index = users.value.findIndex(u => u.id === user.id);

    if (index !== -1) {
        users.value[index].status = newStatus;
        showToastMessage(`用户状态已${newStatus === 'active' ? '启用' : '禁用'}`);
    }
};

// 重置密码
const resetPassword = (user) => {
    currentUser.value = { ...user };
    newPassword.value = '';
    confirmNewPassword.value = '';
    showResetPasswordDialog.value = true;
};

// 确认重置密码
const confirmResetPassword = () => {
    if (newPassword.value !== confirmNewPassword.value) {
        showToastMessage('两次输入的密码不一致', 'error');
        return;
    }

    // 在实际应用中，这里会调用API重置密码
    showToastMessage('密码已成功重置');
    showResetPasswordDialog.value = false;
};

// 确认删除
const confirmDelete = (user) => {
    currentUser.value = { ...user };
    showDeleteDialog.value = true;
};

// 删除用户
const deleteUser = () => {
    users.value = users.value.filter(u => u.id !== currentUser.value.id);
    showDeleteDialog.value = false;
    showToastMessage('用户已成功删除');
};

// 保存用户
const saveUser = () => {
    if (editingUser.value) {
        // 更新现有用户
        const index = users.value.findIndex(u => u.id === currentUser.value.id);
        if (index !== -1) {
            users.value[index] = { ...currentUser.value };
        }
        showToastMessage('用户信息已更新');
    } else {
        // 添加新用户
        if (currentUser.value.password !== currentUser.value.confirmPassword) {
            showToastMessage('两次输入的密码不一致', 'error');
            return;
        }

        const newId = Math.max(...users.value.map(u => u.id), 0) + 1;
        users.value.push({
            id: newId,
            name: currentUser.value.name,
            username: currentUser.value.username,
            email: currentUser.value.email,
            role: currentUser.value.role,
            status: currentUser.value.status,
            phone: currentUser.value.phone,
            lastLogin: null,
            createdAt: new Date().toISOString().split('T')[0],
            avatar: currentUser.value.name.charAt(0)
        });
        showToastMessage('新用户已成功添加');
    }
    closeDialog();
};

// 关闭对话框
const closeDialog = () => {
    showAddUserDialog.value = false;
    showResetPasswordDialog.value = false;
    showDeleteDialog.value = false;
    editingUser.value = false;
    currentUser.value = {
        id: null,
        name: '',
        username: '',
        email: '',
        role: 'user',
        status: 'active',
        phone: '',
        password: '',
        confirmPassword: ''
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
    // 模拟从API加载数据
    setTimeout(() => {
        // 这里可以添加实际的数据加载逻辑
    }, 500);
});
</script>

<style scoped>
.user-management {
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