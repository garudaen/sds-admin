<template>
  <div class="records-page">
    <div class="page-header">
      <div class="header-left">
        <el-button @click="goBack">
          <el-icon><ArrowLeft /></el-icon>
          返回域名列表
        </el-button>
        <div v-if="domain" class="domain-info">
          <span class="domain-name">{{ domain.domainName }}</span>
          <el-tag :type="domain.status === 'active' ? 'success' : 'info'" size="small">
            {{ domain.status === 'active' ? '已启用' : '已禁用' }}
          </el-tag>
        </div>
      </div>
      <el-button type="primary" @click="openAddModal">
        <el-icon><Plus /></el-icon>
        添加解析
      </el-button>
    </div>

    <div class="toolbar">
      <el-input
        v-model="searchKeyword"
        placeholder="搜索主机记录..."
        clearable
        class="search-input"
        :prefix-icon="Search"
      />
      <el-select v-model="typeFilter" placeholder="全部类型" clearable class="filter-select">
        <el-option
          v-for="rt in recordTypes"
          :key="rt.id"
          :label="rt.name"
          :value="rt.id"
        />
      </el-select>
      <el-select v-model="statusFilter" placeholder="全部状态" clearable class="filter-select">
        <el-option label="已启用" value="enabled" />
        <el-option label="已禁用" value="disabled" />
      </el-select>
    </div>

    <div class="table-wrap">
      <el-table
        v-loading="loading"
        :data="filteredRecords"
        stripe
        style="width: 100%"
        :empty-text="'暂无解析记录，点击右上角「添加解析」创建'"
      >
        <el-table-column prop="host" label="主机记录" width="120">
          <template #default="{ row }">{{ row.host || '@' }}</template>
        </el-table-column>
        <el-table-column label="记录类型" width="100">
          <template #default="{ row }">
            <el-tag size="small" :type="recordTypeTagType(row.recordType)">{{ row.recordType }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="记录值" min-width="200">
          <template #default="{ row }">
            <div class="value-list">
              <div v-for="v in row.values" :key="v.id" class="value-item">
                <span class="value-text">{{ v.value }}</span>
                <el-tag v-if="v.isDefault" type="primary" size="small">默认</el-tag>
                <el-tag v-if="v.clientCidr" type="info" size="small">{{ v.clientCidr }}</el-tag>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="ttl" label="TTL" width="80">
          <template #default="{ row }">{{ row.ttl }}s</template>
        </el-table-column>
        <el-table-column label="状态" width="90" align="center">
          <template #default="{ row }">
            <el-tag :type="row.disabled ? 'info' : 'success'" size="small">
              {{ row.disabled ? '已禁用' : '已启用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" min-width="120" show-overflow-tooltip>
          <template #default="{ row }">{{ row.remark || '-' }}</template>
        </el-table-column>
        <el-table-column label="操作" width="200" align="center" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="openEditModal(row)">编辑</el-button>
            <template v-if="!row.disabled">
              <el-button type="warning" link size="small" @click="disableRecord(row.id)">禁用</el-button>
            </template>
            <template v-else>
              <el-button type="success" link size="small" @click="enableRecord(row.id)">启用</el-button>
            </template>
            <el-button
              type="danger"
              link
              size="small"
              :loading="deletingRecord === row.id"
              @click="deleteRecord(row.id)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- 添加/编辑解析记录 -->
    <el-dialog
      v-model="showModal"
      :title="isEdit ? '编辑解析记录' : '添加解析记录'"
      width="620px"
      destroy-on-close
      class="record-form-dialog"
      @close="closeModal"
    >
      <el-form
        :model="formData"
        label-width="96px"
        label-position="right"
        class="record-form"
        @submit.prevent="submitForm"
      >
        <el-row :gutter="20">
          <el-col :xs="24" :sm="12">
            <el-form-item label="记录类型" required>
              <el-select
                v-model="formData.recordTypeId"
                placeholder="请选择"
                style="width: 100%"
                :disabled="isEdit"
                @change="onRecordTypeChange"
              >
                <el-option
                  v-for="rt in recordTypes"
                  :key="rt.id"
                  :label="rt.name"
                  :value="rt.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12">
            <el-form-item label="主机记录" required>
              <el-input v-model="formData.host" placeholder="如：www、@、*" clearable />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :xs="24" :sm="12">
            <el-form-item label="TTL（秒）">
              <el-input-number
                v-model="formData.ttl"
                :min="1"
                :max="86400"
                placeholder="默认 300"
                style="width: 100%"
                controls-position="right"
              />
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12">
            <el-form-item label="备注">
              <el-input v-model="formData.remark" placeholder="选填" clearable />
            </el-form-item>
          </el-col>
        </el-row>

        <!-- 记录值区域 -->
        <el-form-item class="record-values-item">
          <template #label>
            <span>记录值</span>
            <el-tooltip v-if="needDefaultValue" placement="top">
              <template #content>A/AAAA/CNAME 必须勾选一个默认值</template>
              <el-icon class="form-tip-icon"><QuestionFilled /></el-icon>
            </el-tooltip>
            <el-tooltip v-else-if="isMXType" placement="top">
              <template #content>可为每条记录设置不同优先级</template>
              <el-icon class="form-tip-icon"><QuestionFilled /></el-icon>
            </el-tooltip>
          </template>
          <div class="record-values-block">
            <div
              v-for="(value, index) in formData.values"
              :key="index"
              class="record-value-row"
            >
              <span class="value-row-label">#{{ index + 1 }}</span>
              <!-- MX：优先级 + 值 -->
              <template v-if="isMXType">
                <el-input-number
                  v-model="value.mxPriority"
                  placeholder="优先级"
                  :min="0"
                  :max="65535"
                  controls-position="right"
                  class="value-field mx-priority"
                />
                <el-input
                  v-model="value.value"
                  placeholder="邮件服务器地址"
                  class="value-field value-input-mx"
                  clearable
                />
              </template>
              <!-- A/AAAA/CNAME：默认 + 值 + CIDR -->
              <template v-else-if="needDefaultValue">
                <el-checkbox
                  v-model="value.isDefault"
                  class="value-default-cb"
                  @change="() => onDefaultChange(index)"
                >
                  默认
                </el-checkbox>
                <el-input
                  v-model="value.value"
                  :placeholder="getValuePlaceholder"
                  class="value-field value-input-main"
                  clearable
                />
                <el-input
                  v-model="value.clientCidr"
                  placeholder="CIDR 可选"
                  class="value-field value-input-cidr"
                  :disabled="value.isDefault"
                  clearable
                />
              </template>
              <!-- 其他类型：仅值 -->
              <template v-else>
                <el-input
                  v-model="value.value"
                  :placeholder="getValuePlaceholder"
                  class="value-field value-input-single"
                  clearable
                />
              </template>
              <el-button
                type="danger"
                link
                circle
                size="small"
                class="value-remove-btn"
                :disabled="formData.values.length <= 1"
                @click="removeValue(index)"
              >
                <el-icon><Delete /></el-icon>
              </el-button>
            </div>
            <el-button type="primary" plain size="small" class="add-value-btn" @click="addValue">
              <el-icon><Plus /></el-icon>
              添加一条记录值
            </el-button>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="closeModal">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submitForm">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useStore } from 'vuex'
import { ArrowLeft, Plus, Search, Delete, QuestionFilled } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '../api/request'

const route = useRoute()
const router = useRouter()
const store = useStore()

const domain = ref(null)
const records = ref([])
const recordTypes = ref([])
const loading = ref(false)
const searchKeyword = ref('')
const typeFilter = ref('')
const statusFilter = ref('')
const deletingRecord = ref(null)

const showModal = ref(false)
const isEdit = ref(false)
const editingRecord = ref(null)
const submitting = ref(false)

const formData = ref({
  recordTypeId: '',
  host: '',
  ttl: 300,
  remark: '',
  values: [{ value: '', isDefault: true, clientCidr: '', mxPriority: null }],
})

const filteredRecords = computed(() => {
  let list = records.value
  if (typeFilter.value) list = list.filter((r) => r.recordTypeId === typeFilter.value)
  if (statusFilter.value === 'enabled') list = list.filter((r) => !r.disabled)
  else if (statusFilter.value === 'disabled') list = list.filter((r) => r.disabled)
  if (searchKeyword.value) {
    const k = searchKeyword.value.toLowerCase()
    list = list.filter(
      (r) =>
        (r.host && r.host.toLowerCase().includes(k)) ||
        r.values.some((v) => v.value && v.value.toLowerCase().includes(k))
    )
  }
  return list
})

const needDefaultValue = computed(() => {
  const rt = recordTypes.value.find((r) => r.id == formData.value.recordTypeId)
  return rt && ['A', 'AAAA', 'CNAME'].includes(rt.name)
})

const isMXType = computed(() => {
  const rt = recordTypes.value.find((r) => r.id == formData.value.recordTypeId)
  return rt && rt.name === 'MX'
})

const getValuePlaceholder = computed(() => {
  const rt = recordTypes.value.find((r) => r.id == formData.value.recordTypeId)
  if (!rt) return '请输入记录值'
  const map = {
    A: 'IPv4 地址，如 192.168.1.1',
    AAAA: 'IPv6 地址',
    CNAME: '目标域名',
    MX: '邮件服务器地址',
  }
  return map[rt.name] || '请输入记录值'
})

function recordTypeTagType(type) {
  const t = (type || '').toLowerCase()
  if (['a', 'aaaa'].includes(t)) return 'success'
  if (t === 'cname') return 'warning'
  if (t === 'mx') return 'danger'
  return 'info'
}

async function loadDomain() {
  const id = route.params.id
  const current = store.getters['app/currentDomain']
  if (current && current.id == id) {
    domain.value = current
    return
  }
  try {
    const res = await request.get(`/domains/${id}`)
    if (res.data?.code === 200 && res.data.data) {
      domain.value = res.data.data
    }
  } catch (_) {}
}

async function fetchRecords() {
  if (!domain.value) return
  loading.value = true
  try {
    const res = await request.get(`/domains/${domain.value.id}/records`)
    if (res.data?.code === 200) {
      records.value = res.data.data?.records || []
    }
  } finally {
    loading.value = false
  }
}

async function fetchRecordTypes() {
  try {
    const res = await request.get('/record-types')
    if (res.data?.code === 200) {
      recordTypes.value = res.data.data || []
    }
  } catch (_) {}
}

function goBack() {
  store.dispatch('app/backToDomains')
  router.push('/domains')
}

function openAddModal() {
  isEdit.value = false
  editingRecord.value = null
  formData.value = {
    recordTypeId: '',
    host: '',
    ttl: 300,
    remark: '',
    values: [{ value: '', isDefault: true, clientCidr: '', mxPriority: null }],
  }
  showModal.value = true
}

function openEditModal(record) {
  isEdit.value = true
  editingRecord.value = record
  formData.value = {
    recordTypeId: record.recordTypeId,
    host: record.host,
    ttl: record.ttl,
    remark: record.remark || '',
    values: record.values.map((v) => ({
      value: v.value,
      isDefault: v.isDefault,
      clientCidr: v.clientCidr || '',
      mxPriority: v.mxPriority,
    })),
  }
  showModal.value = true
}

function closeModal() {
  showModal.value = false
  isEdit.value = false
  editingRecord.value = null
}

function onRecordTypeChange() {
  formData.value.values.forEach((v, i) => {
    v.isDefault = i === 0 && needDefaultValue.value
    if (v.isDefault) v.clientCidr = ''
    if (!isMXType.value) v.mxPriority = null
  })
}

function onDefaultChange(index) {
  if (formData.value.values[index].isDefault) {
    formData.value.values.forEach((v, i) => {
      if (i !== index) v.isDefault = false
    })
    formData.value.values[index].clientCidr = ''
  }
}

function addValue() {
  formData.value.values.push({
    value: '',
    isDefault: false,
    clientCidr: '',
    mxPriority: null,
  })
}

function removeValue(index) {
  if (formData.value.values.length > 1) formData.value.values.splice(index, 1)
}

async function submitForm() {
  if (needDefaultValue.value && !formData.value.values.some((v) => v.isDefault)) {
    ElMessage.warning('A/AAAA/CNAME 记录必须有一个默认值')
    return
  }
  if (isMXType.value && !formData.value.values.some((v) => v.mxPriority != null && v.mxPriority !== '')) {
    ElMessage.warning('MX 记录必须填写优先级')
    return
  }
  submitting.value = true
  try {
    if (isEdit.value) {
      await request.put(
        `/domains/${domain.value.id}/records/${editingRecord.value.id}`,
        formData.value
      )
      ElMessage.success('记录更新成功')
    } else {
      await request.post(`/domains/${domain.value.id}/records`, formData.value)
      ElMessage.success('记录添加成功')
    }
    await fetchRecords()
    closeModal()
  } finally {
    submitting.value = false
  }
}

async function disableRecord(id) {
  await request.post(`/domains/${domain.value.id}/records/${id}/disable`)
  await fetchRecords()
  ElMessage.success('已禁用')
}

async function enableRecord(id) {
  await request.post(`/domains/${domain.value.id}/records/${id}/enable`)
  await fetchRecords()
  ElMessage.success('已启用')
}

async function deleteRecord(id) {
  await ElMessageBox.confirm('确定要删除该解析记录吗？', '提示', { type: 'warning' }).catch(() => {})
  deletingRecord.value = id
  try {
    await request.delete(`/domains/${domain.value.id}/records/${id}`)
    await fetchRecords()
    ElMessage.success('删除成功')
  } finally {
    deletingRecord.value = null
  }
}

watch(
  () => route.params.id,
  async (id) => {
    if (id) {
      await loadDomain()
      if (domain.value) await fetchRecords()
    }
  },
  { immediate: true }
)

onMounted(() => {
  fetchRecordTypes()
  window.addEventListener('layout-add-record', openAddModal)
})

onUnmounted(() => {
  window.removeEventListener('layout-add-record', openAddModal)
})
</script>

<style scoped>
.records-page {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);
  overflow: hidden;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid #e4e7ed;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.domain-name {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.domain-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.toolbar {
  display: flex;
  gap: 12px;
  padding: 16px 20px;
  border-bottom: 1px solid #e4e7ed;
}

.search-input {
  max-width: 280px;
}

.filter-select {
  width: 120px;
}

.table-wrap {
  flex: 1;
  overflow: auto;
  padding: 0 20px;
}

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
  font-family: monospace;
  font-size: 13px;
}

/* 解析记录弹窗 */
.record-form-dialog :deep(.el-dialog__body) {
  padding-top: 8px;
  max-height: 60vh;
  overflow-y: auto;
}

.record-form {
  padding-right: 8px;
}

.record-form .el-form-item {
  margin-bottom: 18px;
}

.record-values-item :deep(.el-form-item__label) {
  align-items: flex-start;
  padding-top: 8px;
}

.form-tip-icon {
  margin-left: 4px;
  color: var(--el-color-info-light-5);
  font-size: 14px;
  vertical-align: middle;
  cursor: help;
}

.record-values-block {
  width: 100%;
  padding: 12px 16px;
  background: var(--el-fill-color-light);
  border-radius: 8px;
  border: 1px solid var(--el-border-color-lighter);
}

.record-value-row {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 10px;
}

.record-value-row:last-of-type {
  margin-bottom: 0;
}

.value-row-label {
  flex-shrink: 0;
  width: 28px;
  font-size: 12px;
  color: var(--el-text-color-secondary);
  text-align: right;
}

.value-field {
  flex: 1;
  min-width: 0;
}

.value-field.mx-priority {
  flex: 0 0 100px;
}

.value-field.value-input-mx {
  flex: 1;
  min-width: 160px;
}

.value-default-cb {
  flex-shrink: 0;
  margin-right: 4px;
}

.value-field.value-input-main {
  flex: 1;
  min-width: 140px;
}

.value-field.value-input-cidr {
  flex: 0 0 110px;
}

.value-field.value-input-single {
  flex: 1;
  min-width: 200px;
}

.value-remove-btn {
  flex-shrink: 0;
}

.add-value-btn {
  margin-top: 10px;
}

/* 小屏下记录值行改为纵向 */
@media (max-width: 520px) {
  .record-value-row {
    flex-wrap: wrap;
  }

  .value-row-label {
    width: 100%;
    text-align: left;
  }

  .value-field.mx-priority,
  .value-field.value-input-cidr {
    flex: 1 1 100%;
  }

  .value-remove-btn {
    margin-left: auto;
  }
}

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;
  }

  .toolbar {
    flex-direction: column;
  }

  .search-input,
  .filter-select {
    max-width: none;
    width: 100%;
  }
}
</style>
