<template>
  <div class="admin-layout">
    <!-- 移动端遮罩层 -->
    <div
      v-if="isMobileMenuOpen"
      class="mobile-overlay"
      @click="closeMobileMenu"
    ></div>

    <!-- 移动端顶部栏 -->
    <header class="mobile-header">
      <button class="hamburger-btn" @click="toggleMobileMenu">
        <span class="hamburger-icon" :class="{ open: isMobileMenuOpen }">
          <span></span>
          <span></span>
          <span></span>
        </span>
      </button>
      <h1 class="mobile-title">DNS管理系统</h1>
      <button class="mobile-add-btn" @click="handleAdd">+</button>
    </header>

    <!-- 左侧导航栏 -->
    <aside class="sidebar" :class="{ collapsed: isCollapsed, 'mobile-open': isMobileMenuOpen }">
      <div class="logo">
        <h2 v-if="!isCollapsed">DNS管理系统</h2>
        <span v-else class="logo-icon">DNS</span>
      </div>
      <nav class="nav-menu">
        <div
          v-for="item in menuItems"
          :key="item.key"
          class="nav-item"
          :class="{ active: activeMenu === item.key }"
          @click="handleMenuClick(item.key)"
        >
          <span class="nav-icon">{{ item.icon }}</span>
          <span v-if="!isCollapsed" class="nav-text">{{ item.label }}</span>
        </div>
      </nav>
      <div class="sidebar-toggle" @click="toggleSidebar">
        <span>{{ isCollapsed ? '→' : '←' }}</span>
      </div>
    </aside>

    <!-- 右侧主内容区 -->
    <main class="main-content">
      <!-- 顶部栏 -->
      <header class="header">
        <div class="header-left">
          <h3>{{ currentMenuLabel }}</h3>
        </div>
        <div class="header-right" v-if="showAddButton">
          <button class="btn btn-primary" @click="handleAdd">
            <span class="btn-icon">+</span>
            {{ addButtonText }}
          </button>
        </div>
      </header>

      <!-- 内容区 -->
      <div class="content">
        <slot></slot>
      </div>
    </main>

    <!-- 添加域名弹窗 -->
    <div v-if="showModal" class="modal-overlay" @click.self="closeModal">
      <div class="modal">
        <div class="modal-header">
          <h3>添加域名</h3>
          <button class="modal-close" @click="closeModal">&times;</button>
        </div>
        <div class="modal-body">
          <form @submit.prevent="submitForm">
            <div class="form-group">
              <label for="domainName">域名名称 <span class="required">*</span></label>
              <input
                type="text"
                id="domainName"
                v-model="formData.domainName"
                required
                placeholder="请输入域名，如：example.com"
              />
            </div>
            <div class="form-group">
              <label for="description">描述</label>
              <input
                type="text"
                id="description"
                v-model="formData.description"
                placeholder="请输入域名描述"
              />
            </div>
            <div class="form-group checkbox-group">
              <label class="checkbox-label">
                <input
                  type="checkbox"
                  v-model="formData.recursive"
                />
                <span>是否递归查询</span>
              </label>
            </div>
            <div class="form-actions">
              <button type="button" class="btn btn-default" @click="closeModal">取消</button>
              <button type="submit" class="btn btn-primary" :disabled="submitting">
                {{ submitting ? '提交中...' : '确定' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'

const props = defineProps({
  activeMenu: {
    type: String,
    default: 'domain'
  },
  showAddButton: {
    type: Boolean,
    default: true
  },
  addButtonText: {
    type: String,
    default: '添加域名'
  }
})

const emit = defineEmits(['menuClick', 'addDomain', 'addClick'])

// 侧边栏折叠状态
const isCollapsed = ref(false)

// 移动端菜单状态
const isMobileMenuOpen = ref(false)

// 是否是移动端
const isMobile = ref(false)

// 检测屏幕宽度
const checkMobile = () => {
  isMobile.value = window.innerWidth <= 768
  if (!isMobile.value) {
    isMobileMenuOpen.value = false
  }
}

// 菜单项
const menuItems = [
  { key: 'domain', label: '域名管理', icon: '🌐' },
  { key: 'record', label: '记录管理', icon: '📋' },
  { key: 'settings', label: '系统设置', icon: '⚙️' }
]

// 当前菜单标签
const currentMenuLabel = computed(() => {
  const item = menuItems.find(item => item.key === props.activeMenu)
  return item ? item.label : ''
})

// 弹窗显示状态
const showModal = ref(false)

// 表单数据
const formData = ref({
  domainName: '',
  description: '',
  recursive: false
})

// 提交状态
const submitting = ref(false)

// 切换侧边栏
const toggleSidebar = () => {
  isCollapsed.value = !isCollapsed.value
}

// 处理菜单点击
const handleMenuClick = (key) => {
  emit('menuClick', key)
  // 移动端点击菜单后关闭侧边栏
  if (isMobile.value) {
    closeMobileMenu()
  }
}

// 切换移动端菜单
const toggleMobileMenu = () => {
  isMobileMenuOpen.value = !isMobileMenuOpen.value
}

// 关闭移动端菜单
const closeMobileMenu = () => {
  isMobileMenuOpen.value = false
}

// 处理添加按钮点击
const handleAdd = () => {
  emit('addClick')
  showModal.value = true
}

// Open add domain modal (exposed to parent)
// 打开添加域名弹窗（暴露给父组件）
const openAddDomainModal = () => {
  showModal.value = true
}

// 关闭弹窗
const closeModal = () => {
  showModal.value = false
  // 重置表单
  formData.value = {
    domainName: '',
    description: '',
    recursive: false
  }
}

// 提交表单
const submitForm = async () => {
  submitting.value = true
  try {
    await emit('addDomain', { ...formData.value })
    closeModal()
  } finally {
    submitting.value = false
  }
}

// 暴露方法给父组件
defineExpose({
  closeModal,
  openAddDomainModal
})

// 监听窗口大小变化
onMounted(() => {
  checkMobile()
  window.addEventListener('resize', checkMobile)
})

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
})
</script>

<style scoped>
/* 主布局 - 填满整个视口 */
.admin-layout {
  display: flex;
  height: 100vh;
  width: 100vw;
  overflow: hidden;
  background-color: #f5f7fa;
}

/* 移动端遮罩层 */
.mobile-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  z-index: 998;
  animation: fadeIn 0.2s ease;
}

/* 移动端顶部栏 */
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
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.mobile-title {
  font-size: 16px;
  font-weight: 600;
  margin: 0;
}

.hamburger-btn {
  width: 40px;
  height: 40px;
  border: none;
  background: transparent;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  transition: background-color 0.3s ease;
}

.hamburger-btn:hover {
  background-color: rgba(255, 255, 255, 0.1);
}

.hamburger-icon {
  width: 20px;
  height: 14px;
  position: relative;
}

.hamburger-icon span {
  display: block;
  position: absolute;
  height: 2px;
  width: 100%;
  background-color: #fff;
  border-radius: 2px;
  transition: all 0.3s ease;
}

.hamburger-icon span:nth-child(1) {
  top: 0;
}

.hamburger-icon span:nth-child(2) {
  top: 6px;
}

.hamburger-icon span:nth-child(3) {
  top: 12px;
}

.hamburger-icon.open span:nth-child(1) {
  transform: rotate(45deg);
  top: 6px;
}

.hamburger-icon.open span:nth-child(2) {
  opacity: 0;
}

.hamburger-icon.open span:nth-child(3) {
  transform: rotate(-45deg);
  top: 6px;
}

.mobile-add-btn {
  width: 36px;
  height: 36px;
  border: none;
  border-radius: 50%;
  background-color: #409eff;
  color: #fff;
  font-size: 24px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s ease;
}

.mobile-add-btn:hover {
  background-color: #66b1ff;
}

/* 左侧导航栏 */
.sidebar {
  width: 200px;
  min-width: 60px;
  height: 100%;
  background: linear-gradient(180deg, #1a1a2e 0%, #16213e 100%);
  color: #fff;
  display: flex;
  flex-direction: column;
  box-shadow: 2px 0 8px rgba(0, 0, 0, 0.1);
  flex-shrink: 0;
  transition: width 0.3s ease;
  position: relative;
}

.sidebar.collapsed {
  width: 60px;
}

.logo {
  padding: 16px 12px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  flex-shrink: 0;
  text-align: center;
  min-height: 56px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.logo h2 {
  font-size: 16px;
  font-weight: 600;
  color: #fff;
  margin: 0;
  white-space: nowrap;
}

.logo-icon {
  font-size: 14px;
  font-weight: 600;
  color: #fff;
}

.nav-menu {
  padding: 8px 0;
  flex: 1;
  overflow-y: auto;
}

.nav-item {
  display: flex;
  align-items: center;
  padding: 12px;
  cursor: pointer;
  transition: all 0.3s ease;
  border-left: 3px solid transparent;
  margin: 2px 6px;
  border-radius: 4px;
  justify-content: flex-start;
}

.sidebar.collapsed .nav-item {
  justify-content: center;
  padding: 12px 8px;
}

.nav-item:hover {
  background-color: rgba(255, 255, 255, 0.05);
}

.nav-item.active {
  background-color: rgba(64, 158, 255, 0.15);
  border-left-color: #409eff;
}

.nav-icon {
  font-size: 18px;
  min-width: 24px;
  text-align: center;
  flex-shrink: 0;
}

.nav-text {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.85);
  margin-left: 8px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.nav-item.active .nav-text {
  color: #409eff;
  font-weight: 500;
}

.sidebar-toggle {
  padding: 12px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
  cursor: pointer;
  text-align: center;
  color: rgba(255, 255, 255, 0.6);
  transition: all 0.3s ease;
  flex-shrink: 0;
}

.sidebar-toggle:hover {
  background-color: rgba(255, 255, 255, 0.05);
  color: #fff;
}

/* 右侧主内容区 */
.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
  min-width: 0;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 20px;
  background-color: #fff;
  border-bottom: 1px solid #e4e7ed;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.05);
  flex-shrink: 0;
}

.header-left h3 {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
  margin: 0;
}

.content {
  flex: 1;
  padding: 16px 20px;
  overflow: auto;
  background-color: #f5f7fa;
}

/* 按钮样式 */
.btn {
  display: inline-flex;
  align-items: center;
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
}

.btn-primary {
  background-color: #409eff;
  color: #fff;
}

.btn-primary:hover {
  background-color: #66b1ff;
}

.btn-primary:disabled {
  background-color: #a0cfff;
  cursor: not-allowed;
}

.btn-default {
  background-color: #fff;
  color: #606266;
  border: 1px solid #dcdfe6;
}

.btn-default:hover {
  color: #409eff;
  border-color: #c6e2ff;
  background-color: #ecf5ff;
}

.btn-icon {
  margin-right: 6px;
  font-size: 16px;
}

/* 弹窗样式 */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
  animation: fadeIn 0.2s ease;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.modal {
  background-color: #fff;
  border-radius: 8px;
  width: 480px;
  max-width: 90%;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  animation: slideUp 0.3s ease;
}

@keyframes slideUp {
  from {
    transform: translateY(20px);
    opacity: 0;
  }
  to {
    transform: translateY(0);
    opacity: 1;
  }
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid #e4e7ed;
}

.modal-header h3 {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
  margin: 0;
}

.modal-close {
  background: none;
  border: none;
  font-size: 24px;
  color: #909399;
  cursor: pointer;
  padding: 0;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  transition: all 0.3s ease;
}

.modal-close:hover {
  background-color: #f5f7fa;
  color: #409eff;
}

.modal-body {
  padding: 20px;
}

.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-size: 14px;
  color: #606266;
  font-weight: 500;
}

.form-group input[type="text"] {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  font-size: 14px;
  transition: all 0.3s ease;
}

.form-group input[type="text"]:focus {
  outline: none;
  border-color: #409eff;
  box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.2);
}

.form-group input[type="text"]::placeholder {
  color: #c0c4cc;
}

.required {
  color: #f56c6c;
}

.checkbox-group {
  display: flex;
  align-items: center;
}

.checkbox-label {
  display: flex;
  align-items: center;
  cursor: pointer;
  font-size: 14px;
  color: #606266;
}

.checkbox-label input[type="checkbox"] {
  margin-right: 8px;
  width: 16px;
  height: 16px;
  cursor: pointer;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 24px;
  padding-top: 16px;
  border-top: 1px solid #e4e7ed;
}

/* 响应式布局 */
@media (max-width: 1200px) {
  .sidebar {
    width: 180px;
  }

  .content {
    padding: 14px 16px;
  }
}

@media (max-width: 992px) {
  .sidebar {
    width: 60px;
  }

  .sidebar .nav-text {
    display: none;
  }

  .sidebar .logo h2 {
    display: none;
  }

  .sidebar .logo-icon {
    display: block;
  }

  .sidebar .nav-item {
    justify-content: center;
    padding: 12px 8px;
  }

  .sidebar .nav-icon {
    margin-right: 0;
  }

  .sidebar-toggle {
    display: none;
  }
}

@media (max-width: 768px) {
  /* 显示移动端头部 */
  .mobile-header {
    display: flex;
  }

  /* 侧边栏改为覆盖层形式 */
  .sidebar {
    position: fixed;
    top: 0;
    left: 0;
    height: 100%;
    z-index: 999;
    width: 240px;
    transform: translateX(-100%);
    transition: transform 0.3s ease;
  }

  .sidebar.mobile-open {
    transform: translateX(0);
  }

  /* 显示导航文字 */
  .sidebar .nav-text {
    display: inline;
  }

  .sidebar .logo h2 {
    display: block;
  }

  .sidebar .logo-icon {
    display: none;
  }

  .sidebar .nav-item {
    justify-content: flex-start;
    padding: 12px 16px;
  }

  .sidebar .nav-icon {
    margin-right: 8px;
  }

  /* 主内容区顶部留出空间 */
  .main-content {
    margin-top: 56px;
  }

  .header {
    padding: 10px 16px;
  }

  .header-left h3 {
    font-size: 14px;
  }

  /* 隐藏桌面端添加按钮 */
  .header-right {
    display: none;
  }

  .content {
    padding: 12px;
  }

  .btn {
    padding: 6px 12px;
    font-size: 13px;
  }
}

@media (max-width: 480px) {
  .mobile-title {
    font-size: 14px;
  }

  .content {
    padding: 8px;
  }

  /* 弹窗移动端优化 */
  .modal {
    width: 95%;
    max-width: 95%;
    max-height: 85vh;
    margin: 0 10px;
    border-radius: 8px;
  }

  .modal-header {
    padding: 14px 16px;
  }

  .modal-header h3 {
    font-size: 15px;
  }

  .modal-body {
    padding: 16px;
  }

  .form-group label {
    font-size: 13px;
  }

  .form-group input[type="text"] {
    padding: 8px 10px;
    font-size: 14px;
  }

  .form-actions {
    flex-direction: column;
    gap: 10px;
  }

  .form-actions .btn {
    width: 100%;
    justify-content: center;
  }
}
</style>
