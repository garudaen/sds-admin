<template>
  <div class="admin-layout">
      <el-container>
        <!-- 移动端遮罩 -->
        <div
          v-if="isMobileMenuOpen"
          class="mobile-overlay"
          @click="closeMobileMenu"
        />

        <!-- 移动端顶栏 -->
        <header class="mobile-header">
          <el-button
            text
            class="hamburger-btn"
            @click="toggleMobileMenu"
          >
            <el-icon :size="20">
              <component :is="isMobileMenuOpen ? 'Close' : 'Fold'" />
            </el-icon>
          </el-button>
          <span class="mobile-title">DNS 管理</span>
          <el-button
            v-if="showAddButton"
            type="primary"
            circle
            class="mobile-add-btn"
            @click="handleAdd"
          >
            <el-icon><Plus /></el-icon>
          </el-button>
        </header>

        <!-- 侧边栏 -->
        <el-aside
          class="sidebar"
          :class="{ collapsed: isCollapsed, 'mobile-open': isMobileMenuOpen }"
          :width="sidebarWidth"
        >
          <div class="logo">
            <span v-if="!isCollapsed">DNS 管理</span>
            <span v-else class="logo-icon">DNS</span>
          </div>
          <el-menu
            :default-active="activeMenu"
            :collapse="isCollapsed"
            background-color="#1a1a2e"
            text-color="rgba(255,255,255,0.85)"
            active-text-color="#409eff"
            router
            @select="onMenuSelect"
          >
            <el-menu-item index="/domains">
              <el-icon><Odometer /></el-icon>
              <template #title>域名管理</template>
            </el-menu-item>
          </el-menu>
          <div class="sidebar-toggle" @click="toggleSidebar">
            <el-icon><DArrowLeft v-if="!isCollapsed" /><DArrowRight v-else /></el-icon>
          </div>
        </el-aside>

        <!-- 主内容 -->
        <el-container class="main-wrap">
          <el-header class="main-header">
            <h3 class="header-title">{{ headerTitle }}</h3>
            <el-button
              v-if="showAddButton"
              type="primary"
              @click="handleAdd"
            >
              <el-icon><Plus /></el-icon>
              {{ addButtonText }}
            </el-button>
          </el-header>
          <el-main class="main-content">
            <router-view v-slot="{ Component }">
              <component :is="Component" />
            </router-view>
          </el-main>
        </el-container>
      </el-container>

      <!-- 添加域名弹窗 -->
      <el-dialog
        v-model="showModal"
        title="添加域名"
        width="480px"
        destroy-on-close
        @close="closeModal"
      >
        <el-form
          ref="formRef"
          :model="formData"
          label-width="120px"
          @submit.prevent="submitForm"
        >
          <el-form-item label="域名名称" required>
            <el-input
              v-model="formData.domainName"
              placeholder="如：example.com"
              clearable
            />
          </el-form-item>
          <el-form-item label="描述">
            <el-input
              v-model="formData.description"
              placeholder="请输入域名描述"
              clearable
            />
          </el-form-item>
          <el-form-item label="递归查询">
            <el-switch v-model="formData.recursive" />
          </el-form-item>
        </el-form>
        <template #footer>
          <el-button @click="closeModal">取消</el-button>
          <el-button type="primary" :loading="submitting" @click="submitForm">
            确定
          </el-button>
        </template>
      </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useStore } from 'vuex'
import { Close, Fold, Plus, Odometer, DArrowLeft, DArrowRight } from '@element-plus/icons-vue'
import request from '../api/request'
import { ElMessage } from 'element-plus'

const route = useRoute()
const router = useRouter()
const store = useStore()

const isCollapsed = ref(false)
const isMobileMenuOpen = ref(false)
const isMobile = ref(false)
const showModal = ref(false)
const formRef = ref(null)
const submitting = ref(false)

const formData = ref({
  domainName: '',
  description: '',
  recursive: false,
})

const checkMobile = () => {
  isMobile.value = window.innerWidth <= 768
  if (!isMobile.value) isMobileMenuOpen.value = false
}

const sidebarWidth = computed(() => {
  if (isMobile.value) return '240px'
  return isCollapsed.value ? '64px' : '200px'
})

const activeMenu = computed(() => {
  if (route.path.startsWith('/domains/') && route.path.endsWith('/records')) return '/domains'
  return route.path
})

const headerTitle = computed(() => {
  if (route.path.startsWith('/domains/') && route.path.endsWith('/records')) return '解析记录'
  return '域名管理'
})

const showAddButton = computed(() => {
  if (route.path.startsWith('/domains/') && route.path.endsWith('/records')) return true
  return route.path === '/domains' || route.path === '/'
})

const addButtonText = computed(() => {
  if (route.path.startsWith('/domains/') && route.path.endsWith('/records')) return '添加解析'
  return '添加域名'
})

function toggleSidebar() {
  isCollapsed.value = !isCollapsed.value
}

function onMenuSelect(index) {
  if (index === '/domains') {
    store.dispatch('app/setActiveMenu', 'domain')
    router.push('/domains')
  }
  if (isMobile.value) closeMobileMenu()
}

function toggleMobileMenu() {
  isMobileMenuOpen.value = !isMobileMenuOpen.value
}

function closeMobileMenu() {
  isMobileMenuOpen.value = false
}

function handleAdd() {
  if (route.path.startsWith('/domains/') && route.path.endsWith('/records')) {
    // 添加解析由 DomainRecords 内部处理，通过 provide/inject 或事件
    window.dispatchEvent(new CustomEvent('layout-add-record'))
    return
  }
  showModal.value = true
}

function closeModal() {
  showModal.value = false
  formData.value = { domainName: '', description: '', recursive: false }
}

async function submitForm() {
  if (!formData.value.domainName?.trim()) {
    return
  }
  submitting.value = true
  try {
    await request.post('/domains', formData.value)
    ElMessage.success('域名添加成功')
    closeModal()
    window.dispatchEvent(new CustomEvent('domains-refresh'))
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  checkMobile()
  window.addEventListener('resize', checkMobile)
})

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
})
</script>

<style scoped>
.admin-layout {
  height: 100vh;
  background: #f5f7fa;
}

.mobile-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  z-index: 998;
}

.mobile-header {
  display: none;
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 56px;
  background: linear-gradient(180deg, #1a1a2e 0%, #16213e 100%);
  color: #fff;
  align-items: center;
  justify-content: space-between;
  padding: 0 12px;
  z-index: 997;
}

.mobile-title {
  font-size: 16px;
  font-weight: 600;
}

.hamburger-btn {
  color: #fff;
}

.mobile-add-btn {
  color: #fff;
}

.sidebar {
  background: linear-gradient(180deg, #1a1a2e 0%, #16213e 100%);
  transition: width 0.2s;
}

.sidebar.collapsed .sidebar-toggle .el-icon {
  transform: rotate(180deg);
}

.sidebar.mobile-open {
  transform: translateX(0);
}

.logo {
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-weight: 600;
  border-bottom: 1px solid rgba(255,255,255,0.1);
}

.logo-icon {
  font-size: 14px;
}

.sidebar-toggle {
  padding: 12px;
  border-top: 1px solid rgba(255,255,255,0.1);
  color: rgba(255,255,255,0.6);
  text-align: center;
  cursor: pointer;
}

.sidebar-toggle:hover {
  background: rgba(255,255,255,0.05);
  color: #fff;
}

.main-wrap {
  flex-direction: column;
  min-width: 0;
  overflow: hidden;
}

.main-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: #fff;
  border-bottom: 1px solid #e4e7ed;
  height: 56px;
  padding: 0 20px;
}

.header-title {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.main-content {
  flex: 1;
  padding: 16px 20px;
  overflow: auto;
  background: #f5f7fa;
}

@media (max-width: 768px) {
  .mobile-header {
    display: flex;
  }

  .sidebar {
    position: fixed;
    top: 0;
    left: 0;
    bottom: 0;
    z-index: 999;
    transform: translateX(-100%);
    transition: transform 0.2s;
  }

  .sidebar.mobile-open {
    transform: translateX(0);
  }

  .main-wrap {
    margin-top: 56px;
  }
}
</style>
