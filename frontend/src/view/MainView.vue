<template>
    <div class="layout-content p-6">
        <div class="mb-6">
            <h2 class="text-2xl font-bold text-gray-800">仪表盘</h2>
            <p class="text-gray-600">欢迎回来，管理员！今天是 {{ currentDate }}</p>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
            <div class="stat-card bg-white rounded-xl shadow-card p-6 flex items-center border border-gray-100">
                <div class="bg-blue-100 p-3 rounded-lg mr-4">
                    <i class="fas fa-th-large text-blue-600 text-xl"></i>
                </div>
                <div>
                    <div class="text-gray-500 text-sm">应用总数</div>
                    <div class="text-2xl font-bold text-gray-800">42</div>
                    <div class="text-xs mt-1 text-green-500"><i class="fas fa-arrow-up mr-1"></i> 较上月增长 12%
                    </div>
                </div>
            </div>

            <div class="stat-card bg-white rounded-xl shadow-card p-6 flex items-center border border-gray-100">
                <div class="bg-green-100 p-3 rounded-lg mr-4">
                    <i class="fas fa-server text-green-600 text-xl"></i>
                </div>
                <div>
                    <div class="text-gray-500 text-sm">集群数量</div>
                    <div class="text-2xl font-bold text-gray-800">28</div>
                    <div class="text-xs mt-1 text-green-500"><i class="fas fa-arrow-up mr-1"></i> 较上月增长 8%
                    </div>
                </div>
            </div>

            <div class="stat-card bg-white rounded-xl shadow-card p-6 flex items-center border border-gray-100">
                <div class="bg-amber-100 p-3 rounded-lg mr-4">
                    <i class="fas fa-folder text-amber-600 text-xl"></i>
                </div>
                <div>
                    <div class="text-gray-500 text-sm">命名空间</div>
                    <div class="text-2xl font-bold text-gray-800">156</div>
                    <div class="text-xs mt-1 text-green-500"><i class="fas fa-arrow-up mr-1"></i> 较上月增长 5%
                    </div>
                </div>
            </div>

            <div class="stat-card bg-white rounded-xl shadow-card p-6 flex items-center border border-gray-100">
                <div class="bg-purple-100 p-3 rounded-lg mr-4">
                    <i class="fas fa-users text-purple-600 text-xl"></i>
                </div>
                <div>
                    <div class="text-gray-500 text-sm">活跃用户</div>
                    <div class="text-2xl font-bold text-gray-800">1,284</div>
                    <div class="text-xs mt-1 text-green-500"><i class="fas fa-arrow-up mr-1"></i> 较上月增长 3%
                    </div>
                </div>
            </div>
        </div>

        <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
            <div class="lg:col-span-2">
                <div class="bg-white rounded-xl shadow-card p-6 border border-gray-100">
                    <div class="flex justify-between items-center mb-6">
                        <h3 class="text-lg font-semibold text-gray-800">最近更新</h3>
                        <router-link to="/recent" class="text-primary-600 text-sm hover:underline">查看全部</router-link>
                    </div>

                    <div class="space-y-4">
                        <div class="flex items-start">
                            <div class="bg-blue-100 p-2 rounded-lg mr-3">
                                <i class="fas fa-edit text-blue-600"></i>
                            </div>
                            <div class="flex-1">
                                <div class="font-medium">订单服务配置更新</div>
                                <div class="text-sm text-gray-600">更新了超时时间和重试策略</div>
                                <div class="text-xs text-gray-500 mt-1">2小时前 • 由张工更新</div>
                            </div>
                        </div>

                        <div class="flex items-start">
                            <div class="bg-green-100 p-2 rounded-lg mr-3">
                                <i class="fas fa-plus-circle text-green-600"></i>
                            </div>
                            <div class="flex-1">
                                <div class="font-medium">新增支付服务集群</div>
                                <div class="text-sm text-gray-600">添加了上海区域的支付服务集群</div>
                                <div class="text-xs text-gray-500 mt-1">5小时前 • 由李工创建</div>
                            </div>
                        </div>

                        <div class="flex items-start">
                            <div class="bg-amber-100 p-2 rounded-lg mr-3">
                                <i class="fas fa-cog text-amber-600"></i>
                            </div>
                            <div class="flex-1">
                                <div class="font-medium">用户服务配置优化</div>
                                <div class="text-sm text-gray-600">优化了缓存策略和数据库连接池</div>
                                <div class="text-xs text-gray-500 mt-1">昨天 • 由王工更新</div>
                            </div>
                        </div>

                        <div class="flex items-start">
                            <div class="bg-purple-100 p-2 rounded-lg mr-3">
                                <i class="fas fa-shield-alt text-purple-600"></i>
                            </div>
                            <div class="flex-1">
                                <div class="font-medium">安全策略更新</div>
                                <div class="text-sm text-gray-600">更新了所有服务的访问权限策略</div>
                                <div class="text-xs text-gray-500 mt-1">2天前 • 由安全团队更新</div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <div>
                <div class="bg-white rounded-xl shadow-card p-6 border border-gray-100">
                    <div class="flex justify-between items-center mb-6">
                        <h3 class="text-lg font-semibold text-gray-800">常用应用</h3>
                        <router-link to="/apps" class="text-primary-600 text-sm hover:underline">管理</router-link>
                    </div>

                    <div class="space-y-4">
                        <div class="app-card bg-gray-50 rounded-lg p-4 flex items-center border border-gray-100">
                            <div class="bg-blue-500 w-10 h-10 rounded-lg flex items-center justify-center mr-3">
                                <i class="fas fa-shopping-cart text-white"></i>
                            </div>
                            <div class="flex-1">
                                <div class="font-medium">订单服务</div>
                                <div class="text-sm text-gray-600">3个集群，12个命名空间</div>
                            </div>
                            <div class="text-primary-600">
                                <i class="fas fa-chevron-right"></i>
                            </div>
                        </div>

                        <div class="app-card bg-gray-50 rounded-lg p-4 flex items-center border border-gray-100">
                            <div class="bg-green-500 w-10 h-10 rounded-lg flex items-center justify-center mr-3">
                                <i class="fas fa-wallet text-white"></i>
                            </div>
                            <div class="flex-1">
                                <div class="font-medium">支付服务</div>
                                <div class="text-sm text-gray-600">2个集群，8个命名空间</div>
                            </div>
                            <div class="text-primary-600">
                                <i class="fas fa-chevron-right"></i>
                            </div>
                        </div>

                        <div class="app-card bg-gray-50 rounded-lg p-4 flex items-center border border-gray-100">
                            <div class="bg-purple-500 w-10 h-10 rounded-lg flex items-center justify-center mr-3">
                                <i class="fas fa-user text-white"></i>
                            </div>
                            <div class="flex-1">
                                <div class="font-medium">用户服务</div>
                                <div class="text-sm text-gray-600">4个集群，15个命名空间</div>
                            </div>
                            <div class="text-primary-600">
                                <i class="fas fa-chevron-right"></i>
                            </div>
                        </div>

                        <div class="app-card bg-gray-50 rounded-lg p-4 flex items-center border border-gray-100">
                            <div class="bg-amber-500 w-10 h-10 rounded-lg flex items-center justify-center mr-3">
                                <i class="fas fa-box text-white"></i>
                            </div>
                            <div class="flex-1">
                                <div class="font-medium">库存服务</div>
                                <div class="text-sm text-gray-600">2个集群，10个命名空间</div>
                            </div>
                            <div class="text-primary-600">
                                <i class="fas fa-chevron-right"></i>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted,computed } from 'vue'

// 当前日期
const currentDate = computed(() => {
    return new Date().toLocaleDateString('zh-CN', {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
        weekday: 'long'
    });
});

onMounted(() => {
    console.log('Component mounted.')
})
</script>

<script>
export default { name: 'MainView' }
</script>

<style scoped></style>