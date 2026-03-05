<template>
  <div class="record-management">
    <!-- Header with domain info and back button -->
    <!-- 头部：域名信息和返回按钮 -->
    <div class="page-header">
      <div class="header-left">
        <button class="btn-back" @click="goBack">
          <span class="back-icon">←</span>
          <span>返回域名列表</span>
        </button>
        <div class="domain-info" v-if="domain">
          <span class="domain-icon">🌐</span>
          <span class="domain-name">{{ domain.domainName }}</span>
          <span class="status-badge" :class="domain.status">
            {{ domain.status === 'active' ? '已启用' : '已禁用' }}
          </span>
        </div>
      </div>
      <div class="header-right">
        <button class="btn btn-primary" @click="openAddModal">
          <span class="btn-icon">+</span>
          添加解析
        </button>
      </div>
    </div>

    <!-- Search bar -->
    <!-- 搜索栏 -->
    <div class="search-bar">
      <div class="search-input-wrapper">
        <span class="search-icon">🔍</span>
        <input
          type="text"
          v-model="searchKeyword"
          placeholder="搜索主机记录..."
          class="search-input"
        />
      </div>
      <div class="filter-group">
        <select v-model="typeFilter" class="filter-select">
          <option value="">全部类型</option>
          <option v-for="rt in recordTypes" :key="rt.id" :value="rt.id">
            {{ rt.name }}
          </option>
        </select>
        <select v-model="statusFilter" class="filter-select">
          <option value="">全部状态</option>
          <option value="enabled">已启用</option>
          <option value="disabled">已禁用</option>
        </select>
      </div>
    </div>

    <!-- Record list -->
    <!-- 记录列表 -->
    <div class="table-container">
      <div v-if="loading" class="loading-wrapper">
        <div class="loading-spinner"></div>
        <p>加载中...</p>
      </div>

      <div v-else-if="filteredRecords.length === 0" class="empty-wrapper">
        <div class="empty-icon">📭</div>
        <p>暂无解析记录</p>
        <span class="empty-tip">点击右上角"添加解析"按钮创建</span>
      </div>

      <div v-else class="table-scroll-wrapper">
        <!-- Desktop table -->
        <!-- 桌面端表格 -->
        <table class="data-table desktop-table">
          <thead>
            <tr>
              <th class="col-host">主机记录</th>
              <th class="col-type">记录类型</th>
              <th class="col-value">记录值</th>
              <th class="col-ttl">TTL</th>
              <th class="col-status">状态</th>
              <th class="col-remark">备注</th>
              <th class="col-actions">操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="record in filteredRecords" :key="record.id" class="table-row">
              <td class="col-host">
                <span class="host-name">{{ record.host || '@' }}</span>
              </td>
              <td class="col-type">
                <span class="type-tag" :class="record.recordType.toLowerCase()">
                  {{ record.recordType }}
                </span>
              </td>
              <td class="col-value">
                <div class="value-list">
                  <div v-for="v in record.values" :key="v.id" class="value-item">
                    <span class="value-text">{{ v.value }}</span>
                    <span v-if="v.isDefault" class="default-badge">默认</span>
                    <span v-if="v.clientCidr" class="cidr-badge">{{ v.clientCidr }}</span>
                  </div>
                </div>
              </td>
              <td class="col-ttl">{{ record.ttl }}s</td>
              <td class="col-status">
                <span class="status-badge" :class="record.disabled ? 'disabled' : 'active'">
                  {{ record.disabled ? '已禁用' : '已启用' }}
                </span>
              </td>
              <td class="col-remark">
                <span class="remark-text" :title="record.remark">
                  {{ record.remark || '-' }}
                </span>
              </td>
              <td class="col-actions">
                <div class="action-buttons">
                  <button
                    class="btn-icon-action edit"
                    @click="openEditModal(record)"
                    title="编辑"
                  >
                    ✏️
                  </button>
                  <button
                    v-if="!record.disabled"
                    class="btn-icon-action disable"
                    @click="disableRecord(record.id)"
                    title="禁用"
                  >
                    🚫
                  </button>
                  <button
                    v-else
                    class="btn-icon-action enable"
                    @click="enableRecord(record.id)"
                    title="启用"
                  >
                    ✅
                  </button>
                  <button
                    class="btn-icon-action delete"
                    @click="deleteRecord(record.id)"
                    :disabled="deletingRecord === record.id"
                    title="删除"
                  >
                    🗑️
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>

        <!-- Mobile card list -->
        <!-- 移动端卡片列表 -->
        <div class="mobile-card-list">
          <div v-for="record in filteredRecords" :key="record.id" class="mobile-card">
            <div class="mobile-card-header">
              <div class="mobile-card-title">
                <span class="type-tag" :class="record.recordType.toLowerCase()">
                  {{ record.recordType }}
                </span>
                <span class="host-name">{{ record.host || '@' }}</span>
              </div>
              <span class="status-badge" :class="record.disabled ? 'disabled' : 'active'">
                {{ record.disabled ? '已禁用' : '已启用' }}
              </span>
            </div>
            <div class="mobile-card-body">
              <div class="mobile-card-row">
                <span class="mobile-card-label">记录值:</span>
                <div class="mobile-card-value">
                  <div v-for="v in record.values" :key="v.id" class="value-item-small">
                    {{ v.value }}
                    <span v-if="v.isDefault" class="default-badge-small">默认</span>
                  </div>
                </div>
              </div>
              <div class="mobile-card-row">
                <span class="mobile-card-label">TTL:</span>
                <span class="mobile-card-value">{{ record.ttl }}s</span>
              </div>
              <div class="mobile-card-row">
                <span class="mobile-card-label">备注:</span>
                <span class="mobile-card-value">{{ record.remark || '-' }}</span>
              </div>
            </div>
            <div class="mobile-card-footer">
              <div class="mobile-card-actions">
                <button class="btn-card-action edit" @click="openEditModal(record)">
                  ✏️ 编辑
                </button>
                <button
                  v-if="!record.disabled"
                  class="btn-card-action disable"
                  @click="disableRecord(record.id)"
                >
                  🚫 禁用
                </button>
                <button
                  v-else
                  class="btn-card-action enable"
                  @click="enableRecord(record.id)"
                >
                  ✅ 启用
                </button>
                <button
                  class="btn-card-action delete"
                  @click="deleteRecord(record.id)"
                  :disabled="deletingRecord === record.id"
                >
                  🗑️ 删除
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Add/Edit Record Modal -->
    <!-- 添加/编辑记录弹窗 -->
    <div v-if="showModal" class="modal-overlay" @click.self="closeModal">
      <div class="modal modal-large">
        <div class="modal-header">
          <h3>{{ isEdit ? '编辑解析记录' : '添加解析记录' }}</h3>
          <button class="modal-close" @click="closeModal">&times;</button>
        </div>
        <div class="modal-body">
          <form @submit.prevent="submitForm">
            <div class="form-row">
              <div class="form-group">
                <label>记录类型 <span class="required">*</span></label>
                <select
                  v-model="formData.recordTypeId"
                  :disabled="isEdit"
                  required
                  @change="onRecordTypeChange"
                >
                  <option value="">请选择记录类型</option>
                  <option v-for="rt in recordTypes" :key="rt.id" :value="rt.id">
                    {{ rt.name }}
                  </option>
                </select>
              </div>
              <div class="form-group">
                <label>主机记录 <span class="required">*</span></label>
                <input
                  type="text"
                  v-model="formData.host"
                  placeholder="如：www, @, *"
                  required
                />
              </div>
            </div>

            <div class="form-row">
              <div class="form-group">
                <label>TTL (秒)</label>
                <input
                  type="number"
                  v-model.number="formData.ttl"
                  placeholder="默认 300"
                  min="1"
                />
              </div>
            </div>

            <div class="form-group">
              <label>备注</label>
              <input
                type="text"
                v-model="formData.remark"
                placeholder="请输入备注"
              />
            </div>

            <!-- Record values section -->
            <!-- 记录值区域 -->
            <div class="form-group">
              <label>
                记录值 <span class="required">*</span>
                <span v-if="needDefaultValue" class="form-tip">（A/AAAA/CNAME 必须有一个默认值）</span>
                <span v-if="isMXType" class="form-tip">（每个记录值可设置不同优先级）</span>
              </label>
              <div class="values-container">
                <div
                  v-for="(value, index) in formData.values"
                  :key="index"
                  class="value-input-row"
                >
                  <!-- MX 记录：优先级 + 记录值 -->
                  <template v-if="isMXType">
                    <input
                      type="number"
                      v-model.number="value.mxPriority"
                      placeholder="优先级"
                      class="mx-priority-input"
                      min="0"
                      max="65535"
                    />
                    <input
                      type="text"
                      v-model="value.value"
                      placeholder="请输入邮件服务器地址"
                      required
                      class="value-input-mx"
                    />
                  </template>
                  <!-- A/AAAA/CNAME：默认开关 + 记录值 + CIDR -->
                  <template v-else-if="needDefaultValue">
                    <label class="checkbox-inline">
                      <input
                        type="checkbox"
                        v-model="value.isDefault"
                        @change="onDefaultChange(index)"
                      />
                      <span>默认</span>
                    </label>
                    <input
                      type="text"
                      v-model="value.value"
                      :placeholder="getValuePlaceholder"
                      required
                      class="value-input"
                    />
                    <input
                      type="text"
                      v-model="value.clientCidr"
                      placeholder="CIDR (可选)"
                      class="cidr-input"
                      :disabled="value.isDefault"
                    />
                  </template>
                  <!-- TXT等其他记录：仅记录值 -->
                  <template v-else>
                    <input
                      type="text"
                      v-model="value.value"
                      :placeholder="getValuePlaceholder"
                      required
                      class="value-input-full"
                    />
                  </template>
                  <button
                    type="button"
                    class="btn-remove-value"
                    @click="removeValue(index)"
                    :disabled="formData.values.length <= 1"
                  >
                    ×
                  </button>
                </div>
                <button type="button" class="btn-add-value" @click="addValue">
                  + 添加记录值
                </button>
              </div>
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
import { ref, computed, onMounted, watch } from 'vue'
import axios from 'axios'

const API_BASE_URL = 'http://localhost:8080/api/v1'

// Props
const props = defineProps({
  domain: {
    type: Object,
    default: null
  }
})

// Emits
const emit = defineEmits(['back'])

// Data
const records = ref([])
const recordTypes = ref([])
const loading = ref(false)
const searchKeyword = ref('')
const typeFilter = ref('')
const statusFilter = ref('')
const deletingRecord = ref(null)

// Modal state
const showModal = ref(false)
const isEdit = ref(false)
const editingRecord = ref(null)
const submitting = ref(false)

// Form data
const formData = ref({
  recordTypeId: '',
  host: '',
  ttl: 300,
  remark: '',
  values: [{ value: '', isDefault: true, clientCidr: '', mxPriority: null }]
})

// Computed
const filteredRecords = computed(() => {
  let result = records.value

  // Filter by type
  if (typeFilter.value) {
    result = result.filter(r => r.recordTypeId === typeFilter.value)
  }

  // Filter by status
  if (statusFilter.value === 'enabled') {
    result = result.filter(r => !r.disabled)
  } else if (statusFilter.value === 'disabled') {
    result = result.filter(r => r.disabled)
  }

  // Filter by keyword
  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase()
    result = result.filter(r =>
      r.host.toLowerCase().includes(keyword) ||
      r.values.some(v => v.value.toLowerCase().includes(keyword))
    )
  }

  return result
})

// Check if current record type needs default value
// 检查当前记录类型是否需要默认值
const needDefaultValue = computed(() => {
  const rt = recordTypes.value.find(r => r.id == formData.value.recordTypeId)
  return rt && ['A', 'AAAA', 'CNAME'].includes(rt.name)
})

// Check if current record type is MX
// 检查当前记录类型是否为MX
const isMXType = computed(() => {
  const rt = recordTypes.value.find(r => r.id == formData.value.recordTypeId)
  return rt && rt.name === 'MX'
})

// Check if show CIDR input (only for A/AAAA/CNAME)
// 检查是否显示CIDR输入框（仅A/AAAA/CNAME类型）
const showCIDR = computed(() => {
  const rt = recordTypes.value.find(r => r.id == formData.value.recordTypeId)
  return rt && ['A', 'AAAA', 'CNAME'].includes(rt.name)
})

// Check if show MX priority field
// 检查是否显示MX优先级字段（已废弃，使用isMXType）
const showMXPriority = computed(() => {
  return isMXType.value
})

// Get value input placeholder
// 获取值输入框占位符
const getValuePlaceholder = computed(() => {
  const rt = recordTypes.value.find(r => r.id == formData.value.recordTypeId)
  if (!rt) return '请输入记录值'
  switch (rt.name) {
    case 'A': return '请输入IPv4地址，如：192.168.1.1'
    case 'AAAA': return '请输入IPv6地址'
    case 'CNAME': return '请输入目标域名'
    case 'MX': return '请输入邮件服务器地址'
    default: return '请输入记录值'
  }
})

// Methods
const fetchRecords = async () => {
  if (!props.domain) return
  loading.value = true
  try {
    const response = await axios.get(`${API_BASE_URL}/domains/${props.domain.id}/records`)
    if (response.data.code === 200) {
      records.value = response.data.data.records || []
    }
  } catch (error) {
    console.error('Failed to fetch records:', error)
  } finally {
    loading.value = false
  }
}

const fetchRecordTypes = async () => {
  try {
    const response = await axios.get(`${API_BASE_URL}/record-types`)
    if (response.data.code === 200) {
      recordTypes.value = response.data.data || []
    }
  } catch (error) {
    console.error('Failed to fetch record types:', error)
  }
}

const goBack = () => {
  emit('back')
}

const openAddModal = () => {
  isEdit.value = false
  editingRecord.value = null
  formData.value = {
    recordTypeId: '',
    host: '',
    ttl: 300,
    remark: '',
    values: [{ value: '', isDefault: true, clientCidr: '', mxPriority: null }]
  }
  showModal.value = true
}

const openEditModal = (record) => {
  isEdit.value = true
  editingRecord.value = record
  formData.value = {
    recordTypeId: record.recordTypeId,
    host: record.host,
    ttl: record.ttl,
    remark: record.remark || '',
    values: record.values.map(v => ({
      value: v.value,
      isDefault: v.isDefault,
      clientCidr: v.clientCidr || '',
      mxPriority: v.mxPriority
    }))
  }
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
  isEdit.value = false
  editingRecord.value = null
}

const onRecordTypeChange = () => {
  // Reset values when type changes
  // 类型改变时重置值
  formData.value.values.forEach((v, i) => {
    // Reset default checkbox for A/AAAA/CNAME
    v.isDefault = i === 0 && needDefaultValue.value
    if (v.isDefault) {
      v.clientCidr = ''
    }
    // Clear MX priority if not MX type
    if (!isMXType.value) {
      v.mxPriority = null
    }
  })
}

const onDefaultChange = (index) => {
  // Only one default value allowed
  // 只允许一个默认值
  if (formData.value.values[index].isDefault) {
    formData.value.values.forEach((v, i) => {
      if (i !== index) {
        v.isDefault = false
      }
    })
    // Clear CIDR for default value
    // 默认值清除CIDR
    formData.value.values[index].clientCidr = ''
  }
}

const addValue = () => {
  formData.value.values.push({
    value: '',
    isDefault: false,
    clientCidr: ''
  })
}

const removeValue = (index) => {
  if (formData.value.values.length > 1) {
    formData.value.values.splice(index, 1)
  }
}

const submitForm = async () => {
  // Validate
  // 验证
  if (needDefaultValue.value && !formData.value.values.some(v => v.isDefault)) {
    alert('A/AAAA/CNAME 记录必须有一个默认值')
    return
  }

  if (isMXType.value && !formData.value.values.some(v => v.mxPriority !== null && v.mxPriority !== '')) {
    alert('MX 记录必须填写优先级')
    return
  }

  submitting.value = true
  try {
    if (isEdit.value) {
      await axios.put(
        `${API_BASE_URL}/domains/${props.domain.id}/records/${editingRecord.value.id}`,
        formData.value
      )
    } else {
      await axios.post(
        `${API_BASE_URL}/domains/${props.domain.id}/records`,
        formData.value
      )
    }
    await fetchRecords()
    closeModal()
    alert(isEdit.value ? '记录更新成功' : '记录添加成功')
  } catch (error) {
    console.error('Failed to submit:', error)
    const message = error.response?.data?.message || '操作失败'
    alert(message)
  } finally {
    submitting.value = false
  }
}

const disableRecord = async (id) => {
  try {
    await axios.post(`${API_BASE_URL}/domains/${props.domain.id}/records/${id}/disable`)
    await fetchRecords()
  } catch (error) {
    console.error('Failed to disable record:', error)
    alert('禁用记录失败')
  }
}

const enableRecord = async (id) => {
  try {
    await axios.post(`${API_BASE_URL}/domains/${props.domain.id}/records/${id}/enable`)
    await fetchRecords()
  } catch (error) {
    console.error('Failed to enable record:', error)
    alert('启用记录失败')
  }
}

const deleteRecord = async (id) => {
  if (!confirm('确定要删除该解析记录吗？')) {
    return
  }

  deletingRecord.value = id
  try {
    await axios.delete(`${API_BASE_URL}/domains/${props.domain.id}/records/${id}`)
    await fetchRecords()
  } catch (error) {
    console.error('Failed to delete record:', error)
    alert('删除记录失败')
  } finally {
    deletingRecord.value = null
  }
}

// Watch domain change
watch(() => props.domain, (newDomain) => {
  if (newDomain) {
    fetchRecords()
  }
}, { immediate: true })

// Lifecycle
onMounted(() => {
  fetchRecordTypes()
})
</script>

<style scoped>
/* Main container */
/* 主容器 */
.record-management {
  display: flex;
  flex-direction: column;
  height: 100%;
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);
  overflow: hidden;
}

/* Page header */
/* 页面头部 */
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid #e4e7ed;
  background-color: #fff;
  flex-shrink: 0;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.btn-back {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 8px 12px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  background-color: #fff;
  color: #606266;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.btn-back:hover {
  color: #409eff;
  border-color: #c6e2ff;
  background-color: #ecf5ff;
}

.back-icon {
  font-size: 16px;
}

.domain-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.domain-icon {
  font-size: 20px;
}

.domain-name {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

/* Search bar */
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

.filter-group {
  display: flex;
  gap: 12px;
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

/* Table container */
/* 表格容器 */
.table-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  position: relative;
}

.table-scroll-wrapper {
  flex: 1;
  overflow: auto;
}

/* Data table */
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

/* Column widths */
/* 列宽 */
.col-host {
  width: 120px;
}

.col-type {
  width: 80px;
}

.col-value {
  width: 25%;
}

.col-ttl {
  width: 80px;
}

.col-status {
  width: 90px;
}

.col-remark {
  width: 15%;
}

.col-actions {
  width: 100px;
  text-align: center;
}

/* Type tags */
/* 类型标签 */
.type-tag {
  display: inline-block;
  padding: 3px 10px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 600;
}

.type-tag.a {
  background-color: #f0f9eb;
  color: #67c23a;
}

.type-tag.aaaa {
  background-color: #ecf5ff;
  color: #409eff;
}

.type-tag.cname {
  background-color: #fdf6ec;
  color: #e6a23c;
}

.type-tag.mx {
  background-color: #fef0f0;
  color: #f56c6c;
}

.type-tag.txt {
  background-color: #f4f4f5;
  color: #909399;
}

/* Value list */
/* 值列表 */
.value-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.value-item {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-wrap: wrap;
}

.value-text {
  color: #303133;
  font-family: monospace;
  font-size: 13px;
}

.default-badge {
  background-color: #ecf5ff;
  color: #409eff;
  font-size: 11px;
  padding: 1px 6px;
  border-radius: 3px;
}

.cidr-badge {
  background-color: #f4f4f5;
  color: #909399;
  font-size: 11px;
  padding: 1px 6px;
  border-radius: 3px;
}

/* Status badge */
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

/* Action buttons */
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

/* Loading state */
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

/* Empty state */
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

/* Button styles */
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

/* Modal styles */
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
  width: 520px;
  max-width: 95%;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  animation: slideUp 0.3s ease;
}

.modal-large {
  width: 560px;
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

.form-row {
  display: flex;
  gap: 16px;
}

.form-row .form-group {
  flex: 1;
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

.form-group input[type="text"],
.form-group input[type="number"],
.form-group select {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  font-size: 14px;
  transition: all 0.3s ease;
  background-color: #fff;
}

.form-group input:focus,
.form-group select:focus {
  outline: none;
  border-color: #409eff;
  box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.2);
}

.form-group input::placeholder {
  color: #c0c4cc;
}

.form-group select:disabled {
  background-color: #f5f7fa;
  cursor: not-allowed;
}

.required {
  color: #f56c6c;
}

.form-tip {
  font-weight: normal;
  font-size: 12px;
  color: #909399;
  margin-left: 8px;
}

/* Values container */
/* 值容器 */
.values-container {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.value-input-row {
  display: flex;
  align-items: center;
  gap: 12px;
}

.value-input {
  flex: 1;
  min-width: 180px;
}

.value-input-full {
  flex: 1;
}

.cidr-input {
  width: 100px;
  flex-shrink: 0;
}

.mx-priority-input {
  width: 30%;
  flex-shrink: 0;
}

.value-input-mx {
  width: 70%;
  flex-shrink: 0;
}

.checkbox-inline {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  color: #606266;
  white-space: nowrap;
  flex-shrink: 0;
  min-width: 45px;
}

.checkbox-inline input[type="checkbox"] {
  width: 14px;
  height: 14px;
  cursor: pointer;
}

.btn-remove-value {
  width: 28px;
  height: 28px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  background-color: #fff;
  color: #909399;
  font-size: 18px;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
}

.btn-remove-value:hover:not(:disabled) {
  color: #f56c6c;
  border-color: #fbc4c4;
  background-color: #fef0f0;
}

.btn-remove-value:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-add-value {
  padding: 8px 12px;
  border: 1px dashed #dcdfe6;
  border-radius: 4px;
  background-color: #fff;
  color: #409eff;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.3s ease;
  text-align: left;
}

.btn-add-value:hover {
  border-color: #409eff;
  background-color: #ecf5ff;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 24px;
  padding-top: 16px;
  border-top: 1px solid #e4e7ed;
}

/* Mobile styles */
/* 移动端样式 */
.mobile-card-list {
  display: none;
}

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;
    padding: 12px;
  }

  .header-left {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }

  .domain-info {
    width: 100%;
  }

  .header-right {
    width: 100%;
  }

  .header-right .btn {
    width: 100%;
    justify-content: center;
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

  .filter-group {
    flex-direction: column;
    gap: 8px;
  }

  .filter-select {
    width: 100%;
  }

  /* Hide desktop table, show card list */
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

  .mobile-card-title .host-name {
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
    align-items: flex-start;
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
    max-width: 60%;
  }

  .value-item-small {
    display: flex;
    align-items: center;
    gap: 4px;
    margin-bottom: 2px;
  }

  .default-badge-small {
    background-color: #ecf5ff;
    color: #409eff;
    font-size: 10px;
    padding: 1px 4px;
    border-radius: 2px;
  }

  .mobile-card-footer {
    padding: 12px 16px;
    background-color: #f9fafb;
    border-top: 1px solid #ebeef5;
  }

  .mobile-card-actions {
    display: flex;
    gap: 8px;
    flex-wrap: wrap;
  }

  .btn-card-action {
    flex: 1;
    min-width: 80px;
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

  .btn-card-action:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  /* Modal mobile styles */
  /* 弹窗移动端样式 */
  .modal {
    width: 95%;
    max-width: 95%;
    max-height: 85vh;
    margin: 0 10px;
  }

  .modal-large {
    width: 95%;
  }

  .form-row {
    flex-direction: column;
    gap: 0;
  }

  .value-input-row {
    flex-wrap: wrap;
    gap: 8px;
  }

  .checkbox-inline {
    order: 1;
  }

  .cidr-input {
    order: 2;
    width: calc(50% - 30px);
    flex: none;
  }

  .mx-priority-input {
    order: 2;
    width: 30%;
    flex: none;
  }

  .value-input-mx {
    order: 3;
    width: calc(70% - 40px);
    flex: none;
  }

  .btn-remove-value {
    order: 4;
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
