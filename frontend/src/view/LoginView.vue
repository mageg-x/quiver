<template>
    <div class="login-container">
        <!-- Animated Background -->
        <div class="background-animation">
            <div class="floating-shape shape-1"></div>
            <div class="floating-shape shape-2"></div>
            <div class="floating-shape shape-3"></div>
            <div class="floating-shape shape-4"></div>
            <div class="floating-shape shape-5"></div>
        </div>

        <div class="login-content">
            <div class="login-card">
                <!-- Header -->
                <div class="login-header">
                    <div class="logo-container">
                        <div class="logo">
                            <img src="@/assets/img/favicon.svg" alt="Logo" class="logo-image">
                        </div>
                    </div>
                    <h2 class="title">配置中心管理系统</h2>
                    <p class="subtitle">请登录您的账户以继续</p>
                </div>

                <!-- Form -->
                <form @submit.prevent="handleLogin" class="login-form">
                    <div class="form-group">
                        <label for="username" class="form-label">用户名</label>
                        <div class="input-wrapper">
                            <i class="pi pi-user input-icon"></i>
                            <InputText id="username" v-model="form.username" placeholder="请输入用户名" class="form-input  !pl-10"
                                :class="{ 'error': errors.username }" required />
                        </div>
                        <Transition name="error">
                            <small v-if="errors.username" class="error-message">{{ errors.username }}</small>
                        </Transition>
                    </div>

                    <div class="form-group">
                        <label for="password" class="form-label">密码</label>
                        <div class="input-wrapper">
                            <i class="pi pi-lock input-icon"></i>
                            <Password id="password" v-model="form.password" placeholder="请输入密码" class="form-input py-0 px-10"
                                :class="{ 'error': errors.password }" :feedback="false" toggleMask required />
                        </div>
                        <Transition name="error">
                            <small v-if="errors.password" class="error-message">{{ errors.password }}</small>
                        </Transition>
                    </div>

                    <div class="form-group">
                        <label for="environment" class="form-label">环境</label>
                        <div class="input-wrapper">
                            <i class="pi pi-server input-icon"></i>
                            <Dropdown id="environment" v-model="form.environment" :options="environments"
                                optionLabel="label" optionValue="value" placeholder="选择环境"
                                class="form-input environment-dropdown pl-10" />
                        </div>
                    </div>

                    <Button type="submit" label="登录" class="login-button" :loading="loading" :disabled="loading">
                        <template #icon>
                            <i class="pi pi-sign-in"></i>
                        </template>
                    </Button>
                </form>

                <!-- Footer -->
                <div class="login-footer">
                    <p class="footer-text">
                        遇到问题？
                        <a href="#" class="footer-link">联系管理员</a>
                    </p>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import InputText from 'primevue/inputtext'
import Password from 'primevue/password'
import Dropdown from 'primevue/dropdown'
import Button from 'primevue/button'

const environments = [
    { label: '开发环境', value: 'dev' },
    { label: '生产环境', value: 'pro' }
]

const router = useRouter()
const toast = useToast()
const loading = ref(false)

const form = reactive({
    username: '',
    password: '',
    environment: 'dev'
})

const errors = reactive({
    username: '',
    password: ''
})

const validateForm = () => {
    errors.username = ''
    errors.password = ''

    if (!form.username.trim()) {
        errors.username = '请输入用户名'
        return false
    }

    if (!form.password.trim()) {
        errors.password = '请输入密码'
        return false
    }

    return true
}

const handleLogin = async () => {
    if (!validateForm()) return

    loading.value = true

    try {
        // Simulate API call
        await new Promise(resolve => setTimeout(resolve, 1500))

        toast.add({
            severity: 'success',
            summary: '登录成功',
            detail: '欢迎使用配置中心管理系统',
            life: 3000
        })

        router.push('/')
    } catch (error) {
        toast.add({
            severity: 'error',
            summary: '登录失败',
            detail: '网络错误，请稍后重试',
            life: 5000
        })
    } finally {
        loading.value = false
    }
}
</script>

<script>
export default { name: 'LoginView' }
</script>

<style scoped>
::v-deep(.p-password-input) {
    width: 100%;
    border: 0px;
}

.login-container {
    position: relative;
    min-height: 100vh;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 2rem 1rem;
    overflow: hidden;
}

.background-animation {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    pointer-events: none;
    z-index: 1;
}

.floating-shape {
    position: absolute;
    border-radius: 50%;
    background: rgba(255, 255, 255, 0.1);
    backdrop-filter: blur(10px);
    animation: float 6s ease-in-out infinite;
}

.shape-1 {
    width: 80px;
    height: 80px;
    top: 10%;
    left: 10%;
    animation-delay: 0s;
}

.shape-2 {
    width: 120px;
    height: 120px;
    top: 20%;
    right: 10%;
    animation-delay: 2s;
}

.shape-3 {
    width: 60px;
    height: 60px;
    bottom: 30%;
    left: 5%;
    animation-delay: 4s;
}

.shape-4 {
    width: 100px;
    height: 100px;
    bottom: 10%;
    right: 20%;
    animation-delay: 1s;
}

.shape-5 {
    width: 40px;
    height: 40px;
    top: 50%;
    left: 50%;
    animation-delay: 3s;
}

@keyframes float {
    0%,
    100% {
        transform: translateY(0px) rotate(0deg);
        opacity: 0.7;
    }

    50% {
        transform: translateY(-20px) rotate(180deg);
        opacity: 1;
    }
}

.login-content {
    position: relative;
    z-index: 2;
    width: 100%;
    max-width: 420px;
}

.login-card {
    background: rgba(255, 255, 255, 0.95);
    backdrop-filter: blur(20px);
    border: 1px solid rgba(255, 255, 255, 0.2);
    border-radius: 24px;
    padding: 2rem 2rem;
    box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
    animation: slideUp 0.6s ease-out;
}

@keyframes slideUp {
    from {
        opacity: 0;
        transform: translateY(30px);
    }

    to {
        opacity: 1;
        transform: translateY(0);
    }
}

.login-header {
    text-align: center;
    margin-bottom: 1rem;
}

.logo-container {
    margin-bottom: 1rem;
}

.logo {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 80px;
    height: 80px;
    background: linear-gradient(135deg, #3B82F6 0%, #1D4ED8 100%);
    border-radius: 20px;
    box-shadow: 0 10px 30px rgba(59, 130, 246, 0.3);
    animation: logoFloat 3s ease-in-out infinite;
}

@keyframes logoFloat {
    0%,
    100% {
        transform: translateY(0px);
    }

    50% {
        transform: translateY(-5px);
    }
}

.logo i {
    font-size: 2rem;
    color: white;
}

.title {
    font-size: 1.75rem;
    font-weight: 700;
    color: #1F2937;
    margin-bottom: 0.5rem;
    letter-spacing: -0.025em;
}

.subtitle {
    font-size: 0.95rem;
    color: #6B7280;
    font-weight: 400;
}

.login-form {
    margin-bottom: 1rem;
}

.form-group {
    margin-bottom: 1rem;
}

.form-label {
    display: block;
    font-size: 0.875rem;
    font-weight: 600;
    color: #374151;
    margin-bottom: 0.5rem;
}

.input-wrapper {
    position: relative;
    display: flex;
    align-items: center;
}

.input-icon {
    position: absolute;
    left: 1rem;
    color: #9CA3AF;
    z-index: 10;
    transition: color 0.2s ease;
}

.form-input {
    width: 100%;
    /* padding: 0.875rem 1rem 0.875rem 2.75rem; */
    border: 2px solid #E5E7EB;
    border-radius: 12px;
    font-size: 0.95rem;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    background: rgba(255, 255, 255, 0.8);
    backdrop-filter: blur(10px);
}

.form-input:focus {
    outline: none;
    border-color: #3B82F6;
    background: rgba(255, 255, 255, 0.95);
    box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.form-input:focus+.input-icon,
.input-wrapper:focus-within .input-icon {
    color: #3B82F6;
}

.form-input.error {
    border-color: #EF4444;
    background: rgba(254, 242, 242, 0.8);
}

.environment-dropdown {
    border: 2px solid #E5E7EB;
    border-radius: 12px;
}

.error-message {
    display: block;
    color: #EF4444;
    font-size: 0.8rem;
    margin-top: 0.5rem;
    font-weight: 500;
}

.error-enter-active {
    transition: all 0.3s ease-out;
}

.error-leave-active {
    transition: all 0.3s ease-in;
}

.error-enter-from,
.error-leave-to {
    opacity: 0;
    transform: translateY(-10px);
}

.login-button {
    width: 100%;
    padding: 1rem;
    background: linear-gradient(135deg, #3B82F6 0%, #1D4ED8 100%);
    border: none;
    border-radius: 12px;
    color: white;
    font-weight: 600;
    font-size: 1rem;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    box-shadow: 0 10px 30px rgba(59, 130, 246, 0.3);
    position: relative;
    overflow: hidden;
}

.login-button:hover:not(:disabled) {
    transform: translateY(-2px);
    box-shadow: 0 15px 35px rgba(59, 130, 246, 0.4);
}

.login-button:active {
    transform: translateY(0);
}

.login-button:disabled {
    opacity: 0.7;
    cursor: not-allowed;
}

.login-button::before {
    content: '';
    position: absolute;
    top: 0;
    left: -100%;
    width: 100%;
    height: 100%;
    background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
    transition: left 0.5s ease;
}

.login-button:hover::before {
    left: 100%;
}

.login-footer {
    text-align: center;
    border-top: 1px solid rgba(229, 231, 235, 0.5);
    padding-top: 1.5rem;
}

.footer-text {
    font-size: 0.875rem;
    color: #6B7280;
}

.footer-link {
    color: #3B82F6;
    text-decoration: none;
    font-weight: 500;
    transition: color 0.2s ease;
}

.footer-link:hover {
    color: #1D4ED8;
}

/* Mobile Responsiveness */
@media (max-width: 640px) {
    .login-container {
        padding: 1rem;
    }

    .login-card {
        padding: 2rem 1.5rem;
        border-radius: 20px;
    }

    .title {
        font-size: 1.5rem;
    }

    .logo {
        width: 60px;
        height: 60px;
    }

    .logo i {
        font-size: 1.5rem;
    }

    .floating-shape {
        display: none;
    }
}

/* Dark mode compatibility */
@media (prefers-color-scheme: dark) {
    .login-card {
        background: rgba(17, 24, 39, 0.95);
        border: 1px solid rgba(75, 85, 99, 0.3);
    }

    .title {
        color: #F9FAFB;
    }

    .subtitle {
        color: #D1D5DB;
    }

    .form-label {
        color: #E5E7EB;
    }

    .form-input {
        background: rgba(31, 41, 55, 0.8);
        border-color: #374151;
        color: #F9FAFB;
    }

    .form-input:focus {
        background: rgba(31, 41, 55, 0.95);
    }

    .footer-text {
        color: #9CA3AF;
    }
}
</style>