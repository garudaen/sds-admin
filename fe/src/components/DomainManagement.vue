<template>
  <div class="domain-management">
    <!-- 标签页 -->
    <div class="tabs">
      <div
        v-for="tab in tabs"
        :key="tab.key"
        class="tab-item"
        :class="{ active: activeTab === tab.key }"
        @click="activeTab = tab.key"
      >
        {{ tab.label }}
        <span v-if="tab.key === 'all'" class="tab-badge">{{ domains.length }}</span>
      </div>
    </div>

    <!-- 搜索栏 -->
    <div class="search-bar">
      <div class="search-input-wrapper">
        <span class="search-icon">🔍</span>
        <input
          type="text"
          v-model="searchKeyword"
          placeholder="搜索域名..."
          class="search-input"
        />
      </div>
      <div class="filter-group">
        <select v-model="statusFilter" class="filter-select">
          <option value="">全部状态</option>
          <option value="active">已启用</option>
          <option value="disabled">已禁用</option>
        </select>
      </div>
    </div>

    <!-- 域名列表 -->
    <div class="table-container">
      <div v-if="loading" class="loading-wrapper">
        <div class="loading-spinner"></div>
        <p>加载中...</p>
      </div>
      
      <div v-else-if="filteredDomains.length === 0" class="empty-wrapper">
        <div class="empty-icon">📭</div>
        <p>暂无域名数据</p>
        <span class="empty-tip">点击右上角"添加域名"按钮创建</span>
      </div>
      
      <div v-else class="table-scroll-wrapper">
        <!-- 桌面端表格 -->
        <table class="data-table desktop-table">
          <thead>
            <tr>
              <th class="col-checkbox">
                <input type="checkbox" v-model="selectAll" @change="toggleSelectAll" />
              </th>
              <th class="col-id">ID</th>
              <th class="col-domain">域名名称</th>
              <th class="col-desc">描述</th>
              <th class="col-recursive">递归</th>
              <th class="col-status">状态</th>
              <th class="col-time">创建时间</th>
              <th class="col-actions">操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="domain in filteredDomains" :key="domain.id" class="table-row">
              <td class="col-checkbox">
                <input type="checkbox" v-model="selectedDomains" :value="domain.id" />
              </td>
              <td class="col-id">{{ domain.id }}</td>
              <td class="col-domain">
                <div class="domain-info">
                  <span class="domain-icon">🌐</span>
                  <span class="domain-name">{{ domain.domainName }}</span>
                </div>
              </td>
              <td class="col-desc">
                <span class="desc-text" :title="domain.description">
                  {{ domain.description || '-' }}
                </span>
              </td>
              <td class="col-recursive">
                <span class="tag" :class="domain.recursive ? 'tag-success' : 'tag-default'">
                  {{ domain.recursive ? '是' : '否' }}
                </span>
              </td>
              <td class="col-status">
                <span class="status-badge" :class="domain.status">
                  {{ domain.status === 'active' ? '已启用' : '已禁用' }}
                </span>
              </td>
              <td class="col-time">{{ formatDate(domain.createdAt) }}</td>
              <td class="col-actions">
                <div class="action-buttons">
                  <button
                    class="btn-icon-action records"
                    @click="goToRecords(domain)"
                    title="解析"
                  >
                    📋
                  </button>
                  <button
                    class="btn-icon-action edit"
                    @click="openEditModal(domain)"
                    title="编辑"
                  >
                    ✏️
                  </button>
                  <button
                    v-if="domain.status === 'active'"
                    class="btn-icon-action disable"
                    @click="disableDomain(domain.id)"
                    title="禁用"
                  >
                    🚫
                  </button>
                  <button
                    v-else
                    class="btn-icon-action enable"
                    @click="enableDomain(domain.id)"
                    title="启用"
                  >
                    ✅
                  </button>
                  <button
                    class="btn-icon-action delete"
                    @click="deleteDomain(domain.id)"
                    :disabled="deletingDomain === domain.id"
                    title="删除"
                  >
                    🗑️
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>

        <!-- 移动端卡片列表 -->
        <div class="mobile-card-list">
          <div v-for="domain in filteredDomains" :key="domain.id" class="mobile-card">
            <div class="mobile-card-header">
              <div class="mobile-card-title">
                <span class="domain-icon">🌐</span>
                <span class="domain-name">{{ domain.domainName }}</span>
              </div>
              <span class="status-badge" :class="domain.status">
                {{ domain.status === 'active' ? '已启用' : '已禁用' }}
              </span>
            </div>
            <div class="mobile-card-body">
              <div class="mobile-card-row">
                <span class="mobile-card-label">ID:</span>
                <span class="mobile-card-value">{{ domain.id }}</span>
              </div>
              <div class="mobile-card-row">
                <span class="mobile-card-label">描述:</span>
                <span class="mobile-card-value">{{ domain.description || '-' }}</span>
              </div>
              <div class="mobile-card-row">
                <span class="mobile-card-label">递归:</span>
                <span class="tag" :class="domain.recursive ? 'tag-success' : 'tag-default'">
                  {{ domain.recursive ? '是' : '否' }}
                </span>
              </div>
              <div class="mobile-card-row">
                <span class="mobile-card-label">创建时间:</span>
                <span class="mobile-card-value">{{ formatDate(domain.createdAt) }}</span>
              </div>
            </div>
            <div class="mobile-card-footer">
              <div class="mobile-card-actions">
                <button
                  class="btn-card-action records"
                  @click="goToRecords(domain)"
                >
                  📋 解析
                </button>
                <button
                  class="btn-card-action edit"
                  @click="openEditModal(domain)"
                >
                  ✏️ 编辑
                </button>
                <button
                  v-if="domain.status === 'active'"
                  class="btn-card-action disable"
                  @click="disableDomain(domain.id)"
                >
                  🚫 禁用
                </button>
                <button
                  v-else
                  class="btn-card-action enable"
                  @click="enableDomain(domain.id)"
                >
                  ✅ 启用
                </button>
                <button
                  class="btn-card-action delete"
                  @click="deleteDomain(domain.id)"
                  :disabled="deletingDomain === domain.id"
                >
                  🗑️ 删除
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 分页 -->
    <div v-if="filteredDomains.length > 0" class="pagination">
      <div class="pagination-info">
        共 {{ filteredDomains.length }} 条记录
      </div>
      <div class="pagination-actions">
        <button class="btn-page" :disabled="currentPage === 1" @click="currentPage--">
          上一页
        </button>
        <span class="page-number">第 {{ currentPage }} 页</span>
        <button class="btn-page" @click="currentPage++">
          下一页
        </button>
      </div>
    </div>

    <!-- 编辑域名弹窗 -->
    <div v-if="showEditModal" class="modal-overlay" @click.self="closeEditModal">
      <div class="modal">
        <div class="modal-header">
          <h3>编辑域名配置</h3>
          <button class="modal-close" @click="closeEditModal">&times;</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>域名名称</label>
            <input
              type="text"
              :value="editingDomain?.domainName"
              disabled
              class="input-disabled"
            />
          </div>
          <div class="form-group">
            <label for="editDescription">描述</label>
            <input
              type="text"
              id="editDescription"
              v-model="editFormData.description"
              placeholder="请输入域名描述"
            />
          </div>
          <div class="form-group checkbox-group">
            <label class="checkbox-label">
              <input
                type="checkbox"
                v-model="editFormData.recursive"
              />
              <span>是否递归查询</span>
            </label>
          </div>
          <div class="form-actions">
            <button type="button" class="btn btn-default" @click="closeEditModal">取消</button>
            <button type="button" class="btn btn-primary" :disabled="updatingConfig" @click="updateDomainConfig">
              {{ updatingConfig ? '保存中...' : '保存' }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import axios from 'axios'

const API_BASE_URL = 'http://localhost:8080/api/v1'

// 定义 emit
const emit = defineEmits(['goToRecords'])

// 标签页
const tabs = [
  { key: 'all', label: '全部域名' },
  { key: 'active', label: '已启用' },
  { key: 'disabled', label: '已禁用' }
]

// 当前激活的标签
const activeTab = ref('all')

// 搜索关键词
const searchKeyword = ref('')

// 状态筛选
const statusFilter = ref('')

// 域名列表
const domains = ref([])

// 加载状态
const loading = ref(false)

// 当前页码
const currentPage = ref(1)

// 选中项
const selectedDomains = ref([])

// 全选状态
const selectAll = ref(false)

// 删除中的域名ID
const deletingDomain = ref(null)

// 编辑弹窗状态
const showEditModal = ref(false)
const editingDomain = ref(null)
const editFormData = ref({
  description: '',
  recursive: false
})
const updatingConfig = ref(false)

// 过滤后的域名列表
const filteredDomains = computed(() => {
  let result = domains.value

  // 根据标签页筛选
  if (activeTab.value === 'active') {
    result = result.filter(d => d.status === 'active')
  } else if (activeTab.value === 'disabled') {
    result = result.filter(d => d.status === 'disabled')
  }

  // 根据状态筛选
  if (statusFilter.value) {
    result = result.filter(d => d.status === statusFilter.value)
  }

  // 根据关键词搜索
  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase()
    result = result.filter(d => 
      d.domainName.toLowerCase().includes(keyword) ||
      (d.description && d.description.toLowerCase().includes(keyword))
    )
  }

  return result
})

// 获取域名列表
const fetchDomains = async () => {
  loading.value = true
  try {
    const response = await axios.get(`${API_BASE_URL}/domains`)
    if (response.data.code === 200) {
      domains.value = response.data.data || []
    }
  } catch (error) {
    console.error('获取域名列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 删除域名
const deleteDomain = async (id) => {
  if (!confirm('确定要删除该域名吗？')) {
    return
  }

  deletingDomain.value = id
  try {
    await axios.delete(`${API_BASE_URL}/domains/${id}`)
    await fetchDomains()
  } catch (error) {
    console.error('删除域名失败:', error)
    alert('删除域名失败')
  } finally {
    deletingDomain.value = null
  }
}

// 禁用域名
const disableDomain = async (id) => {
  try {
    await axios.post(`${API_BASE_URL}/domains/${id}/disable`)
    await fetchDomains()
  } catch (error) {
    console.error('禁用域名失败:', error)
    alert('禁用域名失败')
  }
}

// 启用域名
const enableDomain = async (id) => {
  try {
    await axios.post(`${API_BASE_URL}/domains/${id}/enable`)
    await fetchDomains()
  } catch (error) {
    console.error('启用域名失败:', error)
    alert('启用域名失败')
  }
}

// 打开编辑弹窗
const openEditModal = (domain) => {
  editingDomain.value = domain
  editFormData.value = {
    description: domain.description || '',
    recursive: domain.recursive
  }
  showEditModal.value = true
}

// 关闭编辑弹窗
const closeEditModal = () => {
  showEditModal.value = false
  editingDomain.value = null
  editFormData.value = {
    description: '',
    recursive: false
  }
}

// 更新域名配置
const updateDomainConfig = async () => {
  if (!editingDomain.value) return

  updatingConfig.value = true
  try {
    await axios.patch(`${API_BASE_URL}/domains/${editingDomain.value.id}/config`, {
      description: editFormData.value.description,
      recursive: editFormData.value.recursive
    })
    await fetchDomains()
    closeEditModal()
    alert('域名配置更新成功')
  } catch (error) {
    console.error('更新域名配置失败:', error)
    const message = error.response?.data?.message || '更新域名配置失败'
    alert(message)
  } finally {
    updatingConfig.value = false
  }
}

// 格式化日期
const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 全选/取消全选
const toggleSelectAll = () => {
  if (selectAll.value) {
    selectedDomains.value = filteredDomains.value.map(d => d.id)
  } else {
    selectedDomains.value = []
  }
}

// Go to records management page.
// 跳转到解析记录管理页面
const goToRecords = (domain) => {
  emit('goToRecords', domain)
}

// 监听选中项变化
watch(selectedDomains, (newVal) => {
  selectAll.value = newVal.length === filteredDomains.value.length && newVal.length > 0
})

// 暴露方法给父组件
defineExpose({
  fetchDomains
})

// 组件挂载时获取数据
onMounted(() => {
  fetchDomains()
})
</script>

<style scoped>
/* 主容器 - 填满父容器 */
.domain-management {
  display: flex;
  flex-direction: column;
  height: 100%;
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);
  overflow: hidden;
}

/* 标签页 */
.tabs {
  display: flex;
  border-bottom: 1px solid #e4e7ed;
  padding: 0 20px;
  flex-shrink: 0;
  background-color: #fff;
}

.tab-item {
  padding: 14px 20px;
  cursor: pointer;
  font-size: 14px;
  color: #606266;
  border-bottom: 2px solid transparent;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 8px;
  white-space: nowrap;
}

.tab-item:hover {
  color: #409eff;
}

.tab-item.active {
  color: #409eff;
  border-bottom-color: #409eff;
  font-weight: 500;
}

.tab-badge {
  background-color: #f0f2f5;
  color: #909399;
  font-size: 12px;
  padding: 2px 8px;
  border-radius: 10px;
}

.tab-item.active .tab-badge {
  background-color: #ecf5ff;
  color: #409eff;
}

/* 搜索栏 */
.search-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid #e4e7ed;
  gap: 16px;
  flex-shrink: 0;
  background-color: #fff;
}

.search-input-wrapper {
  flex: 1;
  max-width: 300px;
  position: relative;
}

.search-icon {
  position: absolute;
  left: 12px;
  top: 50%;
  transform: translateY(-50%);
  color: #909399;
  font-size: 14px;
}

.search-input {
  width: 100%;
  padding: 9px 12px 9px 36px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  font-size: 14px;
  transition: all 0.3s ease;
}

.search-input:focus {
  outline: none;
  border-color: #409eff;
  box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.2);
}

.filter-select {
  padding: 9px 12px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  font-size: 14px;
  background-color: #fff;
  cursor: pointer;
  min-width: 110px;
}

.filter-select:focus {
  outline: none;
  border-color: #409eff;
}

/* 表格容器 - 占据剩余空间 */
.table-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  position: relative;
}

/* 表格滚动区域 */
.table-scroll-wrapper {
  flex: 1;
  overflow: auto;
  padding: 0;
}

/* 数据表格 */
.data-table {
  width: 100%;
  border-collapse: collapse;
  table-layout: fixed;
}

.data-table th,
.data-table td {
  padding: 12px 16px;
  text-align: left;
  border-bottom: 1px solid #ebeef5;
}

.data-table th {
  background-color: #f5f7fa;
  font-weight: 600;
  color: #606266;
  font-size: 13px;
  position: sticky;
  top: 0;
  z-index: 1;
}

.table-row:hover {
  background-color: #f5f7fa;
}

/* 列宽设置 */
.col-checkbox {
  width: 50px;
  text-align: center;
}

.col-id {
  width: 60px;
}

.col-domain {
  width: 20%;
  min-width: 160px;
}

.col-desc {
  width: 25%;
  min-width: 140px;
}

.col-recursive {
  width: 80px;
  text-align: center;
}

.col-status {
  width: 90px;
  text-align: center;
}

.col-time {
  width: 150px;
}

.col-actions {
  width: 130px;
  text-align: center;
}

/* 域名信息 */
.domain-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.domain-icon {
  font-size: 16px;
}

.domain-name {
  font-weight: 500;
  color: #303133;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* 描述文本 */
.desc-text {
  color: #606266;
  display: block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* 标签 */
.tag {
  display: inline-block;
  padding: 3px 10px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

.tag-success {
  background-color: #f0f9eb;
  color: #67c23a;
}

.tag-default {
  background-color: #f4f4f5;
  color: #909399;
}

/* 状态徽章 */
.status-badge {
  display: inline-block;
  padding: 3px 10px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

.status-badge.active {
  background-color: #f0f9eb;
  color: #67c23a;
}

.status-badge.disabled {
  background-color: #f4f4f5;
  color: #909399;
}

/* 操作按钮 */
.action-buttons {
  display: flex;
  justify-content: center;
  gap: 6px;
}

.btn-icon-action {
  width: 28px;
  height: 28px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 13px;
  transition: all 0.3s ease;
  background-color: transparent;
  display: flex;
  align-items: center;
  justify-content: center;
}

.btn-icon-action:hover {
  background-color: #f5f7fa;
}

.btn-icon-action:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-icon-action.disable:hover {
  background-color: #fef0f0;
}

.btn-icon-action.enable:hover {
  background-color: #f0f9eb;
}

.btn-icon-action.delete:hover {
  background-color: #fef0f0;
}

.btn-icon-action.edit:hover {
  background-color: #ecf5ff;
}

.btn-icon-action.records:hover {
  background-color: #fdf6ec;
}

/* 加载状态 */
.loading-wrapper {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  flex: 1;
  color: #909399;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 3px solid #f3f3f3;
  border-top-color: #409eff;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 12px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* 空状态 */
.empty-wrapper {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  flex: 1;
  color: #909399;
  padding: 40px 20px;
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.empty-wrapper p {
  font-size: 16px;
  margin-bottom: 8px;
}

.empty-tip {
  font-size: 14px;
  color: #c0c4cc;
}

/* 分页 */
.pagination {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 20px;
  border-top: 1px solid #e4e7ed;
  background-color: #fff;
  flex-shrink: 0;
}

.pagination-info {
  color: #606266;
  font-size: 14px;
}

.pagination-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.btn-page {
  padding: 7px 14px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  background-color: #fff;
  color: #606266;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.btn-page:hover:not(:disabled) {
  color: #409eff;
  border-color: #c6e2ff;
}

.btn-page:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.page-number {
  color: #606266;
  font-size: 14px;
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
  background-color: #fff;
}

.form-group input[type="text"]:focus {
  outline: none;
  border-color: #409eff;
  box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.2);
}

.form-group input[type="text"]::placeholder {
  color: #c0c4cc;
}

.form-group .input-disabled {
  background-color: #f5f7fa;
  cursor: not-allowed;
  color: #909399;
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

/* 响应式布局 */
/* 移动端卡片默认隐藏 */
.mobile-card-list {
  display: none;
}

/* 桌面端表格默认显示 */
.desktop-table {
  display: table;
}

@media (max-width: 768px) {
  .tabs {
    padding: 0 12px;
    overflow-x: auto;
  }

  .tab-item {
    padding: 12px 16px;
    font-size: 13px;
  }

  .search-bar {
    flex-direction: column;
    align-items: stretch;
    padding: 12px;
    gap: 12px;
  }

  .search-input-wrapper {
    max-width: none;
  }

  .filter-select {
    width: 100%;
  }

  /* 隐藏桌面表格，显示卡片列表 */
  .desktop-table {
    display: none;
  }

  .mobile-card-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
    padding: 12px;
  }

  /* 移动端卡片样式 */
  .mobile-card {
    background-color: #fff;
    border-radius: 8px;
    box-shadow: 0 1px 4px rgba(0, 0, 0, 0.08);
    overflow: hidden;
    border: 1px solid #ebeef5;
  }

  .mobile-card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px 16px;
    background-color: #f9fafb;
    border-bottom: 1px solid #ebeef5;
  }

  .mobile-card-title {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .mobile-card-title .domain-name {
    font-weight: 600;
    color: #303133;
    font-size: 15px;
  }

  .mobile-card-body {
    padding: 12px 16px;
  }

  .mobile-card-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 6px 0;
    border-bottom: 1px solid #f5f7fa;
  }

  .mobile-card-row:last-child {
    border-bottom: none;
  }

  .mobile-card-label {
    color: #909399;
    font-size: 13px;
    flex-shrink: 0;
  }

  .mobile-card-value {
    color: #606266;
    font-size: 13px;
    text-align: right;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    max-width: 60%;
  }

  .mobile-card-footer {
    padding: 12px 16px;
    background-color: #f9fafb;
    border-top: 1px solid #ebeef5;
  }

  .mobile-card-actions {
    display: flex;
    gap: 8px;
  }

  .btn-card-action {
    flex: 1;
    padding: 8px 12px;
    border: 1px solid #dcdfe6;
    border-radius: 4px;
    background-color: #fff;
    color: #606266;
    font-size: 13px;
    cursor: pointer;
    transition: all 0.3s ease;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 4px;
  }

  .btn-card-action:hover:not(:disabled) {
    color: #409eff;
    border-color: #c6e2ff;
    background-color: #ecf5ff;
  }

  .btn-card-action.enable:hover:not(:disabled) {
    color: #67c23a;
    border-color: #c2e7b0;
    background-color: #f0f9eb;
  }

  .btn-card-action.disable:hover:not(:disabled) {
    color: #f56c6c;
    border-color: #fbc4c4;
    background-color: #fef0f0;
  }

  .btn-card-action.delete:hover:not(:disabled) {
    color: #f56c6c;
    border-color: #fbc4c4;
    background-color: #fef0f0;
  }

  .btn-card-action.edit:hover:not(:disabled) {
    color: #409eff;
    border-color: #c6e2ff;
    background-color: #ecf5ff;
  }

  .btn-card-action.records:hover:not(:disabled) {
    color: #e6a23c;
    border-color: #f5dab1;
    background-color: #fdf6ec;
  }

  .btn-card-action:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .pagination {
    flex-direction: column;
    align-items: center;
    gap: 12px;
    padding: 12px;
  }

  .pagination-info {
    order: 2;
  }

  .pagination-actions {
    order: 1;
  }
}

@media (max-width: 480px) {
  .tab-item {
    padding: 10px 12px;
    font-size: 12px;
  }

  .search-input {
    font-size: 13px;
  }

  .empty-wrapper {
    padding: 30px 16px;
  }

  .empty-icon {
    font-size: 40px;
  }

  .btn-page {
    padding: 6px 12px;
    font-size: 12px;
  }

  .mobile-card-header,
  .mobile-card-body,
  .mobile-card-footer {
    padding-left: 12px;
    padding-right: 12px;
  }

  .btn-card-action {
    padding: 6px 8px;
    font-size: 12px;
  }

  /* 弹窗移动端适配 */
  .modal {
    width: 95%;
    max-width: 95%;
    max-height: 85vh;
    margin: 0 10px;
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
