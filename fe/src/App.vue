<template>
  <AdminLayout
    ref="adminLayoutRef"
    :active-menu="activeMenu"
    :show-add-button="showAddButton"
    :add-button-text="addButtonText"
    @menu-click="handleMenuClick"
    @add-click="handleAddClick"
    @add-domain="handleAddDomain"
  >
    <DomainManagement
      v-if="currentPage === 'domain'"
      ref="domainManagementRef"
      @goToRecords="goToRecords"
    />
    <RecordManagement
      v-else-if="currentPage === 'records'"
      :domain="currentDomain"
      @back="backToDomains"
    />
  </AdminLayout>
</template>

<script setup>
import { ref, computed } from 'vue'
import axios from 'axios'
import AdminLayout from './components/AdminLayout.vue'
import DomainManagement from './components/DomainManagement.vue'
import RecordManagement from './components/RecordManagement.vue'

const API_BASE_URL = 'http://localhost:8080/api/v1'

// Current page state
// 当前页面状态
const currentPage = ref('domain')
const currentDomain = ref(null)

// Current active menu
// 当前激活的菜单
const activeMenu = ref('domain')

// Component refs
// 组件引用
const adminLayoutRef = ref(null)
const domainManagementRef = ref(null)

// Computed: show add button
// 计算属性：是否显示添加按钮
const showAddButton = computed(() => {
  return currentPage.value === 'domain'
})

// Computed: add button text
// 计算属性：添加按钮文本
const addButtonText = computed(() => {
  if (currentPage.value === 'domain') {
    return '添加域名'
  }
  return '添加解析'
})

// Handle menu click
// 处理菜单点击
const handleMenuClick = (key) => {
  if (key === 'domain') {
    currentPage.value = 'domain'
    currentDomain.value = null
  }
  activeMenu.value = key
}

// Handle add button click
// 处理添加按钮点击
const handleAddClick = () => {
  if (currentPage.value === 'domain') {
    adminLayoutRef.value?.openAddDomainModal()
  } else if (currentPage.value === 'records') {
    // Will be handled by RecordManagement component
    // 将由 RecordManagement 组件处理
  }
}

// Handle add domain
// 处理添加域名
const handleAddDomain = async (formData) => {
  try {
    await axios.post(`${API_BASE_URL}/domains`, formData)
    // Refresh domain list
    // 刷新域名列表
    if (domainManagementRef.value) {
      await domainManagementRef.value.fetchDomains()
    }
    alert('域名添加成功')
  } catch (error) {
    console.error('Failed to add domain:', error)
    const message = error.response?.data?.message || '添加域名失败，请稍后重试'
    alert(message)
    throw error
  }
}

// Go to records page
// 跳转到解析记录页面
const goToRecords = (domain) => {
  currentDomain.value = domain
  currentPage.value = 'records'
}

// Back to domains page
// 返回域名列表页面
const backToDomains = () => {
  currentDomain.value = null
  currentPage.value = 'domain'
}
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  background-color: #f5f7fa;
  color: #333;
}

#app {
  min-height: 100vh;
}

/* Global scrollbar styles */
/* 全局滚动条样式 */
::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 4px;
}

::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 4px;
}

::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}
</style>
